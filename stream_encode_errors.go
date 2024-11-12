// Copyright © 2015-2017 Go Opus Authors (see AUTHORS file)
//
// License for use of this code is detailed in the LICENSE file

//go:build !nolibopusfile
// +build !nolibopusfile

package opus

/*
#cgo pkg-config: libopusenc
#include <opusenc.h>
*/
import "C"

// StreamError represents an error from libopusenc.
type StreamEncodeError int

var _ error = StreamEncodeError(0)

// Libopusfile errors. The names are copied verbatim from the libopusfile
// library.
const (
	ErrStreamEncoderBadArg         = StreamEncodeError(C.OPE_BAD_ARG)
	ErrStreamEncoderInternalError  = StreamEncodeError(C.OPE_INTERNAL_ERROR)
	ErrStreamEncoderUninplemented  = StreamEncodeError(C.OPE_UNIMPLEMENTED)
	ErrStreamEncoderAllocFailed    = StreamEncodeError(C.OPE_ALLOC_FAIL)
	ErrStreamEncoderCannotOpen     = StreamEncodeError(C.OPE_CANNOT_OPEN)
	ErrStreamEncoderTooLate        = StreamEncodeError(C.OPE_TOO_LATE)
	ErrStreamEncoderInvalidPicture = StreamEncodeError(C.OPE_INVALID_PICTURE)
	ErrStreamEncoderInvalidIcon    = StreamEncodeError(C.OPE_INVALID_ICON)
	ErrStreamEncoderWriteFail      = StreamEncodeError(C.OPE_WRITE_FAIL)
	ErrStreamEncoderCloseFail      = StreamEncodeError(C.OPE_CLOSE_FAIL)
)

func (i StreamEncodeError) Error() string {
	switch i {
	case ErrStreamEncoderBadArg:
		return "OPE_BAD_ARG"
	case ErrStreamEncoderInternalError:
		return "OPE_INTERNAL_ERROR"
	case ErrStreamEncoderUninplemented:
		return "OPE_UNIMPLEMENTED"
	case ErrStreamEncoderAllocFailed:
		return "OPE_ALLOC_FAIL"
	case ErrStreamEncoderCannotOpen:
		return "OPE_CANNOT_OPEN"
	case ErrStreamEncoderTooLate:
		return "OPE_TOO_LATE"
	case ErrStreamEncoderInvalidPicture:
		return "OPE_INVALID_PICTURE"
	case ErrStreamEncoderInvalidIcon:
		return "OPE_INVALID_ICON"
	case ErrStreamEncoderWriteFail:
		return "OPE_WRITE_FAIL"
	case ErrStreamEncoderCloseFail:
		return "OPE_CLOSE_FAIL"
	default:
		return "libopusenc error: %d (unknown code)"
	}
}
