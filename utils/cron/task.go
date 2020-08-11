package cron

import (
	"context"
	"errors"
	"github.com/jinzhu/gorm"
	"math"
)

type taskInterface interface {
	Run(*TaskContext)
	GetName() string
	IsEnableTx() bool
}

type CronTask struct {
	name string
}

func (t *CronTask) Run(taskContext *TaskContext){
	panic(errors.New("not implemented"))
}

func (t *CronTask) GetName() string{
	return t.name
}

func (t *CronTask) SetName(name string) {
	t.name = name
}

func (t *CronTask) IsEnableTx() bool{
	return true
}

func NewCronTask(name string) CronTask{
	t := CronTask{name:name}
	return t
}

type pipeInterface interface {
	AddData(data interface{}) error
	GetData() interface{}
	GetCap() int
	GetConsumerCount() int
	RunConsumer(data interface{}, taskCtx *TaskContext)
	EnableParallel() bool
}

type Pipe struct{
	ch chan interface{}
	chCap int
}

func (p *Pipe) GetData() interface{}{
	return <- p.ch
}

func (p *Pipe) AddData(data interface{}) error{
	select {
	case p.ch <- data:
	default:
		return errors.New("channel is full")
	}
	return nil
}

func (p *Pipe) GetCap() int{
	return p.chCap
}

// GetConsumerCount 消费者数量
// 默认为通道容量十分之一
func (p *Pipe) GetConsumerCount() int{
	return int(math.Ceil(float64(p.GetCap())/10))
}

func (p *Pipe) RunConsumer() error{
	return errors.New("RunConsumer not implemented")
}

// EnableParallel 启用并行，默认启用
func (p *Pipe) EnableParallel() bool{
	return true
}

func NewPipe(chCap int) *Pipe{
	p := &Pipe{}
	p.chCap = chCap
	p.ch = make(chan interface{}, chCap)
	return p
}

type TaskContext struct{
	db *gorm.DB
	ctx context.Context
}

func (this *TaskContext) Init(ctx context.Context, db *gorm.DB){
	this.ctx = ctx
	this.db = db
}

func (this *TaskContext) GetDb() *gorm.DB{
	return this.db
}

func (this *TaskContext) GetCtx() context.Context{
	return this.ctx
}

func (this *TaskContext) SetCtx(ctx context.Context) {
	this.ctx = ctx
}