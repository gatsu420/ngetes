package handlers

import (
	"errors"
	"net/http"

	"github.com/gatsu420/ngetes/models"
	"github.com/go-chi/render"
)

type UserOperations interface {
	CreateUser(u *models.User) error
	GetUserNameExistence(userName string) (isExist bool, err error)

	ListRoles() ([]models.Role, error)
	GetRoleByRoleName(roleName string) (roleID int, err error)
	GetRoleByUserName(name string) (roleID int, err error)
}

type UserHandlers struct {
	Operations UserOperations
}

func NewUserHandlers(operations UserOperations) *UserHandlers {
	return &UserHandlers{
		Operations: operations,
	}
}

func (hd *UserHandlers) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	user := &userRequest{}
	err := render.Bind(r, user)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	roles, err := hd.Operations.ListRoles()
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	isRoleNameFieldExist := false
	if user.User.RoleName != "" {
		isRoleNameFieldExist = true
	}
	if !isRoleNameFieldExist {
		render.Render(w, r, errRender(errors.New("role_name does not exist in request payload")))
		return
	}

	isRoleNameExist := false
	for _, v := range roles {
		if user.User.RoleName == v.Name {
			isRoleNameExist = true
			break
		}
	}
	if !isRoleNameExist {
		render.Render(w, r, errRender(errors.New("role_name is wrong")))
		return
	}

	roleID, err := hd.Operations.GetRoleByRoleName(user.User.RoleName)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	isUserNameExist, err := hd.Operations.GetUserNameExistence(user.User.Name)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}
	if isUserNameExist {
		render.Render(w, r, errRender(errors.New("name already exists")))
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
