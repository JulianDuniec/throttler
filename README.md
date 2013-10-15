throttler
=========

Throttle execution of goroutines


````````````
package main

import (
	"fmt"
	"github.com/julianduniec/throttler"
	"io/ioutil"
	"net/http"
)

func main() {
	//Demonstration: of usage
	//Fetch the contents of a list of urls.
	//Since you are only allowed to make a certain amount of http.Get-requests
	//you can use the throttler to constrain the amount of go-routines 
	//executing http.Get simultaneously

	inputData := []interface{}{
		"http://www.google.com/",
		"http://www.facebook.com/",
		"http://www.twitter.com/",
		"http://www.facebook.com/"}

	//Fetch the body of the url (in)
	worker := func(in interface{}, res chan interface{}) {
		url := in.(string)
		resp, _ := http.Get(url)
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		
		//You need to send the result back to the channel.. or ELSE!!!
		res <- body
	}

  //Handle the result from each worker
	onItemFinished := func(in interface{}) {
		bytes := in.([]byte)
		bodyContent := string(bytes)
		fmt.Println(bodyContent)
	}

	//How many simultaneous workers
	chunkSize := 1

	t := throttler.Throttler{chunkSize, worker, onItemFinished}

	//Will block until finished
	t.Run(inputData)
}

````````````
