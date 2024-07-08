package main

func generatePackagesGraph(gofiles []*goFile) {
	depsMap := make(map[string][]*goPackage)
	for _, goFile := range gofiles {
		packageName := goFile.goPackage.name
		// Use a map to track existing dependencies for each package
		existingDeps := make(map[string]bool)
		for _, depPackage := range depsMap[packageName] {
			existingDeps[depPackage.name] = true
		}
		for _, depPackage := range goFile.deps {
			// Add dependency if it doesn't exist
			if !existingDeps[depPackage.name] {
				depsMap[packageName] = append(depsMap[packageName], depPackage)
				existingDeps[depPackage.name] = true
			}
		}
	}
	for name, deps := range depsMap {
		println("=================")
		println(name)
		for _, d := range deps {
			println(d)
		}
		println("=================")

	}
}
