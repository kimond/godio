package godio

import (
	"math"
	"slices"
	"strings"
	"unicode"

	"github.com/samber/lo"
)

type Chord struct {
	Root         string
	Quality      string
	BassNote     string
	Extensions   []string
	VoicingRules []VoicingRule
	Tones        []int
}

var chordFormulas = map[string][]int{
	"":        {4, 7}, // It's reminiscent of figured bass with chromatic intervals instead of diatonic ones.
	"maj":     {4, 7},
	"aug":     {4, 8},
	"dim":     {3, 6},
	"dim7":    {3, 6, 9},
	"m":       {3, 7},
	"m7":      {3, 7, 10},
	"7":       {4, 7, 10},
	"maj7":    {4, 7, 11},
	"9":       {4, 7, 10, 2},
	"m9":      {3, 7, 10, 2},
	"maj9":    {4, 7, 11, 2},
	"13":      {4, 7, 10, 2, 9},
	"m11":     {3, 7, 10, 14, -7}, // m11 chords are strange, and should be spaced properly. Maybe they shouldn't even have a 9.
	"6":       {4, 7, 9},
	"m6":      {3, 7, 9},
	"69":      {4, 7, 9, 2},
	"m69":     {3, 7, 9, 2},
	"mmaj7":   {3, 7, 11},
	"minmaj7": {3, 7, 11},
}

var extensionFormulas = map[string]int{
	"#1":   1,
	"#15":  1,
	"b9":   1,
	"9":    2,
	"#9":   3,
	"11":   5,
	"ll":   5,
	"#11":  6,
	"b5":   6,
	"#5":   8,
	"b13":  8,
	"13":   9,
	"sus":  5 - 12, // By having "sus" as an extension with a negative value, chords like C9sus will be properly parsed
	"sus2": 2 - 12, // but still put the fourth of the chord lower as is typical for this kind of chord
	"sus4": 5 - 12,
}

func (c *Chord) addTone(tone int) {
	c.Tones = append(c.Tones, tone)
}

func (c *Chord) addTones(tones []int) {
	c.Tones = append(c.Tones, tones...)
}

func (c *Chord) replaceTone(tone int, newTone int) {
	for i, t := range c.Tones {
		if t == tone {
			c.Tones[i] = newTone
		}
	}
}

func (c *Chord) removeTone(tone int) {
	c.Tones = lo.Filter(c.Tones, func(i int, _ int) bool {
		return i != tone
	})
}

func (c Chord) hasInterval(interval int) bool {
	formula := chordFormulas[c.Quality]
	for _, i := range formula {
		if i == interval {
			return true
		}
	}
	return false
}

