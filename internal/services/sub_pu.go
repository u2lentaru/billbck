package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type SubPuService struct {
	storage pgsql.SubPuStorage
}

type ifSubPuStorage interface {
	GetList(ctx context.Context, pg, pgs, gs1, ord int, dsc bool) (models.Pu_count, error)
	Add(ctx context.Context, ea models.SubPu) (int, error)
	Upd(ctx context.Context, eu models.SubPu) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.SubPu_count, error)
	GetPrl(ctx context.Context, pg, pgs, gs1, gs2, ord int, dsc bool) (models.Pu_count, error)
}

//func NewSubPuService(storage pgsql.SubPutorage) *SubPuService
func NewSubPuService(storage pgsql.SubPuStorage) *SubPuService {
	return &SubPuService{storage}
}

//func (esv *SubPuService) GetList(ctx context.Context, pg, pgs, gs1, ord int, dsc bool) (models.Pu_count, error)
func (esv *SubPuService) GetList(ctx context.Context, pg, pgs, gs1, ord int, dsc bool) (models.Pu_count, error) {
	var est ifSubPuStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("SubPuStorage.GetList", err)
		return models.Pu_count{Values: []models.Pu{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *SubPuService) Add(ctx context.Context, ea models.SubPu) (int, error)
func (esv *SubPuService) Add(ctx context.Context, ea models.SubPu) (int, error) {
	var est ifSubPuStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("SubPuStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *SubPuService) Upd(ctx context.Context, eu models.SubPu) (int, error)
func (esv *SubPuService) Upd(ctx context.Context, eu models.SubPu) (int, error) {
	var est ifSubPuStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("SubPuStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *SubPuService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *SubPuService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifSubPuStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("SubPuStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *SubPuService) GetOne(ctx context.Context, i int) (models.SubPu_count, error)
func (esv *SubPuService) GetOne(ctx context.Context, i int) (models.SubPu_count, error) {
	var est ifSubPuStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("SubPuStorage.GetOne", err)
		return models.SubPu_count{Values: []models.SubPu{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *SubPuService) GetPrl(ctx context.Context, pg, pgs, gs1, gs2, ord int, dsc bool) (models.Pu_count, error)
func (esv *SubPuService) GetPrl(ctx context.Context, pg, pgs, gs1, gs2, ord int, dsc bool) (models.Pu_count, error) {
	var est ifSubPuStorage
	est = &esv.storage

	out_count, err := est.GetPrl(ctx, pg, pgs, gs1, gs2, ord, dsc)

	if err != nil {
		log.Println("SubPuStorage.GetList", err)
		return models.Pu_count{Values: []models.Pu{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
