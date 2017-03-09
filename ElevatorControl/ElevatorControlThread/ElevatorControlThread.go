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

func ElevatorControlThred(main_To_Elev_ch <-chan []byte, Elev_To_main chan<- []byte, mutex_Ec_Ch chan bool) {
	setLightCh := make(chan ElevatorDriver.ButtonPlacement)
	setMotorCh := make(chan Elev.MotorDir)
	getButtonCh := make(chan ElevatorDriver.ButtonPlacement)
	getFloorCh := make(chan int)
	mutex_Ed_Ch := make(chan bool,1)
	timerFinishedDoorCh := make(chan bool)
	timerFinishedunconfirmedOrderCh := make(chan bool)

	mutex_Ed_Ch <- true

	go ElevatorDriver.ElevatorDriverThred(setLightCh,setMotorCh,getButtonCh,getFloorCh,mutex_Ed_Ch)
	ElevatorStatus.InitElevatorStatus(timerFinishedDoorCh)
	//var button = ElevatorDriver.ButtonPlacement{Floor: 1,ButtonType: Elev.Comand,Value: 1}
	//setLightCh <- button
	for {
		select{
		case <- mutex_Ec_Ch:
			select{
			case getButton := <- getButtonCh:
				fmt.Println("Buttons presed, floor: ", getButton.Floor, "button type: ", getButton.ButtonType)
				ElevatorStatus.NewUnconfirmedOrder(getButton)
				go timer.TimerThredTwo(timerFinishedunconfirmedOrderCh,3)
				// Place Holder, send mesage over net about new unconfirmed Order
				Elev_To_main <- generateMsg(MessageFormat.NEW_ELEVATOR_REQUEST,ElevatorStatus.Order{Floor: getButton.Floor, OrderDir: ElevatorStatus.Dir(getButton.ButtonType)})
			case getFloor := <- getFloorCh:
				fmt.Println("Floor sensor: ", getFloor)
				motorDir, orderComplet := ElevatorStatus.NewFloor(getFloor)
				if orderComplet == true{
					// Do stof related to order compleat
				}
				<-mutex_Ed_Ch
				setMotorCh <- Elev.MotorDir(motorDir)
				mutex_Ed_Ch<- true
				//Send melding om ny etasje, og retning. 
				Elev_To_main <- generateMsg(MessageFormat.ELEVATOR_STATUS_DATA,ElevatorStatus.GetPosition())
			case <- mutex_Ed_Ch:
				select{	 
				case /*doorTimer :=*/ <- timerFinishedDoorCh:
					motorDir := ElevatorStatus.DoorTimeOut()
					setMotorCh <- Elev.MotorDir(motorDir)
					Elev_To_main <- generateMsg(MessageFormat.ELEVATOR_STATUS_DATA,ElevatorStatus.GetPosition())
				}
				mutex_Ed_Ch <- true
			case /*UnconfirmedOrderTimer :=*/ <- timerFinishedunconfirmedOrderCh:
				Elev_To_main <- generateMsg(MessageFormat.NEW_ELEVATOR_REQUEST,ElevatorStatus.GetUnconfirmedOrder())
				//Logikk for å sende melding på nytt on ny ordere. 
			}
			mutex_Ec_Ch <- true
		//Tar mutexen til ElevatorDriver, da er det kun du som få lov til å skriv til den. 
		case <- mutex_Ed_Ch:
			select{
			// Skjekker om Main prøver å snakke med deg. 	
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
					motorDir, err, orderComplet := ElevatorStatus.NewCurentOrder(newOrder)
					if err != nil || orderComplet == false{
						fmt.Println("An error hapene i New_ORDER_TO_ELEVATOR")
					}
					setMotorCh <- Elev.MotorDir(motorDir)
					fmt.Println("New Order: ", newOrder)


				} else if msgHead.MsgType == MessageFormat.SET_LIGHT{
					var button ElevatorDriver.ButtonPlacement
					if err := json.Unmarshal(data, &button); err != nil{
						fmt.Println("Error in MsgType SET_LITGHT", err)
					}
					setLightCh <- button
					fmt.Println("New button: ", button)
				}
			}
			mutex_Ed_Ch <- true

		default:
			//Kan være nødvendig, vet ikke helt hvorfor.  
		}
	}
}


func generateMsg(msgType MessageFormat.MsgType_t, inputStruct interface{}) []byte{
	msgHeader := MessageFormat.MessageHeader_t{To: MessageFormat.MASTER,From: MessageFormat.ELEVATOR,MsgType:msgType}
	data, err := json.Marshal(inputStruct)
	if err != nil{
		fmt.Println("Error in GereateMsg")
	}
	msg, err := MessageFormat.Encode_msg(msgHeader, data)
	if err != nil{
		fmt.Println("Error in GenerateMsg Msg Encoding")
	}
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


