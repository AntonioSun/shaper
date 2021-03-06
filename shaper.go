////////////////////////////////////////////////////////////////////////////
// Porgram: shaper.go
// Purpose: mold strings into shape
// authors: Antonio Sun (c) 2016, All rights reserved
// Credits: Howard C. Shaw III
//          https://groups.google.com/d/msg/golang-nuts/snoIyANd-8c/V_IC57y4AwAJ
//          https://groups.google.com/d/msg/golang-nuts/snoIyANd-8c/hhOnu-lFAgAJ
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

// Copy returns a copy of the original object, instead of editing in-place,
// so make sure you've already got a reference to the original
// This should NEVER be hung off of a NewFilter string, or the original NewFilter will be lost
func (me *Shaper) Copy() *Shaper {
	return &Shaper{
		ShaperStack: me.ShaperStack,
	}
}

// Call this on the returned object to actually process a string
func (me *Shaper) Process(s string) string {
	return me.ShaperStack(s)
}

// Use this to apply arbitrary filters
func (me *Shaper) AddFilter(f func(string) string) *Shaper {
	me.ShaperStack = func(a func(string) string, b func(string) string) func(string) string {
		return func(s string) string {
			return a(b(s))
		}
	}(f, me.ShaperStack)
	return me
}

func (me *Shaper) ApplyToLower() *Shaper {
	me.AddFilter(strings.ToLower)
	return me
}

func (me *Shaper) ApplyToUpper() *Shaper {
	me.AddFilter(strings.ToUpper)
	return me
}

func (me *Shaper) ApplyReplace(old, new string, times int) *Shaper {
	me.AddFilter(func(s string) string {
		return strings.Replace(s, old, new, times)
	})
	return me
}

func (me *Shaper) ApplyRegexpReplaceAll(rexp, repl string) *Shaper {
	me.AddFilter(func(s string) string {
		return regexp.MustCompile(rexp).ReplaceAllString(s, repl)
	})
	return me
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
