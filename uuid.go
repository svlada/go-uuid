package uuid

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net"
	"sync"
	time "time"
)

type UUID [16]byte
type ClockSequence uint16

// Increment When the sequence reaches its 14-bit maximum, it wraps back around to 0.
func (cs *ClockSequence) Increment() {
	*cs = (*cs + 1) & 0x3FFF
}

var clockSequence ClockSequence
var clockMutex sync.Mutex
var lastTimestamp time.Time

func init() {
	randClockSeq, _ := rand.Int(rand.Reader, big.NewInt(1<<14))
	clockSequence = ClockSequence(randClockSeq.Int64())
}

func UUIDv1() (UUID, error) {
	var uuid UUID

	timestamp := time.Now().UnixNano()/100 + 122192928000000000
	timeLow := uint32(timestamp & 0xFFFFFFFF)
	timeMid := uint16((timestamp >> 32) & 0xFFFF)
	timeHiAndVersion := uint16((timestamp >> 48) & 0x0FFF)
	timeHiAndVersion |= 0x1000

	binary.BigEndian.PutUint32(uuid[0:], timeLow)
	binary.BigEndian.PutUint16(uuid[4:], timeMid)
	binary.BigEndian.PutUint16(uuid[6:], timeHiAndVersion)

	clockMutex.Lock()
	defer clockMutex.Unlock()

	currentTime := time.Now()
	if currentTime.Before(lastTimestamp) || currentTime.Equal(lastTimestamp) {
		clockSequence.Increment()
	}

	clockSeqLow := uint8(clockSequence & 0xFF)
	clockSeqHi := uint8((clockSequence >> 8) & 0x3F)
	clockSeqHi |= 0x80
	binary.BigEndian.PutUint16(uuid[8:10], uint16(clockSeqHi)<<8|uint16(clockSeqLow))

	node, err := getHardwareAddr()
	if err != nil {
		fmt.Println("Could not get MAC address:", err)
		node = make([]byte, 6)
		_, randErr := rand.Read(node)
		if randErr != nil {
			fmt.Println("Error generating random node:", err)
			return UUID{}, err
		}
	}

	copy(uuid[10:], node)

	return uuid, nil
}

func getHardwareAddr() ([]byte, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, inter := range interfaces {
		if inter.Flags&net.FlagUp == 0 || inter.Flags&net.FlagLoopback != 0 {
			continue
		}
		if len(inter.HardwareAddr) == 6 {
			return inter.HardwareAddr, nil
		}
	}
	return nil, errors.New("could not find a suitable MAC address")
}

func (u UUID) ToString() string {
	return fmt.Sprintf(
		"%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		u[0], u[1], u[2], u[3], u[4], u[5], u[6], u[7], u[8], u[9], u[10], u[11], u[12], u[13], u[14], u[15],
	)
}

func (u UUID) ToStringRaw() string {
	return hex.EncodeToString(u[:])
}
