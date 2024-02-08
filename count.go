package main

func count(resources map[string]int, processes []process, goal goal) int {
	finishedAt := 0
	curr := make([]*process, 0, len(processes))
	for i := range processes {
		processes[i].count = processes[i].minCount.numerator
		if processes[i].initial {
			processes[i].start = 0
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

	// k is a stopper to avoid infinite loops; adjust as needed
	for k := 0; k < 100; k++ {
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
					resources[ingredient.name] -= ingredient.quantity * curr[i].count
				}
				for _, product := range curr[i].products {
					resources[product.name] += product.quantity * curr[i].count
				}
				if curr[i].successor != nil {
					end := curr[i].start + curr[i].time
					if end > curr[i].successor.start {
						curr[i].successor.start = end
					}
					if !added[curr[i].successor.name] {
						added[curr[i].successor.name] = true
						next = append(next, curr[i].successor)
					}
				}
				if curr[i].final {
					finishedAt = curr[i].start + curr[i].time
				}

			}
			curr = next
		}
	}
	return finishedAt
}
