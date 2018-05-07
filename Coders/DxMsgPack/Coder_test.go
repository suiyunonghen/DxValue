package DxMsgPack

import (
	"testing"
	"os"
	"fmt"
	"io/ioutil"
	"bytes"
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

func TestMsgpackCoder_Custom(t *testing.T)  {
	buf, err := ioutil.ReadFile("test.msgpack")
	if err != nil {
		fmt.Println("ReadFile Err:",err)
		return
	}
	var myclass *classes
	coder := NewDecoder(bytes.NewReader(buf))
	coder.OnStartMap = func(mapLen int, keyIsStr bool) (mapInterface interface{}) {
		myclass = new(classes)
		return myclass
	}
	coder.OnStartStrMapArray = func(arrLen int, Key string, mapInterface interface{}) (arrInterface interface{}) {
		mclass := mapInterface.(*classes)
		switch Key {
		case "老师列表":
			mclass.Teachers = make([]*teacher,arrLen)
			return mclass.Teachers
		case "学生列表":
			mclass.Students = make([]*student,arrLen)
			return mclass.Students
		}
		return nil
	}
	coder.OnParserArrElement = func(arrInterface interface{}, index int, v interface{}) {
		if teachers,ok := arrInterface.([]*teacher);ok{
			if v!=nil{
				teachers[index] = v.(*teacher)
			}else{
				teachers[index] = nil
			}
		}else {
			if v!=nil{
				arrInterface.([]*student)[index] = v.(*student)
			}else{
				arrInterface.([]*student)[index] = nil
			}

		}
	}

	coder.OnParserArrObject = func(arrInterface interface{}, index int) (object interface{}) {
		if mtechers,ok := arrInterface.([]*teacher);ok{
			mtechers[index] = new(teacher)
			object = mtechers[index]
			return
		}else{
			students := arrInterface.([]*student)
			students[index] = new(student)
			object = students[index]
			return
		}
	}

	coder.OnParserStrMapKv= func(mapInterface interface{}, key string, v interface{}) {
		if mteacher,ok := mapInterface.(*teacher);ok{
			switch key {
			case "教师ID": mteacher.TeacherID = v.(string)
			case "课程": mteacher.Lesson = v.(string)
			}
		}else{

		}
	}
	coder.DecodeCustom()

	fmt.Println(myclass.Teachers[0])
}