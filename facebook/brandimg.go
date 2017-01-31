package main

import (
	"fmt"
	"strings"
	"time"

	fb "github.com/huandu/facebook"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var accessToken = "1380374468852725|d5zZcwcUyaF7i1VMUa31kD1Pz6k"

func getImages(fbid string) string {
	res, _ := fb.Get(fmt.Sprintf("/v2.8/%s/", fbid), fb.Params{
		"fields":       "images",
		"access_token": accessToken})

	large := res["images"].([]interface{})[0].(map[string]interface{})
	return large["source"].(string)
}

func retrieve(data interface{}) []string {
	var result []string
	for _, v := range data.([]interface{}) {
		id := v.(map[string]interface{})["id"].(string)
		result = append(result, getImages(id))
	}
	return result
}

func album(fbid string) []string {
	// fmt.Println(fbid)
	res, _ := fb.Get(fmt.Sprintf("/v2.8/%s/photos?limit=100", fbid), fb.Params{
		"access_token": accessToken})
	return retrieve(res.GetField("data"))
}

func retrieveImages(fbid string) []string {
	res, _ := fb.Get(fmt.Sprintf("/v2.8/%s", fbid), fb.Params{
		"fields":       "id,name,albums",
		"access_token": accessToken})

	if res["albums"] == nil {
		fmt.Println("skip")
		return []string{}
	}

	data := res["albums"].(map[string]interface{})["data"]
	for _, v := range data.([]interface{}) {
		obj := v.(map[string]interface{})
		name := obj["name"].(string)
		if strings.Compare(name, "Timeline Photos") == 0 {
			id := obj["id"].(string)
			images := album(id)
			// fmt.Println(id)
			// fmt.Println(len(images))
			return images
		}
	}
	var result = make([]string, 100)

	return result
}

func main() {
	var mongohost = "olympus:oft7,harrows@brandflask.com/brandflask"

	session, err := mgo.Dial(mongohost)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	db := session.DB("brandflask")
	// img := db.C("imgref")
	// count, _ := img.Find(bson.M{"Parentid": "155715798525"}).Count()
	// fmt.Println(count)

	var result []interface{}

	db.C("brandrefx").Find(nil).All(&result)
	db.C("brandimages").DropCollection()
	bimages := db.C("brandimages")

	// retrieveImages("127539803930077")
	for _, v := range result {
		fbid := v.(bson.M)["fbid"]
		images := retrieveImages(fbid.(string))
		fmt.Println(fbid)
		var item = make(map[string]interface{})
		item["fbid"] = fbid.(string)
		item["images"] = images
		bimages.Insert(bson.M(item))
		time.Sleep(1 * time.Second)
	}

}
