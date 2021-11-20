package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import "os"
import "strconv"

//
// example to show how to declare the arguments
// and reply for an RPC.
//
// message type
const (
	REUQEST 		uint8 = 0
	REPLY			uint8 = 1
)
// task state in worker
const (
	FINISHED		uint8 = 0
	ABORTED			uint8 = 1
)

type WorkerRequest struct {
	Kind			uint8
	Info			TaskInfo

	State			uint8
}

type WorkerResponse struct {
	Info 			TaskInfo

	State			uint8
}

// Add your RPC definitions here.


// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
