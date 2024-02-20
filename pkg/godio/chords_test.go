package godio

import (
	"fmt"
	"testing"
)

func TestChord(t *testing.T) {
	parameters := []struct {
		input    string
		expected []float64
	}{
		{"Cmaj", []float64{NoteFrequencies["C2"], NoteFrequencies["E4"], NoteFrequencies["G4"]}},
		{"Cm7", []float64{NoteFrequencies["C2"], NoteFrequencies["D#4"], NoteFrequencies["G4"], NoteFrequencies["A#4"]}},
		{"Cm7b5", []float64{NoteFrequencies["C2"], NoteFrequencies["D#4"], NoteFrequencies["F#4"], NoteFrequencies["A#4"]}},
		{"Cmadd11", []float64{NoteFrequencies["C2"], NoteFrequencies["D#4"], NoteFrequencies["G4"], NoteFrequencies["F5"]}},
		{"C11", []float64{NoteFrequencies["C2"], NoteFrequencies["E4"], NoteFrequencies["A#4"], NoteFrequencies["D5"], NoteFrequencies["F5"]}},
		{"C11b5", []float64{NoteFrequencies["C2"], NoteFrequencies["E4"], NoteFrequencies["F#4"], NoteFrequencies["A#4"], NoteFrequencies["D5"], NoteFrequencies["F5"]}},
	}

	for i := range parameters {
		t.Run(fmt.Sprintf("Testing %v", parameters[i].input), func(t *testing.T) {
			chord := ParseChord(parameters[i].input)
			frequencies := chord.GetFrequencies()
			for j, frequency := range frequencies {
				if frequency != parameters[i].expected[j] {
					t.Errorf("Expected %f, but got %f", parameters[i].expected[j], frequency)
				}
			}
		})
	}
}
