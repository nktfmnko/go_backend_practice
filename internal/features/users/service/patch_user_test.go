package users_service

import (
	"context"
	"errors"
	"practice/internal/core/domain"
	"testing"
)

func TestPatchUser(t *testing.T) {
	ctx := context.Background()

	fullName := "Nikita Updated"
	phoneNumber := "+12345678901"
	invalidFullName := "N"

	tests := []struct {
		name         string
		id           int
		patch        domain.UserPatch
		mockBehavior func(m *mockUsersRepository)
		expectedUser domain.User
		expectedErr  bool
	}{
		{
			name: "Success - update full name",
			id:   123,
			patch: domain.NewUserPatch(
				domain.Nullable[string]{
					Value: &fullName,
					Set:   true,
				},
				domain.Nullable[string]{},
			),
			mockBehavior: func(m *mockUsersRepository) {
				m.getUserFunc = func(ctx context.Context, id int) (domain.User, error) {
					return domain.User{
						ID:       123,
						FullName: "Nikita",
					}, nil
				}

				m.patchUserFunc = func(ctx context.Context, id int, user domain.User) (domain.User, error) {
					if user.FullName != fullName {
						t.Errorf("expected full name %s, got %s", fullName, user.FullName)
					}

					return user, nil
				}
			},
			expectedUser: domain.User{
				ID:       123,
				FullName: fullName,
			},
			expectedErr: false,
		},
		{
			name: "Success - update phone number",
			id:   123,
			patch: domain.NewUserPatch(
				domain.Nullable[string]{},
				domain.Nullable[string]{
					Value: &phoneNumber,
					Set:   true,
				},
			),
			mockBehavior: func(m *mockUsersRepository) {
				m.getUserFunc = func(ctx context.Context, id int) (domain.User, error) {
					return domain.User{
						ID:       123,
						FullName: "Nikita",
					}, nil
				}

				m.patchUserFunc = func(ctx context.Context, id int, user domain.User) (domain.User, error) {
					if user.PhoneNumber == nil {
						t.Error("expected phone number to be updated")
					}

					if *user.PhoneNumber != phoneNumber {
						t.Errorf(
							"expected phone number %s, got %s",
							phoneNumber,
							*user.PhoneNumber,
						)
					}

					return user, nil
				}
			},
			expectedUser: domain.User{
				ID:          123,
				FullName:    "Nikita",
				PhoneNumber: &phoneNumber,
			},
			expectedErr: false,
		},
		{
			name: "Error - get user",
			id:   123,
			patch: domain.NewUserPatch(
				domain.Nullable[string]{
					Value: &fullName,
					Set:   true,
				},
				domain.Nullable[string]{},
			),
			mockBehavior: func(m *mockUsersRepository) {
				m.getUserFunc = func(ctx context.Context, id int) (domain.User, error) {
					return domain.User{}, errors.New("db error")
				}
			},
			expectedUser: domain.User{},
			expectedErr:  true,
		},
		{
			name: "Error - apply patch",
			id:   123,
			patch: domain.NewUserPatch(
				domain.Nullable[string]{
					Value: &invalidFullName,
					Set:   true,
				},
				domain.Nullable[string]{},
			),
			mockBehavior: func(m *mockUsersRepository) {
				m.getUserFunc = func(ctx context.Context, id int) (domain.User, error) {
					return domain.User{
						ID:       123,
						FullName: "Nikita",
					}, nil
				}

				m.patchUserFunc = func(ctx context.Context, id int, user domain.User) (domain.User, error) {
					t.Error("PatchUser should not be called when ApplyPatch returns an error")
					return domain.User{}, nil
				}
			},
			expectedUser: domain.User{},
			expectedErr:  true,
		},
		{
			name: "Error - patch user",
			id:   123,
			patch: domain.NewUserPatch(
				domain.Nullable[string]{
					Value: &fullName,
					Set:   true,
				},
				domain.Nullable[string]{},
			),
			mockBehavior: func(m *mockUsersRepository) {
				m.getUserFunc = func(ctx context.Context, id int) (domain.User, error) {
					return domain.User{
						ID:       123,
						FullName: "Nikita",
					}, nil
				}

				m.patchUserFunc = func(ctx context.Context, id int, user domain.User) (domain.User, error) {
					return domain.User{}, errors.New("db error")
				}
			},
			expectedUser: domain.User{},
			expectedErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockUsersRepository{}

			if tt.mockBehavior != nil {
				tt.mockBehavior(mockRepo)
			}

			service := NewUsersService(mockRepo)

			user, err := service.PatchUser(ctx, tt.id, tt.patch)

			if (err != nil) != tt.expectedErr {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
				return
			}

			if user.ID != tt.expectedUser.ID {
				t.Errorf("expected user ID %d, got %d", tt.expectedUser.ID, user.ID)
			}

			if user.FullName != tt.expectedUser.FullName {
				t.Errorf(
					"expected full name %s, got %s",
					tt.expectedUser.FullName,
					user.FullName,
				)
			}
		})
	}
}
