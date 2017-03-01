package MessageFormat

import
(
	//"errors"
)

type Address_t uint8
const (
	MASTER Address_t = iota
	BACKUP
	ELEVATOR
	NODE_COM
)


type AddressID_t uint8


type MsgType_t uint8
const (
	HEARTHBEAT = iota
	BACKUP_DATA_TRANSFER
	NEW_ORDER_TO_ELEVATOR
	SET_LIGHT
	CLEAR_LIGHT
	ELEVATOR_STATUS_DATA
	NEW_ELEVATOR_REQUEST
	ORDER_FINISHED_BY_ELEVATOR
	PING
)

type MessageHeader_t struct {
	to 			Address_t
	toNodeID 	AddressID_t
	from 		Address_t
	fromNodeID 	AddressID_t
	msgType 	MsgType_t
}

func Gen_Msg(msgHead MessageHeader_t, data string) ([]byte, error) {
	var msg []byte
	var err error
	err = nil
	
	msg = append(msg, byte(msgHead.to), byte(msgHead.toNodeID), byte(msgHead.from), byte(msgHead.fromNodeID), byte(msgHead.msgType))
	msg = append(msg, []byte(data)...)
	return msg, err
}


func De_Gen_Msg(msg []byte) (MessageHeader_t, string, error) {
	var msgHead MessageHeader_t
	var data string
	var err error
	err = nil
	
	msgHead.to 			= Address_t(msg[0])
	msgHead.toNodeID 	= AddressID_t(msg[1])
	msgHead.from 		= Address_t(msg[2])
	msgHead.fromNodeID  = AddressID_t(msg[3])
	msgHead.msgType 	= MsgType_t(msg[4])
	data 				= string(msg[5:])
	
	return msgHead, data, err
}