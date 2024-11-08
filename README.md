[![Test](https://github.com/hraban/opus/workflows/Test/badge.svg)](https://github.com/hraban/opus/actions?query=workflow%3ATest)

## Differences

This fork adds the capability of encoding ogg streams.

## Go wrapper for Opus

This package provides Go bindings for the xiph.org C libraries libopus and
libopusfile.

The C libraries and docs are hosted at https://opus-codec.org/. This package
just handles the wrapping in Go, and is unaffiliated with xiph.org.

Features:

- ✅ encode and decode raw PCM data to raw Opus data
- ✅ useful when you control the recording device, _and_ the playback
- ✅ decode .opus and .ogg files into raw audio data ("PCM")
- ✅ reuse the system libraries for opus decoding (libopus)
- ✅ works easily on Linux, Mac and Docker; needs libs on Windows
- ✅ does  _create_ .opus or .ogg files and streams
- ❌ does not work with .wav files (you need a separate .wav library for that)
- ❌ no self-contained binary (you need the xiph.org libopus lib, e.g. through a package manager)
- ❌ no cross compiling (because it uses CGo)

Good use cases:

- 👍 you are writing a music player app in Go, and you want to play back .opus files
- 👍 you record raw wav in a web app or mobile app, you encode it as Opus on the client, you send the opus to a remote webserver written in Go, and you want to decode it back to raw audio data on that server

## Details

This wrapper provides a Go translation layer for three elements from the
xiph.org opus libs:

* encoders
* decoders
* files & streams

### Import

```go
import "github.com/vdepedraza/opus-go"
```

### Encoding

To encode raw audio to the Opus format, create an encoder first:

```go
const sampleRate = 48000
const channels = 1 // mono; 2 for stereo

enc, err := opus.NewEncoder(sampleRate, channels, opus.AppVoIP)
if err != nil {
    ...
}
```

Then pass it some raw PCM data to encode.

Make sure that the raw PCM data you want to encode has a legal Opus frame size.
This means it must be exactly 2.5, 5, 10, 20, 40 or 60 ms long. The number of
bytes this corresponds to depends on the sample rate (see the [libopus
documentation](https://www.opus-codec.org/docs/opus_api-1.1.3/group__opus__encoder.html)).

```go
var pcm []int16 = ... // obtain your raw PCM data somewhere
const bufferSize = 1000 // choose any buffer size you like. 1k is plenty.

// Check the frame size. You don't need to do this if you trust your input.
frameSize := len(pcm) // must be interleaved if stereo
frameSizeMs := float32(frameSize) / channels * 1000 / sampleRate
switch frameSizeMs {
case 2.5, 5, 10, 20, 40, 60:
    // Good.
default:
    return fmt.Errorf("Illegal frame size: %d bytes (%f ms)", frameSize, frameSizeMs)
}

data := make([]byte, bufferSize)
n, err := enc.Encode(pcm, data)
if err != nil {
    ...
}
data = data[:n] // only the first N bytes are opus data. Just like io.Reader.
```

Note that you must choose a target buffer size, and this buffer size will affect
the encoding process:

> Size of the allocated memory for the output payload. This may be used to
> impose an upper limit on the instant bitrate, but should not be used as the
> only bitrate control. Use `OPUS_SET_BITRATE` to control the bitrate.

-- https://opus-codec.org/docs/opus_api-1.1.3/group__opus__encoder.html

### Decoding

To decode opus data to raw PCM format, first create a decoder:

```go
dec, err := opus.NewDecoder(sampleRate, channels)
if err != nil {
    ...
}
```

Now pass it the opus bytes, and a buffer to store the PCM sound in:

```go
var frameSizeMs float32 = ...  // if you don't know, go with 60 ms.
frameSize := channels * frameSizeMs * sampleRate / 1000
pcm := make([]int16, int(frameSize))
n, err := dec.Decode(data, pcm)
if err != nil {
    ...
}

// To get all samples (interleaved if multiple channels):
pcm = pcm[:n*channels] // only necessary if you didn't know the right frame size

// or access sample per sample, directly:
for i := 0; i < n; i++ {
    ch1 := pcm[i*channels+0]
    // For stereo output: copy ch1 into ch2 in mono mode, or deinterleave stereo
    ch2 := pcm[(i*channels)+(channels-1)]
}
```

To handle packet loss from an unreliable network, see the
[DecodePLC](https://godoc.org/gopkg.in/hraban/opus.v2#Decoder.DecodePLC) and
[DecodeFEC](https://godoc.org/gopkg.in/hraban/opus.v2#Decoder.DecodeFEC)
options.

### Streams (and Files)

To decode a .opus file (or .ogg with Opus data), or to decode a "Opus stream"
(which is a Ogg stream with Opus data), use the `Stream` interface. It wraps an
io.Reader providing the raw stream bytes and returns the decoded Opus data.

A crude example for reading from a .opus file:

```go
f, err := os.Open(fname)
if err != nil {
    ...
}
s, err := opus.NewStream(f)
if err != nil {
    ...
}
defer s.Close()
pcmbuf := make([]int16, 16384)
for {
    n, err = s.Read(pcmbuf)
    if err == io.EOF {
        break
    } else if err != nil {
        ...
    }
    pcm := pcmbuf[:n*channels]

    // send pcm to audio device here, or write to a .wav file

}
```

To encode a .opus file, or a Stream, see the example on examples/e0.


### API Docs

## Build & Installation

This package requires libopus ,libopusfile and libopusenc development packages to be
installed on your system. These are available on Debian based systems from
aptitude as `libopus-dev` , `libopusfile-dev` and `libopusenc-dev`, and on Mac OS X from homebrew.

They are linked into the app using pkg-config.

Debian, Ubuntu, ...:
```sh
sudo apt-get install pkg-config libopus-dev libopusfile-dev libopusenc-dev
```

## License

The licensing terms for the Go bindings are found in the LICENSE file. The
authors and copyright holders are listed in the AUTHORS file.

The copyright notice uses range notation to indicate all years in between are
subject to copyright, as well. This statement is necessary, apparently. For all
those nefarious actors ready to abuse a copyright notice with incorrect
notation, but thwarted by a mention in the README. Pfew!
