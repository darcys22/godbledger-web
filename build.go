// +build ignore

package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	//"encoding/json"
	"encoding/base64"
	"flag"
	"fmt"
	//"go/build"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/cespare/cp"

	"github.com/darcys22/godbledger-web/internal/build"
)

const (
	windows = "windows"
	linux   = "linux"
)

var (
	//versionRe = regexp.MustCompile(`-[0-9]{1,3}-g[0-9a-f]{5,10}`)
	goarch  string
	goos    string
	gocc    string
	cgo     bool
	libc    string
	pkgArch string
	version string = "v1"
	// deb & rpm does not support semver so have to handle their version a little differently
	linuxPackageVersion   string = "v1"
	linuxPackageIteration string = ""
	race                  bool
	workingDir            string
	includeBuildId        bool     = true
	buildId               string   = "0"
	serverBinary          string   = "godbledger-web"
	binaries              []string = []string{serverBinary}
	isDev                 bool     = false
	enterprise            bool     = false
	skipRpmGen            bool     = false
	skipDebGen            bool     = false
	printGenVersion       bool     = false
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	var buildIdRaw string

	flag.StringVar(&goarch, "goarch", runtime.GOARCH, "GOARCH")
	flag.StringVar(&goos, "goos", runtime.GOOS, "GOOS")
	flag.StringVar(&gocc, "cc", "", "CC")
	flag.StringVar(&libc, "libc", "", "LIBC")
	flag.BoolVar(&cgo, "cgo-enabled", cgo, "Enable cgo")
	flag.StringVar(&pkgArch, "pkg-arch", "", "PKG ARCH")
	flag.BoolVar(&race, "race", race, "Use race detector")
	flag.BoolVar(&includeBuildId, "includeBuildId", includeBuildId, "IncludeBuildId in package name")
	flag.BoolVar(&enterprise, "enterprise", enterprise, "Build enterprise version of GoDBLedger")
	flag.StringVar(&buildIdRaw, "buildId", "0", "Build ID from CI system")
	flag.BoolVar(&isDev, "dev", isDev, "optimal for development, skips certain steps")
	flag.BoolVar(&skipRpmGen, "skipRpm", skipRpmGen, "skip rpm package generation (default: false)")
	flag.BoolVar(&skipDebGen, "skipDeb", skipDebGen, "skip deb package generation (default: false)")
	flag.BoolVar(&printGenVersion, "gen-version", printGenVersion, "generate GoDBLedger version and output (default: false)")
	flag.Parse()

	buildId = shortenBuildId(buildIdRaw)

	readVersionFromPackageJson()

	if pkgArch == "" {
		pkgArch = goarch
	}

	if printGenVersion {
		printGeneratedVersion()
		return
	}

	log.Printf("Version: %s, Linux Version: %s, Package Iteration: %s\n", version, linuxPackageVersion, linuxPackageIteration)

	if flag.NArg() == 0 {
		log.Println("Usage: go run build.go build")
		return
	}

	workingDir, _ = os.Getwd()

	for _, cmd := range flag.Args() {
		switch cmd {
		case "setup":
			setup()

		case "build-srv", "build-server":
			doBuild("godbledger-web", "./pkg/cmd/godbledger-web", []string{})

		case "build":
			for _, binary := range binaries {
				doBuild(binary, "./pkg/cmd/"+binary, []string{})
			}

		case "build-frontend":
			grunt(gruntBuildArg("build")...)

		case "test":
			test("./pkg/...")
			grunt("test")

		case "package":
			grunt(gruntBuildArg("build")...)
			grunt(gruntBuildArg("package")...)
			if goos == linux {
				createLinuxPackages()
			}

		case "package-only":
			grunt(gruntBuildArg("package")...)
			if goos == linux {
				createLinuxPackages()
			}
		case "pkg-archive":
			grunt(gruntBuildArg("package")...)

		case "pkg-rpm":
			grunt(gruntBuildArg("release")...)
			createRpmPackages()

		case "pkg-deb":
			doDebianSource([]string{})

		case "sha-dist":
			shaFilesInDist()

		case "latest":
			makeLatestDistCopies()

		default:
			log.Fatalf("Unknown command %q", cmd)
		}
	}
}

