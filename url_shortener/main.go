package main

import (
	"flag"
	"fmt"
	"golang/practice/handler"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gomodule/redigo/redis"
)

func initializeHandlers(mux *http.ServeMux, yamlFile string, jsonFile string, c redis.Conn) http.HandlerFunc {

	h := handler.RedisHanlder(c, mux)
	if yamlFile != "" {
		yaml, err := ioutil.ReadFile(yamlFile)
		if err != nil {
			fmt.Println(err)
		}
		h, _ = handler.YAMLHandler(yaml, h)
	}

	if jsonFile != "" {
		json, err := ioutil.ReadFile(jsonFile)
		if err != nil {
			fmt.Println(err)
		}
		h, _ = handler.JsonHandler(json, h)
	}

	return h
}

func main() {

	yamlFile := flag.String("yamlFile", "", "YAML file name for URLs")
	jsonFile := flag.String("jsonFile", "", "JSON file name for URLs")

	flag.Parse()
	mux := defaultMux()

	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	_, err = conn.Do("HMSET", "urls", "/", "localhost:8080", "/urlshort3", "https://github.com/gophercises/urlshort", "/urlshort4", "https://github.com/gophercises/urlshort")
	if err != nil {
		log.Fatal(err)
	}
	finalHandler := initializeHandlers(mux, *yamlFile, *jsonFile, conn)
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", finalHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
