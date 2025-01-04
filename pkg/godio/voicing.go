package godio

type VoicingRule struct {
	Condition func(*Chord) bool
	Action    func(*Chord)
}

// Common voicing rules based on jazz theory
var defaultVoicingRules = []VoicingRule{
	{
		// In ninth chords, fifth is often omitted
		Condition: func(c *Chord) bool {
			return c.hasInterval(MajorNinth) || c.hasInterval(MinorNinth)
		},
		Action: func(c *Chord) {
			c.removeTone(PerfectFifth)
		},
	},
	{
		// In thirteenth chords, fifth and ninth are often omitted
		Condition: func(c *Chord) bool {
			return c.hasInterval(MajorThirteenth) || c.hasInterval(MinorThirteenth)
		},
		Action: func(c *Chord) {
			c.removeTone(PerfectFifth)
			c.removeTone(MajorNinth)
		},
	},
	{
		// In altered dominant chords, natural fifth is omitted
		Condition: func(c *Chord) bool {
			return (c.hasInterval(DiminishedFifth) || c.hasInterval(AugmentedFifth)) && !c.hasInterval(MajorSeventh)
		},
		Action: func(c *Chord) {
			c.removeTone(PerfectFifth)
		},
	},
	{
		// In eleventh chords, third is often omitted (except in maj11)
		Condition: func(c *Chord) bool {
			return c.hasInterval(PerfectEleventh) && !c.hasInterval(MajorSeventh)
		},
		Action: func(c *Chord) {
			c.removeTone(MajorThird)
			c.removeTone(MinorThird)
		},
	},
	{
		// Remove upper root if there are 5 or more tones
		Condition: func(c *Chord) bool {
			return len(c.Tones) >= 5
		},
		Action: func(c *Chord) {
			c.removeTone(Root)
		},
	},
	{
		// Remove the fifth if chord has more than 5 notes and no alteration of the fifth
		Condition: func(c *Chord) bool {
			return len(c.Tones) >= 5 && !c.hasInterval(AugmentedFifth) && !c.hasInterval(DiminishedFifth)
		},
		Action: func(c *Chord) {
			c.removeTone(PerfectFifth)
		},
	},
}
