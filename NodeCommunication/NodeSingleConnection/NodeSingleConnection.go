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

const readDeadlineTime 	= 1000*time.Millisecond
const writeDeadlineTime = 1000*time.Millisecond


func HandleConnection(	conn 							net.Conn,
						thisConnectsToNodeID 			uint8	,
						from_node_Ch 			<-chan 	[]byte	,
						to_node_Ch 				chan<- 	[]byte	,
						connection_Mutex_Ch		chan 	bool	) {
	var receiveMsg 	[]byte
	var sendMsg		[]byte
	//var receiveErr 	error
	var sendErr		error
	var connBroke 			= false
	var connBrokeMsgSent 	= false

	forLoop:
	for {
		time.Sleep(time.Millisecond*400)
		select {
		case sendMsg = <- from_node_Ch:
			fmt.Println("Melding over nett:", string(sendMsg))
			_ 		= conn.SetWriteDeadline(time.Now().Add(writeDeadlineTime))
			fmt.Println(sendMsg)
			_, sendErr	= conn.Write(sendMsg)
			fmt.Println(sendErr)
			if sendErr != nil {
				fmt.Println(sendErr)
				connBroke = true
			}

			if connBroke == true && connBrokeMsgSent == true {
				head, _, _ := MessageFormat.Decode_msg(sendMsg)
				if head.MsgType == MessageFormat.NODE_DISCONNECTED && head.From == MessageFormat.NODE_COM {
					conn.Close()
					break forLoop
				}
			}

		case <-connection_Mutex_Ch:
			_ 		= conn.SetReadDeadline(time.Now().Add(readDeadlineTime))
			n, receiveErr	:= conn.Read(receiveMsg)
			fmt.Println("receiveMsg:", receiveMsg, "receiveErr:", receiveErr, "n:", n)
			if receiveErr == nil && n > 0 {
				to_node_Ch <- receiveMsg
			}else if e, ok := receiveErr.(net.Error); ok && !e.Timeout() {
				fmt.Println(e)
				connBroke = true
			}

			if connBroke == true {
				fmt.Println("A connection error accured")
				sendHeader := MessageFormat.MessageHeader_t{To: 		MessageFormat.NODE_COM			,
												FromNodeID: thisConnectsToNodeID			,
												MsgType: 	MessageFormat.NODE_DISCONNECTED	}
				connBrokeMsg, _ := MessageFormat.Encode_msg(sendHeader, "")
				to_node_Ch <- connBrokeMsg
				connBrokeMsgSent = true
			}
			connection_Mutex_Ch <- true
		}
	}

}