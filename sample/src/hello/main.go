package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"poetry"
	_ "sort"
	_ "strconv"
	"sync"
)

var cacheMutex sync.Mutex
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

	var wg sync.WaitGroup
	cache = make(map[string]poetry.Poem)

	// Pre-load (in parallel!) all valid poems into a cache so we don't have to
	// load a poem each time it is requested.
	for _, name := range c.ValidPoems {
		wg.Add(1)
		go func(n string) {
			cacheMutex.Lock()                  // to protect a map shared between goroutines
			cache[n], err = poetry.LoadPoem(n) // has to be n, not name!
			cacheMutex.Unlock()
			if err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}(name)
	}

	wg.Wait() // wait for the cache to be ready before starting the server

	// Setup and start web server.
	http.HandleFunc(c.Route, poemHandler)
	http.ListenAndServe(c.BindAddress, nil)
}
