package main

import (
	"encoding/json"
	"log"
	"net/http"
	"poetry"
)

type poemWithTitle struct {
	Title       string
	Body        poetry.Poem
	NumWords    int
	NumTheLines int
}

func poemHandler(w http.ResponseWriter, r *http.Request) {
	// curl localhost:8080/poem?name=wordsworth | jq
	r.ParseForm()
	poemName := r.Form["name"][0]

	p, err := poetry.LoadPoem(poemName)
	if err != nil {
		log.Fatal(err)
	}

	//pwt := poemWithTitle{Title: poemName, Body: p, NumWords: p.NumWords()}
	pwt := poemWithTitle{poemName, p, p.NumWords(), p.NumThe()}

	enc := json.NewEncoder(w)
	err = enc.Encode(pwt)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	http.HandleFunc("/poem", poemHandler)
	http.ListenAndServe(":8080", nil)
}
