package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	models "github.com/ld100/goblet/domain/users"
	"github.com/ld100/goblet/domain/users/repository/orm"
	"github.com/ld100/goblet/domain/users/service"
	"github.com/ld100/goblet/persistence"
	httperrors "github.com/ld100/goblet/server/rest/errors"
	"github.com/ld100/goblet/util/log"
)

func UserRouter() chi.Router {
	// Persistence/Data layers wiring
	dbConn := persistence.GormDB
	userRepo := orm.NewOrmUserRepository(dbConn)
	userService := service.NewUserService(userRepo)
	handler := &RESTUserHandler{UService: userService}

	// RESTful routes
	r := chi.NewRouter()
	r.Get("/", handler.GetAll)
	r.Post("/", handler.Store) // POST /user
	r.Route("/{userID}", func(r chi.Router) {
		r.Use(handler.userCtx)        // Load the *User on the request context
		r.Get("/", handler.GetByID)   // GET /users/123
		r.Put("/", handler.Update)    // PUT /users/123
		r.Delete("/", handler.Delete) // DELETE /users/123
	})
	return r
}

type RESTUserHandler struct {
	UService service.UserService
}

func (handler *RESTUserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Assume if we've reach this far, we can access the user
	// context because this handler is a child of the userCtx
	// middleware. The worst case, the recoverer middleware will save us.
	user := r.Context().Value("user").(*models.User)

	if err := render.Render(w, r, NewUserResponse(user)); err != nil {
		render.Render(w, r, httperrors.ErrRender(err))
		return
	}
}

func (handler *RESTUserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, error := handler.UService.GetAll()
	if error != nil {
		render.Render(w, r, httperrors.ErrRender(error))
		return
	}

	if err := render.RenderList(w, r, NewUserListResponse(users)); err != nil {
		render.Render(w, r, httperrors.ErrRender(err))
		return
	}
}

func (handler *RESTUserHandler) Update(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)

	data := &UserRequest{User: user}

	//textbody, _ := strconv.Unquote(r.Body)
	//jsonErr := json.NewDecoder(textbody).Decode(data)
	jsonErr := json.NewDecoder(r.Body).Decode(data)
	log.Debug("Custom JSON parsing:", jsonErr)

	if err := render.Bind(r, data); err != nil {
		log.Debug("Could not map user request to model", err)
		render.Render(w, r, httperrors.ErrInvalidRequest(err))
		return
	}

	user = data.User
	user, err := handler.UService.Update(user)
	if err != nil {
		log.Debug("Could not update user entity")
		render.Render(w, r, httperrors.ErrInvalidRequest(err))
		return
	}

	render.Render(w, r, NewUserResponse(user))
}

func (handler *RESTUserHandler) Store(w http.ResponseWriter, r *http.Request) {
	data := &UserRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, httperrors.ErrInvalidRequest(err))
		return
	}

	var err error
	user := data.User
	user, err = handler.UService.Store(user)
	if err != nil {
		render.Render(w, r, httperrors.ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewUserResponse(user))
}

func (handler *RESTUserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)

	deleted, err := handler.UService.Delete(user.ID)
	if !deleted {
		render.Render(w, r, httperrors.ErrInvalidRequest(err))
		return
	}

	render.Render(w, r, NewUserResponse(user))
}

// UserCtx middleware is used to load an User object from
// the URL parameters passed through as the request. In case
// the User could not be found, we stop here and return a 404.
func (handler *RESTUserHandler) userCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *models.User
		var err error

		if userID := chi.URLParam(r, "userID"); userID != "" {
			// TODO: Check whether userID is integer actually
			intUserId, _ := strconv.Atoi(userID)
			user, err = handler.UService.GetByID(uint(intUserId))
		} else {
			render.Render(w, r, httperrors.ErrNotFound)
			return
		}

		if err != nil {
			render.Render(w, r, httperrors.ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
