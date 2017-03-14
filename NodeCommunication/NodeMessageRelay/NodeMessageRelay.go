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
		if masterShouldTakeMutexNext == false { // remove for if no work
			routingTable_ptr = <- routingTable_Ch
		}
		
		//fmt.Printf("%+v", *routingTable_ptr)
		//fmt.Println()


		for _, tableEntry := range (*routingTable_ptr) {
			//-------------------------------------- remove if no work
			if masterShouldTakeMutexNext == true {
				fmt.Println("I am in the second if")
				if tableEntry.IsMaster {
					masterShouldTakeMutexNext = false
				}else{
					continue
				}
			}
			//-------------------------------------- remove if no work
			select {
			case tableEntry.Mutex_Ch <- true:
				continueFor := true
				//fmt.Println("masterShouldTakeMutexNext: ", masterShouldTakeMutexNext)
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
										fmt.Println("I am in if")
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
								} else {
									fmt.Println("We are in else")
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
		if masterShouldTakeMutexNext == false { // remove for if no work
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

