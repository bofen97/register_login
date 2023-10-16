package main

import (
	"database/sql"
	"log"
	"time"
)

type UserTable struct {
	db *sql.DB
}

func (ut *UserTable) Connect(url string) (err error) {

	ut.db, err = sql.Open("mysql", url)
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = ut.db.Ping()
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
func (ut *UserTable) CreateTable() error {

	query := `
	CREATE TABLE IF NOT EXISTS userTable (
		uid INT(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
		email TEXT NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME,
		validlink TEXT
	);`

	if _, err := ut.db.Exec(query); err != nil {
		return err
	}
	return nil
}

func (ut *UserTable) InsertEmailPasswdAndValidlink(email string, password string, validlink string) error {

	insertStr := `
		insert into userTable(
			email,
			password,
			validlink
		) values(?,?,?)`

	_, err := ut.db.Exec(insertStr, email, password, validlink)
	if err != nil {
		return err
	}
	return nil
}

func (ut *UserTable) CreateUserCommit(validLink string) error {

	updateStr := `
		update userTable set created_at = ? where validLink = ?
	`

	_, err := ut.db.Exec(updateStr, time.Now(), validLink)
	if err != nil {
		return err
	}
	return nil
}
func (ut *UserTable) GetCommitEmail(validLink string) (string, error) {
	var email string
	updateStr := `
		select email from userTable where created_at is not null and validLink = ?
	`

	row := ut.db.QueryRow(updateStr, validLink)
	err := row.Scan(&email)
	if err != nil {
		return "", err
	}
	return email, nil
}
