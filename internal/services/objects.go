package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type ObjectService struct {
	storage pgsql.ObjectStorage
}

type ifObjectStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, gs3, gs3f bool, ord int, dsc bool) (models.Object_count, error)
	Add(ctx context.Context, ea models.Object) (int, error)
	Upd(ctx context.Context, eu models.Object) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Object_count, error)
	GetObjContract(ctx context.Context, i int, a string) (models.ObjContract_count, error)
	GetMff(ctx context.Context, i int) (models.Object_count, error)
}

//NewObjectService(storage pg.ObjectStorage) *ObjectService
func NewObjectService(storage pgsql.ObjectStorage) *ObjectService {
	return &ObjectService{storage}
}

//func (esv *ObjectService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, gs3, gs3f bool, ord int, dsc bool) (models.Object_count, error)
func (esv *ObjectService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, gs3, gs3f bool, ord int, dsc bool) (models.Object_count, error) {
	var est ifObjectStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, gs2, gs3, gs3f, ord, dsc)

	if err != nil {
		log.Println("ObjectStorage.GetList", err)
		return models.Object_count{Values: []models.Object{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *ObjectService) Add(ctx context.Context, ea models.Object) (int, error)
func (esv *ObjectService) Add(ctx context.Context, ea models.Object) (int, error) {
	var est ifObjectStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("ObjectStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *ObjectService) Upd(ctx context.Context, eu models.Object) (int, error)
func (esv *ObjectService) Upd(ctx context.Context, eu models.Object) (int, error) {
	var est ifObjectStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("ObjectStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *ObjectService) Del(ctx context.Context, ed []int) ([]int, error)
func (esv *ObjectService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifObjectStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("ObjectStorage.Del", err)
		return nil, err
	}

	return res, nil
}

//func (esv *ObjectService) GetOne(ctx context.Context, i int) (models.Object_count, error)
func (esv *ObjectService) GetOne(ctx context.Context, i int) (models.Object_count, error) {
	var est ifObjectStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("ObjectStorage.GetOne", err)
		return models.Object_count{Values: []models.Object{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *ObjectService) GetObjContract(ctx context.Context, i int, s string) (models.ObjContract_count, error)
func (esv *ObjectService) GetObjContract(ctx context.Context, i int, s string) (models.ObjContract_count, error) {
	var est ifObjectStorage
	est = &esv.storage

	out_count, err := est.GetObjContract(ctx, i, s)

	if err != nil {
		log.Println("ObjectStorage.GetObjContract", err)
		return models.ObjContract_count{}, err
	}

	return out_count, nil
}

//func (esv *ObjectService) GetMff(ctx context.Context, i int) (models.Object_count, error)
func (esv *ObjectService) GetMff(ctx context.Context, i int) (models.Object_count, error) {
	var est ifObjectStorage
	est = &esv.storage

	out_count, err := est.GetMff(ctx, i)

	if err != nil {
		log.Println("ObjectStorage.GetMff", err)
		return models.Object_count{}, err
	}

	return out_count, nil
}
