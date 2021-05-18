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
	userPass        = "registedUserPass"
	cognitoClientID = "cognitoClientID"
)

type LoginUserDetail struct {
	Name        string `json:"name"`
	AccessToken string `json:"access_token"`
}

func login(name string, pass string, clientID string) (LoginUserDetail, error) {
	svc := cognitoidentityprovider.New(session.New(), &aws.Config{
		Region: aws.String("ap-northeast-1"),
	})

	params := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(name),
			"PASSWORD": aws.String(pass),
		},
		ClientId: aws.String(clientID),
	}

	res, err := svc.InitiateAuth(params)
	if err != nil {
		return LoginUserDetail{}, err
	}

	loginUser := LoginUserDetail{
		Name:        name,
		AccessToken: aws.StringValue(res.AuthenticationResult.AccessToken),
	}

	return loginUser, nil
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	res, err := login(userName, userPass, cognitoClientID)
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
