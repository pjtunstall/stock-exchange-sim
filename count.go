package main

func count(resources []resource, processes []process, goal goal) {
	for i := range processes {
		if processes[i].initial {
			processes[i].start = 0
			for j := range processes[i].ingredients {
				for k := range resources {
					if processes[i].ingredients[j].name == resources[k].name && resources[k].quantity > 0 && resources[k].quantity >= processes[i].ingredients[j].quantity {
						resources[k].quantity -= processes[i].ingredients[j].quantity
						processes[i].count = processes[i].minCount.numerator
					}
				}
			}
		}
	}
}

// // Work in progress: count number of times each process must be run.
// func count(resources []resource, processes []process, goal goal) {
// 	curr := make([]*process, 0, 64)

// 	// This map will be used to keep track of the remaining resources
// 	// at each step, starting with initial resources, and at each step
// 	// containing the products of the previous step.
// 	rMap := make(map[string]int)
// 	for i := range resources {
// 		rMap[resources[i].name] = resources[i].quantity
// 	}

// 	// Seed the curr slice with the initial processes.
// 	for i := range processes {
// 		if processes[i].initial {
// 			curr = append(curr, &processes[i])
// 			processes[i].count = 1
// 			for _, ingredient := range processes[i].ingredients {
// 				rMap[ingredient.name] -= ingredient.quantity
// 				if rMap[ingredient.name] == 0 {
// 					delete(rMap, ingredient.name)
// 				}
// 				if rMap[ingredient.name] < 0 {
// 					panic("insufficient resources")
// 				}
// 			}
// 		}
// 	}

// 	// While there are processes to be run, decrement the resources
// 	// and add the products to the next slice of processes to be run.
// 	for len(curr) > 0 {
// 		var next []*process
// 		for k := range curr {
// 			s := strength(curr[k])

// 		}
// 		// for len(rMap) > 0 {
// 		// 	for k := range curr {
// 		// 		for j := range curr[k].ingredients {
// 		// 			rMap[curr[k].ingredients[j].name] -= curr[k].ingredients[j].quantity
// 		// 			if rMap[curr[k].ingredients[j].name] <= 0 {
// 		// 				delete(rMap, curr[k].ingredients[j].name)
// 		// 			}
// 		// 		}
// 		// 		d := curr[k].minCount.Denominator
// 		// 		if d > 1 {
// 		// 			curr[k].count *= d
// 		// 		}
// 		// 		next = append(next, curr[k].successor)
// 		// 	}
// 		// }
// 		curr = next
// 	}
// }

// func min(a []int) (int, error) {
// 	if len(a) == 0 {
// 		return 0, fmt.Errorf("empty slice")
// 	}
// 	m := a[0]
// 	for i := 1; i < len(a); i++ {
// 		if a[i] < m {
// 			m = a[i]
// 		}
// 	}
// 	return m, nil
// }
