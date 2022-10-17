package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type UzoService struct {
	storage pgsql.UzoStorage
}

type ifUzoStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Uzo_count, error)
	Add(ctx context.Context, ea models.Uzo) (int, error)
	Upd(ctx context.Context, eu models.Uzo) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Uzo_count, error)
}

//func NewUzoService(storage pgsql.Uzotorage) *UzoService
func NewUzoService(storage pgsql.UzoStorage) *UzoService {
	return &UzoService{storage}
}

//func (esv *UzoService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Uzo_count, error)
func (esv *UzoService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Uzo_count, error) {
	var est ifUzoStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("UzoStorage.GetList", err)
		return models.Uzo_count{Values: []models.Uzo{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *UzoService) Add(ctx context.Context, ea models.Uzo) (int, error)
func (esv *UzoService) Add(ctx context.Context, ea models.Uzo) (int, error) {
	var est ifUzoStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("UzoStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *UzoService) Upd(ctx context.Context, eu models.Uzo) (int, error)
func (esv *UzoService) Upd(ctx context.Context, eu models.Uzo) (int, error) {
	var est ifUzoStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("UzoStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *UzoService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *UzoService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifUzoStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("UzoStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *UzoService) GetOne(ctx context.Context, i int) (models.Uzo_count, error)
func (esv *UzoService) GetOne(ctx context.Context, i int) (models.Uzo_count, error) {
	var est ifUzoStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("UzoStorage.GetOne", err)
		return models.Uzo_count{Values: []models.Uzo{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
