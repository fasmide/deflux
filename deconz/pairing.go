package deconz

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type APIKey string

type pairRequest struct {
	DeviceType string `json:"devicetype"`
}

type pairResponse []pair

type pair struct {
	Success struct {
		Username string
	}
}

type pairFailureResponse []pairFailure

type pairFailure struct {
	Error struct {
		Description string
	}
}

// Pair tries to pair with deconz and returns a pairing with an API key
func Pair(u url.URL) (APIKey, error) {
	// to pair we must send a POST request to "/api" containing a pairRequest
	u.Path = "/api"

	pr := pairRequest{
		DeviceType: "Deflux",
	}

	var buff bytes.Buffer
	enc := json.NewEncoder(&buff)
	err := enc.Encode(pr)
	if err != nil {
		return "", fmt.Errorf("unable to marshal pair request: %s", err)
	}

	// send POST request and read body
	response, err := http.Post(u.String(), "application/json", &buff)
	if err != nil {
		return "", fmt.Errorf("unable to send post request: %s", err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read body: %s", err)
	}

	// this error happens if the gateway is locked
	if response.StatusCode == http.StatusForbidden {
		var failure pairFailureResponse
		err = json.Unmarshal(body, &failure)
		if err != nil {
			return "", fmt.Errorf("unable to parse failure message drom deconz: %s", err)
		}
		return "", fmt.Errorf("unable to pair with deconz: %s", failure[0].Error.Description)
	}

	// we dont know what to do with other then 200 statuses
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected statuscode from deconz: %d\n%s", response.StatusCode, body)
	}

	// finally, if none of the above failed, extract APIKey from response
	var pairResp pairResponse
	err = json.Unmarshal(body, &pairResp)
	if err != nil {
		return "", fmt.Errorf("unable to parse json from pair response: %s", err)
	}

	return APIKey(pairResp[0].Success.Username), nil
}
