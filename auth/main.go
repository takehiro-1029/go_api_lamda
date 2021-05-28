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

type authUserDetail struct {
	Name     string                                       `json:"name"`
	Responce *cognitoidentityprovider.ConfirmSignUpOutput `json:"responce"`
}

func auth(req *common.Request) (*authUserDetail, error) {
	svc := cognitoidentityprovider.New(session.New(), &aws.Config{
		Region: aws.String("ap-northeast-1"),
	})

	params := &cognitoidentityprovider.ConfirmSignUpInput{
		Username:         aws.String(req.UserName),
		ConfirmationCode: aws.String(req.UserAuthCode),
		ClientId:         aws.String(req.CognitoClientID),
	}

	res, err := svc.ConfirmSignUp(params)
	if err != nil {
		return nil, err
	}

	authUser := &authUserDetail{
		Name:     req.UserName,
		Responce: res,
	}

	return authUser, nil

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

	res, err := auth(&req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
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
