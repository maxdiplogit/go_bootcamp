package structs

import (
	"fmt"
)

type MaxTempCelsius float64
type MaxTempFarhenheit float64

type Vehicle struct {
	Company string
	FuelType string
	SerialNumber uint16
}

func (v *Vehicle) PrintVehicleData() {
	fmt.Printf("Company: %s\nFuelType: %s\nSerialNumber: %d\n", v.Company, (*v).FuelType, *(&v.SerialNumber))
}

func (v *Vehicle) UpdateVehicleData(company string, fuelType string, serialNumber uint16) {
	(*v).Company = company
	v.FuelType = fuelType
	(*v).SerialNumber = serialNumber
}
