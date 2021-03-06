package NodeMessageRelay
/*
||	File: NodeMessageRelay.go
||
||	Authors:  
||
||	Date: 	 Spring 2017
||	Course:  TTK4145 - Real-time Programming, NTNU
||	
||	Summary of File:
||		Contains a goroutine function which gets messages from channels and
||		forwards them based on lookup in the routing table in NodeRoutingTable
||		which is updated by NodeConnectionManager when necessary.
||
*/

import
(
	"../NodeRoutingTable"
	"../../MessageFormat"
	
	"fmt"
)


var routingTable_ptr *NodeRoutingTable.RoutingTable_t
var masterShouldTakeMutexNext = false




func Thread (routingTable_Ch chan *NodeRoutingTable.RoutingTable_t) {
	for {
		if masterShouldTakeMutexNext == false {
			routingTable_ptr = <- routingTable_Ch
		}

		for _, tableEntry := range (*routingTable_ptr) {
			if masterShouldTakeMutexNext == true {
				
				if tableEntry.IsMaster {
					masterShouldTakeMutexNext = false
				}else{
					continue
				}
			}
			select {
			case tableEntry.Mutex_Ch <- true:
				continueFor := true
				for continueFor {
					select {
					case receivedMsg := <- tableEntry.Receive_Ch:
						msgHeader, data, err := MessageFormat.Decode_msg(receivedMsg)
						if err == nil {
							for j, searchTableEntry := range (*routingTable_ptr) {
								//-----------------------------------------------------------------
								// Routing Entries
								if msgHeader.To == MessageFormat.MASTER && searchTableEntry.IsMaster == true {

									if msgHeader.From == MessageFormat.ELEVATOR && msgHeader.FromNodeID == 0 {
										msgHeader.FromNodeID = tableEntry.NodeID
										sendMsg, _ := MessageFormat.Encode_msg(msgHeader, data)
										searchTableEntry.Send_Ch <- sendMsg
										masterShouldTakeMutexNext = true
										break
									}else{
										searchTableEntry.Send_Ch <- receivedMsg
										break
									}
									
								}else if msgHeader.To == MessageFormat.ELEVATOR && (searchTableEntry.IsElev == true || searchTableEntry.IsExtern == true) && msgHeader.ToNodeID == searchTableEntry.NodeID {
									searchTableEntry.Send_Ch <- receivedMsg
									break
									
								}else if msgHeader.To == MessageFormat.BACKUP && (searchTableEntry.IsBackup == true || searchTableEntry.IsExtern == true) {
									searchTableEntry.Send_Ch <- receivedMsg
									break
									
								}else if msgHeader.To == MessageFormat.NODE_COM && searchTableEntry.IsNet == true {
									searchTableEntry.Send_Ch <- receivedMsg
									break
									
								}else if msgHeader.To == MessageFormat.ORDER_DIST && searchTableEntry.IsOrderDist == true {
									searchTableEntry.Send_Ch <- receivedMsg
									break
									
								}else if j == len(*routingTable_ptr)-1{
									// Drop package if it doesn't match any of the above filters/masks
									fmt.Println("Package with header", msgHeader, "dropped because of no matching filters!")
								} 
								//-----------------------------------------------------------------
							}
						}
						
					default:
					}
					select {
					case <- tableEntry.Mutex_Ch:
						continueFor = false
					default:
					}
				}
			default:
			}
			
		}
		if masterShouldTakeMutexNext == false {
			routingTable_Ch <- routingTable_ptr
			routingTable_ptr = nil
		}
		
	}
}




func eval_error(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

