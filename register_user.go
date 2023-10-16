package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ses"
)

// register module
// POST register email
// generate link  to send && insert userTable
type RegisterUser struct {
	Ut  *UserTable
	Svc *ses.SES
}
type RegisterUserData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (register *RegisterUser) CheckUserExist(email string) (bool, error) {
	//check user registed ?
	//should be 0;
	row := register.Ut.db.QueryRow(" select COUNT(email) from userTable where email = ? and created_at is not null ", email)
	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}
	if count != 0 {
		return true, nil
	}
	return false, nil

}

func (register *RegisterUser) GenValidLinkSendToUser(email string, password string) error {
	salt := time.Now().String()

	validLink := fmt.Sprintf("%x", md5.Sum([]byte(email+salt)))
	log.Printf("gen valid link %s \n", validLink)

	//insert into userTable (validlink)
	err := register.Ut.InsertEmailPasswdAndValidlink(email, password, validLink)
	if err != nil {
		return err
	}

	//send valid link to user

	str, err := BuildEmailTemplate(fmt.Sprintf("http://localhost:8080/valid?hash=%s", validLink))
	if err != nil {
		return err
	}

	register.SendMessage(email, emailfrom, str)
	fmt.Printf(" %s send ValidLink :[%s] to User [%s]\n", emailfrom, validLink, email)
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
			log.Printf("get registe request %v\n", registerData)
			// check user email exist ?
			isExist, err := register.CheckUserExist(registerData.Email)
			if err != nil {
				log.Print(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if isExist {
				w.WriteHeader(http.StatusCreated)
				return
			}
			//user not exist.
			// generate link send to user;
			err = register.GenValidLinkSendToUser(registerData.Email, registerData.Password)
			if err != nil {
				log.Printf("%v\n", err)
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}

			w.WriteHeader(http.StatusOK)
			return
		}

	}

	w.WriteHeader(http.StatusBadRequest)

}

// precommit is validlink exist and created_at is null , count  should be 1
func (register *RegisterUser) ValidLinkIsPreCommit(validLink string) (bool, error) {

	row := register.Ut.db.QueryRow(" select COUNT(email) from userTable where validlink = ? and created_at is null ", validLink)
	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}
	if count == 1 {
		return true, nil
	}
	return false, nil
}

func (register *RegisterUser) ValidLinkCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		validlink := r.URL.Query().Get("hash")
		log.Printf("GET VALIDLINK HASH %s\n", validlink)
		//validlink exist ?
		preCommit, err := register.ValidLinkIsPreCommit(validlink)
		if err != nil {
			log.Printf("%v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !preCommit {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// precommit status

		err = register.Ut.CreateUserCommit(validlink)
		if err != nil {
			log.Printf("%v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		//get mail
		email, err := register.Ut.GetCommitEmail(validlink)
		if err != nil {
			log.Printf("%v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		BuildRegisteOk(w, email)
		return
	}

	w.WriteHeader(http.StatusBadRequest)

}

func (register *RegisterUser) SendMessage(to, from string, registMessage string) {
	input := BuildMessage(to, from, registMessage)
	result, err := register.Svc.SendEmail(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			case ses.ErrCodeConfigurationSetSendingPausedException:
				fmt.Println(ses.ErrCodeConfigurationSetSendingPausedException, aerr.Error())
			case ses.ErrCodeAccountSendingPausedException:
				fmt.Println(ses.ErrCodeAccountSendingPausedException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
	log.Println(result)

}
