package main

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"sync"

	check "github.com/jaeyeom/sugo/errors/must"
	mgo "gopkg.in/mgo.v2"
	resty "gopkg.in/resty.v0"
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
	"superLotto": "superlotto-plus",
	"megaLotto":  "mega-millions",
	"powerLotto": "powerball"}

type formatRange struct {
	Start int /* inclusive */
	End   int /* exclusive */
}

const stringFormat = "^(\\d+)\\s+(\\w\\w\\w\\.\\s\\w\\w\\w\\s\\d\\d,\\s\\d\\d\\d\\d)\\s+(\\d+)\\s+(\\d+)\\s+(\\d+)\\s+(\\d+)\\s+(\\d+)\\s+(\\d+)\\s*"

var lineFormat *regexp.Regexp = regexp.MustCompile(stringFormat)

const trimSet = " \r\n"

func newDraw(line string) *Draw {
	if len(line) > 0 {
		result := new(Draw)
		items := lineFormat.FindStringSubmatch(line)
		if len(items) != 9 {
			return nil
		}
		val := check.Int(strconv.Atoi(items[1]))
		result.Id = val
		result.Date = items[2]

		for i := 0; i < 5; i++ {
			val = check.Int(strconv.Atoi(items[3+i]))
			result.LuckyNo = append(result.LuckyNo, val)
		}

		val = check.Int(strconv.Atoi(items[8]))
		result.Mega = val

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

func fetchAndUpdate(wg *sync.WaitGroup, kind string) {
	defer wg.Done()
	fmt.Printf("%s is fetching\n", kind)

	resp := check.Any(resty.R().Get(getEndPoint(kind))).(*resty.Response)
	data := string(resp.Body())

	// rawdata, err := ioutil.ReadFile("output.txt")
	// if err != nil {
	// 	panic(err)
	// }
	// data := string(rawdata)

	draws := make([]*Draw, 0)

	for _, line := range strings.Split(data, "\n")[5:] {
		item := newDraw(line)
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

	lotteriaDB := session.DB("lotteria")

	names, err := lotteriaDB.CollectionNames()
	if err != nil {
		panic(err)
	}

	sort.Strings(names)
	i := sort.SearchStrings(names, kind)

	collection := session.DB("lotteria").C(kind)
	count := 0

	if i < len(names) && strings.Compare(names[i], kind) == 0 {
		var latest Draw
		err = collection.Find(nil).Sort("-id").One(&latest)
		if err != nil {
			panic(err)
		}

		for _, item := range draws {
			if item.Id > latest.Id {
				collection.Insert(item)
				count++
			} else {
				break
			}
		}
	} else {
		collection.Create(&mgo.CollectionInfo{})
		for _, item := range draws {
			collection.Insert(item)
			count++
		}
	}
	fmt.Printf("%s : %v items are inserted\n", kind, count)
}

func main() {
	// gocron.Every(1).Thursday().Do(fetchAndUpdate)
	// gocron.Every(1).Sunday().Do(fetchAndUpdate)
	// <- gocron.Start()

	var wg sync.WaitGroup

	for kind := range kinds {
		wg.Add(1)
		go fetchAndUpdate(&wg, kind)
	}

	wg.Wait()

}
