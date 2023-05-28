package htmltest

import (
	"github.com/theunrepentantgeek/htmltest/issues"
	"testing"
)

func BenchmarkExternal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tTestDirectoryOpts("/home/will/local/history-project/_site/",
			map[string]interface{}{"LogLevel": issues.LevelInfo, "CheckExternal": false})
	}
}
