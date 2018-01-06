package logx_test

import (
	"bytes"
	"strconv"
	"sync"
	"testing"
)

func BenchmarkParamAllocations(b *testing.B) {
	b.Run(`strings`, func(b *testing.B) {
		var a string
		fn := func(param string) {
			a = param
		}
		for i := 1; i < b.N; i++ {
			func() {
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
		for i := 1; i < b.N; i++ {
			func() {
				p := []byte(strconv.Itoa(i))
				fn(p)
			}()
		}
	})
}

func BenchmarkBufferPool(b *testing.B) {
	b.Run(`byte slice`, func(b *testing.B) {
		var res bytes.Buffer
		for i := 0; i < b.N; i++ {
			func() {
				var chunk []byte
				for j := 0; j < 10; j++ {
					chunk = append(chunk, []byte(strconv.Itoa(j))...)
				}
				res.Write(chunk)
			}()
		}
	})
	b.Run(`buffer pool`, func(b *testing.B) {
		var res bytes.Buffer
		pool := &sync.Pool{
			New: func() interface{} {
				return &bytes.Buffer{}
			},
		}
		for i := 0; i < b.N; i++ {
			func() {
				chunk := pool.Get().(*bytes.Buffer)
				for j := 0; j < 10; j++ {
					chunk.Write([]byte(strconv.Itoa(j)))
				}
				chunk.WriteTo(&res)
				chunk.Reset()
				pool.Put(chunk)
			}()
		}
	})
}
