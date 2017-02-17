package OrderDistributerThread

import
(
	//"fmt"
)



func OD_thr(NC_to_OD_Ch <-chan []byte, OD_to_NC_Ch chan<- []byte, OD_exit_Ch chan<- bool) {
	byteArray := <- NC_to_OD_Ch
	str := string(byteArray)
	str = str+"...durp!"
	byteArray = []byte(str)
	OD_to_NC_Ch <- byteArray
	
	OD_exit_Ch <- true
}