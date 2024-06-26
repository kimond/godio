package godio

import (
	"regexp"
	"strings"

	"github.com/samber/lo"
)

type Chord struct {
	Root        string
	Type        string
	Extension   string
	Alterations []string
	Addition    string
	BassNote    string
}

var ChordFormulas = map[string][]int{
	"maj":  {0, 4, 7},
	"maj7": {0, 4, 7, 11},
	"m":    {0, 3, 7},
	"dim":  {0, 3, 6},
	"aug":  {0, 4, 8},
	"sus2": {0, 2, 7},
	"sus4": {0, 5, 7},
}

var ExtensionIntervals = map[string][]int{
	"6":    {9},
	"7":    {10},
	"maj7": {11},
	"9":    {10, 14},
	"11":   {10, 14, 17},
	"13":   {10, 14, 21},
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
	"addb9":  13,
	"add9":   14,
	"add#9":  15,
	"addb11": 16,
	"add11":  17,
	"add#11": 18,
	"addb13": 20,
	"add13":  21,
	"add#13": 22,
}

func (c Chord) GetFrequencies() []float64 {
	formula := ChordFormulas[c.Type]
	var frequencies []float64
	var bassNoteName string
	if c.BassNote != "" {
		bassNoteName = c.BassNote + "2"
	} else {
		bassNoteName = c.Root + "2"
	}
	rootNoteName := c.Root + "3"

	frequencies = append(frequencies, NoteFrequencies[bassNoteName])
	if c.Extension != "" {
		formula = append(formula, ExtensionIntervals[c.Extension]...)
		if c.Extension == "11" {
			// remove the 3rd if it's an 11
			formula = lo.Filter(formula, func(i int, _ int) bool {
				return i != 3 && i != 4
			})
		}
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

	if len(formula)+1 >= 6 {
		// Remove the upper root note if chord has more than 5 notes
		formula = lo.Filter(formula, func(i int, _ int) bool {
			return i != 0
		})
	}

	if len(formula)+1 >= 6 && !hasFifthAlteration {
		// Remove the fifth if chord has more than 6 notes and no alteration of the fifth
		formula = lo.Filter(formula, func(i int, _ int) bool {
			return i != 7
		})
	}

	for _, interval := range formula {
		intervalFrequency := NoteFrequencies.GetNoteFromInterval(rootNoteName, interval)
		if intervalFrequency > NoteFrequencies["F#4"] {
			intervalFrequency = NoteFrequencies.GetNoteFromInterval(rootNoteName, interval-12)
		}
		if intervalFrequency < NoteFrequencies["G3"] {
			intervalFrequency = NoteFrequencies.GetNoteFromInterval(rootNoteName, interval+12)
		}
		frequencies = append(frequencies, intervalFrequency)
	}
	return frequencies
}

func ParseChord(chordStr string) Chord {
	regex := regexp.MustCompile(`([A-G][#b]?)((?:maj7?|m|dim|aug|sus2|sus4)?)((?:6|7|maj7|9|11|13)?)((?:[#b]\d{1,2})*)((?:add[#b]?\d{1,2})?)((/[A-G][#b]?)?)`)

	matches := regex.FindStringSubmatch(chordStr)

	root := matches[1]
	if strings.Contains(root, "b") {
		root = flatToSharp[root]
	}
	chordType := matches[2]
	if chordType == "" {
		chordType = "maj"
	}
	extension := matches[3]
	alterationMatch := matches[4]
	addition := matches[5]
	bassNote := ""
	if matches[6] != "" {
		bassNote = matches[6][1:]
		if strings.Contains(bassNote, "b") {
			bassNote = flatToSharp[bassNote]
		}
	}

	alterations := regexp.MustCompile(`[#b]\d{1,2}`).FindAllString(alterationMatch, -1)

	return Chord{
		Root:        root,
		Type:        chordType,
		Extension:   extension,
		Alterations: alterations,
		Addition:    addition,
		BassNote:    bassNote,
	}
}
