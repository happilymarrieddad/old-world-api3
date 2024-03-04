package utils

func Ref[T any](V T) *T {
	return &V
}
