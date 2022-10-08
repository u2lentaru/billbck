package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type PuService struct {
	storage pgsql.PuStorage
}

type ifPuStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2, gs3, gs4, gs5, gs6, gs7 string, ord int, dsc bool) (models.Pu_count, error)
	Add(ctx context.Context, ea models.Pu) (int, error)
	Upd(ctx context.Context, eu models.Pu) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Pu_count, error)
	GetObj(ctx context.Context, gs1, gs2 string) (models.Pu_count, error)
}

//func NewPuService(storage pgsql.Putorage) *PuService
func NewPuService(storage pgsql.PuStorage) *PuService {
	return &PuService{storage}
}

//func (esv *PuService) GetList(ctx context.Context, pg, pgs int, gs1, gs2, gs3, gs4, gs5, gs6, gs7 string, ord int, dsc bool) (models.Pu_count, error)
func (esv *PuService) GetList(ctx context.Context, pg, pgs int, gs1, gs2, gs3, gs4, gs5, gs6, gs7 string, ord int, dsc bool) (models.Pu_count, error) {
	var est ifPuStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, gs2, gs3, gs4, gs5, gs6, gs7, ord, dsc)

	if err != nil {
		log.Println("PuStorage.GetList", err)
		return models.Pu_count{Values: []models.Pu{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *PuService) Add(ctx context.Context, ea models.Pu) (int, error)
func (esv *PuService) Add(ctx context.Context, ea models.Pu) (int, error) {
	var est ifPuStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("PuStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *PuService) Upd(ctx context.Context, eu models.Pu) (int, error)
func (esv *PuService) Upd(ctx context.Context, eu models.Pu) (int, error) {
	var est ifPuStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("PuStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *PuService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *PuService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifPuStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("PuStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *PuService) GetOne(ctx context.Context, i int) (models.Pu_count, error)
func (esv *PuService) GetOne(ctx context.Context, i int) (models.Pu_count, error) {
	var est ifPuStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("PuStorage.GetOne", err)
		return models.Pu_count{Values: []models.Pu{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *PuService) GetObj(ctx context.Context, gs1, gs2 string) (models.Pu_count, error)
func (esv *PuService) GetObj(ctx context.Context, gs1, gs2 string) (models.Pu_count, error) {
	var est ifPuStorage
	est = &esv.storage

	out_count, err := est.GetObj(ctx, gs1, gs2)

	if err != nil {
		log.Println("PuStorage.GetOne", err)
		return models.Pu_count{Values: []models.Pu{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
