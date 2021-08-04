package service

import "wq-fotune-backend/pkg/evaluation"

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
