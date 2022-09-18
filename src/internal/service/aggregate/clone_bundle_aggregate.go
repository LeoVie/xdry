package aggregate

import "x-dry-go/src/internal/clone_detect"

type cloneInstance struct {
	Path     string
	Language string
	Index    int
}

type aggregatedClone struct {
	Content   string
	Language  string
	Instances []cloneInstance
}

func (aggrClone *aggregatedClone) hasInstance(instance cloneInstance) bool {
	for _, x := range aggrClone.Instances {
		if x == instance {
			return true
		}
	}
	return false
}

func (aggrClone *aggregatedClone) addInstanceIfNotAlreadyHaving(instance cloneInstance) {
	if !aggrClone.hasInstance(instance) {
		aggrClone.Instances = append(aggrClone.Instances, instance)
	}
}

type CloneBundle struct {
	CloneType        int
	AggregatedClones []aggregatedClone
}

func newCloneBundle(cloneType int, aggregatedClones []aggregatedClone) *CloneBundle {
	if aggregatedClones == nil {
		aggregatedClones = []aggregatedClone{}
	}
	return &CloneBundle{
		CloneType:        cloneType,
		AggregatedClones: aggregatedClones,
	}
}

func AggregateCloneBundles(clones map[int][]clone_detect.Clone) []CloneBundle {
	var cloneBundles []CloneBundle

	for cloneType, clones := range clones {
		var aggregatedClones []aggregatedClone

		for _, clone := range clones {
			for _, match := range clone.Matches {
				instanceA := cloneInstance{
					Path:     clone.A,
					Language: clone.Language,
					Index:    match.IndexA,
				}
				instanceB := cloneInstance{
					Path:     clone.B,
					Language: clone.Language,
					Index:    match.IndexB,
				}

				addedToExistingAggregatedClone := false

				for aggregatedCloneKey, aggrClone := range aggregatedClones {
					if aggrClone.Content == match.Content {
						addedToExistingAggregatedClone = true
						aggregatedClones[aggregatedCloneKey].addInstanceIfNotAlreadyHaving(instanceA)
						aggregatedClones[aggregatedCloneKey].addInstanceIfNotAlreadyHaving(instanceB)
						continue
					}
				}

				if !addedToExistingAggregatedClone {
					aggregatedClones = append(aggregatedClones, aggregatedClone{
						Content:  match.Content,
						Language: clone.Language,
						Instances: []cloneInstance{
							instanceA,
							instanceB,
						},
					})
				}
			}
		}

		cloneBundles = append(cloneBundles, *newCloneBundle(cloneType, aggregatedClones))
	}

	return cloneBundles
}
