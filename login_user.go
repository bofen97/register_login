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
	Ut       *UserTable
	Session  *SessionTable
	Register *RegisterUser
}
type LoginUserData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// check passwd
func (login *LoginUser) CheckUserIsOK(email string, password string) (bool, error) {
	row := login.Ut.db.QueryRow("select COUNT(email) from userTable where email= ? and password = ? and created_at is not null ", email, password)
	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}
	if count == 1 {
		return true, nil
	}
	return false, nil
}

type LoginResp struct {
	Session string `json:"session"`
}

func (login *LoginUser) GenUserCurrentSession(email string) ([]byte, error) {

	row := login.Ut.db.QueryRow("select uid from userTable where email= ?  and created_at is not null", email)

	var loginresp LoginResp
	var uid int
	row.Scan(&uid)

	//gen session
	sessionStr := strconv.Itoa(uid) + time.Now().String() + email
	sessionByte := sha256.Sum256([]byte(sessionStr))
	loginresp.Session = fmt.Sprintf("%x", sessionByte)

	err := login.Session.InsertSessionAndUid(loginresp.Session, uid)
	if err != nil {
		return nil, err
	}
	data, err := json.MarshalIndent(loginresp, " ", " ")
	if err != nil {
		return nil, err
	}
	log.Printf("Insert session [%s] and id [%d] \n", loginresp.Session, uid)
	return data, nil

}
func (login *LoginUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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

		//check user exist ?
		ok, err := login.CheckUserIsOK(loginData.Email, loginData.Password)
		if err != nil {
			log.Printf("%v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !ok {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
		//is ok
		log.Printf("USER [%s]   PASSWD [%s] login \n", loginData.Email, loginData.Password)

		if r.Method == "POST" {
			sess, err := login.GenUserCurrentSession(loginData.Email)
			if err != nil {
				log.Printf("%v\n", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.Write(sess)
			// return session
			return

		} else if r.Method == "DELETE" {
			//delete account
			err := login.Ut.DeleteUserAccount(loginData.Email)
			if err != nil {
				log.Printf("%v\n", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			str, err := BuildEmailTemplateForDeleteAccount()
			if err != nil {
				log.Printf("%v\n", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			err = login.Register.SendMessage(loginData.Email, str)
			if err != nil {
				log.Printf("%v\n", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			fmt.Printf(" %s send delete account msg to User [%s]\n", "no-reply@easypapertracker.com", loginData.Email)
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)

}
