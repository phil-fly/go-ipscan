package tools

import (
	"sync"
	"fmt"
	"log"
)

type Workdist struct {
	Host	string
}

const (
	taskload		    = 255
	tasknum			= 255
)
var wg sync.WaitGroup

func Task(ip string,debugLog *log.Logger){
	tasks := make(chan Workdist,taskload)
	wg.Add(tasknum)
	//创建chan消费者worker
	for gr:=1;gr<=tasknum;gr++ {
		go worker(tasks,debugLog)
	}

	//创建chan生产者
	for i:=1;i<256;i++ {
		host:=fmt.Sprintf("%s.%d",ip,i)
		task := Workdist{
			Host:host,
		}
		tasks <- task
	}
	close(tasks)
	wg.Wait()
}

func worker(tasks chan Workdist,debugLog *log.Logger){
	defer wg.Done()
	task,ok := <- tasks
	if !ok {
		return
	}
	host := task.Host
	ping(host,debugLog)
}

