package server

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/darcys22/godbledger-web/pkg/api"
	"github.com/darcys22/godbledger-web/pkg/registry"
	"github.com/darcys22/godbledger-web/pkg/setting"

	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"golang.org/x/sync/errgroup"
)

var log = logrus.WithField("prefix", "server")

// Config contains parameters for the New function.
type Config struct {
	ConfigFile  string
	HomePath    string
	PidFile     string
	Version     string
	Commit      string
	BuildBranch string
	Listener    net.Listener
}

// New returns a new instance of Server.
func New(cfg Config) (*Server, error) {
	rootCtx, shutdownFn := context.WithCancel(context.Background())
	childRoutines, childCtx := errgroup.WithContext(rootCtx)

	s := &Server{
		context:       childCtx,
		shutdownFn:    shutdownFn,
		childRoutines: childRoutines,
		cfg:           setting.NewCfg(),

		configFile:  cfg.ConfigFile,
		homePath:    cfg.HomePath,
		pidFile:     cfg.PidFile,
		version:     cfg.Version,
		commit:      cfg.Commit,
		buildBranch: cfg.BuildBranch,
	}
	if cfg.Listener != nil {
		if err := s.init(&cfg); err != nil {
			return nil, err
		}
	}

	return s, nil
}

// Server is responsible for managing the lifecycle of services.
type Server struct {
	context            context.Context
	shutdownFn         context.CancelFunc
	childRoutines      *errgroup.Group
	cfg                *setting.Cfg
	shutdownReason     string
	shutdownInProgress bool
	isInitialized      bool
	mtx                sync.Mutex

	configFile  string
	homePath    string
	pidFile     string
	version     string
	commit      string
	buildBranch string

	httpSrv *http.Server
}

// init initializes the server and its services.
func (s *Server) init(cfg *Config) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.isInitialized {
		return nil
	}
	s.isInitialized = true

	s.loadConfiguration()
	s.writePIDFile()

	//login.Init()

	return nil
}

// Run initializes and starts services. This will block until all services have
// exited. To initiate shutdown, call the Shutdown method in another goroutine.
func (s *Server) Run() (err error) {

	if err = s.init(nil); err != nil {
		return
	}

	services := registry.GetServices()

	// Start background services.
	for _, svc := range services {
		service, ok := svc.Instance.(registry.BackgroundService)
		if !ok {
			continue
		}

		if registry.IsDisabled(svc.Instance) {
			continue
		}

		// Variable is needed for accessing loop variable in callback
		descriptor := svc
		s.childRoutines.Go(func() error {
			// Don't start new services when server is shutting down.
			if s.shutdownInProgress {
				return nil
			}

			err := service.Run(s.context)
			// Mark that we are in shutdown mode
			// So no more services are started
			s.shutdownInProgress = true
			if err != nil {
				if err != context.Canceled {
					// Server has crashed.
					log.Error("Stopped "+descriptor.Name, "reason", err)
				} else {
					log.Debug("Stopped "+descriptor.Name, "reason", err)
				}

				return err
			}

			return nil
		})
	}

	defer func() {
		log.Debug("Waiting on services...")
		if waitErr := s.childRoutines.Wait(); waitErr != nil && !errors.Is(waitErr, context.Canceled) {
			log.Error("A service failed", "err", waitErr)
			if err == nil {
				err = waitErr
			}
		}
	}()

	s.notifySystemd("READY=1")

	m := newGin()
	api.Register(m)

	listenAddr := fmt.Sprintf("%s:%s", setting.HttpAddr, setting.HttpPort)
	log.Infof("Listen: %v://%s%s", setting.Protocol, listenAddr, setting.AppSubUrl)

	s.httpSrv = &http.Server{
		Addr:    listenAddr,
		Handler: m,
	}

	var wg sync.WaitGroup
	wg.Add(1)

	// handle http shutdown on server context done
	go func() {
		defer wg.Done()

		<-s.context.Done()
		if err := s.httpSrv.Shutdown(context.Background()); err != nil {
			log.WithField("err", err).Error("Failed to shutdown server", "error", err)
		}
	}()

	switch setting.Protocol {
	case setting.HTTP:
		if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	case setting.HTTPS:
		if err := s.httpSrv.ListenAndServeTLS(setting.CertFile, setting.KeyFile); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	default:
		log.Fatalf("Invalid protocol: %s", setting.Protocol)
	}

	return nil
}

