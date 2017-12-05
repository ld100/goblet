package articles

import (
	"context"
	"strings"
	"fmt"
	"errors"
	"math/rand"
	"net/http"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	rest "github.com/ld100/goblet/server/rest/errors"
)

func ArticleRouter() chi.Router {
	r := chi.NewRouter()
	r.With(paginate).Get("/", ListArticles)
	r.Post("/", CreateArticle)       // POST /articles
	r.Get("/search", SearchArticles) // GET /articles/search

	r.Route("/{articleID}", func(r chi.Router) {
		r.Use(ArticleCtx)            // Load the *Article on the request context
		r.Get("/", GetArticle)       // GET /articles/123
		r.Put("/", UpdateArticle)    // PUT /articles/123
		r.Delete("/", DeleteArticle) // DELETE /articles/123
	})

	// GET /articles/whats-up
	r.With(ArticleCtx).Get("/{articleSlug:[a-z-]+}", GetArticle)
	return r
}

// paginate is a stub, but very possible to implement middleware logic
// to handle the request params for handling a paginated request.
func paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// just a stub.. some ideas are to look at URL query params for something like
		// the page number, or the limit, and send a query cursor down the chain
		next.ServeHTTP(w, r)
	})
}

func ListArticles(w http.ResponseWriter, r *http.Request) {
	if err := render.RenderList(w, r, NewArticleListResponse(articles)); err != nil {
		render.Render(w, r, rest.ErrRender(err))
		return
	}
}

// ArticleCtx middleware is used to load an Article object from
// the URL parameters passed through as the request. In case
// the Article could not be found, we stop here and return a 404.
func ArticleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var article *Article
		var err error

		if articleID := chi.URLParam(r, "articleID"); articleID != "" {
			article, err = dbGetArticle(articleID)
		} else if articleSlug := chi.URLParam(r, "articleSlug"); articleSlug != "" {
			article, err = dbGetArticleBySlug(articleSlug)
		} else {
			render.Render(w, r, rest.ErrNotFound)
			return
		}
		if err != nil {
			render.Render(w, r, rest.ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "article", article)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// SearchArticles searches the Articles data for a matching article.
// It's just a stub, but you get the idea.
func SearchArticles(w http.ResponseWriter, r *http.Request) {
	render.RenderList(w, r, NewArticleListResponse(articles))
}

// CreateArticle persists the posted Article and returns it
// back to the client as an acknowledgement.
func CreateArticle(w http.ResponseWriter, r *http.Request) {
	data := &ArticleRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, rest.ErrInvalidRequest(err))
		return
	}

	article := data.Article
	dbNewArticle(article)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewArticleResponse(article))
}

// GetArticle returns the specific Article. You'll notice it just
// fetches the Article right off the context, as its understood that
// if we made it this far, the Article must be on the context. In case
// its not due to a bug, then it will panic, and our Recoverer will save us.
func GetArticle(w http.ResponseWriter, r *http.Request) {
	// Assume if we've reach this far, we can access the article
	// context because this handler is a child of the ArticleCtx
	// middleware. The worst case, the recoverer middleware will save us.
	article := r.Context().Value("article").(*Article)

	if err := render.Render(w, r, NewArticleResponse(article)); err != nil {
		render.Render(w, r, rest.ErrRender(err))
		return
	}
}

// UpdateArticle updates an existing Article in our persistent store.
func UpdateArticle(w http.ResponseWriter, r *http.Request) {
	article := r.Context().Value("article").(*Article)

	data := &ArticleRequest{Article: article}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, rest.ErrInvalidRequest(err))
		return
	}
	article = data.Article
	dbUpdateArticle(article.ID, article)

	render.Render(w, r, NewArticleResponse(article))
}

// DeleteArticle removes an existing Article from our persistent store.
func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	var err error

	// Assume if we've reach this far, we can access the article
	// context because this handler is a child of the ArticleCtx
	// middleware. The worst case, the recoverer middleware will save us.
	article := r.Context().Value("article").(*Article)

	article, err = dbRemoveArticle(article.ID)
	if err != nil {
		render.Render(w, r, rest.ErrInvalidRequest(err))
		return
	}

	render.Render(w, r, NewArticleResponse(article))
}

// ArticleRequest is the request payload for Article data model.
//
// NOTE: It's good practice to have well defined request and response payloads
// so you can manage the specific inputs and outputs for clients, and also gives
// you the opportunity to transform data on input or output, for example
// on request, we'd like to protect certain fields and on output perhaps
// we'd like to include a computed field based on other values that aren't
// in the data model. Also, check out this awesome blog post on struct composition:
// http://attilaolah.eu/2014/09/10/json-and-struct-composition-in-go/
type ArticleRequest struct {
	*Article

	User *UserPayload `json:"user,omitempty"`

	ProtectedID string `json:"id"` // override 'id' json to have more control
}

