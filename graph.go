package main

func generatePackagesGraph(gofiles []*goFile) map[string][]*goPackage {
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

	return depsMap
}
