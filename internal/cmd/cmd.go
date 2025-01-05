package cmd

import (
	"github.com/kimond/godio/pkg/godio"
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "godio",
	Short: "Sound generation library",
	Long:  `godio is a library for generating sound. It is a work in progress.`,
}

func Execute() {
	rootCmd.AddCommand(noteCmd)
	rootCmd.AddCommand(chordCmd)
	rootCmd.AddCommand(sequenceCmd)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	addCommonFlags(noteCmd)
	addCommonFlags(chordCmd)
	addCommonFlags(sequenceCmd)
}

func addCommonFlags(cmd *cobra.Command) {
	cmd.Flags().Float64P("duration", "d", 1, "Duration in seconds")
	cmd.Flags().StringP("waveform", "w", string(godio.WaveformTriangle), "Waveform to use (Sine, Square, Sawtooth, Triangle)")
	cmd.Flags().StringP("output", "o", "note.wav", "Output file name")
}

var noteCmd = &cobra.Command{
	Use:        "note [frequency]",
	Short:      "Generate a note",
	Long:       `Generate a note with a given frequency.`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"frequency"},
	Run: func(cmd *cobra.Command, args []string) {
		frequency := args[0]
		duration, err := cmd.Flags().GetFloat64("duration")
		if err != nil {
			panic(err)
		}
		waveform, err := cmd.Flags().GetString("waveform")
		if err != nil {
			panic(err)
		}
		output, err := cmd.Flags().GetString("output")
		if err != nil {
			panic(err)
		}

		sb := godio.NewSoundBuffer()
		sb.AppendNote(godio.NoteFrequencies[frequency], duration, godio.Waveform(waveform))

		wavFile, err := os.Create(output)
		if err != nil {
			panic(err)
		}

		if err := sb.Write(wavFile); err != nil {
			panic(err)
		}
	},
}

var chordCmd = &cobra.Command{
	Use:        "chord [chord]",
	Short:      "Generate a chord",
	Long:       `Generate a chord.`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"chord"},
	Run: func(cmd *cobra.Command, args []string) {
		chordString := args[0]

		duration, err := cmd.Flags().GetFloat64("duration")
		if err != nil {
			panic(err)
		}
		waveform, err := cmd.Flags().GetString("waveform")
		if err != nil {
			panic(err)
		}
		output, err := cmd.Flags().GetString("output")
		if err != nil {
			panic(err)
		}

		chord := godio.ChalkParseChord(chordString)
		sb := godio.NewSoundBuffer()
		sb.AppendChord(chord, duration, godio.Waveform(waveform))
		sb.ApplyADSR(godio.ADSREnvelope{
			Attack:  10,
			Decay:   0,
			Sustain: 1,
			Release: 100,
		})

		wavFile, err := os.Create(output)
		if err != nil {
			panic(err)
		}

		if err := sb.Write(wavFile); err != nil {
			panic(err)
		}
	},
}

var sequenceCmd = &cobra.Command{
	Use:   "sequence",
	Short: "Generate a sequence of chords",
	Long:  `Generate a sequence of chords.`,
	Run: func(cmd *cobra.Command, args []string) {
		duration, err := cmd.Flags().GetFloat64("duration")
		if err != nil {
			panic(err)
		}
		waveform, err := cmd.Flags().GetString("waveform")
		if err != nil {
			panic(err)
		}
		output, err := cmd.Flags().GetString("output")
		if err != nil {
			panic(err)
		}
		chords := args

		sb := godio.NewSoundBuffer()
		for _, chordStr := range chords {
			chord := godio.ParseChord(chordStr)
			sb.AppendChord(chord.GetFrequencies(), duration, godio.Waveform(waveform))
		}
		sb.ApplyADSR(godio.ADSREnvelope{
			Attack:  1,
			Decay:   int(duration * 1000),
			Sustain: 0,
			Release: 0,
		})

		wavFile, err := os.Create(output)
		if err != nil {
			panic(err)
		}

		if err := sb.Write(wavFile); err != nil {
			panic(err)
		}
	},
}
