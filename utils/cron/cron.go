package cron

import (
	"context"
	"fmt"
	"github.com/limoxi/ghost"
	"runtime/debug"
	"time"
)

type cron struct {
	name string
	spec string
	taskFunc TaskFunc
}

var name2task = make(map[string]*cron)

func newTaskCtx() *TaskContext{
	inst := new(TaskContext)
	ctx := context.Background()
	o := ghost.GetDB()
	ctx = context.WithValue(ctx, "orm", o)

	inst.Init(ctx, o)
	return inst
}

func taskWrapper(task taskInterface) TaskFunc{

	return func() error{
		taskCtx := newTaskCtx()
		o := taskCtx.GetOrm()
		ctx := taskCtx.GetCtx()

		defer ghost.RecoverFromCronTaskPanic(ctx)
		var fnErr error
		taskName := task.GetName()
		startTime := time.Now()
		ghost.Info(fmt.Sprintf("[%s] run...", taskName))
		if o != nil && task.IsEnableTx(){
			o.Begin()
			fnErr = task.Run(taskCtx)
			o.Commit()
		}else{
			fnErr = task.Run(taskCtx)
		}
		dur := time.Since(startTime)
		ghost.Info(fmt.Sprintf("[%s] done, cost %g s", taskName, dur.Seconds()))
		return fnErr
	}
}

func fetchData(pi pipeInterface){
	taskName := pi.(taskInterface).GetName()
	go func(){
		defer func(){
			if err := recover(); err!=nil{
				ghost.Warn(string(debug.Stack()))
				fetchData(pi)
				errMsg := err.(error).Error()
				ghost.CaptureTaskErrorToSentry(context.Background(), errMsg)
			}
		}()
		for{
			data := pi.GetData()
			if data != nil{
				taskCtx := newTaskCtx()
				ghost.Info(fmt.Sprintf("[%s] consume data...", taskName))
				startTime := time.Now()
				pi.RunConsumer(data, taskCtx)
				dur := time.Since(startTime)
				ghost.Info(fmt.Sprintf("[%s] consume done, cost %g s !", taskName, dur.Seconds()))
			}
		}
	}()
}

func RegisterPipeTask(pi pipeInterface, spec string){
	task := RegisterTask(pi.(taskInterface), spec)
	if task != nil{
		if pi.EnableParallel(){ // 并行模式下，开启通道容量十分之一的goroutine消费通道
			for i := pi.GetConsumerCount(); i>0; i--{
				fetchData(pi)
			}
		}else{
			fetchData(pi)
		}
	}
}

func RegisterTask(task taskInterface, spec string, args ...bool) *cron {
	taskInRest := false
	switch len(args) {
	case 1:
		taskInRest = args[0]
	}
	if ghost.Config.GetString("run_mod") == "cron" || taskInRest{
		tname := task.GetName()
		wrappedFn := taskWrapper(task)
		cronTask := &cron{
			name: tname,
			spec: spec,
			taskFunc: wrappedFn,
		}
		name2task[tname] = cronTask
		
		return cronTask
	}
	return nil
}

func StartCronTasks() {
	for _, cronTask := range name2task {
		ghost.Info("[cron] create cron task ", cronTask.name, cronTask.spec)
		task := NewTask(cronTask.name, cronTask.spec, cronTask.taskFunc)
		AddTask(cronTask.name, task)
	}
	StartTask()
}

func StopCronTasks() {
	StopTask()
}
