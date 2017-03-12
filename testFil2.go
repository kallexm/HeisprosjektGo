package main

import
(
	"fmt"
	"time"
	//"net"
	"math/rand"
)

const nodeID = 34



func main(){
	rand.Seed(nodeID*int64(time.Now().Second()))
	
	for {
		randomNumber := (upperRandomValue - lowerRandomValue)*rand.Float32() + lowerRandomValue
		fmt.Println(randomNumber)
		time.Sleep(1*time.Second)
	}

}


func getNanoSecTime() int64 {
	return (time.Now().UnixNano() - (time.Now().UnixNano()/100000)*100000)
}