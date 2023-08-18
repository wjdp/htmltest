package refcache

import (
	"testing"
	"time"

	"github.com/daviddengcn/go-assert"
)

func TestRefCacheNew(t *testing.T) {
	rS := NewRefCache("does-not-exist", "2s")
	assert.NotEquals(t, "RefCache", rS, nil)
}

func TestRefSaveGet(t *testing.T) {
	// refcache should store and return cached items by key
	rS := NewRefCache("does-not-exist", "2s")
	URLSTR := "http://example.com/page.html"
	// check nil is returned when key is not in store
	_, ok1 := rS.Get(URLSTR)
	assert.IsFalse(t, "url not in store", ok1)
	// Save a status code in store
	rS.Save(URLSTR, 200)
	// Retrieve cached response
	cR, ok2 := rS.Get(URLSTR)
	assert.IsTrue(t, "url in store", ok2)
	assert.Equals(t, "get from RefCache (actual)", cR.StatusCode, 200)
}

func TestRefCacheWriteRead(t *testing.T) {
	// write cache, read back in again, preserves state
	rS1 := NewRefCache("does-not-exist", "2s")
	URLSTR := "http://example.com/page.html"
	rS1.Save(URLSTR, 200)
	STOREPATH := ".htmltest/refcache-test-writeread.json"
	rS1.WriteStore(STOREPATH)
	rS2 := NewRefCache(STOREPATH, "2s")
	_, okN := rS2.Get("xyz")
	cRY, okY := rS2.Get(URLSTR)
	assert.IsFalse(t, "url not in cache", okN)
	assert.IsTrue(t, "url in cache", okY)
	assert.Equals(t, "url status in cache", cRY.StatusCode, 200)
}

func TestRefCacheExpiry(t *testing.T) {
	// does the cache invalidate?
	rS := NewRefCache("does-not-exist", "1s")
	URLSTR := "http://example.com/page.html"
	rS.Save(URLSTR, 200)
	_, okY := rS.Get(URLSTR)
	assert.IsTrue(t, "cache valid", okY)
	d, _ := time.ParseDuration("1s")
	time.Sleep(d)
	_, okN := rS.Get(URLSTR)
	assert.IsFalse(t, "cache invalid", okN)
}
