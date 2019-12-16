package ghost

import (
	"context"
	"github.com/gin-gonic/gin"
	"reflect"
)

type ContextObject struct{
	ctx context.Context
}

func (c *ContextObject) SetCtx(ctx context.Context){
	c.ctx = ctx
}

func (c *ContextObject) GetCtx() context.Context{
	return c.ctx
}

// DomainObject 领域对象(可以表示聚合根、聚合、实体、值对象和领域服务)
type DomainObject struct{
	ContextObject
}

// DomainModel 领域模型（表示实体）
type DomainModel struct{
	DomainObject
}

func (this *DomainModel) GetDbModel() interface{}{
	dbModel, ok := this.ctx.(*gin.Context).Get("dbModel")
	if ok{
		return dbModel
	}else{
		return nil
	}
}

func (this *DomainModel) handleEmbedStruct(stField reflect.StructField, svField reflect.Value, dv reflect.Value){
	for i:=0; i<stField.Type.NumField(); i++{
		field := stField.Type.Field(i)
		fieldName := field.Name
		dvField := dv.FieldByName(fieldName)
		if field.Type.Kind() == reflect.Struct && dvField.Kind() != reflect.Struct{
			this.handleEmbedStruct(field, svField.FieldByName(fieldName), dv)
		}else{
			if dvField.CanSet(){
				dvField.Set(svField.Field(i))
			}
		}
	}
}

// NewFromDbModel
// 使用反射机制将dbModel中的field值复制到domainObject中
func (this *DomainModel) NewFromDbModel(do interface{}, dbModel interface{}){
	siType := reflect.TypeOf(dbModel)
	siValue := reflect.ValueOf(dbModel)
	if siType.Kind() == reflect.Ptr{
		siType = siType.Elem()
		siValue = siValue.Elem()
	}
	diValue := reflect.ValueOf(do).Elem()
	for i:=0; i<siType.NumField(); i++{
		field := siType.Field(i)
		fieldName := field.Name
		if field.Type.Kind() == reflect.Struct{
			this.handleEmbedStruct(field, siValue.FieldByName(fieldName), diValue)
		}else{
			diField := diValue.FieldByName(fieldName)
			if diField.CanSet(){
				diField.Set(siValue.Field(i))
			}
		}
	}
	(this.ctx).(*gin.Context).Set("dbModel", dbModel)
}