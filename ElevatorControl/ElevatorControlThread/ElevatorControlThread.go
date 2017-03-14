package ElevatorControlThread 

import
(
	"../ElevatorDriver"
	//"../ElevatorDriver/Elev"
	"../ElevatorDriver/simulator/client"
	"../ElevatorStatus"
	"../ElevatorStatus/timer"
	"../ElevatorStructs"
	"../../MessageFormat"
	"fmt"
	"encoding/json"
)

/*type OrderCompletStruck struct{
	orderComplet bool
}*/

var timerFinishedunconfirmedOrderCh (chan bool)
var orderComplet_Ch (chan bool)
var main_To_Elev_ch (<-chan []byte)
var Elev_To_main (chan <-[]byte)

func Thread(main_To_Elev_ch_ <-chan []byte, Elev_To_main_ chan<- []byte, mutex_Ec_Ch chan bool, ElevCtrl_exit_Ch chan bool) {
	getButtonCh := make(chan ElevatorStructs.ButtonPlacement)
	getFloorCh := make(chan int)
	timerFinishedDoorCh := make(chan bool)
	timerFinishedunconfirmedOrderCh := make(chan bool)
	orderComplet_Ch := make(chan bool,1)
	//Stykkt finnpå bedre variable navn
	main_To_Elev_ch = main_To_Elev_ch_
	Elev_To_main = Elev_To_main_
	Elev.ElevInit()
	go ElevatorDriver.ElevatorPullingThred(getButtonCh,getFloorCh)
	ElevatorStatus.InitElevatorStatus(timerFinishedDoorCh)
	ElevatorDriver.SetMotor(Elev.DirDown)
	<- mutex_Ec_Ch
	getFloor := <- getFloorCh
	newFloorReached(getFloor)
	mutex_Ec_Ch<- true
	fmt.Println("Ec started foor loop")
	for {
		select{
		case <- mutex_Ec_Ch:
			//fmt.Println("Controll has the Ec mutex")
			select{	
			case <- orderComplet_Ch:
				Elev_To_main <- generateMsg(MessageFormat.ORDER_FINISHED_BY_ELEVATOR,ElevatorStructs.OrderCompletStruck{OrderComplet: true})
			case getButton := <- getButtonCh:
				newButtonPresed(getButton)
			case getFloor := <- getFloorCh:
				newFloorReached(getFloor)
			case <- timerFinishedDoorCh:
				cloaseDoor()
			case <- timerFinishedunconfirmedOrderCh:
				Elev_To_main <- generateMsg(MessageFormat.NEW_ELEVATOR_REQUEST,ElevatorStatus.GetUnconfirmedOrder())
			default:
				//is nessasary so to prevent dedlocks. If it is not presant, we can be stuck in the case statement.   
			}
			//fmt.Println("Controll releases the Ec mutex")
			mutex_Ec_Ch <- true
		//See if there is any mesage from the network taht nees to be handeled. 
		case msg := <-main_To_Elev_ch:
				msgHead, data, err := MessageFormat.Decode_msg(msg)
				fmt.Println("msgHead MsgType: ", msgHead.MsgType)
				if err != nil{
					fmt.Println("Error in decoding: ", err)
				}
				if msgHead.MsgType == MessageFormat.NEW_ORDER_TO_ELEVATOR{
					newOrderToElevatorHandeler(data)
				} else if msgHead.MsgType == MessageFormat.SET_LIGHT{
					setLightHandeler(data)
				} 
		default:
			//is nessasary so to prevent dedlocks. If it is not presant, we can be stuck in the case statement.   
		}
	}
	ElevCtrl_exit_Ch <- true
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

func newOrderToElevatorHandeler(data []byte){
	var newOrder ElevatorStructs.Order
	if err := json.Unmarshal(data, &newOrder); err != nil{
		fmt.Println("Error in Unmarshal: ", err)
	}
	motorDir, err, orderComplet := ElevatorStatus.NewCurentOrder(newOrder)
	if err != nil{
		fmt.Println("An error hapene i New_ORDER_TO_ELEVATOR")
	} else if orderComplet == true{
		// Do stof related to order compleat
		//Må legg til tilstandsendring, og starting av lys
		orderComplet_Ch <- true
		ElevatorDriver.SetLight(ElevatorStructs.ButtonPlacement{Floor:0,ButtonType:Elev.Door,Value:1})
	}
	ElevatorDriver.SetMotor(Elev.MotorDir(motorDir))
	fmt.Println("Ec recived new order: ", newOrder)
}

func setLightHandeler(data []byte){
	var button ElevatorStructs.ButtonPlacement
	if err := json.Unmarshal(data, &button); err != nil{
		fmt.Println("Error in MsgType SET_LITGHT", err)
	}
	ElevatorDriver.SetLight(button)
	fmt.Println("Ec handelde new set ligt request ", button)
}

func newButtonPresed(button ElevatorStructs.ButtonPlacement){
	fmt.Println("Buttons presed, floor: ", button.Floor, "button type: ", button.ButtonType)
	ElevatorStatus.NewUnconfirmedOrder(button)
	go timer.TimerThredTwo(timerFinishedunconfirmedOrderCh,3)
	Elev_To_main <- generateMsg(MessageFormat.NEW_ELEVATOR_REQUEST,ElevatorStructs.Order{Floor: button.Floor, OrderDir: ElevatorStructs.Dir(button.ButtonType)})
}

func newFloorReached(floor int){
	fmt.Println("Floor sensor: ", floor)
	motorDir, orderComplet := ElevatorStatus.NewFloor(floor)
	if orderComplet == true{
		// Do stof related to order compleat
		//Må legg til tilstandsendring, og starting av lys.
		orderComplet_Ch <- true 
		ElevatorDriver.SetLight(ElevatorStructs.ButtonPlacement{Floor:0,ButtonType:Elev.Door,Value:1})
	}
	ElevatorDriver.SetMotor(Elev.MotorDir(motorDir))
	//Send melding om ny etasje, og retning. 
	Elev_To_main <- generateMsg(MessageFormat.ELEVATOR_STATUS_DATA,ElevatorStatus.GetPosition())
	fmt.Println("Elevator control has send ELVATOR STATUS DATA") 
}


func cloaseDoor(){
	ElevatorDriver.SetLight(ElevatorStructs.ButtonPlacement{Floor:0,ButtonType:Elev.Door,Value:0})
	motorDir := ElevatorStatus.DoorTimeOut()
	ElevatorDriver.SetMotor(Elev.MotorDir(motorDir))
	Elev_To_main <- generateMsg(MessageFormat.ELEVATOR_STATUS_DATA,ElevatorStatus.GetPosition())
}