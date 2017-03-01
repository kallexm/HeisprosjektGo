package main 

import
(
	"./ElevatorDriver"
	"./ElevatorDriver/Elev"
	"fmt"
)

func main() {
	setLightCh := make(chan ElevatorDriver.ButtonPlacement)
	setMotorCh := make(chan int)
	getButtonCh := make(chan ElevatorDriver.ButtonPlacement)
	getFloorCh := make(chan int)

	go ElevatorDriver.ElevatorDriverThred(setLightCh,setMotorCh,getButtonCh,getFloorCh)

	testButton := ElevatorDriver.ButtonPlacement{1,Elev.Up}

	setLightCh <- testButton
	fmt.Println("Dett gikk bra")

}


