package DxJson

import (
	"testing"
	"bytes"
	"fmt"
	"encoding/json"
)

type Testmm struct {
	MP			string 	`json:"血量"`
	HP			string 	`json:"魔法"`
	TG			[]interface{} `json:"tg"`
}

func TestJsonEncoder_EncodeStand(t *testing.T) {
	var buffer bytes.Buffer
	encoder := NewEncoder(&buffer)
	g := Testmm{MP:"23423",HP:"asdfasdf"}
	g.TG = make([]interface{},4)
	g.TG[0] = "TG0"
	g.TG[1] = 234
	g.TG[2] = true
	g.TG[3] = "TG3"
	for i := 0;i<1;i++{
		encoder.EncodeStand(g)
	}
	fmt.Println(string(buffer.Bytes()))
}

func BenchmarkJsonEncoder_EncodeStand(b *testing.B) {
	//var buffer bytes.Buffer
	g := Testmm{MP:"23423",HP:"asdfasdf"}
	g.TG = make([]interface{},4)
	g.TG[0] = "TG0"
	g.TG[1] = 234
	g.TG[2] = true
	g.TG[3] = "TG3"
	//encoder := NewEncoder(&buffer)
	//var bt []byte
	for i := 0;i<b.N;i++{
		//buffer.Reset()
		//encoder.EncodeStand(g)
		Marshal(g)

	}
	//fmt.Println(string(bt))
}

func BenchmarkJsonEncoder_Stand(b *testing.B) {
	//var buffer bytes.Buffer
	g := Testmm{MP:"23423",HP:"asdfasdf"}
	g.TG = make([]interface{},4)
	g.TG[0] = "TG0"
	g.TG[1] = 234
	g.TG[2] = true
	g.TG[3] = "TG3"
	//encoder := NewEncoder(&buffer)
	//var bt []byte
	for i := 0;i<b.N;i++{
		//buffer.Reset()
		//encoder.EncodeStand(g)
		json.Marshal(g)

	}
	//fmt.Println(string(bt))
}

