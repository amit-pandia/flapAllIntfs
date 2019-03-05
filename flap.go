package main

import (
	"fmt"
	"os"
	"time"
	"os/exec"
	"strconv"
	"sync"
)

var wg sync.WaitGroup
var flapFrequency int
var subIntfCount int

func main() {
	var err error

	if len(os.Args)<3{
		fmt.Println("./flap <sub-interface-count> <frequency>")
		return
	}

	subIntfCount, err = strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("subIntfCount parse error")
		return
	}

	flapFrequency, err = strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("flapFrequency parse error")
		return
	}

	if subIntfCount != 1 && subIntfCount != 2 && subIntfCount != 4 {
		fmt.Println("subIntfCount invalid")
		return
	}

	fmt.Println("flapFrequency", flapFrequency, "subIntfCount=", subIntfCount)

	startTime :=time.Now().UnixNano()
	flapper()
	endTime:=time.Now().UnixNano()
	executionTime := endTime-startTime

	fmt.Println("Executed Successfully..\nExecution time=",executionTime/1000000, "mili seconds")
}

func flapper() {
	fmt.Println("Start..")
	if subIntfCount == 1 {
		for intf := 1; intf <= 32; intf += 2 {
			wg.Add(1)
			go flap(intf, 0)
		}
	} else if subIntfCount > 1 {
		for intf := 1; intf <= 32; intf += 2 {
			wg.Add(1)
			go func() {
				for subport := 1; subport <= subIntfCount; subport++ {
					wg.Add(1)
					go flap(intf, subport)
				}
				wg.Done()
			}()
		}
	}
	wg.Wait()
}

func flap(intf int, subport int) {
	var cmdargsIfup string
	var cmdargsIfdown string

	if subport == 0 {
		cmdargsIfup = "ifup xeth" + fmt.Sprint(intf)
		cmdargsIfdown = "ifdown xeth" + fmt.Sprint(intf)
	}

	if subport > 0 {
		cmdargsIfup = "ifup xeth" + fmt.Sprint(intf) + "-" + fmt.Sprint(subport)
		cmdargsIfdown = "ifdown xeth" + fmt.Sprint(intf) + "-" + fmt.Sprint(subport)
	}

	for i := 1; i <= flapFrequency; i++ {

		// ifdown
		fmt.Println(time.Now().UTC(),"Flap=", i, "Executing:", cmdargsIfdown)
		cmd := exec.Command("bash", "-c", cmdargsIfdown)
		_, err := cmd.Output()
		if err != nil {
			fmt.Println("Bash:", err)
		}

		// ifup
		fmt.Println(time.Now().UTC(),"Flap=", i, "Executing:", cmdargsIfup)
		cmd = exec.Command("bash", "-c", cmdargsIfup)
		_, err = cmd.Output()
		if err != nil {
			fmt.Println("Bash:", err)
		}

	}

	wg.Done()
	return
}
