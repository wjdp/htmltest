package htmltest

import (
	"path"
	"testing"

	"github.com/theunrepentantgeek/htmltest/issues"
)

var genericTests = []struct {
	fixture    string
	errorCount int
}{
	{"areaValid.html", 0},
	{"areaBroken.html", 1},
	{"areaBlank.html", 2},
	{"areaMissing.html", 0},
	{"audioValid.html", 0},
	{"audioBroken.html", 5},
	{"audioBlank.html", 5},
	{"audioMissing.html", 0},
	{"citeValid.html", 0},
	{"citeBroken.html", 4},
	{"citeBlank.html", 4},
	{"citeMissing.html", 0},
	{"embedValid.html", 0},
	{"embedBroken.html", 1},
	{"embedBlank.html", 2},
	{"embedMissing.html", 0},
	{"iframeValid.html", 0},
	{"iframeBroken.html", 1},
	{"iframeBrokenButIgnored.html", 0},
	{"iframeBlank.html", 2},
	{"iframeMissing.html", 0},
	{"inputSrcValid.html", 0},
	{"inputSrcBroken.html", 1},
	{"inputSrcBlank.html", 2},
	{"inputSrcMissing.html", 0},
	{"objectValid.html", 0},
	{"objectBroken.html", 2},
	{"objectBlank.html", 2},
	{"objectMissing.html", 0},
	{"videoValid.html", 0},
	{"videoBroken.html", 9},
	{"videoBlank.html", 9},
	{"videoMissing.html", 0},
}

func TestCheckGenericTable(t *testing.T) {
	for _, gt := range genericTests {
		hT := tTestFileOpts(path.Join("fixtures/generic", gt.fixture),
			map[string]interface{}{"VCREnable": true})
		c := hT.issueStore.Count(issues.LevelError)
		if c != gt.errorCount {
			t.Error("error count", c, "!=", gt.errorCount, "in", gt.fixture)
			hT.issueStore.DumpIssues(true)
		}
	}
}
