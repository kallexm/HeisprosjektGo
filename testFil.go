package main

import
(
	"fmt"
	//"time"
	"errors"
)

type durp struct {
	str string
	integer int
	boolean bool
}

func main() {
	fmt.Println("Hello Sunshine!")
	
	data := "...durp!"
	mh := MessageHeader{sendTo: "MASTER", sendToNodeID: 0, sendFrom: "ELEVATOR", sendFromNodeID: 65, msgType: "NEW_ELEVATOR_REQUEST"}
	fmt.Printf("%+v%s", mh, "\n")
	
	ba1, err := GenerateMessage(mh, data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ba1)
	fmt.Println(string(ba1))
	
	mh1, data1, err1 := DegenerateMessage(ba1)
	fmt.Printf("%+v%s", mh1, "\n")
	fmt.Println(data1)
	fmt.Println(err1)
}


type MessageHeader struct {
	sendTo string
	sendToNodeID uint8
	sendFrom string
	sendFromNodeID uint8
	msgType string
}


func GenerateMessage(msgHead MessageHeader, data string) ([]byte, error) {
	var msg []byte
	var err error
	err = nil
	
	switch(msgHead.sendFrom) {
	case "MASTER":
		switch(msgHead.sendTo) {
		case "BACKUP":
			switch(msgHead.msgType) {
			case "HEARTHBEAT":
				msg, err = genMsg("B", 0, "M", 0, "HEARTH", data)
			case "BACKUP_DATA_TRANSFER":
				msg, err = genMsg("B", 0, "M", 0, "BKUPDA", data)
			default:
				err = errors.New("Error 001")
			}

		case "ELEVATOR":
			switch(msgHead.msgType) {
			case "NEW_ORDER_TO_ELEVATOR":
				msg, err = genMsg("E", msgHead.sendToNodeID, "M", 0, "ORTOEL", data)
			case "SET_LIGHT":
				msg, err = genMsg("E", msgHead.sendToNodeID, "M", 0, "SETLIG", data)
			case "CLEAR_LIGHT":
				msg, err = genMsg("E", msgHead.sendToNodeID, "M", 0, "CLRLIG", data)
			default:
				err = errors.New("Error 002")
			}

		case "NODE_COM":
			switch(msgHead.msgType) {
			default:
				err = errors.New("Error 003")
			}
	
		default:
			err = errors.New("Error 004")
		}
		
	case "ELEVATOR":
		switch(msgHead.sendTo) {
		case "MASTER":
			switch(msgHead.msgType) {
			case "ELEVATOR_STATUS_DATA":
				msg, err = genMsg("M", 0, "E", msgHead.sendFromNodeID, "ELSTDA", data)
			case "NEW_ELEVATOR_REQUEST":
				msg, err = genMsg("M", 0, "E", msgHead.sendFromNodeID, "NEWREQ", data)
			case "ORDER_FINISHED_BY_ELEVATOR":
				msg, err = genMsg("M", 0, "E", msgHead.sendFromNodeID, "ORDFIN", data)
			default:
				err = errors.New("Error 005")
			}
			
		default:
			err = errors.New("Error 006")
		}
		
	case "NODE_COM":
		switch(msgHead.sendTo) {
		case "NODE_COM":
			switch(msgHead.msgType) {
			case "PING":
				msg, err = genMsg("N", msgHead.sendToNodeID, "N", msgHead.sendFromNodeID, "PINGNC", data)
			default:
				err = errors.New("Error 007")
			}
			
		default:
			err = errors.New("Error 008")
		}
		
	default:
		err = errors.New("Error 009")
	}
	
	return msg, err
}


func genMsg(to string, toNodeID uint8, from string, fromNodeID uint8, msgType string, data string) ([]byte, error) {
	var err error
	err = nil
	message := []byte(to)
	message = append(message, toNodeID)
	message = append(message, []byte(from)...)
	message = append(message, fromNodeID)
	message = append(message, []byte(msgType+data)...)
	return message, err
}

func DeGenMsg(message []byte) (MessageHeader, string, error) {
	var msgHead MessageHeader
	var err error
	var data string
	msgHead.sendTo = string(message[0])
	msgHead.sendToNodeID = message[1]
	msgHead.sendFrom = string(message[2])
	msgHead.sendFromNodeID = message[3]
	msgHead.msgType = string(message[4:10])
	data = string(message[10:])
	
	return msgHead, data, err
}

