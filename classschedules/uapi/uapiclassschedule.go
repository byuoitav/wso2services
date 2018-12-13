package uapi

import (
	"fmt"
	"strings"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/wso2services/classschedules/calendar"
	"github.com/byuoitav/wso2services/wso2requests"
)

//cache by yearterm -> room

var classScheduleCache map[string]map[string][]ClassSchedule
var updateTimesByBuilding map[string]time.Time
var ttl time.Duration = (24 * time.Hour) * -1

func init() {
	updateTimesByBuilding = map[string]time.Time{}
	classScheduleCache = map[string]map[string][]ClassSchedule{}

}

//GetClassScheduleForRoomAndTime .
func GetClassScheduleForRoomAndTime(roomname string, classtime time.Time) (ClassSchedule, *nerr.E) {
	//first thing is to get the schedule

	_, err := getClassSchedule(roomname, classtime)
	if err != nil {
		return ClassSchedule{}, err.Addf("Couldn't get current class information")
	}

	//we need to go through and figure out what schedule we're looking for - check to see if it's a special day
	return ClassSchedule{}, nil
}

func getClassSchedule(roomname string, classtime time.Time) ([]ClassSchedule, *nerr.E) {

	t, err := calendar.GetYearTermForDate(classtime)
	if err != nil {
		return []ClassSchedule{}, err.Addf("Couldn't get class schedule for room %v and time %v", roomname, classtime)
	}

	rmsplit := strings.Split(roomname, "-")

	//check to see if we have the class schedule cached for that term
	if termmap, ok := classScheduleCache[t.YearTerm]; ok {
		if time.Now().Add(ttl).Before(updateTimesByBuilding[rmsplit[0]]) {
			//check for the cache
			if class, ok := termmap[roomname]; ok {
				//check to see if it's up to date
				return class, nil
			}
			return []ClassSchedule{}, nerr.Create(fmt.Sprintf("Cannot get schedule for room %v, even after fetch", roomname), "invalid-room")
		}
	} else {
		//nothing for this term, we need to initialze the map
		classScheduleCache[t.YearTerm] = map[string][]ClassSchedule{}
	}

	classes, err := fetchClassSchedules(roomname, t.YearTerm)
	if err != nil {
		return []ClassSchedule{}, err.Addf("Couldn't fetch class schedule")
	}

	updateTimesByBuilding[rmsplit[0]] = time.Now()

	m := classScheduleCache[t.YearTerm]

	for i := range classes {
		//we go through and update the map
		rmname := fmt.Sprintf("%v-%v", classes[i].AssignedSchedules.Values[0].Building.Value, classes[i].AssignedSchedules.Values[0].Room.Value)

		m[rmname] = append(m[rmname], classes[i])
	}
	classScheduleCache[t.YearTerm] = m

	if v, ok := classScheduleCache[t.YearTerm][roomname]; ok {
		return v, nil
	}

	return []ClassSchedule{}, nerr.Create(fmt.Sprintf("Cannot get schedule for room %v, even after fetch", roomname), "invalid-room")

}

//assume roomname is in format BUILDING-ROOM
func fetchClassSchedules(roomname, term string) ([]ClassSchedule, *nerr.E) {

	//we figure out the building
	br := strings.Split(roomname, "-")
	var toReturn []ClassSchedule

	var resp ClassResponse

	err := wso2requests.MakeWSO2Request("GET", fmt.Sprintf("https://api.byu.edu/byuapi/classes/v1?year_term=%v&building=%v&context=class_schedule", term, br[0]), []byte{}, &resp)

	if err != nil {
		return toReturn, err.Addf("Couldn't fetch class scheudle")
	}

	for i := range resp.Values {
		toReturn = append(toReturn, resp.Values[i])
	}

	for resp.Metadata.PageEnd < resp.Metadata.CollectionSize {

		err := wso2requests.MakeWSO2Request("GET", fmt.Sprintf("https://api.byu.edu/byuapi/classes/v1?year_term=%v&building=%v&context=class_schedule&page_start=%v", term, br[0], resp.Metadata.PageEnd+1), []byte{}, &resp)

		if err != nil {
			return toReturn, err.Addf("Couldn't fetch class scheudle")
		}

		for i := range resp.Values {
			toReturn = append(toReturn, resp.Values[i])
		}
		log.L.Debugf("Have %v classes", len(toReturn))
	}

	return toReturn, nil
}
