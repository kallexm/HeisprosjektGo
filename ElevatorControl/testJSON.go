package main 

import(
	"./testTh"
	"encoding/json"
	"fmt"
)


func main() {
	printCh := make(chan []byte)
	temp := test.TestStruct{Foo:"1",Bar:"Hello World!"}
	/*b, err := json.Marshal(temp)
	if err != nil{
		fmt.Println("error: ", err)
	}*/
	b := generateMsg(temp)
	go test.PrintTh(printCh)
	printCh <- b

}


/*func main() {
	var jsonBlob = []byte(`[
		{"Name": "Platypus", "Order": "Monotremata"},
		{"Name": "Quoll",    "Order": "Dasyuromorphia"}
	]`)
	type Animal struct {
		Name  string
		Order string
	}
	var animals []Animal
	err := json.Unmarshal(jsonBlob, &animals)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v", animals)
}*/


func generateMsg(v interface{}) []byte{
	b, err := json.Marshal(v)
	if err != nil{
		fmt.Printf("Error i generateMsg", err)
	}
	return b
}