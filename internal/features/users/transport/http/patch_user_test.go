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

func TestPatchUser(t *testing.T) {
	phone := "+358401234567"
	tests := []struct {
		name           string
		userID         string
		body           string
		mockBehavior   func(m *mockUsersService)
		expectedStatus int
	}{
		{
			name:   "Success - patch user",
			userID: "123",
			body: `{
				"full_name": "Nikita",
				"phone_number": "+358401234567"
			}`,
			mockBehavior: func(m *mockUsersService) {
				m.patchUserFunc = func(
					ctx context.Context,
					id int,
					patch domain.UserPatch,
				) (domain.User, error) {
					if id != 123 {
						t.Errorf("expected user id 123, got %d", id)
					}

					if *patch.FullName.Value != "Nikita" || *patch.PhoneNumber.Value != phone {
						t.Errorf("expected name Nikita, got %v", patch.FullName.Value)
						t.Errorf("expected phone %v, got %v", phone, patch.PhoneNumber.Value)
					}

					return domain.User{
						ID:          123,
						FullName:    "Nikita",
						PhoneNumber: &phone,
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Success - patch only full name",
			userID: "123",
			body: `{
				"full_name": "Alex"
			}`,
			mockBehavior: func(m *mockUsersService) {
				m.patchUserFunc = func(
					ctx context.Context,
					id int,
					patch domain.UserPatch,
				) (domain.User, error) {
					if id != 123 {
						t.Errorf("expected user id 123, got %d", id)
					}

					return domain.User{
						ID:       123,
						FullName: "Alex",
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Success - patch only phone number",
			userID: "123",
			body: `{
				"phone_number": "+358401234567"
			}`,
			mockBehavior: func(m *mockUsersService) {
				m.patchUserFunc = func(
					ctx context.Context,
					id int,
					patch domain.UserPatch,
				) (domain.User, error) {
					if id != 123 {
						t.Errorf("expected user id 123, got %d", id)
					}

					return domain.User{
						ID:          123,
						PhoneNumber: &phone,
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Error - invalid user id",
			userID: "abc",
			body: `{
				"full_name": "Nikita"
			}`,
			mockBehavior:   func(m *mockUsersService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Error - invalid JSON",
			userID: "123",
			body: `{
				"full_name": "Nikita"`,
			mockBehavior:   func(m *mockUsersService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Error - invalid full name",
			userID: "123",
			body: `{
				"full_name": "Ni"
			}`,
			mockBehavior:   func(m *mockUsersService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Error - invalid phone number",
			userID: "123",
			body: `{
				"phone_number": "123"
			}`,
			mockBehavior:   func(m *mockUsersService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Error - service error",
			userID: "123",
			body: `{
				"full_name": "Nikita"
			}`,
			mockBehavior: func(m *mockUsersService) {
				m.patchUserFunc = func(
					ctx context.Context,
					id int,
					patch domain.UserPatch,
				) (domain.User, error) {
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
			mux.HandleFunc("PATCH /users/{id}", handler.PatchUser)

			req := httptest.NewRequest(
				http.MethodPatch,
				"/users/"+tt.userID,
				strings.NewReader(tt.body),
			)

			req.Header.Set("Content-Type", "application/json")
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
