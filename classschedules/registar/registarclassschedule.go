package registar

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/byuoitav/common/log"
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

	startTime time.Time //the zero time with the hours/minutes of start of class set
	endTime   time.Time //the zero time with the hours/minutes of end of class set
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

func init() {

	cache = map[string]map[string]Room{}
	updateTimes = map[string]time.Time{}

	classtimereg = regexp.MustCompile(datere)
}

var cache map[string]map[string]Room //yearterm -> room
var updateTimes map[string]time.Time
var classtimereg *regexp.Regexp

//GetClassScheduleForTime .
func GetClassScheduleForTime(roomname string, classtime time.Time) (ClassSchedule, *nerr.E) {

	t, err := calendar.GetYearTermForDate(classtime)
	if err != nil {
		return ClassSchedule{}, err.Addf("Couldn't get class schedule for room %v and time %v", roomname, classtime)
	}

	//check to see if we have it in the cache
	term, ok := cache[t.YearTerm]
	if !ok {

		rm, err := fetchClassSchedule(roomname, t.YearTerm)
		if err != nil {
			updateTimes[roomname] = time.Now()
			return ClassSchedule{}, err.Addf("Couldn't get class schedule for time.")
		}
		//we need to fetch the whole thing
		term = map[string]Room{
			roomname: rm,
		}
		cache[t.YearTerm] = term
		updateTimes[roomname] = time.Now()
	}

	//check for the room
	rm, ok := term[roomname]
	if !ok {
		//check to see when the last time we checked was, if it was less than 30 minutes we don't go try again, we just assume it's a bad room
		ck, ok := updateTimes[roomname]
		if !ok || time.Now().Sub(ck) > 5*time.Minute {

			//we go get it .
			updateTimes[roomname] = time.Now()
			rm, err := fetchClassSchedule(roomname, t.YearTerm)
			if err != nil {
				return ClassSchedule{}, err.Addf("Couldn't get class schedule for time.")
			}
			term[roomname] = rm
		}
	} else {
		//we need to check update times here, too
		ck, ok := updateTimes[roomname]
		if !ok || time.Now().Sub(ck) > 5*time.Minute {

			//we go get it .
			rm, err = fetchClassSchedule(roomname, t.YearTerm)
			updateTimes[roomname] = time.Now()
			if err != nil {
				return ClassSchedule{}, err.Addf("Couldn't get class schedule for time.")
			}
			term[roomname] = rm
		}
	}

	//rm will be set here.

	//figure out what class classtime falls into
	//check to see if there are any exceptions for that date
	exceptions, err := calendar.CheckExceptionsForDate(classtime)
	if err != nil {
		return ClassSchedule{}, err.Addf("Couldn't get class schedule for time.")
	}

	if len(exceptions) > 0 {
		//we have some exeptions, check to see what kind?
		return FindScheduleWithExceptions(exceptions[0], rm, classtime)
	}
	return FindSchedule(rm, classtime, classtime.Weekday())
}

//FindScheduleWithExceptions .
func FindScheduleWithExceptions(exception calendar.Exception, rm Room, t time.Time) (ClassSchedule, *nerr.E) {

	switch exception.Category {
	case calendar.Mondayinstruction:
		return FindSchedule(rm, t, time.Monday)
	case calendar.Fridayinstruction:
		return FindSchedule(rm, t, time.Friday)
	default: //holiday or no class
		return ClassSchedule{}, nil
	}

}

//FindSchedule .
func FindSchedule(rm Room, t time.Time, weekday time.Weekday) (ClassSchedule, *nerr.E) {
	log.L.Debugf("Finding schedule for %v at %v. it's a %v", rm.Building+rm.Room, t, weekday)
	log.L.Debugf("%v has %v classes", rm.Building+rm.Room, len(rm.Schedules))

	for _, v := range rm.Schedules {
		wkdays := FindDaysOfWeek(v.Days)

		var ok bool
		for i := range wkdays {
			if wkdays[i] == weekday {
				ok = true
				break
			}
		}
		if !ok {
			continue
		}
		log.L.Debugf("Class %v occurs on %v", v.DeptName+v.CatalogNumber, weekday)

		//we're on the right day of the week, check to see if we're in the right time
		dayclassStart := time.Date(t.Year(), t.Month(), t.Day(), v.startTime.Hour(), v.startTime.Minute(), 0, 0, time.UTC)
		dayclassStart = dayclassStart.Add(-5 * time.Minute) //for a buffer

		log.L.Debugf("Class %v started at %v ", v.DeptName+v.CatalogNumber, dayclassStart)

		if t.Before(dayclassStart) {
			continue
		}
		//we're on the right day of the week, check to see if we're in the right time
		dayclassEnd := time.Date(t.Year(), t.Month(), t.Day(), v.endTime.Hour(), v.endTime.Minute(), 0, 0, time.UTC)
		dayclassEnd = dayclassEnd.Add(5 * time.Minute) //for a buffer

		log.L.Debugf("Class %v ended at %v ", v.DeptName+v.CatalogNumber, dayclassEnd)

		if t.After(dayclassEnd) {
			continue
		}

		//bingo.
		return v, nil
	}

	//didn't find a class

	return ClassSchedule{}, nil
}

