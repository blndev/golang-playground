package main

import (
	"fmt"
)

//you can play around here: https://play.golang.org/

//define our eventhandler
type callback func(int) bool

//main function
func main() {
	fmt.Println("Hello, playground")
	interruptable(nil)
	interruptable(userFeedback)
  interruptable(noCanceling)
}

//one of our eventhandlers
func userFeedback(data int) bool {
	fmt.Println("callback ", data)
	// you can also ask the user here via input or whatever
  // in our case the user like to cancel the funtion
	return false
}

func noCanceling(data int) bool {
	fmt.Println("ignore callback ", data)
	// must always return true to continue
	return true
}


func interruptable(feedback callback) {
	fmt.Println("Interruptable function")
	if feedback != nil {
		if feedback(123) {
			fmt.Println("callback returns true --> continue call")
		} else {
			fmt.Println("callback returns false --> cancel function")
			return
		}
	} else {
		fmt.Println("no userfeedback defined")
	}

	fmt.Println("function continued")

}

