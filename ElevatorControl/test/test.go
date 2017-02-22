package main

import
(
	"fmt"
	"errors"
)
//var innerMap = initInnerMap //map[string]int
var OuterMap = initOuterMap() //map[int]map[string]int

func initOuterMap() map[string]map[int]int{
	up := map[string]map[int]int{
		"Up":map[int]int{1:0x2000,2:0x2001},
		"Down":map[int]int{1:0x2000,2:0x2001},
		"Comand":map[int]int{1:0x2003,2:0x004}}
	return up
}

func getButtonChannle(floor int, button string) (int, error){
	if OuterMap[button][floor] == 0{
		return 0, errors.New("Index out of bounds error: 0001")
	}
	return OuterMap[button][floor], nil
}


func main() {
	/*innerMap = make(map[string]int)
	OuterMap = make(map[int]map[string]int)
	innerMap["hei"] = 1
	innerMap["hade"] = 2
	OuterMap[1] = innerMap
	fmt.Println("OuterMap: ", OuterMap)
	fmt.Println("test hei", OuterMap[1]["hei"])
	*/
	hex, err := getButtonChannle(5,"hei")
	if err != nil{
		fmt.Println("Error: ", err)
	}else{
		fmt.Println("value: ", hex)
	}
	fmt.Println(OuterMap)
}
