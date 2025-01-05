package godio

import (
	"math"
	"regexp"
	"slices"
	"strings"
	"unicode"

	"github.com/samber/lo"
)

type Chord struct {
	Root         string
	Type         string
	Extension    string
	Alterations  []string
	Addition     string
	BassNote     string
	VoicingRules []VoicingRule
	Tones        []Interval
}

type ChalkChord struct { //All of my python code copied here will have its types and functions preceeded with "Chalk"
	Root       int
	BassNote   int
	Quality    string
	Extensions []string
	Tones      []int
}

var ChalkChordFormulas = map[string][]int{ //This kind of chord formula uses numbers instead of an intermediary "interval" type.
	"":        {4, 7}, //It's reminiscent of figured bass with chromatic intervals instead of diatonic ones.
	"maj":     {4, 7},
	"aug":     {4, 8},
	"dim":     {3, 6},
	"dim7":    {3, 6, 9},
	"m":       {3, 7},
	"m7":      {3, 10, 7},
	"maj7":    {4, 11, 7},
	"7":       {4, 10, 7},
	"9":       {4, 10, 2, 7},
	"m9":      {3, 10, 2, 7},
	"maj9":    {4, 11, 2, 7},
	"13":      {4, 10, 2, 9, 7},
	"m11":     {3, 10, 14, -7, 7}, //m11 chords are strange, and should be spaced properly. Maybe they shouldn't even have a 9.
	"6":       {4, 7, 9},
	"m6":      {3, 7, 9},
	"69":      {4, 7, 9, 2},
	"m69":     {3, 7, 9, 2},
	"mmaj7":   {3, 11, 7},
	"minmaj7": {3, 11, 7},
}

var ChalkExtensionFormulas = map[string]int{
	"#1":   1,
	"#15":  1,
	"b9":   1,
	"9":    2,
	"#9":   3,
	"ll":   5,
	"#11":  6,
	"b5":   6,
	"#5":   8,
	"b13":  8,
	"13":   9,
	"sus":  5 - 12, // By having "sus" as an extension with a negative value, chords like C9sus will be properly parsed
	"sus2": 2 - 12, //but still put the fourth of the chord lower as is typical for this kind of chord
	"sus4": 5 - 12,
}

var ChordFormulas = map[string][]Interval{
	"maj":  {Root, MajorThird, PerfectFifth},
	"maj7": {Root, MajorThird, PerfectFifth, MajorSeventh},
	"m":    {Root, MinorThird, PerfectFifth},
	"dim":  {Root, MinorThird, DiminishedFifth},
	"aug":  {Root, MajorThird, AugmentedFifth},
	"sus2": {Root, MajorSecond, PerfectFifth},
	"sus4": {Root, PerfectFourth, PerfectFifth},
}

var ExtensionIntervals = map[string][]Interval{
	"6":    {MajorSixth},
	"7":    {MinorSeventh},
	"maj7": {MajorSeventh},
	"9":    {MinorSeventh, MajorNinth},
	"11":   {MinorSeventh, MajorNinth, PerfectEleventh},
	"13":   {MinorSeventh, MajorNinth, MajorThirteenth},
}

var AlterationMap = map[string]Interval{
	"5":  PerfectFifth,
	"6":  MajorSixth,
	"7":  MinorSeventh,
	"9":  MajorNinth,
	"11": PerfectEleventh,
	"13": MajorThirteenth,
}

var AdditionalNotes = map[string]Interval{
	"addb9":  MinorNinth,
	"add9":   MajorNinth,
	"add#9":  AugmentedNinth,
	"addb11": MinorEleventh,
	"add11":  PerfectEleventh,
	"add#11": AugmentedEleven,
	"addb13": MinorThirteenth,
	"add13":  MajorThirteenth,
	"add#13": AugmentedThirteenth,
}

func (c *Chord) addTone(tone Interval) {
	c.Tones = append(c.Tones, tone)
}

func (c *Chord) addTones(tones []Interval) {
	c.Tones = append(c.Tones, tones...)
}

func (c *Chord) replaceTone(tone Interval, newTone Interval) {
	for i, t := range c.Tones {
		if t == tone {
			c.Tones[i] = newTone
		}
	}
}

func (c *Chord) removeTone(tone Interval) {
	c.Tones = lo.Filter(c.Tones, func(i Interval, _ int) bool {
		return i != tone
	})
}

