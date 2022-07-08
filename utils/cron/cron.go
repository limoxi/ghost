package cron

import (
	"context"
	"fmt"
	"github.com/limoxi/ghost"
	"runtime/debug"
	"time"
)

type cron struct {
	name     string
	spec     string
	taskFunc TaskFunc
}

var name2task = make(map[string]*cron)

func newTaskCtx() *TaskContext {
	inst := new(TaskContext)
	ctx := context.Background()
	db := ghost.GetDB()
	ctx = context.WithValue(ctx, "db", db)

	inst.Init(ctx, db)
	return inst
}

func taskWrapper(task taskInterface) TaskFunc {

	return func() error {
		taskCtx := newTaskCtx()
		db := taskCtx.GetDb()
		ctx := taskCtx.GetCtx()

		taskName := task.GetName()
		startTime := time.Now()
		ghost.Info(fmt.Sprintf("[%s] run...", taskName))
		if db != nil && task.IsEnableTx() {
			tx := db.Begin()
			if err := tx.Error; err != nil {
				panic(err)
			}
			ctx = context.WithValue(ctx, "db", tx)
			defer ghost.RecoverFromCronTaskPanic(ctx)
			taskCtx.SetCtx(ctx)

			task.Run(taskCtx)
			if err := tx.Commit().Error; err != nil {
				ghost.Error(err)
			}
		} else {
			defer ghost.RecoverFromCronTaskPanic(ctx)
			task.Run(taskCtx)
		}
		dur := time.Since(startTime)
		ghost.Info(fmt.Sprintf("[%s] done, cost %g s", taskName, dur.Seconds()))
		return nil
	}
}

func fetchData(pi pipeInterface) {
	taskName := pi.(taskInterface).GetName()
	go func() {
		defer func() {
			if err := recover(); err != nil {
				ghost.Warn(string(debug.Stack()))
				fetchData(pi)
				errMsg := err.(error).Error()
				ghost.CaptureTaskErrorToSentry(context.Background(), errMsg)
			}
		}()
		for {
			data := pi.GetData()
			if data != nil {
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

func RegisterPipeTask(pi pipeInterface, spec string, runInRest ...bool) {
	task := RegisterTask(pi.(taskInterface), spec, runInRest...)
	if task != nil {
		if pi.EnableParallel() { // 并行模式下，开启通道容量十分之一的goroutine消费通道
			for i := pi.GetConsumerCount(); i > 0; i-- {
				fetchData(pi)
			}
		} else {
			fetchData(pi)
		}
	}
}

func RegisterTask(task taskInterface, spec string, runInRest ...bool) *cron {
	taskInRest := false
	switch len(runInRest) {
	case 1:
		taskInRest = runInRest[0]
	}
	if ghost.Config.GetString("run_mod") == "cron" || taskInRest {
		tname := task.GetName()
		wrappedFn := taskWrapper(task)
		cronTask := &cron{
			name:     tname,
			spec:     spec,
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
