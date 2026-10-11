package core_transport_grpc

import (
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

func PtrConvert[From, To Number](val *From) *To {
	if val == nil {
		return nil
	}
	v := To(*val)
	return &v
}

func ConvertProtoToTimePtr(protoTime *timestamppb.Timestamp) *time.Time {
	if protoTime == nil {
		return nil
	}
	t := protoTime.AsTime()
	return &t
}

func ConvertDurationToProto(d *time.Duration) *durationpb.Duration {
	if d == nil {
		return nil
	}
	return durationpb.New(*d)
}
