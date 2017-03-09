package test

import(
	"fmt"
	"encoding/json"
)


type TestStruct struct {
	Foo string
	Bar string
}


func PrintTh(PringCh <-chan []byte){
	b := <- PringCh
	var temp TestStruct
	err := json.Unmarshal(b,&temp)
	if err != nil{
		fmt.Println("Error: ", err)
	}
	fmt.Println(temp) 

}