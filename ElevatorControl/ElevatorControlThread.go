package main 

import
(
	"./ElevatorDriver"
	"./ElevatorDriver/Elev"
	"./Elevatorstatus"
	"./Elevatorstatus/timer"
	"fmt"
)

func main() {
	setLightCh := make(chan ElevatorDriver.ButtonPlacement)
	setMotorCh := make(chan Elev.MotorDir)
	getButtonCh := make(chan ElevatorDriver.ButtonPlacement)
	getFloorCh := make(chan int)
	timerFinishedDoorCh := make(chan bool)
	timerFinishedunconfirmedOrderCh := make(chan bool)

	go ElevatorDriver.ElevatorDriverThred(setLightCh,setMotorCh,getButtonCh,getFloorCh)
	Elevatorstatus.InitElevatorStatus(timerFinishedDoorCh)
	//var button = ElevatorDriver.ButtonPlacement{Floor: 1,ButtonType: Elev.Comand,Value: 1}
	//setLightCh <- button
	for {
		select{
		case getButton := <- getButtonCh:
			fmt.Println("Buttons presed, floor: ", getButton.Floor, "button type: ", getButton.ButtonType)
			Elevatorstatus.NewUnconfirmedOrder(getButton)
			go 
			// Place Holder, send mesage over net about new unconfirmed Order
		case getFloor := <- getFloorCh:
			fmt.Println("Floor sensor: ", getFloor)
			if motorDir, err := Elevatorstatus.NewFloor(getFloor); err != nil{
				fmt.Println("Noe gikk galt i ElevatorControlThred, err: ",err)
			}
			setMotorCh <- Elev.MotorDir(motorDir)
			//Send melding om ny etasje.
			 
		case doorTimer := <- timerFinishedDoorCh:
			motorDir = Elevatorstatus.DoorTimeOut()
			setMotorCh <- Elev.MotorDir(motorDir)
		case UnconfirmedOrderTimer := <- timerFinishedunconfirmedOrderCh:
			//Logikk for å sende melding på nytt on ny ordere.  


		}



	}

}


