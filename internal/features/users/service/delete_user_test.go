package users_service

import (
	"context"
	"errors"
	"testing"
)

func TestDeleteUser(t *testing.T) {
	bgCtx := context.Background()

	tests := []struct {
		name         string
		id           int
		mockBehavior func(m *mockUsersRepository)
		expectedErr  bool
	}{
		{
			name: "Success - user delete",
			id:   123,
			mockBehavior: func(m *mockUsersRepository) {
				m.deleteUserFunc = func(ctx context.Context, id int) error {
					if id != 123 {
						t.Errorf("expected id 123, got %d", id)
					}

					return nil
				}
			},
			expectedErr: false,
		},
		{
			name: "Error - repository",
			id:   123,
			mockBehavior: func(m *mockUsersRepository) {
				m.deleteUserFunc = func(ctx context.Context, id int) error {
					return errors.New("db error")
				}
			},
			expectedErr: true,
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

			err := service.DeleteUser(bgCtx, tt.id)

			if tt.expectedErr {
				if err == nil {
					t.Error("expected error, got nil")
				}

				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
