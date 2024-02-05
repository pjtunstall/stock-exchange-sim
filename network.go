package main

func buildNetwork(g goal, processes []process) {
	for _, process := range processes {
		for _, product := range process.products {
			if product.name == g.product {
				process.final = true
			}
		}
	}

	// Work back through the processes to build precedence relationships.
}
