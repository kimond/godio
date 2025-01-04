package godio

import (
	"fmt"
	"testing"
)

func TestChord(t *testing.T) {
	parameters := []struct {
		input    string
		expected []string
	}{
		{"Cmaj", []string{"C2", "C4", "E4", "G3"}},
		{"Cmaj7", []string{"C2", "C4", "E4", "G3", "B3"}},
		{"Cm7", []string{"C2", "C4", "D#4", "G3", "A#3"}},
		{"Cmmaj7", []string{"C2", "C4", "D#4", "G3", "B3"}},
		{"Cm7b5", []string{"C2", "C4", "D#4", "F#4", "A#3"}},
		{"Cmadd11", []string{"C2", "C4", "D#4", "G3", "F4"}},
		{"Cmadd#11", []string{"C2", "C4", "D#4", "G3", "F#4"}},
		{"C11", []string{"C2", "E4", "A#3", "D4", "F4"}},
		{"C11b5", []string{"C2", "E4", "F#4", "A#3", "D4", "F4"}},
		{"Dbm7", []string{"C#2", "C#4", "E4", "G#3", "B3"}},
		{"G9", []string{"G2", "B3", "D4", "F4", "A3"}},
		{"F#13", []string{"F#2", "A#3", "E4", "G#3", "D#4"}},
		{"Cm7/G", []string{"G2", "C4", "D#4", "G3", "A#3"}},
		{"Cm7/Gb", []string{"F#2", "C4", "D#4", "G3", "A#3"}},
	}

	for i := range parameters {
		t.Run(fmt.Sprintf("Testing %v", parameters[i].input), func(t *testing.T) {
			chord := ParseChord(parameters[i].input)
			frequencies := chord.GetFrequencies()
			for j, frequency := range frequencies {
				expectedFrequency := NoteFrequencies[parameters[i].expected[j]]
				if frequency != expectedFrequency {
					t.Errorf("Expected %f, but got %f in position %d", expectedFrequency, frequency, j+1)
				}
			}
		})
	}
}
