package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

// APIResponse struct with numverify Json structure
type APIResponse struct {
	Valid    bool
	Number   string
	LineType string `json:"line_type"`
}

// ValidateMobilePhone function checks if passed phone number is valid
func ValidateMobilePhone(phoneNumber string) (*APIResponse, error) {
	response, err := getPhoneValidation(phoneNumber)
	if err != nil {
		return nil, err
	}

	if response.LineType != "mobile" {
		response.Valid = false
	}

	return response, nil
}

func createRequest(phoneNumber string) *http.Request {
	const url = "http://apilayer.net/api/validate"
	req, _ := http.NewRequest("GET", url, nil)

	q := req.URL.Query()
	q.Add("access_key", os.Getenv("NUMVERIFY_API_KEY"))
	q.Add("country_code", "")
	q.Add("format", "1")
	q.Add("number", phoneNumber)

	req.URL.RawQuery = q.Encode()

	return req
}

var getPhoneValidation = func(phoneNumber string) (response *APIResponse, err error) {
	const t = time.Duration(3 * time.Second)
	client := http.Client{
		Timeout: t,
	}

	req := createRequest(phoneNumber)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}
