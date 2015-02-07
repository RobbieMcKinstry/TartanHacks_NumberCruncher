package smoother

import "fmt"
import "strconv"

type ReturnInfo struct {
	Current          string `json:"current"`
	Tomorrow         string `json:"tomorrow"`
	AmountIncrease   string `json:"amount_increase"`
	PercentIncrease  string `json:"percent_increase"`
	MeanSquaredError string `json:"mse"`
}

func Smooth(input []float64) *ReturnInfo {
	fmt.Println("Goodbye world")

	history := &History{Data: input}
	_ = history
	alpha := 0.5
	beta := 0.5
	//alpha, beta := OptimizeAB(history)
	next, mse := history.forecast(alpha, beta)
	fmt.Println(next)
	fmt.Println(mse)

	today := input[len(input)-1]
	return &ReturnInfo{
		Current:          strconv.FormatFloat(today, 'f', 2, 64),
		Tomorrow:         strconv.FormatFloat(next, 'f', 2, 64),
		AmountIncrease:   strconv.FormatFloat(next-today, 'f', 2, 64),
		PercentIncrease:  strconv.FormatFloat(next/today-1, 'f', 2, 64),
		MeanSquaredError: strconv.FormatFloat(mse, 'f', 2, 64),
	}
}
