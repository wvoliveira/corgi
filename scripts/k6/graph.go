//go:build never
// +build never

// CLI tool to get values from stdout and print a simple ascii graphic
// Its really good to test API local with k6 + asciigraph

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/guptarohit/asciigraph"
)

func main() {
	data := [][]float64{{0}, {0}}

	nextFlushTime := time.Now()
	flushInterval := time.Duration(float64(time.Second) * 1)

	s := bufio.NewScanner(os.Stdin)
	s.Split(bufio.ScanWords)

	valueA := float64(0)
	valueB := float64(0)

	for s.Scan() {
		word := s.Text()
		n, err := strconv.ParseFloat(word, 64)

		if err != nil {
			log.Printf("ignore %q: cannot parse value", word)
			continue
		}

		// If http status is 200 append value in first line
		// and keep the last value for second line
		// valueA == http 200
		// valueB == http !200

		if n == 200 {
			valueA++
		} else {
			valueB++
		}

		if currentTime := time.Now(); currentTime.After(nextFlushTime) || currentTime.Equal(nextFlushTime) {
			data[0] = append(data[0], valueA)
			data[1] = append(data[1], valueB)

			fmt.Println("200: ", valueA)
			fmt.Println("not 200: ", valueB)

			// Reset values to get count per second
			valueA = float64(0)
			valueB = float64(0)

			graph := asciigraph.PlotMany(
				data,
				asciigraph.Precision(0),
				asciigraph.SeriesColors(asciigraph.Green, asciigraph.Red),
				asciigraph.Height(int(10)),
				asciigraph.Width(int(0)),
				asciigraph.Offset(int(3)),
				asciigraph.Precision(2),
				asciigraph.Caption("Green: 200 / Red: !200"),
				asciigraph.CaptionColor(asciigraph.White),
				asciigraph.AxisColor(asciigraph.White),
				asciigraph.LabelColor(asciigraph.White),
			)

			asciigraph.Clear()
			fmt.Println(graph)
			nextFlushTime = time.Now().Add(flushInterval)
		}
	}

}
