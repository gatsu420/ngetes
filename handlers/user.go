package handlers

import (
	"net/http"

	"github.com/gatsu420/ngetes/models"
	"github.com/go-chi/render"
)

type UserOperations interface {
	CreateUser(u *models.User) error

	GetUserRole(roleModel *models.Role, roleName string) (roleID int, err error)
}

type UserHandlers struct {
	Operations UserOperations
}

func NewUserHandlers(operations UserOperations) *UserHandlers {
	return &UserHandlers{
		Operations: operations,
	}
}

type userResponse struct {
	User *models.User `json:"user"`
}

func newUserResponse(u *models.User) *userResponse {
	return &userResponse{
		User: u,
	}
}

type userRequest struct {
	User *models.User `json:"user"`
}

func (ur *userRequest) Bind(r *http.Request) error {
	return nil
}

func (hd *UserHandlers) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	user := &userRequest{}
	err := render.Bind(r, user)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	role := &models.Role{}
	roleID, err := hd.Operations.GetUserRole(role, user.User.RoleName)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	user.User.RoleID = roleID
	err = hd.Operations.CreateUser(user.User)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	render.Respond(w, r, newUserResponse(user.User))
}
