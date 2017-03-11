package ElevatorControlThread

import
(
	"../../MessageFormat"
	
	"fmt"
)

var msg string

func Thread(from_NodeComm_Ch 			<-chan []byte	,
			to_NodeComm_Ch 				chan<- []byte	,
			ElevCtrl_NodeComm_Mutex_Ch	chan 	bool	,
			ElevCtrl_exit_Ch 			chan<-	bool	) {
	for {
		
		_, err := fmt.Scanln(&msg)
		CheckError(err)
		
		
		sendMsgHeader := MessageFormat.MessageHeader_t{To: MessageFormat.MASTER, From: MessageFormat.ELEVATOR, MsgType: MessageFormat.NEW_ELEVATOR_REQUEST}
		msgToSend, err := MessageFormat.Encode_msg(sendMsgHeader, msg)
		CheckError(err)
		<- ElevCtrl_NodeComm_Mutex_Ch
		to_NodeComm_Ch <- msgToSend
		ElevCtrl_NodeComm_Mutex_Ch <- true
	}
	
	ElevCtrl_exit_Ch <- true
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}