package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gopayments/models"
	"gopayments/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"github.com/fatih/structs"
	"github.com/gorilla/mux"
)

func init() {
	fmt.Println("This is the beggining of something great")
}

func getaccesstoken() string {
	auth_params := utils.Get_parameters{ConsumerKey: os.Getenv("CONSUMER_KEY"), ConsumerSecret: os.Getenv("CONSUMER_SECRET"), ApiUrl: os.Getenv("MPESA_URL") + os.Getenv("MPESA_AUTH_URL")}
	resp, err := utils.Get_request(auth_params)
	if err != nil {
		log.Fatalln(" Error found sending auth request" + err.Error())
	}
	access_token := resp["access_token"].(string)
	return access_token
}
func stkpush(w http.ResponseWriter, r *http.Request) {
	fmt.Println("This is the MPESA URL " + os.Getenv("MPESA_URL"))
	access_token := getaccesstoken()
	reqBody, _ := ioutil.ReadAll(r.Body)
	var postBody map[string]interface{}
	json.Unmarshal(reqBody, &postBody)

	ispaybill := true
	fmt.Println(postBody)
	phoneNumber := postBody["phoneNumber"].(string)
	amount := postBody["amount"].(string)
	account := postBody["account"].(string)
	fmt.Println(os.Getenv("MPESA_URL"))
	apiurl := os.Getenv("MPESA_URL") + "/mpesa/stkpush/v1/processrequest"
	mpesa_shortcode := os.Getenv("MPESA_SHORTCODE")
	lnm_pass_key := os.Getenv("LNM_PASSKEY")
	transactiondesc := "just an mpesa transaction"

	time := time.Now().Add(time.Hour * 3)
	timestamp := time.Format("20060102150405")

	encode_data := mpesa_shortcode + lnm_pass_key + timestamp
	password := base64.StdEncoding.EncodeToString([]byte(encode_data))
	phoneint, err := strconv.Atoi(phoneNumber)

	if err != nil {
		log.Fatalln("Error converting phone number to integer")
	}
	wallet := models.Wallet{ID: phoneint, Balance: 0}
	var transactiontype string
	if ispaybill {
		transactiontype = "CustomerPayBillOnline"
	}

	// create wallet if does not exits
	db := utils.ConnectDB()
	db.Create(&wallet)
	init_reqBody := structs.Map(&models.StkPush{
		BusinessShortCode: mpesa_shortcode,
		Password:          password,
		Timestamp:         timestamp,
		TransactionType:   transactiontype,
		Amount:            amount,
		PartyA:            phoneNumber,
		PartyB:            os.Getenv("C2B_ONLINE_PARTY_B"),
		PhoneNumber:       phoneNumber,
		CallBackURL:       os.Getenv("C2B_ONLINE_CHECKOUT_CALLBACK_URL"),
		AccountReference:  account,
		TransactionDesc:   transactiondesc,
	})

	fmt.Println(init_reqBody)
	resp, err := utils.Post_request(apiurl, init_reqBody, access_token)
	if err != nil {
		log.Fatalln("There exists an error" + err.Error())
	}
	json.Marshal(&resp)
	json.NewEncoder(w).Encode(resp)

}

func stkhandlecb(w http.ResponseWriter, r *http.Request) {
	fmt.Println("This is the skpush handlecallback endpoint")
}
func c2b(w http.ResponseWriter, r *http.Request) {
	fmt.Println("This is the c2b endpoint")

}
func c2bhandlecb(w http.ResponseWriter, r *http.Request) {
	fmt.Println("This is the c2b handlecallback endpoint ")

}

func b2c(w http.ResponseWriter, r *http.Request) {
	fmt.Println("This is the c2b endpoint")

}

