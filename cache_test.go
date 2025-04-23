package weblib

import (
	"fmt"
	"testing"
	"time"
)

func TestCache_PutAndGet(t *testing.T) {
	cache := NewCache(2*time.Second, 1*time.Second)
	defer cache.Close()

	cache.Put("key1", "value1")

	val := cache.Get("key1")
	if val != "value1" {
		t.Errorf("expected value1, got %v", val)
	}
}

func TestCache_Expiration(t *testing.T) {
	cache := NewCache(1*time.Second, 5*time.Second)
	defer cache.Close()

	cache.Put("key2", "value2")
	time.Sleep(1500 * time.Millisecond)

	val := cache.Get("key2")
	if val != nil {
		t.Errorf("expected nil, got %v", val)
	}
}

func TestCache_CleanerRemovesExpired(t *testing.T) {
	cache := NewCache(500*time.Millisecond, 300*time.Millisecond)
	defer cache.Close()

	cache.Put("key3", "value3")
	time.Sleep(1 * time.Second)

	cache.mu.Lock()
	_, exists := cache.data["key3"]
	cache.mu.Unlock()

	if exists {
		t.Errorf("expected key3 to be removed by cleaner")
	}
}

func TestCache_Delete(t *testing.T) {
	cache := NewCache(2*time.Second, 1*time.Second)
	defer cache.Close()

	cache.Put("key4", "value4")
	cache.Delete("key4")

	val := cache.Get("key4")
	if val != nil {
		t.Errorf("expected nil after delete, got %v", val)
	}
}

func TestCache_Close(t *testing.T) {
	cache := NewCache(1*time.Second, 100*time.Millisecond)
	cache.Close()

	// Calling Close again should not panic or block
	done := make(chan bool)
	go func() {
		cache.Close()
		done <- true
	}()

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Error("Close() blocked on second call")
	}
}

func TestCache_Concurrency(t *testing.T) {
	cache := NewCache(5*time.Second, 1*time.Second)
	defer cache.Close()

	const numGoroutines = 100
	const opsPerGoroutine = 100

	done := make(chan bool)

	// writers
	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < opsPerGoroutine; j++ {
				key := fmt.Sprintf("key%d", j)
				cache.Put(key, j)
			}
			done <- true
		}()
	}

	// readers
	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < opsPerGoroutine; j++ {
				value := cache.Get(fmt.Sprintf("key%d", j)).(int)
				if value != j {
					t.Errorf("failed concurrency, expected: %d, got: %v", j, value)
				}
			}
			done <- true
		}()
	}

	// wait
	for i := 0; i < numGoroutines*2; i++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("Timeout waiting for goroutines")
		}
	}
}
