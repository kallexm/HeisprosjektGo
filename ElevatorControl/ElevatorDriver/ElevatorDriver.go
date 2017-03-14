//package main
package ElevatorDriver

import(
	//./Elev"
	"./simulator/client"
	"fmt"
	"time"
	"../ElevatorStructs"
)

/*type ButtonPlacement struct{
	Floor int
	ButtonType Elev.ButtonType
	Value int 
}*/

var buttonStatusMap = initButtonStatusMap()
var lastMeasuredFloor int

func initButtonStatusMap() map[Elev.ButtonType]map[int]int {
	button_channel_map := map[Elev.ButtonType]map[int]int{
		Elev.Up:map[int]int{1:0,2:0,3:0,4:0},
		Elev.Down:map[int]int{1:0,2:0,3:0,4:0},
		Elev.Comand:map[int]int{1:0,2:0,3:0,4:0}}
	return button_channel_map
}


func pullButtons() ElevatorStructs.ButtonPlacement{
	for f := 1; f <= Elev.N_FLOORS; f ++{
			var b Elev.ButtonType
			for b =  0; b < Elev.N_BUTTONS; b ++{
				value,_ := Elev.ElevGetButtonSignal(Elev.ButtonType(b),f)
				/*if value, err := Elev.ElevGetButtonSignal(Elev.ButtonType(b),f); err != nil{
					fmt.Println("Noeg gikk galt i button pulling err: ", err)*/
				/*} else*/ if value == 1 && buttonStatusMap[b][f] != value{
					buttonPressed := ElevatorStructs.ButtonPlacement{Floor: f, ButtonType: b,Value: 1}
					buttonStatusMap[b][f] = 1
					return buttonPressed
				} else if value == 0{
					buttonStatusMap[b][f] = 0
				}
			} 
	}
	return ElevatorStructs.ButtonPlacement{} 
}


func ElevatorPullingThred(getButtonCh chan<- ElevatorStructs.ButtonPlacement, getFloorCh chan<- int){
	lastMeasuredFloor = 0
	fmt.Println("ElevatorPullingThred started foor loop")
	for{
		if buttonPressed := pullButtons(); buttonPressed != (ElevatorStructs.ButtonPlacement{}){
			fmt.Println("En knapp ble trykket inn")
			getButtonCh <- buttonPressed
		}
		currentFloor := Elev.ElevGetFloorSensorSignal()
		if currentFloor != 0 && currentFloor != lastMeasuredFloor {
			fmt.Println("Vi kom til en etasje")
			Elev.ElevSetFloorIdicator(currentFloor)
			getFloorCh <- currentFloor
		}
		lastMeasuredFloor = currentFloor
		time.Sleep(time.Millisecond*10)
	}
}

func SetLight(button ElevatorStructs.ButtonPlacement){
	if button.ButtonType < 3{
		if err := Elev.ElevSetButtonLamp(button.ButtonType,button.Floor, button.Value); err != nil{
			fmt.Println(err)
		} 
	} else{
		Elev.ElevSetDoorOpenLamp(button.Value)
	}
}

func SetMotor(motor Elev.MotorDir){
	Elev.ElevSetMotorDirection(motor)
}