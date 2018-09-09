package main

/*
int mydiv(int a, int b){
    return a/b;
    }
*/
import "C"
import (
	"fmt"
	"runtime"
	"sync"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered call to div", r)
		}
	}()

	div()
}

func div() {
	var wg sync.WaitGroup

	wg.Add(1)
	runtime.LockOSThread()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered div panic", r)
			}
		}()

		cgoPanic(&wg)
	}()

	wg.Wait()
}

func cgoPanic(wg *sync.WaitGroup) {

	//Calling C code which panics
	fmt.Println(C.mydiv(10, 0))
	//panic("Panicking in Go...")
}
