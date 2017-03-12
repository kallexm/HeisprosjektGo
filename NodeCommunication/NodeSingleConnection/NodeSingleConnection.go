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
	"bytes"
	//"errors"
	
)

const readDeadlineTime 			= 50*time.Millisecond
const writeDeadlineTime 		= 50*time.Millisecond

const keepAlive 				= true
const keepAliveTime				= 400*time.Millisecond
const numberOfAllowedTimeouts 	= 3


func HandleConnection(	conn 							net.Conn,
						thisConnectsToNodeID 			uint8	,
						from_node_Ch 			<-chan 	[]byte	,
						to_node_Ch 				chan<- 	[]byte	,
						connection_Mutex_Ch		chan 	bool	) {

	var sendMsg		[]byte
	var sendErr		error

	var connBroke 			= false
	var connBrokeMsgSent 	= false

	var numberOfTimeouts uint8
	var keepAliveTicker *time.Ticker
	var keepAliveMessage = []byte{255,255,255,255,255}

	if keepAlive == true {
			keepAliveTicker = time.NewTicker(keepAliveTime)
			numberOfTimeouts = 0
		}

	forLoop:
	for {
		select {
		case sendMsg = <- from_node_Ch:
			_ 		= conn.SetWriteDeadline(time.Now().Add(writeDeadlineTime))
			_, sendErr	= conn.Write(sendMsg)
			if sendErr != nil {
				//fmt.Println(sendErr)
				connBroke = true
			}

			if connBroke == true && connBrokeMsgSent == true {
				head, _, _ := MessageFormat.Decode_msg(sendMsg)
				if head.MsgType == MessageFormat.NODE_DISCONNECTED && head.From == MessageFormat.NODE_COM {
					closeError := conn.Close()
					if closeError != nil {
						fmt.Println("Problem closing connection to node:", thisConnectsToNodeID)
					}
					break forLoop
				}
			}

		case <-connection_Mutex_Ch:
			receiveMsg := make([]byte, 1024)
			_ 		= conn.SetReadDeadline(time.Now().Add(readDeadlineTime))
			n, receiveErr	:= conn.Read(receiveMsg)

			if receiveErr == nil && n > 0 {
				if n >= 4 && bytes.Compare(receiveMsg[0:5], keepAliveMessage) == 0 {
					numberOfTimeouts = 0
				}else{
					to_node_Ch <- receiveMsg[0:n]
				}
			}else if e, ok := receiveErr.(net.Error); !ok || ( ok && !e.Timeout() ) {
				//fmt.Println(receiveErr)
				connBroke = true
			}

			if connBroke == true && connBrokeMsgSent == false {
				fmt.Println("A connection error accured")
				sendHeader := MessageFormat.MessageHeader_t{To: 		MessageFormat.NODE_COM			,
												FromNodeID: thisConnectsToNodeID			,
												MsgType: 	MessageFormat.NODE_DISCONNECTED	}
				connBrokeMsg, _ := MessageFormat.Encode_msg(sendHeader, "")
				to_node_Ch <- connBrokeMsg
				connBrokeMsgSent = true
			}
			connection_Mutex_Ch <- true

		case <- keepAliveTicker.C:
			//fmt.Println("Ticker ticked with numberOfTimeouts: ", numberOfTimeouts)
			if numberOfTimeouts >= numberOfAllowedTimeouts {
				connBroke = true
				keepAliveTicker.Stop()
			}else{
				_ 		= conn.SetWriteDeadline(time.Now().Add(writeDeadlineTime))
				_, sendErr	= conn.Write(keepAliveMessage)
				if sendErr != nil {
					//fmt.Println(sendErr)
					connBroke = true
				}
				numberOfTimeouts++
			}
			
		}
	}
}



