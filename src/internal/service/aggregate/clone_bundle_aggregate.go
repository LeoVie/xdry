package aggregate

import "x-dry-go/src/internal/clone_detect"

type CloneInstance struct {
	Path     string
	Language string
	Index    int
}

type AggregatedClone struct {
	Content   string
	Language  string
	Instances []CloneInstance
}

type CloneBundle struct {
	CloneType        int
	AggregatedClones []AggregatedClone
}

func AggregateCloneBundles(clones map[int][]clone_detect.Clone) []CloneBundle {
	cloneBundles := []CloneBundle{}

	for cloneType, clones := range clones {
		aggregatedClones := []AggregatedClone{}

		for _, clone := range clones {
			for _, match := range clone.Matches {

				addedToExistingAggregatedClone := false

				for aggregatedCloneKey, aggregatedClone := range aggregatedClones {
					if aggregatedClone.Content == match.Content {
						addedToExistingAggregatedClone = true

						cloneInstance := CloneInstance{
							Path:     clone.A,
							Language: clone.Language,
							Index:    match.IndexA,
						}
						if !containsCloneInstance(cloneInstance, aggregatedClones[aggregatedCloneKey].Instances) {
							aggregatedClones[aggregatedCloneKey].Instances = append(aggregatedClones[aggregatedCloneKey].Instances, cloneInstance)
						}

						cloneInstance = CloneInstance{
							Path:     clone.B,
							Language: clone.Language,
							Index:    match.IndexB,
						}
						if !containsCloneInstance(cloneInstance, aggregatedClones[aggregatedCloneKey].Instances) {
							aggregatedClones[aggregatedCloneKey].Instances = append(aggregatedClones[aggregatedCloneKey].Instances, cloneInstance)
						}
					}
				}

				if !addedToExistingAggregatedClone {
					aggregatedClones = append(aggregatedClones, AggregatedClone{
						Content:  match.Content,
						Language: clone.Language,
						Instances: []CloneInstance{
							{
								Path:     clone.A,
								Language: clone.Language,
								Index:    match.IndexA,
							},
							{
								Path:     clone.B,
								Language: clone.Language,
								Index:    match.IndexB,
							},
						},
					})
				}
			}
		}

		cloneBundles = append(cloneBundles, CloneBundle{
			CloneType:        cloneType,
			AggregatedClones: aggregatedClones,
		})
	}

	return cloneBundles
}

func containsCloneInstance(cloneInstance CloneInstance, list []CloneInstance) bool {
	for _, x := range list {
		if x == cloneInstance {
			return true
		}
	}
	return false
}
