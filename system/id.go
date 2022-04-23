package system

import (
	"log"
	"math/rand"
	"net"
	"unsafe"
)

var id string

// Generates a 23-character random string using a MAC address (as byte array) as the seed.
func generateClientID(mac []byte) string {
	const idLen = 23
	const letterBytes = "abcdefghijklmnopqrstuvwxyz"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)

	seed := int64(0)
	for _, b := range mac {
		seed += int64(b)
	}

	src := rand.NewSource(seed)
	b := make([]byte, idLen)
	for i, cache, remain := idLen-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

// Returns (or generates if needed) the 23-character ID for this node using
// the MAC address of the first available network interface as seed.
func ID() string {
	if id != "" {
		return id
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	for _, iface := range ifaces {
		mac := iface.HardwareAddr
		if len(mac) > 0 {
			id = generateClientID(iface.HardwareAddr)
			break
		}
	}

	return id
}
