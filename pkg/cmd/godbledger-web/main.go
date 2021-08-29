package main
import (
	"flag"
	"fmt"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"

	"github.com/darcys22/godbledger-web/pkg/server"
	"github.com/darcys22/godbledger-web/pkg/setting"
)

var version = "0.0.1"
var commit = "NA"
var buildBranch = "master"
var buildstamp string

func main() {
	customFormatter := new(prefixed.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)
	var (
		configFile = flag.String("config", "", "path to config file")
		homePath   = flag.String("homepath", "", "path to grafana install/home path, defaults to working directory")
		pidFile    = flag.String("pidfile", "", "path to pid file")

		v = flag.Bool("v", false, "prints current version and exits")
	)

	flag.Parse()

	if *v {
		fmt.Printf("Version %s (commit: %s, branch: %s)\n", version, commit, buildBranch)
		os.Exit(0)
	}

	buildstampInt64, _ := strconv.ParseInt(buildstamp, 10, 64)
	if buildstampInt64 == 0 {
		buildstampInt64 = time.Now().Unix()
	}

	setting.BuildVersion = version
	setting.BuildCommit = commit
	setting.BuildStamp = buildstampInt64
	setting.BuildBranch = buildBranch

	s, err := server.New(server.Config{
		ConfigFile: *configFile, HomePath: *homePath, PidFile: *pidFile,
		Version: version, Commit: commit, BuildBranch: buildBranch,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	go listenToSystemSignals(s)

	err = s.Run()
	code := 0
	if err != nil {
		code = s.ExitCode(err)
	}

	os.Exit(code)
}

func listenToSystemSignals(s *server.Server) {
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case sig := <-signalChan:
			s.Shutdown(fmt.Sprintf("System signal: %s", sig))
		}
	}
}
