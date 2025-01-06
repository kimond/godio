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
			return c.hasInterval(2) || c.hasInterval(1) || c.hasInterval(14) || c.hasInterval(13)
		},
		Action: func(c *Chord) {
			c.removeTone(7)
		},
	},
	{
		// In thirteenth chords, fifth and ninth are often omitted
		Condition: func(c *Chord) bool {
			return c.hasInterval(9) || c.hasInterval(8) || c.hasInterval(21) || c.hasInterval(20)
		},
		Action: func(c *Chord) {
			c.removeTone(7)
			c.removeTone(14)
		},
	},
	{
		// In altered dominant chords, natural fifth is omitted
		Condition: func(c *Chord) bool {
			return (c.hasInterval(6) || c.hasInterval(8)) && !c.hasInterval(11)
		},
		Action: func(c *Chord) {
			c.removeTone(7)
		},
	},
	{
		// In eleventh chords, third is often omitted (except in maj11)
		Condition: func(c *Chord) bool {
			return c.hasInterval(17) || c.hasInterval(5) && !c.hasInterval(11)
		},
		Action: func(c *Chord) {
			c.removeTone(4)
			c.removeTone(3)
		},
	},
	{
		// Remove upper root if there are 5 or more tones
		Condition: func(c *Chord) bool {
			return len(c.Tones) >= 5
		},
		Action: func(c *Chord) {
			c.removeTone(0)
		},
	},
	{
		// Remove the fifth if chord has more than 5 notes and no alteration of the fifth
		Condition: func(c *Chord) bool {
			return len(c.Tones) >= 5 && !c.hasInterval(8) && !c.hasInterval(6)
		},
		Action: func(c *Chord) {
			c.removeTone(7)
		},
	},
}
