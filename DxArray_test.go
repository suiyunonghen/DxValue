package DxValue

import (
	"testing"
	"fmt"
)

func TestDxArray_JsonParserFromByte(t *testing.T) {
	arr := NewArray()
	_,err := arr.JsonParserFromByte([]byte(`[  32  ,  "2342"  ,[ 2 , true , false  ,{ "Name" : "DxSoft" , "Age"  :  32 } ] ]`))
	if err == nil {
		fmt.Println(arr.ToString())
	}else{
		fmt.Println("Paser Error")
	}
}
