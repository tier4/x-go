package stringx

import "time"

func Ref(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func Deref(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func FormatDateTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