func send_stkpush_query(checkout_id string) (map[string]interface{}, error) {
	access_token := getaccesstoken()
	mpesa_shortcode := os.Getenv("MPESA_SHORTCODE")
	// fmt.Println(postBody)
	// checkout_id := postBody["checkout_id"].(string)
	apiurl := os.Getenv("MPESA_URL") + "/mpesa/stkpushquery/v1/query"

	time := time.Now().Add(time.Hour * 3)
	timestamp := time.Format("20060102150405")
	lnm_pass_key := os.Getenv("LNM_PASSKEY")
	encode_data := mpesa_shortcode + lnm_pass_key + timestamp
	password := base64.StdEncoding.EncodeToString([]byte(encode_data))
	init_reqBody := structs.Map(&models.StkPushQuery{
		BusinessShortCode: mpesa_shortcode,
		Password:          password,
		Timestamp:         timestamp,
		CheckoutRequestID: checkout_id,
	})

	resp, err := utils.Post_request(apiurl, init_reqBody, access_token)

	if err != nil {

		log.Fatalln("There exists an error" + err.Error())

	}
	json.Marshal(&resp)
	return resp, err
}
func stkpushquery(w http.ResponseWriter, r *http.Request) {
	fmt.Println("This is the c2b endpoint")
	reqBody, _ := ioutil.ReadAll(r.Body)
	var postBody map[string]interface{}
	json.Unmarshal(reqBody, &postBody)
	access_token := getaccesstoken()
	mpesa_shortcode := os.Getenv("MPESA_SHORTCODE")

	checkout_id := postBody["checkout_id"].(string)
	apiurl := os.Getenv("MPESA_URL") + "/mpesa/stkpushquery/v1/query"

	time := time.Now().Add(time.Hour * 3)
	timestamp := time.Format("20060102150405")
	lnm_pass_key := os.Getenv("LNM_PASSKEY")

	encode_data := mpesa_shortcode + lnm_pass_key + timestamp
	password := base64.StdEncoding.EncodeToString([]byte(encode_data))

	init_reqBody := structs.Map(&models.StkPushQuery{
		BusinessShortCode: mpesa_shortcode,
		Password:          password,
		Timestamp:         timestamp,
		CheckoutRequestID: checkout_id,
	})

	resp, err := utils.Post_request(apiurl, init_reqBody, access_token)

	if err != nil {
		log.Fatalln("There exists an error" + err.Error())
	}
	json.Marshal(&resp)
	json.NewEncoder(w).Encode(resp)

}

func stkpushquery_conv(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var postBody map[string]interface{}
	json.Unmarshal(reqBody, &postBody)
	fmt.Println("====This is the post body=====")
	fmt.Println(postBody)
	init_resp := make(map[string]interface{})
	init_resp["Status"] = 200
	init_resp["Response"] = "Request received processing ongoing"
	json.Marshal(&init_resp)
	json.NewEncoder(w).Encode(init_resp)
	go loop_check_request(postBody)

}
func send_conversation_update(status string, conversational_id string, access_token string) {
	api_url := "http://localhost:5005/conversations/" + conversational_id + "/trigger_intent?output_channel=latest&token=thisismysecret"
	conv_response := make(map[string]interface{})
	conv_response["name"] = "EXTERNAL_payment_notification"
	conv_response["entities"] = map[string]string{"payment_status": status}
	_, err := utils.Post_request(api_url, conv_response, access_token)
	if err != nil {
		log.Fatalln("Something went wrong" + err.Error())
	}

}
func loop_check_request(postBody map[string]interface{}) {
	init_time := 0
	fmt.Println(postBody)
	conversational_id := postBody["conversation_id"].(string)
	access_token := ""
	checkout_id := postBody["checkout_id"].(string)
loop:
	for {
		resp, err := send_stkpush_query(checkout_id)
		fmt.Println(resp)
		switch resp["ResultCode"] {
		case "0":
			send_conversation_update("successful", conversational_id, access_token)
			break loop
		default:
			fmt.Println("This is the default log")
		}
		fmt.Println(err)
		time.Sleep(10 * time.Second)
		init_time = init_time + 1
		if init_time > 10 {
			send_conversation_update("not successful", conversational_id, access_token)
			break
		}

	}
}
func b2b(w http.ResponseWriter, r *http.Request) {
	fmt.Println("This is the b2b endpoint")

}
func main() {
	Router := mux.NewRouter().StrictSlash(true)
	Router.HandleFunc("/stkpush", stkpush).Methods("POST")
	Router.HandleFunc("/stkpushquery", stkpushquery).Methods("Post")
	Router.HandleFunc("/stkhandlecb", stkhandlecb).Methods("POST")
	Router.HandleFunc("/c2b", c2b).Methods("POST")
	Router.HandleFunc("/stkpushqueryconv", stkpushquery_conv).Methods("POST")
	Router.HandleFunc("/c2bhandlecb", c2bhandlecb).Methods("POST")
	Router.HandleFunc("/b2c", b2c).Methods("POST")
	Router.HandleFunc("/b2b", b2b).Methods("POST")
	log.Fatal(http.ListenAndServe(":10000", Router))

}
