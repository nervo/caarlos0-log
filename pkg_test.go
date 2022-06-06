package log_test

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/caarlos0/log"
)

type Pet struct {
	Name string
	Age  int
}

func (p *Pet) Fields() log.Fields {
	return log.Fields{
		"name": p.Name,
		"age":  p.Age,
	}
}

func TestRootLogOptions(t *testing.T) {
	var out bytes.Buffer
	log.SetHandler(log.New(&out))
	log.SetLevel(log.DebugLevel)
	log.SetLevelFromString("info")
	log.WithError(fmt.Errorf("here")).Info("a")
	log.Debug("debug")
	log.Debugf("warn %d", 1)
	log.Info("info")
	log.Infof("warn %d", 1)
	log.Warn("warn")
	log.Warnf("warn %d", 1)
	log.Error("error")
	log.Errorf("warn %d", 1)
	log.WithField("foo", "bar").Info("foo")
	pet := &Pet{"Tobi", 3}
	log.WithFields(pet).Info("add pet")
	requireEqualOutput(t, out.Bytes())
}

// Unstructured logging is supported, but not recommended since it is hard to query.
func Example_unstructured() {
	log.Infof("%s logged in", "Tobi")
}

// Structured logging is supported with fields, and is recommended over the formatted message variants.
func Example_structured() {
	log.WithField("user", "Tobo").Info("logged in")
}

// Errors are passed to WithError(), populating the "error" field.
func Example_errors() {
	err := errors.New("boom")
	log.WithError(err).Error("upload failed")
}

// Multiple fields can be set, via chaining, or WithFields().
func Example_multipleFields() {
	log.WithFields(log.Fields{
		"user": "Tobi",
		"file": "sloth.png",
		"type": "image/png",
	}).Info("upload")
}

var update = flag.Bool("update", false, "update .golden files")

func requireEqualOutput(tb testing.TB, bts []byte) {
	tb.Helper()

	golden := "testdata/" + tb.Name() + ".golden"
	if *update {
		if err := os.MkdirAll(filepath.Dir(golden), 0o755); err != nil {
			tb.Fatal(err)
		}
		if err := os.WriteFile(golden, bts, 0o600); err != nil {
			tb.Fatal(err)
		}
	}

	gbts, err := os.ReadFile(golden)
	if err != nil {
		tb.Fatal(err)
	}

	sg := format(string(gbts))
	so := format(string(bts))
	if sg != so {
		tb.Fatalf("output do not match:\ngot:\n%s\n\nexpected:\n%s\n\n", so, sg)
	}
}

func format(str string) string {
	eol := "\n"
	if runtime.GOOS == "windows" {
		eol = "\r\n"
	}
	return strings.ReplaceAll(strings.ReplaceAll(str, "\x1b", "\\x1b"), "\n", eol)
}
