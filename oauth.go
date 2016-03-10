// Authored by Allonsy.
// A library that acquires an oauth access for the user.
package fboauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

//crendential type
type credentials struct {
	Clientid string
}

type response struct {
	clientid        string
	Code            string
	UserCode        string `json:"user_code"`
	VerificationUri string `json:"verification_uri"`
	Expiration      int    `json:"expires_in"`
	Interval        int
}

type accessResponse struct {
	AccessCode string `json:"access_token"`
}

func readClientID(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		panic("No file for credentials specified")
	}

	//creds := make(map[string]string)
	var creds credentials
	err = json.Unmarshal(data, &creds)
	if err != nil {
		fmt.Println(err)
		panic("Invalid json in credential file!")
	}
	return creds.Clientid
}

// Grabs the authorization code and returns the verificationURI and accessCode.
// path: the path to the json encoded credential file containing the clientid.
// the json should have the key: "clientid".
// returns the verificationURI and accessCode and error if any occurs.
func GetAuthUrl(path string, permissions []string) (string, string, error) {
	r, e := getAuthCode(path, permissions)
	return r.VerificationUri, r.UserCode, e
}

func getAuthCode(path string, permissions []string) (response, error) {
	clientid := readClientID(path)

	params := url.Values{}
	params.Add("type", "device_code")
	params.Add("client_id", clientid)
	params["scope"] = permissions
	fmt.Println(params.Encode())

	resp, err := http.PostForm("https://graph.facebook.com/oauth/device", params)
	if err != nil {
		fmt.Println(err)
		return response{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println(resp.Status)
		return response{}, errors.New("Bad Request")
	}

	var r response
	r.clientid = clientid
	json.NewDecoder(resp.Body).Decode(&r)

	return r, nil
}

//GetAccessCode: Performs an entire OAuth handshake.
//path: path to credentials file with clientid (see GetAuthCode).
//permissions: a slice of string scopes.
//Returns the access token for queries (valid for ~60 days) or an error.
func GetAccessCode(path string, permissions []string) (string, error) {
	r, e := getAuthCode("credentials.json", []string{"public_profile"})
	if e != nil {
		return "", e
	}

	fmt.Println("Please visit:", r.VerificationUri, "and enter the code:", r.UserCode)

	totalWait := 0
	params := url.Values{}
	params.Add("type", "device_token")
	params.Add("client_id", r.clientid)
	params.Add("code", r.Code)

	interval := time.Tick(time.Duration(int64(time.Second) * int64(r.Interval)))
	for _ = range interval {
		if totalWait >= r.Expiration {
			return "", errors.New("Request Timeout")
		}
		resp, err := http.PostForm("https://graph.facebook.com/oauth/device", params)
		if err != nil {
			return "", err
		}
		if resp.StatusCode == 200 {
			var responseMap accessResponse
			err := json.NewDecoder(resp.Body).Decode(&responseMap)
			if err != nil {
				fmt.Println(err)
				return "", err
			}
			return responseMap.AccessCode, nil
		}
		totalWait += r.Interval
	}
	return "", errors.New("Request Timeout")
}
