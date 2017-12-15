package rest

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"

	"github.com/ld100/goblet/domain/users/models"
	"github.com/ld100/goblet/domain/users/repository/orm"
	"github.com/ld100/goblet/domain/users/service"
	"github.com/ld100/goblet/persistence"
	httperrors "github.com/ld100/goblet/server/rest/errors"
	"github.com/ld100/goblet/util/log"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte(os.Getenv("SECRET_KEY")), nil)
}

func SessionRouter() chi.Router {
	// Persistence/Data layers wiring
	dbConn := persistence.GormDB

	sessionRepo := orm.NewOrmSessionRepository(dbConn)
	sessionService := service.NewSessionService(sessionRepo)

	userRepo := orm.NewOrmUserRepository(dbConn)
	userService := service.NewUserService(userRepo)

	handler := &RESTSessionHandler{SService: sessionService, UService: userService}

	// RESTful routes
	r := chi.NewRouter()
	r.Post("/", handler.Store) // POST /session
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))

		// Handle valid / invalid tokens. In this example, we use
		// the provided authenticator middleware, but you can write your
		// own very easily, look at the Authenticator method in jwtauth.go
		// and tweak it, its not scary.
		r.Use(jwtauth.Authenticator)

		r.Delete("/", handler.Delete) // DELETE /session
	})
	return r
}

type RESTSessionHandler struct {
	SService service.SessionService
	UService service.UserService
}

// TODO: Move login and token issuing logic out of controller
func (handler *RESTSessionHandler) Store(w http.ResponseWriter, r *http.Request) {
	// Issue access token based on session
	// Generate response with access token

	data := &EmailLoginRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, httperrors.ErrInvalidRequest(err))
		return
	}

	user, err := handler.UService.GetByEmail(data.Email)
	if err != nil {
		render.Render(w, r, httperrors.ErrNotFound)
		return
	}

	if !models.CheckPasswordHash(data.Password, user.Password) {
		err := errors.New("provided password is wrong")
		render.Render(w, r, httperrors.ErrUnauthorized(err))
		return
	}

	hours, _ := strconv.Atoi(os.Getenv("SESSION_TTL_HOURS"))
	hoursDuration := time.Duration(hours)
	session := &models.Session{
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(hoursDuration * time.Hour),
	}
	session, err = handler.SService.Store(session)
	if err != nil {
		render.Render(w, r, httperrors.ErrInvalidRequest(err))
		return
	}

	claims := jwtauth.Claims{
		"iss":         "Goblet",
		"exp":         session.ExpiresAt.Unix(),
		"userID":      session.UserID,
		"userUUID":    user.Uuid,
		"sessionID":   session.ID,
		"sessionUUID": session.Uuid,
		"email":       user.Email,
		"firstName":   user.FirstName,
		"lastName":    user.LastName,
	}
	// Standard claims: https://tools.ietf.org/html/rfc7519#section-4.1
	_, tokenString, _ := tokenAuth.Encode(claims)

	log.Info("Created session", session)
	resp := &LoginResponse{Token: tokenString}
	render.Status(r, http.StatusCreated)
	render.Render(w, r, resp)
}

func (handler *RESTSessionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())

	idValue := claims["sessionID"]
	var id uint
	switch v := idValue.(type) {
	case float64:
		id = uint(int64(v))
	case int64:
		id = uint(v)
	default:
	}

	deleted, err := handler.SService.Delete(id)
	if !deleted {
		render.Render(w, r, httperrors.ErrInvalidRequest(err))
		return
	}
	log.Info("Deleted session: ", id)

	render.Status(r, http.StatusNoContent)
}
