package ElevatorControlThread 

import
(
	"../ElevatorDriver"
	"../ElevatorDriver/Elev"
	"../ElevatorStatus"
	"../ElevatorStatus/timer"
	"../../MessageFormat"
	"fmt"
	"encoding/json"
)

func ElevatorControlThred(main_To_Elev_ch chan<- []bool, Elev_To_main <-chan []bool) {
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
			go timer.TimerThredTow(timerFinishedunconfirmedOrderCh)
			// Place Holder, send mesage over net about new unconfirmed Order
			Elev_To_main <- generateMsg(MessageFormat.NEW_ELEVATOR_REQUEST,ElevatorDriver.Order{Floor: getButton.Floor, OrderDir: dir(getButton.ButtonType)})
		case getFloor := <- getFloorCh:
			fmt.Println("Floor sensor: ", getFloor)
			motorDir, orderComplet := Elevatorstatus.NewFloor(getFloor)
			setMotorCh <- Elev.MotorDir(motorDir)
			//Send melding om ny etasje.
			//Send meldign om status. 
			 
		case doorTimer := <- timerFinishedDoorCh:
			motorDir = Elevatorstatus.DoorTimeOut()
			setMotorCh <- Elev.MotorDir(motorDir)
		case UnconfirmedOrderTimer := <- timerFinishedunconfirmedOrderCh:
			//Logikk for å sende melding på nytt on ny ordere.  

		case msg := <-main_To_Elev_ch:
			//Loggik for dekoding av melding og ta riktig beslutning.
			msgHead, data, err := MessageFormat.Decode_msg(msg)
			fmt.Println("msgHead MsgType: ", msgHead.MsgType)
			if err != nil{
				fmt.Println("Error in decoding: ", err)
			}
			if msgHead.MsgType == MessageFormat.NEW_ORDER_TO_ELEVATOR{
				var newOrder ElevatorStatus.Order
				if err := json.Unmarshal(data, &newOrder); err != nil{
					fmt.Println("Error in Unmarshal: ", err)
				}
				fmt.Println("New Order: ", newOrder)


			} else if msgHead.MsgType == MessageFormat.SET_LIGHT{
				var button ElevatorDriver.ButtonPlacement
				if err := json.Unmarshal(data, &button); err != nil{
					fmt.Println("Error in MsgType SET_LITGHT", err)
				}
				//setLightCh <- button
				fmt.Println("New button: ", button)
			}
		default:
			//Kan være nødvendig, vet ikke helt hvorfor.  
		}



	}
}


func generateMsg(msgType MessageFormat.MsgType_t, inputStruct interface{}){
	msgHeader := MessageFormat.MessageHeader_t{To: MessageFormat.MASTER,ToNodeID:nil,From: MessageFormat.ELEVATOR, FromNodeID:nil,MsgType:MsgType}
	data, err := json.Marshal(inputStruct)
	if err != nil{
		fmt.Println("Error in GereateMsg")
	}
	msg, err := MessageFormat.Encode_msg(msgHeader, data)
	if err != nill{
		fmt.Println("Error in GenerateMsg Msg Encoding")
	} else{
	return msg
	
}


/*func ElevatorControlThread(main_To_Elev_ch <-chan []byte, Elev_To_main chan<- []byte) {
	for {
		select{
		case msg := <- main_To_Elev_ch:
			msgHead, data, err := MessageFormat.Decode_msg(msg)
			fmt.Println("msgHead MsgType: ", msgHead.MsgType)
			if err != nil{
				fmt.Println("Error in decoding: ", err)
			}
			if msgHead.MsgType == MessageFormat.NEW_ORDER_TO_ELEVATOR{
				var newOrder ElevatorStatus.Order
				if err := json.Unmarshal(data, &newOrder); err != nil{
					fmt.Println("Error in Unmarshal: ", err)
				}
				fmt.Println("New Order: ", newOrder)


			} else if msgHead.MsgType == MessageFormat.SET_LIGHT{
				var button ElevatorDriver.ButtonPlacement
				if err := json.Unmarshal(data, &button); err != nil{
					fmt.Println("Error in MsgType SET_LITGHT", err)
				}
				//setLightCh <- button
				fmt.Println("New button: ", button)
			}
		//default:
			//fmt.Println("Vi er i for løkka")
		}
	}
	fmt.Println("Hade")	
}*/


