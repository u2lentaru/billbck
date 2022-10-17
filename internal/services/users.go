package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type UserService struct {
	storage pgsql.UserStorage
}

type ifUserStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.User_count, error)
	Add(ctx context.Context, ea models.User) (int, error)
	Upd(ctx context.Context, eu models.User) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.User_count, error)
}

//func NewUserService(storage pgsql.Usertorage) *UserService
func NewUserService(storage pgsql.UserStorage) *UserService {
	return &UserService{storage}
}

//func (esv *UserService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.User_count, error)
func (esv *UserService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.User_count, error) {
	var est ifUserStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("UserStorage.GetList", err)
		return models.User_count{Values: []models.User{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *UserService) Add(ctx context.Context, ea models.User) (int, error)
func (esv *UserService) Add(ctx context.Context, ea models.User) (int, error) {
	var est ifUserStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("UserStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *UserService) Upd(ctx context.Context, eu models.User) (int, error)
func (esv *UserService) Upd(ctx context.Context, eu models.User) (int, error) {
	var est ifUserStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("UserStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *UserService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *UserService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifUserStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("UserStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *UserService) GetOne(ctx context.Context, i int) (models.User_count, error)
func (esv *UserService) GetOne(ctx context.Context, i int) (models.User_count, error) {
	var est ifUserStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("UserStorage.GetOne", err)
		return models.User_count{Values: []models.User{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
