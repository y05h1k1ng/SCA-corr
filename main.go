package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"
	"sync"

	"github.com/montanaflynn/stats"
	"github.com/schollz/progressbar/v3"
)

func ReadCsv(filePath string) ([][]float64, error) {
	// TODO: speed-up
	// TODO: support progressbar
	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(f)
	stringValues, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// convert 2D-string to 2D-float64 array
	var table [][]float64
	for _, line := range stringValues {
		var nums []float64
		for _, v := range line {
			if n, err := strconv.ParseFloat(v, 64); err == nil {
				nums = append(nums, n)
			} else {
				return nil, err
			}
		}
		table = append(table, nums)
	}

	return table, nil
}

func WriteCsv(filePath string, data [][]float64) error {
	// convert 2D-float64 to 2D-string array
	var table [][]string
	for _, line := range data {
		var nums []string
		for _, v := range line {
			s := strconv.FormatFloat(v, 'f', -1, 64)
			nums = append(nums, s)
		}
		table = append(table, nums)
	}

	f, err := os.Create(filePath)
	defer f.Close()
	if err != nil {
		return err
	}

	writer := csv.NewWriter(f)
	writer.WriteAll(table)

	return nil
}

func corr(x, y []float64) (float64, error) {
	cor, err := stats.Correlation(x, y)
	return math.Abs(cor), err
}

func transpose(M [][]float64) [][]float64 {
	xl := len(M[0])
	yl := len(M)

	res := make([][]float64, xl)
	for i := range res {
		res[i] = make([]float64, yl)
	}

	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			res[i][j] = M[j][i]
		}
	}
	return res
}

// TODO: error handling
func main() {
	cpun := runtime.NumCPU()
	fmt.Println("[*] Number of your CPU:", cpun)
	// runtime.GOMAXPROCS(cpun)

	// TODO: more useful, support --verbose flag(logging)
	waveFile := flag.String("wave", "", "a csv file of wave")
	intermFile := flag.String("interm", "", "a csv file of intermidiate values")
	outputFile := flag.String("output", "", "output csv file")
	flag.Parse()

	waves, err := ReadCsv(*waveFile)
	if err != nil {
		fmt.Println(err)
	}
	interm, err := ReadCsv(*intermFile)
	if err != nil {
		fmt.Println(err)
	}

	tl := len(waves[0])
	corrTable := make([][]float64, 256)
	for i := range corrTable {
		corrTable[i] = make([]float64, tl)
	}

	wavesT := transpose(waves)
	intermT := transpose(interm)

	bar := progressbar.Default(256)

	var lock sync.Mutex
	var wg sync.WaitGroup

	for k_idx := 0; k_idx < 256; k_idx++ {
		bar.Add(1) // TODO: support goroutine
		for t := 0; t < tl; t++ {
			wg.Add(1)
			go func(t, k_idx int) {
				defer wg.Done()
				lock.Lock()
				defer lock.Unlock()
				corrTable[k_idx][t], err = corr(wavesT[t], intermT[k_idx])
				if err != nil {
					fmt.Println(err)
					os.Exit(0)
				}
			}(t, k_idx)
		}
	}

	wg.Wait()
	fmt.Println("[+] done.")

	if err := WriteCsv(*outputFile, corrTable); err != nil {
		fmt.Println(err)
	}
}
