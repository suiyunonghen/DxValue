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
	_,err = rc.JsonParserFromByte(buf,false)
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
		_,err = rc.JsonParserFromByte(buf,false)
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
	rc.JsonParserFromByte([]byte(`{"BoolValue":  true  ,"object":{"objBool":  false  }}`),false)
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
	_,err = rc.JsonParserFromByte(buf,true)
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
				if i == 1{
					fmt.Println("input.remark=",rc.AsStringByPath("input.remark","unkonwn"))
				}
				fmt.Println(rc.ToString())
			}
		}
	}
}

func TestEsCapteStr(t *testing.T)  {
	str := "\u6821\u56ed\u7f51"
	//DxCommonLib.FastByte2String([]byte(str))
	str = fmt.Sprintf("%v",str)
}


func TestDxRecord_AsString(t *testing.T) {
	rc := NewRecord()
	rc.JsonParserFromByte([]byte(`{"StringValue":   "String    3"  ,  "object":  {"objStr":"  ObjStr1  ",  "ObjName"  :  "   I  nnerObje  ct  "}}`),false)
	fmt.Println("StringValue=",rc.AsString("StringValue",""))
	fmt.Println("object.objStr=",rc.AsStringByPath("object.objStr",""))
	fmt.Println("object.ObjName=",rc.AsStringByPath("object.ObjName",""))
}