package ghost

import (
	"reflect"
)

var docTmpl = `
### 资源名

- &rq;{{.Resource}}&rq;

##### 资源描述

- &rq;{{.ResourceDesc}}&rq;

{{range .Methods}}
#### 请求方式

- &rq;{{.Name}}&rq;

#### 请求参数

|参数名|必选|类型|说明|
|:----    |:---|:----- |-----   |
{{range .RequestParams}}
|{{.Field}} |{{.Required}}  |{{.FieldType}} |{{.FieldDesc}}   |
{{end}}

#### 响应数据

|参数名|类型|说明|
|:----    |:----- |-----   |
{{range .RespParams}}
|{{.Field}} |{{.FieldType}} |{{.FieldDesc}}   |
{{end}}

#### 响应示例
&rq;&rq;&rq;json
{{.RespExample}}
&rq;&rq;&rq;
{{end}}
`

type docMethodParams struct {
	Field string
	FieldType string
	FieldDesc string
	Required string
}

type docMethods struct {
	Name string
	RequestParams []*docMethodParams
	RespParams []*docMethodParams
	RespExample string
}

type docParams struct {
	Resource string
	ResourceDesc string
	Methods []*docMethods

}

// 生成markdown格式的接口文档
func genApiDoc() {
	for _, inst := range registeredApis{
		instType := reflect.TypeOf(inst).Elem()
		Info(instType.Name(), inst.Resource(), "00000000000000", instType.PkgPath())
		//instVal := reflect.ValueOf(inst).Elem()

		for num := 0; num < instType.NumField(); num++{
			field := instType.Field(num)
			Info(field.Name, field.Tag, "1111111111111")
		}
	}
}