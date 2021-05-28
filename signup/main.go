package main

import (
	"context"
	"encoding/json"
	"hello-world/common"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type signUpUserDetail struct {
	Name     string                                `json:"name"`
	Email    string                                `json:"email"`
	Responce *cognitoidentityprovider.SignUpOutput `json:"responce"`
}

func signUp(req *common.Request) (*signUpUserDetail, error) {
	mySession := session.Must(session.NewSession())
	svc := cognitoidentityprovider.New(mySession, &aws.Config{
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

	var req common.Request
	err := common.ConvertRequestToJSON(&req, request.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	res, err := signUp(&req)
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
