package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ericdaugherty/hipchat-go/hipchat"
)

const s3Bucket = "karmacoprooms"

type karmaCop struct {
	BaseURL string
}

func (k *karmaCop) healthcheck(w http.ResponseWriter, r *http.Request) {
	log.Println("Processing Request to healthcheck")
	json.NewEncoder(w).Encode([]string{"OK"})
}

func (k *karmaCop) atlassianConnect(w http.ResponseWriter, r *http.Request) {
	log.Println("Processing Request to atlassianConect")
	tmpl, err := template.New("config").Parse(descriptorTemplate)
	if err != nil {
		log.Fatalln("Error Parsing Atlassian Connect Template.", err.Error())
	}

	vals := map[string]string{
		"LocalBaseUrl": k.BaseURL,
	}
	err = tmpl.Execute(w, vals)
	if err != nil {
		log.Fatalln("Error Executing Atlassian Connect Template.", err.Error())
	}
}

func (k *karmaCop) installable(w http.ResponseWriter, r *http.Request) {
	log.Println("Processing request to installable")

	var payLoad map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payLoad)
	if err != nil {
		log.Println("Error decoding install request body.", err.Error())
		returnError(w, r)
		return
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
		returnError(w, r)
		return
	}

	a := newAWS()
	//a.writeRoomAuthTokenToDB(roomName, tok)
	a.writeTokenToS3(oAuthID, tok)
	if err != nil {
		returnError(w, r)
		return
	}

	returnOK(w, r)
}

func (k *karmaCop) test(w http.ResponseWriter, r *http.Request) {
	log.Println("Processing Request to test.")

	var err error
	var payLoad map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&payLoad)
	if err != nil {
		log.Println("Error decoding hook request body.", err.Error())
		returnError(w, r)
		return
	}

	log.Println("Test Called, Payload:")
	log.Println(payLoad)

	oAuthID := payLoad["oauth_client_id"].(string)

	a := newAWS()
	//tok, err := a.getRoomAuthTokenFromDB(roomName)
	tok, err := a.getTokenFromS3(oAuthID)
	if err != nil {
		returnError(w, r)
		return
	}

	notifRq := &hipchat.NotificationRequest{
		Message: "KarmaCop is Operational.",
	}

	hipChat := hipchat.NewClient(tok.AccessToken)

	roomName := strconv.Itoa(int((payLoad["item"].(map[string]interface{}))["room"].(map[string]interface{})["id"].(float64)))
	hipChat.Room.Notification(roomName, notifRq)

	returnOK(w, r)
}

func (k *karmaCop) ninja(w http.ResponseWriter, r *http.Request) {
	log.Println("Processing Request to ninja.")

	var err error
	var payLoad map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&payLoad)
	if err != nil {
		log.Println("Error decoding hook request body.", err.Error())
		returnError(w, r)
		return
	}

	log.Println("Ninja Called, Payload:")
	log.Println(payLoad)

	if !checkMessage(payLoad["item"].(map[string]interface{})["message"].(map[string]interface{})["message"].(string)) {
		returnOK(w, r)
		return
	}

	oAuthID := payLoad["oauth_client_id"].(string)

	a := newAWS()
	//tok, err := a.getRoomAuthTokenFromDB(roomName)
	tok, err := a.getTokenFromS3(oAuthID)
	if err != nil {
		returnError(w, r)
		return
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

	returnOK(w, r)
}

func returnOK(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
}

func returnError(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusInternalServerError)
}

func checkMessage(message string) bool {
	return strings.Contains(message, "@") && (strings.Contains(message, "--") || strings.Contains(message, "++"))
}
