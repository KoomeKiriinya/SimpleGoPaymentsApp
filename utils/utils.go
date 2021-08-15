package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gopayments/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//  This variables are used by the redirect policy fun
var consumerKey string
var consumerSecret string
var accessToken string
var getresult map[string]interface{}
var postresult map[string]interface{}

type Get_parameters struct {
	ConsumerKey    string
	ConsumerSecret string
	ApiUrl         string
}

// connect to db function
func ConnectDB() *gorm.DB {

	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dataBaseName := os.Getenv("DB_NAME")
	dataBaseHost := os.Getenv("DB_HOST")

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dataBaseHost, username, dataBaseName, password)

	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})

	if err != nil {
		log.Fatalln("There exist an error" + err.Error())
	}

	db.AutoMigrate(&models.Wallet{}, &models.Transaction{})
	// defer db.Close()

	return db

}

// Go lang drops headers on redirect polict thus best to specify this
func getredirectPolicyFunc(req *http.Request, via []*http.Request) error {
	consumerKey = os.Getenv("CONSUMER_KEY")
	consumerSecret = os.Getenv("CONSUMER_SECRET")
	req.Header.Add("Authorization", "Basic "+basicAuth(consumerKey, consumerSecret))
	return nil
}

func postredirectPolicyFunc(req *http.Request, via []*http.Request) error {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	return nil

}

func basicAuth(consumerKey, consumerSecret string) string {
	auth := consumerKey + ":" + consumerSecret
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// get request function
func Get_request(p Get_parameters) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", p.ApiUrl, nil)
	if err != nil {
		log.Fatal("The exist and issue initializing the GET request" + err.Error())
	}

	fmt.Println(p.ApiUrl)
	req.Header.Add("Authorization", "Basic "+basicAuth(p.ConsumerKey, p.ConsumerSecret))
	// Initialize a http client
	client := &http.Client{
		CheckRedirect: getredirectPolicyFunc,
	}
	resp, err := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("The exist and issue sending the GET request" + err.Error())
	}
	json.Unmarshal(body, &getresult)
	return getresult, err
}

// Post request function

func Post_request(ApiUrl string, Body map[string]interface{}, AccessToken string) (map[string]interface{}, error) {
	json_Body, _ := json.Marshal(Body)
	req_Body := bytes.NewBuffer(json_Body)
	accessToken = AccessToken
	req, err := http.NewRequest("POST", ApiUrl, req_Body)
	if err != nil {
		log.Fatalln("There exists and error creating the request " + err.Error())
	}
	req.Header.Add("Content-Type", "application/json")
	if accessToken != "" {
		req.Header.Add("Authorization", "Bearer "+AccessToken)
	}
	fmt.Println(accessToken)
	client := &http.Client{
		CheckRedirect: postredirectPolicyFunc,
	}
	resp, err := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("The exist and issue sending the POST request" + err.Error())
	}
	json.Unmarshal(body, &postresult)
	return postresult, err

}