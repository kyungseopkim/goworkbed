package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

/* Config
database configuration
*/
type Config struct {
	PublishingPlus map[string]interface{} `yaml:"publishing-plus"`
	Directory      map[string]interface{} `yaml:"directory"`
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func getDB(env map[interface{}]interface{}) (*sql.DB, error) {
	host := env["host"]
	user := env["user"]
	passwd := env["password"]
	dbname := env["database"]
	connstr := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", user, passwd, host, dbname)
	return sql.Open("mysql", connstr)
}

func modifyPageGroup(account string, env interface{}, confirm string) {
	isexec := false
	if len(confirm) > 0 {
		isexec = true
	}

	db, err := getDB(env.(map[interface{}]interface{}))
	checkError(err)
	defer db.Close()

	rows, err := db.Query("select int_id, pretext from account_page_groups where account_name=?", account)
	checkError(err)
	defer rows.Close()

	var id int
	var pretext string

	r, _ := regexp.Compile("^http:")

	for rows.Next() {
		err = rows.Scan(&id, &pretext)
		checkError(err)
		if r.MatchString(pretext) {
			updated := r.ReplaceAllString(pretext, "https:")
			sql := fmt.Sprintf("UPDATE account_page_groups SET pretext='%s' WHERE int_id=%d;", updated, id)
			fmt.Println(sql)
			if isexec {
				_, err = db.Query(sql)
				checkError(err)
				fmt.Println("Executed")
			}
		}
	}
}

func modifyPublishedURL(account string, env interface{}, confirm string) {
	isexec := false
	if len(confirm) > 0 {
		isexec = true
	}

	db, err := getDB(env.(map[interface{}]interface{}))
	checkError(err)
	defer db.Close()

	rows, err := db.Query("select c.campaign_id,c.published_url from captora_object  o, campaign c where o.account_name = ? and o.type = 'CAMPAIGN' and c.object_key = o.object_key ", account)
	checkError(err)
	defer rows.Close()

	var id int
	var url string

	r, _ := regexp.Compile("^http:")

	for rows.Next() {
		err = rows.Scan(&id, &url)
		checkError(err)
		if r.MatchString(url) {
			updated := r.ReplaceAllString(url, "https:")
			sql := fmt.Sprintf("UPDATE campaign SET published_url='%s' WHERE campaign_id=%d;", updated, id)
			fmt.Println(sql)
			if isexec {
				_, err = db.Query(sql)
				checkError(err)
				fmt.Println("Executed")
			}
		}
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("args : Account Env [confirm]")
		os.Exit(1)
	}

	account := os.Args[1]
	env := os.Args[2]

	var confirm string
	if len(os.Args) > 3 {
		confirm = os.Args[3]
	}

	databaseYaml := "./database.yml"

	if _, err := os.Stat(databaseYaml); err != nil {
		fmt.Printf("%s should be here!!!\n", databaseYaml)
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(databaseYaml)
	checkError(err)
	var database Config

	err = yaml.Unmarshal(data, &database)
	checkError(err)

	if database.Directory[env] == nil {
		fmt.Printf("%s not exists", env)
		os.Exit(1)
	}

	modifyPageGroup(account, database.Directory[env], confirm)
	modifyPublishedURL(account, database.PublishingPlus[env], confirm)
}
