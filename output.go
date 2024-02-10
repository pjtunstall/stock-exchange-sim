package main

import (
	"fmt"
	"strings"
)

func buildOutput(resources map[string]int, processes []process, end int, finite bool, starts map[int][]string, t string) string {
	var builder strings.Builder
	builder.WriteString("Processes scheduled:\n")
	var someScheduled bool
	if finite {
		for _, process := range processes {
			for i := 0; i <= end; i++ {
				if process.start == i {
					for j := 0; j < process.iterations; j++ {
						someScheduled = true
						builder.WriteString(fmt.Sprintf(" %d:%s\n", process.start, process.name))
					}
				}
			}
		}
	} else {
		builder.WriteString(t)
		if len(t) > 0 {
			someScheduled = true
		}
		// for i := 0; i <= 11; i++ {
		// 	if _, ok := starts[i]; ok {
		// 		for _, process := range starts[i] {
		// 			someScheduled = true
		// 			builder.WriteString(fmt.Sprintf(" %d:%s\n", i, process))
		// 		}
		// 	}
		// }
		//
		// for start := range starts {
		// 	for _, process := range starts[start] {
		// 		someScheduled = true
		// 		builder.WriteString(fmt.Sprintf(" %d:%s\n", start, process))
		// 	}
		// }
	}
	if !someScheduled {
		builder.WriteString(" none\n")
	}
	builder.WriteString(fmt.Sprintf("No more process doable at cycle %d.\n", end+1))
	builder.WriteString("Stock left:\n")
	for resource, quantity := range resources {
		builder.WriteString(fmt.Sprintf(" %s:%d\n", resource, quantity))
	}
	s := builder.String()
	return s
}

// func writeOutput(s string) {
// 	file, err := os.Create(os.Args[1] + ".log")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()

// 	_, err = io.WriteString(file, s)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
