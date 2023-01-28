package scenario

import requestpb "github.com/ksusonic/alice-coffee/protos/request"

func coffeeTypeToEnum(value string) requestpb.MakeCoffeeAction_CoffeeEnum {
	switch value {
	case "espresso":
		return requestpb.MakeCoffeeAction_ESPRESSO
	case "americano":
		return requestpb.MakeCoffeeAction_AMERICANO
	case "cappuccino":
		return requestpb.MakeCoffeeAction_CAPPUCCINO
	case "latte":
		return requestpb.MakeCoffeeAction_LATTE
	case "cacao":
		return requestpb.MakeCoffeeAction_CACAO
	}
	return requestpb.MakeCoffeeAction_SIMPLE
}

func coffeeTypeToHumanReadable(value string) string {
	switch value {
	case "espresso":
		return "эспрессо"
	case "americano":
		return "американо"
	case "cappuccino":
		return "капучино"
	case "latte":
		return "латте"
	case "cacao":
		return "какао"
	}
	return "кофе"
}
