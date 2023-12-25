package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/rpc"
)

const addForm = `
<html><body>
<form method="POST" action="/add">
URL: <input type="text" name="url">
<input type="submit" value="Add">
</form>
<\html><\body>
`

var store Store

var (
	listenAddr = flag.String("http", ":8080", "HTTP listen address")
	dataFile   = flag.String("file", "store.json", "json store filename")
	rpcEnabled = flag.Bool("rpc", false, "enable RPC server")
	hostname   = flag.String("hostname", "localhost", "hostname and port")
	masterAddr = flag.String("master", "", "RPC master address")
)

func main() {
	flag.Parse()
	store = NewURLStore(*dataFile)

	if *masterAddr != "" {
		store = NewProxyStore(*masterAddr)
	} else {
		store = NewURLStore(*dataFile)
	}

	if *rpcEnabled {
		rpc.RegisterName("Store", store)
		rpc.HandleHTTP()
	}

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
	var key string
	if err := store.Put(&url, &key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "Added url %s as %s", url, key)
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]
	if key == "" {
		http.NotFound(w, r)
		return
	}
	var url string
	if err := store.Get(&key, &url); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}
