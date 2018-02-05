package DxMsgPack

import (
	"testing"
	"io/ioutil"
	"fmt"
	"github.com/vmihailenco/msgpack"
	"os"
	"io"
	"time"
)


type testPkg struct{
	Age			int
	Name		string
	Double		float64
	Float		float32
	Now			time.Time
}

func Benchmark_DecodeMsgPack(b *testing.B) {
	f, err := os.Open("DataProxy.config.msgPack")
	if err != nil {
		fmt.Println("BenchmarkDecodeMsgPack Error:",err)
		return
	}
	defer f.Close()
	coder := NewDecoder(f)
	mp := make(map[string]interface{},32)
	//mp := testPkg{}
	//Coders.RegisterType(reflect.TypeOf(mp))
	for i:=0;i<b.N;i++{
		if err = coder.DecodeStand(&mp);err!=nil{
			fmt.Println("BenchmarkDecodeMsgPack Error:",err)
		}
		f.Seek(0,io.SeekStart)
	}
	//fmt.Println(mp)
}

func Benchmark_vmihailenco_decode(b *testing.B)  {
	if databytes, err := ioutil.ReadFile("DataProxy.config.msgPack");err!=nil{
		fmt.Println("BenchmarkDecodeMsgPack Error:",err)
		return
	}else{
		mp := make(map[string]interface{},32)
		//mp := testPkg{}
		for i:=0;i<b.N;i++{
			if err = msgpack.Unmarshal(databytes,&mp);err!=nil{
				fmt.Println("Benchmark_vmihailenco_decode Error:",err)
			}
		}
	}
}
