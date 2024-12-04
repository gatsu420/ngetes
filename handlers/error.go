package handlers

import (
	"net/http"

	"github.com/go-chi/render"
)

type errResponse struct {
	Err            error  `json:"-"`
	HTTPStatusCode int    `json:"-"`
	StatusText     string `json:"status"`
	ErrorText      string `json:"error,omitempty"`
}

func (e *errResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)

	return nil
}

func errRender(err error) render.Renderer {
	return &errResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnprocessableEntity,
		StatusText:     "Error rendering response",
		ErrorText:      err.Error(),
	}
}

func errUnauthorizedRender() render.Renderer {
	return &errResponse{
		HTTPStatusCode: http.StatusUnauthorized,
		StatusText:     "wrong name",
	}
}

func errForbiddenRender() render.Renderer {
	return &errResponse{
		HTTPStatusCode: http.StatusForbidden,
		StatusText:     "access forbidden",
	}
}
