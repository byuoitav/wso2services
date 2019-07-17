package wso2requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
)

//MakeWSO2Request makes a generic WSO2 request
//toReturn should be a pointer
func MakeWSO2Request(method, url string, body interface{}, toReturn interface{}) *nerr.E {
	return MakeWSO2RequestWithHeaders(method, url, body, toReturn, nil)
}

//MakeWSO2RequestWithHeaders makes a generic WSO2 request with headers
//toReturn should be a pointer
func MakeWSO2RequestWithHeaders(method, url string, body interface{}, toReturn interface{}, headers map[string]string) *nerr.E {
	//	log.L.Debugf("Making %v request against %v at %v", method, url, time.Now())

	key, er := GetAccessKey()
	if er != nil {
		return er.Addf("Couldn't make WSO2 request")
	}

	//attach key
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

	for {
		hasRetried := false

		req, err := http.NewRequest(method, url, bytes.NewBuffer(b))
		if err != nil {
			return nerr.Translate(err).Addf("Couldn't build WSO2 request")
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", key))

		for k, v := range headers {
			log.L.Debugf("Setting header %v to %v", k, v)
			req.Header.Set(k, v)
		}

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
			if resp.StatusCode == 400 && len(rb) == 0 && !hasRetried {
				//if we get a 400 and a blank body and we haven't retried, then just try again
				log.L.Debugf("Retrying WSO2 request")
				hasRetried = true
				continue
			}

			return nerr.Create(fmt.Sprintf("response code %v: %s", resp.StatusCode, rb), "request-error")
		}

		err = json.Unmarshal(rb, toReturn)
		if err != nil {
			return nerr.Translate(err).Addf("Couldn't unmarshal response %s", "unmarshal error")
		}

		return nil
	}
}
