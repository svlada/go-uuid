package v1

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
	"sync"
	time "time"

	"svlada.com/uuid"
)

func init() {
	gob.Register(UUIDGeneratorState{})
	fileStrategy := &FilePersistence{Filename: "uuid_state.bin"}
	var err error
	g, err = NewUUIDGenerator(fileStrategy)
	if err != nil {
		panic("Failed to initialize UUID generator: " + err.Error())
	}
}

type ClockSequence uint16

// Increment When the sequence reaches its 14-bit maximum, it wraps back around to 0.
func (cs *ClockSequence) Increment() {
	*cs = (*cs + 1) & 0x3FFF
}

type UUIDGenerator struct {
	State       UUIDGeneratorState
	Mu          sync.Mutex
	Persistence StateStore
}

type UUIDGeneratorState struct {
	ClockSeq ClockSequence
	LastTime time.Time
	Node     []byte
}

func (gen *UUIDGenerator) Generate() (uuid.UUID, error) {
	var uuidValue uuid.UUID

	timestamp := time.Now().UnixNano()/100 + 122192928000000000
	timeLow := uint32(timestamp & 0xFFFFFFFF)
	timeMid := uint16((timestamp >> 32) & 0xFFFF)
	timeHiAndVersion := uint16((timestamp >> 48) & 0x0FFF)
	timeHiAndVersion |= 0x1000

	binary.BigEndian.PutUint32(uuidValue[0:], timeLow)
	binary.BigEndian.PutUint16(uuidValue[4:], timeMid)
	binary.BigEndian.PutUint16(uuidValue[6:], timeHiAndVersion)

	gen.Mu.Lock()
	defer gen.Mu.Unlock()

	currentTime := time.Now()
	if currentTime.Before(gen.State.LastTime) || currentTime.Equal(gen.State.LastTime) {
		gen.State.ClockSeq.Increment()
	}

	clockSeqLow := uint8(gen.State.ClockSeq & 0xFF)
	clockSeqHi := uint8((gen.State.ClockSeq >> 8) & 0x3F)
	clockSeqHi |= 0x80
	binary.BigEndian.PutUint16(uuidValue[8:10], uint16(clockSeqHi)<<8|uint16(clockSeqLow))

	node, err := getHardwareAddr()
	if err != nil {
		fmt.Println("Could not get MAC address:", err)
		node = make([]byte, 6)
		if _, randErr := rand.Read(node); randErr != nil {
			return uuidValue, randErr
		}
	}

	copy(uuidValue[10:], node)

	if fileErr := gen.Persistence.Save(&gen.State); err != nil {
		return uuidValue, fmt.Errorf("error saving generator state: %w", fileErr)
	}

	return uuidValue, nil
}

var g *UUIDGenerator

func UUIDv1() (uuid.UUID, error) {
	return g.Generate()
}

type StateStore interface {
	Load() (*UUIDGeneratorState, error)
	Save(*UUIDGeneratorState) error
}

func NewUUIDGenerator(store StateStore) (*UUIDGenerator, error) {
	state, err := store.Load()
	if err != nil {
		randClockSeq, err := rand.Int(rand.Reader, big.NewInt(1<<14))
		if err != nil {
			panic("Could not initialize clock sequence: " + err.Error())
		}

		node, nodeErr := getHardwareAddr()
		if nodeErr != nil {
			return nil, nodeErr
		}

		state = &UUIDGeneratorState{
			ClockSeq: ClockSequence(randClockSeq.Int64()),
			Node:     node,
			LastTime: time.Now(),
		}

		err = store.Save(state)
		if err != nil {
			return nil, err
		}
	}

	return &UUIDGenerator{
		State:       *state,
		Persistence: store,
	}, nil
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

type FilePersistence struct {
	Filename string
}

func (fp *FilePersistence) Load() (*UUIDGeneratorState, error) {
	file, err := os.Open(fp.Filename)
	if err != nil {
		if os.IsNotExist(err) {
			return &UUIDGeneratorState{}, nil
		}
		return nil, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	state := &UUIDGeneratorState{}
	if err := decoder.Decode(state); err != nil {
		return nil, err
	}

	return state, nil
}

func (fp *FilePersistence) Save(state *UUIDGeneratorState) error {
	file, err := os.Create(fp.Filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(state); err != nil {
		return err
	}

	return nil
}
