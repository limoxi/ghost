package ghost

import (
	"github.com/gin-gonic/gin"
	"github.com/limoxi/ghost/utils"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var Config *config

type config struct {
	GMap
	Mode string
}

// setDefaultConfigData 默认配置
func (this *config) setDefaultConfigData() {
	this.GMap = NewEmptyGMap()
}

// findConfDir 查找配置文件夹conf
// 1、在当前工作目录下查找
// 2、逐级向上查找，查找深度默认3
func findConfDir(p string, args ...int) string {
	count := 0
	switch len(args) {
	case 1:
		count = args[0]
	}
	if count >= 3 {
		return ""
	}
	ds, err := ioutil.ReadDir(p)
	if err != nil {
		panic(err)
	}
	for _, d := range ds {
		if d.Name() == "conf" {
			return p
		}
	}
	p = path.Join(p, "..")
	count += 1
	return findConfDir(p, count)
}

// load
// 固定加载项目路径下conf文件夹中的yaml配置文件
func (this *config) load() {
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	workPath = findConfDir(workPath)
	if workPath == "" {
		this.setDefaultConfigData()
		return
	}
	pre := "prod"
	if this.Mode == gin.DebugMode {
		pre = "dev"
	}
	filename := pre + ".conf.yaml"
	confPath := filepath.Join(workPath, "conf", filename)
	if utils.FileExist(confPath) {
		if err := this.parseYamlFile(confPath); err == nil {
			return
		}
	} else {
		confPath = filepath.Join(workPath, "conf", "conf.yaml")
		if utils.FileExist(confPath) {
			if err := this.parseYamlFile(confPath); err == nil {
				return
			}
		}
	}
	this.setDefaultConfigData()
}

func (this *config) parseYamlFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		Warn(err)
		return err
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		Warn(err)
		return err
	}
	var data Map
	err = yaml.Unmarshal(content, &data)
	if err != nil {
		Warn(err)
		return err
	}
	this.GMap = NewGMapFromData(parseEnvArgs(data))
	return nil
}

// parseEnvArgs 从字符串中解析出环境变量
func parseEnvArgs(data Map) Map {
	for k, v := range data {
		switch v.(type) {
		case string:
			data[k] = parseEnvFromString(v.(string))
		case Map:
			ns := k + "."
			for ik, iv := range parseEnvArgs(v.(Map)) {
				fullKey := ns + ik
				data[fullKey] = iv
			}
		}
	}
	return data
}

func parseEnvFromString(str string) string {
	str = strings.Replace(str, " ", "", -1)
	if !strings.HasPrefix(str, "${") {
		return str
	}
	envV := ""
	defaultV := ""
	sps := strings.Split(str, "||")
	if len(sps) == 2 {
		defaultV = sps[1]
	}
	sqIndex := strings.Index(str, "}")
	envKey := str[2:sqIndex]
	if envKey != "" {
		envV = os.Getenv(envKey)
		if envV == "" {
			envV = defaultV
		}
	}
	return envV
}

func LoadConfigFromFile(path string) *config {
	c := new(config)
	c.parseYamlFile(path)
	return c
}

func init() {
	Config = new(config)
	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = gin.DebugMode
	}
	Infof("loading config in %s mode...", mode)
	Config.Mode = mode
	Config.load()
}
