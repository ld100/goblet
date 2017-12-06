package rest

import (
	"net/http"

	models "github.com/ld100/goblet/domain/users"
	"strings"
)

// UserRequest is the request payload for Article data model.
//
// NOTE: It's good practice to have well defined request and response payloads
// so you can manage the specific inputs and outputs for clients, and also gives
// you the opportunity to transform data on input or output, for example
// on request, we'd like to protect certain fields and on output perhaps
// we'd like to include a computed field based on other values that aren't
// in the data model. Also, check out this awesome blog post on struct composition:
// http://attilaolah.eu/2014/09/10/json-and-struct-composition-in-go/
type UserRequest struct {
	User *models.User

	ProtectedID string `json:"id"` // override 'id' json to have more control
}

func (u *UserRequest) Bind(r *http.Request) error {
	// just u post-process after u decode..
	u.ProtectedID = ""                           // unset the protected ID
	u.User.Email = strings.ToLower(u.User.Email) // as an example, we down-case
	return nil
}