func makeLatestDistCopies() {
	files, err := ioutil.ReadDir("dist")
	if err != nil {
		log.Fatalf("failed to create latest copies. Cannot read from /dist")
	}

	latestMapping := map[string]string{
		"_amd64.deb":               "dist/godbledger_latest_amd64.deb",
		".x86_64.rpm":              "dist/godbledger-latest-1.x86_64.rpm",
		".linux-amd64.tar.gz":      "dist/godbledger-latest.linux-x64.tar.gz",
		".linux-amd64-musl.tar.gz": "dist/godbledger-latest.linux-x64-musl.tar.gz",
		".linux-armv7.tar.gz":      "dist/godbledger-latest.linux-armv7.tar.gz",
		".linux-armv7-musl.tar.gz": "dist/godbledger-latest.linux-armv7-musl.tar.gz",
		".linux-armv6.tar.gz":      "dist/godbledger-latest.linux-armv6.tar.gz",
		".linux-arm64.tar.gz":      "dist/godbledger-latest.linux-arm64.tar.gz",
		".linux-arm64-musl.tar.gz": "dist/godbledger-latest.linux-arm64-musl.tar.gz",
	}

	for _, file := range files {
		for extension, fullName := range latestMapping {
			if strings.HasSuffix(file.Name(), extension) {
				runError("cp", path.Join("dist", file.Name()), fullName)
			}
		}
	}
}

func readVersionFromPackageJson() {
	//reader, err := os.Open("package.json")
	//if err != nil {
	//log.Fatal("Failed to open package.json")
	//return
	//}
	//defer reader.Close()

	//jsonObj := map[string]interface{}{}
	//jsonParser := json.NewDecoder(reader)

	//if err := jsonParser.Decode(&jsonObj); err != nil {
	//log.Fatal("Failed to decode package.json")
	//}

	//version = jsonObj["version"].(string)
	version = "0.0.1"
	linuxPackageVersion = version
	linuxPackageIteration = ""

	// handle pre version stuff (deb / rpm does not support semver)
	parts := strings.Split(version, "-")

	if len(parts) > 1 {
		linuxPackageVersion = parts[0]
		linuxPackageIteration = parts[1]
	}

	// add timestamp to iteration
	if includeBuildId {
		if buildId != "0" {
			linuxPackageIteration = fmt.Sprintf("%s%s", buildId, linuxPackageIteration)
		} else {
			linuxPackageIteration = fmt.Sprintf("%d%s", time.Now().Unix(), linuxPackageIteration)
		}
	}
}

type linuxPackageOptions struct {
	packageType            string
	packageArch            string
	homeDir                string
	homeBinDir             string
	binPath                string
	serverBinPath          string
	cliBinPath             string
	configDir              string
	ldapFilePath           string
	etcDefaultPath         string
	etcDefaultFilePath     string
	initdScriptFilePath    string
	systemdServiceFilePath string

	postinstSrc         string
	initdScriptSrc      string
	defaultFileSrc      string
	systemdFileSrc      string
	cliBinaryWrapperSrc string

	depends []string
}

var GOBIN, _ = filepath.Abs(filepath.Join("bin"))

var (
	// A debian package is created for all executables listed here.
	debExecutables = []debExecutable{
		{
			BinaryName:  "godbledger-web",
			Description: "Godbledger Web Server",
		},
	}

	// A debian package is created for all executables listed here.
	debGoDBLedger = debPackage{
		Name: "godbledger",
		//Version:     version.Version,
		Version:     "0.0.1",
		Executables: debExecutables,
	}

	// Debian meta packages to build and push to Ubuntu PPA
	debPackages = []debPackage{
		debGoDBLedger,
	}

	// Distros for which packages are created.
	debDistroGoBoots = map[string]string{
		"xenial":  "golang-go",
		"bionic":  "golang-go",
		"focal":   "golang-go",
		"groovy":  "golang-go",
		"hirsute": "golang-go",
	}

	debGoBootPaths = map[string]string{
		"golang-go": "/usr/lib/go",
	}

	dlgoVersion = "1.16"
)

