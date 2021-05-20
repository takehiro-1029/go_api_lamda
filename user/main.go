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
	userToken = "registedUserName"
)

type userDetail struct {
	Name     string                                 `json:"name"`
	Responce *cognitoidentityprovider.GetUserOutput `json:"responce"`
}

func getUser(token string) (userDetail, error) {
	svc := cognitoidentityprovider.New(session.New(), &aws.Config{
		Region: aws.String("ap-northeast-1"),
	})

	params := &cognitoidentityprovider.GetUserInput{
		AccessToken: aws.String(token),
	}

	res, err := svc.GetUser(params)
	if err != nil {
		return userDetail{}, err
	}

	user := userDetail{
		Name:     aws.StringValue(res.Username),
		Responce: res,
	}

	return user, nil
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	res, err := getUser(userToken)
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
