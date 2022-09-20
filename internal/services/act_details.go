package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type ActDetailService struct {
	storage pgsql.ActDetailStorage
}

type ifActDetailStorage interface {
	GetList(ctx context.Context, pg, pgs, nm, ord int, dsc bool) (models.ActDetail_count, error)
	Add(ctx context.Context, ea models.ActDetail) (int, error)
	Upd(ctx context.Context, eu models.ActDetail) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.ActDetail_count, error)
}

//NewActDetailService(storage pg.ActTypeStorage) *ActTypeService
func NewActDetailService(storage pgsql.ActDetailStorage) *ActDetailService {
	return &ActDetailService{storage}
}

//func (esv *ActDetailService) GetList(ctx context.Context, pg, pgs, nm, ord int, dsc bool) (models.ActDetail_count, error)
func (esv *ActDetailService) GetList(ctx context.Context, pg, pgs, nm, ord int, dsc bool) (models.ActDetail_count, error) {
	var est ifActDetailStorage
	est = &esv.storage
	// est = pgsql.NewActDetailStorage(nil)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	out_count, err := est.GetList(ctx, pg, pgs, nm, ord, dsc)

	if err != nil {
		log.Println("ActDetailStorage.GetList", err)
		return models.ActDetail_count{Values: []models.ActDetail{}, Count: 0, Auth: auth}, err
	}

	return out_count, nil
}

//func (esv *ActDetailService) Add(ctx context.Context, ea models.ActDetail) (int, error)
func (esv *ActDetailService) Add(ctx context.Context, ea models.ActDetail) (int, error) {
	var est ifActDetailStorage
	est = &esv.storage
	// est = pgsql.NewActDetailStorage(nil)

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("ActDetailStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *ActDetailService) Upd(ctx context.Context, eu models.ActDetail) (int, error)
func (esv *ActDetailService) Upd(ctx context.Context, eu models.ActDetail) (int, error) {
	var est ifActDetailStorage
	est = &esv.storage
	// est = pgsql.NewActDetailStorage(nil)

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("ActDetailStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *ActDetailService) Del(ctx context.Context, ed []int) ([]int, error)
func (esv *ActDetailService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifActDetailStorage
	est = &esv.storage
	// est = pgsql.NewActDetailStorage(nil)

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("ActDetailStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *ActDetailService) GetOne(ctx context.Context, i int) (models.ActDetail_count, error)
func (esv *ActDetailService) GetOne(ctx context.Context, i int) (models.ActDetail_count, error) {
	var est ifActDetailStorage
	est = &esv.storage
	// est = pgsql.NewActDetailStorage(nil)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("ActDetailStorage.GetOne", err)
		return models.ActDetail_count{Values: []models.ActDetail{}, Count: 0, Auth: auth}, err
	}

	return out_count, nil

}
