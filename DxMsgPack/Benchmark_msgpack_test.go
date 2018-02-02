package DxMsgPack

import (
	"testing"
	"io/ioutil"
	"fmt"
	"github.com/vmihailenco/msgpack"
	"os"
	"io"
	"bufio"
	"github.com/suiyunonghen/DxValue"
)

func Benchmark_DecodeMsgPack(b *testing.B) {
	f, err := os.Open("DataProxy.config.msgPack")
	if err != nil {
		fmt.Println("BenchmarkDecodeMsgPack Error:",err)
		return
	}
	defer f.Close()
	coder := MsgPackDecoder{}
	coder.r = bufio.NewReader(f)
	rec := DxValue.NewRecord()
	for i:=0;i<b.N;i++{
		if err := coder.Decode(f,&rec.DxBaseValue);err!=nil{
			return
		}
		f.Seek(0,io.SeekStart)
	}
}

func Benchmark_vmihailenco_decode(b *testing.B)  {
	if databytes, err := ioutil.ReadFile("DataProxy.config.msgPack");err!=nil{
		fmt.Println("BenchmarkDecodeMsgPack Error:",err)
		return
	}else{
		mp := make(map[string]interface{},32)
		for i:=0;i<b.N;i++{
			if err = msgpack.Unmarshal(databytes,&mp);err!=nil{
				fmt.Println("Benchmark_vmihailenco_decode Error:",err)
			}
		}
	}
}
