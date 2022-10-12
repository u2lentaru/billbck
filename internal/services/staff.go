package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type StaffService struct {
	storage pgsql.StaffStorage
}

type ifStaffStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Staff_count, error)
	Add(ctx context.Context, ea models.Staff) (int, error)
	Upd(ctx context.Context, eu models.Staff) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Staff_count, error)
}

//func NewStaffService(storage pgsql.Stafftorage) *StaffService
func NewStaffService(storage pgsql.StaffStorage) *StaffService {
	return &StaffService{storage}
}

//func (esv *StaffService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Staff_count, error)
func (esv *StaffService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Staff_count, error) {
	var est ifStaffStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("StaffStorage.GetList", err)
		return models.Staff_count{Values: []models.Staff{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *StaffService) Add(ctx context.Context, ea models.Staff) (int, error)
func (esv *StaffService) Add(ctx context.Context, ea models.Staff) (int, error) {
	var est ifStaffStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("StaffStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *StaffService) Upd(ctx context.Context, eu models.Staff) (int, error)
func (esv *StaffService) Upd(ctx context.Context, eu models.Staff) (int, error) {
	var est ifStaffStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("StaffStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *StaffService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *StaffService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifStaffStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("StaffStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *StaffService) GetOne(ctx context.Context, i int) (models.Staff_count, error)
func (esv *StaffService) GetOne(ctx context.Context, i int) (models.Staff_count, error) {
	var est ifStaffStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("StaffStorage.GetOne", err)
		return models.Staff_count{Values: []models.Staff{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
