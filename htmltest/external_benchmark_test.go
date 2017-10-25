package htmltest

import (
	"testing"
	"wjdp.uk/htmltest/issues"
)

func BenchmarkExternal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tTestDirectoryOpts("/home/will/local/history-project/_site/",
			map[string]interface{}{"LogLevel": issues.LevelInfo, "CheckExternal": false})
	}
}
