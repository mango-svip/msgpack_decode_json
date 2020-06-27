package main

import (
	"encoding/json"
	"flag"
	"fmt"
	p "github.com/tidwall/pretty"
	"github.com/ugorji/go/codec"
	"io/ioutil"
	"os"
)

var (
	sourceFile string
	outFile    string
	h          bool
	pretty     bool
)

func init() {
	flag.BoolVar(&h, "h", false, "help")
	flag.StringVar(&sourceFile, "s", "", "需要decode的json文件名")
	flag.StringVar(&outFile, "o", "tmp.json", "输出的文件名")
	flag.BoolVar(&pretty, "p", false, "是否美化输出格式")

	flag.Usage = usage

}

func main() {
	flag.Parse()

	if h {
		flag.Usage()
		return
	}

	arg := flag.Arg(0)

	if sourceFile == "" && arg == "" {

		fmt.Println("需要decode的json文件名不能为空")
		return
	} else if sourceFile == ""{
		sourceFile = arg
		fmt.Println(sourceFile)
	}

	bytes, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		panic(err)
	}
	var m = make(map[string]interface{})

	codec.NewDecoderBytes(bytes, new(codec.MsgpackHandle)).Decode(&m)

	tmp := mapHandler(m)

	marshal, err := json.Marshal(tmp)
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(outFile, p.Pretty(marshal), os.ModePerm)

}

func usage() {
	fmt.Fprintf(os.Stderr, `codec version: ugorji/go/codec/v1.1.7
Usage: msgPackDecode [-h] [-s filename] [-o filename]

Options:
`)
	flag.PrintDefaults()
}

func mapHandler(param interface{}) map[string]interface{} {
	tmp := make(map[string]interface{})
	switch param.(type) {
	case nil:
		return tmp

	case map[string]interface{}:
		for k, v := range param.(map[string]interface{}) {
			switch v.(type) {
			case map[interface{}]interface{}:
				tmp[k] = mapHandler(v)
			case []uint8:
				tmp[k] = B2S(v.([]uint8))

			default:
				tmp[k] = v
			}
		}
	case map[interface{}]interface{}:
		for k, v := range param.(map[interface{}]interface{}) {
			switch k.(type) {
			case string:
				switch v.(type) {
				case map[interface{}]interface{}:
					tmp[k.(string)] = mapHandler(v)
					continue
				case []interface{}:
					//tmpArray := make([]interface{}, 0)
					//for _, n := range v.([]interface{}) {
					//	newArray := make([]interface{}, 0)
					//	switch n.(type) {
					//	case []interface{}:
					//		for _, m := range n.([]interface{}) {
					//			switch m.(type) {
					//			case []uint8:
					//				newArray = append(newArray, B2S(m.([]uint8)))
					//			default:
					//				newArray = append(newArray, m)
					//			}
					//		}
					//	default:
					//		fmt.Println("unknown type", n)
					//	}
					//	tmpArray = append(tmpArray, newArray)
					//}
					//tmp[k.(string)] = tmpArray
					tmp[k.(string)] = B2SinArraySlice(v.([]interface{}))
				default:
					tmp[k.(string)] = v
				}
			default:
				continue
			}
		}
	}
	return tmp
}

func B2S(bs []uint8) string {
	ba := []byte{}
	for _, b := range bs {
		ba = append(ba, byte(b))
	}
	return string(ba)
}

func B2SinArraySlice(array []interface{}) []interface{}{
	tmp := make([]interface{}, 0)
	for _, v := range array {
		switch v.(type) {
		case []interface{}:
			tmp = append(tmp, B2SinArraySlice(v.([]interface{})))
		case []uint8:
			tmp = append(tmp, B2S(v.([]uint8)))
		default:
			tmp = append(tmp, v)

		}
	}
	return tmp
}
