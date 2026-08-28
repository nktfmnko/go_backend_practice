package users_transport_http

import (
	"net/http"
	core_logger "practice/internal/core/logger"
	core_http_response "practice/internal/core/transport/http/response"
	core_http_utils "practice/internal/core/transport/http/utils"
)

func (h *UsersHTTPHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHttpResponseHandler(logger, w)

	userID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "fail to get user id")
		return
	}

	if err := h.usersService.DeleteUser(ctx, userID); err != nil {
		responseHandler.ErrorResponse(err, "fail to delete user")
		return
	}

	responseHandler.NoContentResponse()
}
