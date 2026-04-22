package main

import (
	"go_revision/utils"
	"go_revision/structs"
	"go_revision/arrays"
	"go_revision/slices"
	"go_revision/make_practice"
	"go_revision/maps"

	"fmt"
	"time"
	"encoding/json"
)

func main() {
	var x, y int

	fmt.Print("Enter x: ")
	fmt.Scanf("%d", &x)	

	fmt.Print("Enter y: ")
	fmt.Scanf("%d", &y)

	var res int = utils.AddTwoNumbers(x, y)
	
	fmt.Printf("Result: %d\n", res)

	utils.PrintStars()
	utils.PrintDay()

	var a, b int = utils.ReturnMultipleValues()
	fmt.Printf("a = %d; b = %d\n", a, b)

	var vehicle_0 *structs.Vehicle;
	vehicle_0 = &structs.Vehicle{
		Company: "BMW",
		FuelType: "Petrol",
		SerialNumber: 12,
	}

	vehicle_0.PrintVehicleData()
	vehicle_0.UpdateVehicleData("Audi", "Diesel", 65535)	
	vehicle_0.PrintVehicleData()

	arrays.ArrayPractice()

	slices.SlicesPractice()

	make_practice.MakeExample()

	maps.MapsExample()

	truck_0 := &structs.Truck{
		EngineModel: "XA31",
		Tyres: 6,
		Pistons: 8,
	}
	truck_0.TruckJsonExample()

	fmt.Println()

	userIdentifier := &structs.Identifier{
		ID: 1,
		UUID: "user-abc",
	}

	u := structs.User{
	        UserIdentifier: userIdentifier,
        	Timestamps: structs.Timestamps{
            		CreatedAt: time.Now().Add(-24 * time.Hour),
            		UpdatedAt: time.Now(),
        	},
        	AuditTrail: structs.AuditTrail{
            		CreatedBy: 99,
            		UpdatedBy: 99,
            		Version:   2,
        	},
       		Email:    "asha@example.com",
        	FullName: "Asha Rao",
    	}
    
    	// Promoted access — clean and ergonomic:
    	fmt.Println(u.UserIdentifier)         // from Identifier
    	fmt.Println(u.CreatedAt)  // from Timestamps
    	fmt.Println(u.Version)    // from AuditTrail
    	fmt.Println(u.Email)      // User's own field

   	 // JSON output is flat — looks like one coherent object:
    	data, _ := json.MarshalIndent(u, "", "  ")
    	fmt.Println(string(data))
}
