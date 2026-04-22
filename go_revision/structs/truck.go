package structs

import (
	"encoding/json"
	"fmt"
)

type Truck struct {
	EngineModel string	`json:"engine_model"`
	Tyres uint8		`json:"tyres,string"`		// encode a number or bool as a JSON string 
	Pistons uint8		`json:"-"`			// '-' never include this field in JSON
}

func (truck *Truck) TruckJsonExample() {
	jsonMarshalByteData, ok := json.Marshal(truck)
	
	if ok != nil {
		fmt.Printf("Could not marshal truck into JSON: %q\n", ok)
	}

	fmt.Printf("Marshaled JSON: %#s\n", string(jsonMarshalByteData))

	jsonData := []byte(`{"engine_model": "XAE56", "tyres": "4", "pistons": 3}`)
	var truck_0 Truck
	json.Unmarshal(jsonData, &truck_0)

	fmt.Printf("Truck struct from json: %v\n", truck_0)
}
