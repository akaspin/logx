// +build debug trace

package logx_test

import (
	"bytes"
	"github.com/akaspin/logx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStandaloneLogger_Debug_On(t *testing.T) {
	w := &bytes.Buffer{}
	l := logx.NewLog(logx.NewTextAppender(w, 0), "")
	l.Debug("test")
	assert.Equal(t, "DEBUG test\n", w.String())
}
