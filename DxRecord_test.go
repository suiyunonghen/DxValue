package DxValue

import (
	"testing"
	"io/ioutil"
	"fmt"
	"encoding/json"
	//"github.com/json-iterator/go"
	"github.com/suiyunonghen/DxCommonLib"
	"unsafe"
	"time"
	"bytes"
)

func TestDxRecord_JsonParserFromByte(t *testing.T) {
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
	fmt.Println(rc.ToString())
	rc.SaveMsgPackFile("d:\\testMsgPack.bin")
}

func TestParserTime(t *testing.T)  {
	fmt.Println(time.Now().Format("2006-01-02T15:04:05Z"))
	fmt.Println(time.Parse("2006-01-02T15:04:05Z","2010-07-12T03:05:21Z"))
	at := DxCommonLib.ParserJsonTime("/Date(1402384458000)/")
	fmt.Println(at.ToTime().Format("2006-01-02T15:04:05Z"))
	at = DxCommonLib.ParserJsonTime("/Date(1224043200000+0800)/")
	fmt.Println(at.ToTime().Format("2006-01-02T15:04:05Z"))
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


/*func BenchmarkJsoniterParser(b *testing.B){
	buf, err := ioutil.ReadFile("DataProxy.config.json")
	if err != nil {
		fmt.Println("ReadFile Err:",err)
		return
	}
	mp := make(map[string]interface{})
	for i := 0;i<b.N;i++ {
		jsoniter.Unmarshal(buf,&mp)
	}
}*/

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

func TestDxValue_JsonParserFromByte(t *testing.T) {
	var v DxValue
	buf, err := ioutil.ReadFile("DataProxy.config.json")
	if err != nil {
		fmt.Println("ReadFile Err:",err)
		return
	}
	_,err = v.JsonParserFromByte(buf,false)
	if err != nil{
		fmt.Println("Parser Error: ",err)
	}else{
		switch v.ValueType() {
		case DVT_Record:
			fmt.Println("Is Json Object: ",(*DxRecord)(unsafe.Pointer(v.fValue)).ToString())
		case DVT_Array:
			fmt.Println("Is Json Array: ",(*DxArray)(unsafe.Pointer(v.fValue)).ToString())
		}
	}
}

func TestDxRecord_SaveJsonFile(t *testing.T) {
	rec := NewRecord()
	rec.SetInt("Age",-12)
	rec.SetString("Name","suiyunonghen")
	rec.SetValue("Home",map[string]interface{}{
		"Addres": "湖北武汉",
		"code":"430000",
		"Peoples":4,
	})
	rec.SetDouble("Double",234234234.4564564)
	rec.SetFloat("Float",-34.534)
	rec.SetValue("Now",time.Now())
	//rec.SaveJsonFile("d:\\testJson.json",true)
	rec.SaveMsgPackFile("d:\\msgpack.bin")
}

func TestMsgPackDecode(t *testing.T)  {
	bt, err := ioutil.ReadFile("test.Msgpack")
	if err != nil {
		fmt.Println("ReadFile Err:",err)
		return
	}
	buf := bytes.NewBuffer(bt)
	if rec,err := DecodeMsgPack(buf);err!=nil{
		fmt.Println("Error；",err)
	}else{
		fmt.Println(rec.ToString())
	}
}

func TestDxRecord_LoadMsgPackFile(t *testing.T) {
	rec := NewRecord()
	if err := rec.LoadMsgPackFile("test.Msgpack");err!=nil{
		fmt.Println("Error；",err)
	}else{
		fmt.Println(rec.ToString())
	}

}


func TestDxRecord_AsString(t *testing.T) {
	rc := NewRecord()
	rc.JsonParserFromByte([]byte(`{"StringValue":   "String    3"  ,  "object":  {"objStr":"  ObjStr1  ",  "ObjName"  :  "   I  nnerObje  ct  "}}`),false)
	fmt.Println("StringValue=",rc.AsString("StringValue",""))
	fmt.Println("object.objStr=",rc.AsStringByPath("object.objStr",""))
	fmt.Println("object.ObjName=",rc.AsStringByPath("object.ObjName",""))
}

func TestDxRecord_SetIntRecordValue(t *testing.T) {
	rc := NewRecord()
	inarc := NewIntKeyRecord()
	inarc.SetInt(2,23)
	inarc.SetValue(23,"DxSoft")
	rc.SetIntRecordValue("IntRecord",inarc)
	fmt.Println(rc.Contains("IntRecord.23"))
	if !rc.Contains("IntRecord.ConVertName"){
		rc.ForcePath("IntRecord.ConVertName","Record")
		fmt.Println(rc.Contains("IntRecord.ConVertName"))
	}
	fmt.Println(rc.ToString())
}