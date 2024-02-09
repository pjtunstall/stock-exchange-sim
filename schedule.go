package main

import "fmt"

func schedule(resources map[string]int, processes []process, goal goal, finite bool, c <-chan struct{}) (int, map[int][]string) {
	finishedAt := 0
	curr := make([]*process, 0, len(processes))
	starts := make(map[int][]string)

	for i := range processes {
		processes[i].count = processes[i].minCount.numerator
		if processes[i].initial {
			processes[i].start = 0
			if !finite {
				processes[i].startInfinite = append(processes[i].startInfinite, 0)
			}
			processes[i].count = processes[i].minCount.numerator
			processes[i].doable = true
		}
	}

	for i := range processes {
		arr := make([]int, 0, len(processes))
		for j := range processes {
			arr = append(arr, processes[j].minCount.denominator)
		}
		processes[i].maxCount = lcm(arr)
		processes[i].count = processes[i].minCount.Times(rational{processes[i].maxCount, 1}).numerator
	}

	// Infinite loop limited by timer. See checkArgs() in stock.go
	// for default timer value.
	for {
		for i := range processes {
			if processes[i].initial && processes[i].doable {
				isSufficient := true
				for _, ingredient := range processes[i].ingredients {
					less := resources[ingredient.name] < ingredient.quantity*processes[i].count
					if less || resources[ingredient.name] == 0 {
						isSufficient = false
					}
				}
				if !isSufficient {
					processes[i].doable = false
					continue
				}
				curr = append(curr, &processes[i])
			}
		}

		if len(curr) == 0 {
			break
		}

		for len(curr) > 0 {
			fmt.Println("curr: ", curr[0].name)
			select {
			case <-c:
				return finishedAt, starts
			default:
			}
			next := make([]*process, 0, len(processes))
			added := make(map[string]bool)

			for i := range curr {
				isSufficient := true
				for _, ingredient := range curr[i].ingredients {
					less := resources[ingredient.name] < ingredient.quantity*curr[i].count
					if less || resources[ingredient.name] == 0 {
						isSufficient = false
					}
				}
				if !isSufficient {
					continue
				}
				curr[i].iterations += curr[i].count
				for _, ingredient := range curr[i].ingredients {
					if ingredient.name != "you" {
						resources[ingredient.name] -= ingredient.quantity * curr[i].count
					}
				}
				for _, product := range curr[i].products {
					if product.name != "you" {
						resources[product.name] += product.quantity * curr[i].count
					}
				}
				if curr[i].successor != nil {
					end := curr[i].start + curr[i].time
					firstTime := len(curr[i].startInfinite) == 0
					if (finite || firstTime) && end > curr[i].successor.start {
						curr[i].successor.start = end
					}
					// Sheer chaos! Not working yet.
					if !finite {
						if len(curr[i].successor.startInfinite) == 0 {
							end = curr[i].start + curr[i].time
							curr[i].startInfinite = append(curr[i].startInfinite, end)
						} else {
							if len(curr[i].startInfinite) > 0 {
								end = curr[i].startInfinite[len(curr[i].startInfinite)-1] + curr[i].time
							} else {
								end = curr[i].start + curr[i].time
							}
							if end > curr[i].successor.startInfinite[len(curr[i].successor.startInfinite)-1] {
								curr[i].successor.startInfinite = append(curr[i].successor.startInfinite, end)
							}
						}
					}
					// End of chaos.
					starts[end] = append(starts[end], curr[i].successor.name)
					if !added[curr[i].successor.name] {
						added[curr[i].successor.name] = true
						next = append(next, curr[i].successor)
					}
				}

				if finite && curr[i].final {
					finishedAt = curr[i].start + curr[i].time
				}
				if !finite {
					finishedAt = maxTime(curr)
				}
			}
			curr = next
		}
	}
	return finishedAt, starts
}

func maxTime(processes []*process) int {
	max := 0
	for i := range processes {
		if processes[i].time > max {
			max = processes[i].time
		}
	}
	return max
}
