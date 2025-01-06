package sentinel_config_overrides

func yieldPtrNodes[T any](this, that map[string]*T, comparer func(p, o T)) {
	for name, thisNode := range this {
		if thisNode == nil {
			continue
		}
		thatNode := that[name]
		if thatNode == nil {
			continue
		}
		comparer(*thisNode, *thatNode)
	}
}

func yieldNodes[T any](this, that map[string]T, comparer func(p, o T)) {
	for name, thisNode := range this {
		if thatNode, ok := that[name]; ok {
			comparer(thisNode, thatNode)
		}
	}
}

func mapChanged[T any](this, that map[string]*T, comparer func(p, o T) bool) bool {
	if len(this) == 0 || len(that) == 0 {
		return true
	}

	for name, thatNode := range that {
		if thatNode == nil {
			return true
		}
		thisNode, ok := this[name]
		if !ok || thisNode == nil {
			return true
		}
		if !comparer(*thisNode, *thatNode) {
			return true
		}
	}
	return false
}

func noStringChange(this, that string) bool {
	return this == that && that != ""
}

func noStringListChange(this, that []string) bool {
	if len(this) != len(that) {
		return false
	}
	if len(this) == 0 {
		return false
	}

	for idx := range this {
		if this[idx] != that[idx] {
			return false
		}
	}

	return true
}
