package DxValue

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/json-iterator/go"
	"gitee.com/johng/gf/g/encoding/gparser"
	"testing"
	//"os"
	"os"
)

func BenchmarkGparser(b *testing.B){
	buf, err := ioutil.ReadFile("DataProxy.config.json")
	if err != nil {
		fmt.Println("ReadFile Err:",err)
		return
	}
	var gp *gparser.Parser
	for i := 0;i<b.N;i++ {
		gp,_ = gparser.LoadContent(buf, "json")
		if bt,err := gp.ToJson();err!=nil{
			return
		}else if file,err := os.OpenFile("Gparser.json",os.O_CREATE | os.O_TRUNC,0644);err == nil{
			file.Write(bt)
			file.Close()
		}
	}
}

func BenchmarkDxRecord_JsonParserFromByte(b *testing.B) {
	buf, err := ioutil.ReadFile("DataProxy.config.json")
	if err != nil {
		fmt.Println("ReadFile Err:",err)
		return
	}
	rc := NewRecord()
	for i := 0;i<b.N;i++{
		_,err := rc.JsonParserFromByte(buf,false,false)
		if err != nil{
			fmt.Println("Parser Error: ",err)
			break
		}
		if file,err := os.OpenFile("DxRecord.json",os.O_CREATE | os.O_TRUNC,0644);err == nil{
			rc.SaveJsonWriter(file)
			file.Close()
		}
	}
}


func BenchmarkJsoniterParser(b *testing.B){
	buf, err := ioutil.ReadFile("DataProxy.config.json")
	if err != nil {
		fmt.Println("ReadFile Err:",err)
		return
	}
	mp := make(map[string]interface{})
	for i := 0;i<b.N;i++ {
		jsoniter.Unmarshal(buf,&mp)
		if bt,err := jsoniter.Marshal(&mp);err!=nil{
			return
		}else if file,err := os.OpenFile("Jsoniter.json",os.O_CREATE | os.O_TRUNC,0644);err == nil{
			file.Write(bt)
			file.Close()
		}
	}
}

func BenchmarkStandJsonParser(b *testing.B){
	buf, err := ioutil.ReadFile("DataProxy.config.json")
	if err != nil {
		fmt.Println("ReadFile Err:",err)
		return
	}
	mp := make(map[string]interface{})
	for i := 0;i<b.N;i++ {
		json.Unmarshal(buf, &mp)
		if bt,err := json.Marshal(&mp);err!=nil{
			return
		}else if file,err := os.OpenFile("StandJson.json",os.O_CREATE | os.O_TRUNC,0644);err == nil{
			file.Write(bt)
			file.Close()
		}
	}
}

