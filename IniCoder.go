package DxValue

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/suiyunonghen/DxCommonLib"
	"io"
	"unsafe"
)

type decoderReader interface {
	io.Reader
	ReadBytes(delim byte) ([]byte, error)
}

type  DxIniDecoder struct {
	r			io.Reader
	readCM		DxCommonLib.FileCodeMode
}



func (decoder *DxIniDecoder)Decode(record *DxRecord)(error)  {
	var tmpr decoderReader
	reader,ok := decoder.r.(*bufio.Reader)
	if ok{
		tmpr = reader
	}else{
		reader,ok := decoder.r.(*bytes.Buffer)
		if ok{
			tmpr = reader
		}
	}
	if tmpr == nil{
		tmpr = bufio.NewReader(decoder.r)
	}
	curSection := ""
	var curRecord *DxRecord
	for{
		line,err := tmpr.ReadBytes('\n')
		if err == io.EOF{
			break
		}
		if err == nil {
			linelen := len(line)
			if linelen >= 2{
				if line[linelen-2] == '\r'{
					line = line[:linelen - 2]
					linelen = len(line)
				}else if line[linelen - 1] == '\n'{
					line = line[:linelen-1]
					linelen = len(line)
				}
			}
			if linelen>0{
				ifcommit := false
				idx := 0
				for ;idx<linelen;idx++{
					if !IsSpace(line[idx]){
						ifcommit = line[idx] == ';'
						break
					}
				}
				if ifcommit{
					continue
				}
				//判断是否是section
				//判断是否是一个有效的KV结构
				startSectionIdx := -1
				section := ""
				parserFor:
				for i := idx;i<linelen;i++{
					if IsSpace(line[i]){
						continue
					}
					switch line[i] {
					case '[':
						if i == idx{
							startSectionIdx = i+1
						}
					case ']':
						if startSectionIdx == -1{
							return errors.New(fmt.Sprintf("无效的ini格式，%d",i))
						}
						endidx := i
						for ;i>startSectionIdx;i--{
							if !IsSpace(line[endidx]){
								break
							}
						}
						section = DxCommonLib.FastByte2String(line[startSectionIdx:endidx])
					case ';':
						if section == ""{
							return errors.New(fmt.Sprintf("无效的ini格式，%d",i))
						}
					case '=':
						if curSection != ""{
							k := ""
							switch decoder.readCM {
							case DxCommonLib.File_Code_Utf8:
								k = DxCommonLib.FastByte2String(line[idx:i])
							case DxCommonLib.File_Code_GBK,DxCommonLib.File_Code_Unknown:
								if tmpbytes, err := DxCommonLib.GBK2Utf8(line[idx:i]); err == nil {
									k = DxCommonLib.FastByte2String(tmpbytes)
								}else{
									k = DxCommonLib.FastByte2String(line[idx:i])
								}
							case DxCommonLib.File_Code_Utf16LE,DxCommonLib.File_Code_Utf16BE:
								k = DxCommonLib.UTF16Byte2string(line[idx:i],decoder.readCM == DxCommonLib.File_Code_Utf16BE)
							}
							v := line[i+1:]
							for vcommitidx := i+1;vcommitidx<linelen;vcommitidx++{
								if !IsSpace(line[vcommitidx]){
									if line[vcommitidx]==';'{
										v = line[i+1:vcommitidx]
										break
									}
								}
							}
							bt := v
							switch decoder.readCM {
							case DxCommonLib.File_Code_Utf16LE,DxCommonLib.File_Code_Utf16BE:
								bt = ([]byte)(DxCommonLib.UTF16Byte2string(line[idx:i],decoder.readCM == DxCommonLib.File_Code_Utf16BE))
							case DxCommonLib.File_Code_GBK,DxCommonLib.File_Code_Unknown:
								if tmpbytes, err := DxCommonLib.GBK2Utf8(line[idx:i]); err == nil {
									bt = tmpbytes
								}
							}
							var dv DxValue
							_,err := dv.JsonParserFromByte(bt,false,false)
							if err!=nil{
								curRecord.SetString(k,DxCommonLib.FastByte2String(bt))
							}else{
								switch dv.fValue.fValueType {
								case DVT_Record:
									curRecord.SetRecordValue(k,(*DxRecord)(unsafe.Pointer(dv.fValue)))
								case DVT_Array:
									curRecord.SetArray(k,(*DxArray)(unsafe.Pointer(dv.fValue)))
								}
							}
						}
						break parserFor
					default:
						if section != ""{
							return errors.New(fmt.Sprintf("无效的ini格式，%d",i))
						}
					}
				}
				if section != ""{
					curRecord = record.AsRecord(section)
					if curRecord == nil{
						curRecord = record.NewRecord(section)
					}
					curSection = section
				}
			}
		}else{
			return err
		}
	}
	return nil
}

func NewIniDecoder(reader io.Reader,encodeMode DxCommonLib.FileCodeMode)*DxIniDecoder  {
	return &DxIniDecoder{r:reader,readCM:encodeMode}
}