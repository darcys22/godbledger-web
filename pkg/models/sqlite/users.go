package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func New(path string) UserModel {
	database, _ := sql.Open("sqlite3", path)

	statement, _ := database.Prepare(`
	CREATE TABLE IF NOT EXISTS users (
	id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT, name VARCHAR(255) NOT NULL,
	email VARCHAR(255) UNIQUE NOT NULL,
	hashed_password CHAR(60) NOT NULL,
	created DATETIME NOT NULL,
	active BOOLEAN NOT NULL DEFAULT TRUE
	);
	`)
	statement.Exec()

	//TODO insert the default user if in database config

	return UserModel{DB: database}
}

func (m *UserModel) Insert(name, email, password string) error {
	// Create a bcrypt hash of the plain-text password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES(?, ?, ?, UTC_TIMESTAMP())`
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
