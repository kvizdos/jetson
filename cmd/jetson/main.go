package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/kvizdos/jetson/jetson"
)

func main() {
	file := flag.String("file", "", "Path to the NDJSON file")
	key := flag.String("key", "", "Key to search for")
	value := flag.String("value", "", "Substring value to match")
	cpuProfile := flag.String("cpu", "", "Write CPU profile to file")
	memProfile := flag.String("mem", "", "Write memory profile to file")

	flag.Parse()

	if *file == "" || *key == "" || *value == "" {
		fmt.Println("Usage: jetson -file=your.ndjson -key=categories -value=hep [-cpu=cpu.pprof] [-mem=heap.pprof]")
		os.Exit(1)
	}

	// Start CPU profiling if requested
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	scanner := jetson.NewJetsonScanner(*file, *key, *value)
	result, err := scanner.Scan()
	if err != nil {
		log.Fatal(err)
	}

	// Write heap profile if requested
	if *memProfile != "" {
		f, err := os.Create(*memProfile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		pprof.WriteHeapProfile(f)
	}

	result.Print()
}
