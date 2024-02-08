package main

import (
	"fmt"
	"strings"
)

func buildOutput(resources map[string]int, processes []process, end int) string {
	var builder strings.Builder
	builder.WriteString("Main Processes:\n")
	for _, process := range processes {
		for i := 0; i <= end; i++ {
			if process.start == i {
				for j := 0; j < process.iterations; j++ {
					builder.WriteString(fmt.Sprintf(" %d:%s\n", process.start, process.name))
				}
			}
		}
	}
	builder.WriteString(fmt.Sprintf("No more process doable at cycle %d\n", end+1))
	builder.WriteString("Stock:\n")
	for resource, quantity := range resources {
		builder.WriteString(fmt.Sprintf(" %s:%d\n", resource, quantity))
	}
	s := builder.String()
	fmt.Println()
	fmt.Println(s)
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
