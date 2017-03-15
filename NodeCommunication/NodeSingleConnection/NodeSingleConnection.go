package NodeSingleConnection
/*
||	File: NodeSingleConnection
||
||	Authors: 
||			 
||	Date: 	 Spring 2017
||	Course:  TTK4145 - Real-time Programming, NTNU
||	
||	Summary of File:
||		The function HandleConnection(...) should be used as a goroutine.
||		The function takes a TCP connection, an input and an output channel and passes
||		messages between these and sends notification to NodeConnectionManager if the connection breaks.
||
*/

import
(
	"../../MessageFormat"
	
	"fmt"
	"net"
	"time"
	"bytes"
)

const readDeadlineTime 			= 50*time.Millisecond
const writeDeadlineTime 		= 50*time.Millisecond

const keepAlive 				= true
const keepAliveTime				= 700*time.Millisecond
const numberOfAllowedTimeouts 	= 3

const readBufferSize 			= 1024



func HandleConnection(	conn 							net.Conn,
						thisConnectsToNodeID 			uint8	,
						from_node_Ch 			<-chan 	[]byte	,
						to_node_Ch 				chan<- 	[]byte	,
						connection_Mutex_Ch		chan 	bool	) {


	var sendErr						error
	var numberOfTimeouts 			uint8
	var keepAliveTicker 			*time.Ticker

	var connBroke 					= false
	var connBrokeMsgSent 			= false

	var keepAliveMessage 			= []byte{255,255,255,255,255}
	var lengthOfKeepAliveMessage 	= len(keepAliveMessage)

	if keepAlive == true {
			keepAliveTicker = time.NewTicker(keepAliveTime)
			numberOfTimeouts = 0
		}

	forLoop:
	for {
		select {
		case sendMsg := <- from_node_Ch:
			_ 			= conn.SetWriteDeadline(time.Now().Add(writeDeadlineTime))
			_, sendErr	= conn.Write(sendMsg)
			if sendErr != nil {
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
			receiveMsg := make([]byte, readBufferSize)
			_ 		= conn.SetReadDeadline(time.Now().Add(readDeadlineTime))
			n, receiveErr	:= conn.Read(receiveMsg)

			if receiveErr == nil && n > 0 {
				m := n
				for m >= lengthOfKeepAliveMessage && bytes.Compare(receiveMsg[n-m:lengthOfKeepAliveMessage], keepAliveMessage) == 0 {
					numberOfTimeouts = 0
					m = m - lengthOfKeepAliveMessage
				}
				i := n
				for i >= lengthOfKeepAliveMessage && bytes.Compare(receiveMsg[i-lengthOfKeepAliveMessage:i], keepAliveMessage) == 0 {
					numberOfTimeouts = 0
					i = i - lengthOfKeepAliveMessage
				}
				if i-(n-m) > 0 {
					to_node_Ch <- receiveMsg[(n-m):i]
				}
			
			}else if e, ok := receiveErr.(net.Error); !ok || ( ok && !e.Timeout() ) {
				connBroke = true
			}

			if connBroke == true && connBrokeMsgSent == false {
				fmt.Println("A connection error accured")
				sendHeader := MessageFormat.MessageHeader_t{To: 		MessageFormat.NODE_COM			,
												FromNodeID: thisConnectsToNodeID			,
												MsgType: 	MessageFormat.NODE_DISCONNECTED	}
				connBrokeMsg, _ := MessageFormat.Encode_msg(sendHeader, []byte(""))
				to_node_Ch <- connBrokeMsg
				connBrokeMsgSent = true
			}
			connection_Mutex_Ch <- true

		case <- keepAliveTicker.C:
			if numberOfTimeouts >= numberOfAllowedTimeouts {
				connBroke = true
				keepAliveTicker.Stop()
			}else{
				_ 		= conn.SetWriteDeadline(time.Now().Add(writeDeadlineTime))
				_, sendErr	= conn.Write(keepAliveMessage)
				if sendErr != nil {
					connBroke = true
				}
				numberOfTimeouts++
			}
		}
	}
}



