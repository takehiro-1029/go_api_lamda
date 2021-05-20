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
	userAuthCode    = "registedUserAuthCode"
	cognitoClientID = "cognitoClientID"
)

type authUserDetail struct {
	Name     string                                       `json:"name"`
	Responce *cognitoidentityprovider.ConfirmSignUpOutput `json:"responce"`
}

func auth(name string, authCode string, clientId string) (authUserDetail, error) {
	svc := cognitoidentityprovider.New(session.New(), &aws.Config{
		Region: aws.String("ap-northeast-1"),
	})

	params := &cognitoidentityprovider.ConfirmSignUpInput{
		Username:         aws.String(name),
		ConfirmationCode: aws.String(authCode),
		ClientId:         aws.String(clientId),
	}

	res, err := svc.ConfirmSignUp(params)
	if err != nil {
		return authUserDetail{}, err
	}

	authUser := authUserDetail{
		Name:     name,
		Responce: res,
	}

	return authUser, nil

}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	res, err := auth(userName, userAuthCode, cognitoClientID)
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
