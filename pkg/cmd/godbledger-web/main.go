package main

import (
	"flag"
	"fmt"
	//"net/http"
	_ "net/http/pprof"
	"os"
	//"os/signal"
	//"runtime"
	//"runtime/trace"
	"strconv"
	//"syscall"
	"time"

	"github.com/darcys22/godbledger-web/pkg/server"
	"github.com/darcys22/godbledger-web/pkg/setting"
)

var version = "0.0.1"
var commit = "NA"
var buildBranch = "master"
var buildstamp string

func main() {
	var (
		configFile = flag.String("config", "", "path to config file")
		homePath   = flag.String("homepath", "", "path to grafana install/home path, defaults to working directory")
		pidFile    = flag.String("pidfile", "", "path to pid file")
		//packaging  = flag.String("packaging", "unknown", "describes the way Grafana was installed")

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
	//setting.Packaging = validPackaging(*packaging)

	s, err := server.New(server.Config{
		ConfigFile: *configFile, HomePath: *homePath, PidFile: *pidFile,
		Version: version, Commit: commit, BuildBranch: buildBranch,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	err = s.Run()
	code := 0
	if err != nil {
		code = s.ExitCode(err)
	}
	//log.Close()

	os.Exit(code)
}
