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

type request struct {
	UserName        string `json:"user_name"`
	UserPass        string `json:"user_pass"`
	CognitoClientID string `json:"client_id"`
}

type LoginUserDetail struct {
	Name        string `json:"name"`
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
}

func convertRequestJSON(inputs string) (*request, error) {
	var req request
	err := json.Unmarshal([]byte(inputs), &req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func login(req *request) (*LoginUserDetail, error) {
	svc := cognitoidentityprovider.New(session.New(), &aws.Config{
		Region: aws.String("ap-northeast-1"),
	})

	params := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(req.UserName),
			"PASSWORD": aws.String(req.UserPass),
		},
		ClientId: aws.String(req.CognitoClientID),
	}

	res, err := svc.InitiateAuth(params)
	if err != nil {
		return nil, err
	}

	loginUser := &LoginUserDetail{
		Name:        req.UserName,
		AccessToken: aws.StringValue(res.AuthenticationResult.AccessToken),
		IDToken:     aws.StringValue(res.AuthenticationResult.IdToken),
	}

	return loginUser, nil
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	req, err := convertRequestJSON(request.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	res, err := login(req)
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
