package users_service

import (
	"context"
	"practice/internal/core/domain"
	"reflect"
	"testing"
)

func TestCreateUser(t *testing.T) {
	bgCtx := context.Background()
	phone := "+12345678901"

	tests := []struct {
		name         string
		fullName     string
		phoneNumber  *string
		mockBehavior func(m *mockUsersRepository)
		expectedUser domain.User
		expectedErr  bool
	}{
		{
			name:        "Success - user create",
			fullName:    "Nikita",
			phoneNumber: &phone,
			mockBehavior: func(m *mockUsersRepository) {
				m.createUserFunc = func(
					ctx context.Context,
					user domain.User,
				) (domain.User, error) {
					return domain.User{
						ID:          123,
						FullName:    user.FullName,
						PhoneNumber: user.PhoneNumber,
					}, nil
				}
			},
			expectedUser: domain.User{
				ID:          123,
				FullName:    "Nikita",
				PhoneNumber: &phone,
			},
		},
		{
			name:         "Error - validation",
			fullName:     "N",
			phoneNumber:  nil,
			mockBehavior: func(m *mockUsersRepository) {},
			expectedErr:  true,
		},
		{
			name:         "Error - repository",
			fullName:     "Nikita",
			phoneNumber:  &phone,
			mockBehavior: func(m *mockUsersRepository) {},
			expectedErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockUsersRepository{}

			if tt.mockBehavior != nil {
				tt.mockBehavior(mockRepo)
			}

			service := &UsersService{
				usersRepository: mockRepo,
			}

			actualUser, err := service.CreateUser(
				bgCtx,
				tt.fullName,
				tt.phoneNumber,
			)

			if tt.expectedErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(actualUser, tt.expectedUser) {
				t.Errorf(
					"expected user %+v, got %+v",
					tt.expectedUser,
					actualUser,
				)
			}
		})
	}
}
