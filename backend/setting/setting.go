package setting

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	ini "gopkg.in/ini.v1"
)

var (
	log = logrus.WithField("prefix", "setting")
	globalcfg *Cfg
)

type Scheme string

const (
	HTTP              Scheme = "http"
	HTTPS             Scheme = "https"
	DEFAULT_HTTP_ADDR string = "0.0.0.0"
	REDACTED_PASSWORD string = "*********"
)

const (
	DEV      = "development"
	PROD     = "production"
	TEST     = "test"
	APP_NAME = "GoDBLedger-Web"
)

var (
	// App settings.
	Env              = DEV

	// build
	BuildVersion    string
	BuildCommit     string
	BuildBranch     string
	BuildStamp      int64
	IsEnterprise    bool
	ApplicationName string

	// Paths
	HomePath       string
	CustomInitPath = "conf/custom.ini"

	// Http server options
	Protocol           Scheme
	HttpAddr, HttpPort string
	CertFile, KeyFile  string
	StaticRootPath     string
	EnableGzip         bool

	// Global setting objects.
	Raw          *ini.File
	IsWindows    bool

	// for logging purposes
	configFiles                  []string
	appliedCommandLineProperties []string
	appliedEnvOverrides          []string
)

// TODO move all global vars to this struct
type Cfg struct {
	Raw *ini.File

	// HTTP Server Settings
	StaticRootPath   string
	Protocol         Scheme
	Domain           string

	// build
	BuildVersion string
	BuildCommit  string
	BuildBranch  string
	BuildStamp   int64
	IsEnterprise bool

	// security
	DisableInitialAdminCreation bool
	AdminUser string
	AdminPassword string

  // backend godbledger info
  DatabaseType string
  DatabaseLocation string
  GoDBLedgerHost string
  GoDBLedgerPort string
  GoDBLedgerCACert string
  GoDBLedgerCert string
  GoDBLedgerKey string

}

type CommandLineArgs struct {
	Config   string
	HomePath string
	Args     []string
}

func init() {
	IsWindows = runtime.GOOS == "windows"
}

func shouldRedactKey(s string) bool {
	uppercased := strings.ToUpper(s)
	return strings.Contains(uppercased, "PASSWORD") || strings.Contains(uppercased, "SECRET") || strings.Contains(uppercased, "PROVIDER_CONFIG")
}

func shouldRedactURLKey(s string) bool {
	uppercased := strings.ToUpper(s)
	return strings.Contains(uppercased, "DATABASE_URL")
}

func applyEnvVariableOverrides(file *ini.File) error {
	appliedEnvOverrides = make([]string, 0)
	for _, section := range file.Sections() {
		for _, key := range section.Keys() {
			envKey := envKey(section.Name(), key.Name())
			envValue := os.Getenv(envKey)

			if len(envValue) > 0 {
				key.SetValue(envValue)
				if shouldRedactKey(envKey) {
					envValue = REDACTED_PASSWORD
				}
				if shouldRedactURLKey(envKey) {
					u, err := url.Parse(envValue)
					if err != nil {
						return fmt.Errorf("could not parse environment variable. key: %s, value: %s. error: %v", envKey, envValue, err)
					}
					ui := u.User
					if ui != nil {
						_, exists := ui.Password()
						if exists {
							u.User = url.UserPassword(ui.Username(), "-redacted-")
							envValue = u.String()
						}
					}
				}
				appliedEnvOverrides = append(appliedEnvOverrides, fmt.Sprintf("%s=%s", envKey, envValue))
			}
		}
	}

	return nil
}

func envKey(sectionName string, keyName string) string {
	sN := strings.ToUpper(strings.Replace(sectionName, ".", "_", -1))
	sN = strings.Replace(sN, "-", "_", -1)
	kN := strings.ToUpper(strings.Replace(keyName, ".", "_", -1))
	envKey := fmt.Sprintf("GF_%s_%s", sN, kN)
	return envKey
}

