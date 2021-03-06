package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"regexp"
)

/* Config
database configuration
*/
type Config struct {
	PublishingPlus map[string]interface{} `yaml:"publishing-plus"`
	Directory      map[string]interface{} `yaml:"directory"`
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
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("select int_id, pretext from account_page_groups where account_name=?", account)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var id int
	var pretext string

	r, _ := regexp.Compile("^http:")

	for rows.Next() {
		err = rows.Scan(&id, &pretext)
		if err != nil {
			panic(err)
		}
		if r.MatchString(pretext) {
			updated := r.ReplaceAllString(pretext, "https:")
			sql := fmt.Sprintf("UPDATE account_page_groups SET pretext='%s' WHERE int_id=%d;", updated, id)
			fmt.Println(sql)
			if isexec {
				_, err = db.Exec(sql)
				if err != nil {
					panic(err)
				}
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
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("select c.campaign_id,c.published_url from captora_object  o, campaign c where o.account_name = ? and o.type = 'CAMPAIGN' and c.object_key = o.object_key ", account)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var id int
	var url string

	r, _ := regexp.Compile("^http:")

	for rows.Next() {
		err = rows.Scan(&id, &url)
		if err != nil {
			panic(err)
		}
		if r.MatchString(url) {
			updated := r.ReplaceAllString(url, "https:")
			sql := fmt.Sprintf("UPDATE campaign SET published_url='%s' WHERE campaign_id=%d;", updated, id)
			fmt.Println(sql)
			if isexec {
				_, err = db.Exec(sql)
				if err != nil {
					panic(err)
				}
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
	if err != nil {
		panic(err)
	}
	var database Config

	err = yaml.Unmarshal(data, &database)
	if err != nil {
		panic(err)
	}

	if database.Directory[env] == nil {
		fmt.Printf("%s not exists", env)
		os.Exit(1)
	}

	modifyPageGroup(account, database.Directory[env], confirm)
	modifyPublishedURL(account, database.PublishingPlus[env], confirm)
}
