package pokecache

import (
	"reflect"
	"testing"
	"time"
)

func TestCacheCreation(t *testing.T) {
	interval := 5 * time.Second
	cache := NewCache(interval)
	if reflect.TypeOf(cache).String() != "*pokecache.Cache" {
		t.Errorf("expected type *pokecache.Cache, got %T", cache)
	}
}

func TestCacheAdd(t *testing.T) {
	interval := 5 * time.Second
	cache := NewCache(interval)
	cache.Add("test", []byte{})
	entry, exists := cache.cacheEntries["test"]
	if !exists {
		t.Errorf("test entry did not get added to test cache")
	}
	if reflect.TypeOf(entry.createdAt).String() != "time.Time" {
		t.Errorf("expected type time.Time, or cache.cacheEntries['test'].createdAt got %T", entry)
	}

}

func TestCacheGetTrue(t *testing.T) {
	interval := 5 * time.Second
	cache := NewCache(interval)
	cache.Add("test", []byte{})
	entry, result := cache.Get("test")
	if !result && entry != nil {
		t.Errorf("cache.Get() did not return true for an existing entry")
	}
	if result && entry == nil {
		t.Errorf("cache.Get() returned an empty result and marked it true")
	}
}

func TestCacheGetFalse(t *testing.T) {
	interval := 5 * time.Second
	cache := NewCache(interval)
	cache.Add("test", []byte{})
	entry, result := cache.Get("foo")
	if !result && entry != nil {
		t.Errorf("cache.Get() said nothing returned, but there was a result.")
	}
	if result && entry == nil {
		t.Errorf("cache.Get() says something returned, but there was no result.")
	}
}

func TestReapLoop(t *testing.T) {
	interval := 5 * time.Millisecond
	cache := NewCache(interval)

	cache.Add("test", []byte("value"))
	_, exists := cache.Get("test")
	if !exists {
		t.Fatal("entry should exist before reaping")
	}

	time.Sleep(interval * 3)

	_, exists = cache.Get("test")
	if exists {
		t.Error("entry should have been reaped")
	}
}
