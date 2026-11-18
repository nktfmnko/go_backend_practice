package users_transport_http

import (
	"net/http"
	core_logger "practice/internal/core/logger"
	core_auth "practice/internal/core/transport/auth"
	core_http_request "practice/internal/core/transport/http/request"

	core_http_response "practice/internal/core/transport/http/response"
)

func (h *UsersHTTPHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHttpResponseHandler(logger, w)

	claims, err := core_auth.MustUserFromContext(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "unauthenticated")
		return
	}

	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "fail to get user id")
		return
	}

	if err := h.usersService.DeleteUser(ctx, claims.UserID, userID); err != nil {
		responseHandler.ErrorResponse(err, "fail to delete user")
		return
	}

	responseHandler.NoContentResponse()
}
