package ReadJson

import (
	"encoding/json"
	_ "fmt"
	"io/ioutil"
)

type Config struct {
	SourceCategoryLink string
	PageNumber         int
	AllTime            string
	CategoryId         int
	PostUserId         int
	HostName           string
	MysqlUser          string
	MysqlPass          string
	MysqlPort          string
	MysqlHost          string
	MysqlDbName        string
}

//func main() {
//	JsonParse := NewJsonStruct()
//	v := Config{}
//	JsonParse.Load("./config.json", &v)
//	fmt.Println(v.SourceCategoryLink)
//	fmt.Println(v.PageNumber)
//	fmt.Println(v.AllTime)
//	fmt.Println(v.CategoryId)
//}

type JsonStruct struct {
}

func NewJsonStruct() *JsonStruct {
	return &JsonStruct{}
}
func (jst *JsonStruct) Load(filename string, v interface{}) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("config.json 缺失,请在此目录下面增添该文件")
		return
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		return
	}
}
