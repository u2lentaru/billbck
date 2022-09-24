package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type ConclusionService struct {
	storage pgsql.ConclusionStorage
}

type ifConclusionStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Conclusion_count, error)
	Add(ctx context.Context, ea models.Conclusion) (int, error)
	Upd(ctx context.Context, eu models.Conclusion) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Conclusion_count, error)
}

//func NewConclusionService(storage pgsql.ConclusionStorage) *ConclusionService
func NewConclusionService(storage pgsql.ConclusionStorage) *ConclusionService {
	return &ConclusionService{storage}
}

//func (esv *ConclusionService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Conclusion_count, error)
func (esv *ConclusionService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Conclusion_count, error) {
	var est ifConclusionStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("ConclusionStorage.GetList", err)
		return models.Conclusion_count{Values: []models.Conclusion{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *ConclusionService) Add(ctx context.Context, ea models.Conclusion) (int, error)
func (esv *ConclusionService) Add(ctx context.Context, ea models.Conclusion) (int, error) {
	var est ifConclusionStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("ConclusionStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *ConclusionService) Upd(ctx context.Context, eu models.Conclusion) (int, error)
func (esv *ConclusionService) Upd(ctx context.Context, eu models.Conclusion) (int, error) {
	var est ifConclusionStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("ConclusionStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *ConclusionService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *ConclusionService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifConclusionStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("ConclusionStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *ConclusionService) GetOne(ctx context.Context, i int) (models.Conclusion_count, error)
func (esv *ConclusionService) GetOne(ctx context.Context, i int) (models.Conclusion_count, error) {
	var est ifConclusionStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("ConclusionStorage.GetOne", err)
		return models.Conclusion_count{Values: []models.Conclusion{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
