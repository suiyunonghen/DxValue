package DxValue

import (
	"testing"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"github.com/json-iterator/go"
)

func TestDxRecord_JsonParserFromByte(t *testing.T) {
	buf, err := ioutil.ReadFile("DataProxy.config.json")
	if err != nil {
		fmt.Println("ReadFile Err:",err)
		return
	}
	rc := NewRecord()
	_,err = rc.JsonParserFromByte(buf)
	if err != nil{
		fmt.Println("Parser Error: ",err)
	}
	fmt.Println(rc.ToString())
}

func BenchmarkDxRecord_JsonParserFromByte(b *testing.B) {
	buf, err := ioutil.ReadFile("DataProxy.config.json")
	if err != nil {
		fmt.Println("ReadFile Err:",err)
		return
	}
	rc := NewRecord()
	for i := 0;i<b.N;i++{
		_,err = rc.JsonParserFromByte(buf)
		if err != nil{
			fmt.Println("Parser Error: ",err)
			break
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
	}
}


func TestDxRecord_AsBool(t *testing.T) {
	rc := NewRecord()
	rc.JsonParserFromByte([]byte(`{"BoolValue":  true  ,"object":{"objBool":  false  }}`))
	fmt.Println("BoolValue=",rc.AsBool("BoolValue",false))
	fmt.Println("object.objBool=",rc.AsBoolByPath("object.objBool",true))
}

func TestDxRecord_AsArray(t *testing.T) {
	buf, err := ioutil.ReadFile("DataProxy.config.json")
	if err != nil {
		fmt.Println("ReadFile Err:",err)
		return
	}
	rc := NewRecord()
	_,err = rc.JsonParserFromByte(buf)
	if err != nil{
		fmt.Println("Parser Error: ",err)
	}
	array := rc.AsArray("list")
	if array != nil{
		for i := 0;i<array.Length();i++{
			fmt.Print("Array Index ",i,"=")
			switch array.VaueTypeByIndex(0) {
			case DVT_Record:
				rc := array.AsRecord(i)
				fmt.Println(rc.ToString())
			}
		}
	}
}

func TestDxRecord_AsString(t *testing.T) {
	rc := NewRecord()
	rc.JsonParserFromByte([]byte(`{"StringValue":   "String    3"  ,  "object":  {"objStr":"  ObjStr1  ",  "ObjName"  :  "   I  nnerObje  ct  "}}`))
	fmt.Println("StringValue=",rc.AsString("StringValue",""))
	fmt.Println("object.objStr=",rc.AsStringByPath("object.objStr",""))
	fmt.Println("object.ObjName=",rc.AsStringByPath("object.ObjName",""))
}