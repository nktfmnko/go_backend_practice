package core_transport_grpc

func PtrInt64ToInt(val *int64) *int {
	if val == nil {
		return nil
	}
	v := int(*val)
	return &v
}
