package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// register module
// POST register phonenumber
// generate MSGCODE to send && insert userTable
type RegisterUser struct {
	Ut *UserTable
}
type RegisterUserData struct {
	Phonenumber string `json:"phonenumber"`
	Type        int    `json:"type"` //type 0 , msgcode else passwd
	Password    string `json:"password"`
	MsgCode     int    `json:"msgcode"` //check msgcode is right ?

}

func (register *RegisterUser) CheckUserExist(phonenumber string) (bool, error) {
	//check user registed ?

	row := register.Ut.db.QueryRow(" select password from userTable where phonenumber = ? ", phonenumber)
	var passwd string
	if err := row.Scan(&passwd); err != nil {
		return false, nil
	}
	if passwd == "" {
		return false, nil
	}
	return true, nil

}

// TODO message send API,maybe ALIYUN
func (register *RegisterUser) GenMsgCodeAndSendMsgToUser(phonenumber string) error {
	var msgCode int
	for {
		msgCode = rand.Intn(10000)
		if msgCode >= 1000 {
			break
		}
	}

	row := register.Ut.db.QueryRow("select uid from userTable where phonenumber = ? ", phonenumber)
	var uid int
	err := row.Scan(&uid)
	if err != nil {
		err = register.Ut.InsertPhonenumberMsgCode(phonenumber, msgCode)
		if err != nil {
			log.Fatal(err)
			return err
		}
	} else {
		register.Ut.db.Exec("update userTable set msgcode = ? , msgcode_at = ? where phonenumber = ? ", msgCode, time.Now(), phonenumber)

	}
	phonenumberInt, err := strconv.Atoi(phonenumber)
	if err != nil {
		log.Fatal(err)
		return err
	}

	//send user
	fmt.Printf("Send MsgCode :[%d] to User [%d]\n", msgCode, phonenumberInt)
	return nil

}

// check msgcode before 5 mins ago?
// check msgcode is exits？
// and msgcode is right ?
func (register *RegisterUser) CheckMsgCodeIsNotExistAndBefore5MinsAgo(phonenumber string, msgcode int) bool {
	row := register.Ut.db.QueryRow("select msgcode_at , msgcode from userTable where phonenumber= ? ", phonenumber)
	var msgcode_at time.Time
	var msgcodeTmp int
	err := row.Scan(&msgcode_at, &msgcodeTmp)
	if err != nil {
		return true
	}

	if time.Since(msgcode_at) > 5*time.Minute {
		return true
	}
	if msgcodeTmp != msgcode {
		return true
	}

	return false
}
func (register *RegisterUser) RegistePhonenumerAndPassWord(phonenumber string, password string) error {

	err := register.Ut.UpdatePhonenumberPassword(phonenumber, password)
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Printf("Registe Phonenumber [%s]  And Password [%s] \n", phonenumber, password)

	return nil

}

func (register *RegisterUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if r.Header.Get("Content-Type") == "application/json" {

			data, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			var registerData RegisterUserData
			err = json.Unmarshal(data, &registerData)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if registerData.Type == 0 {

				//got phonenumber&code type
				// check user exist ?

				flag, _ := register.CheckUserExist(registerData.Phonenumber)
				if flag {
					w.WriteHeader(http.StatusCreated)
					fmt.Fprintf(w, "User [%s] exist \n", registerData.Phonenumber)
					return
				}
				err = register.GenMsgCodeAndSendMsgToUser(registerData.Phonenumber)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

			} else {
				// check msgcode is not exits？&& msgcode time.since(msgcode_at) > 5 mins?
				// and check msgcode is right ?
				flag := register.CheckMsgCodeIsNotExistAndBefore5MinsAgo(registerData.Phonenumber, registerData.MsgCode)
				if flag {
					w.WriteHeader(http.StatusRequestTimeout)
					fmt.Fprint(w, "Msgcode timeout , plase repeat! \n")
					return
				}

				//got phonenumber& password
				err = register.RegistePhonenumerAndPassWord(registerData.Phonenumber, registerData.Password)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}

			w.WriteHeader(http.StatusOK)
			return

		}

	}

	w.WriteHeader(http.StatusBadRequest)

}