func (s *Server) Shutdown(reason string) {
	log.WithField("reason", reason).Info("Shutdown started")
	s.shutdownReason = reason
	s.shutdownInProgress = true

	// call cancel func on root context
	s.shutdownFn()

	// wait for child routines

	if err := s.childRoutines.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		log.WithField("err", err).Error("Failed waiting for services to shutdown", "error", err)
	}
}

// ExitCode returns an exit code for a given error.
func (s *Server) ExitCode(reason error) int {
	code := 1

	if reason == context.Canceled && s.shutdownReason != "" {
		reason = fmt.Errorf(s.shutdownReason)
		code = 0
	}

	log.Error("Server shutdown", "reason", reason)

	return code
}

// loadConfiguration loads settings and configuration from config files.
func (s *Server) loadConfiguration() {
	args := &setting.CommandLineArgs{
		Config:   s.configFile,
		HomePath: s.homePath,
		Args:     flag.Args(),
	}

	if err := s.cfg.Load(args); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start godbledger-web. error: %s\n", err.Error())
		os.Exit(1)
	}

	log.WithFields(logrus.Fields{
		"version":  s.version,
		"commit":   s.commit,
		"branch":   s.buildBranch,
		"compiled": time.Unix(setting.BuildStamp, 0),
	}).Infof("Starting %s", setting.ApplicationName)
}

// notifySystemd sends state notifications to systemd.
func (s *Server) notifySystemd(state string) {
	notifySocket := os.Getenv("NOTIFY_SOCKET")
	if notifySocket == "" {
		log.Debug(
			"NOTIFY_SOCKET environment variable empty or unset, can't send systemd notification")
		return
	}

	socketAddr := &net.UnixAddr{
		Name: notifySocket,
		Net:  "unixgram",
	}
	conn, err := net.DialUnix(socketAddr.Net, nil, socketAddr)
	if err != nil {
		log.Warn("Failed to connect to systemd", "err", err, "socket", notifySocket)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte(state))
	if err != nil {
		log.Warn("Failed to write notification to systemd", "err", err)
	}
}

func newGin() *gin.Engine {
	//gin.Env = setting.Env
	//macaron.Env = setting.Env

	m := gin.Default()
	//m.Use(middleware.Logger())
	m.Use(gin.Recovery())
	if setting.EnableGzip {
		m.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	mapStatic(m, "", "public")
	mapStatic(m, "app", "app")
	mapStatic(m, "css", "css")
	mapStatic(m, "img", "img")

	m.LoadHTMLGlob(path.Join(setting.StaticRootPath, "views/*.html"))

	return m
}

func mapStatic(m *gin.Engine, dir string, prefix string) {
	headers := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Writer.Header().Set("Cache-Control", "public, max-age=3600")
			c.Next()
		}
	}

	if setting.Env == setting.DEV {
		headers = func() gin.HandlerFunc {
			return func(c *gin.Context) {
				c.Writer.Header().Set("Cache-Control", "max-age=0, must-revalidate, no-cache")
				c.Next()
			}
		}
	}

	m.Static(prefix, path.Join(setting.StaticRootPath, dir))
	m.Use(headers())
}

// writePIDFile retrieves the current process ID and writes it to file.
func (s *Server) writePIDFile() {
	if s.pidFile == "" {
		return
	}

	// Ensure the required directory structure exists.
	err := os.MkdirAll(filepath.Dir(s.pidFile), 0700)
	if err != nil {
		log.Error("Failed to verify pid directory", "error", err)
		os.Exit(1)
	}

	// Retrieve the PID and write it to file.
	pid := strconv.Itoa(os.Getpid())
	if err := ioutil.WriteFile(s.pidFile, []byte(pid), 0644); err != nil {
		log.Error("Failed to write pidfile", "error", err)
		os.Exit(1)
	}

	log.Info("Writing PID file", "path", s.pidFile, "pid", pid)
}
