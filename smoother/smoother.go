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

func Smooth(input []int) *ReturnInfo {
	fmt.Println("Goodbye world")
	return &ReturnInfo{
		Symbol: "AAPL",
		Current: "100",
		Tomorrow: "110",
		AmountIncrease: "10",
		PercentIncrease: "10"
		MeanSquaredError: "2.5"
	}
}
