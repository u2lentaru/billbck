package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type EquipmentService struct {
	storage pgsql.EquipmentStorage
}

type ifEquipmentStorage interface {
	GetList(ctx context.Context, pg, pgs, gs1 int, gs2 string, ord int, dsc bool) (models.Equipment_count, error)
	Add(ctx context.Context, ea models.Equipment) (int, error)
	Upd(ctx context.Context, eu models.Equipment) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Equipment_count, error)
	AddList(ctx context.Context, al models.Equipment_count) ([]int, error)
	DelObj(ctx context.Context, d []int) ([]int, error)
}

//NewEquipmentService(storage pg.EquipmentStorage) *EquipmentService
func NewEquipmentService(storage pgsql.EquipmentStorage) *EquipmentService {
	return &EquipmentService{storage}
}

//func (esv *EquipmentService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Equipment_count, error)
func (esv *EquipmentService) GetList(ctx context.Context, pg, pgs, gs1 int, gs2 string, ord int, dsc bool) (models.Equipment_count, error) {
	var est ifEquipmentStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, gs2, ord, dsc)

	if err != nil {
		log.Println("EquipmentStorage.GetList", err)
		return models.Equipment_count{Values: []models.Equipment{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *EquipmentService) Add(ctx context.Context, ea models.Equipment) (int, error)
func (esv *EquipmentService) Add(ctx context.Context, ea models.Equipment) (int, error) {
	var est ifEquipmentStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("EquipmentStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *EquipmentService) Upd(ctx context.Context, eu models.Equipment) (int, error)
func (esv *EquipmentService) Upd(ctx context.Context, eu models.Equipment) (int, error) {
	var est ifEquipmentStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("EquipmentStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *EquipmentService) Del(ctx context.Context, ed []int) ([]int, error)
func (esv *EquipmentService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifEquipmentStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("EquipmentStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *EquipmentService) GetOne(ctx context.Context, i int) (models.Equipment_count, error)
func (esv *EquipmentService) GetOne(ctx context.Context, i int) (models.Equipment_count, error) {
	var est ifEquipmentStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("EquipmentStorage.GetOne", err)
		return models.Equipment_count{Values: []models.Equipment{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *EquipmentService) AddList(ctx context.Context, al models.Equipment_count) ([]int, error)
func (esv *EquipmentService) AddList(ctx context.Context, al models.Equipment_count) ([]int, error) {
	var est ifEquipmentStorage
	est = &esv.storage

	ai, err := est.AddList(ctx, al)

	if err != nil {
		log.Println("EquipmentStorage.Add", err)
		return []int{}, err
	}

	return ai, nil
}

//func (esv *EquipmentService) DelObj(ctx context.Context, d []int) ([]int, error)
func (esv *EquipmentService) DelObj(ctx context.Context, d []int) ([]int, error) {
	var est ifEquipmentStorage
	est = &esv.storage

	res, err := est.DelObj(ctx, d)

	if err != nil {
		log.Println("EquipmentStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}
