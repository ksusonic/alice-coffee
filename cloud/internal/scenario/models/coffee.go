package models

import "encoding/json"

type CoffeeType struct {
	Value         string
	HumanReadable string
}

var (
	UnknownCoffee = CoffeeType{"unknown", ""}
	Espresso      = CoffeeType{"espresso", "эспрессо"}
	Americano     = CoffeeType{"americano", "американо"}
	Cappuccino    = CoffeeType{"cappuccino", "капучино"}
	Latte         = CoffeeType{"latte", "латте"}
	Tea           = CoffeeType{"tea", "чай"}

	AllCoffeeTypes = []CoffeeType{
		Espresso, Americano, Cappuccino, Latte, Tea,
	}
)

func (t CoffeeType) String() string {
	return t.Value
}

func (t CoffeeType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func ParseCoffee(coffeeType string) *CoffeeType {
	for i := range AllCoffeeTypes {
		if AllCoffeeTypes[i].Value == coffeeType {
			return &AllCoffeeTypes[i]
		}
	}
	return nil
}
