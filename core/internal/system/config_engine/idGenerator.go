package config_engine

func NextRouteID() int {
	max := 0
	for _, r := range candidate.Routes {
		if r.ID > max {
			max = r.ID
		}
	}
	return max + 1
}

func NextPolicyID() int {
	max := 0
	for _, p := range candidate.Policies {
		if p.ID > max {
			max = p.ID
		}
	}
	return max + 1
}

func NextNATID() int {
	max := 0
	for _, n := range candidate.NATRules {
		if n.ID > max {
			max = n.ID
		}
	}
	return max + 1
}
