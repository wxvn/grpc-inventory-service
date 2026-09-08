package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wxvn/grpc-inventory-service/internal/domain"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, dsn string) (*Repository, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &Repository{pool: pool}, nil
}

func (repository *Repository) CreateProduct(ctx context.Context, p domain.Product) (domain.Product, error) {
	query := `
	INSERT INTO product (name, description, price, stock)
	VALUES ($1, $2, $3, $4)
	RETURNING *;
	`

	var product domain.Product

	err := repository.pool.QueryRow(ctx, query, p.Name, p.Description, p.Price, p.Stock).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
	)
	if err != nil {
		return domain.Product{}, err
	}

	return product, nil
}

func (repository *Repository) GetProductByID(ctx context.Context, id int) (domain.Product, error) {
	query := `
	SELECT *
	FROM product
	WHERE id=$1;
	`

	var product domain.Product

	err := repository.pool.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
	)
	if err != nil {
		return domain.Product{}, err
	}

	return product, nil
}

func (repository *Repository) GetProducts(ctx context.Context, limit, offset int) ([]domain.Product, error) {
	query := `
	SELECT *
	FROM product
	ORDER By id
	LIMIT $1 OFFSET $2;
	`

	rows, err := repository.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []domain.Product

	for rows.Next() {
		var product domain.Product

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil

}
