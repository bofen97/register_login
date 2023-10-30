package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	sqlurl := os.Getenv("sqlurl")
	if sqlurl == "" {
		log.Fatal("sqlurl is none")
		return
	}
	serverPort := os.Getenv("serverport")
	if serverPort == "" {
		log.Fatal("serverPort is none")
		return
	}
	emailpasswd := os.Getenv("emailpasswd")
	if emailpasswd == "" {
		log.Fatal("emailpasswd is none")
		return
	}

	register := new(RegisterUser)
	register.InitialGoMail(emailpasswd)
	register.Ut = new(UserTable)
	if err := register.Ut.Connect(sqlurl); err != nil {
		log.Fatal(err)
	}
	if err := register.Ut.CreateTable(); err != nil {
		log.Fatal(err)
	}
	log.Printf("%s\n", "register created !")

	login := new(LoginUser)
	login.Register = register

	login.Ut = new(UserTable)
	if err := login.Ut.Connect(sqlurl); err != nil {
		log.Fatal(err)
	}

	login.Session = new(SessionTable)
	log.Printf("%s\n", "login connect !")
	if err := login.Session.Connect(sqlurl); err != nil {
		log.Fatal(err)
	}
	if err := login.Session.CreateTable(); err != nil {
		log.Fatal(err)
	}
	log.Printf("%s\n", "Session created !")

	subj := new(SubjectServer)
	subj.Session = new(SessionTable)
	if err := subj.Session.Connect(sqlurl); err != nil {
		log.Fatal(err)
	}

	subj.St = new(SubjectTable)
	if err := subj.St.Connect(sqlurl); err != nil {
		log.Fatal(err)
	}
	log.Printf("%s\n", "subject  connected .. !")

	if err := subj.St.CreateTable(); err != nil {
		log.Fatal(err)
	}
	log.Printf("%s\n", "subject  table createTable .. !")

	mux := http.NewServeMux()
	mux.Handle("/register", register)
	mux.Handle("/login", login)
	mux.HandleFunc("/valid", register.ValidLinkCheck)
	mux.HandleFunc("/EasyTrackerLogo", Logo)
	mux.Handle("/subject_server", subj)
	server := &http.Server{
		Addr:    serverPort,
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
