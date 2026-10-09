package domain

type Nullable[T any] struct {
	Value *T
	Set   bool
}

type ValueWrapper[T any] interface {
	GetValue() T
}

func NewNullableFromWrapper[T any, W ValueWrapper[T]](w W) Nullable[T] {
	if any(w) == nil || isNil(w) {
		return Nullable[T]{
			Value: nil,
			Set:   false,
		}
	}

	val := w.GetValue()
	return Nullable[T]{
		Value: &val,
		Set:   true,
	}
}

func isNil(i any) bool {
	if i == nil {
		return false
	}
	type nillable interface{ IsNil() bool }
	if n, ok := i.(nillable); ok {
		return n.IsNil()
	}
	return false
}
