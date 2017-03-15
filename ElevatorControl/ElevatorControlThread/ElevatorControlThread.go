package ElevatorControlThread 

import
(
	//"../ElevatorDriver/Elev"
	"../ElevatorDriver/simulator/client"
	
	"../ElevatorDriver"
	"../ElevatorStatus"
	"../ElevatorStatus/timer"
	"../ElevatorStructs"
	"../../MessageFormat"

	"fmt"
	"encoding/json"
)

var reSendOrderTimer_Ch (chan bool)
var orderComplete_Ch 	(chan bool)
var from_MsgRelay_Ch 	(<-chan []byte)
var to_MsgRelay_Ch 		(chan<- []byte)
var masterOnNet_Ch      (chan bool)
var masterOnNet         bool

func Thread(from_MsgRelay_Ch_ <-chan []byte, to_MsgRelay_Ch_ chan<- []byte, mutex_Ec_Ch chan bool, ElevCtrl_exit_Ch chan bool) {
	getButton_Ch			:= make(chan ElevatorStructs.ButtonPlacement)
	getFloor_Ch 			:= make(chan int)
	timerFinishedDoor_Ch 	:= make(chan bool)
	reSendOrderTimer_Ch 	:= make(chan bool)
	orderComplete_Ch 		= make(chan bool, 1)
	masterOnNet_Ch          = make(chan bool, 1)
	from_MsgRelay_Ch 	= from_MsgRelay_Ch_
	to_MsgRelay_Ch 		= to_MsgRelay_Ch_

	masterOnNet = false 

	Elev.ElevInit()
	go ElevatorDriver.ElevatorPullingThread(getButton_Ch,getFloor_Ch)

	ElevatorStatus.InitElevatorStatus(timerFinishedDoor_Ch)
	ElevatorDriver.SetMotor(Elev.DirDown)

	<- mutex_Ec_Ch
	getFloor := <- getFloor_Ch
	newFloorReached(getFloor)
	mutex_Ec_Ch<- true
	fmt.Println("Ec started for loop")
	for {
		select{
		case <- mutex_Ec_Ch:
			//fmt.Println("Fikk mutex_Ec_Ch")
			select{	
			case <- orderComplete_Ch:
				to_MsgRelay_Ch <- generateMsg(MessageFormat.ORDER_FINISHED_BY_ELEVATOR,ElevatorStructs.OrderCompletStruck{OrderComplet: true})
			case masterOnNet = <- masterOnNet_Ch:
				to_MsgRelay_Ch <- generateMsg(MessageFormat.ELEVATOR_STATUS_DATA,ElevatorStatus.GetPosition())
			case getButton 	:= <- getButton_Ch:
				newButtonPressed(getButton)
			case getFloor 	:= <- getFloor_Ch:
				newFloorReached(getFloor)
			case <- timerFinishedDoor_Ch:
				closeDoor()
			case <- reSendOrderTimer_Ch:
				unconfirmedOrder, left := ElevatorStatus.GetUnconfirmedOrder()
				if(left){
					to_MsgRelay_Ch <- generateMsg(MessageFormat.NEW_ELEVATOR_REQUEST,unconfirmedOrder)
					go timer.TimerThread(reSendOrderTimer_Ch, 3)
				}
			default:
			}
			mutex_Ec_Ch <- true
		//See if there is any message from the network that needs to be handled. 
		case msg := <- from_MsgRelay_Ch:
				msgHead, data, err := MessageFormat.Decode_msg(msg)
				fmt.Println("msgHead MsgType: ", msgHead.MsgType)
				if err != nil{
					fmt.Println("Error in decoding: ", err)
				}
				if msgHead.MsgType == MessageFormat.NEW_ORDER_TO_ELEVATOR {
					newOrderToElevatorHandler(data)
				} else if msgHead.MsgType == MessageFormat.SET_LIGHT {
					setLightHandler(data)
				} else if msgHead.MsgType == MessageFormat.MASTER_ON_NET{
					masterOnNet_Ch <- true
				} else if msgHead.MsgType == MessageFormat.MASTER_NOT_ON_NET{
					masterOnNet = false
				} else if msgHead.MsgType == MessageFormat.NEW_ELEVATOR_REQUEST_ACCEPTED{
					ElevatorStatus.RemoveUnconfirmedOrder()
				}
		}
	}
	ElevCtrl_exit_Ch <- true
}




