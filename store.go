package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"sync"
)

const saveQueueLength = 1000

type URLStore struct {
	urls map[string]string
	mu   sync.RWMutex
	save chan record
}

type record struct {
	Key, URL string
}

func NewURLStore(filename string) *URLStore {
	s := &URLStore{
		urls: make(map[string]string),
		save: make(chan record, saveQueueLength),
	}

	if err := s.load(filename); err != nil {
		fmt.Errorf("error loading %s: %v", filename, err)
	}

	go s.saveLoop(filename)
	return s

}

func (s *URLStore) Get(key string) string {
	s.mu.RLock() // multiple go routines can read at a time but only one go routine can write at a time
	defer s.mu.RUnlock()
	return s.urls[key]
}

func (s *URLStore) Set(key, url string) bool {
	s.mu.Lock() // only one go routine can write at a time
	defer s.mu.Unlock()

	if _, present := s.urls[key]; present {
		return false
	}
	s.urls[key] = url
	return true
}

func (s *URLStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.urls)
}

func (s *URLStore) Put(url string) string {
	for {
		key := genKey(s.Count()) // generate short url
		if s.Set(key, url) {
			s.save <- record{key, url} // save to go routine
			return key
		}
	}
	panic("unreachable")
}

func (s *URLStore) load(filename string) error {

	f, err := os.Open(filename)
	if err != nil {
		PrintHandleError(err)
		return err
	}

	defer f.Close()

	d := gob.NewDecoder(f)

	for err == nil {
		var r record
		if err = d.Decode(&r); err == nil {
			s.Set(r.Key, r.URL)
		}
	}
	if err == io.EOF {
		return nil
	}

	PrintHandleError(err)
	return err
}

func (s *URLStore) saveLoop(filename string) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	FatalHandleError(err)

	defer f.Close()
	e := gob.NewEncoder(f) // return a encoder to io.Writer
	for {
		r := <-s.save // taking a record from the channel and encoding it to the file
		PrintHandleError(e.Encode(r))
	}

}
