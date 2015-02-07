package smoother

import "fmt"

type ReturnInfo struct {
	Symbol string `json:"symbol"`
	Current string `json:"current"`
	Tomorrow string `json:"tomorrow"`
	AmountIncrease string `json:"amount_increase"`
	PercentIncrease string `json:"percent_increase"`
	MeanSquaredError string `json:"mse"`
}

func Smooth() {
	fmt.Println("Goodbye world")
}
