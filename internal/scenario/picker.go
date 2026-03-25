package scenario
import "math/rand"

func PickEndpoint(endpoints []Endpoint) *Endpoint {
	if len(endpoints) == 0 {
		return nil
	}
	if len(endpoints) == 1 {
		return &endpoints[0]
	}
	total := 0
	for _, ep := range(endpoints) {
		total += ep.Weight
	}

	r := rand.Intn(total)
	cumlative := 0
	for i, ep := range endpoints {
		cumlative += ep.Weight
		if r < cumlative {
			return &endpoints[i]
		}
	}

	return &endpoints[len(endpoints)-1]
}

