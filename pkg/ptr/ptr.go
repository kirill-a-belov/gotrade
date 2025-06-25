package ptr

func PtrString(s string) *string {
	return &s
}

func PtrBool(b bool) *bool {
	return &b
}
