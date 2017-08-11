package main

import (
	"errors"
	"fmt"
	"log"
	"sort"

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
func UpdateOperation(dbName string, data []*Operation) {

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

// QueryOperation is
func QueryOperation(dbName string) ([]Operation, error) {
	session := mgoSession.Clone()
	defer session.Close()

	db := session.DB(dbName)
	collection := db.C(OperationCollection)
	if !CheckIfCollection(db, OperationCollection) {
		msg := fmt.Sprintf("Collection[%s] does not exits.", OperationCollection)
		log.Println(msg)
		return nil, errors.New(msg)
	}
	result := make([]Operation, 0)
	err := collection.Find(bson.M{}).Limit(10).All(&result)
	if err != nil {
		return nil, errors.New("Query is wrong")
	}

	return result, nil
}
