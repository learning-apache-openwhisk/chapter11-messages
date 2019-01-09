package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readCred() map[string]interface{} {
	cred, err := ioutil.ReadFile("../cred.json")
	check(err)
	res := make(map[string]interface{})
	err = json.Unmarshal(cred, &res)
	check(err)
	return res
}
func ExampleSend() {
	args := readCred()
	//fmt.Printf("%v", args)
	p := Connect(args)
	fmt.Printf("%v", p)
	Send(p, "queue", []byte("hello"))
	p.Flush(1000)
	// Output:
	// -
}

func ExampleMain() {
	args := readCred()
	args["topic"] = "queue"
	args["message"] = "hello world"
	res := Main(args)
	fmt.Printf("%v\n", res)
	// Output:
	// map[ok:true]
}
