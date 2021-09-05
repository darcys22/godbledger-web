package server

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/darcys22/godbledger-web/backend/api"
	"github.com/darcys22/godbledger-web/backend/setting"

	"github.com/sirupsen/logrus"
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

	s := &Server{
		context:       rootCtx,
		shutdownFn:    shutdownFn,
		cfg:           setting.NewCfg(),

		configFile:  cfg.ConfigFile,
		homePath:    cfg.HomePath,
		pidFile:     cfg.PidFile,
		version:     cfg.Version,
		commit:      cfg.Commit,
		buildBranch: cfg.BuildBranch,
	}
	if cfg.Listener != nil {
		if err := s.init(); err != nil {
			return nil, err
		}
	}

	return s, nil
}

// Server is responsible for managing the lifecycle of services.
type Server struct {
	context            context.Context
	shutdownFn         context.CancelFunc
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
func (s *Server) init() error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.isInitialized {
		return nil
	}
	s.isInitialized = true

	s.loadConfiguration()
	s.writePIDFile()

	return nil
}

// Run initializes and starts services. This will block until all services have
// exited. To initiate shutdown, call the Shutdown method in another goroutine.
func (s *Server) Run() (err error) {

	if err = s.init(); err != nil {
		return
	}

	m := api.NewGin()

	listenAddr := fmt.Sprintf("%s:%s", setting.HttpAddr, setting.HttpPort)
	log.Infof("Listen: %v://%s", setting.Protocol, listenAddr)

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
