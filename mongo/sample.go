package main

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Feed struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Product     string        `bson:product`
	Price       uint64        `bson:price`
	Fmp         uint64        `bson:fmp`
	Description string        `bson:description`
	Source      string        `bson:source`
	Brand       string        `bson:brand`
	Url         string        `bson:url`
	Image       string        `bson:image`
	Status      bool          `bson:status`
	Promotion   string        `bson:promotion`
	Brandid     uint64        `bson:brandid`
	Category    string        `bson:category`
}

var (
	IsDrop = false
)

func main() {
	session, err := mgo.Dial("brandflask.com")
	if err != nil {
		panic(err)
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	// Drop Database
	if IsDrop {
		err = session.DB("test").DropDatabase()
		if err != nil {
			panic(err)
		}
	}

	// Collection People
	c := session.DB("brandflask").C("feeddb")

	var results []Feed
	err = c.Find(bson.M{"$and":bson.M[bson.M{description:bson.RegEx{"/boho/"}},
    bson.M{category: bson.M{"$nin": bson.M[bson.RegEx("/Accessories/"),
      bson.RegEx("/Bikini/"),bson.RegEx("/Activewear/"),bson.RegEx("/Swimwear/")]}},
    bson.M{"brandid":{$ne:0}}]}).All(&results)

	if err != nil {
		panic(err)
	}
	fmt.Println("Results All: ", results)
}
