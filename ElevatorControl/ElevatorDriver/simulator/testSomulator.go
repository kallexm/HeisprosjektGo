package main 

import(
	"fmt"
	"./client"
)

func main(){
	Elev.ElevInit()
	Elev.ElevSetButtonLamp(Elev.Up,1,1)
	Elev.ElevSetMotorDirection(Elev.DirUp)
	button,_ := Elev.ElevGetButtonSignal(Elev.Up,1)
	floor := Elev.ElevGetFloorSensorSignal()
	fmt.Println("The button was: ", button)
	fmt.Println("The floor was", floor)
}