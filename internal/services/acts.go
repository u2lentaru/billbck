package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type ActService struct {
	storage pgsql.ActStorage
}

type ifActStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, gs3, ord int, dsc bool) (models.Act_count, error)
	Add(ctx context.Context, ea models.Act) (int, error)
	Upd(ctx context.Context, eu models.Act) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Act_count, error)
	Activate(ctx context.Context, i int, d string) (int, error)
}

//func NewActService(storage pgsql.ActStorage) *ActService
func NewActService(storage pgsql.ActStorage) *ActService {
	return &ActService{storage}
}

//func (esv *ActService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, gs3, ord int, dsc bool) (models.Act_count, error)
func (esv *ActService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, gs3, ord int, dsc bool) (models.Act_count, error) {
	var est ifActStorage
	est = &esv.storage
	// est = pgsql.NewActStorage(nil)

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	out_count, err := est.GetList(ctx, pg, pgs, gs1, gs2, gs3, ord, dsc)

	if err != nil {
		log.Println("ActStorage.GetList", err)
		return models.Act_count{Values: []models.Act{}, Count: 0, Auth: auth}, err
	}

	return out_count, nil
}

//func (esv *ActService) Add(ctx context.Context, ea models.Act) (int, error)
func (esv *ActService) Add(ctx context.Context, ea models.Act) (int, error) {
	var est ifActStorage
	est = &esv.storage
	// est = pgsql.NewActStorage(nil)

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("ActStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *ActService) Upd(ctx context.Context, eu models.Act) (int, error)
func (esv *ActService) Upd(ctx context.Context, eu models.Act) (int, error) {
	var est ifActStorage
	est = &esv.storage
	// est = pgsql.NewActStorage(nil)

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("ActStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *ActService) Del(ctx context.Context, ed []int) ([]int, error)
func (esv *ActService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifActStorage
	est = &esv.storage
	// est = pgsql.NewActStorage(nil)

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("ActStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *ActService) GetOne(ctx context.Context, i int) (models.Act_count, error)
func (esv *ActService) GetOne(ctx context.Context, i int) (models.Act_count, error) {
	var est ifActStorage
	est = &esv.storage
	// est = pgsql.NewActStorage(nil)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("ActStorage.GetOne", err)
		return models.Act_count{Values: []models.Act{}, Count: 0, Auth: auth}, err
	}

	return out_count, nil

}

//func (esv *ActService) Activate(ctx context.Context, i int, d string) (int, error)
func (esv *ActService) Activate(ctx context.Context, i int, d string) (int, error) {
	var est ifActStorage
	est = &esv.storage
	// est = pgsql.NewActStorage(nil)

	ai, err := est.Activate(ctx, i, d)

	if err != nil {
		log.Println("ActStorage.Activate", err)
		return 0, err
	}

	return ai, nil
}
