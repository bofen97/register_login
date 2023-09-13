package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type LoginUser struct {
	Ut      *UserTable
	Session *SessionTable
}
type LoginUserData struct {
	Phonenumber string `json:"phonenumber"`
	Password    string `json:"password"`
}

// check passwd
func (login *LoginUser) CheckPasswordIsWrong(phonenumber string, password string) bool {
	row := login.Ut.db.QueryRow("select password from userTable where phonenumber= ? ", phonenumber)
	var passwordTmp string
	err := row.Scan(&passwordTmp)

	if err != nil {
		log.Printf("Got User [%s] Login Request && Passwd [%s]  , But user not exist  \n", phonenumber, password)

		return true
	}
	if passwordTmp != password {
		log.Printf("Got User [%s] Login Request && Passwd [%s]  , But Not Equ [%s] \n", phonenumber, password, passwordTmp)
		return true
	}
	return false
}

type LoginResp struct {
	Session string `json:"session"`
}

func (login *LoginUser) ResponseSession(w http.ResponseWriter, phonenumber string) error {

	row := login.Ut.db.QueryRow("select uid from userTable where phonenumber= ?  ", phonenumber)

	var loginresp LoginResp
	var uidTmp int
	row.Scan(&uidTmp)

	//gen session
	sessionStr := strconv.Itoa(uidTmp) + time.Now().String() + phonenumber
	sessionByte := sha256.Sum256([]byte(sessionStr))
	loginresp.Session = fmt.Sprintf("%x", sessionByte)

	err := login.Session.InsertSessionAndUid(loginresp.Session, uidTmp)
	if err != nil {
		log.Fatal(err)
		return err
	}
	data, err := json.MarshalIndent(loginresp, " ", " ")
	if err != nil {
		log.Fatal(err)
		return err
	}
	w.Write(data)

	log.Printf("Insert session [%s] and id [%d] \n", loginresp.Session, uidTmp)
	return nil

}
func (login *LoginUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		if r.Header.Get("Content-Type") == "application/json" {

			data, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			var loginData LoginUserData
			err = json.Unmarshal(data, &loginData)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			//check passwd && user exist ?

			flag := login.CheckPasswordIsWrong(loginData.Phonenumber, loginData.Password)
			if flag {
				w.WriteHeader(http.StatusNonAuthoritativeInfo)
				return
			}
			log.Printf("User [%s] logined  \n", loginData.Phonenumber)

			// return session
			err = login.ResponseSession(w, loginData.Phonenumber)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			//w.WriteHeader(http.StatusOK)
			return
		}

	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
