package main

import (
	"flag"
	"fmt"
	"net/http"
)

const addForm = `
<html><body>
<form method="POST" action="/add">
URL: <input type="text" name="url">
<input type="submit" value="Add">
</form>
<\html><\body>
`

var store *URLStore

var (
	listenAddr = flag.String("http", ":8080", "HTTP listen address")
	dataFile   = flag.String("file", "store.gob", "data store filename")
	hostname   = flag.String("hostname", "localhost", "hostname and port")
)

func main() {
	flag.Parse()
	store = NewURLStore(*dataFile)
	http.HandleFunc("/", Redirect)
	http.HandleFunc("/add", Add)
	err := http.ListenAndServe(*listenAddr, nil)
	FatalHandleError(err)
}

func Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	url := r.FormValue("url")
	if url == "" {
		fmt.Fprint(w, addForm)
		return
	}
	key := store.Put(url)

	fmt.Fprintf(w, "Added url %s as %s", url, key)
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]
	url := store.Get(key)
	if url == "" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}
