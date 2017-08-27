package main

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	// OperationCollection is name
	OperationCollection = "operation"
)

var (
	mgoSession *mgo.Session
)

//CloseMgo is ...
func CloseMgo() {
	if mgoSession != nil {
		mgoSession.Close()
		mgoSession = nil
	}
}

//InitMgo is ..
func InitMgo() {
	CloseMgo()
	connectURL := "mongodb://apiuser:resuipa@innohs.com"
	var err error
	mgoSession, err = mgo.Dial(connectURL)
	if err != nil {
		log.Fatal(err)
	}
	mgoSession.SetMode(mgo.Monotonic, true)
}

// CheckIfCollection is ...
func CheckIfCollection(db *mgo.Database, name string) bool {
	names, err := db.CollectionNames()
	if err != nil {
		log.Fatal(err)
	}
	sort.Strings(names)
	i := sort.SearchStrings(names, name)

	if i < len(names) {
		return true
	}
	return false
}

// UpdateOperation is ...
func UpdateOperation(dbname string, data []*Operation) {

	session := mgoSession.Clone()
	defer session.Close()

	db := session.DB(dbname)

	collection := db.C(OperationCollection)
	if !CheckIfCollection(db, OperationCollection) {
		collection.Create(&mgo.CollectionInfo{})
	}
	count := 0
	for _, item := range data {
		collection.Upsert(bson.M{"SID": item.SID}, *item)
		count++
		if (count % 10) == 0 {
			fmt.Printf("count - %05d\n", count)
		}
	}
}

func doQuery(db *mgo.Database, collection string, query bson.M) (*mgo.Query, error) {
	if !CheckIfCollection(db, collection) {
		msg := fmt.Sprintf("Collection[%s] does not exits.", collection)
		log.Println(msg)
		return nil, errors.New(msg)
	}
	c := db.C(collection)
	return c.Find(bson.M{}), nil
}

func queryOperation(db *mgo.Database, start int, limit int) (map[string]interface{}, error) {
	query, _ := doQuery(db, OperationCollection, bson.M{})
	total, _ := query.Count()
	result := make([]Operation, 0)

	if start > 0 {
		query.Skip(start)
	}
	var err error
	if limit < 1 {
		err = query.All(&result)
	} else {
		err = query.Limit(limit).All(&result)
	}

	if err != nil {
		return nil, errors.New("Query is wrong")
	}

	var output = map[string]interface{}{
		"total": total,
		"start": start,
		"limit": limit,
		"data":  result,
	}
	return output, nil
}

// QueryOperation is
func QueryOperation(dbname string, start int, limit int) (map[string]interface{}, error) {
	session := mgoSession.Clone()
	defer session.Close()

	db := session.DB(dbname)
	return queryOperation(db, start, limit)
}

func doPipe(dbname string, pipes []bson.M) ([]bson.M, error) {
	session := mgoSession.Clone()
	defer session.Close()

	db := session.DB(dbname)
	collection := db.C(OperationCollection)
	if !CheckIfCollection(db, OperationCollection) {
		msg := fmt.Sprintf("Collection[%s] does not exits.", OperationCollection)
		log.Println(msg)
		return nil, errors.New(msg)
	}
	query := collection.Pipe(pipes)
	resp := make([]bson.M, 0)
	query.All(&resp)
	return resp, nil
}

// OperationByDoctorStat returns
func OperationByDoctorStat(dbname string) (interface{}, error) {
	pipeline := []bson.M{
		bson.M{"$group": bson.M{"_id": "$doctorname", "count": bson.M{"$sum": 1}}},
	}
	result, err := doPipe(dbname, pipeline)
	if err != nil {
		return nil, err
	}

	stat := make([]DoctorOperatoionStat, 0)
	for _, submap := range result {
		name := submap["_id"].(string)
		count := submap["count"].(int)
		stat = append(stat, DoctorOperatoionStat{name, count})
	}
	return stat, nil
}

// OperationByWeekdayStat returns group by weekday, department
func OperationByWeekdayStat(dbname string) (interface{}, error) {

	pipeline := []bson.M{
		bson.M{"$group": bson.M{"_id": bson.M{"weekday": "$dayofweek", "department": "$departmentname"}, "count": bson.M{"$sum": 1}}},
	}
	response, err := doPipe(dbname, pipeline)
	if err != nil {
		return nil, err
	}

	type output map[string]interface{}
	result := make([]output, 0)
	for _, submap := range response {
		item := make(output)
		idmap := submap["_id"].(bson.M)
		item["department"] = idmap["department"].(string)
		item["weekday"] = RevDayOfWeek[time.Weekday(idmap["weekday"].(int))]
		item["count"] = submap["count"].(int)
		result = append(result, item)
	}
	return result, nil
}

// GetDistictValue returns ..
func GetDistictValue(db *mgo.Database, field string) (interface{}, error) {
	query, err := doQuery(db, OperationCollection, bson.M{})
	if err != nil {
		return nil, err
	}
	var result interface{}
	query.Distinct(field, &result)

	return result, nil
}

// OperationByTimStat returns group by weekday, department
func OperationByTimStat(dbname string) (interface{}, error) {
	session := mgoSession.Clone()
	defer session.Close()

	db := session.DB(dbname)
	roominfo, err := GetDistictValue(db, "operationroom")
	if err != nil {
		return nil, err
	}

	length := len(roominfo.([]interface{}))
	xcategory := make([]string, length)
	for index, room := range roominfo.([]interface{}) {
		xcategory[index] = room.(string)
	}

	sort.Strings(xcategory)
	// fmt.Println(xcategory)
	roomMapping := make(map[string]int)
	for index, item := range xcategory {
		roomMapping[item] = index
	}
	numTimeSlot := 24 * 2 // 30 minute granularity
	ycategory := make([]string, numTimeSlot)
	timeMapping := make(map[InnoTime]int)

	current := InnoTime{true, 0, 0}
	for i := numTimeSlot - 1; i > -1; i-- {
		ycategory[i] = current.String()
		timeMapping[current] = i
		current.delta(InnoTime{Minute: 30})
	}

	result := make([][]int, numTimeSlot)
	for i := 0; i < numTimeSlot; i++ {
		result[i] = make([]int, length)
		for j := 0; j < length; j++ {
			result[i][j] = 0
		}
	}

	output, _ := queryOperation(db, 0, 10)
	for _, item := range output["data"].([]Operation) {
		//fmt.Println(item)
		if item.RoomEntranceTime.Exist && item.RoomOutTime.Exist {
			timeRange := item.RoomEntranceTime.EnumerateTo(item.RoomOutTime, 30)
			for _, timeSlot := range timeRange {
				result[timeMapping[timeSlot]][roomMapping[item.OperationRoom]]++
			}

		}
	}
	total := numTimeSlot * length
	heatdata := make([][]int, total)
	for j := 0; j < numTimeSlot; j++ {
		for i := 0; i < length; i++ {
			heatdata[j*length+i] = []int{i, j, result[j][i]}
		}
	}
	response := make(map[string]interface{})
	response["xCategory"] = xcategory
	response["yCategory"] = ycategory
	response["data"] = heatdata
	return response, nil
}
