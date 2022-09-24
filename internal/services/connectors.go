package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type ConnectorService struct {
	storage pgsql.ConnectorStorage
}

type ifConnectorStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Connector_count, error)
	Add(ctx context.Context, ea models.Connector) (int, error)
	Upd(ctx context.Context, eu models.Connector) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Connector_count, error)
}

//func NewConnectorService(storage pgsql.ConnectorStorage) *ConnectorService
func NewConnectorService(storage pgsql.ConnectorStorage) *ConnectorService {
	return &ConnectorService{storage}
}

//func (esv *ConnectorService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Connector_count, error)
func (esv *ConnectorService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Connector_count, error) {
	var est ifConnectorStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("ConnectorStorage.GetList", err)
		return models.Connector_count{Values: []models.Connector{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *ConnectorService) Add(ctx context.Context, ea models.Connector) (int, error)
func (esv *ConnectorService) Add(ctx context.Context, ea models.Connector) (int, error) {
	var est ifConnectorStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("ConnectorStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *ConnectorService) Upd(ctx context.Context, eu models.Connector) (int, error)
func (esv *ConnectorService) Upd(ctx context.Context, eu models.Connector) (int, error) {
	var est ifConnectorStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("ConnectorStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *ConnectorService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *ConnectorService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifConnectorStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("ConnectorStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *ConnectorService) GetOne(ctx context.Context, i int) (models.Connector_count, error)
func (esv *ConnectorService) GetOne(ctx context.Context, i int) (models.Connector_count, error) {
	var est ifConnectorStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("ConnectorStorage.GetOne", err)
		return models.Connector_count{Values: []models.Connector{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