//FindDaysOfWeek .
func FindDaysOfWeek(wkstr string) []time.Weekday {
	toReturn := []time.Weekday{}

	for i := 0; i < len(wkstr); i++ {
		switch wkstr[i] {
		case 'M':
			toReturn = append(toReturn, time.Monday)
		case 'T':
			if i+1 != len(wkstr) && wkstr[i+1] == 'h' {
				toReturn = append(toReturn, time.Thursday)
				i++
			} else {
				toReturn = append(toReturn, time.Tuesday)
			}
		case 'W':
			toReturn = append(toReturn, time.Wednesday)
		case 'F':
			toReturn = append(toReturn, time.Friday)
		case 'S':
			if i+1 != len(wkstr) && wkstr[i+1] == 'a' {
				toReturn = append(toReturn, time.Saturday)
				i++
			} else {
				toReturn = append(toReturn, time.Saturday)
			}
		}
	}

	return toReturn
}

const datere = `(\d{1,2}):(\d{2})([ap]) - (\d{1,2}):(\d{2})([ap])`

func fetchClassSchedule(roomname, term string) (Room, *nerr.E) {

	//we figure out the building
	br := strings.Split(roomname, "-")

	var resp ClassScheduleResponse

	err := wso2requests.MakeWSO2Request("GET", fmt.Sprintf("https://api.byu.edu:443/domains/legacy/academic/classschedule/classroom/v1/%v/%v/%v/Schedule", term, br[0], br[1]), []byte{}, &resp)

	if err != nil {
		return resp.ClassRoomService.Room, err.Addf("Couldn't fetch class scheudle")
	}

	//for each time, we'll need to go and parse the slasstime to figure out the deal
	rm := resp.ClassRoomService.Room

	for i := range rm.Schedules {
		matches := classtimereg.FindAllStringSubmatch(rm.Schedules[i].ClassTime, -1)
		if len(matches) != 1 {
			msg := fmt.Sprintf("Unknown class format %v", rm.Schedules[i].ClassTime)
			return rm, nerr.Create(msg, "unknown-format")
		}

		shr, err := strconv.Atoi(matches[0][1])
		if err != nil {
			msg := fmt.Sprintf("Unknown class format %v", rm.Schedules[i].ClassTime)
			return rm, nerr.Create(msg, "unknown-format")
		}

		smin, err := strconv.Atoi(matches[0][2])
		if err != nil {
			msg := fmt.Sprintf("Unknown class format %v", rm.Schedules[i].ClassTime)
			return rm, nerr.Create(msg, "unknown-format")
		}

		rm.Schedules[i].startTime = time.Time{}.Add(time.Duration(shr)*time.Hour + time.Duration(smin)*time.Minute)
		if matches[0][3] == "p" {
			rm.Schedules[i].startTime = rm.Schedules[i].startTime.Add(12 * time.Hour)
		}

		ehr, err := strconv.Atoi(matches[0][4])
		if err != nil {
			msg := fmt.Sprintf("Unknown class format %v", rm.Schedules[i].ClassTime)
			return rm, nerr.Create(msg, "unknown-format")
		}

		emin, err := strconv.Atoi(matches[0][5])
		if err != nil {
			msg := fmt.Sprintf("Unknown class format %v", rm.Schedules[i].ClassTime)
			return rm, nerr.Create(msg, "unknown-format")
		}

		rm.Schedules[i].endTime = time.Time{}.Add(time.Duration(ehr)*time.Hour + time.Duration(emin)*time.Minute)
		if matches[0][6] == "p" {
			rm.Schedules[i].endTime = rm.Schedules[i].endTime.Add(12 * time.Hour)
		}
	}

	return resp.ClassRoomService.Room, nil
}
