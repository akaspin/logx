package logx_test

import (
	"testing"
	"bytes"
	"github.com/tamtam-im/logx"
	"github.com/stretchr/testify/assert"
)

func TestStdSetOutput(t *testing.T) {
	w := &bytes.Buffer{}
	logx.SetOutput(w)
	logx.Info("test")
	assert.Contains(t, w.String(), "INFO std_test.go")
}

func TestStdGetLogger(t *testing.T) {
	w := &bytes.Buffer{}
	logx.SetOutput(w)
	logx.Info("test")

	l2 := logx.GetLog("second")
	l2.Info("test")
	assert.Contains(t, w.String(), "INFO std_test.go")
	assert.Contains(t, w.String(), "INFO second std_test.go")
}
