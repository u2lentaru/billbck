package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type StreetService struct {
	storage pgsql.StreetStorage
}

type ifStreetStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, gs2, ord int, dsc bool) (models.Street_count, error)
	Add(ctx context.Context, ea models.Street) (int, error)
	Upd(ctx context.Context, eu models.Street) (int, error)
	Del(ctx context.Context, d models.StreetClose) (models.Json_id, error)
	GetOne(ctx context.Context, i int) (models.Street_count, error)
}

//NewStreetService(storage pg.StreetStorage) *StreetService
func NewStreetService(storage pgsql.StreetStorage) *StreetService {
	return &StreetService{storage}
}

//func (esv *StreetService) GetList(ctx context.Context, pg, pgs int, gs1 string, gs2, ord int, dsc bool) (models.Street_count, error)
func (esv *StreetService) GetList(ctx context.Context, pg, pgs int, gs1 string, gs2, ord int, dsc bool) (models.Street_count, error) {
	var est ifStreetStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, gs2, ord, dsc)

	if err != nil {
		log.Println("StreetStorage.GetList", err)
		return models.Street_count{Values: []models.Street{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *StreetService) Add(ctx context.Context, ea models.Street) (int, error)
func (esv *StreetService) Add(ctx context.Context, ea models.Street) (int, error) {
	var est ifStreetStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("StreetStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *StreetService) Upd(ctx context.Context, eu models.Street) (int, error)
func (esv *StreetService) Upd(ctx context.Context, eu models.Street) (int, error) {
	var est ifStreetStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("StreetStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *StreetService) Del(ctx context.Context, ed models.StreetClose) (models.Json_id, error)
func (esv *StreetService) Del(ctx context.Context, ed models.StreetClose) (models.Json_id, error) {
	var est ifStreetStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("StreetStorage.Del", err)
		return models.Json_id{}, err
	}

	return res, nil
}

//func (esv *StreetService) GetOne(ctx context.Context, i int) (models.Street_count, error)
func (esv *StreetService) GetOne(ctx context.Context, i int) (models.Street_count, error) {
	var est ifStreetStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("StreetStorage.GetOne", err)
		return models.Street_count{Values: []models.Street{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