// skips archiving for some build configurations.
func maybeSkipArchive(env build.Environment) {
	if env.IsPullRequest {
		log.Printf("skipping because this is a PR build")
		os.Exit(0)
	}
	if env.IsCronJob {
		log.Printf("skipping because this is a cron job")
		os.Exit(0)
	}
	//if env.Branch != "master" && !strings.HasPrefix(env.Tag, "v1.") {
	//log.Printf("skipping because branch %q, tag %q is not on the whitelist", env.Branch, env.Tag)
	//os.Exit(0)
	//}
}

func doDebianSource(cmdline []string) {
	var (
		cachedir = flag.String("cachedir", "./build/cache", `Filesystem path to cache the downloaded Go bundles at`)
		signer   = flag.String("signer", "", `Signing key name, also used as package author`)
		upload   = flag.String("upload", "", `Where to upload the source package (usually "darcys22/godbledger")`)
		sshUser  = flag.String("sftp-user", "", `Username for SFTP upload (usually "darcys22")`)
		workdir  = flag.String("workdir", "", `Output directory for packages (uses temp dir if unset)`)
		now      = time.Now()
	)
	flag.CommandLine.Parse(cmdline)
	*workdir = makeWorkdir(*workdir)
	env := build.Env()
	maybeSkipArchive(env)

	// Import the signing key.
	//gpg --export-secret-key sean@darcyfinancial.com  | base64 | paste -s -d '' > secret-signing-key-base64-encoded.gpg
	if key := getenvBase64("PPA_SIGNING_KEY"); len(key) > 0 {
		gpg := exec.Command("gpg", "--import", "--no-tty", "--batch", "--yes")
		gpg.Stdin = bytes.NewReader(key)
		build.MustRun(gpg)
	}

	// Download and verify the Go source package.
	gobundle := downloadGoSources(*cachedir)

	// Download all the dependencies needed to build the sources and run the ci script
	srcdepfetch := goTool("mod", "download")
	gopath, _ := filepath.Abs(filepath.Join(*workdir, "modgopath"))
	srcdepfetch.Env = append(os.Environ(), "GOPATH="+gopath)
	build.MustRun(srcdepfetch)

	cidepfetch := goTool("run", "./build.go")
	cidepfetch.Env = append(os.Environ(), "GOPATH="+filepath.Join(*workdir, "modgopath"))
	cidepfetch.Run() // Command fails, don't care, we only need the deps to start it

	// Create Debian packages and upload them.
	for _, pkg := range debPackages {
		for distro, goboot := range debDistroGoBoots {
			// Prepare the debian package with the go-ethereum sources.
			meta := newDebMetadata(distro, goboot, *signer, env, now, pkg.Name, pkg.Version, pkg.Executables)
			fmt.Println("Building debian package in: " + *workdir)
			pkgdir := stageDebianSource(*workdir, meta)

			// Add Go source code
			if err := build.ExtractArchive(gobundle, pkgdir); err != nil {
				log.Fatalf("Failed to extract Go sources: %v", err)
			}
			if err := os.Rename(filepath.Join(pkgdir, "go"), filepath.Join(pkgdir, ".go")); err != nil {
				log.Fatalf("Failed to rename Go source folder: %v", err)
			}
			// Add all dependency modules in compressed form
			os.MkdirAll(filepath.Join(pkgdir, ".mod", "cache"), 0755)
			if err := cp.CopyAll(filepath.Join(pkgdir, ".mod", "cache", "download"), filepath.Join(*workdir, "modgopath", "pkg", "mod", "cache", "download")); err != nil {
				log.Fatalf("Failed to copy Go module dependencies: %v", err)
			}
			// Run the packaging and upload to the PPA
			debuild := exec.Command("debuild", "-S", "-sa", "-us", "-uc", "-d", "-Zxz", "-nc")
			debuild.Dir = pkgdir
			build.MustRun(debuild)

			var (
				basename  = fmt.Sprintf("%s_%s", meta.Name(), meta.VersionString())
				source    = filepath.Join(*workdir, basename+".tar.xz")
				dsc       = filepath.Join(*workdir, basename+".dsc")
				changes   = filepath.Join(*workdir, basename+"_source.changes")
				buildinfo = filepath.Join(*workdir, basename+"_source.buildinfo")
			)
			if *signer != "" {
				debsign := exec.Command("debsign", changes)
				build.MustRun(debsign)
			}
			if *upload != "" {
				ppaUpload(*workdir, *upload, *sshUser, []string{source, dsc, changes, buildinfo})
			}
		}
	}
}

