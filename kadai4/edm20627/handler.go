package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, HTTPサーバ")
}

func pserson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	p := Person{Name: "tenten", Age: 32}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(p); err != nil {
		http.Error(w, "hoge", 300)
	}
	fmt.Fprint(w, buf.String())
}

func form(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, r.FormValue("msg"))
}

func api(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, r.FormValue("msg"))
}
