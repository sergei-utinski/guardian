package peas

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/guardian/gardener"
	"code.cloudfoundry.org/guardian/rundmc/depot"
	"code.cloudfoundry.org/guardian/rundmc/runrunc"
	"code.cloudfoundry.org/lager"
	multierror "github.com/hashicorp/go-multierror"
	uuid "github.com/nu7hatch/gouuid"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	errorwrapper "github.com/pkg/errors"
)

var RootfsPath = filepath.Join(os.TempDir(), "pea-empty-rootfs")

//go:generate counterfeiter . Volumizer
type Volumizer interface {
	Create(log lager.Logger, spec garden.ContainerSpec) (specs.Spec, error)
	Destroy(log lager.Logger, handle string) error
}

//go:generate counterfeiter . PidGetter
type PidGetter interface {
	Pid(pidFilePath string) (int, error)
}

//go:generate counterfeiter . PrivilegedGetter
type PrivilegedGetter interface {
	Privileged(bundlePath string) (bool, error)
}

//go:generate counterfeiter . RuncDeleter
type RuncDeleter interface {
	Delete(log lager.Logger, force bool, handle string) error
}

type PeaCreator struct {
	Volumizer              Volumizer
	PidGetter              PidGetter
	PrivilegedGetter       PrivilegedGetter
	BindMountSourceCreator depot.BindMountSourceCreator
	BundleGenerator        depot.BundleGenerator
	BundleSaver            depot.BundleSaver
	ProcessBuilder         runrunc.ProcessBuilder
	ExecRunner             runrunc.ExecRunner
	RuncDeleter            RuncDeleter
}

func (p *PeaCreator) CreatePea(log lager.Logger, spec garden.ProcessSpec, procIO garden.ProcessIO, sandboxHandle, sandboxBundlePath string) (garden.Process, error) {
	errs := func(action string, err error) (garden.Process, error) {
		wrappedErr := errorwrapper.Wrap(err, action)
		log.Error(action, wrappedErr)
		return nil, wrappedErr
	}

	log = log.Session("create-pea", lager.Data{})

	processID, err := generateProcessID(spec.ID)
	if err != nil {
		return nil, err
	}

	log.Info("creating", lager.Data{"process_id": processID})
	defer log.Info("done")

	peaBundlePath := filepath.Join(sandboxBundlePath, "processes", processID)
	if mkdirErr := os.MkdirAll(peaBundlePath, 0700); mkdirErr != nil {
		return nil, err
	}

	privileged, err := p.PrivilegedGetter.Privileged(sandboxBundlePath)
	if err != nil {
		return errs("determining-privileged", err)
	}

	defaultBindMounts, err := p.BindMountSourceCreator.Create(sandboxBundlePath, !privileged)
	if err != nil {
		return errs("creating-bind-mount-sources", err)
	}

	if spec.Dir == "" {
		spec.Dir = "/"
	}
	uid, gid, err := parseUser(spec.User)
	if err != nil {
		return errs("parse-user", err)
	}

	linuxNamespaces, err := p.linuxNamespaces(sandboxBundlePath, privileged)
	if err != nil {
		return errs("determining-namespaces", err)
	}

	runtimeSpec, err := p.Volumizer.Create(log, garden.ContainerSpec{
		Handle:     processID,
		Image:      spec.Image,
		Privileged: privileged,
	})
	if err != nil {
		return errs("creating-volume", err)
	}

	if runtimeSpec.Windows == nil {
		runtimeSpec.Windows = new(specs.Windows)
	}

	runtimeSpec.Windows.Network = &specs.WindowsNetwork{
		NetworkSharedContainerName: sandboxHandle,
	}

	cgroupPath := sandboxHandle
	if spec.OverrideContainerLimits != nil {
		cgroupPath = processID
	}

	bndl, genErr := p.BundleGenerator.Generate(gardener.DesiredContainerSpec{
		Handle:     processID,
		BaseConfig: runtimeSpec,
		CgroupPath: cgroupPath,
		Namespaces: linuxNamespaces,
		BindMounts: append(spec.BindMounts, defaultBindMounts...),
		Privileged: privileged,
	}, sandboxBundlePath)
	if genErr != nil {
		destroyErr := p.Volumizer.Destroy(log, processID)
		return errs("generating-bundle", multierror.Append(genErr, destroyErr))
	}

	preparedProcess := p.ProcessBuilder.BuildProcess(bndl, runrunc.ProcessSpec{
		ProcessSpec:  spec,
		ContainerUID: uid,
		ContainerGID: gid,
	})

	bndl = bndl.WithProcess(preparedProcess.Process)
	if saveErr := p.BundleSaver.Save(bndl, peaBundlePath); saveErr != nil {
		destroyErr := p.Volumizer.Destroy(log, processID)
		return errs("saving-bundle", multierror.Append(saveErr, destroyErr))
	}

	extraCleanup := func() error {
		var result *multierror.Error
		err := p.RuncDeleter.Delete(log, true, processID)
		if err != nil {
			result = multierror.Append(result, err)
		}
		err = p.Volumizer.Destroy(log, processID)
		if err != nil {
			result = multierror.Append(result, err)
		}

		return result.ErrorOrNil()
	}
	proc, runErr := p.ExecRunner.Run(
		log, processID, peaBundlePath, sandboxHandle, sandboxBundlePath,
		preparedProcess.ContainerRootHostUID, preparedProcess.ContainerRootHostGID,
		procIO, preparedProcess.Process.Terminal, nil, extraCleanup,
	)
	if runErr != nil {
		destroyErr := p.Volumizer.Destroy(log, processID)
		return nil, multierror.Append(runErr, destroyErr)
	}

	return proc, nil
}

func (p *PeaCreator) linuxNamespaces(sandboxBundlePath string, privileged bool) (map[string]string, error) {
	originalCtrInitPid, err := p.PidGetter.Pid(filepath.Join(sandboxBundlePath, "pidfile"))
	if err != nil {
		return nil, errorwrapper.Wrap(err, "reading-ctr-pid")
	}

	linuxNamespaces := map[string]string{}
	linuxNamespaces["mount"] = ""
	linuxNamespaces["network"] = fmt.Sprintf("/proc/%d/ns/net", originalCtrInitPid)
	linuxNamespaces["ipc"] = fmt.Sprintf("/proc/%d/ns/ipc", originalCtrInitPid)
	linuxNamespaces["pid"] = fmt.Sprintf("/proc/%d/ns/pid", originalCtrInitPid)
	linuxNamespaces["uts"] = fmt.Sprintf("/proc/%d/ns/uts", originalCtrInitPid)

	if !privileged {
		linuxNamespaces["user"] = fmt.Sprintf("/proc/%d/ns/user", originalCtrInitPid)
	}

	return linuxNamespaces, nil
}

func parseUser(uidgid string) (int, int, error) {
	if uidgid == "" {
		return 0, 0, nil
	}

	errs := func() (int, int, error) {
		return 0, 0, fmt.Errorf("'%s' is not a valid uid:gid", uidgid)
	}

	uidGidComponents := strings.Split(uidgid, ":")
	if len(uidGidComponents) != 2 {
		return errs()
	}
	uid, err := strconv.Atoi(uidGidComponents[0])
	if err != nil {
		return errs()
	}
	gid, err := strconv.Atoi(uidGidComponents[1])
	if err != nil {
		return errs()
	}
	return uid, gid, nil
}

func generateProcessID(existingID string) (string, error) {
	if existingID != "" {
		return existingID, nil
	}
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
