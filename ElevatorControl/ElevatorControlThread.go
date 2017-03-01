package main 

import
(
	"./ElevatorDriver"
	"./ElevatorDriver/Elev"
	"fmt"
)

func main() {
	setLightCh := make(chan ElevatorDriver.ButtonPlacement)
	setMotorCh := make(chan Elev.MotorDir)
	getButtonCh := make(chan ElevatorDriver.ButtonPlacement)
	getFloorCh := make(chan int)

	go ElevatorDriver.ElevatorDriverThred(setLightCh,setMotorCh,getButtonCh,getFloorCh)
	var button = ElevatorDriver.ButtonPlacement{Floor: 1,ButtonType: Elev.Comand,Value: 1}
	setLightCh <- button
	for {
		select{
		case getButton := <- getButtonCh:
			fmt.Println("Buttons presed, floor: ", getButton.Floor, "button type: ", getButton.ButtonType)
		case getFloor := <- getFloorCh:
			fmt.Println("Floor sensor: ", getFloor)
		}


	}

}


