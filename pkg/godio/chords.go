package godio

import (
	"regexp"
	"strings"

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
