package lru

import (
	"reflect"
	"testing"
)

type String string

func (s String) Len() int64 {
	return int64(len(s))
}

func TestGet(t *testing.T) {
	lru := New(0, nil)
	lru.Add("key1", String("12345"))
	if v, has := lru.Get("key1"); !has || string(v.(String)) != "12345" {
		t.Fatalf("cache hit key1=12345 failed")
	}
	if _, has := lru.Get("key2"); has {
		t.Fatalf("cache miss key2 failed")
	}
}

func TestRemoveOldset(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	cap := int64(len(k1 + k2 + v1 + v2))
	lru := New(cap, nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))
	if _, has := lru.Get(k1); has || lru.ll.Len() != 2 {
		t.Fatalf("RemoveOldest key1 failed")
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	onEvictedFunc := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := New(10, onEvictedFunc)
	lru.Add("key1", String("123456"))
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))
	lru.Add("k4", String("k4"))

	expect := []string{"key1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}
