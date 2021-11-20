package mr

import "fmt"
import "log"
import "net/rpc"
import "hash/fnv"
import "time"
import "os"
import "strconv"
import "io/ioutil"
import "encoding/json"
import "path/filepath"
import "sort"


//
// Map functions return a slice of KeyValue.
//
type KeyValue struct {
	Key   string
	Value string
}

// for sorting by key.
type ByKey []KeyValue

// for sorting by key.
func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

//
// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
//
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}


//
// main/mrworker.go calls this function.
//
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	// Your worker implementation here.

	// uncomment to send the Example RPC to the coordinator.
	// CallExample()
	for {
		res := WorkerResponse{}
		req := WorkerRequest{}
		req.Kind = REUQEST

		// ssk
		SendRPC(&req, &res)

		switch res.State{
		case Running:
			if res.Info.Task_type == Mapper{
				req.State = Map(mapf, res.Info)
			} else {
				req.State = Reduce(reducef, res.Info)
			}
		case Wait:
			time.Sleep(time.Duration(time.Second * 5))
			req.State = FINISHED
		case Finished:
			log.Printf("All Task Finished")
			return
		}

		req.Kind = REPLY
		req.Info = res.Info
		// reply
		SendRPC(&req, &res)
	}
}

func Map(mapf func(string, string) []KeyValue, task TaskInfo) uint8{
	id := task.Task_id
	filename := task.Filename
	n_reduce := task.N_reduce
	path := "/home/geekwu/knowledgeSummary/system/dsystem/6.824/src/main/mr-tmp"

	intermediate := []KeyValue{}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("cannot open %v", filename)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("cannot read %v", filename)
	}
	file.Close()

	kva := mapf(filename, string(content))
	intermediate = append(intermediate, kva...)

	outFiles := make([]*os.File, n_reduce)
	fileEncs := make([]*json.Encoder, n_reduce)
	for i:=0; i<n_reduce; i++ {
		outFiles[i], err = ioutil.TempFile(path, "mr-tmp-*")
		if err != nil {
			log.Fatalf("Error %v", err)
		}
		defer outFiles[i].Close()
		defer os.Remove(outFiles[i].Name())
		fileEncs[i] = json.NewEncoder(outFiles[i])
	}

	for _, kv := range intermediate {
		out := ihash(kv.Key) % n_reduce
		err := fileEncs[out].Encode(&kv)
		if err != nil {
			fmt.Printf("File %v Key %v Value %v Error: %v\n", filename, kv.Key, kv.Value, err)
			log.Fatalf("Json encode failed")
			outFiles[out].Close()
			return ABORTED
		}
	}

	outprefix := path + "/mr-" + strconv.Itoa(int(id)) + "-"
	for outindex, file := range outFiles {
		outname := outprefix + strconv.Itoa(outindex)
		oldpath := filepath.Join(file.Name())

		os.Rename(oldpath, outname)
	}

	return FINISHED
}


func Reduce(reducef func(string, []string) string, task TaskInfo) uint8{
	id := task.Task_id
	fileprefix := task.Filename
	filesuffix := "-" + strconv.Itoa(int(id))
	m_map := task.M_map
	outname := "mr-out-" + strconv.Itoa(int(id))
	path := "/home/geekwu/knowledgeSummary/system/dsystem/6.824/src/main/mr-tmp"
	
	intermediate := []KeyValue{}
	for index := 0; index < m_map; index++ {
		inname := fileprefix + strconv.Itoa(index) + filesuffix
		file, err := os.Open(inname)
		if err != nil {
			log.Printf("Open intermediate file %v failed: %v\n", inname, err)
			return ABORTED
		}
		dec := json.NewDecoder(file)
		for {
			var kv KeyValue
			if err := dec.Decode(&kv); err != nil {
				break
			}
			intermediate = append(intermediate, kv)
		}
   		file.Close()
	}

	sort.Sort(ByKey(intermediate))

	ofile, err := ioutil.TempFile(path, "mr-*")
	if err != nil {
		log.Printf("Create output file %v failed: %v\n", outname, err)
		return ABORTED
	}

	for i:=0; i < len(intermediate); {
		j := i + 1
		for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key {
			j++
		}
		values := []string{}
		for k := i; k < j; k++ {
			values = append(values, intermediate[k].Value)
		}
		output := reducef(intermediate[i].Key, values)

		// this is the correct format for each line of Reduce output.
		fmt.Fprintf(ofile, "%v %v\n", intermediate[i].Key, output)

		i = j
	}

	os.Rename(filepath.Join(ofile.Name()), outname)
	ofile.Close()

	return FINISHED
}

//
// example function to show how to make an RPC call to the coordinator.
//
// the RPC argument and reply types are defined in rpc.go.
//
func SendRPC(req *WorkerRequest, res *WorkerResponse) {
	// send the RPC request, wait for the reply.
	call("Coordinator.Response", req, res)

	// reply.Y should be 100.
	// fmt.Printf("reply.Y %v\n", reply.Y)
}

//
// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := coordinatorSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
