package setting

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	ini "gopkg.in/ini.v1"
)

var log = logrus.WithField("prefix", "setting")

type Scheme string

const (
	HTTP              Scheme = "http"
	HTTPS             Scheme = "https"
	HTTP2             Scheme = "h2"
	SOCKET            Scheme = "socket"
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
	AppUrl           string
	AppSubUrl        string
	ServeFromSubPath bool
	InstanceName     string

	// build
	BuildVersion    string
	BuildCommit     string
	BuildBranch     string
	BuildStamp      int64
	IsEnterprise    bool
	ApplicationName string

	// packaging
	Packaging = "unknown"

	// Paths
	HomePath       string
	PluginsPath    string
	CustomInitPath = "conf/custom.ini"

	// Http server options
	Protocol           Scheme
	Domain             string
	HttpAddr, HttpPort string
	CertFile, KeyFile  string
	SocketPath         string
	RouterLogging      bool
	DataProxyLogging   bool
	DataProxyTimeout   int
	StaticRootPath     string
	EnableGzip         bool
	EnforceDomain      bool

	// Global setting objects.
	Raw          *ini.File
	IsWindows    bool

	// for logging purposes
	configFiles                  []string
	appliedCommandLineProperties []string
	appliedEnvOverrides          []string

	ReportingEnabled   bool
	CheckForUpdates    bool
	GoogleAnalyticsId  string
	GoogleTagManagerId string

	// LDAP
	LDAPEnabled           bool
	LDAPConfigFile        string
	LDAPSyncCron          string
	LDAPAllowSignup       bool
	LDAPActiveSyncEnabled bool

	// Alerting
	AlertingEnabled            bool
	ExecuteAlerts              bool
	AlertingRenderLimit        int
	AlertingErrorOrTimeout     string
	AlertingNoDataOrNullValues string

	AlertingEvaluationTimeout   time.Duration
	AlertingNotificationTimeout time.Duration
	AlertingMaxAttempts         int
	AlertingMinInterval         int64

	// Explore UI
	ExploreEnabled bool

	// GoDBLedger.NET URL
	GodbledgerComUrl string

	ImageUploadProvider string
)

// TODO move all global vars to this struct
type Cfg struct {
	Raw *ini.File

	// HTTP Server Settings
	AppUrl           string
	AppSubUrl        string
	ServeFromSubPath bool
	StaticRootPath   string
	Protocol         Scheme

	// build
	BuildVersion string
	BuildCommit  string
	BuildBranch  string
	BuildStamp   int64
	IsEnterprise bool

	// packaging
	Packaging string

	// Paths
	ProvisioningPath   string
	DataPath           string
	BundledPluginsPath string

	// Dataproxy
	SendUserHeader bool

	// DistributedCache
	RemoteCacheOptions *RemoteCacheOptions
}

type CommandLineArgs struct {
	Config   string
	HomePath string
	Args     []string
}

func init() {
	IsWindows = runtime.GOOS == "windows"
}

func parseAppUrlAndSubUrl(section *ini.Section) (string, string, error) {
	appUrl, err := valueAsString(section, "root_url", "http://localhost:3000/")
	if err != nil {
		return "", "", err
	}
	if appUrl[len(appUrl)-1] != '/' {
		appUrl += "/"
	}

	// Check if has app suburl.
	url, err := url.Parse(appUrl)
	if err != nil {
		log.Fatalf("Invalid root_url(%s): %s", appUrl, err)
	}
	appSubUrl := strings.TrimSuffix(url.Path, "/")

	return appUrl, appSubUrl, nil
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

	// update data path and logging config
	dataPath, err := valueAsString(parsedFile.Section("paths"), "data", "")
	if err != nil {
		return nil, err
	}
	cfg.DataPath = makeAbsolute(dataPath, HomePath)

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
	cfg.Packaging = Packaging

	ApplicationName = APP_NAME

	Env, err = valueAsString(iniFile.Section(""), "app_mode", "development")
	if err != nil {
		return err
	}
	InstanceName, err = valueAsString(iniFile.Section(""), "instance_name", "unknown_instance_name")
	if err != nil {
		return err
	}
	plugins, err := valueAsString(iniFile.Section("paths"), "plugins", "")
	if err != nil {
		return err
	}
	PluginsPath = makeAbsolute(plugins, HomePath)
	cfg.BundledPluginsPath = makeAbsolute("plugins-bundled", HomePath)
	provisioning, err := valueAsString(iniFile.Section("paths"), "provisioning", "")
	if err != nil {
		return err
	}
	cfg.ProvisioningPath = makeAbsolute(provisioning, HomePath)
	if err := readServerSettings(iniFile, cfg); err != nil {
		return err
	}

	// read data proxy settings
	dataproxy := iniFile.Section("dataproxy")
	DataProxyLogging = dataproxy.Key("logging").MustBool(false)
	DataProxyTimeout = dataproxy.Key("timeout").MustInt(30)
	cfg.SendUserHeader = dataproxy.Key("send_user_header").MustBool(false)


	analytics := iniFile.Section("analytics")
	ReportingEnabled = analytics.Key("reporting_enabled").MustBool(true)
	CheckForUpdates = analytics.Key("check_for_updates").MustBool(true)
	GoogleAnalyticsId = analytics.Key("google_analytics_ua_id").String()
	GoogleTagManagerId = analytics.Key("google_tag_manager_id").String()

	if err := readAlertingSettings(iniFile); err != nil {
		return err
	}

	explore := iniFile.Section("explore")
	ExploreEnabled = explore.Key("enabled").MustBool(true)

	cfg.Protocol = Protocol

	cfg.readLDAPConfig()
	cfg.readSessionConfig()

	// check old key  name
	GodbledgerComUrl, err = valueAsString(iniFile.Section("godbledger_net"), "url", "")
	if err != nil {
		return err
	}
	if GodbledgerComUrl == "" {
		GodbledgerComUrl, err = valueAsString(iniFile.Section("godbledger_com"), "url", "https://godbledger.com")
		if err != nil {
			return err
		}
	}

	imageUploadingSection := iniFile.Section("external_image_storage")
	ImageUploadProvider, err = valueAsString(imageUploadingSection, "provider", "")
	if err != nil {
		return err
	}

	cacheServer := iniFile.Section("remote_cache")
	dbName, err := valueAsString(cacheServer, "type", "database")
	if err != nil {
		return err
	}
	connStr, err := valueAsString(cacheServer, "connstr", "")
	if err != nil {
		return err
	}
	cfg.RemoteCacheOptions = &RemoteCacheOptions{
		Name:    dbName,
		ConnStr: connStr,
	}

	return nil
}

