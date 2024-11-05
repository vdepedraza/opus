// +build !nolibopusfile

// Copyright © Go Opus Authors (see AUTHORS file)
//
// License for use of this code is detailed in the LICENSE file

// Allocate callback struct in C to ensure it's not managed by the Go GC. This
// plays nice with the CGo rules and avoids any confusion.

#include <opusfile.h>
#include <opusenc.h>
#include <stdint.h>

// Defined in Go. Uses the same signature as Go, no need for proxy function.
int go_readcallback(void *p, unsigned char *buf, int nbytes);

static struct OpusFileCallbacks streamCallbacks = {
    .read = go_readcallback
};

// Proxy function for op_open_callbacks, because it takes a void * context but
// we want to pass it non-pointer data, namely an arbitrary uintptr_t
// value. This is legal C, but go test -race (-d=checkptr) complains anyway. So
// we have this wrapper function to shush it.
// https://groups.google.com/g/golang-nuts/c/995uZyRPKlU
OggOpusFile *
my_open_callbacks(uintptr_t p, int *error)
{
    return op_open_callbacks((void *)p, &streamCallbacks, NULL, 0, error);
}

// encoder callbacks
int go_writecallback(void *p, const unsigned char *buf, int nbytes);

int closeCallback(void *p){
    return 0;
}

static OpusEncCallbacks encoderCallbacks = {
    .write = go_writecallback,
    .close = closeCallback,
};

OggOpusEnc * my_ope_encoder_create_callback(uintptr_t p, int *error){
    // 48khz, mono channel 
    OggOpusComments * comments = ope_comments_create();
    return ope_encoder_create_callbacks(&encoderCallbacks, (void *)p, comments, 48000, 1, 0, error);

}