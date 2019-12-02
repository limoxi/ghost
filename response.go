package ghost

type Response interface {
	GetData() interface{}
	GetDataType() string
}

type JsonResponse struct {
	code int
	data interface{}
}
func (this *JsonResponse) GetData() interface{}{
	return this.data
}
func (this *JsonResponse) GetDataType() string{
	return "json"
}

func NewJsonResponse(data interface{}) *JsonResponse{
	r := new(JsonResponse)
	r.code = 200
	r.data = data
	return r
}
func NewErrorJsonResponse(code int, errCode string, args ...string) *JsonResponse{
	r := new(JsonResponse)
	r.code = 200
	d := map[string]interface{}{
		"code": code,
		"errCode": errCode,
	}
	l := len(args)
	if l >= 1{
		d["errMsg"] = args[0]
	}
	if l >= 2{
		d["errStack"] = args[1]
	}
	r.data = d
	return r
}