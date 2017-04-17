package main

import (
	"fmt"
	"strings"

	"strconv"

	"gopkg.in/mgo.v2"
	"gopkg.in/resty.v0"
	"github.com/jasonlvhit/gocron"
)

type Draw struct {
	Id      int    `json:"id"`
	Date    string `json:"date"`
	LuckyNo []int  `json:"luckyNumber"`
	Mega    int    `json:"mega"`
}

func (d *Draw) String() string {
	return fmt.Sprintf(
		"[%v][%v][%v][%v][%v][%v][%v][%v]",
		d.Id, d.Date, d.LuckyNo[0], d.LuckyNo[1], d.LuckyNo[2], d.LuckyNo[3], d.LuckyNo[4], d.Mega)
}

var kinds = map[string]string{
	"super": "superlotto-plus",
	"mega":  "mega-millions",
	"power": "powerball"}

type formatRange struct {
	Start int /* inclusive */
	End   int /* exclusive */
}

var super_format = map[string]formatRange{
	"draw": formatRange{0, 6},
	"date": formatRange{8, 26},
	"no1":  formatRange{34, 42},
	"no2":  formatRange{46, 54},
	"no3":  formatRange{58, 66},
	"no4":  formatRange{69, 78},
	"no5":  formatRange{80, 90},
	"mega": formatRange{92, 100},
}

func checkError(err error) {

}
func string2Draw(draw *Draw, line string, k string, v formatRange) {

}

const TRIMSET = " \r\n"

func newDraw(line string) *Draw {
	if len(line) >= 100 {
		var result *Draw = new(Draw)
		for k, v := range super_format {
			switch {
			case strings.Compare(k, "draw") == 0:
				val, err := strconv.Atoi(strings.Trim(line[v.Start:v.End], TRIMSET))
				if err != nil {
					panic(err)
				}
				result.Id = val
			case strings.Compare(k, "date") == 0:
				result.Date = line[v.Start:v.End]
			case strings.Compare(k[0:2], "no") == 0:
				val, err := strconv.Atoi(strings.Trim(line[v.Start:v.End], TRIMSET))
				if err != nil {
					panic(err)
				}
				result.LuckyNo = append(result.LuckyNo, val)
			case strings.Compare(k, "mega") == 0:
				val, err := strconv.Atoi(strings.Trim(line[v.Start:v.End], TRIMSET))
				if err != nil {
					panic(err)
				}
				result.Mega = val
			}
		}
		return result
	}
	return nil
}

const endpoint = "http://www.calottery.com/sitecore/content/Miscellaneous/download-numbers/?GameName=%s&Order=No"

func getEndPoint(kind string) string {
	if val, ok := kinds[kind]; ok {
		return fmt.Sprintf(endpoint, val)
	}
	return ""
}

func fetchAndUpdate() {
	resp, err := resty.R().Get(getEndPoint("super"))
	if err != nil {
		panic(err)
	}
	data := string(resp.Body())

	/*
		rawdata, err := ioutil.ReadFile("output.txt")
		if err != nil {
			panic(err)
		}

		data := string(rawdata)
	*/

	var draws []*Draw = make([]*Draw, 0)

	for _, line := range strings.Split(data, "\n")[5:] {
		var item *Draw = newDraw(line)
		if item != nil {
			draws = append(draws, item)
		}
	}

	connStr := "mongodb://fetcher:fetcher@brandflask.com/lotteria"
	session, err := mgo.Dial(connStr)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	collection := session.DB("lotteria").C("superLottoDraw")

	var latest Draw
	err = collection.Find(nil).Sort("-id").One(&latest)
	if err != nil {
		panic(err)
	}

	count := 0
	for _, item := range draws {
		if item.Id > latest.Id {
			collection.Insert(item)
			count++
		}
	}

	fmt.Printf("%v items are inserted\n", count)
}

func main() {
	gocron.Every(1).Thursday().Do(fetchAndUpdate)
	gocron.Every(1).Sunday().Do(fetchAndUpdate)

	<- gocron.Start()
}
