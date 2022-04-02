# ディレクトリ構成

```bash
.
├── Makefile                    <-- Make to automate build
├── README.md                   <-- This instructions file
├── hello-world                 <-- DynamoDBとの連携
├── login                       <-- Cognitoとの連携
├── signup                      <-- Cognitoとの連携
├── user                        <-- Cognitoとの連携
├── dynamoDB.json               <-- DynamoDBテーブル定義
└── template.yaml
```

# DynamoDBとの連携
## DynamoDB ローカルで立ち上げ

1. npm install dynamodb-admin -g
2. export DYNAMO_ENDPOINT=http://localhost:8000
3. dynamodb-admin
- [参考](https://qiita.com/gzock/items/e0225fd71917c234acce)

## DynamoDB ローカルで作成

1. docker-compose.ymlを作成
2. aws configure set aws_access_key_id dummy --profile local
3. aws configure set aws_secret_access_key dummy --profile local
4. aws configure set region ap-northeast-1 --profile local
5. docker-compose up
6. dynamoDB.json作成
7. aws dynamodb --profile local --endpoint-url http://localhost:8000 create-table --cli-input-json file://dynamoDB.json

## DynamoDB 操作方法(CUI)

1. テーブル確認
- aws dynamodb list-tables --profile local --endpoint-url http://localhost:8000
2. テーブル定義確認
- aws dynamodb describe-table --table-name local_company_table --profile local --endpoint-url http://localhost:8000
3. 値をいれる
- aws dynamodb put-item --table-name local_company_table --profile local --endpoint-url http://localhost:8000 --item '{ "company": {"S":"1"}, "year":{"S":"2"}}'
4. 中身見る
- aws dynamodb scan --table-name --profile local local_company_table --endpoint-url http://localhost:8000

## 立ち上げ
1. docker-compose up -d
2. sam local start-api --docker-network go_api_lamda_dynamodb-local-network
3. http://localhost:3000/hello
- 変更反映はmake build


# 参考

## SAM-CLIインストール（APIGateWayをローカルで使う）

0. 事前にdockerとhomebrewのインストールしておく
1. brew tap aws/tap
2. brew install aws-sam-cli
3. sam --version（確認）
- [参考](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install-mac.html)

## APIGateway　ローカル環境導入

1. sam init --runtime go1.x --name go_lamda
- go mod, go.sumをディレクトリに移動させる
2. make build
3. sam local start-api
