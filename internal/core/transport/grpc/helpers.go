package core_transport_grpc

import (
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func StringToWrapper(str *string) *wrapperspb.StringValue {
	if str == nil {
		return nil
	}
	return wrapperspb.String(*str)
}
