package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/hajimehoshi/oto"
	"github.com/vdepedraza/opus-go"
)

func int16ToByteBuffer(int16Buf []int16) []byte {
	byteBuf := make([]byte, len(int16Buf)*2) // Each int16 is 2 bytes
	for i, v := range int16Buf {
		// Convert int16 to two bytes and store them in the byte buffer
		binary.LittleEndian.PutUint16(byteBuf[i*2:], uint16(v))
	}
	return byteBuf
}

func main() {
	const (
		sampleRate     = 48000         // Samples per second
		channelCount   = 1             // Number of audio channels (e.g., 2 for stereo)
		bytesPerSample = 2             // Bytes per sample (16-bit audio)
		bufferSize     = 48000 * 1 * 2 // 1 second buffer for stereo 16-bit audio
	)

	context, err := oto.NewContext(sampleRate, channelCount, bytesPerSample, bufferSize)
	if err != nil {
		panic(err)
	}
	defer context.Close()

	player := context.NewPlayer()
	defer player.Close()

	// opus decoder begins

	f, err := os.Open("sample1_mono.opus")
	if err != nil {
		panic(err)
	}
	s, err := opus.NewStream(f)
	if err != nil {
		panic(err)
	}
	defer s.Close()
	pcmbuf := make([]int16, 16384)

	// encoder

	fo, err := os.Create("output.opus")
	if err != nil {
		panic(err)
	}
	so, err := opus.NewEncoderStream(fo)
	so.SetBitrate(6000)
	so.SetApplication(opus.AppVoIP)
	so.SetComplexity(10)
	if err != nil {
		panic(err)
	}
	defer so.Close()

	for {
		n, err := s.Read(pcmbuf)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		pcm := pcmbuf[:n*1] // n * number of channels

		// send pcm to audio device here, or write to a .wav file
		//player.Write(int16ToByteBuffer(pcm))

		// write reencoded opus
		so.Write(pcm)
	}

	fmt.Printf("Encoding\n")
}
