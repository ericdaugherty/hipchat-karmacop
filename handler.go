package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var k karmaCop

// Handler processes requests from AWS Lambda
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	switch request.Path {
	case "/healthcheck":
		return k.healthcheck(request)
	case "/atlassian-connect.json":
		return k.atlassianConnect(request)
	case "/":
		return k.atlassianConnect(request)
	case "/installable":
		return k.installable(request)
	case "/test":
		return k.test(request)
	case "/ninja":
		return k.ninja(request)
	default:
		return events.APIGatewayProxyResponse{
			Body:       "",
			StatusCode: 400,
		}, nil
	}
}

func main() {
	k = karmaCop{"https://bz96q93bj3.execute-api.us-east-1.amazonaws.com/Prod"}

	lambda.Start(Handler)
}
