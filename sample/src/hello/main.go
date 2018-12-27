package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"poetry"
	_ "sort"
	_ "strconv"
)

type config struct {
	Route       string
	BindAddress string   `json:"addr"`
	ValidPoems  []string `json:"valid"`
}

type poemWithTitle struct {
	Title string
	Body  poetry.Poem
	//NumWords    string
	NumWords    int
	NumTheLines int
}

// global variable; so I can access it in poemHandler()
var c config

func poemHandler(w http.ResponseWriter, r *http.Request) {
	// Get the poem name (curl localhost:8080/poem?name=wordsworth).
	r.ParseForm()
	poemName := r.Form["name"][0]

	// Check the poem is valid as defined by the config file.
	found := false
	for _, v := range c.ValidPoems {
		if v == poemName {
			found = true
			break
		}
	}
	if !found {
		http.Error(w, "File not found (invalid)", http.StatusNotFound)
		return
	}

	// Load the poem from the disk.
	p, err := poetry.LoadPoem(poemName)
	if err != nil {
		//log.Fatal(err)
		http.Error(w, "File not found", http.StatusInternalServerError)
	}

	// sort first stanza by line length
	//sort.Sort(p[0])

	//pwt := poemWithTitle{Title: poemName, Body: p, NumWords: p.NumWords()}
	//pwt := poemWithTitle{poemName, p, strconv.Itoa(p.NumWords()), p.NumThe()}
	pwt := poemWithTitle{poemName, p, p.NumWords(), p.NumThe()}

	// Encode poem as JSON.
	enc := json.NewEncoder(w)
	err = enc.Encode(pwt)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Open config file.
	f, err := os.Open("config")
	if err != nil {
		log.Fatal(err)
	}

	// Read in the configuration from the config file.
	dec := json.NewDecoder(f)
	err = dec.Decode(&c)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Setup and start web server.
	http.HandleFunc(c.Route, poemHandler)
	http.ListenAndServe(c.BindAddress, nil)
}
