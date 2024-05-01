package tests

import (
	"reflect"
	"testing"
)

func assert(tb testing.TB, condition bool, msg string) {
	tb.Helper()
	if !condition {
		tb.Fatal(msg)
	}
}

func ok(tb testing.TB, err error) {
	tb.Helper()
	if err != nil {
		tb.Fatalf("%s failed: %s", tb.Name(), err)
	}
}

func equals(tb testing.TB, exp, act interface{}) {
	tb.Helper()
	if !reflect.DeepEqual(exp, act) {
		tb.Fatalf("%s doesnt equal %s", exp, act)
	}
}
