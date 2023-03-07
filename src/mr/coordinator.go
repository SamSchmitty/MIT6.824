package mr

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

type Coordinator struct {
	// Your definitions here.
	filesToProcess    []string
	nReduce           int
	numProcessedFiles int //no idea if I will use this
}

// Your code here -- RPC handlers for the worker to call.
// Returns a files string based on the number of Processed files. This probably needs a mutex lock around it when called.
// Will also need some way to track coordinator and tasks to stop execution if the task is taking too long.
func (c *Coordinator) GetNextFile(args *WorkerFileNameRequest, reply *WorkerFileNameResponse) error {
	reply.FileName = c.filesToProcess[c.numProcessedFiles]
	c.numProcessedFiles += 1
	return nil
}

// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

// start a thread that listens for RPCs from worker.go
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

// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
func (c *Coordinator) Done() bool {
	ret := false

	// Your code here.

	return ret
}

// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}

	c.filesToProcess = files
	c.nReduce = nReduce
	c.numProcessedFiles = 0

	c.server()
	return &c
}
