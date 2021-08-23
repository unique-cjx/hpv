package main

import (
	"hpv/app/task"
	"hpv/bootstrap"
)

func main() {
	app := bootstrap.NewApp()
	app.AddTask(bootstrap.Task{Values: []interface{}{app.Ctx}, Handler: task.GetActiveRegions})
	app.AddTask(bootstrap.Task{Handler: task.SubscribeDepart})
	app.Run()
}
