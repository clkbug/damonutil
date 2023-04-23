package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	damonuntil "github.com/clkbug/damonutil"
)

type cmdOptions struct {
	input  string
	output string
}

var cmdopt cmdOptions

func init() {
	flag.StringVar(&cmdopt.input, "input", "damon.data", "input file path (default: damon.data)")
	flag.StringVar(&cmdopt.output, "output", "", "output file path (default: stdout)")
}

func main() {
	flag.Parse()

	err := run(cmdopt)
	if err != nil {
		fmt.Printf("error: %s", err.Error())
		os.Exit(1)
	}
}

func run(cmdopt cmdOptions) error {
	damon, err := damonuntil.ParseDamonFile(cmdopt.input)
	if err != nil {
		return err
	}

	if cmdopt.output == "" {
		printDamonResult(damon)
		return nil
	}

	if strings.HasSuffix(cmdopt.output, ".json") {
		fp, err := os.Create(cmdopt.output)
		if err != nil {
			return err
		}
		defer fp.Close()
		buf := bufio.NewWriter(fp)
		if err := json.NewEncoder(buf).Encode(damon); err != nil {
			return err
		}
		if err := buf.Flush(); err != nil {
			return err
		}
	}

	return nil
}

func printDamonResult(damon *damonuntil.Result) {
	baseTime := uint64(0)
	for i, record := range damon.Records {
		for j, snapshot := range record.Snapshots {
			if i == 0 && j == 0 {
				baseTime = snapshot.StartTime
				fmt.Printf("base_time_absolute: %d\n\n", baseTime)
			}
			fmt.Printf("monitoring_start:    %16d\n", snapshot.StartTime-baseTime)
			fmt.Printf("monitoring_end:      %16d\n", snapshot.EndTime-baseTime)
			fmt.Printf("monitoring_duration: %16d\n", snapshot.EndTime-snapshot.StartTime)
			fmt.Printf("target_id: %d\n", snapshot.TargetId)
			fmt.Printf("nr_regions: %d\n", len(snapshot.Regions))
			fmt.Println("# start_addr     end_addr        length  nr_accesses   age")
			for _, region := range snapshot.Regions {
				fmt.Printf("%012x-%012x (%12d) %11d %5d\n", region.StartAddr, region.EndAddr, region.EndAddr-region.StartAddr, region.NumberOfAccesses, region.Age)
			}
			fmt.Println()
		}
	}
}
