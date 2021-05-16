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

	// DBと接続するセッションを作る→DB接続
	sess, err := session.NewSession()
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	db := dynamodb.New(sess)

	// リクエストボディのjsonから、Item構造体(DB用データの構造体)を作成
	reqBody := request.Body
	resBodyJSONBytes := ([]byte)(reqBody)
	item := Item{}
	if err := json.Unmarshal(resBodyJSONBytes, &item); err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	// Item構造体から、inputするデータを用意
	inputAV, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String("user"),
		Item:      inputAV,
	}

	// insert実行
	_, err = db.PutItem(input)
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
