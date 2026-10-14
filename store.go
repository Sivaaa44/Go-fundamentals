package main

import (
	"sync"
	"time"
)

type entry struct {
	value   string
	expires time.Time
}

func (e entry) expired() bool {
	return !e.expires.IsZero() && time.Now().After(e.expires)
}

type SharedStore struct {
	mx   sync.Mutex
	data map[string]entry
}

func NewStore() *SharedStore {
	return &SharedStore{data: make(map[string]entry)}
}

// Set stores a value. ttl <= 0 means no expiry.
func (s *SharedStore) Set(key, value string, ttl time.Duration) {
	s.mx.Lock()
	defer s.mx.Unlock()
	e := entry{value: value}
	if ttl > 0 {
		e.expires = time.Now().Add(ttl)
	}
	s.data[key] = e
}

// Get is "lazy expiry": an expired key is deleted when someone asks for it.
func (s *SharedStore) Get(key string) (string, bool) {
	s.mx.Lock()
	defer s.mx.Unlock()
	e, ok := s.data[key]
	if !ok {
		return "", false
	}
	if e.expired() {
		delete(s.data, key)
		return "", false
	}
	return e.value, true
}

func (s *SharedStore) Del(key string) bool {
	s.mx.Lock()
	defer s.mx.Unlock()
	e, ok := s.data[key]
	delete(s.data, key)
	return ok && !e.expired()
}

// TTL returns seconds left, -1 if no expiry, -2 if the key doesn't exist.
func (s *SharedStore) TTL(key string) int {
	s.mx.Lock()
	defer s.mx.Unlock()
	e, ok := s.data[key]
	if !ok || e.expired() {
		return -2
	}
	if e.expires.IsZero() {
		return -1
	}
	return int(time.Until(e.expires).Seconds())
}

// removes expired keys
func (s *SharedStore) StartSweeper(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			s.mx.Lock()
			for k, e := range s.data {
				if e.expired() {
					delete(s.data, k)
				}
			}
			s.mx.Unlock()
		}
	}()
}
