package DxValue

import (
	"testing"
	"fmt"
)

func TestDxArray_JsonParserFromByte(t *testing.T) {
	arr := NewArray()
	_,err := arr.JsonParserFromByte([]byte(`[  32  ,  "2342"  ,[ 2 , true , false  ,{ "Name" : "DxSoft" , "Age"  :  32 } ] ]`),false)
	if err == nil {
		fmt.Println(arr.ToString())
	}else{
		fmt.Println("Paser Error")
	}
}

func TestDxArray_LoadJsonFile(t *testing.T) {
	var v DxValue
	v.LoadJsonFile("DataProxy.config.json",true)
	fmt.Println(v.AsString())
}

func TestDxArray_SaveJsonFile(t *testing.T) {
	var v DxValue
	v.LoadJsonFile("DataProxy.config.json",true)
	if rec,_ := v.AsRecord();rec != nil{
		if arr := rec.AsArray("list");arr!=nil{
			arr.SaveJsonFile("d:\\1.json",true)
			fmt.Println("SaveJsonOK")
		}
	}
}