// downloadGoSources downloads the Go source tarball.
func downloadGoSources(cachedir string) string {
	csdb := build.MustLoadChecksums("utils/checksums.txt")
	file := fmt.Sprintf("go%s.src.tar.gz", dlgoVersion)
	url := "https://dl.google.com/go/" + file
	dst := filepath.Join(cachedir, file)
	if err := csdb.DownloadFile(url, dst); err != nil {
		log.Fatal(err)
	}
	return dst
}

func ppaUpload(workdir, ppa, sshUser string, files []string) {
	p := strings.Split(ppa, "/")
	if len(p) != 2 {
		log.Fatal("-upload PPA name must contain single /")
	}
	if sshUser == "" {
		sshUser = p[0]
	}
	incomingDir := fmt.Sprintf("~%s/ubuntu/%s", p[0], p[1])
	// Create the SSH identity file if it doesn't exist.
	var idfile string
	//cat ppakey  | base64 | paste -s -d '' > secret-ssh-key-base64-encoded
	if sshkey := getenvBase64("PPA_SSH_KEY"); len(sshkey) > 0 {
		idfile = filepath.Join(workdir, "sshkey")
		if _, err := os.Stat(idfile); os.IsNotExist(err) {
			ioutil.WriteFile(idfile, sshkey, 0600)
		}
	}
	// Upload
	dest := sshUser + "@ppa.launchpad.net"
	if err := build.UploadSFTP(idfile, dest, incomingDir, files); err != nil {
		log.Fatal(err)
	}
}

func getenvBase64(variable string) []byte {
	dec, err := base64.StdEncoding.DecodeString(os.Getenv(variable))
	if err != nil {
		log.Fatal("invalid base64 " + variable)
	}
	return []byte(dec)
}

func makeWorkdir(wdflag string) string {
	var err error
	if wdflag != "" {
		err = os.MkdirAll(wdflag, 0744)
	} else {
		wdflag, err = ioutil.TempDir("", "godbledger-build-")
	}
	if err != nil {
		log.Fatal(err)
	}
	return wdflag
}

func isUnstableBuild(env build.Environment) bool {
	if env.Tag != "" {
		return false
	}
	return true
}

type debPackage struct {
	Name        string          // the name of the Debian package to produce, e.g. "godbledger"
	Version     string          // the clean version of the debPackage, e.g. 1.8.12, without any metadata
	Executables []debExecutable // executables to be included in the package
}

type debMetadata struct {
	Env           build.Environment
	GoBootPackage string
	GoBootPath    string

	PackageName string

	// go-ethereum version being built. Note that this
	// is not the debian package version. The package version
	// is constructed by VersionString.
	Version string

	Author       string // "name <email>", also selects signing key
	Distro, Time string
	Executables  []debExecutable
}

type debExecutable struct {
	PackageName string
	BinaryName  string
	Description string
}

// Package returns the name of the package if present, or
// fallbacks to BinaryName
func (d debExecutable) Package() string {
	if d.PackageName != "" {
		return d.PackageName
	}
	return d.BinaryName
}

func newDebMetadata(distro, goboot, author string, env build.Environment, t time.Time, name string, version string, exes []debExecutable) debMetadata {
	if author == "" {
		// No signing key, use default author.
		author = "Sean Darcy <sean@darcyfinanical.com>"
	}
	return debMetadata{
		GoBootPackage: goboot,
		GoBootPath:    debGoBootPaths[goboot],
		PackageName:   name,
		Env:           env,
		Author:        author,
		Distro:        distro,
		Version:       version,
		Time:          t.Format(time.RFC1123Z),
		Executables:   exes,
	}
}

// Name returns the name of the metapackage that depends
// on all executable packages.
func (meta debMetadata) Name() string {
	if isUnstableBuild(meta.Env) {
		return meta.PackageName + "-unstable"
	}
	return meta.PackageName
}

