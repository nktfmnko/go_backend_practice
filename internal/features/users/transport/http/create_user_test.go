package users_transport_http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"practice/internal/core/domain"
	"strings"
	"testing"
)

func TestCreateUser(t *testing.T) {
	bgCtx := testContext()
	phone := "+12345678901"

	tests := []struct {
		name           string
		body           string
		mockBehavior   func(m *mockUsersService)
		expectedStatus int
	}{
		{
			name: "Success - user create",
			body: `{
				"full_name": "Nikita",
				"phone_number": "+12345678901"
			}`,
			mockBehavior: func(m *mockUsersService) {
				m.createUserFunc = func(ctx context.Context, fullName string, phoneNumber *string) (domain.User, error) {
					if fullName != "Nikita" {
						t.Errorf("expected full name Nikita, got %s", fullName)
					}
					if phoneNumber == nil || *phoneNumber != phone {
						t.Errorf("unexpected phone number")
					}
					return domain.User{ID: 123, FullName: fullName, PhoneNumber: phoneNumber}, nil
				}
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Error - invalid JSON",
			body: `{
				"full_name": "Nikita"
			`,
			mockBehavior:   func(m *mockUsersService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Error - validation",
			body: `{
				"full_name": "N"
			}`,
			mockBehavior:   func(m *mockUsersService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Error - server error",
			body: `{
				"full_name": "Nikita"
			}`,
			mockBehavior: func(m *mockUsersService) {
				m.createUserFunc = func(ctx context.Context, fullName string, phoneNumber *string) (domain.User, error) {
					return domain.User{}, errors.New("db error")
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
			mux.HandleFunc("POST /users", handler.CreateUser)

			req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(tt.body))
			req = req.WithContext(bgCtx)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
