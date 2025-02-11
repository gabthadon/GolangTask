package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
)

func main() {
	keywords := []string{"INFO", "ERROR", "DEBUG"}
	ProcessLogFile("log.txt", keywords)
}

func ProcessLogFile(filePath string, keywords []string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	keywordCounts := make(map[string]int)
	var mu sync.Mutex
	var wg sync.WaitGroup

	chunkSize := 1024 * 4 // 4KB chunks
	buffer := make([]byte, chunkSize)

	reader := bufio.NewReader(file)
	for {
		n, err := reader.Read(buffer)
		if n == 0 {
			break
		}
		if err != nil {
			fmt.Println("Error reading file:", err)
			break
		}

		wg.Add(1)
		go func(data []byte) {
			defer wg.Done()
			counts := CountKeywords(string(data), keywords)
			mu.Lock()
			for k, v := range counts {
				keywordCounts[k] += v
			}
			mu.Unlock()
		}(buffer[:n])
	}

	wg.Wait()

	sortedCounts := make([]struct {
		Keyword string
		Count   int
	}, 0, len(keywordCounts))

	for k, v := range keywordCounts {
		sortedCounts = append(sortedCounts, struct {
			Keyword string
			Count   int
		}{k, v})
	}

	sort.Slice(sortedCounts, func(i, j int) bool {
		return sortedCounts[i].Count > sortedCounts[j].Count
	})

	for _, entry := range sortedCounts {
		fmt.Printf("%s: %d\n", entry.Keyword, entry.Count)
	}
}

func CountKeywords(data string, keywords []string) map[string]int {
	counts := make(map[string]int)
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		for _, keyword := range keywords {
			if strings.Contains(strings.ToUpper(line), strings.ToUpper(keyword)) {
				counts[keyword]++
			}
		}
	}
	return counts
}