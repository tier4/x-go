package ref

func Ref[T any](s T) *T {
	return &s
}

func Deref[T any](p *T) T {
	if p == nil {
		var ret T
		return ret
	}
	return *p
}
