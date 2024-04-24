package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	var url string
	var requests int
	var concurrency int
	flag.StringVar(&url, "url", "", "URL do serviço")
	flag.IntVar(&requests, "requests", 100, "Número total de requests")
	flag.IntVar(&concurrency, "concurrency", 10, "Número de chamadas simultâneas")
	flag.Parse()

	if url == "" {
		fmt.Println("A URL é obrigatória")
		return
	}

	totalTime, successCount, statusDistribution, requestFail, errorDistribution := executeLoadTest(url, requests, concurrency)

	fmt.Printf("Tempo total gasto: %v\n", totalTime)
	fmt.Printf("Requests realizados: %d\n", requests)
	fmt.Printf("Requests com status 200: %d\n", successCount)
	fmt.Println("Distribuição de outros códigos de status HTTP:")
	for status, count := range statusDistribution {
		fmt.Printf(">>\t%d: %d vezes\n", status, count)
	}
	fmt.Printf("Requests não concluidos: %d\n", requestFail)
	fmt.Println("Erros que impediram conclusão:")
	for errMessage, count := range errorDistribution {
		fmt.Printf(">>\t%s\n", errMessage)
		fmt.Printf("  \t%d vezes\n", count)
	}

}

func executeLoadTest(targetURL string, totalRequests, concurrencyLevel int) (time.Duration, int, map[int]int, int, map[string]int) {
	var wg sync.WaitGroup
	var successCount int
	statusDistribution := make(map[int]int)
	errorDistribution := make(map[string]int)
	var requestFail = 0
	start := time.Now()

	stoper := make(chan struct{}, concurrencyLevel)

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		stoper <- struct{}{}
		go func() {
			defer wg.Done()
			response, err := http.Get(targetURL)
			if err == nil {
				if response.StatusCode == 200 {
					successCount++
				} else {
					statusDistribution[response.StatusCode]++
				}
				response.Body.Close()
			} else {
				requestFail++
				errorDistribution[err.Error()]++
			}
			<-stoper
		}()
	}

	wg.Wait()
	totalTime := time.Since(start)
	return totalTime, successCount, statusDistribution, requestFail, errorDistribution
}
