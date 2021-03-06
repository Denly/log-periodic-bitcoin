package worker

import (
	"fmt"
	"log-periodic-bitcoin/config"
	"log-periodic-bitcoin/models"
	"log-periodic-bitcoin/regression"
	"math"
	"os"
	"time"
)

type Result struct {
	N       int64
	J       float64
	Params  *regression.Parameters
	ExeTime int64
}

func (r *Result) WriteResults(f *os.File) error {
	_, err := f.WriteString(fmt.Sprintf(config.ResultFormat, r.J, r.Params.A, r.Params.B, r.Params.Tc, r.Params.Beta, r.Params.C, r.Params.Omega, r.Params.Phi))
	if err != nil {
		return err
	}
	return nil
}

func OpenResultFile(mode, dataFilename string, minDate, maxDate float64) (*os.File, error) {
	f, err := os.OpenFile(config.ResultsFilePrefix+mode+config.ResultFileSufix, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	_, err = f.WriteString(fmt.Sprintf("-----------New run started - %v ------------ filename=%v minDate=%v maxDate=%v \n", time.Now(), dataFilename, minDate, maxDate))
	if err != nil {
		return nil, err
	}
	return f, nil
}

type Worker struct {
	nIter int64
	resultChn chan<- *Result
}

// this is a random search implementation
func New(nIter int64, resultChn chan<- *Result) *Worker {
	return &Worker{nIter, resultChn}
}

func (w *Worker) StartBasicSearch(dataSet []models.DataPoint) {
	for {
		result := w.FindBasicParameters(dataSet)
		w.resultChn <- result
	}
}

func (w *Worker) StartPeriodicSearch(a, b, tc, beta float64, dataSet []models.DataPoint) {
	for {
		result := w.FindPeriodicParameters(a, b, tc, beta, dataSet)
		w.resultChn <- result
	}
}

func (w *Worker) StartFullSearch(dataSet []models.DataPoint) {
	for {
		result := w.FindFullParameters(dataSet)
		w.resultChn <- result
	}
}

func (w *Worker) FindBasicParameters(dataSet []models.DataPoint) *Result {
	result := &Result{J: math.MaxFloat64}

	start := time.Now().Unix()
	for result.N < w.nIter {
		tmpParams := regression.InitRandomBasicParameters(dataSet[len(dataSet)-1].Date)
		cost := regression.J(dataSet, tmpParams)
		if cost < result.J {
			result.J = cost
			result.Params = tmpParams
		}
		result.N++
	}
	result.ExeTime = time.Now().Unix() - start

	return result
}

func (w *Worker) FindPeriodicParameters(a, b, tc, beta float64, dataSet []models.DataPoint) *Result {
	result := &Result{J: math.MaxFloat64}

	start := time.Now().Unix()
	for result.N < w.nIter {
		tmpParams := regression.InitRandomPeriodicParameters(a, b, tc, beta)
		cost := regression.J(dataSet, tmpParams)
		if cost < result.J {
			result.J = cost
			result.Params = tmpParams
		}
		result.N++
	}
	result.ExeTime = time.Now().Unix() - start

	return result
}

func (w *Worker) FindFullParameters(dataSet []models.DataPoint) *Result {
	result := &Result{J: math.MaxFloat64}

	start := time.Now().Unix()
	for result.N < w.nIter {
		tmpParams := regression.InitRandomFullParameters(dataSet[len(dataSet)-1].Date)
		cost := regression.J(dataSet, tmpParams)
		if cost < result.J {
			result.J = cost
			result.Params = tmpParams
		}
		result.N++
	}
	result.ExeTime = time.Now().Unix() - start

	return result
}
