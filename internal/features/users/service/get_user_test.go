package users_service

import (
	"context"
	"errors"
	"practice/internal/core/domain"
	"testing"
)

func TestGetUser(t *testing.T) {
	bgCtx := context.Background()

	tests := []struct {
		name         string
		id           int
		mockBehavior func(m *mockUsersRepository)
		expectedUser domain.User
		expectedErr  bool
	}{
		{
			name: "Success - user get",
			id:   123,
			mockBehavior: func(m *mockUsersRepository) {
				m.getUserFunc = func(
					ctx context.Context,
					id int,
				) (domain.User, error) {
					if id != 123 {
						t.Errorf("expected id 123, got %d", id)
					}

					return domain.User{
						ID:       123,
						FullName: "Nikita",
					}, nil
				}
			},
			expectedUser: domain.User{
				ID:       123,
				FullName: "Nikita",
			},
			expectedErr: false,
		},
		{
			name: "Error - repository",
			id:   123,
			mockBehavior: func(m *mockUsersRepository) {
				m.getUserFunc = func(
					ctx context.Context,
					id int,
				) (domain.User, error) {
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

			service := &UsersService{
				usersRepository: mockRepo,
			}

			actualUser, err := service.GetUser(bgCtx, tt.id)

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

			if actualUser.ID != tt.expectedUser.ID {
				t.Errorf(
					"expected id %d, got %d",
					tt.expectedUser.ID,
					actualUser.ID,
				)
			}

			if actualUser.FullName != tt.expectedUser.FullName {
				t.Errorf(
					"expected full name %s, got %s",
					tt.expectedUser.FullName,
					actualUser.FullName,
				)
			}
		})
	}
}
