package OrderDistributerThread

/*
1. Holde orden på en kø per heis som er koblet på nettverket.
2. En liste med kø objekter.
3. Dynamisk alokere nye kø er hvis nye noder kobler seg på.
4. Hver kø må ha en bit som sier om køen er aktiv eller ikke.
5. Når en heis disconnect'er vil NodeConnectionManager si ifra om at noden har forsvunnet fra nettverket og da må aktiv kø bit'en deaktiveres.
6. Køen tas vare på, men brukes ikke før den samme noden har koblet seg på igjen og NodeConnectionManager har sagt fra om dette. (Det må i samme tilfelle synkroniseres en ny kø.)
*/

import
(
	"../../MessageFormat"
	
	"fmt"
	"time"
)



func Thread(from_NodeComm_Ch 			<-chan 	[]byte	,
			to_NodeComm_Ch 				chan<- 	[]byte	,
			OrderDist_NodeComm_Mutex_Ch chan 	bool	,
			OrderDist_exit_Ch 			chan<- 	bool	) {
	
	for {
		select {
		case msg := <- from_NodeComm_Ch:
			receivedMsgHeader, data, err := MessageFormat.Decode_msg(msg)
			CheckError(err)
			fmt.Println("Message received:")
			fmt.Println(receivedMsgHeader)
			fmt.Println(data)
			
		default:
			time.Sleep(100*time.Millisecond)
		}
	}
	
	OrderDist_exit_Ch <- true
}


func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}