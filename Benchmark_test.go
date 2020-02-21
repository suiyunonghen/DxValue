package DxValue

import (
	"encoding/json"
	"fmt"

	//"github.com/gogf/gf/g/encoding/gparser"
	"github.com/json-iterator/go"
	"github.com/valyala/fastjson"
	"io/ioutil"
	"testing"
)

func BenchmarkJsonParse(b *testing.B){
	buf, err := ioutil.ReadFile("DataProxy.config.json")
	if err != nil{
		return
	}
	b.Run("std", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(buf)))
		b.RunParallel(func(pb *testing.PB) {
			mp := make(map[string]interface{},100)
			for pb.Next() {
				json.Unmarshal(buf, &mp)
				/*if bt,err := json.Marshal(&mp);err!=nil{
					return
				}else if file,err := os.OpenFile("StandJson.json",os.O_CREATE | os.O_TRUNC,0644);err == nil{
					file.Write(bt)
					file.Close()
				}*/
			}
		})
	})
	b.Run("Jsoniter", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(buf)))
		b.RunParallel(func(pb *testing.PB) {
			mp := make(map[string]interface{},100)
			for pb.Next() {
				jsoniter.Unmarshal(buf,&mp)
				/*if bt,err := jsoniter.Marshal(&mp);err!=nil{
					return
				}else if file,err := os.OpenFile("Jsoniter.json",os.O_CREATE | os.O_TRUNC,0644);err == nil{
					file.Write(bt)
					file.Close()
				}*/
			}
		})
	})
	b.Run("DxRecord", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(buf)))
		b.RunParallel(func(pb *testing.PB) {
			rc := NewRecord()
			for pb.Next() {
				rc.JsonParserFromByte(buf,false,false)
			}
		})
	})
	b.Run("fastJson", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(buf)))
		//var parsePool fastjson.ParserPool
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				/*p := parsePool.Get()
				p.ParseBytes(buf)
				parsePool.Put(p)*/
				fastjson.MustParseBytes(buf)
			}
		})
	})
}

func TestFastJson(t *testing.T)  {
	buf, err := ioutil.ReadFile("DataProxy.config.json")
	if err != nil{
		return
	}
	var parsePool fastjson.ParserPool
	p := parsePool.Get()
	v,err := p.ParseBytes(buf)
	if err != nil{
		return
	}
	fmt.Println(v.String())
	parsePool.Put(p)
}