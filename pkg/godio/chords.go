package godio

import (
	"regexp"
)

type Chord struct {
	OctaveMod  string
	Root       string
	Type       string
	Extensions []string
}

var ChordFormulas = map[string][]int{
	"maj":  {0, 4, 7},
	"m":    {0, 3, 7},
	"dim":  {0, 3, 6},
	"aug":  {0, 4, 8},
	"sus2": {0, 2, 7},
	"sus4": {0, 5, 7},
}

var ExtensionInterval = map[string]int{
	"7":   10,
	"b9":  13,
	"9":   14,
	"#9":  15,
	"b11": 16,
	"11":  17,
	"#11": 18,
	"b13": 20,
	"13":  21,
}

func (c Chord) GetFrequencies() []float64 {
	formula := ChordFormulas[c.Type]
	var frequencies []float64
	rootNoteName := c.Root + "5"
	if c.OctaveMod == "l" {
		rootNoteName = c.Root + "4"
	}
	frequencies = append(frequencies, NoteFrequencies[rootNoteName])
	for _, ext := range c.Extensions {
		formula = append(formula, ExtensionInterval[ext])
	}
	for _, interval := range formula {
		frequencies = append(frequencies, NoteFrequencies.GetNoteFromInterval(rootNoteName, interval))
	}
	return frequencies
}

func ParseChord(chordStr string) Chord {
	regex := regexp.MustCompile(`(l?)([A-G][#b]?)((?:m|maj|dim|aug|sus2|sus4)?)(\d*)`)

	matches := regex.FindStringSubmatch(chordStr)

	// Extract components
	octaveModifier := matches[1]
	root := matches[2]
	chordType := matches[3]
	if chordType == "" {
		chordType = "maj"
	}
	extensions := matches[4]

	// Convert extensions to slice of strings
	var exts []string
	for _, char := range extensions {
		exts = append(exts, string(char))
	}

	return Chord{
		OctaveMod:  octaveModifier,
		Root:       root,
		Type:       chordType,
		Extensions: exts,
	}
}
