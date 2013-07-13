package index

import (
	"github.com/daviddengcn/go-algs/ed"
	"strings"
	"testing"
)

func showText(text string) string {
	return text + "."
}

func AssertTextEquals(t *testing.T, name, act, exp string) {
	if exp == act {
		return
	}

	expLines := strings.Split(exp, "\n")
	actLines := strings.Split(act, "\n")

	t.Errorf("%s unexpected(exp: %d lines, act %d lines!", name, len(expLines), len(actLines))
	t.Logf("exp ---  act +++")
	t.Logf("Difference:")
	_, matA, matB := ed.EditDistanceFFull(len(expLines), len(actLines), func(iA, iB int) int {
		sa, sb := expLines[iA], actLines[iB]
		if sa == sb {
			return 0
		}
		return ed.String(sa, sb)
	}, func(iA int) int {
		return len(expLines[iA]) + 1
	}, func(iB int) int {
		return len(actLines[iB]) + 1
	})
	for i, j := 0, 0; i < len(expLines) || j < len(actLines); {
		switch {
		case j >= len(actLines) || i < len(expLines) && matA[i] < 0:
			t.Logf("--- %3d: %s", i+1, showText(expLines[i]))
			i++
		case i >= len(expLines) || j < len(actLines) && matB[j] < 0:
			t.Logf("+++ %3d: %s", j+1, showText(actLines[j]))
			j++
		default:
			if expLines[i] != actLines[j] {
				t.Logf("--- %3d: %s", i+1, showText(expLines[i]))
				t.Logf("+++ %3d: %s", j+1, showText(actLines[j]))
			} // else
			i++
			j++
		}
	} // for i, j
}

