package main

import (
	"fmt"
	"log"
	"net/http"
)

const AddForm = `
<html><body>
<form method="POST" action="/add">
URL: <input type="text" name="url">
<input type="submit" value="Add">
</form>
<\html><\body>`

var store = NewURLStore()

func main() {
	http.HandleFunc("/", Redirect)
	http.HandleFunc("/add", Add)
	//http.ListenAndServe("8080", nil)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}

func Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	url := r.FormValue("url")
	if url == "" {
		fmt.Fprintf(w, AddForm)
		return
	}
	key := store.Put(url)

	fmt.Fprintf(w, "URL %s added with key %s", url, key)
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
