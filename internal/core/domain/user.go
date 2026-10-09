package domain

import (
	"fmt"
	core_error "practice/internal/core/errors"
	core_transport_grpc "practice/internal/core/transport/grpc"
	v1 "practice/pkg/pb/v1"
	"regexp"
)

type User struct {
	ID      int
	Version int

	FullName    string
	PhoneNumber *string
}

func NewUser(id, version int, fullName string, phoneNumber *string) User {
	return User{
		ID:          id,
		Version:     version,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

func NewUserUninitialized(fullName string, phoneNumber *string) User {
	return NewUser(UninitializedID, UninitializedVersion, fullName, phoneNumber)
}

func (u *User) ToProtoUser() *v1.User {
	return &v1.User{
		Id:      int64(u.ID),
		Version: int64(u.Version),
		Info: &v1.UserInfo{
			FullName:    u.FullName,
			PhoneNumber: core_transport_grpc.StringToWrapper(u.PhoneNumber),
		},
	}
}

func (u *User) Validate() error {
	fullNameLen := len([]rune(u.FullName))
	if fullNameLen < 3 || fullNameLen > 100 {
		return fmt.Errorf("invalid fullNameLength %d: %w", fullNameLen, core_error.ErrInvalidArgument)
	}

	if u.PhoneNumber != nil {
		phoneNumberLen := len([]rune(*u.PhoneNumber))
		if phoneNumberLen < 10 || phoneNumberLen > 15 {
			return fmt.Errorf("invalid phoneNumberLen %d: %w", phoneNumberLen, core_error.ErrInvalidArgument)
		}

		re := regexp.MustCompile(`^\+[0-9]+$`)
		if !re.MatchString(*u.PhoneNumber) {
			return fmt.Errorf("invalid phoneNumber format %w", core_error.ErrInvalidArgument)
		}
	}
	return nil
}

type UserPatch struct {
	FullName    Nullable[string]
	PhoneNumber Nullable[string]
}

func NewUserPatch(fullName Nullable[string], phoneNumber Nullable[string]) UserPatch {
	return UserPatch{
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

func NewUserPatchFromProto(protoPatch *v1.PatchUserRequest) UserPatch {
	return UserPatch{
		FullName:    NewNullableFromWrapper(protoPatch.GetPatch().GetFullName()),
		PhoneNumber: NewNullableFromWrapper(protoPatch.GetPatch().GetPhoneNumber()),
	}
}

func (p *UserPatch) Validate() error {
	if p.FullName.Set && p.FullName.Value == nil {
		return fmt.Errorf("FullName can be patched to null: %w", core_error.ErrInvalidArgument)
	}

	return nil
}

func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate user patch: %w", err)
	}

	tmp := *u
	if patch.FullName.Set {
		tmp.FullName = *patch.FullName.Value
	}
	if patch.PhoneNumber.Set {
		tmp.PhoneNumber = patch.PhoneNumber.Value
	}
	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched user: %w", err)
	}
	*u = tmp

	return nil
}
