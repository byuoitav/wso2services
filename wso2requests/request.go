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

//MakeWSO2RequestReturnResponse makes a generic WSO2 request
//toReturn should be a pointer
func MakeWSO2RequestReturnResponse(method, url string, body interface{}, toReturn interface{}) (*nerr.E, *http.Response) {
	err, response := MakeWSO2RequestWithHeadersReturnResponse(method, url, body, toReturn, nil)
	return err, response
}

//MakeWSO2RequestWithHeaders makes a generic WSO2 request with headers
//toReturn should be a pointer
func MakeWSO2RequestWithHeaders(method, url string, body interface{}, toReturn interface{}, headers map[string]string) *nerr.E {
	err, _ := MakeWSO2RequestWithHeadersReturnResponse(method, url, body, toReturn, headers)
	return err
}

//MakeWSO2RequestWithHeadersReturnResponse makes a generic WSO2 request with headers
//toReturn should be a pointer
func MakeWSO2RequestWithHeadersReturnResponse(method, url string, body interface{}, toReturn interface{}, headers map[string]string) (*nerr.E, *http.Response) {
	//	log.L.Debugf("Making %v request against %v at %v", method, url, time.Now())

	key, er := GetAccessKey()
	if er != nil {
		return er.Addf("Couldn't make WSO2 request"), nil
	}

	//attach key
	var b []byte
	var ok bool
	var err error

	if body != nil {
		if b, ok = body.([]byte); !ok {
			b, err = json.Marshal(body)
			if err != nil {
				return nerr.Translate(err).Addf("Couldn't marhsal request"), nil
			}
			log.L.Debugf("Sending %s to WSO2", b)
		}
	}

	for {
		hasRetried := false

		req, err := http.NewRequest(method, url, bytes.NewBuffer(b))
		if err != nil {
			return nerr.Translate(err).Addf("Couldn't build WSO2 request"), nil
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
			return nerr.Translate(err).Addf("Couldn't make WSO2 request"), resp
		}
		defer resp.Body.Close()

		rb, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nerr.Translate(err).Addf("Couldn't read response body"), resp
		}

		//log.L.Debugf("Response body: %s", rb)

		//log.L.Debugf("Response Headers: %v", resp.Header)

		if resp.StatusCode/100 != 2 {
			if resp.StatusCode == 400 && len(rb) == 0 && !hasRetried {
				//if we get a 400 and a blank body and we haven't retried, then just try again
				log.L.Debugf("400 and blank body - retrying WS02 request")
				hasRetried = true
				continue
			}
			return nerr.Create(fmt.Sprintf("Non 200: body [%s] Response code: [%v]", rb, resp.StatusCode), "request-error"), resp
		}

		err = json.Unmarshal(rb, toReturn)
		if err != nil {
			return nerr.Translate(err).Addf("Couldn't unmarshal response %s", "unmarshal error"), resp
		}

		return nil, nil
	}
}
