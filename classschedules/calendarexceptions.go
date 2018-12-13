package classschedules

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/byuoitav/common/nerr"
)

const (
	mondayinstruction = "monday-instruction"
	fridayinstruction = "friday-instruction"
	holiday           = "holiday"
	noclass           = "noclass"
)

//Exception .
type Exception struct {
	StartDate time.Time `json:"StartDateTime"`
	EndDate   time.Time `json:"EndDateTime"`
	Category  string    `json:"CategoryName"`
	Title     string    `json:"Title"`
}

var exceptions map[string][]Exception

var lastCheck time.Time
var interval = 48 * time.Hour
var checkMutex sync.Mutex

func init() {
	checkMutex = sync.Mutex{}
}

//CheckExceptionsForDate .
func CheckExceptionsForDate(date time.Time) (Exception, *nerr.E) {
	checkMutex.Lock()
	//check to see if we've updated since then
	if time.Now().Sub(interval).After(lastCheck) {
		e, err := getExceptions()
		if err != nil {
			return Exception{}, err.Addf("Couldn't check eXceptions...")
		}
		exceptions = e
	}
	checkMutex.Unlock()

	//we check to see if the date we're
	if v, ok := exceptions[date.Format(timeformat)]; ok {
		return v, nil
	}

	return Exception{}, nil
}

var codes = map[int]string{
	420: noclass,
	419: mondayinstruction,
	417: fridayinstruction,
	402: holiday,
}

var timeformat = "2006-01-02"

func getExceptions() (map[string][]Exception, *nerr.E) {

	//make our api call(s)
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	toReturn := map[string][]Exception{}

	for k, v := range codes {

		addr := fmt.Sprintf("https://calendar.byu.edu/apii/Events.json?categories=%v", k)

		resp, err := client.Get(addr)
		if err != nil {
			return map[string][]Exception{}, nerr.Translate(err).Addf("Couldn't get events for type %v", v)
		}

		defer resp.Body.Close()
		var events []Exception

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return map[string][]Exception{}, nerr.Translate(err).Addf("Couldn't get events for type %v", v)
		}

		if resp.StatusCode/100 != 2 {
			return map[string][]Exception{}, nerr.Create(fmt.Spritnf("Couldn't get events for type %v. Response %s", v, b), "non-200")
		}

		err := json.Unmarshal(b, &events)
		if err != nil {
			return map[string][]Exception{}, nerr.Translate(err).Addf("Couldn't get events for type %v. Unkown response %s", v, b)
		}

		for i := range events {
			events[i].Category = v
			toReturn[events.StartDate.Format(timeformat)] = append(toReturn[events.StartDate.Format(timeformat)], events[i])
		}
	}

	lastCheck = time.Now()

	return toReturn, nil
}
