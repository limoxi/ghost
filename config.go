package ghost

import (
	"encoding/json"
	"github.com/limoxi/ghost/utils"
	"io/ioutil"
	"os"
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
	Mode string
	data Map
}

// setDefaultConfigData 默认配置
func (this *config) setDefaultConfigData(){
	this.data = NewEmptyMap()
}

// load
// 固定加载项目路径下conf文件夹中的json配置文件
// 如果找不到，则加载自带配置文件
func (this *config) load(){
	workPath, err := os.Getwd()
	if err != nil{
		panic(err)
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
	var data map[string]interface{}
	err = json.Unmarshal(content, &data)
	if err != nil{
		return err
	}
	this.data = parseEnvArgs(data)
	return nil
}

// parseEnvArgs 从字符串中解析出环境变量
func parseEnvArgs(data map[string]interface{}) map[string]interface{}{
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

func (this *config) get(key string) interface{}{
	if this.data == nil{
		this.load()
	}
	return this.data.Get(key)
}

func (this *config) GetString(key string, defaultArgs ...string) string{
	return this.data.GetString(key, defaultArgs...)
}

func (this *config) GetBool(key string, defaultArgs ...bool) bool{
	return this.data.GetBool(key, defaultArgs...)
}

func (this *config) GetInt(key string, defaultArgs ...int) int{
	return this.data.GetInt(key, defaultArgs...)
}

func (this *config) GetFloat(key string, defaultArgs ...float64) float64{
	return this.data.GetFloat(key, defaultArgs...)
}

func (this *config) GetArray(key string) []interface{}{
	return this.data.GetArray(key)
}

func (this *config) GetMap(key string) Map{
	return this.data.GetMap(key)
}

func init(){
	Config = new(config)
	mode := os.Getenv("GHOST_MODE")
	if mode == ""{
		mode = DEV_MODE
	}
	Config.Mode = mode
}