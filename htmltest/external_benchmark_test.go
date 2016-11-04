package htmltest

import (
	"github.com/wjdp/htmltest/issues"
	"testing"
)

func BenchmarkExternal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t_testDirectoryOpts("/home/will/local/history-project/_site/",
			map[string]interface{}{"LogLevel": issues.NONE})
	}
}
