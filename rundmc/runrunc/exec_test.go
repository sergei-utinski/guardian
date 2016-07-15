package runrunc_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/guardian/rundmc/goci"
	"code.cloudfoundry.org/guardian/rundmc/runrunc"
	fakes "code.cloudfoundry.org/guardian/rundmc/runrunc/runruncfakes"
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagertest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opencontainers/runc/libcontainer/user"
	"github.com/opencontainers/runtime-spec/specs-go"
)

var _ = Describe("Execer", func() {
	var (
		logger *lagertest.TestLogger

		execPreparer *fakes.FakeExecPreparer
		execRunner   *fakes.FakeExecRunner
		execer       *runrunc.Execer
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("test")
		execRunner = new(fakes.FakeExecRunner)
		execPreparer = new(fakes.FakeExecPreparer)

		execer = runrunc.NewExecer(
			execPreparer,
			execRunner,
		)
	})

	It("runs the execRunner with the prepared process spec", func() {
		execPreparer.PrepareStub = func(log lager.Logger, bundlePath string, spec garden.ProcessSpec) (*runrunc.PreparedSpec, error) {
			return &runrunc.PreparedSpec{
				Process: specs.Process{
					Args: []string{spec.Path, bundlePath},
				},
				HostUID: 10,
				HostGID: 10,
			}, nil
		}

		execer.Exec(logger, "some-bundle-path", "some-id", garden.ProcessSpec{
			Path: "potato",
		}, garden.ProcessIO{})

		Expect(execRunner.RunCallCount()).To(Equal(1))
		_, spec, processesPath, id, _, _ := execRunner.RunArgsForCall(0)
		Expect(spec.Args).To(ConsistOf("potato", "some-bundle-path"))
		Expect(processesPath).To(Equal("some-bundle-path/processes"))
		Expect(id).To(Equal("some-id"))
	})
})

