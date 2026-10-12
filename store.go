package main

import (
	"sync"
)

type SharedStore struct {
	Store map[string]string
	mx    sync.Mutex
}

func (s *SharedStore) Set(key, value string) bool {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.Store[key] = value
	//if some error we nee to return false
	return true
}

func (s *SharedStore) Get(key string) (string, bool) {
	s.mx.Lock()
	defer s.mx.Unlock()
	if val, ok := s.Store[key]; ok {
		return val, ok
	} else {
		return "not found", ok
	}
}

func (s *SharedStore) Del(key string) bool {
	s.mx.Lock()
	defer s.mx.Unlock()
	if _, ok := s.Store[key]; ok {
		delete(s.Store, key)
		return ok
	} else {
		return ok
	}
}