func (c Chord) GetFrequencies() []float64 {
	var frequencies []float64
	var bassNoteName string
	if c.BassNote != "" {
		bassNoteName = c.BassNote + "2"
	} else {
		bassNoteName = c.Root + "2"
	}
	rootNoteName := c.Root + "3"

	c.applyVoicingRules()

	frequencies = append(frequencies, NoteFrequencies[bassNoteName])

	for _, interval := range c.Tones {
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

func noteToNumber(note string) int {
	notes := map[byte]int{
		'C': 12,
		'D': 14,
		'E': 16,
		'F': 17,
		'G': 19,
		'A': 21,
		'B': 23,
	}
	accidentals := map[rune]int{
		'b': -1,
		'#': 1,
	}
	value := notes[note[0]]
	for _, item := range note[1:] {
		value += accidentals[item]
	}
	return value
}

func (c Chord) GetFrequenciesV2() []float64 {
	spreadExtensions := false
	lowRoot := true
	// Above is a list of parameters for the chord voicing that change how voicings are done
	extensionOctave := 48
	if spreadExtensions {
		extensionOctave = 60
	}

	root := noteToNumber(c.Root)
	bassNote := noteToNumber(c.BassNote)

	chordTones := chordFormulas[c.Quality]
	if slices.Contains(c.Extensions, "sus") || slices.Contains(c.Extensions, "sus2") || slices.Contains(c.Extensions, "sus4") {
		for index, tone := range chordTones {
			if tone == 3 || tone == 4 {
				// remove at index
				chordTones = slices.Delete(chordTones, index, index+1)
			}
		}
	}
	if slices.Contains(c.Extensions, "b5") || slices.Contains(c.Extensions, "#5") {
		for index, tone := range chordTones {
			if tone == 7 {
				chordTones = slices.Delete(chordTones, index, index+1)
			}
		}
	}

	// The following code generates a voicing with no regard to what the previous chord was.
	// In the python code, a much longer chunk of code uses an iterative algorithm to find
	// a combination of octaves for the notes that minimizes the sum of the distances
	// of each note in the new chord to the nearest note in the previous chord.

	voicing := []int{}

	for _, tone := range chordTones {
		voicing = append(voicing, root+48+tone)
	}
	for _, extension := range c.Extensions {
		voicing = append(voicing, root+extensionOctave+extensionFormulas[extension])
	}
	if lowRoot && !slices.Contains(chordTones, 6) && !slices.Contains(chordTones, 8) && !slices.Contains(c.Extensions, "#5") && !slices.Contains(c.Extensions, "b5") {
		voicing = append(voicing, root+43)
	}
	voicing = append(voicing, root+48)
	if lowRoot {
		voicing = append(voicing, bassNote+36)
		voicing = append(voicing, bassNote+24)
	}

	// The chord is now voiced in MIDI note numbers.
	freqVoicing := []float64{}

	for _, m := range voicing {
		freqVoicing = append(freqVoicing, math.Pow(2, float64((float64(m)-69.0)/12.0))*440.0)
	}
	return freqVoicing
}

func (c *Chord) applyVoicingRules() {
	for _, rule := range c.VoicingRules {
		if rule.Condition(c) {
			rule.Action(c)
		}
	}
}

func ParseChord(chordStr string) *Chord {
	finding := "root"
	chord := &Chord{VoicingRules: defaultVoicingRules}
	accumulator := ""
	extensions := []string{}
	for _, char := range chordStr {

		if finding == "bass" {
			accumulator += string(char) // string concatenation
		}
		if char == '/' { // any part of the chord can suddenly end like this then go to bass
			switch finding {
			case "root":
				chord.Root = accumulator
			case "quality":
				chord.Quality = accumulator
			case "extensions":
				extensions = append(extensions, accumulator)
				chord.Extensions = extensions
			}
			finding = "bass"
			accumulator = ""
		}
		if finding == "root" {
			if (!unicode.IsDigit(char) && !unicode.IsLower(char)) || char == 'b' {
				accumulator += string(char)
			} else if char == 's' { // sus
				finding = "extensions"
				chord.Root = accumulator
				accumulator = ""
			} else {
				finding = "quality"
				chord.Root = accumulator
				accumulator = ""

			}
		}
		if finding == "quality" {
			if len(accumulator) > 2 {
				if accumulator[len(accumulator)-3:] == "add" {
					chord.Quality = accumulator[:len(accumulator)-3]
					finding = "extensions"
					accumulator = ""
				}
			}
			if finding == "quality" {
				if !slices.Contains([]rune{'#', 'b', '(', 's'}, char) {
					accumulator += string(char)
				} else {
					finding = "extensions"
					chord.Quality = accumulator
					accumulator = ""
				}
			}
		}
		if finding == "extensions" {
			if slices.Contains([]rune{'#', 'b'}, char) {
				if accumulator != "" {
					extensions = append(extensions, accumulator)
				}
				accumulator = string(char)
			} else if char == '(' {
				if accumulator != "" {
					extensions = append(extensions, accumulator)
				}
				accumulator = ""
			} else if slices.Contains([]rune{',', ')'}, char) {
				extensions = append(extensions, accumulator)
				accumulator = ""
			} else {
				accumulator += string(char)
			}
		}
	}
	switch finding {
	case "root":
		chord.Root = accumulator
	case "quality":
		chord.Quality = accumulator
	case "extensions":
		if accumulator != "" {
			extensions = append(extensions, accumulator)
		}
		chord.Extensions = extensions
	case "bass": // rare "bass case"
		chord.BassNote = accumulator
	}
	if chord.BassNote == "" {
		chord.BassNote = chord.Root
	}

	if strings.Contains(chord.Root, "b") {
		chord.Root = flatToSharp[chord.Root]
	}

	if strings.Contains(chord.BassNote, "b") {
		chord.BassNote = flatToSharp[chord.BassNote]
	}

	chord.Tones = append(chord.Tones, 0)

	chord.Tones = append(chord.Tones, chordFormulas[chord.Quality]...)

	for _, extension := range chord.Extensions {
		chord.Tones = append(chord.Tones, extensionFormulas[extension])
	}
	return chord
}
