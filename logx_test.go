package logx_test

import (
	"strconv"
	"testing"
)

func BenchmarkAllocs(b *testing.B) {
	b.Run(`strings`, func(b *testing.B) {
		var a string
		fn := func(param string) {
			a = param
		}
		for i:=1; i<b.N; i++ {
			func(){
				p := []byte(strconv.Itoa(i))
				fn(string(p))
			}()
		}
	})
	b.Run(`bytes`, func(b *testing.B) {
		var a string
		fn := func(param []byte) {
			a = string(param)
		}
		for i:=1; i<b.N; i++ {
			func(){
				p := []byte(strconv.Itoa(i))
				fn(p)
			}()
		}
	})
}
