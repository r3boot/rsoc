package logger

import "testing"

func runtestNewLoggerWithParams(t *testing.T, timestamp, debug bool) {
	log := NewLogger(timestamp, debug)
	if log == nil {
		t.Errorf("log == nil")
	}

	if log.UseTimestamp != timestamp {
		t.Errorf("log.UseTimestamp != timestamp")
	}

	if log.UseDebug != debug {
		t.Errorf("log.UseDebug != debug")
	}

	if log.UseVerbose != debug {
		t.Errorf("log.UseVerbose != debug")
	}
}

func TestNewLogger(t *testing.T) {
	runtestNewLoggerWithParams(t, false, false)
	runtestNewLoggerWithParams(t, true, false)
	runtestNewLoggerWithParams(t, false, true)
	runtestNewLoggerWithParams(t, true, true)
}
