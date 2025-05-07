package usecase

import (
	"x-ui-monitor/internal/domain"
	"x-ui-monitor/internal/repository"
)

type UserUsecase struct {
	repo *repository.UserRepository
}

func NewUserUsecase(repo *repository.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (u *UserUsecase) AddUser(inboundTag, ip string) {
	user := domain.User{IP: ip, InboundTag: inboundTag}
	u.repo.AddUser(user)
}

func (u *UserUsecase) GetTotalUsersCount() (*domain.TotalUsersCountResult, error) {
	return u.repo.GetTotalUsersCount()
}

func (u *UserUsecase) GetTotalUsersByInbound(inbound string) (int, error) {
	return u.repo.GetTotalUsersByInbound(inbound)
}

func (u *UserUsecase) GetActiveIpsByInbound(inbound string) ([]string, error) {
	return u.repo.GetActiveIPsByInbound(inbound)
}
