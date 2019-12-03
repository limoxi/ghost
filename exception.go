package ghost

var ErrLevel = map[string]uint8{
	"BUSINESS_ERR": 1,
	"SYSTEM_ERR": 2,
}

type BaseError struct{
	code int
	level uint8
	ErrCode string
	ErrMsg string
}

func (this *BaseError) GetCode() int{
	return this.code
}
func (this *BaseError) IsBusinessError() bool{
	return this.level == ErrLevel["BUSINESS_ERR"]
}
func (this *BaseError) IsSystemError() bool{
	return this.level == ErrLevel["SYSTEM_ERR"]
}
func DefaultError(errMsg string) *BaseError{
	return NewBusinessError(errMsg)
}
func NewBusinessError(args ...string) *BaseError{
	inst := new(BaseError)
	inst.code = 500
	inst.level = ErrLevel["BUSINESS_ERR"]
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
	inst.code = 533
	inst.level = ErrLevel["SYSTEM_ERR"]
	switch len(args) {
	case 1:
		inst.ErrMsg = args[0]
	case 2:
		inst.ErrCode = args[0]
		inst.ErrMsg = args[1]
	}
	return inst
}