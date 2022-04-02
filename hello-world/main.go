package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var dbEndpoint = "http://dynamodb-local:8000"
var region = "ap-northeast-1"
var tableName = "local_company_table"

type Company struct {
	Company *string `json:"company"`
	Year    *string `json:"year"`
}

func CreateCompany(company string, year string) Company {
	return Company{
		Company: aws.String(company),
		Year:    aws.String(year),
	}
}

func DB() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint: aws.String(dbEndpoint),
		Region:   aws.String(region),
	}))

	return dynamodb.New(sess)
}

func putItem(ctx context.Context, db *dynamodb.DynamoDB, tableName string, v interface{}) error {
	av, err := dynamodbattribute.MarshalMap(v)
	if err != nil {
		return fmt.Errorf("dynamodb attribute marshalling map: %w", err)
	}

	i := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	if _, err = db.PutItemWithContext(ctx, i); err != nil {
		return fmt.Errorf("dynamodb put item: %w", err)
	}
	return nil
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	item := CreateCompany("company", "year")

	jsonBytes, _ := json.Marshal(item)

	if err := putItem(ctx, DB(), tableName, item); err != nil {
		fmt.Printf("%v", err)
		return events.APIGatewayProxyResponse{
			Body:       string(jsonBytes),
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       string(jsonBytes),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
