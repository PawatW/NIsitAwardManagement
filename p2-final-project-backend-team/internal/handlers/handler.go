package handlers

import (
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/handlers/admin"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/handlers/auth"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/handlers/award"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/handlers/dev"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/handlers/request"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/handlers/organization"
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/handlers/user"
	academicTermHd "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/handlers/academic_term"
	
)

type Handler struct {
	User    user.Handler
	Admin   admin.AdminHandler
	Auth    auth.AuthHandler
	Award   award.AwardHandler
	Request request.RequestHandler
	Dev     dev.DevHandler
	Organization organization.OrganizationHandler
	AcademicTerm academicTermHd.AcademicTermHandler

}

func NewHandlers(
	admin admin.AdminHandler,
	auth auth.AuthHandler,
	user user.Handler,
	award award.AwardHandler,
	request request.RequestHandler,
	dev dev.DevHandler,
	organization organization.OrganizationHandler,
	academicTerm academicTermHd.AcademicTermHandler,

) *Handler {
	return &Handler{
		User:    user,
		Admin:   admin,
		Auth:    auth,
		Award:   award,
		Request: request,
		Dev:     dev,
		Organization: organization,
		AcademicTerm: academicTerm,
	}
}
