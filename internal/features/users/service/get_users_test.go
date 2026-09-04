package users_service

import (
	"context"
	"errors"
	"practice/internal/core/domain"
	"testing"
)

func TestGetUsers(t *testing.T) {
	bgCtx := context.Background()

	limit := 10
	offset := 20

	tests := []struct {
		name          string
		limit         *int
		offset        *int
		mockBehavior  func(m *mockUsersRepository)
		expectedUsers []domain.User
		expectedErr   bool
	}{
		{
			name:   "Success - users get",
			limit:  nil,
			offset: nil,
			mockBehavior: func(m *mockUsersRepository) {
				m.getUsersFunc = func(
					ctx context.Context,
					limit *int,
					offset *int,
				) ([]domain.User, error) {
					if limit != nil {
						t.Errorf("expected nil limit, got %d", *limit)
					}

					if offset != nil {
						t.Errorf("expected nil offset, got %d", *offset)
					}

					return []domain.User{
						{
							ID:       1,
							FullName: "Nikita",
						},
						{
							ID:       2,
							FullName: "Alex",
						},
					}, nil
				}
			},
			expectedUsers: []domain.User{
				{
					ID:       1,
					FullName: "Nikita",
				},
				{
					ID:       2,
					FullName: "Alex",
				},
			},
		},
		{
			name:   "Success - users get with pagination",
			limit:  &limit,
			offset: &offset,
			mockBehavior: func(m *mockUsersRepository) {
				m.getUsersFunc = func(
					ctx context.Context,
					limit *int,
					offset *int,
				) ([]domain.User, error) {
					if limit == nil || *limit != 10 {
						t.Errorf("expected limit 10")
					}

					if offset == nil || *offset != 20 {
						t.Errorf("expected offset 20")
					}

					return []domain.User{
						{
							ID:       1,
							FullName: "Nikita",
						},
					}, nil
				}
			},
			expectedUsers: []domain.User{
				{
					ID:       1,
					FullName: "Nikita",
				},
			},
		},
		{
			name: "Error - negative limit",
			limit: func() *int {
				v := -1
				return &v
			}(),
			offset: nil,
			mockBehavior: func(m *mockUsersRepository) {
				m.getUsersFunc = func(
					ctx context.Context,
					limit *int,
					offset *int,
				) ([]domain.User, error) {
					t.Error("repository should not be called")
					return nil, nil
				}
			},
			expectedErr: true,
		},
		{
			name:  "Error - negative offset",
			limit: nil,
			offset: func() *int {
				v := -1
				return &v
			}(),
			mockBehavior: func(m *mockUsersRepository) {
				m.getUsersFunc = func(
					ctx context.Context,
					limit *int,
					offset *int,
				) ([]domain.User, error) {
					t.Error("repository should not be called")
					return nil, nil
				}
			},
			expectedErr: true,
		},
		{
			name:   "Error - repository",
			limit:  &limit,
			offset: &offset,
			mockBehavior: func(m *mockUsersRepository) {
				m.getUsersFunc = func(
					ctx context.Context,
					limit *int,
					offset *int,
				) ([]domain.User, error) {
					return nil, errors.New("db error")
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

			actualUsers, err := service.GetUsers(
				bgCtx,
				tt.limit,
				tt.offset,
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

			if len(actualUsers) != len(tt.expectedUsers) {
				t.Errorf(
					"expected %d users, got %d",
					len(tt.expectedUsers),
					len(actualUsers),
				)

				return
			}

			for i := range tt.expectedUsers {
				if actualUsers[i].ID != tt.expectedUsers[i].ID {
					t.Errorf(
						"expected user ID %d, got %d",
						tt.expectedUsers[i].ID,
						actualUsers[i].ID,
					)
				}

				if actualUsers[i].FullName != tt.expectedUsers[i].FullName {
					t.Errorf(
						"expected full name %s, got %s",
						tt.expectedUsers[i].FullName,
						actualUsers[i].FullName,
					)
				}
			}
		})
	}
}
