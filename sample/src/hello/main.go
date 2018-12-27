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

var cache map[string]poetry.Poem

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

	// Get the poem from the cache.
	p, ok := cache[poemName]
	if !ok {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// sort first stanza by line length
	//sort.Sort(p[0])

	//pwt := poemWithTitle{Title: poemName, Body: p, NumWords: p.NumWords()}
	//pwt := poemWithTitle{poemName, p, strconv.Itoa(p.NumWords()), p.NumThe()}
	pwt := poemWithTitle{poemName, p, p.NumWords(), p.NumThe()}

	// Encode poem as JSON.
	enc := json.NewEncoder(w)
	err := enc.Encode(pwt)
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

	// Pre-load all valid poems into a cache so we don't have to load a poem
	// each time it is requested.
	cache = make(map[string]poetry.Poem)
	for _, name := range c.ValidPoems {
		cache[name], err = poetry.LoadPoem(name)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Setup and start web server.
	http.HandleFunc(c.Route, poemHandler)
	http.ListenAndServe(c.BindAddress, nil)
}
