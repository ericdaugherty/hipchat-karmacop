package main

import (
	"regexp"
	"testing"
)

// var secretKarmaReg = regexp.MustCompile(`^s/[[:print:]]*\s*[[:print:]]*@([^/])+(--|\Q++\E)`)
var secretKarmaReg = regexp.MustCompile(`^s/.*`)

var shouldMatch = []string{
	"s/@Eric--/",
	"s/@eric++",
	"s/@Eric --/",
	"s/@eric ++",
	"s/@bob ++/",
	"s/@pop ++/",
	"s/@bbbbb ++/",
	"s/@eric ++/",
	"s/@Eric --/qwerty",
	"s/@bob--/asdf",
	"s/@eric ++/123e2s*^%",
	"s/stupid/@Eric++",
	"s/@Eric--/stupid",
	"s/@bob -------/asdf",
	"s/@eric +++++/more",
	"s/@✆̩̺ͧͨͯ̋̌̈̉☠̳̫͎̗̗̒̅̃͒♔̱̻̤ͭ̊̉̓ ++/test",
	`s/asdasdas/

@Eric---`,
	`s/afdsadfsafda/

ü
@TheRealDeal -----`,
}
var shouldNotMatch = []string{
	"s/hello",
	"goodbye",
}

// Doesn't test live code but used to test regex used in descriptor.
func Test_regex(t *testing.T) {

	for _, val := range shouldMatch {
		if !secretKarmaReg.MatchString(val) {
			t.Errorf("String %s should have matched but didn't.", val)
		}
	}

	// for _, val := range shouldNotMatch {
	// 	if secretKarmaReg.MatchString(val) {
	// 		t.Errorf("String %s should not have have matched but did.", val)
	// 	}
	// }
}

func Test_checkMessage(t *testing.T) {

	for _, val := range shouldMatch {
		if !checkMessage(val) {
			t.Errorf("String %s should have matched but didn't.", val)
		}
	}

	for _, val := range shouldNotMatch {
		if checkMessage(val) {
			t.Errorf("String %s should not have have matched but did.", val)
		}
	}
}
