package database

import (
	"context"

	"github.com/buemura/event-driven-commerce/product-svc/internal/domain/product"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxProductRepository struct {
	conn *pgxpool.Pool
}

func NewPgxProductRepository(conn *pgxpool.Pool) *PgxProductRepository {
	return &PgxProductRepository{
		conn: conn,
	}
}

func (r *PgxProductRepository) FindById(id int) (*product.Product, error) {
	rows, err := r.conn.Query(context.Background(), `SELECT * FROM product WHERE id = $1`, id)
	p, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[product.Product])
	if err != nil {
		return nil, err
	}
	if len(p) == 0 {
		return nil, nil
	}
	return p[0], nil
}

func (r *PgxProductRepository) FindMany(in *product.GetManyProductsIn) (*product.ProductRepositoryPaginatedOut, error) {
	limit := in.Items
	offset := (in.Page - 1) * in.Items

	rows, err := r.conn.Query(context.Background(), `SELECT * FROM product LIMIT $1 OFFSET $2`, limit, offset)
	p, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[product.Product])
	if err != nil {
		return nil, err
	}

	var totalCount int
	err = r.conn.QueryRow(context.Background(), `SELECT count(id) as total_count FROM product`).Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	return &product.ProductRepositoryPaginatedOut{
		ProductList: p,
		TotalCount:  totalCount,
	}, nil
}

func (r *PgxProductRepository) Update(newP *product.Product) (*product.Product, error) {
	res, err := r.conn.Exec(
		context.Background(),
		`
		UPDATE product 
		SET name = $1, description = $2, price = $3, quantity = $4, image_url = $5 
		WHERE id = $6
		`,
		newP.Name, newP.Description, newP.Price, newP.Quantity, newP.ImageUrl,
		newP.ID,
	)
	if err != nil {
		return nil, err
	}

	if res.RowsAffected() == 0 {
		return nil, product.ErrProductNotFound
	}
	return newP, nil
}
