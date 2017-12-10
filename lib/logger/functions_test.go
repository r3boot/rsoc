package logger

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"
)

const (
	TMP_DIR    = "/tmp"
	TMP_PREFIX = "rsoc_lib_logger_functions_test"
)

var (
	MSGS_WITHOUT_TIMESTAMP map[string]string = map[string]string{
		LOG_INFO:    "I: info test without vars",
		LOG_DEBUG:   "D: debug test without vars",
		LOG_WARNING: "W: warning test without vars",
		LOG_FATAL:   "F: fatal test without vars",
	}

	reRFC3339TS = regexp.MustCompile("^(20[0-9]{2}-[01][0-9]-[0-2][0-9]T[012][0-9]:[0-6][0-9]:[0-6][0-9][A-Z0-9:\\+]{1,3})")
)

func NewTestLogger(t *testing.T, timestamp, debug bool) *Logger {
	fd, err := ioutil.TempFile(TMP_DIR, TMP_PREFIX)
	if err != nil {
		t.Errorf("ioutil.Tempfile failed: %v", err)
	}

	log := NewLogger(timestamp, debug)
	if log.UseTimestamp != timestamp {
		t.Errorf("log.UseTimestamp != timestamp")
	}
	if log.UseVerbose != debug {
		t.Errorf("log.UseVerbose != debug")
	}
	if log.UseDebug != debug {
		t.Errorf("log.UseDebug != debug")
	}

	log.TestFd = fd

	return log
}

func CleanupTestLogFile(t *testing.T, fd *os.File) {
	fname := fd.Name()
	err := fd.Close()
	if err != nil {
		t.Errorf("fd.Close() failed on test log")
	}

	err = os.Remove(fname)
	if err != nil {
		t.Errorf("os.Remove() failed on test log")
	}
}

func HasLine(t *testing.T, content []byte, wanted string, timestamp bool) bool {
	for _, line := range strings.Split(string(content), "\n") {
		if line == "" {
			continue
		}

		if timestamp {
			result := reRFC3339TS.FindAllStringSubmatch(line, -1)
			if len(result) == 0 {
				t.Errorf("timestamp == %v, but no timestamp found in line", timestamp)
			}
			line = line[len(result[0][0])+1:]
		}

		if line == wanted {
			return true
		}
	}

	return false
}

func RunLoggerTestsWith(t *testing.T, timestamp, debug bool) {
	log := NewTestLogger(t, timestamp, debug)
	defer CleanupTestLogFile(t, log.TestFd)

	log.Infof("info test without vars")
	log.Debugf("debug test without vars")
	log.Warningf("warning test without vars")
	log.Fatalf("fatal test without vars")

	content, err := ioutil.ReadFile(log.TestFd.Name())
	if err != nil {
		t.Errorf("ioutil.ReadFile failed on test log")
	}

	for key, value := range MSGS_WITHOUT_TIMESTAMP {
		switch key {
		case LOG_DEBUG:
			{
				if debug {
					if !HasLine(t, content, value, timestamp) {
						t.Errorf("debug == %v, but did not find message for debug level", debug)
					}
				} else {
					if HasLine(t, content, value, timestamp) {
						t.Errorf("debug == %v, but found a message for info debug", debug)
					}
				}
			}
		default:
			{
				if !HasLine(t, content, value, timestamp) {
					t.Errorf("did not find message for %s", value)
				}
			}
		}
	}
}

func TestLoggerFunctions(t *testing.T) {
	RunLoggerTestsWith(t, false, false)
	RunLoggerTestsWith(t, true, false)
	RunLoggerTestsWith(t, false, true)
	RunLoggerTestsWith(t, true, true)
}
