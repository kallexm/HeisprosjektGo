package ElevatorDriver

import(
	//./Elev"
	"./simulator/client"
	
	"../ElevatorStructs"

	"fmt"
	"time"
)


var buttonStatusMap = initButtonStatusMap()
var lastMeasuredFloor int

func initButtonStatusMap() map[ElevatorStructs.ButtonType]map[int]int {
	button_channel_map := map[ElevatorStructs.ButtonType]map[int]int{
		ElevatorStructs.Up:		map[int]int{1:0, 2:0, 3:0, 4:0},
		ElevatorStructs.Down:		map[int]int{1:0, 2:0, 3:0, 4:0},
		ElevatorStructs.Comand:	map[int]int{1:0, 2:0, 3:0, 4:0}}
	return button_channel_map
}


func pullButtons() ElevatorStructs.ButtonPlacement{
	for f := 1; f <= Elev.N_FLOORS; f ++{
		var b ElevatorStructs.ButtonType
		for b =  0; b < Elev.N_BUTTONS; b ++{
			value,_ := Elev.ElevGetButtonSignal(ElevatorStructs.ButtonType(b),f)
			if value == 1 && buttonStatusMap[b][f] != value{
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


func ElevatorPullingThread(getButton_Ch chan<- ElevatorStructs.ButtonPlacement, getFloor_Ch chan<- int){
	lastMeasuredFloor = 0
	fmt.Println("ElevatorPullingThread started foor loop")
	for{
		if buttonPressed := pullButtons(); buttonPressed != (ElevatorStructs.ButtonPlacement{}){
			fmt.Println("En knapp ble trykket inn")
			getButton_Ch <- buttonPressed
		}
		currentFloor := Elev.ElevGetFloorSensorSignal()
		if currentFloor != 0 && currentFloor != lastMeasuredFloor {
			fmt.Println("Vi kom til en etasje")
			Elev.ElevSetFloorIdicator(currentFloor)
			getFloor_Ch <- currentFloor
		}
		lastMeasuredFloor = currentFloor
		time.Sleep(time.Millisecond*10)
	}
}

func SetLight(button ElevatorStructs.ButtonPlacement){
	if button.ButtonType < 3{
		if err := Elev.ElevSetButtonLamp(button.ButtonType, button.Floor, button.Value); err != nil{
			fmt.Println(err)
		} 
	} else{
		Elev.ElevSetDoorOpenLamp(button.Value)
	}
}

func SetMotor(motor Elev.MotorDir){
	Elev.ElevSetMotorDirection(motor)
}