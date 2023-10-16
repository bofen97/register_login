package main

import (
	"io"
	"net/http"
	"os"
	"text/template"
)

func BuildRegisteOk(w http.ResponseWriter, email string) error {
	var code = `
	<html>
    <head>
        <meta charset="utf-8">
        <title>
            注册成功
        </title>
        <style >
            body{
                background-size: cover;
                background-attachment: fixed;
                background-color: #ffffff;
                background-size: 500;
            }
        </style>
        
    
    </head>
    <body bgcolor="black" >
        <div style="align-items: center; font-size: 45px"  >
            
            <p align="center">
                <img src="EasyTrackerLogo" title="Welcome" width="500" height="500" />
            </p>
            <p align="center"   style="color:rgb(57, 114, 206) ;" >Hi , {{.}} .</p>
            <p align="center"   style="color:rgba(56, 84, 129, 0.81) ;" >You have successfully registered your account. Please log in.</p>
            


        </div>
            
    </body>
</html>`
	tmpl, err := template.New("html").Parse(code)

	if err != nil {
		return err
	}
	err = tmpl.Execute(w, email)
	if err != nil {
		return err
	}
	return nil

}

func Logo(w http.ResponseWriter, r *http.Request) {
	f, err := os.OpenFile("./EasyTracker.jpg", os.O_RDONLY, 0755)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	io.Copy(w, f)
	defer f.Close()

}