func ChalkRemove(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func (c Chord) hasInterval(interval Interval) bool {
	formula := ChordFormulas[c.Type]
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

func (c *Chord) applyVoicingRules() {
	for _, rule := range c.VoicingRules {
		if rule.Condition(c) {
			rule.Action(c)
		}
	}
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

func ParseChord(chordStr string) *Chord {
	regex := regexp.MustCompile(`([A-G][#b]?)((?:maj7?|m|dim|aug|sus2|sus4)?)((?:6|7|maj7|9|11|13)?)((?:[#b]\d{1,2})*)((?:add[#b]?\d{1,2})?)((/[A-G][#b]?)?)`)

	matches := regex.FindStringSubmatch(chordStr)

	chord := &Chord{
		Tones:        []Interval{},
		VoicingRules: defaultVoicingRules,
	}

	root := matches[1]
	if strings.Contains(root, "b") {
		root = flatToSharp[root]
	}
	chord.Root = root

	chordType := matches[2]
	if chordType == "" {
		chordType = "maj"
	}
	chord.Type = chordType
	chord.addTones(ChordFormulas[chordType])

	extension := matches[3]
	if extension != "" {
		chord.Extension = extension
		chord.addTones(ExtensionIntervals[extension])
	}

	addition := matches[5]
	chord.Addition = addition
	if addition != "" {
		chord.addTone(AdditionalNotes[addition])
	}

	alterationMatch := matches[4]
	alterations := regexp.MustCompile(`[#b]\d{1,2}`).FindAllString(alterationMatch, -1)
	chord.Alterations = alterations
	for _, alteration := range alterations {
		mod, degree := alteration[0], alteration[1:]
		var newInterval Interval
		switch mod {
		case '#':
			newInterval = AlterationMap[degree] + 1
		case 'b':
			newInterval = AlterationMap[degree] - 1
		}

		chord.replaceTone(AlterationMap[degree], newInterval)
	}

	bassNote := ""
	if matches[6] != "" {
		bassNote = matches[6][1:]
		if strings.Contains(bassNote, "b") {
			bassNote = flatToSharp[bassNote]
		}
	}
	chord.BassNote = bassNote

	chord.applyVoicingRules()

	return chord
}

func ChalkParseChord(inChord string) []float64 {
	finding := "root"
	chord := ChalkChord{}
	accumulator := ""
	extensions := []string{}
	for _, char := range inChord {

		if finding == "bass" {
			accumulator += string(char) //string concatenation
		}
		if char == '/' { //any part of the chord can suddenly end like this then go to bass
			switch finding {
			case "root":
				chord.Root = noteToNumber(accumulator)
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
			} else if char == 's' { //sus
				finding = "extensions"
				chord.Root = noteToNumber(accumulator)
				accumulator = ""
			} else {
				finding = "quality"
				chord.Root = noteToNumber(accumulator)
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
		chord.Root = noteToNumber(accumulator)
	case "quality":
		chord.Quality = accumulator
	case "extensions":
		if accumulator != "" {
			extensions = append(extensions, accumulator)
		}
		chord.Extensions = extensions
	case "bass": //rare "bass case"
		chord.BassNote = noteToNumber(accumulator)
	}
	if chord.BassNote == 0 {
		chord.BassNote = chord.Root
	}
	//
	// This concludes the section of the code that parses the chord input into its
	// root, quality, bass note, and extensions. Below, this information is translated
	// into a list of integers corresponding to the MIDI notes for the chord.
	//

	spreadExtensions := false
	lowRoot := true
	// Above is a list of parameters for the chord voicing that change how voicings are done
	extensionOctave := 48
	if spreadExtensions {
		extensionOctave = 60
	}

	chordTones := ChalkChordFormulas[chord.Quality]

	if slices.Contains(chord.Extensions, "sus") || slices.Contains(chord.Extensions, "sus2") || slices.Contains(chord.Extensions, "sus4") {
		for index, tone := range chordTones {
			if tone == 3 || tone == 4 {
				chordTones = ChalkRemove(chordTones, index)
			}
		}
	}
	if slices.Contains(chord.Extensions, "b5") || slices.Contains(chord.Extensions, "#5") {
		for index, tone := range chordTones {
			if tone == 7 {
				chordTones = ChalkRemove(chordTones, index)
			}
		}
	}

	// The following code generates a voicing with no regard to what the previous chord was.
	// In the python code, a much longer chunk of code uses an iterative algorithm to find
	// a combination of octaves for the notes that minimizes the sum of the distances
	// of each note in the new chord to the nearest note in the previous chord.

	voicing := []int{}

	for _, tone := range chordTones {
		voicing = append(voicing, chord.Root+48+tone)
	}
	for _, extension := range extensions {
		voicing = append(voicing, chord.Root+extensionOctave+ChalkExtensionFormulas[extension])
	}
	if lowRoot && !slices.Contains(chordTones, 6) && !slices.Contains(chordTones, 8) && !slices.Contains(extensions, "#5") && !slices.Contains(extensions, "b5") {
		voicing = append(voicing, chord.Root+43)
	}
	voicing = append(voicing, chord.Root+48)
	if lowRoot {
		voicing = append(voicing, chord.BassNote+36)
		voicing = append(voicing, chord.BassNote+24)
	}

	// The chord is now voiced in MIDI note numbers.
	freqVoicing := []float64{}

	for _, m := range voicing {
		freqVoicing = append(freqVoicing, math.Pow(2, float64((float64(m)-69.0)/12.0))*440.0)
	}
	return freqVoicing
}
