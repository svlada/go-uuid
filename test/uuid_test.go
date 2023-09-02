package uuid

import (
	"fmt"
	v1 "svlada.com/uuid/v1"
	"testing"
)

func TestUUIDv1(t *testing.T) {
	var a, _ = v1.UUIDv1()
	fmt.Println(a.ToString())
}
