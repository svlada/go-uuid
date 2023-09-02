package uuid

import (
	"encoding/hex"
	"fmt"
)

type UUID [16]byte

func (u UUID) ToString() string {
	return fmt.Sprintf(
		"%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		u[0], u[1], u[2], u[3], u[4], u[5], u[6], u[7], u[8], u[9], u[10], u[11], u[12], u[13], u[14], u[15],
	)
}

func (u UUID) ToStringRaw() string {
	return hex.EncodeToString(u[:])
}
