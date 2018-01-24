package main

import (
	"github.com/suiyunonghen/DxCommonLib"
	"fmt"
	"github.com/suiyunonghen/DxValue"
	//"encoding/json"
)

type DxPeople struct {
	name  string
	age   int
	Sex   string
}

func main()  {
	mb := "aBCA得闲"
	bt := DxCommonLib.Binary2Hex(([]byte)(mb))
	fmt.Println(bt)
	mbt := DxCommonLib.Hex2Binary(string(bt))
	fmt.Println(mbt)
	fmt.Println(string(mbt))

	mrec := DxValue.NewRecord()
	_,err := mrec.JsonParserFromByte([]byte(`{"DxSoft":{"Name":"不得闲"},"Age":32,"Name":true,"testArray":["gg",23,"gasdf"]}`))//
	if err != nil{
		panic(err)
	}
	fmt.Println(mrec.ToString())

	mrec.SetValue("Name","DxSoft")
	mrec.SetValue("Age",23)
	mrec.SetString("Sex","男")
	crec := mrec.NewRecord("Home")
	crec.SetString("Father","ParentF")
	crec.SetString("Mother","ParentM")



	mp := make(map[string]*DxPeople,6)
	father := new(DxPeople)
	father.name = "HuPing"
	father.age = 30
	mp["Father"] = father

	mrec.SetValue("Sun",mp)

	mrec.SetValue("testFather",father)

	fmt.Println(mrec.ToString())
	//json.Marshal()

	mrec.ClearValue()
	mrec.ForcePath("DxSoft.Name","不得闲")
	fmt.Println(mrec.ToString())

	ma := DxValue.NewArray()

	_,err = ma.JsonParserFromByte([]byte(`[null,["gg",23,"gasdf"],20,null ,null ,12,{ "Name": "不得闲" }]`))
	if err != nil{
		panic(err)
	}
	fmt.Println(ma.ToString())

	ma.SetInt(1,20)
	ma.SetInt(4,12)
	fmt.Println(ma.ToString())

}
