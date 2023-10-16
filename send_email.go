package main

import (
	"bytes"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// region = "us-east-1"
func NewSession(region string) (*ses.SES, error) {
	creds := credentials.NewEnvCredentials()

	// Retrieve the credentials value
	_, err := creds.Get()
	if err != nil {
		return nil, err
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "default",
		Config: aws.Config{
			Credentials: creds,
			Region:      aws.String(region),
		},
	})

	if err != nil {
		return nil, err
	}

	svc := ses.New(sess)
	return svc, nil
}

func BuildEmailTemplate(validlink string) (string, error) {
	code := "Hello, this is the registration email for EasyTracker, please <a class=\"ulink\" href=\"{{.}}\" target=\"_blank\">click the link to complete registration</a>."
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

func BuildMessage(to string, from string, registMessage string) *ses.SendEmailInput {
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(to),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(registMessage),
				},
				// Text: &ses.Content{
				// 	Charset: aws.String("UTF-8"),
				// 	Data:    aws.String("Hello, this is the registration email for EasyTracker, please click the link to complete the registration."),
				// },
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String("EasyTracker"),
			},
		},

		Source: aws.String(from),
	}
	return input
}
