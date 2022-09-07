package validation_test

import (
	"testing"

	"github.com/PacodiazDG/Backend-blog/Modules/validation"
)

type SliceStringContainsTest struct {
	arg1     []string
	arg2     string
	expected bool
}

func TestSliceStringContains(t *testing.T) {
	var addTests = []SliceStringContainsTest{
		{[]string{"Banned", "MLA", "Locked"}, "Bannned", false},
		{[]string{"Banned", "MLA", "Locked"}, "Banned", true},
		{[]string{"Banned", "MLA", "Locked"}, "Locked", true},
		{[]string{"Banned", "MLA", "Locked"}, "MLAA", false},
	}
	for _, test := range addTests {
		if output := validation.SliceStringContains(test.arg1, test.arg2); output != test.expected {
			t.Errorf("Output %t not equal to expected %t", output, test.expected)
		}
	}
}
