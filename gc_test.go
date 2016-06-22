package logx

import (
	"runtime"
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
)

type gcMock struct {
	any int
}

func runWithGC(res *bool) {
	mock := &gcMock{}
	runtime.SetFinalizer(mock, func(m *gcMock) {
		*res = true
	})
}

func TestGC_Finalizer(t *testing.T) {
	var res bool
	runWithGC(&res)
	runtime.GC()
	time.Sleep(10 * time.Microsecond)
	assert.True(t, res)
}