// VersionString returns the debian version of the packages.
func (meta debMetadata) VersionString() string {
	vsn := meta.Version
	if meta.Env.Buildnum != "" {
		vsn += "+build" + meta.Env.Buildnum
	}
	if meta.Distro != "" {
		vsn += "+" + meta.Distro
	}
	return vsn
}

// ExeList returns the list of all executable packages.
func (meta debMetadata) ExeList() string {
	names := make([]string, len(meta.Executables))
	for i, e := range meta.Executables {
		names[i] = meta.ExeName(e)
	}
	return strings.Join(names, ", ")
}

// ExeName returns the package name of an executable package.
func (meta debMetadata) ExeName(exe debExecutable) string {
	if isUnstableBuild(meta.Env) {
		return exe.Package() + "-unstable"
	}
	return exe.Package()
}

// ExeConflicts returns the content of the Conflicts field
// for executable packages.
func (meta debMetadata) ExeConflicts(exe debExecutable) string {
	if isUnstableBuild(meta.Env) {
		// Set up the conflicts list so that the *-unstable packages
		// cannot be installed alongside the regular version.
		//
		// https://www.debian.org/doc/debian-policy/ch-relationships.html
		// is very explicit about Conflicts: and says that Breaks: should
		// be preferred and the conflicting files should be handled via
		// alternates. We might do this eventually but using a conflict is
		// easier now.
		return "godbledger, " + exe.Package()
	}
	return ""
}

func stageDebianSource(tmpdir string, meta debMetadata) (pkgdir string) {
	pkg := meta.Name() + "-" + meta.VersionString()
	pkgdir = filepath.Join(tmpdir, pkg)
	if err := os.Mkdir(pkgdir, 0755); err != nil {
		log.Fatal(err)
	}
	// Copy the source code.
	build.MustRunCommand("git", "checkout-index", "-a", "--prefix", pkgdir+string(filepath.Separator))

	// Put the debian build files in place.
	debian := filepath.Join(pkgdir, "debian")
	build.Render("utils/deb/deb.rules", filepath.Join(debian, "rules"), 0755, meta)
	build.Render("utils/deb/deb.changelog", filepath.Join(debian, "changelog"), 0644, meta)
	build.Render("utils/deb/deb.control", filepath.Join(debian, "control"), 0644, meta)
	build.Render("utils/deb/deb.copyright", filepath.Join(debian, "copyright"), 0644, meta)
	build.RenderString("8\n", filepath.Join(debian, "compat"), 0644, meta)
	build.RenderString("3.0 (native)\n", filepath.Join(debian, "source/format"), 0644, meta)
	for _, exe := range meta.Executables {
		install := filepath.Join(debian, meta.ExeName(exe)+".install")
		build.Render("utils/deb/deb.install", install, 0644, exe)

		docs := filepath.Join(debian, meta.ExeName(exe)+".docs")
		build.Render("utils/deb/deb.docs", docs, 0644, exe)

		if exe.PackageName == "godbledger-web" {
			preinst := filepath.Join(debian, meta.ExeName(exe)+".preinst")
			build.Render("utils/deb/godbledger-web.preinst", preinst, 0644, meta)

			postinst := filepath.Join(debian, meta.ExeName(exe)+".postinst")
			build.Render("utils/deb/godbledger-web.postinst", postinst, 0644, meta)

			prerm := filepath.Join(debian, meta.ExeName(exe)+".prerm")
			build.Render("utils/deb/godbledger-web.prerm", prerm, 0644, meta)

			postrm := filepath.Join(debian, meta.ExeName(exe)+".postrm")
			build.Render("utils/deb/godbledger-web.postrm", postrm, 0644, meta)

			servicefile := filepath.Join(debian, meta.ExeName(exe)+".service")
			build.Render("utils/deb/godbledger-web.service", servicefile, 0644, meta)
		}
	}
	return pkgdir
}

