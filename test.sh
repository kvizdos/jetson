#!/bin/bash

RUNS=10
OUTFILE="results.tmp"
> "$OUTFILE"

echo "Building.."
go build -gcflags="all=-B -l" -o benchmark ./cmd/jetson/main.go

echo "Running $RUNS iterations..."
for ((i = 1; i <= RUNS; i++)); do
    vmtouch -e ./demo/arxiv-metadata-oai-snapshot.json > /dev/null
    GOGC=off nice -n -20 ./benchmark -file=./demo/arxiv-metadata-oai-snapshot.json -key=categories -value=hep >> "$OUTFILE"
    echo "Iteration $i completed"
done

rm benchmark

extract_times() {
    local field="$1"
    awk -v key="$field=" '
    {
        for (i = 1; i <= NF; i++) {
            if ($i ~ "^"key) {
                gsub("s", "", $i)
                split($i, a, "=")
                print a[2]
            }
        }
    }' "$OUTFILE"
}

calculate_summary() {
    local label=$1
    shift
    local arr=("$@")

    local sorted=($(printf '%s\n' "${arr[@]}" | sort -n))
    local count=${#sorted[@]}
    local min=${sorted[0]}
    local max=${sorted[$((count-1))]}
    local sum=0

    for val in "${sorted[@]}"; do
        sum=$(echo "$sum + $val" | bc)
    done

    local avg=$(echo "scale=6; $sum / $count" | bc)

    local med
    if (( count % 2 == 0 )); then
        m1=${sorted[$((count / 2 - 1))]}
        m2=${sorted[$((count / 2))]}
        med=$(echo "scale=6; ($m1 + $m2) / 2" | bc)
    else
        med=${sorted[$((count / 2))]}
    fi

    echo "$label : $min / $max / $avg / $med"
}

read_times=($(extract_times "read"))
process_times=($(extract_times "process"))
total_times=($(extract_times "total_time"))
throughputs=($(extract_times "throughput"))

echo "----- Summary -----"
echo "Iterations: $RUNS"
echo "Format: Min / Max / Avg / Median"
calculate_summary "Read     (s)" "${read_times[@]}"
calculate_summary "Process  (s)" "${process_times[@]}"
calculate_summary "Total    (s)" "${total_times[@]}"
calculate_summary "Throughput (MB/s)" "${throughputs[@]}"
