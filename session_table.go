package main

import (
	"database/sql"
	"log"
	"time"
)

type SessionTable struct {
	db *sql.DB
}

func (sess *SessionTable) Connect(url string) (err error) {

	sess.db, err = sql.Open("mysql", url)
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = sess.db.Ping()
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
func (sess *SessionTable) CreateTable() error {

	query := `
	CREATE TABLE IF NOT EXISTS sessionTable (
		session TEXT NOT NULL,
		uid INT(11) NOT NULL ,
		created_at DATETIME NOT NULL
	);`

	if _, err := sess.db.Exec(query); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (sess *SessionTable) InsertSessionAndUid(session string, uid int) error {

	query := ` insert into sessionTable(session,uid,created_at) values(?,?,?)`

	if _, err := sess.db.Exec(query, session, uid, time.Now()); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
