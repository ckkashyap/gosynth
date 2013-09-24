package main

import (
	"bytes"
	"encoding/binary"
	"github.com/ckkashyap/wav"
	"os"
	"math"
	"log"
)


type Track chan int16


func GetTrack(notes []Note, samplingFrequency int32) Track {
	t := make(chan int16)
	go func() {
		for _, n := range notes {
			log.Printf("%v", n)
			numberOfSamples := int32(samplingFrequency * int32(n.Duration) / 100.0)
			step := float64((360.0 / float64(samplingFrequency)) * n.Freq)
			angle := float64(0)
			for i := int32(0); i < numberOfSamples; i++ {
				rad := math.Pi * angle / 180.0
				sample := int16(32767 * math.Sin(rad))
				t <- sample
				angle += step
				if angle > 360.0 {
					angle -= 360.0
				}
			}

		}
		close(t)
	}()
	return t
}


var song = []string{
	"C-3", "C#3", "D-3", "D#3", "E-3", "F-3",
	"F#3", "G-3", "G#3", "A-3", "A#3", "B-3",
	"C-4", "C#4", "D-4", "D#4", "E-4", "F-4",
	"F#4", "G-4", "G#4", "A-4", "A#4", "B-4",
	"C-5", "C#5", "D-5", "D#5", "E-5", "F-5",
	"F#5", "G-5", "G#5", "A-5", "A#5", "B-5",
	"C-6",
}

const (
	noteDuration = waveHz / 0.1
	songLength   = waveHz * 10
)

func main() {
	// build song
	var q []Note
	for _, n := range song {
		q = append(q, Note{NoteToFreq(n), 1.0, 20}, Note{0, 0, 10})
	}

//	osc := NewModule(&Oscillator{}, &NoteLane{Q: q})
//	env := NewModule(&AmpEnvelope{200, 1000}, &NoteLane{Q: q})
//
//	in := make(chan Buf)
//	out := env.Run(osc.Run(in))
//
//	in <- make(Buf, 1000)
//
//	// play song, writing to buf
	var buf bytes.Buffer
//	played := 0
	for b := range GetTrack(q, waveHz) {
		binary.Write(&buf, binary.LittleEndian, b)
	}

	f, err := os.Create("test.wav")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := &wav.File{
		SampleRate:      waveHz,
		SignificantBits: 16,
		Channels:        2,
	}
//	b := make([]byte, 44100*4*10)
	if err := w.WriteData(f, buf.Bytes()); err != nil {
		panic(err)
	}
	log.Println("done")
}
