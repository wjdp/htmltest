package refcache

import (
  "time"
  "encoding/json"
)

type cachedPath struct {
  LastSeen time.Time
  StatusCode int
}

var pathStore map[string]cachedPath

func init() {
  pathStore = make(map[string]cachedPath)
}

func CachedURLStatus(urlStr string) int {
  return pathStore[urlStr].StatusCode
}

func SetCachedURLStatus(urlStr string, status int) {
  pathStore[urlStr] = cachedPath{
    LastSeen: time.Now(),
    StatusCode: status,
  }
  b, _ := json.Marshal(pathStore)
  _ = b
  // log.Println( string(b) )
}
