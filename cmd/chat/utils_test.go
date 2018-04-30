package main

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

var tests = []struct {
	s      []string
	search string
	want   bool
}{
	{[]string{"neil", "matt", "adam"}, "neil", true},
	{[]string{"neil", "matt", "adam"}, "neil1", false},
	{[]string{"neil", "matt", "adam"}, "adam", true},
}

func TestStringInSlice(t *testing.T) {
	for _, c := range tests {
		got, err := StringInSlice(c.s, c.search)
		Ok(t, err)
		Equals(t, c.want, got)
	}
}

/////////////////////////////////////////////////////////////////////////////////////////

// Assert fails the test if the condition is false.
func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d: "+msg+"\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// Ok fails the test if an err is not nil.
func Ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d: unexpected error: %s\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// Equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
