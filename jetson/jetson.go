package jetson

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

type JetsonScanner struct {
	Filename    string
	SearchKey   []byte
	SearchValue []byte
	Workers     int
	Data        []byte
	readTime    time.Duration
}

type ScanResult struct {
	Matches         int64
	TotalLines      int
	ReadDuration    time.Duration
	ProcessDuration time.Duration
	TotalDuration   time.Duration
	DataSizeBytes   int
}

func (r *ScanResult) Print() {
	mb := float64(r.DataSizeBytes) / 1024.0 / 1024.0
	speed := mb / r.TotalDuration.Seconds()

	fmt.Printf("matches=%d workers=%d read=%.6fs process=%.6fs total_time=%.6fs throughput=%.2f MB/s\n",
		r.Matches,
		runtime.NumCPU(),
		r.ReadDuration.Seconds(),
		r.ProcessDuration.Seconds(),
		r.TotalDuration.Seconds(),
		speed,
	)
}

func (s *JetsonScanner) ReadFile() error {
	start := time.Now()
	data, err := os.ReadFile(s.Filename)
	if err != nil {
		return err
	}
	s.Data = data
	s.readTime = time.Since(start)
	return nil
}

func findNextNewline(data []byte, start int) int {
	i := bytes.IndexByte(data[start:], '\n')
	if i == -1 {
		return len(data)
	}
	return start + i + 1
}

func (s *JetsonScanner) Scan() (*ScanResult, error) {
	if s.Data == nil {
		if err := s.ReadFile(); err != nil {
			return nil, err
		}
	}

	start := time.Now()
	n := len(s.Data)

	if s.Workers < 1 {
		s.Workers = runtime.NumCPU()
	}
	chunkSize := n / s.Workers

	var wg sync.WaitGroup
	matches := make([]int64, s.Workers)

	for i := range s.Workers {
		wg.Add(1)
		var startIdx, endIdx int

		// prevents reprocessing partial line
		if i == 0 {
			startIdx = 0
		} else {
			startIdx = findNextNewline(s.Data, i*chunkSize)
		}

		// prevents truncating a line
		if i == s.Workers-1 {
			endIdx = n
		} else {
			endIdx = findNextNewline(s.Data, (i+1)*chunkSize)
		}

		go func(idx, from, to int) {
			defer wg.Done()
			var count int64
			for i := from; i < to; {
				next := findNextNewline(s.Data, i)
				line := s.Data[i:next]

				if start := bytes.Index(line, s.SearchKey); start != -1 {
					valueStart := start + len(s.SearchKey)
					lineTail := line[valueStart:]
					if end := bytes.IndexByte(lineTail, '"'); end != -1 {
						if bytes.Contains(lineTail[:end], s.SearchValue) {
							count++
						}
					}
				}

				i = next
			}
			matches[idx] = count
		}(i, startIdx, endIdx)
	}

	wg.Wait()

	var totalMatches int64
	for _, m := range matches {
		totalMatches += m
	}

	return &ScanResult{
		Matches:         totalMatches,
		ReadDuration:    s.readTime,
		ProcessDuration: time.Since(start),
		TotalDuration:   s.readTime + time.Since(start),
		DataSizeBytes:   len(s.Data),
	}, nil
}

func NewJetsonScanner(fileName string, searchKey string, searchValue string) *JetsonScanner {
	return &JetsonScanner{
		Filename:    fileName,
		SearchKey:   fmt.Appendf([]byte{}, `"%s":"`, searchKey),
		SearchValue: []byte(searchValue),
	}
}
