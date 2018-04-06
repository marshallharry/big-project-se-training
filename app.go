package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/tokopedia/big-project/helper"
)

var (
	getRedis  = helper.GetRedis
	setRedis  = helper.SetRedis
	getUsers  = helper.GetUsers
	closeDB   = helper.CloseConnection
	publish   = helper.Publish
	subscribe = helper.Subscribe
)

func main() {

	http.HandleFunc("/html", handleHTML)
	http.HandleFunc("/search", handleSearch)
	http.ListenAndServe(":8080", nil)
	closeDB()
}

func handleHTML(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		log.Println(err)
		w.Write([]byte(fmt.Sprint(err)))
	}

	users := getUsers("")
	visitor := int64(1)

	val, err := getRedis("mh0318-redis-visitor")
	if err != nil {
		log.Print(err)
	}

	if val == "" {
		publish("mh0318-nsq-visitor-set", []byte("test"))
	} else {
		visitor, err = strconv.ParseInt(val, 0, 64)
	}
	publish("mh0318-nsq-visitor-incr", []byte("test"))

	data := map[string]interface{}{
		"Users":   users,
		"Visitor": visitor,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err)
		w.Write([]byte(fmt.Sprint(err)))
	}
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	key := queryValues.Get("key")

	w.Header().Set("Content-Type", "application/json")

	users := getUsers(key)

	jsonbyte, err := json.Marshal(users)
	if err != nil {
		w.Write([]byte(fmt.Sprint(err)))
	}

	w.Write(jsonbyte)
}