var _ = Describe("ExecPreparer", func() {
	var (
		spec         *runrunc.PreparedSpec
		bundleLoader *fakes.FakeBundleLoader
		users        *fakes.FakeUserLookupper
		mkdirer      *fakes.FakeMkdirer
		bundlePath   string
		logger       lager.Logger

		preparer runrunc.ExecPreparer
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("test")
		bundleLoader = new(fakes.FakeBundleLoader)
		users = new(fakes.FakeUserLookupper)
		mkdirer = new(fakes.FakeMkdirer)

		var err error
		bundlePath, err = ioutil.TempDir("", "bundle")
		Expect(err).NotTo(HaveOccurred())

		bundleLoader.LoadStub = func(path string) (goci.Bndl, error) {
			bndl := goci.Bndl{}
			return bndl, nil
		}

		users.LookupReturns(&user.ExecUser{}, nil)

		Expect(ioutil.WriteFile(filepath.Join(bundlePath, "pidfile"), []byte("999"), 0644)).To(Succeed())
		preparer = runrunc.NewExecPreparer(bundleLoader, users, mkdirer, []string{"foo", "bar", "brains"})
	})

	It("passes a process.json with the correct path and args", func() {
		spec, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{Path: "to enlightenment", Args: []string{"infinity", "and beyond"}})
		Expect(err).NotTo(HaveOccurred())

		Expect(spec.Args).To(Equal([]string{"to enlightenment", "infinity", "and beyond"}))
	})

	It("returns the HostUID and HostGID in the returned spec", func() {
		users.LookupReturns(&user.ExecUser{Uid: 234, Gid: 567}, nil)

		spec, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{Path: "to enlightenment", Args: []string{}})
		Expect(err).NotTo(HaveOccurred())

		Expect(spec.HostUID).To(BeEquivalentTo(234))
		Expect(spec.HostGID).To(BeEquivalentTo(567))
	})

	It("sets the rlimits correctly", func() {
		ptr := func(n uint64) *uint64 { return &n }
		spec, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{
			Limits: garden.ResourceLimits{
				As:         ptr(12),
				Core:       ptr(24),
				Cpu:        ptr(36),
				Data:       ptr(99),
				Fsize:      ptr(101),
				Locks:      ptr(111),
				Memlock:    ptr(987),
				Msgqueue:   ptr(777),
				Nice:       ptr(111),
				Nofile:     ptr(222),
				Nproc:      ptr(1234),
				Rss:        ptr(888),
				Rtprio:     ptr(254),
				Sigpending: ptr(101),
				Stack:      ptr(44),
			},
		})
		Expect(err).ToNot(HaveOccurred())

		Expect(spec.Process.Rlimits).To(ConsistOf(
			specs.Rlimit{Type: "RLIMIT_AS", Hard: 12, Soft: 12},
			specs.Rlimit{Type: "RLIMIT_CORE", Hard: 24, Soft: 24},
			specs.Rlimit{Type: "RLIMIT_CPU", Hard: 36, Soft: 36},
			specs.Rlimit{Type: "RLIMIT_DATA", Hard: 99, Soft: 99},
			specs.Rlimit{Type: "RLIMIT_FSIZE", Hard: 101, Soft: 101},
			specs.Rlimit{Type: "RLIMIT_LOCKS", Hard: 111, Soft: 111},
			specs.Rlimit{Type: "RLIMIT_MEMLOCK", Hard: 987, Soft: 987},
			specs.Rlimit{Type: "RLIMIT_MSGQUEUE", Hard: 777, Soft: 777},
			specs.Rlimit{Type: "RLIMIT_NICE", Hard: 111, Soft: 111},
			specs.Rlimit{Type: "RLIMIT_NOFILE", Hard: 222, Soft: 222},
			specs.Rlimit{Type: "RLIMIT_NPROC", Hard: 1234, Soft: 1234},
			specs.Rlimit{Type: "RLIMIT_RSS", Hard: 888, Soft: 888},
			specs.Rlimit{Type: "RLIMIT_RTPRIO", Hard: 254, Soft: 254},
			specs.Rlimit{Type: "RLIMIT_SIGPENDING", Hard: 101, Soft: 101},
			specs.Rlimit{Type: "RLIMIT_STACK", Hard: 44, Soft: 44},
		))
	})

	It("sets Terminal to true iff a TTY is configured", func() {
		spec, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{
			TTY: &garden.TTYSpec{
				WindowSize: &garden.WindowSize{
					Columns: 80,
					Rows:    24,
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())

		Expect(spec.Process.Terminal).To(BeTrue())

		spec, err = preparer.Prepare(logger, bundlePath, garden.ProcessSpec{
			TTY: nil,
		})
		Expect(err).ToNot(HaveOccurred())

		Expect(spec.Process.Terminal).To(BeFalse())
	})

	Describe("passing the correct uid and gid", func() {
		Context("when the bundle can be loaded", func() {
			BeforeEach(func() {
				users.LookupReturns(&user.ExecUser{Uid: 9, Gid: 7}, nil)

				var err error
				spec, err = preparer.Prepare(logger, bundlePath, garden.ProcessSpec{User: "spiderman"})
				Expect(err).ToNot(HaveOccurred())
			})

			It("looks up the user and group IDs of the user in the right rootfs", func() {
				Expect(users.LookupCallCount()).To(Equal(1))
				actualRootfsPath, actualUserName := users.LookupArgsForCall(0)
				Expect(actualRootfsPath).To(Equal("/proc/999/root"))
				Expect(actualUserName).To(Equal("spiderman"))
			})

			It("passes a process.json with the correct user and group ids", func() {
				Expect(spec.Process.User).To(Equal(specs.User{UID: 9, GID: 7}))
			})
		})

		Context("when the bundle can't be loaded", func() {
			BeforeEach(func() {
				bundleLoader.LoadReturns(goci.Bndl{}, errors.New("whoa! Hold them horses!"))
			})

			It("fails", func() {
				_, err := preparer.Prepare(logger, bundlePath,
					garden.ProcessSpec{User: "spiderman"})
				Expect(err).To(MatchError(ContainSubstring("Hold them horses")))
			})
		})

		Context("when User Lookup returns an error", func() {
			It("passes a process.json with the correct user and group ids", func() {
				users.LookupReturns(&user.ExecUser{Uid: 0, Gid: 0}, errors.New("bang"))

				_, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{User: "spiderman"})
				Expect(err).To(MatchError(ContainSubstring("bang")))
			})
		})

		Context("when the pidfile can't be read", func() {
			It("returns an appropriate error", func() {
				Expect(os.Remove(filepath.Join(bundlePath, "pidfile"))).To(Succeed())

				_, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{User: "spiderman"})
				Expect(err).To(MatchError(ContainSubstring("pidfile")))
			})
		})
	})

	Context("when the user is specified in the process spec", func() {
		Context("when the environment does not contain a USER", func() {
			It("appends a default user", func() {
				spec, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{
					User: "spiderman",
					Env:  []string{"a=1", "b=3", "c=4", "PATH=a", "HOME=/spidermanhome"},
				})
				Expect(err).NotTo(HaveOccurred())

				Expect(spec.Process.Env).To(ConsistOf("a=1", "b=3", "c=4", "PATH=a", "USER=spiderman", "HOME=/spidermanhome"))
			})
		})

		Context("when the environment does contain a USER", func() {
			It("appends a default user", func() {
				spec, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{
					User: "spiderman",
					Env:  []string{"a=1", "b=3", "c=4", "PATH=a", "USER=superman"},
				})
				Expect(err).NotTo(HaveOccurred())

				Expect(spec.Process.Env).To(Equal([]string{"a=1", "b=3", "c=4", "PATH=a", "USER=superman"}))
			})
		})
	})

	Context("when the user is not specified in the process spec", func() {
		Context("when the environment does not contain a USER", func() {
			It("passes the environment variables", func() {
				spec, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{
					Env: []string{"a=1", "b=3", "c=4", "PATH=a"},
				})
				Expect(err).NotTo(HaveOccurred())

				Expect(spec.Process.Env).To(Equal([]string{"a=1", "b=3", "c=4", "PATH=a", "USER=root"}))
			})
		})

		Context("when the environment already contains a USER", func() {
			It("passes the environment variables", func() {
				spec, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{
					Env: []string{"a=1", "b=3", "c=4", "PATH=a", "USER=yo"},
				})
				Expect(err).NotTo(HaveOccurred())

				Expect(spec.Process.Env).To(Equal([]string{"a=1", "b=3", "c=4", "PATH=a", "USER=yo"}))
			})
		})
	})

	Context("when the environment already contains a PATH", func() {
		It("passes the environment variables", func() {
			spec, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{
				Env: []string{"a=1", "b=3", "c=4", "PATH=a"},
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(spec.Process.Env).To(Equal([]string{"a=1", "b=3", "c=4", "PATH=a", "USER=root"}))
		})
	})

	Context("when the environment does not already contain a PATH", func() {
		It("appends a default PATH for the root user", func() {
			users.LookupReturns(&user.ExecUser{Uid: 0, Gid: 0}, nil)
			spec, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{
				Env:  []string{"a=1", "b=3", "c=4"},
				User: "root",
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(spec.Process.Env).To(Equal([]string{"a=1", "b=3", "c=4",
				"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin", "USER=root"}))
		})

		It("appends a default PATH for non-root users", func() {
			users.LookupReturns(&user.ExecUser{Uid: 1000, Gid: 1000}, nil)
			spec, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{
				Env:  []string{"a=1", "b=3", "c=4"},
				User: "alice",
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(spec.Process.Env).To(Equal([]string{"a=1", "b=3", "c=4",
				"PATH=/usr/local/bin:/usr/bin:/bin", "USER=alice"}))
		})
	})

	Context("when the container has environment variables", func() {
		var (
			processEnv   []string
			containerEnv []string
			bndl         goci.Bndl

			spec *runrunc.PreparedSpec
		)

		BeforeEach(func() {
			containerEnv = []string{"ENV_CONTAINER_NAME=garden"}
			processEnv = []string{"ENV_PROCESS_ID=1"}
		})

		JustBeforeEach(func() {
			bndl = goci.Bndl{}
			bndl.Spec.Process.Env = containerEnv
			bundleLoader.LoadReturns(bndl, nil)

			var err error
			spec, err = preparer.Prepare(logger, bundlePath, garden.ProcessSpec{
				Env: processEnv,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("appends the process vars into container vars", func() {
			envWContainer := make([]string, len(spec.Process.Env))
			copy(envWContainer, spec.Process.Env)

			bndl.Spec.Process.Env = []string{}
			bundleLoader.LoadReturns(bndl, nil)

			spec, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{
				Env: processEnv,
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(envWContainer).To(Equal(append(containerEnv, spec.Process.Env...)))
		})

		Context("and the container environment contains PATH", func() {
			BeforeEach(func() {
				containerEnv = append(containerEnv, "PATH=/test")
			})

			It("should not apply the default PATH", func() {
				Expect(spec.Process.Env).To(Equal([]string{
					"ENV_CONTAINER_NAME=garden",
					"PATH=/test",
					"ENV_PROCESS_ID=1",
					"USER=root",
				}))
			})
		})
	})

	Context("when the container has capabilities", func() {
		BeforeEach(func() {
			bndl := goci.Bndl{}
			bndl.Spec.Process.Capabilities = []string{"foo", "bar", "baz"}
			bundleLoader.LoadReturns(bndl, nil)
		})

		Context("and the user is root", func() {
			It("passes them on to the process", func() {
				spec, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{})
				Expect(err).NotTo(HaveOccurred())
				Expect(spec.Process.Capabilities).To(Equal([]string{"foo", "bar", "baz"}))
			})
		})

		Context("and the user is not root", func() {
			It("removes any caps not in nonRootMaxCaps list", func() {
				users.LookupReturns(&user.ExecUser{Uid: 1234, Gid: 0}, nil)
				spec, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{})

				Expect(err).NotTo(HaveOccurred())
				Expect(spec.Process.Capabilities).To(Equal([]string{"foo", "bar"}))
			})
		})
	})

	Describe("working directory", func() {
		Context("when the working directory is specified", func() {
			It("passes the correct cwd to the spec", func() {
				spec, err := preparer.Prepare(
					logger, bundlePath,
					garden.ProcessSpec{Dir: "/home/dir"},
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(spec.Process.Cwd).To(Equal("/home/dir"))
			})

			Describe("Creating the working directory", func() {
				JustBeforeEach(func() {
					users.LookupReturns(&user.ExecUser{Uid: 1012, Gid: 1013}, nil)

					_, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{
						Dir: "/path/to/banana/dir",
					})
					Expect(err).NotTo(HaveOccurred())
				})

				Context("when the container is privileged", func() {
					It("creates the working directory", func() {
						Expect(mkdirer.MkdirAsCallCount()).To(Equal(1))
						rootfs, uid, gid, mode, recreate, dirs := mkdirer.MkdirAsArgsForCall(0)
						Expect(rootfs).To(Equal("/proc/999/root"))
						Expect(dirs).To(ConsistOf("/path/to/banana/dir"))
						Expect(mode).To(BeNumerically("==", 0755))
						Expect(recreate).To(BeFalse())
						Expect(uid).To(BeEquivalentTo(1012))
						Expect(gid).To(BeEquivalentTo(1013))
					})
				})

				Context("when the container is unprivileged", func() {
					BeforeEach(func() {
						bundleLoader.LoadStub = func(path string) (goci.Bndl, error) {
							bndl := goci.Bndl{}
							bndl.Spec.Linux.UIDMappings = []specs.IDMapping{{
								HostID:      1712,
								ContainerID: 1012,
								Size:        1,
							}}
							bndl.Spec.Linux.GIDMappings = []specs.IDMapping{{
								HostID:      1713,
								ContainerID: 1013,
								Size:        1,
							}}
							return bndl, nil
						}
					})

					It("creates the working directory as the mapped user", func() {
						Expect(mkdirer.MkdirAsCallCount()).To(Equal(1))
						rootfs, uid, gid, mode, recreate, dirs := mkdirer.MkdirAsArgsForCall(0)
						Expect(rootfs).To(Equal("/proc/999/root"))
						Expect(dirs).To(ConsistOf("/path/to/banana/dir"))
						Expect(mode).To(BeEquivalentTo(0755))
						Expect(recreate).To(BeFalse())
						Expect(uid).To(BeEquivalentTo(1712))
						Expect(gid).To(BeEquivalentTo(1713))
					})
				})
			})
		})

		Context("when the working directory is not specified", func() {
			It("defaults to the user's HOME directory", func() {
				users.LookupReturns(&user.ExecUser{Home: "/the/home/dir"}, nil)

				spec, err := preparer.Prepare(
					logger, bundlePath,
					garden.ProcessSpec{Dir: ""},
				)
				Expect(err).NotTo(HaveOccurred())

				Expect(spec.Process.Cwd).To(Equal("/the/home/dir"))
			})

			It("creates the directory", func() {
				users.LookupReturns(&user.ExecUser{Uid: 1012, Gid: 1013, Home: "/some/dir"}, nil)

				_, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{})
				Expect(err).NotTo(HaveOccurred())

				Expect(mkdirer.MkdirAsCallCount()).To(Equal(1))
				_, _, _, _, _, dirs := mkdirer.MkdirAsArgsForCall(0)
				Expect(dirs).To(ConsistOf("/some/dir"))
			})
		})

		Context("when the working directory creation fails", func() {
			It("returns an error", func() {
				mkdirer.MkdirAsReturns(errors.New("BOOOOOM"))
				_, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{})
				Expect(err).To(MatchError(ContainSubstring("create working directory: BOOOOOM")))
			})
		})
	})

	Context("when an ApparmorProfile is defined in the base process", func() {
		BeforeEach(func() {
			bundleLoader.LoadStub = func(path string) (goci.Bndl, error) {
				bndl := goci.Bndl{}
				bndl = bndl.WithProcess(specs.Process{
					ApparmorProfile: "default-profile",
				})
				return bndl, nil
			}
		})

		It("should pass it to the process spec", func() {
			spec, err := preparer.Prepare(logger, bundlePath, garden.ProcessSpec{})
			Expect(err).ToNot(HaveOccurred())

			Expect(spec.Process.ApparmorProfile).To(Equal("default-profile"))
		})
	})
})

var _ = Describe("WaitWatcher", func() {
	It("calls Wait only once process.Wait returns", func() {
		waiter := new(fakes.FakeWaiter)
		waitReturns := make(chan struct{})
		waiter.WaitStub = func() (int, error) {
			<-waitReturns
			return 0, nil
		}

		runner := new(fakes.FakeRunner)

		watcher := runrunc.Watcher{}
		go watcher.OnExit(lagertest.NewTestLogger("test"), waiter, runner)

		Consistently(runner.RunCallCount).ShouldNot(Equal(1))
		close(waitReturns)
		Eventually(runner.RunCallCount).Should(Equal(1))
	})
})

var _ = Describe("RemoveFiles", func() {
	It("removes all the paths", func() {
		a := tmpFile("testremovefiles")
		b := tmpFile("testremovefiles")
		runrunc.RemoveFiles([]string{a, b}).Run(lagertest.NewTestLogger("test"))
		Expect(a).NotTo(BeAnExistingFile())
		Expect(b).NotTo(BeAnExistingFile())
	})
})

func tmpFile(name string) string {
	tmp, err := ioutil.TempFile("", name)
	Expect(err).NotTo(HaveOccurred())
	Expect(tmp.Close()).To(Succeed())
	return tmp.Name()
}
