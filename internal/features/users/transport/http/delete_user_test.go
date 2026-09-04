package users_transport_http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteUser(t *testing.T) {
	bgCtx := testContext()

	tests := []struct {
		name           string
		userIdPath     string
		mockBehavior   func(m *mockUsersService)
		expectedStatus int
	}{
		{
			name:       "Success - user delete",
			userIdPath: "123",
			mockBehavior: func(m *mockUsersService) {
				m.deleteUserFunc = func(ctx context.Context, id int) error {
					return nil
				}
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:       "Error - server error",
			userIdPath: "404",
			mockBehavior: func(m *mockUsersService) {
				m.deleteUserFunc = func(ctx context.Context, id int) error {
					return errors.New("db error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
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
			mux.HandleFunc("DELETE /users/{id}", handler.DeleteUser)

			req := httptest.NewRequest(http.MethodDelete, "/users/"+tt.userIdPath, nil)
			req = req.WithContext(bgCtx)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
