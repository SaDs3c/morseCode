package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

var morseCodeDict = map[rune]string{
	'A': ".-", 'B': "-...", 'C': "-.-.", 'D': "-..", 'E': ".", 'F': "..-.", 'G': "--.", 'H': "....",
	'I': "..", 'J': ".---", 'K': "-.-", 'L': ".-..", 'M': "--", 'N': "-.", 'O': "---", 'P': ".--.",
	'Q': "--.-", 'R': ".-.", 'S': "...", 'T': "-", 'U': "..-", 'V': "...-", 'W': ".--", 'X': "-..-",
	'Y': "-.--", 'Z': "--..",
	'1': ".----", '2': "..---", '3': "...--", '4': "....-", '5': ".....", '6': "-....", '7': "--...",
	'8': "---..", '9': "----.", '0': "-----",
}

func textToMorse(text string) string {
	var morse []string
	for _, char := range strings.ToUpper(text) {
		if code, exists := morseCodeDict[char]; exists {
			morse = append(morse, code)
		}
	}
	return strings.Join(morse, " ")
}

func generateTone(frequency, duration, sampleRate float64, amplitude float64) []float64 {
	var tone []float64
	sampleCount := int(sampleRate * duration)
	for i := 0; i < sampleCount; i++ {
		t := float64(i) / sampleRate
		sample := amplitude * math.Sin(2*math.Pi*frequency*t)
		tone = append(tone, sample)
	}
	return tone
}

func appendSilence(duration, sampleRate float64) []float64 {
	return make([]float64, int(sampleRate*duration))
}

func textToMorseAudio(text, filename string) {
	morseCode := textToMorse(text)

	dotDuration := 0.1        // Duration of a dot in seconds
	dashDuration := 3 * dotDuration // Duration of a dash
	frequency := 1000.0       // Frequency of the tone
	sampleRate := 44100.0     // Samples per second
	amplitude := 0.5          // Amplitude of the tone

	var audioSequence []float64

	for _, char := range morseCode {
		switch char {
		case '.':
			audioSequence = append(audioSequence, generateTone(frequency, dotDuration, sampleRate, amplitude)...)
		case '-':
			audioSequence = append(audioSequence, generateTone(frequency, dashDuration, sampleRate, amplitude)...)
		case ' ':
			audioSequence = append(audioSequence, appendSilence(dotDuration, sampleRate)...)
		}
		audioSequence = append(audioSequence, appendSilence(dotDuration, sampleRate)...)
	}

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	enc := wav.NewEncoder(f, int(sampleRate), 16, 1, 1)
	intBuf := &audio.IntBuffer{Data: make([]int, len(audioSequence)), Format: &audio.Format{SampleRate: int(sampleRate), NumChannels: 1}}

	for i, v := range audioSequence {
		intBuf.Data[i] = int(v * (1 << 15-1))
	}

	if err := enc.Write(intBuf); err != nil {
		panic(err)
	}
	if err := enc.Close(); err != nil {
		panic(err)
	}

	fmt.Printf("Saved Morse code audio to %s\n", filename)
}

func main() {
	textInput := "WEB3"
	filename := strings.ReplaceAll(textInput, " ", "_") + ".wav"
	textToMorseAudio(textInput, filename)
}
