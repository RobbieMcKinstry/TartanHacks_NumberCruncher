package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"log"

	"github.com/gorilla/mux"

	"./smoother"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)

	http.Handle("/", r)
	fmt.Println("Hello world")
	fmt.Println("Listening and serving on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Data struct {
	BloombergData []int `json:"data"`
}

func indexHandler(w http.ResponseWriter, req *http.Request) {

	var data *Data
	req.ParseForm()
	for key, _ := range req.Form {
		log.Println(key)         //LOG: {"test": "that"}
		data = &Data{}
		err := json.Unmarshal([]byte(key), data)
		if err != nil {
			log.Println(err.Error())
		}
	}

	result := smoother.Smooth(data.BloombergData)
	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
	}
	w.Write(jsonResult)

}
