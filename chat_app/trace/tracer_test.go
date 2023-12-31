package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)

	if tracer == nil {
		t.Errorf("Newからの戻り値がnilです")
	} else {
		tracer.Trace("こんにちわ、traceパッケージ")
		if buf.String() != "こんにちわ、traceパッケージ\n" {
			t.Errorf("'%s'という誤った文字列が出力されました", buf.String())
		}
	}
}

func TestOff(t *testing.T) {
	var silentTracer Tracer = Off()

	silentTracer.Trace("データ")
}
