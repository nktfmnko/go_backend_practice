package domain

type Nullable[T any] struct {
	Value *T
	Set   bool
}

func NewNullableFromPtr[T any](ptr *T) Nullable[T] {
	if ptr == nil {
		return Nullable[T]{
			Value: nil,
			Set:   false,
		}
	}
	return Nullable[T]{
		Value: ptr,
		Set:   true,
	}
}
