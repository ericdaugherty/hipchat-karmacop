package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ericdaugherty/hipchat-go/hipchat"
)

const s3Bucket = "karmacoprooms"

type karmaCop struct {
	BaseURL string
}

func (k *karmaCop) healthcheck(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Processing Request to healthcheck")
	b, err := json.Marshal([]string{"OK"})

	if err != nil {
		return errorResponse(err)
	}
	return okResponse(string(b))
}

func (k *karmaCop) atlassianConnect(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Processing Request to atlassianConect")
	tmpl, err := template.New("config").Parse(descriptorTemplate)
	if err != nil {
		return errorResponse(err)
	}

	vals := map[string]string{
		"LocalBaseUrl": k.BaseURL,
	}

	var w bytes.Buffer
	err = tmpl.Execute(&w, vals)
	if err != nil {
		return errorResponse(err)
	}
	return okResponse(w.String())
}

func (k *karmaCop) installable(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Processing request to installable")

	var payLoad map[string]interface{}
	decoder := json.NewDecoder(strings.NewReader(request.Body))
	err := decoder.Decode(&payLoad)
	if err != nil {
		log.Println("Error decoding install request body.", err.Error())
		return errorResponse(err)
	}

	log.Println("Installable Called, Payload:")
	log.Println(payLoad)

	oAuthID := payLoad["oauthId"].(string)
	oAuthSecret := payLoad["oauthSecret"].(string)

	credentials := hipchat.ClientCredentials{
		ClientID:     oAuthID,
		ClientSecret: oAuthSecret,
	}

	hc := hipchat.NewClient("")
	tok, _, err := hc.GenerateToken(credentials, []string{hipchat.ScopeSendNotification})
	if err != nil {
		log.Fatalln("Error generating HipChat Client Token.", err.Error())
		return errorResponse(err)
	}

	a := newAWS()
	//a.writeRoomAuthTokenToDB(roomName, tok)
	a.writeTokenToS3(oAuthID, tok)
	if err != nil {
		return errorResponse(err)
	}

	return okResponse("")
}

func (k *karmaCop) test(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Processing Request to test.")

	var err error
	var payLoad map[string]interface{}
	decoder := json.NewDecoder(strings.NewReader(request.Body))
	err = decoder.Decode(&payLoad)
	if err != nil {
		log.Println("Error decoding hook request body.", err.Error())
		return errorResponse(err)
	}

	log.Println("Test Called, Payload:")
	log.Println(payLoad)

	oAuthID := payLoad["oauth_client_id"].(string)

	a := newAWS()
	//tok, err := a.getRoomAuthTokenFromDB(roomName)
	tok, err := a.getTokenFromS3(oAuthID)
	if err != nil {
		return errorResponse(err)
	}

	notifRq := &hipchat.NotificationRequest{
		Message: "KarmaCop is Operational. Now with Native Go!",
	}

	hipChat := hipchat.NewClient(tok.AccessToken)

	roomName := strconv.Itoa(int((payLoad["item"].(map[string]interface{}))["room"].(map[string]interface{})["id"].(float64)))
	hipChat.Room.Notification(roomName, notifRq)

	return okResponse("")
}

func (k *karmaCop) ninja(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Processing Request to ninja.")

	var err error
	var payLoad map[string]interface{}
	decoder := json.NewDecoder(strings.NewReader(request.Body))
	err = decoder.Decode(&payLoad)
	if err != nil {
		log.Println("Error decoding hook request body.", err.Error())
		return errorResponse(err)
	}

	log.Println("Ninja Called, Payload:")
	log.Println(payLoad)

	if !checkMessage(payLoad["item"].(map[string]interface{})["message"].(map[string]interface{})["message"].(string)) {
		return okResponse("")
	}

	oAuthID := payLoad["oauth_client_id"].(string)

	a := newAWS()
	//tok, err := a.getRoomAuthTokenFromDB(roomName)
	tok, err := a.getTokenFromS3(oAuthID)
	if err != nil {
		return errorResponse(err)
	}

	fromItem := payLoad["item"].(map[string]interface{})["message"].(map[string]interface{})["from"].(map[string]interface{})
	userName := fromItem["name"].(string)
	// userMention := fromItem["mention_name"].(string)

	message := fmt.Sprintf("Hey %s, Ninjas may be cool, but you are not. Don't be a Karma Ninja", userName)

	notifRq := &hipchat.NotificationRequest{
		Message:       message,
		MessageFormat: "html",
		Color:         "red",
	}

	hipChat := hipchat.NewClient(tok.AccessToken)

	roomName := strconv.Itoa(int((payLoad["item"].(map[string]interface{}))["room"].(map[string]interface{})["id"].(float64)))
	hipChat.Room.Notification(roomName, notifRq)

	return okResponse("")
}

func checkMessage(message string) bool {
	return strings.Contains(message, "@") && (strings.Contains(message, "--") || strings.Contains(message, "++"))
}

func okResponse(body string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       body,
		StatusCode: 200,
	}, nil
}

func errorResponse(err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       err.Error(),
		StatusCode: 500,
	}, nil
}
