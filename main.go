package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var emailfrom string

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
	region := os.Getenv("region")
	if region == "" {
		region = "us-east-1"

	}
	emailfrom = os.Getenv("emailfrom")
	if emailfrom == "" {
		log.Fatal("emailfrom is none")
		return
	}

	log.Printf("region :%s ", region)
	register := new(RegisterUser)
	svc, err := NewSession(region)
	if err != nil {
		log.Fatal(err)
	}
	register.Svc = svc
	log.Print("aws svc ready .")

	register.Ut = new(UserTable)
	if err := register.Ut.Connect(sqlurl); err != nil {
		log.Fatal(err)
	}
	if err := register.Ut.CreateTable(); err != nil {
		log.Fatal(err)
	}
	log.Printf("%s\n", "register created !")

	login := new(LoginUser)
	login.Ut = new(UserTable)
	login.Session = new(SessionTable)
	if err := login.Ut.Connect(sqlurl); err != nil {
		log.Fatal(err)
	}
	log.Printf("%s\n", "login connect !")
	if err := login.Session.Connect(sqlurl); err != nil {
		log.Fatal(err)
	}
	if err := login.Session.CreateTable(); err != nil {
		log.Fatal(err)
	}
	log.Printf("%s\n", "Session created !")

	mux := http.NewServeMux()
	mux.Handle("/register", register)
	mux.Handle("/login", login)
	mux.HandleFunc("/valid", register.ValidLinkCheck)
	mux.HandleFunc("/EasyTrackerLogo", Logo)
	server := &http.Server{
		Addr:    serverPort,
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
