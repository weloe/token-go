package token_go

import "testing"

func TestFirst(t *testing.T) {
	First()
}

func BenchmarkFirst(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		First()
	}
}
