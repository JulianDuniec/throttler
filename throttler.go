package main

import "fmt"

/*
	Throttler takes a list of data and
	executes the worker-function as a goroutine 
	on each item, but at a maximum of simultaneous 
	goroutines defined by chunkSize
*/
type Throttler struct {
	chunkSize    int
	worker       func(interface{}, chan interface{})
	finishedItem func(interface{})
}

/*
	Processed worker on each item in the list,

	Will block until all data has been processed
*/
func (t *Throttler) Run(data []interface{}) {

	//Adapt chunksize 
	maxLength := t.chunkSize
	if len(data) < maxLength {
		maxLength = len(data)
	}

	//Define a channel for the results
	//that we will listen to when synchronizing
	resultChannel := make(chan interface{})

	//Execute the worker-processes
	for i := 0; i < maxLength; i++ {
		go t.worker(data[i], resultChannel)
	}

	//Synchronize the results from each
	//worker process, and execute finishedItem on 
	//each item
	for i := 0; i < maxLength; i++ {
		currentResult := <-resultChannel
		t.finishedItem(currentResult)
	}

	//Check if we have any remaining data that needs to be processed
	//and recurse 
	remainder := data[maxLength:len(data)]
	if len(remainder) > 0 {
		t.Run(remainder)
	}
}
