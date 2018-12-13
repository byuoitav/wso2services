package classschedules

import (
	"fmt"
	"strings"
	"time"

	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/wos2services/wso2requests"
)

//cache by yearterm -> room

var classScheduleCache map[string]map[string][]UAPIClassSchedule
var updateTimesByBuildiing map[string]time.Time
var ttl time.Duration = (24 * time.Hour)

//GetClassScheduleForRoomAndTime .
func GetClassScheduleForRoomAndTime(roomname string, classtime time.Time) (UAPIClassSchedule, *nerr.E) {
	//first thing is to get the schedule

	classes, err := getClassSchedule(roomane, classtime)
	if err != nil {
		return schedule, err.Addf("Couldn't get current class information")
	}

	//we need to go through and figure out what schedule we're looking for - check to see if it's a special day
}

func getClassSchedule(roomname string, classtime time.Time) ([]UAPIClassSchedule, *nerr.E) {

	t, err := GetYearTermForDate(classtime)
	if err != nil {
		return UAPIClassSchedul{}, err.Addf("Couldn't get class schedule for room %v and time %v", roomname, classtime)
	}

	bu, rm := strings.Split(roomname, "-")

	//check to see if we have the class schedule cached for that term
	if termmap, ok := classScheduleCache[t]; ok {
		if time.Now().Sub(ttl).Before(updateTimesByBuilding[bu]) {
			//check for the cache
			if class, ok := termmap[roomname]; ok {
				//check to see if it's up to date
				return class, nil
			}
			return UAPIClassSchedule{}, nerr.Create(fmt.Sprintf("Cannot get schedule for room %v, even after fetch", roomname), "invalid-room")
		}
	} else {
		//nothing for this term, we need to initialze the map
		classScheduleCache[t] = map[string]UAPIClassSchedule{}
	}

	classes, err := fetchClassSchedules(roomname, t)
	if err != nil {
		return UAPIClassSchedule, err.Addf("Couldn't fetch class schedule")
	}

	updatTimeByBuilding[bu] = time.Now()

	m = classScheduleCache[t]

	for i := range classes {
		//we go through and update the map
		rmname := fmt.Sprintf("%v-%v", classes[i].AssignedSchedules.Values.Building.Value, classes[i].AssignedSchedules.Values.Room.Value)

		m[rmname] = append(m[rmname], classes[i])
	}
	classScheduleCache[t] = m

	if v, ok := classScheduleCache[t][roomname]; ok {
		return v, nil
	}

	return UAPIClassSchedule{}, nerr.Create(fmt.Sprintf("Cannot get schedule for room %v, even after fetch", roomname), "invalid-room")

}

//assume roomname is in format BUILDING-ROOM
func fetchClassSchedules(roomname, term string) ([]UAPIClassSchedule, *nerr.E) {

	//we figure out the building
	br := strings.Split(roomname, "-")

	var resp UAPIClassResponse

	err := wso2requests.MakeWSO2Request("GET", fmt.Sprintf("https://api.byu.edu/byuapi/classes/v1?year_term=%v&building=%v&context=class_schedule", term, br[0]), []byte{}, &resp)

	if err != nil {
		return toReturn, err.Addf("Couldn't fetch class scheudle")
	}

	var toReturn []UAPIClassSchedule
	for i := range resp.Values {
		toReturn = append(toReturn, resp.Values[i])
	}

	for resp.PageEnd < resp.CollectionSize {

		err := wso2requests.MakeWSO2Request("GET", fmt.Sprintf("https://api.byu.edu/byuapi/classes/v1?year_term=%v&building=%v&context=class_schedule&page_start=%v", term, br[0], resp.PageEnd+1), []byte{}, &resp)

		if err != nil {
			return toReturn, err.Addf("Couldn't fetch class scheudle")
		}

		for i := range resp.Values {
			toReturn = append(toReturn, resp.Values[i])
		}
	}

	return toReturn, nil
}