func (a *ArticleRequest) Bind(r *http.Request) error {
	// just a post-process after a decode..
	a.ProtectedID = ""                                 // unset the protected ID
	a.Article.Title = strings.ToLower(a.Article.Title) // as an example, we down-case
	return nil
}

// ArticleResponse is the response payload for the Article data model.
// See NOTE above in ArticleRequest as well.
//
// In the ArticleResponse object, first a Render() is called on itself,
// then the next field, and so on, all the way down the tree.
// Render is called in top-down order, like a http handler middleware chain.
type ArticleResponse struct {
	*Article

	User *UserPayload `json:"user,omitempty"`

	// We add an additional field to the response here.. such as this
	// elapsed computed property
	Elapsed int64 `json:"elapsed"`
}

func NewArticleResponse(article *Article) *ArticleResponse {
	resp := &ArticleResponse{Article: article}

	if resp.User == nil {
		if user, _ := dbGetUser(resp.UserID); user != nil {
			resp.User = NewUserPayloadResponse(user)
		}
	}

	return resp
}

func (rd *ArticleResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	rd.Elapsed = 10
	return nil
}

type ArticleListResponse []*ArticleResponse

func NewArticleListResponse(articles []*Article) []render.Renderer {
	list := []render.Renderer{}
	for _, article := range articles {
		list = append(list, NewArticleResponse(article))
	}
	return list
}

// NOTE: as a thought, the request and response payloads for an Article could be the
// same payload type, perhaps will do an example with it as well.
// type ArticlePayload struct {
//   *Article
// }

// Article data model. I suggest looking at https://upper.io for an easy
// and powerful data persistence adapter.
type Article struct {
	ID     string `json:"id"`
	UserID int64  `json:"user_id"` // the author
	Title  string `json:"title"`
	Slug   string `json:"slug"`
}

// Article fixture data
var articles = []*Article{
	{ID: "1", UserID: 100, Title: "Hi", Slug: "hi"},
	{ID: "2", UserID: 200, Title: "sup", Slug: "sup"},
	{ID: "3", UserID: 300, Title: "alo", Slug: "alo"},
	{ID: "4", UserID: 400, Title: "bonjour", Slug: "bonjour"},
	{ID: "5", UserID: 500, Title: "whats up", Slug: "whats-up"},
}

func dbNewArticle(article *Article) (string, error) {
	article.ID = fmt.Sprintf("%d", rand.Intn(100)+10)
	articles = append(articles, article)
	return article.ID, nil
}

func dbGetArticle(id string) (*Article, error) {
	for _, a := range articles {
		if a.ID == id {
			return a, nil
		}
	}
	return nil, errors.New("article not found")
}

func dbGetArticleBySlug(slug string) (*Article, error) {
	for _, a := range articles {
		if a.Slug == slug {
			return a, nil
		}
	}
	return nil, errors.New("article not found")
}

func dbUpdateArticle(id string, article *Article) (*Article, error) {
	for i, a := range articles {
		if a.ID == id {
			articles[i] = article
			return article, nil
		}
	}
	return nil, errors.New("article not found")
}

func dbRemoveArticle(id string) (*Article, error) {
	for i, a := range articles {
		if a.ID == id {
			articles = append((articles)[:i], (articles)[i+1:]...)
			return a, nil
		}
	}
	return nil, errors.New("article not found")
}

func dbGetUser(id int64) (*User, error) {
	for _, u := range users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

//--
// Request and Response payloads for the REST api.
//
// The payloads embed the data model objects an
//
// In a real-world project, it would make sense to put these payloads
// in another file, or another sub-package.
//--

type UserPayload struct {
	*User
	Role string `json:"role"`
}

func NewUserPayloadResponse(user *User) *UserPayload {
	return &UserPayload{User: user}
}

// Bind on UserPayload will run after the unmarshalling is complete, its
// a good time to focus some post-processing after a decoding.
func (u *UserPayload) Bind(r *http.Request) error {
	return nil
}

func (u *UserPayload) Render(w http.ResponseWriter, r *http.Request) error {
	u.Role = "collaborator"
	return nil
}

//--
// Data model objects and persistence mocks:
//--

// User data model
type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// User fixture data
var users = []*User{
	{ID: 100, Name: "Peter"},
	{ID: 200, Name: "Julia"},
}
