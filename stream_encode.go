// Copyright © Go Opus Authors (see AUTHORS file)
//
// License for use of this code is detailed in the LICENSE file

package opus

/*
#cgo pkg-config: libopusenc

#include <opusfile.h>
#include <opusenc.h>
#include <stdint.h>
#include <string.h>

OggOpusEnc *my_ope_encoder_create_callback(uintptr_t p, int *error);
int setBitrate( OggOpusEnc * enc,uint32_t bitrate);
int setApplication (OggOpusEnc * enc, uint32_t app);
int setComplexity (OggOpusEnc * enc, uint32_t complexity);

*/
import "C"
import (
	"fmt"
	"io"
	"unsafe"
)

type encoderStreamParameters struct {
	encoderApplication C.int
	encoderBitRate     C.int
	encoderComplexity  C.int
	encoderFrameSize   C.int
	encoderSampleRate  C.int
	maxFamesPerPage    C.int
	resampleQuality    C.int
}

type EncoderStream struct {
	id         uintptr
	oggencoder *C.OggOpusEnc
	write      io.Writer
	// Preallocated buffer to pass to the reader
	buf []byte
}

var streamsEnc = newEncoderStreamsMap()

//export go_writecallback
func go_writecallback(p unsafe.Pointer, cbuf *C.uchar, nbytes C.int) C.int {
	streamId := uintptr(p)
	stream := streamsEnc.Get(streamId)
	if stream == nil {
		return 1
	}

	maxbytes := int(nbytes)
	if maxbytes > cap(stream.buf) {
		maxbytes = cap(stream.buf)
	}

	C.memcpy(unsafe.Pointer(&stream.buf[0]), unsafe.Pointer(cbuf), C.size_t(nbytes))

	_, _ = stream.write.Write(stream.buf[:nbytes])
	//fmt.Printf("Written %v bytes\n", n)
	return 0
}

// NewStream creates and initializes a new stream. Don't call .Init() on this.
func NewEncoderStream(write io.Writer) (*EncoderStream, error) {
	var s EncoderStream
	err := s.Init(write)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// Init initializes a stream with an io.Writer to wich write PCM decoded data.
func (s *EncoderStream) Init(write io.Writer) error {
	if s.oggencoder != nil {
		return fmt.Errorf("opus encoder stream is already initialized")
	}
	if write == nil {
		return fmt.Errorf("Writer must be non-nil")
	}

	s.write = write
	s.buf = make([]byte, maxDecodedFrameSize)
	s.id = streamsEnc.NextId()
	var errno C.int

	// Immediately delete the stream after .Init to avoid leaking if the
	// caller forgets to (/ doesn't want to) call .Close(). No need for that,
	// since the callback is only ever called during a .Read operation; just
	// Save and Delete from the map around that every time a reader function is
	// called.
	streamsEnc.Save(s)
	defer streamsEnc.Del(s)
	oggencoder := C.my_ope_encoder_create_callback(C.uintptr_t(s.id), &errno)
	if errno != 0 {
		return StreamEncodeError(errno)
	}
	s.oggencoder = oggencoder
	return nil
}

func (s *EncoderStream) SetBitrate(bitrate int32) error {
	int_err := C.setBitrate(s.oggencoder, C.uint(bitrate))
	if int_err != 0 {
		return fmt.Errorf("setBitrate error: %v", int_err)
	}
	return nil
}

func (s *EncoderStream) SetApplication(app Application) error {
	int_err := C.setApplication(s.oggencoder, C.uint(app))
	if int_err != 0 {
		return fmt.Errorf("setApplication error: %v", int_err)
	}
	return nil
}

func (s *EncoderStream) SetComplexity(complexity int32) error {
	int_err := C.setComplexity(s.oggencoder, C.uint(complexity))
	if int_err != 0 {
		return fmt.Errorf("setComplexity error: %v", int_err)
	}
	return nil
}

func (s *EncoderStream) Write(pcm []int16) error {
	if s.oggencoder == nil {
		return fmt.Errorf("opus encoder stream is uninitialized or already closed")
	}
	if len(pcm) == 0 {
		return nil
	}
	streamsEnc.Save(s)
	defer streamsEnc.Del(s)

	n := C.ope_encoder_write(
		s.oggencoder,
		(*C.opus_int16)(&pcm[0]),
		C.int(len(pcm)))

	if n == 0 {
		return nil
	} else {
		return StreamEncodeError(n)
	}
}

func (s *EncoderStream) WriteFloat32(pcm []float32) error {
	if s.oggencoder == nil {
		return fmt.Errorf("opus encoder stream is uninitialized or already closed")
	}
	if len(pcm) == 0 {
		return nil
	}
	streamsEnc.Save(s)
	defer streamsEnc.Del(s)
	n := C.ope_encoder_write_float(
		s.oggencoder,
		(*C.float)(&pcm[0]),
		C.int(len(pcm)))

	if n == 0 {
		return nil
	} else {
		return StreamEncodeError(n)
	}
}

func (s *EncoderStream) Close() error {
	if s.oggencoder == nil {
		return fmt.Errorf("opus encoder stream is uninitialized or already closed")
	}
	streamsEnc.Save(s)
	defer streamsEnc.Del(s)
	int_err := C.ope_encoder_drain(s.oggencoder)

	if int_err != 0 {
		return fmt.Errorf("ope_encoder_drain error code: %v", int_err)
	}

	C.ope_encoder_destroy(s.oggencoder)

	if closer, ok := s.write.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
