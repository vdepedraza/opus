// Copyright © Go Opus Authors (see AUTHORS file)
//
// License for use of this code is detailed in the LICENSE file

//go:build !nolibopusfile
// +build !nolibopusfile

package opus

import (
	"sync"
	"sync/atomic"
)

type encoderStreamsMap struct {
	sync.RWMutex
	m       map[uintptr]*EncoderStream
	counter uintptr
}

func (sm *encoderStreamsMap) Get(id uintptr) *EncoderStream {
	sm.RLock()
	defer sm.RUnlock()
	return sm.m[id]
}

func (sm *encoderStreamsMap) Del(s *EncoderStream) {
	sm.Lock()
	defer sm.Unlock()
	delete(sm.m, s.id)
}

// NextId returns a unique ID for each call.
func (sm *encoderStreamsMap) NextId() uintptr {
	return atomic.AddUintptr(&sm.counter, 1)
}

func (sm *encoderStreamsMap) Save(s *EncoderStream) {
	sm.Lock()
	defer sm.Unlock()
	sm.m[s.id] = s
}

func newEncoderStreamsMap() *encoderStreamsMap {
	return &encoderStreamsMap{
		counter: 0,
		m:       map[uintptr]*EncoderStream{},
	}
}
