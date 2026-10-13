package main

import "testing"

func TestSetGet(t *testing.T) {
	s := SharedStore{Store: make(map[string]string)}

	s.Set("name", "siva")

	val, ok := s.Get("name")
	if !ok || val != "siva" {
		t.Errorf("expected siva, got %q (ok=%v)", val, ok)
	}
}

func TestGetMissing(t *testing.T) {
	s := SharedStore{Store: make(map[string]string)}

	_, ok := s.Get("nope")
	if ok {
		t.Errorf("expected missing key to return ok=false")
	}
}

func TestDel(t *testing.T) {
	s := SharedStore{Store: make(map[string]string)}
	s.Set("a", "1")

	if !s.Del("a") {
		t.Errorf("expected Del to return true for existing key")
	}
	if _, ok := s.Get("a"); ok {
		t.Errorf("key should be gone after Del")
	}
	if s.Del("a") {
		t.Errorf("expected Del to return false the second time")
	}
}
