package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

const POST = "POST"

var (
	baseURL  = os.Getenv("KAIROS_API_URL")
	clientID = os.Getenv("KAIROS_CLIENT_ID")
	//clientSecret = os.Getenv("KAIROS_CLIENT_SECRET")
)

type LoginResponse struct {
	Challenge string `json:"challenge"`
}

func prettyPrint(data map[string]interface{}) string {
	jsondata, _ := json.MarshalIndent(data, "", "    ")
	return string(jsondata)
}

func login() (challenge string, err error) {
	url, err := url.JoinPath(baseURL, "login")
	if err != nil {
		return "", err
	}
	fmt.Println("Attempt login to:", url)

	data := map[string]interface{}{
		"client_id": clientID,
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	j := prettyPrint(data)
	fmt.Println("With payload:", j)

	client := &http.Client{}
	req, err := http.NewRequest(POST, url, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var response LoginResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	return response.Challenge, nil
}

func main() {
	fmt.Println("Kairos API usage example in Go")
	fmt.Println("------------------------------")

	challenge, err := login()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Challenge:", challenge)
}
