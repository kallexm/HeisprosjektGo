package ElevatorControlThread

import
(
	"../../MessageFormat"
	
	"fmt"
)

var msg string

func ElevatorControl_thread(from_NodeComm_Ch <-chan []byte, to_NodeComm_Ch chan<- []byte, ElevCtrl_exit_Ch chan<- bool) {
	for {
		
		_, err := fmt.Scanf("%q", &msg)
		CheckError(err)
		
		
		sendMsgHeader := MessageFormat.MessageHeader_t{To: MessageFormat.ORDER_DIST, From: MessageFormat.ELEVATOR, MsgType: MessageFormat.NEW_ELEVATOR_REQUEST}
		msgToSend, err := MessageFormat.Encode_msg(sendMsgHeader, msg)
		CheckError(err)
		to_NodeComm_Ch <- msgToSend
	}
	
	ElevCtrl_exit_Ch <- true
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}