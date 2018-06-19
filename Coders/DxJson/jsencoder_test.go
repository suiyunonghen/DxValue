package DxJson

import (
	"testing"
	"bytes"
	"fmt"
)

type Testmm struct {
	MP			string 	`json:"血量"`
	HP			string 	`json:"魔法"`
}

func TestJsonEncoder_EncodeStand(t *testing.T) {
	var buffer bytes.Buffer
	encoder := NewEncoder(&buffer)
	g := Testmm{"23423","asdfasdf"}
	encoder.EncodeStand(g)
	fmt.Println(string(buffer.Bytes()))
}