func createRpmPackages() {
	rpmPkgArch := pkgArch
	switch {
	case pkgArch == "armv7":
		rpmPkgArch = "armhfp"
	case pkgArch == "arm64":
		rpmPkgArch = "aarch64"
	}
	createPackage(linuxPackageOptions{
		packageType:            "rpm",
		packageArch:            rpmPkgArch,
		homeDir:                "/usr/share/godbledger/",
		homeBinDir:             "/usr/share/godbledger/bin",
		binPath:                "/usr/sbin",
		configDir:              "/etc/godbledger",
		etcDefaultPath:         "/etc/sysconfig",
		etcDefaultFilePath:     "/etc/sysconfig/godbledger-web",
		initdScriptFilePath:    "/etc/init.d/godbledger-web",
		systemdServiceFilePath: "/usr/lib/systemd/system/godbledger-web.service",

		postinstSrc:         "packaging/rpm/control/postinst",
		initdScriptSrc:      "packaging/rpm/init.d/godbledger-web",
		defaultFileSrc:      "packaging/rpm/sysconfig/godbledger-web",
		systemdFileSrc:      "packaging/rpm/systemd/godbledger-web.service",
		cliBinaryWrapperSrc: "packaging/wrappers/godbledger-cli",

		depends: []string{"/sbin/service", "fontconfig", "freetype", "urw-fonts"},
	})
}

func createLinuxPackages() {
	if !skipDebGen {
		doDebianSource([]string{})
	}

	if !skipRpmGen {
		createRpmPackages()
	}
}

func createPackage(options linuxPackageOptions) {
	packageRoot, _ := ioutil.TempDir("", "godbledger-web-linux-pack")

	// create directories
	runPrint("mkdir", "-p", filepath.Join(packageRoot, options.homeDir))
	runPrint("mkdir", "-p", filepath.Join(packageRoot, options.configDir))
	runPrint("mkdir", "-p", filepath.Join(packageRoot, "/etc/init.d"))
	runPrint("mkdir", "-p", filepath.Join(packageRoot, options.etcDefaultPath))
	runPrint("mkdir", "-p", filepath.Join(packageRoot, "/usr/lib/systemd/system"))
	runPrint("mkdir", "-p", filepath.Join(packageRoot, "/usr/sbin"))

	// copy godbledger-cli wrapper
	//runPrint("cp", "-p", options.cliBinaryWrapperSrc, filepath.Join(packageRoot, "/usr/sbin/"+cliBinary))

	// copy godbledger-web binary
	runPrint("cp", "-p", filepath.Join(workingDir, "tmp/bin/"+serverBinary), filepath.Join(packageRoot, "/usr/sbin/"+serverBinary))

	// copy init.d script
	runPrint("cp", "-p", options.initdScriptSrc, filepath.Join(packageRoot, options.initdScriptFilePath))
	// copy environment var file
	runPrint("cp", "-p", options.defaultFileSrc, filepath.Join(packageRoot, options.etcDefaultFilePath))
	// copy systemd file
	runPrint("cp", "-p", options.systemdFileSrc, filepath.Join(packageRoot, options.systemdServiceFilePath))
	// copy release files
	runPrint("cp", "-a", filepath.Join(workingDir, "tmp")+"/.", filepath.Join(packageRoot, options.homeDir))
	// remove bin path
	runPrint("rm", "-rf", filepath.Join(packageRoot, options.homeDir, "bin"))

	// create /bin within home
	runPrint("mkdir", "-p", filepath.Join(packageRoot, options.homeBinDir))
	// The GoDBLedger-cli binary is exposed through a wrapper to ensure a proper
	// configuration is in place. To enable that, we need to store the original
	// binary in a separate location to avoid conflicts.
	//runPrint("cp", "-p", filepath.Join(workingDir, "tmp/bin/"+cliBinary), filepath.Join(packageRoot, options.homeBinDir, cliBinary))

	args := []string{
		"-s", "dir",
		"--description", "godbledger",
		"-C", packageRoot,
		"--url", "https://godbledger.com/",
		"--maintainer", "sean@darcyfinancial.com",
		"--config-files", options.initdScriptFilePath,
		"--config-files", options.etcDefaultFilePath,
		"--config-files", options.systemdServiceFilePath,
		"--after-install", options.postinstSrc,

		"--version", linuxPackageVersion,
		"-p", "./dist",
	}

	name := "godbledger-web"
	if enterprise {
		name += "-enterprise"
		args = append(args, "--replaces", "godbledger")
	}
	fmt.Printf("pkgArch is set to '%s', generated arch is '%s'\n", pkgArch, options.packageArch)
	if pkgArch == "armv6" {
		name += "-rpi"
		args = append(args, "--replaces", "godbledger")
	}
	args = append(args, "--name", name)

	description := "GoDBLedger-Web"
	if enterprise {
		description += " Enterprise"
	}

	if !enterprise {
		args = append(args, "--license", "\"Apache 2.0\"")
	}

	if options.packageType == "rpm" {
		args = append(args, "--rpm-posttrans", "packaging/rpm/control/posttrans")
	}

	if options.packageType == "deb" {
		args = append(args, "--deb-no-default-config-files")
	}

	if options.packageArch != "" {
		args = append(args, "-a", options.packageArch)
	}

	if linuxPackageIteration != "" {
		args = append(args, "--iteration", linuxPackageIteration)
	}

	// add dependencies
	for _, dep := range options.depends {
		args = append(args, "--depends", dep)
	}

	args = append(args, ".")

	fmt.Println("Creating package: ", options.packageType)
	runPrint("fpm", append([]string{"-t", options.packageType}, args...)...)
}

