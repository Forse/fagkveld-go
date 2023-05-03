package main

import (
	"errors"
	"fagkveld/internal/slides/12_ddd/cargo"
	"fmt"
)

type shipID struct {
	value string
}

func main() {
	trackingID, err := cargo.NewTrackingID("123")
	if err != nil {
		if errors.Is(err, cargo.ErrEmptyTrackingID) {
			fmt.Println("empty tracking id")
			return
		}
		panic(err)
	}
	trackingID2, err := cargo.NewTrackingID("123")
	if err != nil {
		if errors.Is(err, cargo.ErrEmptyTrackingID) {
			fmt.Println("empty tracking id 2")
			return
		}
		panic(err)
	}

	// var id shipID = trackingID // Fortunately, this doesn't compile

	fmt.Printf("tracking id: %s\n", trackingID)              // 123
	fmt.Printf("tracking id 2: %s\n", trackingID2)           // 123
	fmt.Printf("are equal: %t\n", trackingID == trackingID2) // true
	// fmt.Printf("are equal: %t\n", trackingID == "123")    // doesn't compile
}
