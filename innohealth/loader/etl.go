package main

import (
	"fmt"
	"log"
	"sort"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	OperationCollection = "operation"
)

var (
	mgoSession *mgo.Session
	dbName     string
)

//CloseMgo is ...
func CloseMgo() {
	if mgoSession != nil {
		mgoSession.Close()
		mgoSession = nil
	}
}

//InitMgo is ..
func InitMgo(dbname string) {
	CloseMgo()
	dbName = dbname
	connectURL := fmt.Sprintf("mongodb://testuser:innohealth@innohs.com/%s", dbname)
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

// UpdateOPeration is ...
func UpdateOPeration(data []*Operation) {

	session := mgoSession.Clone()
	defer session.Close()

	db := session.DB(dbName)

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
