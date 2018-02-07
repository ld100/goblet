package rest

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"

	"github.com/ld100/goblet/pkg/domain/user/model"
	"github.com/ld100/goblet/pkg/domain/user/repository/orm"
	"github.com/ld100/goblet/pkg/domain/user/service"
	httperrors "github.com/ld100/goblet/pkg/server/rest/error"
	"github.com/ld100/goblet/pkg/util/log"
	"github.com/ld100/goblet/pkg/server/env"
)

func init() {
	tokenAuth = jwtauth.New("HS256", []byte(os.Getenv("SECRET_KEY")), nil)
}

func UserRouter(env *env.Env) chi.Router {
	// Persistence/Data layers wiring
	dbConn, err := env.DB.ORMConnection()
	if err != nil {
		log.Fatal("cannot connect UserRouter to the database", err)
	}

	userRepo := orm.NewOrmUserRepository(dbConn)
	userService := service.NewUserService(userRepo)
	handler := &RESTUserHandler{UService: userService}

	// RESTful routes
	r := chi.NewRouter()
	r.Get("/", handler.GetAll)
	r.Post("/", handler.Store) // POST /user
	r.Route("/{userID}", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Use(handler.userCtx)        // Load the *User on the request context
		r.Get("/", handler.GetByID)   // GET /user/123
		r.Put("/", handler.Update)    // PUT /user/123
		r.Delete("/", handler.Delete) // DELETE /user/123
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
	user := r.Context().Value("user").(*model.User)

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
	user := r.Context().Value("user").(*model.User)

	data := &UserRequest{User: user}

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
	user := r.Context().Value("user").(*model.User)

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
		var user *model.User
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
