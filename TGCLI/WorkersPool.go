package main

//https://github.com/wangyaofenghist/go-worker-base/blob/master/worker/workerManage.go
import (
	//	"sync/atomic"
	"time"

	tcp "NetworkFramework"
)

type Job func([]interface{})
type taskWork struct {
	Run       Job
	startBool bool
	params    []interface{}
}

var WorkMaxTask int
var WorkTaskPool chan taskWork
var WorkTaskReturn chan []interface{}

//启动任务
func (t *taskWork) start() {
	go func() {
		for {
			select {
			case funcRun := <-WorkTaskPool:
				if funcRun.startBool == true {
					funcRun.Run(funcRun.params)
				} else {
					//fmt.Println("task  stop!")

					tcp.Logger.Debug("task  stop!")

					return
				}
			case <-time.After(time.Millisecond * 1000):

				tcp.Logger.Debug("time out")
			}
		}
	}()
}

func (t *taskWork) stop() {

	tcp.Logger.Debug("task  stop!")
	t.startBool = false
}
func createTask() taskWork {
	var funcJob Job
	var paramSlice []interface{}
	return taskWork{funcJob, true, paramSlice}
}

//循环启动协程池
func StartPool(maxTask int) {
	WorkMaxTask = maxTask
	WorkTaskPool = make(chan taskWork, maxTask)
	WorkTaskReturn = make(chan []interface{}, maxTask)

	for i := 0; i < maxTask; i++ {
		var t = createTask()
		//tcp.Logger.Debug("start task: %d", i)
		t.start()
	}
}

//消费任务
func Dispatch(funcJob Job, params ...interface{}) {
	WorkTaskPool <- taskWork{funcJob, true, params}
}

//停止协程池
func StopPool() {
	var funcJob Job
	var paramSlice []interface{}
	for i := 0; i < WorkMaxTask; i++ {
		WorkTaskPool <- taskWork{funcJob, false, paramSlice}
	}
}

var workerNumDefault int = 50

var workerNumMax int = workerNumDefault * 2

type WorkPool struct {
	taskPool   chan taskWork
	workNum    int
	maxNum     int
	defaultNum int
	stopTopic  bool
	//tasks      int32
	//暂时没有用，考虑后期 作为冗余队列使用
	taskQue chan taskWork
}

//得到一个线程池并返回 句柄
func (p *WorkPool) InitPool(count int) {
	workerNumDefault = count
	//p.tasks = 0
	*p = WorkPool{defaultNum: workerNumDefault,
		maxNum: workerNumMax, stopTopic: false,
		taskPool: make(chan taskWork, workerNumDefault*2),
		taskQue:  nil}

	(p).start()
	go (p).workerRemoveConf()
}

//开始work
func (p *WorkPool) start() {
	for i := 0; i < p.defaultNum; i++ {

		//tcp.Logger.Critical("worker %d  ", i)
		p.workInit(i)

	}
}

//初始化 work池 后期应该考虑如何 自动 增减协程数，以达到最优
func (p *WorkPool) workInit(id int) {
	//tcp.Logger.Debug("start pool task worker id:%d", id)
	p.workNum++
	go func(idNum int) {
		defer func() {
			if err := recover(); err != nil {
				stacks := tcp.PanicTrace(4)
				tcp.Logger.Error("worker %d  exit panics: %v call:%v", idNum, err, string(stacks))
				p.workNum--
				p.workInit(idNum)

			}

		}()

		for {
			select {
			case task := <-p.taskPool:
				if task.startBool == true && task.Run != nil {
					//fmt.Print("this is pool ", idNum, "---")
					//atomic.AddInt32(&p.tasks, 1)
					//start := time.Now()
					task.Run(task.params)
					//cost := time.Since(start)
					//atomic.AddInt32(&p.tasks, -1)
					//tcp.Logger.Emergency("Task Done By Worker Id:%d Cost Time:%s TotalworkersNum:%d", idNum, cost.String(), p.workNum)
				}
				//单个结束任务
				if task.startBool == false {
					//fmt.Print("this is pool -- ", idNum, "---")
					return
				}
				//防止从channal 中读取数据超时
			case <-time.After(time.Millisecond * 1000):
				//fmt.Println("time out init")
				if p.stopTopic == true && len(p.taskPool) == 0 {
					//fmt.Println("topic=", p.stopTopic)
					//work数递减
					p.workNum--
					return
				}
				//从备用队列读取数据
			case queTask := <-p.taskQue:
				if queTask.startBool == true && queTask.Run != nil {
					//fmt.Print("this is que ", idNum, "---")
					queTask.Run(queTask.params)
				}
			}

		}
		tcp.Logger.Error("worker  exit %d", idNum)
	}(id)

}

//停止一个workPool
func (p *WorkPool) Stop() {
	p.stopTopic = true
}

//普通运行实例，非自动扩充
func (p *WorkPool) Run(funcJob Job, params ...interface{}) {
	p.taskPool <- taskWork{funcJob, true, params}
}

//用select 去做
func (p *WorkPool) RunAuto(funcJob Job, params ...interface{}) {
	task := taskWork{funcJob, true, params}
	select {
	//正常写入
	case p.taskPool <- task:
		//写入超时 说明队列满了 写入备用队列
	case <-time.After(time.Millisecond * 1000):
		p.taskQueInit()
		p.workerAddConf()
		//task 入备用队列
		p.taskQue <- task
	}
}

//自动初始化备用队列
func (p *WorkPool) taskQueInit() {
	//扩充队列
	if p.taskQue == nil {
		p.taskQue = make(chan taskWork, p.maxNum*2)
	}
}

//自动扩充协程
func (p *WorkPool) workerAddConf() {
	//说明需要扩充进程  协程数量小于 100 协程数量成倍增长
	if p.workNum < 1000 {
		p.workerAdd(p.workNum)
	} else if p.workNum < p.maxNum {
		tmpNum := p.maxNum - p.workNum
		tmpNum = tmpNum / 10
		if tmpNum == 0 {
			tmpNum = 1
		}
		p.workerAdd(1)
	}
}

//自动扩充协程
func (p *WorkPool) workerRemoveConf() {
	//说明需要扩充进程  协程数量小于 100 协程数量成倍增长
	for {
		select {
		case <-time.After(time.Millisecond * 1000 * 600):
			if p.workNum > p.defaultNum && len(p.taskPool) == 0 && len(p.taskQue) == 0 {
				rmNum := (p.workNum - p.defaultNum) / 5
				if rmNum == 0 {
					rmNum = 1
				}
				p.workerRemove(rmNum)
			}
		}
	}

}
func (p *WorkPool) workerAdd(num int) {
	for i := 0; i < num; i++ {
		p.workNum++
		p.workInit(p.workNum)
	}
}
func (p *WorkPool) workerRemove(num int) {
	for i := 0; i < num; i++ {
		task := taskWork{startBool: false}
		p.taskPool <- task
		p.workNum--
	}
}
