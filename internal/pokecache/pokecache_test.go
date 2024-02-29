package pokecache

import (
	"fmt"
	"sync"
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
				t.Errorf("expected to find key")
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
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}

func TestConcurrency(t *testing.T) {
	cache := NewCache(5 * time.Second)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			val := []byte(fmt.Sprintf("value%d", i))
			cache.Add(key, val)
			retrievedVal, ok := cache.Get(key)
			if !ok || string(retrievedVal) != string(val) {
				t.Errorf("Concurrency test failed for key %s", key)
			}
		}(i)
	}
	wg.Wait()
}

func TestReapLoopAdvanced(t *testing.T) {
	cache := NewCache(1 * time.Millisecond)

	cache.Add("expired", []byte("expired"))
	time.Sleep(10 * time.Millisecond)

	cache.Add("notExpired", []byte("notExpired"))

	_, ok := cache.Get("expired")
	if ok {
		t.Errorf("Expected expired item to be removed")
	}

	_, ok = cache.Get("notExpired")
	if !ok {
		t.Errorf("Expected notExpired item to still be in the cache")
	}

	cache.StopReaping()
}

func TestEdgeCases(t *testing.T) {
	cache := NewCache(5 * time.Second)
	cache.Add("key", []byte(""))

	_, ok := cache.Get("nonExistentKey")
	if ok {
		t.Errorf("Expected to not find item with non-existent key")
	}
}
