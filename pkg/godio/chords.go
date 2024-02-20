package godio

import (
	"github.com/samber/lo"
	"regexp"
)

type Chord struct {
	OctaveMod   string
	Root        string
	Type        string
	Extension   string
	Alterations []string
	Addition    string
}

var ChordFormulas = map[string][]int{
	"maj":  {4, 7},
	"m":    {3, 7},
	"dim":  {3, 6},
	"aug":  {4, 8},
	"sus2": {2, 7},
	"sus4": {5, 7},
}

var ExtensionIntervals = map[string][]int{
	"6":  {9},
	"7":  {10},
	"9":  {10, 14},
	"11": {10, 14, 17},
	"13": {10, 14, 17, 21},
}

var AlterationMap = map[string]int{
	"5":  7,
	"6":  9,
	"7":  10,
	"9":  14,
	"11": 17,
	"13": 21,
}

var AdditionalNotes = map[string]int{
	"add9":  14,
	"add11": 17,
	"add13": 21,
}

func (c Chord) GetFrequencies() []float64 {
	formula := ChordFormulas[c.Type]
	var frequencies []float64
	bassNoteName := c.Root + "2"
	rootNoteName := c.Root + "4"
	frequencies = append(frequencies, NoteFrequencies[bassNoteName])
	if c.Extension != "" {
		formula = append(formula, ExtensionIntervals[c.Extension]...)
	}
	if c.Addition != "" {
		formula = append(formula, AdditionalNotes[c.Addition])
	}

	hasFifthAlteration := false
	for _, alteration := range c.Alterations {
		mod, degree := alteration[0], alteration[1:]
		if degree == "5" {
			hasFifthAlteration = true
		}
		var newValue int
		if mod == '#' {
			newValue = AlterationMap[degree] + 1
		} else if mod == 'b' {
			newValue = AlterationMap[degree] - 1
		}

		for i, interval := range formula {
			if interval == AlterationMap[degree] {
				formula[i] = newValue
			}
		}
	}

	if len(formula)+1 >= 6 && !hasFifthAlteration {
		// Remove the fifth if chord has more than 6 notes and no alteration of the fifth
		formula = lo.Filter(formula, func(i int, _ int) bool {
			return i != 7
		})
	}

	for _, interval := range formula {
		frequencies = append(frequencies, NoteFrequencies.GetNoteFromInterval(rootNoteName, interval))
	}
	return frequencies
}

func ParseChord(chordStr string) Chord {
	regex := regexp.MustCompile(`(l?)([A-G][#b]?)((?:maj|m|dim|aug|sus2|sus4)?)((?:6|7|9|11|13)?)((?:[#b]\d{1,2})*)((?:add\d{1,2})?)`)

	matches := regex.FindStringSubmatch(chordStr)

	octaveModifier := matches[1]
	root := matches[2]
	chordType := matches[3]
	if chordType == "" {
		chordType = "maj"
	}
	extension := matches[4]
	alterationMatch := matches[5]
	addition := matches[6]

	alterations := regexp.MustCompile(`[#b]\d{1,2}`).FindAllString(alterationMatch, -1)

	return Chord{
		OctaveMod:   octaveModifier,
		Root:        root,
		Type:        chordType,
		Extension:   extension,
		Alterations: alterations,
		Addition:    addition,
	}
}
