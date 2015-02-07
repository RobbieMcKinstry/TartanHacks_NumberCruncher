package smoother

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
)

// a wrapper for a slice of float64s, with a cached current value, but remember to use `json`
type History struct {
	Data        []float64 `json: "data"`
}

type Result struct {
	Alpha float64
	Beta  float64
	Mse   float64
}



// this function reads the json string in from the request, before passing it to the NewHistory func
func HandleComputationRequest(w http.ResponseWriter, r *http.Request) {
	jsonString := r.FormValue("data")
	var h *History = NewHistory(jsonString)
	alpha, beta := OptimizeAB(h)
	next, mse := h.forecast(alpha, beta)
	fmt.Println(next)
	fmt.Println(mse)
	// TODO: make the request back to the python server.
}

// this function parses the string passed to it as a history object, and returns that object.
// remember to use the `` notatoin to specify which elements map to which from json to Go structs.
func NewHistory(myjson string) *History {
	var h *History = new(History)
	if err := json.Unmarshal([]byte(myjson), h); err != nil {
		log.Fatal(err)
	}
	return h
}

// this function accepts a pointer to a History object, and returns the optimum values of alpha and beta.
func OptimizeAB(history *History) (alpha, beta float64) {

	var (
		coverage    float64      = 0.2 // value of how frequently we run simulated annealing
		resultsChan chan *Result = make(chan *Result, 30)
		best        *Result
	)

	// alternative code:
	/*
		alpha = .5
		beta = .5
		simulatedAnnealing(alpha, beta, history)
	*/

	for alpha = 0.0; alpha < 1.0; alpha += coverage {
		for beta = 0.0; beta < 1.0; beta += coverage {
			go func() {
				var r *Result = simulatedAnnealing(alpha, beta, history) // put the simulated annealing algo here.
				resultsChan <- r
			}()
		}
	}

	close(resultsChan)

	result, ok := <-resultsChan
	best = result

	// block on results chan
	for ok {
		result, ok = <-resultsChan
		if result.Mse < best.Mse {
			best = result
		}
	}

	alpha = best.Alpha
	beta = best.Beta
	return alpha, beta
}

func simulatedAnnealing(alpha, beta float64, h *History) *Result {

	var (
		f                     = calculateMSE
		coolingRate   float64 = 0.03  //  can be augmented/optimized -- currently arbitrary
		temperature   float64 = 10000 // can be augmented/optimized -- currently arbitrary
		currentEnergy float64         // the fitness metric
		energy        float64         // the new metric
		delta         float64 = 0.1
	)

	observed := h.Data
	expected := h.forecastedArray(alpha, beta)
	currentEnergy = f(observed, expected)

	for ; temperature > 1.0; temperature *= 1 - coolingRate {

		/// how do i define neighbor??? Do i look at the location directly next to where I'm starting, and call that a neighbor? or do I hill climb, then have another neighbor based on a range?
		/*
		 * This is what I do. I hill climb from where I currently am. Easy.
		 * Enter loop:
		 * 	Take current - delta and calculate acceptance chance.
		 * 	Accept or reject probabilty1
		 * 	Take current + delta and calculate acceptance chance.
		 *      Accept or reject probabitity2.
		 * 	Hill Climb the rest of the way.
		 */

		// generate a number between 1 and 8 inclusive.
		// switch on that number

		gen := rand.Intn(8)
		switch gen {

		case 0:
			expected = h.forecastedArray(alpha+delta, beta)
			energy = f(observed, expected)
			prob := acceptanceProbability(currentEnergy, energy, temperature)
			if prob > rand.Float64() {
				alpha = alpha - delta
				currentEnergy = energy
			}

		case 1:
			expected = h.forecastedArray(alpha-delta, beta)
			energy = f(observed, expected)
			prob := acceptanceProbability(currentEnergy, energy, temperature)
			if prob > rand.Float64() {
				alpha = alpha - delta
				currentEnergy = energy
			}

		case 2:
			expected = h.forecastedArray(alpha, beta+delta)
			energy = f(observed, expected)
			prob := acceptanceProbability(currentEnergy, energy, temperature)
			if prob > rand.Float64() {
				beta = beta + delta
				currentEnergy = energy
			}

		case 3:
			expected = h.forecastedArray(alpha, beta-delta)
			energy = f(observed, expected)
			prob := acceptanceProbability(currentEnergy, energy, temperature)
			if prob > rand.Float64() {
				beta = beta - delta
				currentEnergy = energy
			}

		case 4:
			expected = h.forecastedArray(alpha+delta, beta+delta)
			energy = f(observed, expected)
			prob := acceptanceProbability(currentEnergy, energy, temperature)
			if prob > rand.Float64() {
				alpha = alpha + delta
				beta = beta + delta
				currentEnergy = energy
			}

		case 5:
			expected = h.forecastedArray(alpha+delta, beta-delta)
			energy = f(observed, expected)
			prob := acceptanceProbability(currentEnergy, energy, temperature)
			if prob > rand.Float64() {
				alpha = alpha + delta
				beta = beta - delta
				currentEnergy = energy
			}

		case 6:
			expected = h.forecastedArray(alpha-delta, beta+delta)
			energy = f(observed, expected)
			prob := acceptanceProbability(currentEnergy, energy, temperature)
			if prob > rand.Float64() {
				alpha = alpha - delta
				beta = beta + delta
				currentEnergy = energy
			}

		case 7:
			expected = h.forecastedArray(alpha-delta, beta-delta)
			energy = f(observed, expected)
			prob := acceptanceProbability(currentEnergy, energy, temperature)
			if prob > rand.Float64() {
				alpha = alpha - delta
				beta = beta - delta
				currentEnergy = energy
			}
		} // end switch

		delta *= 1 - coolingRate

	} // end loop

	n := new(Result)
	n.Mse = currentEnergy
	n.Alpha = alpha
	n.Beta = beta
	return n

} // end function

