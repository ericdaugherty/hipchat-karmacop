package main

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ericdaugherty/hipchat-go/hipchat"
)

// AWS Provides an interface to AWS
type AWS struct {
	AWSSession *session.Session
	S3         *s3.S3
	DynamoDB   *dynamodb.DynamoDB
}

func newAWS() *AWS {
	a := &AWS{nil, nil, nil}
	a.AWSSession = session.Must(session.NewSession())
	a.S3 = s3.New(a.AWSSession)
	a.DynamoDB = dynamodb.New(a.AWSSession)

	return a
}

func (a *AWS) writeTokenToDB(roomName string, tok *hipchat.OAuthAccessToken) error {

	tokBytes, err := json.Marshal(tok)
	if err != nil {
		log.Println("Error Marshling HipChat Token into String.", err.Error())
		return err
	}

	params := &dynamodb.PutItemInput{
		TableName: aws.String("KarmaCopAuthentication"),
		Item: map[string]*dynamodb.AttributeValue{
			"RoomID": &dynamodb.AttributeValue{
				S: aws.String(roomName),
			},
			"Data": &dynamodb.AttributeValue{
				S: aws.String(string(tokBytes)),
			},
		},
	}
	_, err = a.DynamoDB.PutItem(params)
	if err != nil {
		log.Printf("Error writing Token to DynamoDB for Room %s. %s\n", roomName, err.Error())
		return err
	}

	return nil
}

func (a *AWS) getTokenFromDB(roomName string) (*hipchat.OAuthAccessToken, error) {
	log.Println("Getting Room Auth Token from DynamoDB for Room", roomName)

	params := &dynamodb.GetItemInput{
		TableName: aws.String("KarmaCopAuthentication"),
		Key: map[string]*dynamodb.AttributeValue{
			"RoomID": {
				S: aws.String(roomName),
			},
		},
	}

	resp, err := a.DynamoDB.GetItem(params)
	if err != nil {
		log.Println("Error querying DynamoDB for room", roomName)
		return nil, err
	}

	data := resp.Item["Data"].S
	dataString := *data
	var tok hipchat.OAuthAccessToken
	err = json.Unmarshal([]byte(dataString), &tok)
	if err != nil {
		log.Println("Error Unmarshling", err.Error())
		return nil, err
	}

	return &tok, nil
}

func (a *AWS) writeTokenToS3(tokenName string, tok *hipchat.OAuthAccessToken) error {

	tokBytes, err := json.Marshal(tok)
	if err != nil {
		log.Println("Error Marshling HipChat Token into String.", err.Error())
		return err
	}

	putInput := &s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(tokenName),
		Body:   bytes.NewReader(tokBytes),
	}
	resp, err := a.S3.PutObject(putInput)
	if err != nil {
		log.Println("Error putting object on S3.", err.Error())
		return err
	}
	log.Println(resp)

	return nil
}

func (a *AWS) getTokenFromS3(tokenName string) (*hipchat.OAuthAccessToken, error) {

	getInput := &s3.GetObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(tokenName),
	}

	resp, err := a.S3.GetObject(getInput)
	if err != nil {
		log.Printf("Reqeusted token '%s' not found on S3. %s\n", tokenName, err.Error())
		return nil, err
	}

	var tok hipchat.OAuthAccessToken
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&tok)
	if err != nil {
		log.Println("Error decoding token loaded from s3.", err.Error())
	}

	return &tok, err
}
