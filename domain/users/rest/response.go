package rest

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/ld100/goblet/domain/users/models"
)

type omit *struct{}

// UserResponse is the response payload for the User data model.
// See NOTE above in UserRequest as well.
//
// In the UserResponse object, first a Render() is called on itself,
// then the next field, and so on, all the way down the tree.
// Render is called in top-down order, like a http handler middleware chain.
type UserResponse struct {
	*models.User

	// Lets omit password in the response
	Password omit `json:"password,omitempty"`

	// We add an additional field to the response here.. such as this
	// elapsed computed property
	Elapsed int64 `json:"elapsed"`
}

func (ur *UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	ur.Elapsed = 10
	return nil
}

func NewUserResponse(user *models.User) *UserResponse {
	resp := &UserResponse{User: user}
	// Could attach any shit to our payload here
	return resp
}

type UserListResponse []*UserResponse

func NewUserListResponse(users []*models.User) []render.Renderer {
	list := []render.Renderer{}
	for _, user := range users {
		list = append(list, NewUserResponse(user))
	}
	return list
}

type LoginResponse struct {
	Token string `json:"token,omitempty"`
}

func (lr *LoginResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}