// returns the probability that you should accept the given temperature.
func acceptanceProbability(energy, newEnergy, temperature float64) float64 {
	result := 0.0
	if newEnergy < energy {
		result = 1.0
	} else {
		result = math.Exp((energy - newEnergy) / temperature)
	}
	return result
}

// this METHOD will return the forecasted value for tomorrow, and the mean squared error received when calculating the forecasted value.
// NOTE: if there is ever a reason why we would need the forecast without needing the mse, then mse should not be called from within this function, since its O(n) on the data set.
func (me *History) forecast(alpha, beta float64) (next, mse float64) {
	a := me.forecastedArray(alpha, beta)
	mse = calculateMSE(me.Data, a)
	next = a[len(a)-1]
	return next, mse
}

// this function returns the percent change from the current price to the forecasted price.
func percentChange(current, forecasted float64) float64 {
	return 0.0
}

// this function takes two slices of floats and returns the mse between them.
func calculateMSE(s1, s2 []float64) float64 {
	if len(s2) > len(s1) {
		temp := s1
		s2 = s1
		s1 = temp
	}

	mse := 0.0
	for i, _ := range s1 {
		statError := math.Abs(s1[i] - s2[i])
		mse += statError * statError
	}
	return mse
}

// returns the forecasted array.
func (me *History) forecastedArray(alpha, beta float64) []float64 {
	observed := me.Data
	S := make([]float64, len(observed))
	B := make([]float64, len(observed))

	size := len(observed)
	S[0] = me.Data[0]

	for i := 1; i < size; i++ {
		temp1 := alpha * observed[i]
		temp2 := (1 - alpha)
		temp3 := S[i-1]
		temp4 := B[i-1]
		S[i] = temp1 + temp2*(temp3+temp4)
		S[i] = alpha*observed[i] + (1-alpha)*(S[i-1]+B[i-1])
		// S[i] = alpha * observed[i] + (1−alpha) * (S[i-1] + B[i−1])	// forecast = (alpha * today's observed) + (1 -alpha) * (yesterday's forecast + yesterday's slope)
		B[i] = beta*(S[i]-S[i-1]) + (1-beta)*B[i-1] // deltaSlope = beta times (the change in forecast) + (1 - beta)  times (the last deltaSlope)
	}
	S[size] = S[size-1] + B[size-1] // tomorrow's forecast is yestday's plus yesterday's close
	return S
}
