package crawler

func Contains[Type comparable](s []Type, e Type) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
