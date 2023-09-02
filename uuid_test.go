package uuid

import (
	"fmt"
	"testing"
)

func TestUUIDv1(t *testing.T) {
	var a, _ = UUIDv1()
	fmt.Println(a.ToString())
}
