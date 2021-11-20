package mr

import "log"
import "net"
import "os"
import "net/rpc"
import "net/http"
import "sync"
import "time"
import "strconv"

// time interval
const (
	WAIT_TIME = 10
)

// woker type
const (
	Mapper 		uint8 = 0
	Reducer		uint8 = 1
)
// task type
const (
	MAP_STATE		uint8 = 0
	REDUCE_STATE	uint8 = 1
)
// task state in master
const (
	Running		uint8 = 0
	Wait		uint8 = 1
	Finished	uint8 = 2
)
// task Information
type TaskInfo struct {
	Filename	string
	Task_type	uint8

	Task_id		uint8

	M_map		int
	N_reduce 	int
}

type TaskQueue struct {
	tq			[]*TaskInfo
}

func (q *TaskQueue) Push(task *TaskInfo){
	if task == nil {
		return
	}
	q.tq = append(q.tq, task)
}

func (q *TaskQueue) Pop() *TaskInfo{
	if len(q.tq) > 0{
		res := q.tq[len(q.tq)-1]
		q.tq = q.tq[: len(q.tq)-1]
		return res
	}

	return nil
}

func (q *TaskQueue) Empty() bool{
	return len(q.tq) == 0
}

type Coordinator struct {
	// Your definitions here.
	Queue		TaskQueue
	Mu 			sync.Mutex
	Task_states	map[string]uint8

	State		uint8

	M_map		int
	N_reduce	int
}

// Your code here -- RPC handlers for the worker to call.

func ToStr(task *TaskInfo) string {
	if task == nil {
		return ""
	} else {
		kind := strconv.Itoa(int(task.Task_type))
		id := strconv.Itoa(int(task.Task_id))
		return  kind + "*" + id
	}
}

//
// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
//
func (c *Coordinator) Response(req *WorkerRequest, res *WorkerResponse) error {
	switch req.Kind {
	case REUQEST: 
		res.Info, res.State = c.generateTask()
	case REPLY:
		var tp string

		var task *TaskInfo = &req.Info
		if task.Task_type == Mapper{
			tp = "Map "
		} else {
			tp = "Reduce "
		}
		if req.State == FINISHED {
			log.Printf(tp + "Task %v Finished", task.Task_id)
			c.markTask(&(req.Info))
		} else if req.State == ABORTED{
			log.Printf(tp + "Task %v Aborted", task.Task_id)
			c.Mu.Lock()
			c.Task_states[ToStr(task)] = Wait
			c.Queue.Push(task)
			c.Mu.Unlock()
		} else {
			log.Fatalf("Wrong request state")
		}
	default:
		panic("Wrong Response")
	}

	return nil
}

func (c *Coordinator) markTask(task *TaskInfo) {
	defer c.Mu.Unlock()
	c.Mu.Lock()

	ts := ToStr(task)
	if c.Task_states[ts] == Running {
		delete(c.Task_states, ts)
	} else if c.Task_states[ts] == Wait{
		log.Printf("Old Task %v Finished", task.Task_id)
	} else {
		panic("Wrong Mark")
	}
}

func (c *Coordinator) timer(task *TaskInfo) {
	time.Sleep(time.Second * WAIT_TIME)
	defer c.Mu.Unlock()
	c.Mu.Lock()

	var tp string

	if task.Task_type == Mapper{
		tp = "Map "
	} else {
		tp = "Reduce "
	}
	ts := ToStr(task)
	if _, ok := c.Task_states[ts]; ok {
		log.Printf(tp + "Task id %v Time out", (*task).Task_id)
		c.Task_states[ts] = Wait
		c.Queue.Push(task)
	}
}

func (c *Coordinator) generateTask() (TaskInfo, uint8) {
	defer c.Mu.Unlock()
	c.Mu.Lock()
	
	res := new(TaskInfo)
	var state uint8

	if c.Queue.Empty() && len(c.Task_states)==0 && c.State==REDUCE_STATE {
		state = Finished
	} else if c.Queue.Empty() && len(c.Task_states)==0 && c.State==MAP_STATE {
		c.State = REDUCE_STATE
		c.generateReduce()

		res = c.Queue.Pop()
		c.Task_states[ToStr(res)] = Running
		go c.timer(res)
	} else if !c.Queue.Empty() {
		res = c.Queue.Pop()
		c.Task_states[ToStr(res)] = Running
		go c.timer(res)
	} else if c.Queue.Empty() && len(c.Task_states)!=0 {
		state = Wait
	} else {
		panic("Wrong Task State")
	}

	return *res, state
}

func (c *Coordinator) init(files []string, nReduce int) {
	c.N_reduce = nReduce
	c.M_map = len(files)
	c.Task_states = make(map[string]uint8)
	c.generateMap(files)
}

func (c *Coordinator) generateReduce() {
	log.Println("Map task Finished Reduce task Start")
	path := "/home/geekwu/knowledgeSummary/system/dsystem/6.824/src/main/mr-tmp"
	for i:=0; i<c.N_reduce; i++ {
		newTask := new(TaskInfo)
		newTask.Filename = path + "/mr-"
		newTask.Task_type = Reducer
		newTask.Task_id = uint8(i)
		newTask.N_reduce = c.N_reduce
		newTask.M_map = c.M_map

		c.Queue.Push(newTask)
		c.Task_states[ToStr(newTask)] = Wait
	}
}

func (c *Coordinator) generateMap(files []string) {
	for id, file := range files {
		newTask := new(TaskInfo)
		newTask.Filename = file
		newTask.Task_type = Mapper
		newTask.Task_id = uint8(id)
		newTask.N_reduce = c.N_reduce
		newTask.M_map = c.M_map

		c.Queue.Push(newTask)
		c.Task_states[ToStr(newTask)] = Wait
	}
}

//
// start a thread that listens for RPCs from worker.go
//
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
//
func (c *Coordinator) Done() bool {
	defer c.Mu.Unlock()
	c.Mu.Lock()
	// Your code here.
	ret := (len(c.Task_states)==0) && (c.State == REDUCE_STATE)

	return ret
}

//
// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}

	// Your code here.
	c.init(files, nReduce)

	c.server()
	return &c
}
