package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type EsoService struct {
	storage pgsql.EsoStorage
}

type ifEsoStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Eso_count, error)
	Add(ctx context.Context, ea models.Eso) (int, error)
	Upd(ctx context.Context, eu models.Eso) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Eso_count, error)
}

//func NewEsoService(storage pgsql.EsoStorage) *EsoService
func NewEsoService(storage pgsql.EsoStorage) *EsoService {
	return &EsoService{storage}
}

//func (esv *EsoService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Eso_count, error)
func (esv *EsoService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Eso_count, error) {
	var est ifEsoStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("EsoStorage.GetList", err)
		return models.Eso_count{Values: []models.Eso{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *EsoService) Add(ctx context.Context, ea models.Eso) (int, error)
func (esv *EsoService) Add(ctx context.Context, ea models.Eso) (int, error) {
	var est ifEsoStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("EsoStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *EsoService) Upd(ctx context.Context, eu models.Eso) (int, error)
func (esv *EsoService) Upd(ctx context.Context, eu models.Eso) (int, error) {
	var est ifEsoStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("EsoStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *EsoService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *EsoService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifEsoStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("EsoStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *EsoService) GetOne(ctx context.Context, i int) (models.Eso_count, error)
func (esv *EsoService) GetOne(ctx context.Context, i int) (models.Eso_count, error) {
	var est ifEsoStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("EsoStorage.GetOne", err)
		return models.Eso_count{Values: []models.Eso{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
