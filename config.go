package ghost

import (
	"encoding/json"
	"github.com/limoxi/ghost/utils"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var Config *config

const (
	DEV_MODE = "dev"
	TEST_MODE = "test"
	PROD_MODE = "prod"
)

type config struct{
	Map
	Mode string
}

// setDefaultConfigData 默认配置
func (this *config) setDefaultConfigData(){
	this.Map = NewEmptyMap()
}

// findConfDir 查找配置文件夹conf
// 1、在当前工作目录下查找
// 2、逐级向上查找，查找深度默认3
func findConfDir(p string, args ...int) string{
	count := 0
	switch len(args) {
	case 1:
		count = args[0]
	}
	if count >= 3{
		return ""
	}
	ds, err := ioutil.ReadDir(p)
	if err != nil{
		Panic(err)
	}
	for _, d := range ds{
		if d.Name() == "conf"{
			return p
		}
	}
	p = path.Join(p, "..")
	count += 1
	return findConfDir(p, count)
}

// load
// 固定加载项目路径下conf文件夹中的json配置文件
// 如果找不到，则加载自带配置文件
func (this *config) load(){
	workPath, err := os.Getwd()
	if err != nil{
		panic(err)
	}
	workPath = findConfDir(workPath)
	if workPath == ""{
		this.setDefaultConfigData()
		return
	}
	filename := this.Mode + ".conf.json"
	confPath := filepath.Join(workPath, "conf", filename)
	if utils.FileExist(confPath){
		if err := this.parseJsonFile(confPath); err == nil{
			return
		}
	}else{
		confPath = filepath.Join(workPath, "conf", "conf.json")
		if utils.FileExist(confPath) {
			if err := this.parseJsonFile(confPath); err == nil {
				return
			}
		}
	}
	this.setDefaultConfigData()
}

func (this *config) parseJsonFile(filename string) error{
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	var data Map
	err = json.Unmarshal(content, &data)
	if err != nil{
		return err
	}
	this.Map = parseEnvArgs(data)
	return nil
}

// parseEnvArgs 从字符串中解析出环境变量
func parseEnvArgs(data Map) Map{
	for k, v := range data{
		switch v.(type) {
		case string:
			data[k] = parseEnvFromString(v.(string))
		case map[string]interface{}:
			ns := k + "."
			for ik, iv := range parseEnvArgs(v.(map[string]interface{})){
				fullKey := ns + ik
				data[fullKey] = iv
			}
		}
	}
	return data
}

func parseEnvFromString(str string) string{
	str = strings.Replace(str, " ", "", -1)
	if !strings.HasPrefix(str, "${"){
		return str
	}
	envV := ""
	defaultV := ""
	sps := strings.Split(str, "||")
	if len(sps) == 2{
		defaultV = sps[1]
	}
	sqIndex := strings.Index(str, "}")
	envKey := str[2: sqIndex]
	if envKey != ""{
		envV = os.Getenv(envKey)
		if envV == ""{
			envV = defaultV
		}
	}
	return envV
}

func (this *config) GetArray(key string) []interface{}{
	return this.Get(key).([]interface{})
}

func (this *config) GetMap(key string) Map{
	return this.Get(key).(map[string]interface{})
}

func init(){
	Config = new(config)
	mode := os.Getenv("GHOST_MODE")
	if mode == ""{
		mode = DEV_MODE
	}
	log.Printf("loding config in %s mode...", mode)
	Config.Mode = mode
	Config.load()
}