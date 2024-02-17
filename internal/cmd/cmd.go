package cmd

import (
	"github.com/kimond/godio/pkg/godio"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "godio",
	Short: "Sound generation library",
	Long:  `godio is a library for generating sound. It is a work in progress.`,
}

func Execute() {
	rootCmd.AddCommand(noteCmd)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	noteCmd.Flags().StringP("waveform", "w", string(godio.WaveformSine), "Waveform to use (Sine, Square, Sawtooth, Triangle)")
	noteCmd.Flags().StringP("output", "o", "note.wav", "Output file name")
}

var noteCmd = &cobra.Command{
	Use:        "note [frequency] [duration]",
	Short:      "Generate a note",
	Long:       `Generate a note with a given frequency and duration in seconds.`,
	Args:       cobra.ExactArgs(2),
	ArgAliases: []string{"frequency", "duration"},
	Run: func(cmd *cobra.Command, args []string) {
		frequency := args[0]
		duration, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			log.Fatalf("error parsing duration: %v", err)
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
