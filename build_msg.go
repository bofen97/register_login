package main

import (
	"bytes"
	"text/template"
)

// region = "us-east-1"

func BuildEmailTemplate(validlink string) (string, error) {
	code := "Hello, this is the registration email for EasyPaperTracker, please <a class=\"ulink\" href=\"{{.}}\" target=\"_blank\">click the link to complete registration</a>."
	tmpl, err := template.New("html").Parse(code)

	if err != nil {
		return "", err
	}

	buffer := bytes.NewBuffer([]byte{})
	err = tmpl.Execute(buffer, validlink)

	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}
func BuildEmailTemplateForDeleteAccount() (string, error) {
	code := "We will delete all data we save in accordance with the \"User Agreement\" and \"Privacy Policy\", including user email addresses, passwords, and backup data. Thank you."
	tmpl, err := template.New("html").Parse(code)

	if err != nil {
		return "", err
	}

	buffer := bytes.NewBuffer([]byte{})
	err = tmpl.Execute(buffer, nil)

	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}
