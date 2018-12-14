package calendar

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
)

const (
	//Mondayinstruction .
	Mondayinstruction = "monday-instruction"
	//Fridayinstruction .
	Fridayinstruction = "friday-instruction"
	//Holiday .
	Holiday = "holiday"
	//Noclass .
	Noclass = "noclass"
)

//Exception .
type Exception struct {
	StartDate ExceptionTime `json:"StartDateTime"`
	EndDate   ExceptionTime `json:"EndDateTime"`
	Category  string        `json:"CategoryName"`
	Title     string        `json:"Title"`
}

//ExceptionTime .
type ExceptionTime time.Time

const (
	exceptionDateFormat = "2006-01-02 15:04:05"
)

//UnmarshalJSON .
func (m *ExceptionTime) UnmarshalJSON(p []byte) error {
	t, err := time.Parse(exceptionDateFormat, strings.TrimSpace(strings.Replace(
		string(p),
		"\"",
		"",
		-1,
	)))

	if err != nil {
		return err
	}

	*m = ExceptionTime(t)

	return nil
}

var exceptions map[string][]Exception

var lastCheck time.Time
var interval = 48 * time.Hour
var checkMutex sync.Mutex

func init() {
	checkMutex = sync.Mutex{}
}

//CheckExceptionsForDate .
func CheckExceptionsForDate(date time.Time) ([]Exception, *nerr.E) {
	checkMutex.Lock()
	//check to see if we've updated since then
	if time.Now().Add(-1 * interval).After(lastCheck) {
		e, err := getExceptions()
		if err != nil {
			return []Exception{}, err.Addf("Couldn't check exceptions...")
		}
		exceptions = e
	}
	checkMutex.Unlock()

	//we check to see if the date we're
	if v, ok := exceptions[date.Format(timeformat)]; ok {
		return v, nil
	}

	return []Exception{}, nil
}

var codes = map[int]string{
	420: Noclass,
	419: Mondayinstruction,
	417: Fridayinstruction,
	402: Holiday,
}

var timeformat = "2006-01-02"

func getExceptions() (map[string][]Exception, *nerr.E) {

	//make our api call(s)
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	toReturn := map[string][]Exception{}

	for k, v := range codes {
		log.L.Debugf("getting exceptions for type %v-%v", k, v)

		addr := fmt.Sprintf("https://calendar.byu.edu/api/Events.json?categories=%v", k)

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
			return map[string][]Exception{}, nerr.Create(fmt.Sprintf("Couldn't get events for type %v. Response %s", v, b), "non-200")
		}

		err = json.Unmarshal(b, &events)
		if err != nil {
			return map[string][]Exception{}, nerr.Translate(err).Addf("Couldn't get events for type %v. Unkown response %s", v, b)
		}

		for i := range events {
			events[i].Category = v
			toReturn[time.Time(events[i].StartDate).Format(timeformat)] = append(toReturn[time.Time(events[i].StartDate).Format(timeformat)], events[i])
		}
	}

	lastCheck = time.Now()

	return toReturn, nil
}
