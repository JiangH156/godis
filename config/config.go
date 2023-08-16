package config

import (
	"bufio"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type PropertyHolder struct {
	Bind           string   `cfg:"bind"`           //绑定地址
	Port           int      `cfg:"port"`           //监听端口
	AppendOnly     bool     `cfg:"appendOnly"`     //是否启用AOF（Append-Only File）持久化
	AppendFilename string   `cfg:"appendFilename"` //AOF文件的文件名
	MaxClients     int      `cfg:"maxClients"`     //最大客户端数量
	Peers          []string `cfg:"peers"`          //其他节点的地址列表
	Self           string   `cfg:"self"`           //本身的地址
}

var Properties *PropertyHolder

func init() {
	Properties = &PropertyHolder{
		Bind:       "127.0.0.1",
		Port:       6379,
		AppendOnly: false, //默认关闭AOF持久化
	}
}

func loadConfig(filename string) *PropertyHolder {
	config := Properties
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
		return config
	}
	//读取文件，获取kv对
	kvMap := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || (len(line) > 0 && line[0] == '#') { // 改行为空或为注释
			continue
		}
		index := strings.IndexAny(line, " ")
		if index > 0 && index < len(line)-1 {
			key := line[0 : index-1]
			value := strings.Trim(line[index+1:], " ")
			kvMap[strings.ToLower(key)] = value
		}
	}
	t := reflect.TypeOf(Properties) //Properties是指针
	v := reflect.ValueOf(Properties)
	n := t.Elem().NumField()
	for i := 0; i < n; i++ {
		field := t.Elem().Field(i)
		fieldVal := v.Elem().Field(i)
		key, ok := field.Tag.Lookup("cfg")
		if !ok {
			key = field.Name
		}
		value, ok := kvMap[strings.ToLower(key)]
		if ok { // 存在cfg tag
			switch field.Type.Kind() {
			case reflect.String:
				fieldVal.SetString(value)
			case reflect.Int:
				intValue, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					log.Fatalln(err)
					return config
				}
				fieldVal.SetInt(intValue)
			case reflect.Bool:
				boolValue, err := strconv.ParseBool(value)
				if err != nil {
					log.Fatalln(err)
					return config
				}
				fieldVal.SetBool(boolValue)
			case reflect.Slice:
				if field.Type.Elem().Kind() == reflect.String {
					sliceValue := strings.Split(value, ",")
					fieldVal.Set(reflect.ValueOf(sliceValue))
				}
			}
		}
	}
	return config
}

func SetupConfig(filename string) {
	Properties = loadConfig(filename)
}
