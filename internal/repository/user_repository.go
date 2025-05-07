package repository

import (
	"x-ui-monitor/internal/domain"
	"x-ui-monitor/pkg/bolt"
)

type UserRepository struct {
	db         *bolt.BoltDB
	bucketName string
}

func NewUserRepository(db *bolt.BoltDB, bucketName string) *UserRepository {
	return &UserRepository{
		db:         db,
		bucketName: bucketName,
	}
}

func (r *UserRepository) AddUser(user domain.User) {
	r.db.AddIP(user.InboundTag, user.IP)
}

func (r *UserRepository) GetTotalUsersCount() (*domain.TotalUsersCountResult, error) {
	return r.db.GetTotalUsersCount()
}

func (r *UserRepository) GetTotalUsersByInbound(inbound string) (int, error) {
	return r.db.GetActiveIPCount(inbound)
}

func (r *UserRepository) GetActiveIPsByInbound(inbound string) ([]string, error) {
	return r.db.ListActiveIPs(inbound)
}
