//package main
package ElevatorDriver

import(
	"./Elev"
	"fmt"
	"time"
)

type ButtonPlacement struct{
	Floor int
	ButtonType Elev.ButtonType
	Value int 
}

var buttonStatusMap = initButtonStatusMap()

func initButtonStatusMap() map[Elev.ButtonType]map[int]int {
	button_channel_map := map[Elev.ButtonType]map[int]int{
		Elev.Up:map[int]int{1:0,2:0,3:0,4:0},
		Elev.Down:map[int]int{1:0,2:0,3:0,4:0},
		Elev.Comand:map[int]int{1:0,2:0,3:0,4:0}}
	return button_channel_map
}
/*
func main() {
	err := Elev.ElevInit()
	if(err != nil){
		fmt.Println(err)
	}
	//fmt.Println("great suckes")
	if err := Elev.ElevSetButtonLamp(Elev.Up, 1,1); err != nil{
		fmt.Println("Noe gikk galt i lape setting", err)
	} else {
		fmt.Println("Klarte å sette knappen")
	}
	if err := Elev.ElevSetFloorIndicator(2); err != nil{
		fmt.Println("Noe gikk glat i floor indicator setting", err)
	} else {
		fmt.Println("Klare å sette flor indicator til 2")
	}
	Elev.ElevSetDoorOpenLamp(1)
	if value, err := Elev.ElevGetButtonSignal(Elev.Up, 1); err != nil{
		fmt.Println("NOe gikk galt i button get", err)
	} else {
		fmt.Println("Klarte å button get: ", value)
	}
	value := Elev.ElevGetFloorSensorSignal()
	fmt.Println("Klaret å less floor: ", value)


}*/

func pullButons() ButtonPlacement{
	for f := 1; f <= Elev.N_FLOORS; f ++{
			var b Elev.ButtonType
			for b =  0; b < Elev.N_BUTTONS; b ++{
				if value, err := Elev.ElevGetButtonSignal(b,f); err != nil{
					fmt.Println("Noeg gikk galt i button pulling err: ", err)
				} else if value == 1 && buttonStatusMap[b][f] != value{
					buttonPresed := ButtonPlacement{Floor: f, ButtonType: b,Value: 1}
					buttonStatusMap[b][f] = 1
					return buttonPresed
				} else if value == 0{
					buttonStatusMap[b][f] = 0
				}
			} 
	}
	return ButtonPlacement{} 
}



func ElevatorDriverThred(setLightCh <-chan ButtonPlacement, setMotorCh <-chan Elev.MotorDir, getButtonCh chan<- ButtonPlacement, getFloorCh chan<- int) {
	if err := Elev.ElevInit(); err != nil{
		fmt.Println(err)
	}
	fmt.Println("Vi har initialiert heisen")
	for {
		fmt.Println("Vi starte for løka")
		select{
		case setLight := <- setLightCh:
			if setLight.ButtonType < 3{
				if err := Elev.ElevSetButtonLamp(setLight.ButtonType,setLight.Floor, setLight.Value); err != nil{
					fmt.Println(err)
				} 
			} else{
				Elev.ElevSetDoorOpenLamp(setLight.Value)
			}
		case setMotor := <- setMotorCh:
			Elev.ElevSetMotorDirection(setMotor)
		default:
			if buttonPresed := pullButons(); buttonPresed != (ButtonPlacement{}){
			fmt.Println("En knapp ble trykket inn")
			getButtonCh <- buttonPresed
			}
			if curentFloor := Elev.ElevGetFloorSensorSignal(); curentFloor != 0{
				fmt.Println("Vi kom til en etasje")
				getFloorCh <- curentFloor
			}
			time.Sleep(time.Second * 1)
		}
	}	

	
}