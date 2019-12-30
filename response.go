package ghost

type Response interface {
	GetCode() int
	GetData() interface{}
	GetDataType() string
}

type JsonResponse struct {
	code int
	data interface{}
}

func (this *JsonResponse) GetCode() int{
	return this.code
}
func (this *JsonResponse) GetData() interface{}{
	state := "error"
	if this.code == SERVICE_SUCCESS_CODE{
		state = "success"
	}
	return Map{
		"code": this.code,
		"state": state,
		"data": this.data,
	}
}
func (this *JsonResponse) GetDataType() string{
	return "json"
}

func NewJsonResponse(data interface{}) *JsonResponse{
	r := new(JsonResponse)
	r.code = SERVICE_SUCCESS_CODE
	r.data = data
	return r
}
func NewErrorJsonResponse(errCode string, args ...string) *JsonResponse{
	r := new(JsonResponse)
	r.code = SERVICE_BUSINESS_ERROR_CODE
	d := map[string]interface{}{
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

type RawResponse struct {
	code int
	rawString string
}

func (this *RawResponse) GetCode() int{
	return this.code
}
func (this *RawResponse) GetData() interface{}{
	return this.rawString
}
func (this *RawResponse) GetDataType() string{
	return "string"
}

func NewRawResponse(msg string) *RawResponse{
	r := new(RawResponse)
	r.code = SERVICE_SUCCESS_CODE
	r.rawString = msg
	return r
}

func NewErrorRawResponse(msg string) *RawResponse{
	r := new(RawResponse)
	r.code = SERVICE_BUSINESS_ERROR_CODE
	r.rawString = msg
	return r
}