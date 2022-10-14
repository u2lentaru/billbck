package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type SubjectService struct {
	storage pgsql.SubjectStorage
}

type ifSubjectStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2, gs3 string, ord int, dsc bool) (models.Subject_count, error)
	Add(ctx context.Context, ea models.Subject) (int, error)
	Upd(ctx context.Context, eu models.Subject) (int, error)
	Del(ctx context.Context, ed models.SubjectClose) (int, error)
	GetOne(ctx context.Context, i int) (models.Subject_count, error)
	GetHist(ctx context.Context, i int) ([]string, error)
}

//NewSubjectService(storage pg.SubjectStorage) *SubjectService
func NewSubjectService(storage pgsql.SubjectStorage) *SubjectService {
	return &SubjectService{storage}
}

//func (esv *SubjectService) GetList(ctx context.Context, pg, pgs int, gs1, gs2, gs3 string, ord int, dsc bool) (models.Subject_count, error)
func (esv *SubjectService) GetList(ctx context.Context, pg, pgs int, gs1, gs2, gs3 string, ord int, dsc bool) (models.Subject_count, error) {
	var est ifSubjectStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, gs2, gs3, ord, dsc)

	if err != nil {
		log.Println("SubjectStorage.GetList", err)
		return models.Subject_count{Values: []models.Subject{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *SubjectService) Add(ctx context.Context, ea models.Subject) (int, error)
func (esv *SubjectService) Add(ctx context.Context, ea models.Subject) (int, error) {
	var est ifSubjectStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("SubjectStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *SubjectService) Upd(ctx context.Context, eu models.Subject) (int, error)
func (esv *SubjectService) Upd(ctx context.Context, eu models.Subject) (int, error) {
	var est ifSubjectStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("SubjectStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *SubjectService) Del(ctx context.Context, ed models.SubjectClose) (int, error)
func (esv *SubjectService) Del(ctx context.Context, ed models.SubjectClose) (int, error) {
	var est ifSubjectStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("SubjectStorage.Del", err)
		return 0, err
	}

	return res, nil
}

//func (esv *SubjectService) GetOne(ctx context.Context, i int) (models.Subject_count, error)
func (esv *SubjectService) GetOne(ctx context.Context, i int) (models.Subject_count, error) {
	var est ifSubjectStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("SubjectStorage.GetOne", err)
		return models.Subject_count{Values: []models.Subject{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *SubjectService) GetHist(ctx context.Context, i int) ([]string, error)
func (esv *SubjectService) GetHist(ctx context.Context, i int) ([]string, error) {
	var est ifSubjectStorage
	est = &esv.storage

	out_count, err := est.GetHist(ctx, i)

	if err != nil {
		log.Println("SubjectStorage.GetOne", err)
		return []string{}, err
	}

	return out_count, nil
}