func valueAsString(section *ini.Section, keyName string, defaultValue string) (string, error) {
	return section.Key(keyName).MustString(defaultValue), nil
}

type RemoteCacheOptions struct {
	Name    string
	ConnStr string
}

func (cfg *Cfg) readLDAPConfig() {
	ldapSec := cfg.Raw.Section("auth.ldap")
	LDAPConfigFile = ldapSec.Key("config_file").String()
	LDAPSyncCron = ldapSec.Key("sync_cron").String()
	LDAPEnabled = ldapSec.Key("enabled").MustBool(false)
	LDAPActiveSyncEnabled = ldapSec.Key("active_sync_enabled").MustBool(false)
	LDAPAllowSignup = ldapSec.Key("allow_sign_up").MustBool(true)
}

func (cfg *Cfg) readSessionConfig() {
	sec, _ := cfg.Raw.GetSection("session")

	if sec != nil {
		log.Warn(
			"[Removed] Session setting was removed in v6.2, use remote_cache option instead",
		)
	}
}

func readAlertingSettings(iniFile *ini.File) error {
	alerting := iniFile.Section("alerting")
	AlertingEnabled = alerting.Key("enabled").MustBool(true)
	ExecuteAlerts = alerting.Key("execute_alerts").MustBool(true)
	AlertingRenderLimit = alerting.Key("concurrent_render_limit").MustInt(5)
	var err error
	AlertingErrorOrTimeout, err = valueAsString(alerting, "error_or_timeout", "alerting")
	if err != nil {
		return err
	}
	AlertingNoDataOrNullValues, err = valueAsString(alerting, "nodata_or_nullvalues", "no_data")
	if err != nil {
		return err
	}

	evaluationTimeoutSeconds := alerting.Key("evaluation_timeout_seconds").MustInt64(30)
	AlertingEvaluationTimeout = time.Second * time.Duration(evaluationTimeoutSeconds)
	notificationTimeoutSeconds := alerting.Key("notification_timeout_seconds").MustInt64(30)
	AlertingNotificationTimeout = time.Second * time.Duration(notificationTimeoutSeconds)
	AlertingMaxAttempts = alerting.Key("max_attempts").MustInt(3)
	AlertingMinInterval = alerting.Key("min_interval_seconds").MustInt64(1)

	return nil
}

func readServerSettings(iniFile *ini.File, cfg *Cfg) error {
	server := iniFile.Section("server")
	var err error
	AppUrl, AppSubUrl, err = parseAppUrlAndSubUrl(server)
	if err != nil {
		return err
	}
	ServeFromSubPath = server.Key("serve_from_sub_path").MustBool(false)

	cfg.AppUrl = AppUrl
	cfg.AppSubUrl = AppSubUrl
	cfg.ServeFromSubPath = ServeFromSubPath

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
	if protocolStr == "h2" {
		Protocol = HTTP2
		CertFile = server.Key("cert_file").String()
		KeyFile = server.Key("cert_key").String()
	}
	if protocolStr == "socket" {
		Protocol = SOCKET
		SocketPath = server.Key("socket").String()
	}

	Domain, err = valueAsString(server, "domain", "localhost")
	if err != nil {
		return err
	}
	HttpAddr, err = valueAsString(server, "http_addr", DEFAULT_HTTP_ADDR)
	if err != nil {
		return err
	}
	HttpPort, err = valueAsString(server, "http_port", "3000")
	if err != nil {
		return err
	}
	RouterLogging = server.Key("router_logging").MustBool(false)

	EnableGzip = server.Key("enable_gzip").MustBool(false)
	EnforceDomain = server.Key("enforce_domain").MustBool(false)
	staticRoot, err := valueAsString(server, "static_root_path", "")
	if err != nil {
		return err
	}
	StaticRootPath = makeAbsolute(staticRoot, HomePath)
	cfg.StaticRootPath = StaticRootPath

	return nil
}
