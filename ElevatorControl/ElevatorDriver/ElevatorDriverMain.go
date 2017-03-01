//package main
package ElevatorDriver

import(
	"./Elev"
	"fmt"
)

type ButtonPlacement struct{
	floor int
	buttonType Elev.ButtonType 
}

/*func main() {
	err := Elev.ElevInit()
	if(err != nil){
		fmt.Println(err)
	}
	fmt.Println("great suckes")
	if err := Elev.ElevSetButtonLamp(Elev.Up, 1,1); err != nil{
		fmt.Println("Noe gikk galt i lape setting", err)
	} else {
		fmt.Println("Klarte å sette knappen")
	}
	if err := Elev.ElevSetFloorIndicator(1); err != nil{
		fmt.Println("Noe gikk glat i floor indicator setting", err)
	} else {
		fmt.Println("Klare å sette flor indicator")
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

func ElevatorDriverThred(setLightCh <-chan ButtonPlacement, setMotorCh <-chan int, getButtonCh chan<- ButtonPlacement, getFloorCh chan<- int) {
	var lightBuff ButtonPlacement
	lightBuff =<- setLightCh 
	fmt.Println("Du Klarte det floor er: ", lightBuff.floor, "uttonType er: ", lightBuff.buttonType)
	
}
