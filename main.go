package main

import (
	"hpv/app/task"
	"hpv/bootstrap"
)

func main() {
	app := bootstrap.NewApp()

	app.AddTask(bootstrap.Task{Values: []interface{}{app.Ctx}, Handler: task.DispatchMess})
	app.AddTask(bootstrap.Task{Handler: task.SendMess})
	app.AddTask(bootstrap.Task{Handler: task.RefreshToken})

	app.Run()
}
