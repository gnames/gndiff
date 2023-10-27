package fuzzy

import (
	"strings"

	"github.com/gnames/levenshtein/ent/editdist"
)

const (
	charsPerED      = 5
	maxEditDistance = 3
)

// EditDistance calculates edit distance (**ed**) according to Levenshtein algorithm.
// It also runs additional checks and if they fail, returns -1.
//
// Checks:
// - result should not exceed maxEditDistance
// - number of characters divided by ed should be bigger than charsPerED
//
// It assumes that checks have to be applied only to the second string:
//
//	EditDistance("Pomatomus", "Pom atomus")
//
// returns -1
//
//	EditDistance("Pom atomus", "Pomatomus")
//
// returns 1
//
// It also assumes that number of spaces between words was already
// normalized to 1 space, and that s1 and s2 always have the same number of
// words.
func EditDistance(str, ref string, noCheck bool) int {
	ed, _, _ := editdist.ComputeDistance(str, ref, false)
	if ed == 0 {
		return ed
	}

	if ed > maxEditDistance {
		return -1
	}

	if noCheck {
		return ed
	}

	return checkED(str, ref, ed)
}

func checkED(str, ref string, ed int) int {
	curEd := ed
	wRef := strings.Split(str, " ")
	wStr := strings.Split(ref, " ")

	// if strings have different number of words, we
	// cannot compare them word by word
	if len(wRef) != len(wStr) {
		return ed
	}

	// check only words from the second string
	for i, w := range wStr {
		r := []rune(w)

		// check edit distance for a word, if number
		// of letters per edit distance is smaller than allowed
		// number.:
		if len(r)/curEd < charsPerED {
			wordEd, _, _ := editdist.ComputeDistance(w, wRef[i], false)
			if wordEd > 0 && len(r)/wordEd < charsPerED {
				return -1
			}
			curEd -= wordEd
		}
	}
	return ed
}
