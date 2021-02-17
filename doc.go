package ghost

import (
	"bytes"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"text/template"
)

var docTmpl = `
# {{.ServiceName}} API文档 - 基于Ghost框架

[[toc]]

{{range .Resources}}
## 资源名

- &rq;{{.Name}}&rq;

### 资源描述

- &rq;{{.Desc}}&rq;

{{range .Methods}}
### 请求方法

- &rq;{{.Name}}&rq;

### 请求参数

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
{{range .RequestParams}}
|{{.Field}} |{{.Required}}  |{{.FieldType}} |{{.FieldDesc}}   |
{{end}}

### 响应数据

|参数名|类型|说明|
|:----    |:----- |-----   |
{{range .RespParams}}
|{{.Field}} |{{.FieldType}} |{{.FieldDesc}}   |
{{end}}

### 响应示例
&rq;&rq;&rq;json
{{.RespExample}}
&rq;&rq;&rq;
{{end}}

{{end}}
`

type structFieldInfo struct {
	Field string
	FieldType string
	FieldDesc string
	Required string
}

type methodInfo struct {
	Name string
	RequestParams []*structFieldInfo
	RespParams []*structFieldInfo
	RespExample string
}

type resourceInfo struct {
	Name string
	Desc string
	Methods []*methodInfo

}

type docParams struct {
	ServiceName string
	Resources []*resourceInfo
}

var methodParamNames = []string{"GetParams", "PutParams", "PostParams", "DeleteParams"}

func getStructInfo(sf reflect.StructField) []*structFieldInfo{
	return nil
}

// 生成markdown格式的接口文档
func GenApiDoc() {
	resources := collectResources()
	content := parseContent(resources)
	genMdFile(content)
	build()
}

func collectResources() []*resourceInfo {
	resources := make([]*resourceInfo, 0)
	for _, inst := range getAllApis() {
		instType := reflect.TypeOf(inst).Elem()
		methods := make([]*methodInfo, 0)
		resource := &resourceInfo{
			Name: inst.Resource(),
			Desc: instType.PkgPath() + instType.Name(),
		}

		for _, methodParamName := range methodParamNames {
			if field, exist := instType.FieldByName(methodParamName); exist{
				methods = append(methods, &methodInfo{
					Name: field.Name,
					RequestParams: getStructInfo(field),
					RespParams: getStructInfo(field),
					RespExample: "",
				})
			}
		}

		resource.Methods = methods
		resources = append(resources, resource)
	}
	return resources
}

func parseContent(resources []*resourceInfo) string{
	tmpl, err := template.New("doc").Parse(docTmpl)
	if err != nil{
		Error(err)
		panic(NewBusinessError("invalid_tmpl", "模板字符串错误"))
	}
	var contentBuffer bytes.Buffer
	err = tmpl.Execute(&contentBuffer, &docParams{
		ServiceName: "xxx",
		Resources: resources,
	})
	if err != nil{
		Error(err)
		panic(NewBusinessError("parse_tmpl_failed", "渲染模板失败"))
	}
	content := contentBuffer.String()
	return strings.ReplaceAll(content, "&rq;", "`")
}

func genMdFile(content string){
	f,err := os.Create( "README.md" )

	defer f.Close()

	if err != nil {
		Error(err)
		panic(NewSystemError("create_file:failed", "创建文件失败"))
	}

	_, err = f.Write([]byte(content))
	if err != nil {
		Error(err)
		panic(NewSystemError("write_file:failed", "写文件失败"))
	}
}

func build(){
	err := exec.Command("yarn", "run", "build").Run()
	if err != nil {
		Error(err)
		panic(NewSystemError("run_command:failed", "执行命令失败"))
	}
}