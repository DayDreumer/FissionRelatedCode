package main

import (
	"log"
	"net/http"
	"io/ioutil"
	"os"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	var result string
	
	
	req1, err := http.NewRequest("GET", "http://router.fission/partA", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	q1 := req1.URL.Query()
	q1.Add("partA_name", "Bill")
	req1.URL.RawQuery = q1.Encode()
	
	var resp1 *http.Response
    resp1, err = http.DefaultClient.Do(req1)
	if err != nil {
        log.Print(err)
    }
    defer resp1.Body.Close()
	
	var resp2 *http.Response
	req2, err := http.NewRequest("GET", "http://router.fission/partB", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	q2 := req2.URL.Query()
	q2.Add("partB_name", "Mike")
	req2.URL.RawQuery = q2.Encode()
	
    resp2, err = http.DefaultClient.Do(req2)
	if err != nil {
        log.Print(err)
    }
    defer resp2.Body.Close()
	
	
	body1, err := ioutil.ReadAll(resp1.Body)
	if err != nil {
		// handle error
		log.Print(err)
	}
	body2, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		// handle error
		log.Print(err)
	}
	
	result = string(body1) + string(body2)
	w.Write([]byte(result))
}