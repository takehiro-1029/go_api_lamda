package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

type signUpUserDetail struct {
	Name     string                                `json:"name"`
	Email    string                                `json:"email"`
	Responce *cognitoidentityprovider.SignUpOutput `json:"responce"`
}

func signUp(name string, pass string, mail string, clientID string) (signUpUserDetail, error) {
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

	res, err := svc.SignUp(params)
	if err != nil {
		return signUpUserDetail{}, err
	}

	signUpUser := signUpUserDetail{
		Name:     name,
		Email:    mail,
		Responce: res,
	}

	return signUpUser, nil
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	res, err := signUp(userName, userPass, userEmail, cognitoClientID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       string(err.Error()),
			StatusCode: 500,
		}, nil
	}

	jsonBytes, _ := json.Marshal(res)

	return events.APIGatewayProxyResponse{
		Body:       string(jsonBytes),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
