package bootstrap

import (
	"hpv/bootstrap/context"
	"os"
	"sync"
)

type ConsoleApp struct {
	Ctx   *context.Context
	Tasks []Task
}

type Task struct {
	Values  []interface{}
	Handler TaskHandler
}

type TaskHandler func(values ...interface{})

// initCommon _
func (col *ConsoleApp) initCommon() {
	col.Ctx = context.NewContext()

	// root dir
	rootDir, _ := os.Getwd()
	col.Ctx.Set("root_dir", rootDir)

	// log dir
	logDir := rootDir + "/logs"
	col.Ctx.Set("log_dir", logDir)
}

func (col *ConsoleApp) AddTask(task Task) {
	col.Tasks = append(col.Tasks, task)
}

func (col *ConsoleApp) Run() {
	var wg sync.WaitGroup

	wg.Add(len(col.Tasks))
	for _, task := range col.Tasks {
		task.Values = append(task.Values, &wg)
		go task.Handler(task.Values...)
	}
	wg.Wait()
}
