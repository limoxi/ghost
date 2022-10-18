package ghost

import (
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"reflect"
)

type ContextObject struct {
	ctx context.Context
}

func (c *ContextObject) SetCtx(ctx context.Context) {
	c.ctx = ctx
}

func (c *ContextObject) GetCtx() context.Context {
	return c.ctx
}

func (c *ContextObject) GetGinCtx() *gin.Context {
	ictx := c.ctx.Value("ginCtx")
	if ictx == nil {
		return nil
	}
	return ictx.(*gin.Context)
}

// GetDB 尝试从context中获取db(事务), 没有则返回默认
func (c *ContextObject) GetDB() *gorm.DB {
	if c.ctx != nil {
		return GetDBFromCtx(c.ctx)
	} else {
		return GetDB()
	}
}

// DomainModel 领域模型
type DomainModel struct {
	ContextObject
}

// DomainService 领域服务
type DomainService struct {
	ContextObject
}

func (this *DomainModel) handleEmbedStruct(stField reflect.StructField, svField reflect.Value, dv reflect.Value) {
	for i := 0; i < stField.Type.NumField(); i++ {
		field := stField.Type.Field(i)
		fieldName := field.Name
		dvField := dv.FieldByName(fieldName)
		if field.Type.Kind() == reflect.Struct && dvField.Kind() != reflect.Struct {
			this.handleEmbedStruct(field, svField.FieldByName(fieldName), dv)
		} else {
			if dvField.CanSet() {
				dvField.Set(svField.Field(i))
			}
		}
	}
}

// NewFromDbModel
// 使用反射机制将dbModel中的field值复制到domainObject中
func (this *DomainModel) NewFromDbModel(do interface{}, dbModel interface{}) {
	siType := reflect.TypeOf(dbModel)
	siValue := reflect.ValueOf(dbModel)
	if siType.Kind() == reflect.Ptr {
		siType = siType.Elem()
		siValue = siValue.Elem()
	}
	diValue := reflect.ValueOf(do).Elem()
	for i := 0; i < siType.NumField(); i++ {
		field := siType.Field(i)
		fieldName := field.Name
		diField := diValue.FieldByName(fieldName)
		if field.Type.Kind() == reflect.Struct && diField.Kind() != reflect.Struct {
			this.handleEmbedStruct(field, siValue.FieldByName(fieldName), diValue)
		} else {
			if diField.CanSet() {
				diField.Set(siValue.Field(i))
			}
		}
	}
}

type BaseDomainRepository struct {
	ContextObject
	Paginator *Paginator
}

func (this *BaseDomainRepository) GetPaginator() *Paginator {
	return this.Paginator
}

func (this *BaseDomainRepository) SetPaginator(paginator *Paginator) {
	this.Paginator = paginator
}

func (this *BaseDomainRepository) SetPage(curPage, pageSize int) {
	paginator := NewPaginator(curPage, pageSize)
	this.SetPaginator(paginator)
}
