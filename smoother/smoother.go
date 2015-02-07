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

func Smooth(input []float64) *ReturnInfo {
	fmt.Println("Goodbye world")

	history := &History{Data: input}
	alpha, beta := OptimizeAB(history)
	next, mse := history.forecast(alpha, beta)
	fmt.Println(next)
	fmt.Println(mse)

	return &ReturnInfo{
		Symbol: "AAPL",
		Current: "100",
		Tomorrow: "110",
		AmountIncrease: "10",
		PercentIncrease: "10",
		MeanSquaredError: "2.5",
	}
}
