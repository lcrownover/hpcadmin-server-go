package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/lcrownover/hpcadmin-server/internal/types"
)

func Run(ctx context.Context) {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	dbConn := ctx.Value(types.DBKey).(*sql.DB)
	u := &UserHandler{dbConn: dbConn}

	r.Route("/users", func(r chi.Router) {
		r.Get("/", u.GetAllUsers)
		r.Post("/", u.CreateUser)
		r.Route("/{username}", func(r chi.Router) {
			r.Use(u.UserCtx)
			r.Get("/", u.GetUserByUsername)
		})
	})

	fmt.Println("Listening on :3333")
	http.ListenAndServe(":3333", r)

}

type APIResponse struct {
	Results any `json:"results"`
}

func NewAPIResponse(d any) *APIResponse {
	resp := &APIResponse{Results: d}
	return resp
}

func (a *APIResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

//--
// Error response payloads & renderers
//--

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}
var ErrInternalServer = &ErrResponse{HTTPStatusCode: 500, StatusText: "Internal server error."}
