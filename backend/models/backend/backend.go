package backend

import (
  "database/sql"
  "fmt"
  "regexp"
  "strings"
  "errors"

	"github.com/darcys22/godbledger-web/backend/setting"

	"github.com/sirupsen/logrus"

  mysql "github.com/go-sql-driver/mysql"
  _ "github.com/mattn/go-sqlite3"

)

var log = logrus.WithField("prefix", "godbledger-backend")
var dsnRegex = regexp.MustCompile(`\:(.+?)\@`)

type GodbledgerBackendModel struct {
	DB *sql.DB
	Cfg *setting.Cfg
}

var (
	backend GodbledgerBackendModel
)

func InitBackendConnection() error {
  cfg := setting.GetConfig()
	switch strings.ToLower(cfg.DatabaseType) {
	case "sqlite3", "memorydb":
		log.Debug("Using Sqlite3")
		mode := "rwc"
		if strings.ToLower(cfg.DatabaseType) == "memorydb" {
			log.Debug("In Memory only Mode")
			mode = "memory"
		}
    datafile := fmt.Sprintf("%s?_foreign_keys=true&parseTime=true&mode=%s", cfg.DatabaseLocation, mode)
    if mode == "memory" {
      datafile = fmt.Sprintf("%s?_foreign_keys=true&parseTime=true&mode=%s", ":memory:", mode)
    }
    log.WithField("datafile", datafile).Debug("Opening SQLite3 Datafile")
    SqliteDB, err := sql.Open("sqlite3", datafile)
    if err != nil {
      return err
    }
    backend = GodbledgerBackendModel{DB: SqliteDB, Cfg: cfg}

	case "mysql":
		log.Debug("Using MySQL")
    validatedString, err := ValidateConnectionString(cfg.DatabaseLocation)
    if err != nil {
      log.Fatal(err.Error())
      return err
    }
    MySQLDB, err := sql.Open("mysql", validatedString)
    if err != nil {
      log.Fatal(err.Error())
      return err
    }
    backend = GodbledgerBackendModel{DB: MySQLDB, Cfg: cfg}
	default:
		log.Fatal(cfg.DatabaseType)
		log.Fatal("No implementation available for that database.")
	}


	log.Debug("Initialised database configuration")

  return nil
}

func ValidateConnectionString(dsn string) (string, error) {

	if dsn == "" {
		return "", errors.New("Connection string not provided")
	}

	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		log.Warnf("Connection string could not be parsed: %s", err.Error())
		return "", err
	}
	log.Debugf("DB_ADDR := %s", cfg.Addr)
	log.Debugf("DB_NET := %s", cfg.Net)
	log.Debugf("DB_DBNAME := %s", cfg.DBName)
	log.Debugf("DB_USER := %s", cfg.User)
	log.Debugf("PARAMS := %v", cfg.Params)
	if !cfg.ParseTime {
		cfg.ParseTime = true
	}
	charset, ok := cfg.Params["charset"]
	if !(ok && charset == "utf8") {
		if cfg.Params == nil {
			cfg.Params = make(map[string]string)
		}
		cfg.Params["charset"] = "utf8"
	}

	log.Debugf("ParseTime := %v", cfg.ParseTime)
	log.Debugf("Charset := %s", cfg.Params["charset"])

	dsnString := cfg.FormatDSN()
	log.Debugf("DSN := %s", redactPassword(dsnString))

	return dsnString, nil
}

func redactPassword(rawDSNString string) string {
	cleanedDSNString := dsnRegex.ReplaceAll([]byte(rawDSNString), []byte(":**REDACTED**@"))
	return string(cleanedDSNString)
}

func GetConnection() *sql.DB {
  return backend.DB
}
