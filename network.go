package main

func buildNetwork(resources []resource, processes []process, g goal) {
	var curr []process
	for _, process := range processes {
		for _, product := range process.products {
			if product.name == g.product {
				process.final = true
				curr = append(curr, process)
			}
		}
	}

	// Work back through the processes to build precedence relationships.
	for len(curr) > 0 {
		var next []process
		for _, c := range curr {
			for _, ingredient := range c.ingredients {
				for _, resource := range resources {
					if ingredient.name == resource.name {
						c.initial = true
					}
				}
				for _, process := range processes {
					for _, product := range process.products {
						if ingredient.name == product.name {
							process.successor = &c
							next = append(next, process)
						}
					}
				}
			}
		}
		curr = next
	}
}
