package models

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
	Cacao         = CoffeeType{"cacao", "какао"}

	AllCoffeeTypes = []CoffeeType{
		Espresso, Americano, Cappuccino, Latte, Cacao,
	}
)

func (t CoffeeType) String() string {
	return t.Value
}

func ParseCoffee(coffeeType string) *CoffeeType {
	for i := range AllCoffeeTypes {
		if AllCoffeeTypes[i].Value == coffeeType {
			return &AllCoffeeTypes[i]
		}
	}
	return nil
}
