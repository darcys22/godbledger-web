package sqlite

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"

	"github.com/sirupsen/logrus"

	"github.com/darcys22/godbledger-web/backend/models"
)

var log = logrus.WithField("prefix", "SqliteUsers")

var ErrNoRows = errors.New("sql: no rows in result set")

type UserModel struct {
	DB *sql.DB
}

func New(path string) UserModel {
	database, _ := sql.Open("sqlite3", path)

	log.Info("Users database at path: ", path)

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
	usersdb := UserModel{DB: database}

	//TODO this should be conditionally run from config
	defaultUserID := 0
	err = database.QueryRow(`SELECT id FROM users WHERE email = ? LIMIT 1`, "test@godbledger.com").Scan(&defaultUserID)
	if err != nil {
		if err.Error() == ErrNoRows.Error() {
			log.Info("Inserting default user into users table")
			err = usersdb.Insert("defaultuser", "test@godbledger.com", "password")
			if err != nil {
				log.Error("Error in adding default user: ", err)
			}

		} else {
			log.Error("Error in searching for default user: ", err)
		}
	}

	return usersdb
}

func (m *UserModel) Insert(name, email, password string) error {
	// Create a bcrypt hash of the plain-text password.
	log.Infof("Inserting user into users table, Name: %s Email: %s", name, email)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES(?, ?, ?, datetime('now'))`
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		return err
	}
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
func (m *UserModel) Get(name string) (*models.User, error) {
	var user models.User 
	stmt := "SELECT * FROM users WHERE email = ? AND active = TRUE"
	row := m.DB.QueryRow(stmt, name)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.HashedPassword, &user.Created, &user.Active, &user.Currency, &user.DateLocale, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil 
}

func (m *UserModel) Save(user *models.User) (error) {
	stmt := "UPDATE users SET currency = ?, locale = ? WHERE email = ? AND active = TRUE"
	_, err := m.DB.Exec(stmt, user.Currency, user.DateLocale, user.Email)
	log.Info("Users Saved: ", user)
	if err != nil {
		return err
	}
	return nil
}
