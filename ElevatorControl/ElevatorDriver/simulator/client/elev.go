package Elev
/*
#include "elev.h"
#include <assert.h>
#include <stdlib.h>
#include <sys/socket.h>
#include <netdb.h>
#include <stdio.h>
#include <pthread.h>
#include "channels.h"
#include "io.h"
#include "con_load.h"
*/
import "C"
import "../../../ElevatorStructs"
//import "errors"
import "fmt"



type MotorDir int
const(
	DirDown = iota -1 
	DirStop
	DirUp

)


const N_FLOORS = 4
const N_BUTTONS = 3
const MOTOR_SPEED = 2800 



func ElevInit() error{
	C.elev_init()
	return nil

}

func ElevSetMotorDirection(dir MotorDir){
	C.elev_set_motor_direction(C.elev_motor_direction_t(dir))	
}

func ElevSetButtonLamp(button ElevatorStructs.ButtonType, floor int, value int) error{
	var cButton int
	if (button == ElevatorStructs.Up){
		cButton = 0
	} else if button == ElevatorStructs.Down{
		cButton = 1 
	} else{
		cButton = 2
	}
	fmt.Println("cButton in set: ", cButton)
	C.elev_set_button_lamp(C.elev_button_type_t(cButton),C.int(floor-1),C.int(value))
	return nil
}

func ElevSetDoorOpenLamp(value int){
	C.elev_set_door_open_lamp(C.int(value))
	
}

func ElevGetButtonSignal(button ElevatorStructs.ButtonType, floor int)(int, error){
	var cButton int
	if (button == ElevatorStructs.Up){
		cButton = 0
	} else if button == ElevatorStructs.Down{
		cButton = 1 
	} else{
		cButton = 2
	}
	//fmt.Println("cButton in get: ", cButton)
	value := int(C.elev_get_button_signal(C.elev_button_type_t(cButton),C.int(floor-1)))
	return value, nil
}
 

func ElevGetFloorSensorSignal() int{
	floor := int(C.elev_get_floor_sensor_signal())
	return floor+1
}

func ElevSetFloorIdicator(floor int) error{
	C.elev_set_floor_indicator(C.int(floor-1))
	return nil
}

