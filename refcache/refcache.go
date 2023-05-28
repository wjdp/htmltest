// Package refcache : caches results of external link checking.
package refcache

import (
	"encoding/json"
	"github.com/theunrepentantgeek/htmltest/output"
	"os"
	"path"
	"sync"
	"time"
)

// RefCache struct : store of cached references.
type RefCache struct {
	refStore     map[string]CachedRef
	rwMutex      *sync.RWMutex
	cacheExpires time.Duration
}

// NewRefCache : Create a cached reference.
func NewRefCache(storePath string, cacheExpiresStr string) *RefCache {
	rS := &RefCache{}
	_ = storePath
	rS.rwMutex = &sync.RWMutex{}
	rS.cacheExpires, _ = time.ParseDuration(cacheExpiresStr)

	if !rS.ReadStore(storePath) {
		rS.refStore = make(map[string]CachedRef)
	}

	return rS
}

// ReadStore : Read a saved store from storePath.
func (rS *RefCache) ReadStore(storePath string) bool {
	// Read in RefCache
	f, err := os.Open(storePath)
	if err != nil {
		if !os.IsNotExist(err) {
			output.CheckErrorPanic(err)
		}
		return false
	}
	defer f.Close()

	var refStore map[string]CachedRef
	err = json.NewDecoder(f).Decode(&refStore)
	output.CheckErrorPanic(err)

	rS.refStore = refStore
	return true
}

// WriteStore : Write store to storePath.
func (rS *RefCache) WriteStore(storePath string) {
	// Write out RefCache
	os.MkdirAll(path.Dir(storePath), 0777)
	f, err := os.Create(storePath)
	output.CheckErrorPanic(err)
	defer f.Close()

	err = json.NewEncoder(f).Encode(&rS.refStore)
	output.CheckErrorPanic(err)
}

// CachedRef struct : Single cached result
type CachedRef struct {
	StatusCode int
	LastSeen   time.Time
	// Body byte[] // For when we do hash checking on external documents
}

// Get a cached result, thread safe.
func (rS *RefCache) Get(urlStr string) (*CachedRef, bool) {
	rS.rwMutex.RLock()
	val, ok := rS.refStore[urlStr]
	rS.rwMutex.RUnlock()
	if ok {
		// In cache, check if cache has expired
		if time.Now().Before(val.LastSeen.Add(rS.cacheExpires)) {
			// All ok!
			return &val, true
		}
		// Nope, cache has expired
		return nil, false
	}
	return nil, false
}

// Save a result to the cache, thread safe.
func (rS *RefCache) Save(urlStr string, statusCode int) {
	cR := CachedRef{
		StatusCode: statusCode,
		LastSeen:   time.Now(),
	}
	rS.rwMutex.Lock()
	rS.refStore[urlStr] = cR
	rS.rwMutex.Unlock()
}
