package main

import (
	"container/list"
	"math"
	"math/rand"
	"os"
	"runtime"
	"time"
	"fmt"
)

func Main(obj map[string]interface{}) map[string]interface{} {
	// do your work
	name, _ := obj["name"].(string)
	msg := make(map[string]interface{})
	msg["message"] = "Hello, " + name + "!"

	seed := 42               // default seed value
	ARRAY_SIZE := 10000      // default array size value
	REQ_NUM := math.MaxInt32 // default request number
	
	fmt.Printf("%v\n", mainLogic(seed, ARRAY_SIZE, REQ_NUM))

	return msg
  }  

func mainLogic(seed int, ARRAY_SIZE int, REQ_NUM int) (map[string]interface{}) {
	start := time.Now().UnixMicro()

	rand.Seed(int64(seed))

	lst := list.New()

	for i := 0; i < ARRAY_SIZE; i++ {
		// Inserting integers directly, assuming payload simulation isn't the focus
		lst.PushFront(rand.Intn(seed)) // Use integers for direct summation
		// Stress GC with nested list
		if i%5 == 0 {
			nestedList := list.New()
			for j := 0; j < rand.Intn(5); j++ {
				nestedList.PushBack(rand.Intn(seed))
			}
			lst.PushBack(nestedList)
		}
		// Immediate removal after insertion to stress GC
		if i%5 == 0 {
			e := lst.PushFront(rand.Intn(seed))
			lst.Remove(e)
		}

	}

	// Sum values and return result
	var sum int64 = 0
	for e := lst.Front(); e != nil; e = e.Next() {
		if val, ok := e.Value.(int); ok {
			sum += int64(val)
		}
	}

	executionTime := time.Now().UnixMicro() - start

	response := map[string]interface{}{
		"sum":           sum,
		"executionTime": executionTime, // Include raw execution time in microseconds
		"requestNumber": REQ_NUM,
		"arraysize":     ARRAY_SIZE,
	}

	gogcValue := os.Getenv("GOGC")
	gomemlimitValue := os.Getenv("GOMEMLIMIT")
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	response["heapAlloc"] = m.HeapAlloc
	// response["heapSys"] = m.HeapSys
	response["heapIdle"] = m.HeapIdle
	// response["heapInuse"] = m.HeapInuse
	response["NextGC"] = m.NextGC
	response["NumGC"] = m.NumGC
	response["GOGC"] = gogcValue
	response["GOMEMLIMIT"] = gomemlimitValue
	
	return response
}