package main

import (
	"hpv/app/task"
	"hpv/bootstrap"
)

func main() {
	app := bootstrap.NewApp()

	app.AddTask(bootstrap.Task{Handler: task.SendMess})
	app.AddTask(bootstrap.Task{Handler: task.RefreshToken})
	app.AddTask(bootstrap.Task{Handler: task.RunCorn})

	app.Run()
}
