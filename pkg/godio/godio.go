package godio

import (
	"fmt"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"io"
	"math"
)

const (
	sampleRate = 44100
	volume     = 0.2
)

type Waveform string

const (
	WaveformSine     Waveform = "Sine"
	WaveformSquare   Waveform = "Square"
	WaveformSawtooth Waveform = "Sawtooth"
	WaveformTriangle Waveform = "Triangle"
)

// ADSREnvelope defines the structure for our ADSR envelope
type ADSREnvelope struct {
	Attack  int     // Duration of the attack phase in milliseconds
	Decay   int     // Duration of the decay phase in milliseconds
	Sustain float64 // Sustain level (0 to 1)
	Release int     // Duration of the release phase in milliseconds
}

type SoundBuffer struct {
	buffers [][]int
}

func NewSoundBuffer() *SoundBuffer {
	return &SoundBuffer{}
}

// ApplyADSR applies the ADSR envelope to a buffer
func (sb *SoundBuffer) ApplyADSR(env ADSREnvelope) {
	attackLength := (env.Attack * sampleRate) / 1000
	decayLength := (env.Decay * sampleRate) / 1000
	releaseLength := (env.Release * sampleRate) / 1000

	for _, buffer := range sb.buffers {
		totalLength := len(buffer)
		sustainLength := totalLength - attackLength - decayLength - releaseLength
		for i := range buffer {
			var amplitude float64
			switch {
			case i < attackLength:
				amplitude = float64(i) / float64(attackLength)
			case i < attackLength+decayLength:
				amplitude = 1 - (1-env.Sustain)*float64(i-attackLength)/float64(decayLength)
			case i < attackLength+decayLength+sustainLength:
				amplitude = env.Sustain
			case i < totalLength:
				amplitude = env.Sustain * (1 - float64(i-attackLength-decayLength-sustainLength)/float64(releaseLength))
			}
			buffer[i] = int(float64(buffer[i]) * amplitude)
		}
	}
}

func (sb *SoundBuffer) Write(seeker io.WriteSeeker) error {
	intBuf := &audio.IntBuffer{Data: sb.combineBuffers(), Format: &audio.Format{SampleRate: sampleRate, NumChannels: 1}}
	encoder := wav.NewEncoder(seeker, sampleRate, 16, 1, 1)
	if err := encoder.Write(intBuf); err != nil {
		return fmt.Errorf("error writing buffer to wav: %v", err)
	}
	if err := encoder.Close(); err != nil {
		return fmt.Errorf("error closing encoder: %v", err)
	}
	_, err := (seeker).Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("error seeking to start of buffer: %v", err)
	}

	return nil
}

// AppendNote appends a note to a SoundBuffer
func (sb *SoundBuffer) AppendNote(frequency float64, durationSec float64, waveform Waveform) {
	numSamples := int(float64(sampleRate) * durationSec)
	buf := make([]int, numSamples)

	for i := 0; i < numSamples; i++ {
		var sample float64
		t := float64(i) / float64(sampleRate)

		switch waveform {
		case WaveformSine:
			sample = math.Sin(2 * math.Pi * frequency * t)
		case WaveformSquare:
			if int(math.Floor(t*frequency))%2 == 0 {
				sample = 1.0
			} else {
				sample = -1.0
			}
		case WaveformSawtooth:
			sample = 2.0 * (t*frequency - math.Floor(t*frequency+0.5))
		case WaveformTriangle:
			sample = 2.0*math.Abs(2.0*(t*frequency-math.Floor(t*frequency+0.5))) - 1.0
		}

		buf[i] = int(volume * 32767 * sample)
	}
	sb.buffers = append(sb.buffers, buf)
}

// AppendChord append a chord buffer for a given set of frequencies and waveform type.
func (sb *SoundBuffer) AppendChord(frequencies []float64, durationSec float64, waveform Waveform) {
	numSamples := int(float64(sampleRate) * durationSec)
	// Create a buffer for each note in the chord
	buffers := make([][]float64, len(frequencies))

	// Generate the wave for each frequency and store in buffers
	for i, freq := range frequencies {
		buffers[i] = make([]float64, numSamples)
		for j := 0; j < numSamples; j++ {
			t := float64(j) / float64(sampleRate)
			// Choose waveform type
			switch waveform {
			case WaveformSine:
				buffers[i][j] = math.Sin(2 * math.Pi * freq * t)
			case WaveformSquare:
				if int(math.Floor(t*freq))%2 == 0 {
					buffers[i][j] = 1.0
				} else {
					buffers[i][j] = -1.0
				}
			case WaveformSawtooth:
				buffers[i][j] = 2.0 * (t*freq - math.Floor(t*freq+0.5))
			case WaveformTriangle:
				buffers[i][j] = 2.0*math.Abs(2.0*(t*freq-math.Floor(t*freq+0.5))) - 1.0
			}
		}
	}

	// Mix the buffers together
	chordBuffer := make([]int, numSamples)
	for i := 0; i < numSamples; i++ {
		var sample float64
		for _, buffer := range buffers {
			sample += buffer[i]
		}
		// Normalize the sample to prevent clipping
		sample = sample / float64(len(buffers))

		chordBuffer[i] = int(volume * 32767 * sample)
	}

	sb.buffers = append(sb.buffers, chordBuffer)
}

func (sb *SoundBuffer) combineBuffers() []int {
	var combinedBuffer []int
	for _, buffer := range sb.buffers {
		if len(combinedBuffer) == 0 {
			combinedBuffer = buffer
		} else {
			combinedBuffer = append(combinedBuffer, buffer...)
		}
	}

	return combinedBuffer
}