func applyCommandLineDefaultProperties(props map[string]string, file *ini.File) {
	appliedCommandLineProperties = make([]string, 0)
	for _, section := range file.Sections() {
		for _, key := range section.Keys() {
			keyString := fmt.Sprintf("default.%s.%s", section.Name(), key.Name())
			value, exists := props[keyString]
			if exists {
				key.SetValue(value)
				if shouldRedactKey(keyString) {
					value = REDACTED_PASSWORD
				}
				appliedCommandLineProperties = append(appliedCommandLineProperties, fmt.Sprintf("%s=%s", keyString, value))
			}
		}
	}
}

func applyCommandLineProperties(props map[string]string, file *ini.File) {
	for _, section := range file.Sections() {
		sectionName := section.Name() + "."
		if section.Name() == ini.DefaultSection {
			sectionName = ""
		}
		for _, key := range section.Keys() {
			keyString := sectionName + key.Name()
			value, exists := props[keyString]
			if exists {
				appliedCommandLineProperties = append(appliedCommandLineProperties, fmt.Sprintf("%s=%s", keyString, value))
				key.SetValue(value)
			}
		}
	}
}

func getCommandLineProperties(args []string) map[string]string {
	props := make(map[string]string)

	for _, arg := range args {
		if !strings.HasPrefix(arg, "cfg:") {
			continue
		}

		trimmed := strings.TrimPrefix(arg, "cfg:")
		parts := strings.Split(trimmed, "=")
		if len(parts) != 2 {
			log.Fatalf("Invalid command line argument. argument: %v", arg)
			return nil
		}

		props[parts[0]] = parts[1]
	}
	return props
}

func makeAbsolute(path string, root string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(root, path)
}

func loadSpecifiedConfigFile(configFile string, masterFile *ini.File) error {
	if configFile == "" {
		configFile = filepath.Join(HomePath, CustomInitPath)
		// return without error if custom file does not exist
		if !pathExists(configFile) {
			return nil
		}
	}

	userConfig, err := ini.Load(configFile)
	if err != nil {
		return fmt.Errorf("Failed to parse %v, %v", configFile, err)
	}

	userConfig.BlockMode = false

	for _, section := range userConfig.Sections() {
		for _, key := range section.Keys() {
			if key.Value() == "" {
				continue
			}

			defaultSec, err := masterFile.GetSection(section.Name())
			if err != nil {
				defaultSec, _ = masterFile.NewSection(section.Name())
			}
			defaultKey, err := defaultSec.GetKey(key.Name())
			if err != nil {
				defaultKey, _ = defaultSec.NewKey(key.Name(), key.Value())
			}
			defaultKey.SetValue(key.Value())
		}
	}

	configFiles = append(configFiles, configFile)
	return nil
}

