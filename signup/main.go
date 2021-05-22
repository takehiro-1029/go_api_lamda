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
	UserEmail       string `json:"user_email"`
	UserPass        string `json:"user_pass"`
	CognitoClientID string `json:"client_id"`
}

type signUpUserDetail struct {
	Name     string                                `json:"name"`
	Email    string                                `json:"email"`
	Responce *cognitoidentityprovider.SignUpOutput `json:"responce"`
}

func convertRequestJSON(inputs string) (*request, error) {
	var req request
	err := json.Unmarshal([]byte(inputs), &req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func signUp(req *request) (*signUpUserDetail, error) {
	svc := cognitoidentityprovider.New(session.New(), &aws.Config{
		Region: aws.String("ap-northeast-1"),
	})
	ua := &cognitoidentityprovider.AttributeType{
		Name:  aws.String("email"),
		Value: aws.String(req.UserEmail),
	}
	params := &cognitoidentityprovider.SignUpInput{
		Username: aws.String(req.UserName),
		Password: aws.String(req.UserPass),
		ClientId: aws.String(req.CognitoClientID),
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			ua,
		},
	}

	res, err := svc.SignUp(params)
	if err != nil {
		return nil, err
	}

	signUpUser := &signUpUserDetail{
		Name:     req.UserName,
		Email:    req.UserEmail,
		Responce: res,
	}

	return signUpUser, nil
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	req, err := convertRequestJSON(request.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	res, err := signUp(req)
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
