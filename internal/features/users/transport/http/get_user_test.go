package users_transport_http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"practice/internal/core/domain"
	core_error "practice/internal/core/errors"
	"testing"
)

func TestGetUser(t *testing.T) {
	bgCtx := testContext()

	tests := []struct {
		name           string
		userIdPath     string
		mockBehavior   func(m *mockUsersService)
		expectedStatus int
	}{
		{
			name:       "Success - user found",
			userIdPath: "123",
			mockBehavior: func(m *mockUsersService) {
				m.getUserFunc = func(ctx context.Context, id int) (domain.User, error) {
					return domain.User{
						ID:       id,
						Version:  1,
						FullName: "John doe",
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:       "Error - user not found",
			userIdPath: "404",
			mockBehavior: func(m *mockUsersService) {
				m.getUserFunc = func(ctx context.Context, id int) (domain.User, error) {
					return domain.User{}, core_error.ErrNotFound
				}
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Error - invalid id",
			userIdPath:     "abc",
			mockBehavior:   func(m *mockUsersService) {},
			expectedStatus: http.StatusBadRequest,
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
			mux.HandleFunc("GET /users/{id}", handler.GetUser)

			req := httptest.NewRequest(http.MethodGet, "/users/"+tt.userIdPath, nil)
			req = req.WithContext(bgCtx)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
