package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

func indexHandler(w http.ResponseWriter, req *http.Request) {

	var data *smoother.History
	req.ParseForm()
	for key, _ := range req.Form {
		log.Println(key) //LOG: {"test": "that"}
		data = &smoother.History{}
		err := json.Unmarshal([]byte(key), data)
		if err != nil {
			log.Println(err.Error())
		}
	}
	fmt.Printf("%v", data)
	result := smoother.Smooth(data.Data)
	_ = result
	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
	}
	w.Write(jsonResult)

}
