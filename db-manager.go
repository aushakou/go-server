package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"database/sql"
	"errors"
	"log"
)

const database = "a1"
const userTable = "users"
const staffTable = "staff_data"
const accessTable = "system_access_data"

type DBManager struct {
	Username string
	Password string
	Host     string

	db          *sql.DB
	isConnected bool
}

func (m *DBManager) Connect() error {
	db, err := sql.Open("mysql", m.Username+":"+m.Password+"@tcp("+m.Host+")/"+database)

	if err == nil {
		m.db = db
	}

	return err
}

func (m *DBManager) Authenticate(user User) error {
	invalidError := errors.New("Invalid credentials provided")

	dbUser := user

	row := m.db.QueryRow("select * from "+userTable+" where username = ?", user.Username)
	if err := row.Scan(&dbUser.Username, &dbUser.Password, &dbUser.Salt); err != nil {
		return invalidError
	}

	user.Salt = dbUser.Salt

	err := m.generateSaltedPassword(&user)

	if err != nil {
		return err
	} else if !bytes.Equal(user.Password, dbUser.Password) {
		return invalidError
	} else {
		return nil
	}
}

func (m *DBManager) Register(user User) error {

	// make sure that salt is empty
	user.Salt = nil
	err := m.generateSaltedPassword(&user)

	if err != nil {
		return err
	}

	log.Println(user)
	p := user.Password
	s := user.Salt
	u := user.Username

	res, err := m.db.Exec("insert into "+userTable+" values(?,?,?)", u, p, s)

	if err != nil {
		return err
	}

	if affected, err := res.RowsAffected(); err != nil {
		return err
	} else if affected <= 0 {
		return errors.New("Failed to insert records")
	} else {
		return nil
	}
}

// returns error or nil
func (m *DBManager) generateSaltedPassword(user *User) error {
	// generate salt
	tmp := make([]byte, 10)
	if _, err := rand.Read(tmp); err != nil {
		return errors.New("Unable to generate salt")

	}

	var salt []byte

	// if salt is not present, generate it.
	// if present, just use it
	if user.Salt == nil || len(user.Salt) == 0 {
		a := sha256.Sum256(tmp)
		salt = a[:]
	} else {
		salt = user.Salt
	}

	hasher := sha512.New()
	hasher.Write(append(user.Password, salt...))
	hashedPassword := hasher.Sum(nil)

	user.Password = hashedPassword
	user.Salt = salt

	return nil
}

func (m *DBManager) CloseConnection() error {
	return m.db.Close()
}
