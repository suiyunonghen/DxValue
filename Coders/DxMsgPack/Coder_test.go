package DxMsgPack

import (
	"testing"
	"os"
	"fmt"
	"github.com/suiyunonghen/DxValue/Coders"
)


type(
	people   struct{
		sex   string		`msgpack:"性别"`
		Name  string		`msgpack:"姓名"`
		Age   int			`msgpack:"年龄"`
	}

	student  struct{
		people
		ID    string		`msgpack:"学号"`
		Level string		`msgpack:"年级"`
	}

	teacher  struct{
		people
		TeacherID   string	`msgpack:"教师ID"`
		Lesson		string	`msgpack:"课程"`
	}

	classes   struct{
		Teachers   []*teacher	`msgpack:"老师列表"`
		Students   []*student	`msgpack:"学生列表"`
	}
)



func (p *people)Encode(c Coders.Encoder) error  {
	switch c.Name() {
	case "msgpack":
		encoder :=c.(*MsgPackEncoder)
		encoder.EncodeMapLen(3)
		encoder.EncodeString("sex")
		encoder.EncodeString(p.sex)

		encoder.EncodeString("Name")
		encoder.EncodeString(p.Name)

		encoder.EncodeString("Age")
		encoder.EncodeInt(int64(p.Age))
	}
	return nil
}

func (p *people)Decode(d Coders.Decoder) error  {
	switch d.Name() {
	case "msgpack":
		decoder := d.(*MsgPackDecoder)
		maplen,_ := decoder.DecodeMapLen(CodeUnkonw)



		if maplen == 3{
			fmt.Println(decoder.DecodeString(CodeUnkonw))
			b,_ := decoder.DecodeString(CodeUnkonw)
			p.sex = string(b)

			fmt.Println(decoder.DecodeString(CodeUnkonw))
			b,_ = decoder.DecodeString(CodeUnkonw)
			p.Name = string(b)


			fmt.Println(decoder.DecodeString(CodeUnkonw))
			v,_ := decoder.DecodeInt(CodeUnkonw)
			p.Age = int(v)

			fmt.Println(p)
		}
	}
	return nil
}

func TestMsgPackCode_ValueCocer(t *testing.T)  {
	b, err := Marshal(&people{sex: "男", Name: "DxSoft456456",Age:23})
	if err != nil {
		panic(err)
	}

	var v *people
	err = Unmarshal(b, &v)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v", v)
}

func TestMsgpackCoder_Stand(t *testing.T)  {
	myclass := new(classes)
	myclass.Students =  make([]*student,10)
	myclass.Teachers = make([]*teacher,10)
	for i := 0;i<5;i++{
		myclass.Teachers[i] = new(teacher)
		myclass.Teachers[i].Name = "teacher1"
		myclass.Teachers[i].Age = 28
		myclass.Teachers[i].Lesson = "语文"
		myclass.Teachers[i].TeacherID = "8902-42"
	}
	if file,err := os.OpenFile("test.msgpack",os.O_CREATE | os.O_TRUNC,0644);err == nil{
		defer file.Close()
		if bt,err := Marshal(myclass);err==nil{
			file.Write(bt)

			mclass2 := new(classes)
			Unmarshal(bt,&mclass2)
			fmt.Println(mclass2)
		}else{
			fmt.Println("Marshal Error ",err)
		}
	}else{
		fmt.Println("Error",err)
	}
}