func grunt(params ...string) {
	if runtime.GOOS == windows {
		runPrint(`.\node_modules\.bin\grunt`, params...)
	} else {
		runPrint("./node_modules/.bin/grunt", params...)
	}
}

func genPackageVersion() string {
	if includeBuildId {
		return fmt.Sprintf("%v-%v", linuxPackageVersion, linuxPackageIteration)
	} else {
		return version
	}
}

func gruntBuildArg(task string) []string {
	args := []string{task}
	args = append(args, fmt.Sprintf("--pkgVer=%v", genPackageVersion()))
	if pkgArch != "" {
		args = append(args, fmt.Sprintf("--arch=%v", pkgArch))
	}
	if libc != "" {
		args = append(args, fmt.Sprintf("--libc=%s", libc))
	}
	if enterprise {
		args = append(args, "--enterprise")
	}

	args = append(args, fmt.Sprintf("--platform=%v", goos))

	return args
}

func setup() {
	runPrint("go", "install", "-v", "./pkg/cmd/godbledger-web")
}

func printGeneratedVersion() {
	fmt.Print(genPackageVersion())
}

func test(pkg string) {
	setBuildEnv()
	runPrint("go", "test", "-short", "-timeout", "60s", pkg)
}

func doBuild(binaryName, pkg string, tags []string) {
	libcPart := ""
	if libc != "" {
		libcPart = fmt.Sprintf("-%s", libc)
	}
	binary := fmt.Sprintf("./bin/%s-%s%s/%s", goos, goarch, libcPart, binaryName)
	if isDev {
		//don't include os/arch/libc in output path in dev environment
		binary = fmt.Sprintf("./bin/%s", binaryName)
	}

	if goos == windows {
		binary += ".exe"
	}

	if !isDev {
		rmr(binary, binary+".md5")
	}
	args := []string{"build", "-ldflags", ldflags()}
	if len(tags) > 0 {
		args = append(args, "-tags", strings.Join(tags, ","))
	}
	if race {
		args = append(args, "-race")
	}

	args = append(args, "-o", binary)
	args = append(args, pkg)

	if !isDev {
		setBuildEnv()
		runPrint("go", "version")
		libcPart := ""
		if libc != "" {
			libcPart = fmt.Sprintf("/%s", libc)
		}
		fmt.Printf("Targeting %s/%s%s\n", goos, goarch, libcPart)
	}

	runPrint("go", args...)

	if !isDev {
		// Create an md5 checksum of the binary, to be included in the archive for
		// automatic upgrades.
		err := md5File(binary)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func ldflags() string {
	var b bytes.Buffer
	b.WriteString("-w")
	b.WriteString(fmt.Sprintf(" -X main.version=%s", version))
	b.WriteString(fmt.Sprintf(" -X main.commit=%s", getGitSha()))
	b.WriteString(fmt.Sprintf(" -X main.buildstamp=%d", buildStamp()))
	b.WriteString(fmt.Sprintf(" -X main.buildBranch=%s", getGitBranch()))
	if v := os.Getenv("LDFLAGS"); v != "" {
		b.WriteString(fmt.Sprintf(" -extldflags \"%s\"", v))
	}
	return b.String()
}

func rmr(paths ...string) {
	for _, path := range paths {
		log.Println("rm -r", path)
		os.RemoveAll(path)
	}
}

func setBuildEnv() {
	os.Setenv("GOOS", goos)
	if goos == windows {
		// require windows >=7
		os.Setenv("CGO_CFLAGS", "-D_WIN32_WINNT=0x0601")
	}
	if goarch != "amd64" || goos != linux {
		// needed for all other archs
		cgo = true
	}
	if strings.HasPrefix(goarch, "armv") {
		os.Setenv("GOARCH", "arm")
		os.Setenv("GOARM", goarch[4:])
	} else {
		os.Setenv("GOARCH", goarch)
	}
	if goarch == "386" {
		os.Setenv("GO386", "387")
	}
	if cgo {
		os.Setenv("CGO_ENABLED", "1")
	}
	if gocc != "" {
		os.Setenv("CC", gocc)
	}
}

func getGitBranch() string {
	v, err := runError("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "master"
	}
	return string(v)
}

func getGitSha() string {
	v, err := runError("git", "rev-parse", "--short", "HEAD")
	if err != nil {
		return "unknown-dev"
	}
	return string(v)
}

func buildStamp() int64 {
	// use SOURCE_DATE_EPOCH if set.
	if s, _ := strconv.ParseInt(os.Getenv("SOURCE_DATE_EPOCH"), 10, 64); s > 0 {
		return s
	}

	bs, err := runError("git", "show", "-s", "--format=%ct")
	if err != nil {
		return time.Now().Unix()
	}
	s, _ := strconv.ParseInt(string(bs), 10, 64)
	return s
}

func runError(cmd string, args ...string) ([]byte, error) {
	ecmd := exec.Command(cmd, args...)
	bs, err := ecmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return bytes.TrimSpace(bs), nil
}

func runPrint(cmd string, args ...string) {
	log.Println(cmd, strings.Join(args, " "))
	ecmd := exec.Command(cmd, args...)
	ecmd.Env = append(os.Environ(), "GO111MODULE=on")
	ecmd.Stdout = os.Stdout
	ecmd.Stderr = os.Stderr
	err := ecmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func md5File(file string) error {
	fd, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fd.Close()

	h := md5.New()
	_, err = io.Copy(h, fd)
	if err != nil {
		return err
	}

	out, err := os.Create(file + ".md5")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(out, "%x\n", h.Sum(nil))
	if err != nil {
		return err
	}

	return out.Close()
}

func shaFilesInDist() {
	filepath.Walk("./dist", func(path string, f os.FileInfo, err error) error {
		if path == "./dist" {
			return nil
		}

		if !strings.Contains(path, ".sha256") {
			err := shaFile(path)
			if err != nil {
				log.Printf("Failed to create sha file. error: %v\n", err)
			}
		}
		return nil
	})
}

func shaFile(file string) error {
	fd, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fd.Close()

	h := sha256.New()
	_, err = io.Copy(h, fd)
	if err != nil {
		return err
	}

	out, err := os.Create(file + ".sha256")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(out, "%x\n", h.Sum(nil))
	if err != nil {
		return err
	}

	return out.Close()
}

func shortenBuildId(buildId string) string {
	buildId = strings.Replace(buildId, "-", "", -1)
	if len(buildId) < 9 {
		return buildId
	}
	return buildId[0:8]
}

func goTool(subcmd string, args ...string) *exec.Cmd {
	return goToolArch(runtime.GOARCH, os.Getenv("CC"), subcmd, args...)
}

func goToolArch(arch string, cc string, subcmd string, args ...string) *exec.Cmd {
	cmd := build.GoTool(subcmd, args...)
	if arch == "" || arch == runtime.GOARCH {
		cmd.Env = append(cmd.Env, "GOBIN="+GOBIN)
	} else {
		cmd.Env = append(cmd.Env, "CGO_ENABLED=1")
		cmd.Env = append(cmd.Env, "GOARCH="+arch)
	}
	if cc != "" {
		cmd.Env = append(cmd.Env, "CC="+cc)
	}
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "GOBIN=") {
			continue
		}
		cmd.Env = append(cmd.Env, e)
	}
	return cmd
}
