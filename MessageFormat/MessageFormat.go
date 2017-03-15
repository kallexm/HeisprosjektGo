package MessageFormat
/*
||	File: MessageFormat.go
||
||	Authors: 
||
||	Date: 	 Spring 2017
||	Course:  TTK4145 - Real-time Programming, NTNU
||	
||	Summary of File:
||		Defines the message format which should be used on the internal network.
||		Contains functions to encode and decode messages to be sent on the network.
||
*/

import
(
	"fmt"
)

type Address_t uint8
const (
	MASTER Address_t = iota
	BACKUP
	ELEVATOR
	NODE_COM
	ORDER_DIST
)



type MsgType_t uint8
const (
	BACKUP_DATA_TRANSFER	= iota
	NEW_ORDER_TO_ELEVATOR			//Order stuct fra ElevatorControl/ElevatorStatus
	SET_LIGHT						//ButtonPlacement struc fra ElevatorControl/ElvatorDriver, vurder nytt nav siden den kan skru av og på 
	CLEAR_LIGHT
	ELEVATOR_STATUS_DATA			//position struct, fra ElevatorControl/ElevatorStatus 
	NEW_ELEVATOR_REQUEST			//Order struct fra ElevatorControl/ElevatorStatus 
	ORDER_FINISHED_BY_ELEVATOR		//position struct, fra ElevatorControl/ElevatorStatus
	NEW_ELEVATOR_REQUEST_ACCEPTED 
	NODE_DISCONNECTED
	NODE_CONNECTED
	CHANGE_TO_MASTER
	CHANGE_TO_SLAVE
	MASTER_ON_NET
	MASTER_NOT_ON_NET
	MERGE_ORDERS_REQUEST
)

type MessageHeader_t struct {
	To 			Address_t
	ToNodeID 	uint8
	From 		Address_t
	FromNodeID 	uint8
	MsgType 	MsgType_t
}




func Encode_msg(msgHead MessageHeader_t, data []byte) ([]byte, error) {
	var msg []byte
	var err error
	err = nil
	
	msg = append(msg, byte(msgHead.To), byte(msgHead.ToNodeID), byte(msgHead.From), byte(msgHead.FromNodeID), byte(msgHead.MsgType))
	msg = append(msg, data...)
	return msg, err
}




func Decode_msg(msg []byte) (MessageHeader_t, []byte, error) {
	var msgHead MessageHeader_t
	var data []byte
	var err error
	err = nil
	if len(msg) < 5 {
		err = fmt.Errorf("Error: Message to short to decode. Message was: %v", msg)
		return msgHead, data, err
	}
	
	msgHead.To 			= Address_t(msg[0])
	msgHead.ToNodeID 	= uint8(msg[1])
	msgHead.From 		= Address_t(msg[2])
	msgHead.FromNodeID  = uint8(msg[3])
	msgHead.MsgType 	= MsgType_t(msg[4])
	data 				= msg[5:]
	
	return msgHead, data, err
}