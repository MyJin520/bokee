package base

import (
	"gin-admin/internal/mods/request"
)

type UserService struct{}

func (s *UserService) Register(req request.UserRegisterOrEditReq) error {
	return nil
}
