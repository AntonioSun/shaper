////////////////////////////////////////////////////////////////////////////
// Porgram: shaper.go
// Purpose: mold strings into shape
// authors: Antonio Sun (c) 2016, All rights reserved
// Credits: https://groups.google.com/d/msg/golang-nuts/snoIyANd-8c/V_IC57y4AwAJ
////////////////////////////////////////////////////////////////////////////

package shaper

import (
	"fmt"
	"regexp"
	"strings"
)

////////////////////////////////////////////////////////////////////////////
// Constant and data type/structure definitions

type Shaper struct {
	ShaperStack func(string) string
}

////////////////////////////////////////////////////////////////////////////
// Global variables definitions

////////////////////////////////////////////////////////////////////////////
// Function definitions

func PassThrough(s string) string {
	return s
}

// Make a new Shaper filter and start adding bits
func NewFilter() *Shaper {
	return &Shaper{ShaperStack: PassThrough}
}

// Call this on the returned object to actually process a string
func (m *Shaper) Process(s string) string {
	return m.ShaperStack(s)
}

// Use this to apply arbitrary filters
func (m *Shaper) AddFilter(f func(string) string) *Shaper {
	m.ShaperStack = func(a func(string) string, b func(string) string) func(string) string {
		return func(s string) string {
			return a(b(s))
		}
	}(f, m.ShaperStack)
	return m
}

func (m *Shaper) ApplyToLower() *Shaper {
	m.AddFilter(strings.ToLower)
	return m
}

func (m *Shaper) ApplyToUpper() *Shaper {
	m.AddFilter(strings.ToUpper)
	return m
}

func (m *Shaper) ApplyReplace(old, new string, times int) *Shaper {
	m.AddFilter(func(s string) string {
		return strings.Replace(s, old, new, times)
	})
	return m
}

func (m *Shaper) ApplyRegexpReplaceAll(rexp, repl string) *Shaper {
	m.AddFilter(func(s string) string {
		return regexp.MustCompile(rexp).ReplaceAllString(s, repl)
	})
	return m
}

func ShaperDemo() {
	// Construct pipelines
	UpCase := NewFilter().ApplyToUpper()
	LCase := NewFilter().ApplyToLower()
	Replace := NewFilter().ApplyReplace("test", "biscuit", -1)
	RU := NewFilter().ApplyReplace("test", "biscuit", -1).ApplyToUpper()

	// Test pipelines
	fmt.Printf("%s\n", UpCase.Process("This is a test."))
	fmt.Printf("%s\n", LCase.Process("This is a test."))
	fmt.Printf("%s\n", Replace.Process("This is a test."))
	fmt.Printf("%s\n", RU.Process("This is a test."))

	// Note that we can reuse these stacks as many times as we like
	fmt.Printf("%s\n", Replace.Process("This is also a test. Testificate."))

	// We can also add stages later on - though we cannot remove stages using this style
	Replace.ApplyToUpper()
	fmt.Printf("%s\n", Replace.Process("This is also a test. Testificate."))
	LCase.ApplyReplace("test", "biscuit", -1)
	fmt.Printf("%s\n", LCase.Process("This is also a test. Testificate."))

	// Regexp.ReplaceAll
	RegReplace := NewFilter().ApplyRegexpReplaceAll("(?i)ht(ml)", "X$1")
	fmt.Printf("%s\n", RegReplace.Process("This is html Html HTML."))

	fmt.Printf("Finished.\n")
	/*

		Output :

		THIS IS A TEST.
		this is a test.
		This is a biscuit.
		THIS IS A BISCUIT.
		This is also a biscuit. Testificate.
		THIS IS ALSO A BISCUIT. TESTIFICATE.
		this is also a biscuit. biscuitificate.
		This is Xml Xml XML.
		Finished.

	*/

}
