package sqlite

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"

	"github.com/sirupsen/logrus"

	"github.com/darcys22/godbledger-web/backend/models"
	"github.com/darcys22/godbledger-web/backend/setting"
)

var log = logrus.WithField("prefix", "SqliteUsers")

var ErrNoRows = errors.New("sql: no rows in result set")

type UserModel struct {
	DB *sql.DB
	Cfg *setting.Cfg
}

func New(path string, cfg *setting.Cfg) UserModel {
	database, _ := sql.Open("sqlite3", path)

	statement, err := database.Prepare(`
  CREATE TABLE IF NOT EXISTS users (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
		name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
		currency VARCHAR(10) NOT NULL DEFAULT "USD",
		locale VARCHAR(5) NOT NULL DEFAULT "en-AU",
		role VARCHAR(10) NOT NULL DEFAULT "standard"
  );
	`)
	if err != nil {
		log.Error("Error in prepare statement: ", err)
	}
	statement.Exec()
	usersdb := UserModel{DB: database, Cfg: cfg}

	if !cfg.DisableInitialAdminCreation {
		err = usersdb.CreateDefaultUser(cfg.AdminUser, cfg.AdminPassword)
		if err != nil {
			log.Error("Error creating default user: ", err)
		}
	}

	return usersdb
}

func (m *UserModel) CreateDefaultUser(email, password string) error {
	defaultUserID := 0
	err := m.DB.QueryRow(`SELECT id FROM users WHERE email = ? LIMIT 1`, email).Scan(&defaultUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			if err.Error() == ErrNoRows.Error() {
				err = m.Insert("defaultuser", email , password)
				if err != nil {
					return err
				}
				defaultUser, err := m.Get(email)
				if err != nil {
					return err
				}
				err = m.ChangePermissions(defaultUser, "admin")
				if err != nil {
					return err
				}

			} else {
				return err
			}
			log.Info("Created default user")
		}
	}
	return nil
}

func (m *UserModel) Insert(name, email, password string) error {
	// Create a bcrypt hash of the plain-text password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES(?, ?, ?, datetime('now'))`
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		return err
	}
	log.Infof("Inserted user into users table, Name: %s Email: %s", name, email)
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {

	// Retrieve the id and hashed password associated with the given email. If no
	// matching email exists, or the user is not active, we return the
	// ErrInvalidCredentials error.
	var id int
	var hashedPassword []byte
	stmt := "SELECT id, hashed_password FROM users WHERE email = ? AND active = TRUE"
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return 0, err
	}
	// Check whether the hashed password and plain-text password provided match. // If they don't, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return 0, err
	}
	// Otherwise, the password is correct. Return the user ID.
	return id, nil
}

func (m *UserModel) NewUser(email, password string) (int, error) {

	// Retrieve the id and hashed password associated with the given email. If no
	// matching email exists, or the user is not active, we return the
	// ErrInvalidCredentials error.
	var id int
	var hashedPassword []byte
	stmt := "INSERT into uid, hashed_password FROM users WHERE email = ? AND active = TRUE"
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return 0, err
	}
	// Check whether the hashed password and plain-text password provided match. // If they don't, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return 0, err
	}
	// Otherwise, the password is correct. Return the user ID.
	return id, nil
}

// We'll use the Get method to fetch details for a specific user based on their email/name
func (m *UserModel) Get(email string) (*models.User, error) {
	var user models.User 
	stmt := "SELECT * FROM users WHERE email = ? AND active = TRUE"
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.HashedPassword, &user.Created, &user.Active, &user.Currency, &user.DateLocale, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil 
}

func (m *UserModel) Save(user *models.User) (error) {
	stmt := "UPDATE users SET currency = ?, locale = ? WHERE email = ? AND active = TRUE"
	_, err := m.DB.Exec(stmt, user.Currency, user.DateLocale, user.Email)
	if err != nil {
		return err
	}
	log.Info("Users Saved: ", user)
	return nil
}

func (m *UserModel) ChangePermissions(user *models.User, role string) (error) {
	stmt := "UPDATE users SET role = ? WHERE email = ? AND active = TRUE"
	_, err := m.DB.Exec(stmt, role, user.Email)
	log.Info("Users Permissions Updated: ", user)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) ChangePassword(user *models.User, password string) (error) {
	// Create a bcrypt hash of the plain-text password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `UPDATE users SET hashed_password = ? WHERE email = ? AND active = TRUE`
	_, err = m.DB.Exec(stmt, hashedPassword, user.Email)
	if err != nil {
		return err
	}
	log.Infof("Updated user password, Name: %s Email: %s", user.Name , user.Email)
	return nil
}
