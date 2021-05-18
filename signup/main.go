package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

//TODO:user情報はpostのBodyから受け取る
const (
	userName        = "registedUserName"
	userEmail       = "registedUserEmail"
	userPass        = "registedUserPass"
	cognitoClientID = "cognitoClientID"
)

func signUp(name string, pass string, mail string, clientID string) error {
	svc := cognitoidentityprovider.New(session.New(), &aws.Config{
		Region: aws.String("ap-northeast-1"),
	})
	ua := &cognitoidentityprovider.AttributeType{
		Name:  aws.String("email"),
		Value: aws.String(mail),
	}
	params := &cognitoidentityprovider.SignUpInput{
		Username: aws.String(name),
		Password: aws.String(pass),
		ClientId: aws.String(clientID),
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			ua,
		},
	}

	_, err := svc.SignUp(params)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	if err := signUp(userName, userEmail, userPass, cognitoClientID); err != nil {
		log.Fatal(err)
	}
}
