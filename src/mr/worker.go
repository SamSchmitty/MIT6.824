package mr

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
)

// Map functions return a slice of KeyValue.
type KeyValue struct {
	Key   string
	Value string
}

// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

// main/mrworker.go calls this function.
// worker should be reworked to ask for next task, this is another sequential worker.
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	fileName := CallNextFile()

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("cannot open %v", fileName)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("cannot read %v", fileName)
	}
	file.Close()

	kv := mapf(fileName, string(content))

	intermediate, _ := os.Create("mr-0-0")
	enc := json.NewEncoder(intermediate)
	for _, kvp := range kv {
		//need to break out from the loop and log error if one occurs
		err := enc.Encode(kvp)
		if err != nil {
			log.Fatalf("Failed to encode intermediate file %v. err: %v", fileName, err.Error())
		}
	}
}

func CallNextFile() string {
	args := WorkerFileNameRequest{}

	args.HasMapped = false

	reply := WorkerFileNameResponse{}

	call("Coordinator.GetNextFile", &args, &reply)

	return reply.FileName
}

// example function to show how to make an RPC call to the coordinator.
//
// the RPC argument and reply types are defined in rpc.go.
func CallExample() {

	// declare an argument structure.
	args := ExampleArgs{}

	// fill in the argument(s).
	args.X = 99

	// declare a reply structure.
	reply := ExampleReply{}

	// send the RPC request, wait for the reply.
	// the "Coordinator.Example" tells the
	// receiving server that we'd like to call
	// the Example() method of struct Coordinator.
	ok := call("Coordinator.Example", &args, &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Printf("reply.Y %v\n", reply.Y)
	} else {
		fmt.Printf("call failed!\n")
	}
}

// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
// rpcname: class name and method to call. ex Coordinator.GetNextFile()
// args: arguments to send to rpc
// reply: the response from the rpc. In this case the next file that has not been processed.
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
