package rest

import (
	"net/http"
	"context"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/ld100/goblet/domain/users/persistence"
	"github.com/ld100/goblet/domain/users/service"
	"github.com/ld100/goblet/domain/users/repository"
	"github.com/ld100/goblet/domain/users/repository/orm"
	httperrors "github.com/ld100/goblet/server/rest/errors"
	models "github.com/ld100/goblet/domain/users"
)

func UserRouter() chi.Router {
	// Persistence/Data layers wiring
	dbConn := persistence.GormDB
	userRepo = orm.NewOrmUserRepository(dbConn)
	userService = service.NewUserService(userRepo)
	handler := &RESTUserHandler{UService: userService}

	// RESTful routes
	r := chi.NewRouter()
	r.Get("/", handler.GetAll)
	r.Post("/", handler.Store) // POST /articles
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
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, httperrors.ErrInvalidRequest(err))
		return
	}
	user = data.User
	user, err := handler.UService.Update(&user)
	if err != nil {
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

	deleted, err := handler.UService.Delete(&user.ID)
	if !deleted {
		render.Render(w, r, rest.ErrInvalidRequest(err))
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
			user, err = handler.UService.GetById(userID)
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
