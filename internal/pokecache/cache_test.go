package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key ")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	// Create cache with a very short interval for testing
	interval := 1 * time.Second
	t.Logf("Creating cache with interval: %v", interval)
	cache := NewCache(interval)

	// Add test item
	t.Log("Adding item to cache")
	cache.Add("test", []byte("test data"))

	// Verify item was added
	val, exists := cache.Get("test")
	if !exists {
		t.Fatal("Item should exist immediately after adding")
	}
	t.Logf("Item exists in cache with value: %s", val)

	// Wait longer than the interval
	waitTime := 2 * time.Second
	t.Logf("Waiting %v for reaper...", waitTime)
	time.Sleep(waitTime)

	// Check if item was reaped
	_, exists = cache.Get("test")
	if exists {
		t.Error("Expected item to be reaped from cache")
	} else {
		t.Log("Item was successfully reaped")
	}
}
