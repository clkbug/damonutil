package main

import (
	"fmt"
	"os"

	damonuntil "github.com/clkbug/damonutil"
)

func main() {
	err := run()
	if err != nil {
		fmt.Printf("error: %s", err.Error())
		os.Exit(1)
	}
}

func run() error {
	input := "damon.data"

	if len(os.Args) > 1 {
		input = os.Args[1]
	}

	damon, err := damonuntil.ParseDamonFile(input)
	if err != nil {
		return err
	}

	printDamonResult(damon)

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
