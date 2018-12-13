package wso2requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
)

//MakeWSO2Request makes a generic WSO2 request
//toReturn should be a pointer
func MakeWSO2Request(method, url string, body interface{}, toReturn interface{}) *nerr.E {

	log.L.Debugf("Making %v request against %v", method, url)
	//attach key
	key := os.Getenv("WSO2_TOKEN")
	if len(key) == 0 {
		return nerr.Create("No WSO2 token defined", "no-creds")
	}
	var b []byte
	var ok bool
	var err error

	if body != nil {
		if b, ok = body.([]byte); !ok {
			b, err = json.Marshal(body)
			if err != nil {
				return nerr.Translate(err).Addf("Couldn't marhsal request")
			}
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(b))
	if err != nil {
		return nerr.Translate(err).Addf("Couldn't build WSO2 request")
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", key))

	c := http.Client{
		Timeout: 20 * time.Second, //I wish we could make this shorter... but alas.
	}

	resp, err := c.Do(req)
	if err != nil {
		return nerr.Translate(err).Addf("Couldn't make WSO2 request")
	}
	defer resp.Body.Close()

	rb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nerr.Translate(err).Addf("Couldn't read response body")
	}

	if resp.StatusCode/100 != 2 {
		return nerr.Create(fmt.Sprintf("Non 200: body %s", rb), "request-error")
	}

	err = json.Unmarshal(rb, toReturn)
	if err != nil {
		return nerr.Translate(err).Addf("Couldn't unmarshal response %s", "unmarshal error")
	}

	return nil
}