func generateMsg(msgType MessageFormat.MsgType_t, inputStruct interface{}) []byte{
	msgHeader := MessageFormat.MessageHeader_t{	To:			MessageFormat.MASTER	,
												From: 		MessageFormat.ELEVATOR 	,
												MsgType:	msgType 				}
	
	data, err := json.Marshal(inputStruct)
	if err != nil{
		fmt.Println("Error in GenerateMsg")
	}
	msg, err  := MessageFormat.Encode_msg(msgHeader, data)
	if err != nil{
		fmt.Println("Error in GenerateMsg Msg Encoding")
	}
	return msg
	
}




func newOrderToElevatorHandler(data []byte){
	var newOrder ElevatorStructs.Order
	if err := json.Unmarshal(data, &newOrder); err != nil{
		fmt.Println("Error in Unmarshal: ", err)
	}
	motorDir, err, orderComplete := ElevatorStatus.NewCurrentOrder(newOrder)
	if err != nil{
		fmt.Println("An error hapene i New_ORDER_TO_ELEVATOR")
	} else if orderComplete == true{
		// Do stuff related to order complete
		//Må legg til tilstandsendring, og starting av lys
		orderComplete_Ch <- true
		ElevatorDriver.SetLight(ElevatorStructs.ButtonPlacement{Floor: 0, ButtonType: ElevatorStructs.Door, Value: 1})
	}
	ElevatorDriver.SetMotor(Elev.MotorDir(motorDir))
	fmt.Println("Ec recived new order: ", newOrder)
}




func setLightHandler(data []byte){
	var button ElevatorStructs.ButtonPlacement
	if err := json.Unmarshal(data, &button); err != nil{
		fmt.Println("Error in MsgType SET_LITGHT", err)
	}
	ElevatorDriver.SetLight(button)
	fmt.Println("Ec handelde new set ligt request ", button)
}




func newButtonPressed(button ElevatorStructs.ButtonPlacement){
	if masterOnNet == false{
		return
	}
	fmt.Println("Buttons presed, floor: ", button.Floor, "button type: ", button.ButtonType)
	ElevatorStatus.NewUnconfirmedOrder(button)
	go timer.TimerThread(reSendOrderTimer_Ch, 3)
	to_MsgRelay_Ch <- generateMsg(MessageFormat.NEW_ELEVATOR_REQUEST,ElevatorStructs.Order{Floor: button.Floor, OrderDir: ElevatorStructs.Dir(button.ButtonType)})
}




func newFloorReached(floor int){
	fmt.Println("Floor sensor: ", floor)
	motorDir, orderComplete := ElevatorStatus.NewFloor(floor)
	if orderComplete == true{
		// Do stuff related to order complete
		//Må legg til tilstandsendring, og starting av lys.
		orderComplete_Ch <- true
		ElevatorDriver.SetLight(ElevatorStructs.ButtonPlacement{Floor:0,ButtonType:ElevatorStructs.Door,Value:1})
	}
	fmt.Println("Motor dir is: ", motorDir)
	ElevatorDriver.SetMotor(Elev.MotorDir(motorDir))
	//Send melding om ny etasje, og retning.
	if masterOnNet{
		to_MsgRelay_Ch <- generateMsg(MessageFormat.ELEVATOR_STATUS_DATA,ElevatorStatus.GetPosition())
		fmt.Println("Elevator control has send ELVATOR STATUS DATA") 		
	}
} 
	




func closeDoor(){
	ElevatorDriver.SetLight(ElevatorStructs.ButtonPlacement{Floor: 0,ButtonType:ElevatorStructs.Door,Value: 0})
	motorDir := ElevatorStatus.DoorTimeOut()
	ElevatorDriver.SetMotor(Elev.MotorDir(motorDir))
	if masterOnNet{
		to_MsgRelay_Ch <- generateMsg(MessageFormat.ELEVATOR_STATUS_DATA,ElevatorStatus.GetPosition())
	}	
}