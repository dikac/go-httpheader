package httpheader_test

import (
	"fmt"
	"github.com/dikac/go-httpheader"
	"net/http"
	"sort"
	"time"
)

func ExampleHeader() {
	type Options struct {
		ContentType  string `header:"Content-Type"`
		Length       int
		Bool         bool
		BoolInt      bool      `header:"Bool-Int,int"`
		XArray       []string  `header:"X-Array"`
		TestHide     string    `header:"-"`
		IgnoreEmpty  string    `header:"X-Empty,omitempty"`
		IgnoreEmptyN string    `header:"X-Empty-N,omitempty"`
		CreatedAt    time.Time `header:"Created-At"`
		UpdatedAt    time.Time `header:"Update-At,unix"`
		CustomHeader http.Header
	}

	opt := Options{
		ContentType:  "application/json",
		Length:       2,
		Bool:         true,
		BoolInt:      true,
		XArray:       []string{"test1", "test2"},
		TestHide:     "hide",
		IgnoreEmptyN: "n",
		CreatedAt:    time.Date(2000, 1, 1, 12, 34, 56, 0, time.UTC),
		UpdatedAt:    time.Date(2001, 1, 1, 12, 34, 56, 0, time.UTC),
		CustomHeader: http.Header{
			"X-Test-1": []string{"233"},
			"X-Test-2": []string{"666"},
		},
	}
	h, err := httpheader.Header(opt)
	fmt.Println(err)
	printHeader(h)
	// Output:
	// <nil>
	// Bool: []string{"true"}
	// Bool-Int: []string{"1"}
	// Content-Type: []string{"application/json"}
	// Created-At: []string{"Sat, 01 Jan 2000 12:34:56 GMT"}
	// Length: []string{"2"}
	// Update-At: []string{"978352496"}
	// X-Array: []string{"test1", "test2"}
	// X-Empty-N: []string{"n"}
	// X-Test-1: []string{"233"}
	// X-Test-2: []string{"666"}
}

func printHeader(h http.Header) {
	var keys []string
	for k := range h {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s: %#v\n", k, h[k])
	}
}

func ExampleDecode() {
	type Options struct {
		ContentType  string `header:"Content-Type"`
		Length       int
		Bool         bool
		BoolInt      bool      `header:"Bool-Int,int"`
		XArray       []string  `header:"X-Array"`
		TestHide     string    `header:"-"`
		IgnoreEmpty  string    `header:"X-Empty,omitempty"`
		IgnoreEmptyN string    `header:"X-Empty-N,omitempty"`
		CreatedAt    time.Time `header:"Created-At"`
		UpdatedAt    time.Time `header:"Update-At,unix"`
		CustomHeader http.Header
	}
	h := http.Header{
		"Bool":         []string{"true"},
		"Bool-Int":     []string{"1"},
		"X-Test-1":     []string{"233"},
		"X-Test-2":     []string{"666"},
		"Content-Type": []string{"application/json"},
		"Length":       []string{"2"},
		"X-Array":      []string{"test1", "test2"},
		"X-Empty-N":    []string{"n"},
		"Update-At":    []string{"978352496"},
		"Created-At":   []string{"Sat, 01 Jan 2000 12:34:56 GMT"},
	}
	var opt Options
	err := httpheader.Decode(h, &opt, "header", []string{http.TimeFormat})
	fmt.Println(err)
	fmt.Println(opt.ContentType)
	fmt.Println(opt.Length)
	fmt.Println(opt.BoolInt)
	fmt.Println(opt.XArray)
	fmt.Println(opt.UpdatedAt)
	// Output:
	// <nil>
	// application/json
	// 2
	// true
	// [test1 test2]
	// 2001-01-01 12:34:56 +0000 UTC
}

type EncodedArgs []string

func (e EncodedArgs) EncodeHeader(key string, v *http.Header) error {
	for i, arg := range e {
		v.Set(fmt.Sprintf("%s.%d", key, i), arg)
	}
	return nil
}

type DecodeArg struct {
	arg string
}

func (d *DecodeArg) DecodeHeader(header http.Header, key string) error {
	value := header.Get(key)
	d.arg = value
	return nil
}

func ExampleEncoder() {
	// type EncodedArgs []string
	//
	// func (e EncodedArgs) EncodeHeader(key string, v *http.Header) error {
	// 	for i, arg := range e {
	// 	v.Set(fmt.Sprintf("%s.%d", key, i), arg)
	// }
	// 	return nil
	// }

	s := struct {
		Args EncodedArgs `header:"Args"`
	}{Args: EncodedArgs{"a", "b", "c"}}

	h, err := httpheader.Header(s)
	fmt.Println(err)
	printHeader(h)
	// Output:
	// <nil>
	// Args.0: []string{"a"}
	// Args.1: []string{"b"}
	// Args.2: []string{"c"}
}

func ExampleDecoder() {
	// type DecodeArg struct {
	// 	arg string
	// }
	//
	// func (d *DecodeArg) DecodeHeader(header http.Header, key string) error {
	// 	value := header.Get(key)
	// 	d.arg = value
	// 	return nil
	// }

	var s struct {
		Arg DecodeArg
	}
	h := http.Header{}
	h.Set("Arg", "foobar")
	err := httpheader.Decode(h, &s, "header", []string{http.TimeFormat})
	fmt.Println(err)
	fmt.Println(s.Arg.arg)
	// Output:
	// <nil>
	// foobar
}
