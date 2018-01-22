package main

import (
	"github.com/suiyunonghen/DxCommonLib"
	"fmt"
	"github.com/suiyunonghen/DxValue"
)

func main()  {
	mb := "aBCA得闲"
	bt := DxCommonLib.Binary2Hex(([]byte)(mb))
	fmt.Println(bt)
	mbt := DxCommonLib.Hex2Binary(string(bt))
	fmt.Println(mbt)
	fmt.Println(string(mbt))

	mrec := DxValue.NewRecord()
	mrec.SetValue("Name","DxSoft")
	mrec.SetValue("Age",23)
	mrec.SetString("Sex","男")
	crec := mrec.NewRecord("Home")
	crec.SetString("Father","ParentF")
	crec.SetString("Mother","ParentM")
	fmt.Print(mrec.ToString())

}
