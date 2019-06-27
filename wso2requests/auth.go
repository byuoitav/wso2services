package wso2requests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
)

//we take a client key and secret (as defined in ENV Variables, and handle the token request and refresh of that token.
var (
	curAccessKey = ""
	curTTL       = time.Time{}

	clientKey    = os.Getenv("CLIENT_KEY")
	clientSecret = os.Getenv("CLIENT_SECRET")

	tokenRefreshURL = os.Getenv("TOKEN_REFRESH_URL")

	secretMu = &sync.Mutex{}
)

type accessKeyResp struct {
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
	Expires     int    `json:"expires_in"`
	Refresh     string `json:"refresh_token"`
	AccessToken string `json:"access_token"`
}

//GetAccessKey .
func GetAccessKey() (string, *nerr.E) {

	secretMu.Lock()
	defer secretMu.Unlock()

	//check to see if we have a key and if the curTTL is after now
	if curTTL.After(time.Now()) {
		return curAccessKey, nil
	}

	return getAccessKey()
}

func getAccessKey() (string, *nerr.E) {
	//grant_type = client_credentials
	req, err := http.NewRequest("POST", tokenRefreshURL, nil)
	if err != nil {
		log.L.Errorf("Couldn't build request")
		return "", nerr.Translate(err).Addf("Couldn't build request for getting access key")
	}

	req.SetBasicAuth(clientKey, clientSecret)

	q := req.URL.Query()
	q.Add("grant_type", "client_credentials")
	req.URL.RawQuery = q.Encode()

	c := http.Client{
		Timeout: 5 * time.Second,
	}

	log.L.Debugf("Getting access key token from %s", req.URL.String())

	resp, err := c.Do(req)
	if err != nil {
		return "", nerr.Translate(err).Addf("error executing access key request")
	}

	//read the body
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nerr.Translate(err).Addf("couldn't read the access key body")
	}

	if resp.StatusCode/100 != 2 {
		return "", nerr.Create(fmt.Sprintf("Couldn't get the access key. Non 200: body [%s]", b), "request-error")
	}

	var acr accessKeyResp
	err = json.Unmarshal(b, &acr)
	if err != nil {
		return "", nerr.Translate(err).Addf("couldn't get the access key. Couldn't unmarshal response %s", "unmarshal error")
	}

	log.L.Debugf("access key response: %v", acr)

	//we set the ttl
	curTTL = time.Now().Add(time.Duration(acr.Expires)*time.Second - 20)
	curAccessKey = acr.AccessToken

	return curAccessKey, nil
}
