package index

import (
	"github.com/daviddengcn/go-algs/ed"
	"strings"
	"testing"
)

func AssertStringsEquals(t *testing.T, name string, act, exp []string) {
	t.Errorf("%s unexpected(exp: %d lines, act %d lines!", name, len(exp), len(act))
	t.Logf("exp ---  act +++")
	t.Logf("Difference:")
	_, matA, matB := ed.EditDistanceFFull(len(exp), len(act), func(iA, iB int) int {
		sa, sb := exp[iA], act[iB]
		if sa == sb {
			return 0
		}
		return ed.String(sa, sb)
	}, func(iA int) int {
		return len(exp[iA]) + 1
	}, func(iB int) int {
		return len(act[iB]) + 1
	})
	for i, j := 0, 0; i < len(exp) || j < len(act); {
		switch {
		case j >= len(act) || i < len(exp) && matA[i] < 0:
			t.Logf("--- %3d: %s", i+1, showText(exp[i]))
			i++
		case i >= len(exp) || j < len(act) && matB[j] < 0:
			t.Logf("+++ %3d: %s", j+1, showText(act[j]))
			j++
		default:
			if exp[i] != act[j] {
				t.Logf("--- %3d: %s", i+1, showText(exp[i]))
				t.Logf("+++ %3d: %s", j+1, showText(act[j]))
			} // else
			i++
			j++
		}
	} // for i, j
}

func showText(text string) string {
	return text + "."
}

func AssertTextEquals(t *testing.T, name, act, exp string) {
	if exp == act {
		return
	}

	expLines := strings.Split(exp, "\n")
	actLines := strings.Split(act, "\n")
	
	AssertStringsEquals(t, name, actLines, expLines)
}

