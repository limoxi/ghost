package ghost

import "fmt"

const SERVICE_INNER_SUCCESS_CODE = 200
const SERVICE_SUCCESS_CODE = 200 // 业务成功
const SERVICE_BUSINESS_ERROR_CODE = 520 // 业务错误
const SERVICE_SYSTEM_ERROR_CODE = 530 // 系统错误

type BaseError struct{
	code int
	ErrCode string
	ErrMsg string
}

func (this *BaseError) GetCode() int{
	return this.code
}
func (this *BaseError) GetData() interface{}{
	return Map{
		"code": this.code,
		"state": "error",
		"data": Map{
			"errCode": this.ErrCode,
			"errMsg": this.ErrMsg,
		},
	}
}
func (this *BaseError) ToString() string{
	return fmt.Sprintf("code: %s, msg: %s", this.ErrCode, this.ErrMsg)
}
func (this *BaseError) IsBusinessError() bool{
	return this.code == SERVICE_BUSINESS_ERROR_CODE
}
func (this *BaseError) IsSystemError() bool{
	return this.code == SERVICE_SYSTEM_ERROR_CODE
}
func DefaultError(errMsg string) *BaseError{
	return NewBusinessError(errMsg)
}
func NewBusinessError(args ...string) *BaseError{
	inst := new(BaseError)
	inst.code = SERVICE_BUSINESS_ERROR_CODE
	switch len(args) {
	case 1:
		inst.ErrMsg = args[0]
	case 2:
		inst.ErrCode = args[0]
		inst.ErrMsg = args[1]
	}
	return inst
}

func NewSystemError(args ...string) *BaseError{
	inst := new(BaseError)
	inst.code = SERVICE_SYSTEM_ERROR_CODE
	switch len(args) {
	case 1:
		inst.ErrMsg = args[0]
	case 2:
		inst.ErrCode = args[0]
		inst.ErrMsg = args[1]
	}
	return inst
}