package refcache

// cSpell:ignore URLSTR STOREPATH

import (
	"os"
	"strings"
	"testing"

	"github.com/daviddengcn/go-assert"
)

// TestRefCacheQueryParameterEscaping: Tests that URLs with query parameters
// containing special characters like '&' and '+' are NOT escaped when written
// to the JSON cache file. See: https://github.com/wjdp/htmltest/issues/239
func TestRefCacheQueryParameterEscaping(t *testing.T) {
	rS := NewRefCache("does-not-exist", "2s")
	URLSTR := "https://grpc.io?param1=a+b&param2=Hello,+World"
	rS.Save(URLSTR, 206)

	// Write the cache to disk
	STOREPATH := ".htmltest/refcache-test-query-escape.json"
	rS.WriteStore(STOREPATH)
	defer func() { _ = os.Remove(STOREPATH) }()

	// Read the raw JSON file content to inspect how the URL was encoded
	rawJSON, err := os.ReadFile(STOREPATH)
	if err != nil {
		t.Fatalf("failed to read cache file: %v", err)
	}
	jsonContent := string(rawJSON)

	// The '&' character should appear literally in the JSON, not HTML-escaped
	if !strings.Contains(jsonContent, "param1=a+b&param2=Hello,+World") {
		t.Errorf("Expected literal '&' in JSON output for query parameters")
		t.Logf("JSON content: %s", jsonContent)
	}

	// The URL should NOT be HTML-escaped to '\u0026'
	if strings.Contains(jsonContent, "\\u0026") {
		t.Errorf("Found HTML-escaped '\\u0026' in JSON, but '&' should not be escaped")
	}

	// Verify that the cache can still functionally read the URL back
	rS2 := NewRefCache(STOREPATH, "2s")
	cR, ok := rS2.Get(URLSTR)
	assert.IsTrue(t, "url should be retrievable from cache", ok)
	assert.Equals(t, "status code should match", cR.StatusCode, 206)
}
