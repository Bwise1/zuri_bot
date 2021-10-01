package users

import (
	"net/http"
)

type UserHandler struct {
	UserService
}

func NewUserHandler(us UserService) *UserHandler {
	return &UserHandler{us}
}

func (uh *UserHandler) NewUser(w http.ResponseWriter, r *http.Request) {
	user := userFromContext(r.Context())
	_ = uh.UserService.CreateUser(r.Context(), user)
}
