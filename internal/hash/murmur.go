package hash

import (
	"fmt"
	"math/rand"
	"unsafe"
)

/*
The below hasher type is a is based on the Murmur3 algorithm but the handling of
tail bytes (eg bytes that do not immediately fill out a four byte block) has been
modified to reduce the need for copying data.

The implementation is highly inspired by the 32 bit hash found in
github.com/spaolacci/murmur3. See license below.

LICENSE
-------
Copyright 2013, Sébastien Paolacci.
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:
    * Redistributions of source code must retain the above copyright
      notice, this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright
      notice, this list of conditions and the following disclaimer in the
      documentation and/or other materials provided with the distribution.
    * Neither the name of the library nor the
      names of its contributors may be used to endorse or promote products
      derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

const (
	c1_32   uint32 = 0xcc9e2d51
	c2_32   uint32 = 0x1b873593
	tailLen        = 4
)

type Murm32 struct {
	totLen   int
	hash     uint32
	tail     [tailLen]byte
	tailSize int8
}

// Write adds a byte as input to the hash.
func (m *Murm32) WriteByte(b byte) {
	m.tail[m.tailSize] = b
	m.tailSize++
	m.flushBufIfNeeded()
}

func (m *Murm32) flushBufIfNeeded() {
	if m.tailSize == tailLen {
		m.tailSize = 0
		m.Write(m.tail[:])
	}
}

func intMin(x, y int) int {
	if y < x {
		return y
	}
	return x
}

// Write adds data as input to the hash.
func (m *Murm32) Write(data []byte) {
	if len(data) == 0 {
		return
	}

	if m.tailSize > 0 {
		// If a previous tail exists we first want to fill that up
		// and hash it if full before hashing any remaining bytes.
		copyCount := copy(m.tail[m.tailSize:], data)
		m.tailSize += int8(copyCount)
		m.flushBufIfNeeded()
		m.Write(data[copyCount:])
		return
	}

	h1 := m.hash
	nblocks := len(data) / 4
	m.totLen += len(data)
	p := uintptr(unsafe.Pointer(&data[0]))

	// Hash full 4-byte blocks
	p1 := p + uintptr(4*nblocks)
	for ; p < p1; p += 4 {
		k1 := *(*uint32)(unsafe.Pointer(p))

		k1 *= c1_32
		k1 = (k1 << 15) | (k1 >> 17) // rotl32(k1, 15)
		k1 *= c2_32

		h1 ^= k1
		h1 = (h1 << 13) | (h1 >> 19) // rotl32(h1, 13)
		h1 = h1*4 + h1 + 0xe6546b64
	}

	// Store any remaining bytes in tail
	tail := data[nblocks*4:]
	m.tailSize = int8(copy(m.tail[:], tail))

	m.hash = h1
}

// Reset clears internal state so that the hasher can be reused
func (m *Murm32) Reset() {
	m.tailSize = 0
	m.hash = 0
	m.totLen = 0
}

// Hash returns the final hash value.
func (m *Murm32) Hash() uint32 {
	var k1 uint32
	h1 := m.hash

	// Process any bytes remaining
	switch m.tailSize {
	case 3:
		k1 ^= uint32(m.tail[2]) << 16
		fallthrough
	case 2:
		k1 ^= uint32(m.tail[1]) << 8
		fallthrough
	case 1:
		k1 ^= uint32(m.tail[0])
		k1 *= c1_32
		k1 = (k1 << 15) | (k1 >> 17) // rotl32(k1, 15)
		k1 *= c2_32
		h1 ^= k1
	case 0:
		// Nothing to do
	default:
		panic(fmt.Sprintf("Unexpected tail length %d, this is an implementation bug, please report it!", m.tailSize))
	}

	h1 ^= uint32(m.totLen)

	h1 ^= h1 >> 16
	h1 *= 0x85ebca6b
	h1 ^= h1 >> 13
	h1 *= 0xc2b2ae35
	h1 ^= h1 >> 16
	return h1
}

// WriteRand32 adds 32 random bits to the input data of the hash.
func (m *Murm32) WriteRand32() {
	var nullHashBytes [4]byte
	hashBytes := nullHashBytes[:]
	rand.Read(hashBytes)
	m.Write(hashBytes)
}
