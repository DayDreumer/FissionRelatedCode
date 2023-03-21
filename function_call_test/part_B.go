package main

import (
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	v := r.Form.Get("partB_name")
	if v == "" {
		v = "none"
	}
	w.Write([]byte(v + " did job B."))
}