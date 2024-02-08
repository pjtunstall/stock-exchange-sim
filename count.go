package main

func count(resources map[string]int, processes []process, goal goal) int {
	finishedAt := 0
	curr := make([]*process, 0, len(processes))
	for i := range processes {
		processes[i].count = processes[i].minCount.numerator
		if processes[i].initial {
			processes[i].start = 0
			processes[i].count = processes[i].minCount.numerator
			curr = append(curr, &processes[i])
		}
	}

	// Working on this bit. The idea is to calculate the ratios
	// that each process can be performed in, given enough resources.
	// Needs checking and needs logic to respect resource availability.
	for i := range processes {
		arr := make([]int, 0, len(processes))
		for j := range processes {
			arr = append(arr, processes[j].minCount.denominator)
		}
		processes[i].maxCount = lcm(arr)
		processes[i].count = processes[i].minCount.Times(rational{processes[i].maxCount, 1}).numerator
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
	return finishedAt
}
