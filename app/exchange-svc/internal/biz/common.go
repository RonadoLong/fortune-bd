package biz

import "fortune-bd/app/exchange-svc/internal/evaluation"

func getEvaDirection(direction string) string {
	switch direction {
	case "buy":
		return evaluation.CloseSell
	case "sell":
		return evaluation.CloseBuy
	}
	return ""
}

func getPosDirection(direction string) string {
	switch direction {
	case "buy":
		return "sell"
	case "sell":
		return "buy"
	}
	return ""
}
