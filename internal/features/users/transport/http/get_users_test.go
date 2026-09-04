package users_transport_http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"practice/internal/core/domain"
	"testing"
)

func TestGetUsers(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		mockBehavior   func(m *mockUsersService)
		expectedStatus int
	}{
		{
			name:  "Success - get users",
			query: "?limit=10&offset=20",
			mockBehavior: func(m *mockUsersService) {
				m.getUsersFunc = func(
					ctx context.Context,
					limit *int,
					offset *int,
				) ([]domain.User, error) {
					if *limit != 10 {
						t.Errorf("expected limit 10, got %d", *limit)
					}

					if *offset != 20 {
						t.Errorf("expected offset 20, got %d", *offset)
					}

					return []domain.User{
						{
							ID:       1,
							FullName: "Nikita",
						},
						{
							ID:       2,
							FullName: "Nikita",
						},
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error - invalid limit",
			query:          "?limit=abc&offset=20",
			mockBehavior:   func(m *mockUsersService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Error - invalid offset",
			query:          "?limit=10&offset=abc",
			mockBehavior:   func(m *mockUsersService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "Error - service error",
			query: "?limit=10&offset=20",
			mockBehavior: func(m *mockUsersService) {
				m.getUsersFunc = func(
					ctx context.Context,
					limit *int,
					offset *int,
				) ([]domain.User, error) {
					return nil, errors.New("db error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockUsersService{}
			if tt.mockBehavior != nil {
				tt.mockBehavior(mockSvc)
			}

			handler := NewUsersHTTPHandler(mockSvc)
			mux := http.NewServeMux()
			mux.HandleFunc("GET /users", handler.GetUsers)

			req := httptest.NewRequest(
				http.MethodGet,
				"/users"+tt.query,
				nil,
			)

			req = req.WithContext(testContext())

			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf(
					"expected status %d, got %d",
					tt.expectedStatus,
					w.Code,
				)
			}
		})
	}
}
