package registar

import (
	"fmt"
	"strings"
	"time"

	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/wso2services/classschedules/calendar"
	"github.com/byuoitav/wso2services/wso2requests"
)

//ClassScheduleResponse .
type ClassScheduleResponse struct {
	ClassRoomService struct {
		Request struct {
			Method        string      `json:"method"`
			Resource      string      `json:"resource"`
			Attributes    interface{} `json:"attributes"`
			Status        int         `json:"status"`
			StatusMessage string      `json:"statusMessage"`
		} `json:"request"`
		Room Room `json:"response"`
	} `json:"ClassRoomService"`
}

//ClassSchedule .
type ClassSchedule struct {
	DeptName       string      `json:"dept_name"`
	CatalogNumber  string      `json:"catalog_number"`
	CatalogSuffix  interface{} `json:"catalog_suffix"`
	LabQuizSection interface{} `json:"lab_quiz_section"`
	Honors         interface{} `json:"honors"`
	ServLearning   interface{} `json:"serv_learning"`
	CreditHours    float64     `json:"credit_hours"`
	SectionType    interface{} `json:"section_type"`
	ClassTime      string      `json:"class_time"`
	Days           string      `json:"days"`
	InstructorName string      `json:"instructor_name"`
	SectionSize    int         `json:"section_size"`
	TotalEnr       int         `json:"total_enr"`
	SchedType      string      `json:"sched_type"`
	AssignTo       interface{} `json:"assign_to"`
	StartDate      interface{} `json:"start_date"`
	EndDate        interface{} `json:"end_date"`
}

//Room .
type Room struct {
	Room              string          `json:"room"`
	Building          string          `json:"building"`
	SchedulingType    string          `json:"scheduling_type"`
	RoomDesc          string          `json:"room_desc"`
	Capacity          int             `json:"capacity"`
	EffectiveYearTerm string          `json:"effective_year_term"`
	ExpiredYearTerm   string          `json:"expired_year_term"`
	ROOMDATA          interface{}     `json:"ROOMDATA"`
	Schedules         []ClassSchedule `json:"schedules"`

	lastUpdate time.Time //when we last fetched this
}

const (
	registarRecordTTL = -24 * time.Hour // once a day
)

var cache map[string]map[string]Room //yearterm -> room
var updateTimes map[string]time.Time

//GetClassScheduleFortime .
func GetClassScheduleFortime(roomname string, classtime time.Time) (ClassSchedule, *nerr.E) {

	t, err := calendar.GetYearTermForDate(classtime)
	if err != nil {
		return ClassSchedule{}, err.Addf("Couldn't get class schedule for room %v and time %v", roomname, classtime)
	}

	var room Room
	//check to see if we have it in the cache
	term, ok := cache[t.YearTerm]
	if ok {
		room, ok = term[roomname]
		if !ok {
			//check to see if we've tried recently, only try once every 60 minutes if its not there...
			if v, ok := updateTimes[roomname]; ok {
				if time.Now().Add(-60 * time.Minute).After(v) { //it's been m

				}
			}
		}
	} else {
		cache[t.YearTerm] = map[string]Room{}
	}

	return ClassSchedule{}, nil

}

func fetchClassSchedule(roomname, term string) (Room, *nerr.E) {

	//we figure out the building
	br := strings.Split(roomname, "-")

	var resp ClassScheduleResponse

	err := wso2requests.MakeWSO2Request("GET", fmt.Sprintf("https://api.byu.edu:443/domains/legacy/academic/classschedule/classroom/v1/%v/%v/%v/Schedule", term, br[0], br[1]), []byte{}, &resp)

	if err != nil {
		return resp.ClassRoomService.Room, err.Addf("Couldn't fetch class scheudle")
	}

	return resp.ClassRoomService.Room, nil
}
