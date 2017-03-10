package NodeSingleConnection
/*
||	File: NodeSingleConnection
||
||	Author:  Andreas Hanssen Moltumyr	
||	Partner: Martin Mostad
||	Date: 	 Spring 2017
||	Course:  TTK4145 - Real-time Programming, NTNU
||	
||	Summary of File:
||		The function HandleConnection(...) should be used as a goroutine.
||		The function takes a TCP connection, an input and an output channel and passes
||		messages between these and sends notification to NodeConnectionManager if the connection breaks.
||
*/

/*[FFF]
1. Inneholder en tråd som skal bli kalt hver gang en ny Node vil koble seg på.
2. Denne vil motta data den skal sende over en channel hvor den andre enden ligger lagret i NodeRoutingTable listen sammen med all den andre data om hvilken node som denne tråden kobler til.
3. Har i oppdrag å pinge og sørge for at forbindelsen opprettholdes.
4. Har i oppdrag å motta meldinger fra annen node og sende til NodeMessageRelay
5. Har i oppdrag å sende meldinger til annen node som er mottatt fra NodeMessageRelay 
6. Må varsle NodeConnectionManager hvis den mister forbindelsen/timer ut.
(7.) Føre statestikk over hvor mye som sendes.
*/

import
(
	"../../MessageFormat"
	
	"fmt"
	"net"
	"time"
	
)

const readDeadlineTime 	= 50*time.Millisecond
const writeDeadlineTime = 50*time.Millisecond


func HandleConnection(	conn 							net.Conn,
						thisConnectsToNodeID 			uint8	,
						from_node_Ch 			<-chan 	[]byte	,
						to_node_Ch 				chan<- 	[]byte	,
						connection_Mutex_Ch		chan 	bool	) {
	var msg []byte
	var err error
	
	for {
		select {
		case msg = <- from_node_Ch:
			_ 		= conn.SetWriteDeadline(time.Now().Add(writeDeadlineTime))
			_, err	= conn.Write(msg)
			if err != nil {
				fmt.Println("A connection error accured")
				fmt.Println(err)
				sendConnErrorToNodeConnManager(thisConnectsToNodeID, from_node_Ch, to_node_Ch)
				conn.Close()
			}
		default:
			_ 		= conn.SetReadDeadline(time.Now().Add(readDeadlineTime))
			_, err	= conn.Read(msg)
			if err == nil {
				to_node_Ch <- msg
			}else if e, ok := err.(net.Error); ok && !e.Timeout() {
				fmt.Println("A connection error accured")
				fmt.Println(e)
				sendConnErrorToNodeConnManager(thisConnectsToNodeID, from_node_Ch, to_node_Ch)
				conn.Close()
			}
		}
		
		
		
	}
	
	/* Implement the notes above */
}


func sendConnErrorToNodeConnManager(thisConnectsToNodeID 			uint8,
									from_node_Ch 			<-chan []byte,
									to_node_Ch 				chan<- []byte) {
	// Her bør det legges inn at en venter og prøver å sende på nytt igjen hvis meldingen bruker for lang tid før en får svar.
	sendHeader := MessageFormat.MessageHeader_t{To: 		MessageFormat.NODE_COM			,
												FromNodeID: thisConnectsToNodeID			,
												MsgType: 	MessageFormat.NODE_DISCONNECTED	}
	sendMsg, _ := MessageFormat.Encode_msg(sendHeader, "")
	
	select {
	case to_node_Ch <- sendMsg:
		
	}
	
	
	receiveMsg := <- from_node_Ch
	for {
		receiveHeader, _, _ := MessageFormat.Decode_msg(receiveMsg)
		if receiveHeader.MsgType == MessageFormat.NODE_DISCONNECTED && receiveHeader.From == MessageFormat.NODE_COM {
			break
		}
		// Messages can be dropped here
	}
}






