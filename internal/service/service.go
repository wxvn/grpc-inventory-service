package service

import (
	"context"
	"log/slog"

	"github.com/wxvn/grpc-inventory-service/internal/domain"
)

type Service struct {
	repository Repository
	log        *slog.Logger
}

type Repository interface {
	CreateProduct(ctx context.Context, p domain.Product) (domain.Product, error)
	GetProductByID(ctx context.Context, id int) (domain.Product, error)
	GetProducts(ctx context.Context, limit, offset int) ([]domain.Product, error)
}

func New(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repository: repo,
		log:        logger,
	}
}

func (s *Service) CreateProduct(ctx context.Context, p domain.Product) (domain.Product, error) {
	product, err := s.repository.CreateProduct(ctx, p)
	if err != nil {
		s.log.Error(
			"failed to create product",
			slog.String("error", err.Error()),
		)
		return domain.Product{}, err
	}

	return product, nil
}

func (s *Service) GetProduct(ctx context.Context, id int) (domain.Product, error) {
	product, err := s.repository.GetProductByID(ctx, id)
	if err != nil {
		s.log.Error(
			"failed to get product",
			slog.Int("id", id),
			slog.String("error", err.Error()),
		)
		return domain.Product{}, err
	}

	return product, nil
}

func (s *Service) GetProducts(ctx context.Context, page, pageSize int32) ([]domain.Product, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	limit := pageSize
	offset := (page - 1) * pageSize

	products, err := s.repository.GetProducts(ctx, int(limit), int(offset))
	if err != nil {
		s.log.Error(
			"failed to get products",
			slog.String("error", err.Error()),
		)
		return []domain.Product{}, err
	}

	return products, nil

}
