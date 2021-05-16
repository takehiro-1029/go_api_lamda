package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Item DBに入れるデータ
type Item struct {
	UserID  int    `dynamodbav:"userid" json:"userid"`
	Address string `dynamodbav:"address" json:"address"`
	Email   string `dynamodbav:"email" json:"email"`
	Gender  string `dynamodbav:"gender" json:"gender"`
	Name    string `dynamodbav:"name" json:"name"`
}

// Response Lambdaが返答するデータ
type Response struct {
	RequestMethod string `json:"RequestMethod"`
	Result        Item   `json:"Result"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	method := request.HTTPMethod
	pathparam := request.PathParameters["userid"]

	// DB接続
	sess, err := session.NewSession()
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	db := dynamodb.New(sess)

	// 検索条件を用意
	getParam := &dynamodb.GetItemInput{
		TableName: aws.String("user"),
		Key: map[string]*dynamodb.AttributeValue{
			"userid": {
				N: aws.String(pathparam),
			},
		},
	}

	// 検索
	result, err := db.GetItem(getParam)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 404,
		}, err
	}

	// 結果を構造体にパース
	item := Item{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	// httpレスポンス作成
	res := Response{
		RequestMethod: method,
		Result:        item,
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
