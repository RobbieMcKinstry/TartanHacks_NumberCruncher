package main

import (
	"fmt"
	"net/http"
	"log"

	"./smoother"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)

	http.Handle("/", r)
	fmt.Println("Hello world")
	log.Fatal(http.ListenAndServe())
}

type Data struct {
	BloombergData []int `json:"data"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	req.ParseForm()
	for key, _ := range req.Form {
		log.Println(key)         //LOG: {"test": "that"}
		data := &Data{}
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
	return jsonResult

}