func (cfg *Cfg) loadConfiguration(args *CommandLineArgs) (*ini.File, error) {
	var err error

	// load config defaults
	defaultConfigFile := path.Join(HomePath, "conf/defaults.ini")
	configFiles = append(configFiles, defaultConfigFile)
	fmt.Println("Loading config from ", defaultConfigFile)

	// check if config file exists
	if _, err := os.Stat(defaultConfigFile); os.IsNotExist(err) {
		fmt.Println("GoDBLedger Init Failed: Could not find config defaults, make sure homepath command line parameter is set or working directory is homepath")
		os.Exit(1)
	}

	// load defaults
	parsedFile, err := ini.Load(defaultConfigFile)
	if err != nil {
		fmt.Printf("Failed to parse defaults.ini, %v\n", err)
		os.Exit(1)
		return nil, err
	}

	parsedFile.BlockMode = false

	// command line props
	commandLineProps := getCommandLineProperties(args.Args)
	// load default overrides
	applyCommandLineDefaultProperties(commandLineProps, parsedFile)

	// load specified config file
	err = loadSpecifiedConfigFile(args.Config, parsedFile)

	// apply environment overrides
	err = applyEnvVariableOverrides(parsedFile)
	if err != nil {
		return nil, err
	}

	// apply command line overrides
	applyCommandLineProperties(commandLineProps, parsedFile)

	return parsedFile, err
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func setHomePath(args *CommandLineArgs) {
	if args.HomePath != "" {
		HomePath = args.HomePath
		return
	}

	HomePath, _ = filepath.Abs(".")
	// check if homepath is correct
	if pathExists(filepath.Join(HomePath, "conf/defaults.ini")) {
		return
	}

	// try down one path
	if pathExists(filepath.Join(HomePath, "../conf/defaults.ini")) {
		HomePath = filepath.Join(HomePath, "../")
	}
}

func GetConfig() *Cfg {
	return globalcfg
}

func NewCfg() *Cfg {
	return &Cfg{
		Raw: ini.Empty(),
	}
}

func (cfg *Cfg) Load(args *CommandLineArgs) error {
	setHomePath(args)

	iniFile, err := cfg.loadConfiguration(args)
	if err != nil {
		return err
	}

	cfg.Raw = iniFile

	// Temporary keep global, to make refactor in steps
	Raw = cfg.Raw

	cfg.BuildVersion = BuildVersion
	cfg.BuildCommit = BuildCommit
	cfg.BuildStamp = BuildStamp
	cfg.BuildBranch = BuildBranch
	cfg.IsEnterprise = IsEnterprise

	ApplicationName = APP_NAME

	Env, err = valueAsString(iniFile.Section(""), "app_mode", "development")
	if err != nil {
		return err
	}
	if err := readServerSettings(iniFile, cfg); err != nil {
		return err
	}
	if err := readSecuritySettings(iniFile, cfg); err != nil {
		return err
	}
	if err := readBackendSettings(iniFile, cfg); err != nil {
		return err
	}


	cfg.Protocol = Protocol

	globalcfg = cfg

	return nil
}

func valueAsString(section *ini.Section, keyName string, defaultValue string) (string, error) {
	return section.Key(keyName).MustString(defaultValue), nil
}

func readServerSettings(iniFile *ini.File, cfg *Cfg) error {
	server := iniFile.Section("server")
	var err error

	Protocol = HTTP
	protocolStr, err := valueAsString(server, "protocol", "http")
	if err != nil {
		return err
	}
	if protocolStr == "https" {
		Protocol = HTTPS
		CertFile = server.Key("cert_file").String()
		KeyFile = server.Key("cert_key").String()
	}

	HttpAddr, err = valueAsString(server, "http_addr", DEFAULT_HTTP_ADDR)
	if err != nil {
		return err
	}
	HttpPort, err = valueAsString(server, "http_port", "3000")
	if err != nil {
		return err
	}

	cfg.Domain, err = valueAsString(server, "domain", "localhost")
	if err != nil {
		return err
	}

	EnableGzip = server.Key("enable_gzip").MustBool(false)
	staticRoot, err := valueAsString(server, "static_root_path", "")
	if err != nil {
		return err
	}
	StaticRootPath = makeAbsolute(staticRoot, HomePath)
	cfg.StaticRootPath = StaticRootPath

	return nil
}

func readSecuritySettings(iniFile *ini.File, cfg *Cfg) error {
	server := iniFile.Section("security")
	var err error

	cfg.DisableInitialAdminCreation = server.Key("disable_initial_admin_creation").MustBool(false)

	cfg.AdminUser, err = valueAsString(server, "admin_user", "test@godbledger.com")
	if err != nil {
		return err
	}

	cfg.AdminPassword, err = valueAsString(server, "admin_password", "password")
	if err != nil {
		return err
	}

	return nil
}

func readBackendSettings(iniFile *ini.File, cfg *Cfg) error {
	server := iniFile.Section("backend")
	var err error


	cfg.DatabaseType, err = valueAsString(server, "database_type", "sqlite3")
	if err != nil {
		return err
	}

	cfg.DatabaseLocation, err = valueAsString(server, "database_location", "~/.ledger/ledgerdata/ledger.db")
	if err != nil {
		return err
	}
  
	cfg.GoDBLedgerHost, err = valueAsString(server, "godbledger_host", "127.0.0.1")
	if err != nil {
		return err
	}

	cfg.GoDBLedgerPort, err = valueAsString(server, "godbledger_port", "50051")
	if err != nil {
		return err
	}

	cfg.GoDBLedgerCACert, err = valueAsString(server, "godbledger_cacert", "")
	if err != nil {
		return err
	}

	cfg.GoDBLedgerCert, err = valueAsString(server, "godbledger_cert", "")
	if err != nil {
		return err
	}

	cfg.GoDBLedgerKey, err = valueAsString(server, "godbledger_key", "")
	if err != nil {
		return err
	}

	return nil
}
