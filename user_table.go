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
		phonenumber TEXT NOT NULL,
		password TEXT ,
		created_at DATETIME NOT NULL,
		msgcode INT(11) ,
		msgcode_at DATETIME
	);`

	if _, err := ut.db.Exec(query); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (ut *UserTable) InsertPhonenumberMsgCode(phonenumber string, msgcode int) error {

	insertStr := `
		insert into userTable(
			phonenumber,
			msgcode,
			created_at,
			msgcode_at
		) values(?,?,?,?)`

	_, err := ut.db.Exec(insertStr, phonenumber, msgcode, time.Now(), time.Now())
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
func (ut *UserTable) UpdatePhonenumberPassword(phonenumber string, password string) error {

	updateStr := `
		update userTable set password = ? where phonenumber = ?
	`

	_, err := ut.db.Exec(updateStr, password, phonenumber)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (ut *UserTable) UpdatePhonenumberMsgCode(phonenumber string, msgcode int) error {
	/*
		UPDATE table_name
		SET column1 = value1, column2 = value2, ...
		WHERE condition;

	*/
	updateStr := `
		update userTable set msgcode = ? , msgcode_at = ? where phonenumber= ?
	`
	_, err := ut.db.Exec(updateStr, msgcode, time.Now(), phonenumber)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (ut *UserTable) DeleteUserInfo(phonenumber string) error {

	//DELETE FROM table_name
	//WHERE condition;
	deleteStr := `
		delete from userTable where phonenumber = ?
	`
	_, err := ut.db.Exec(deleteStr, phonenumber)
	if err != nil {
		log.Fatal(err)
	}
	return nil

}
