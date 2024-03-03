package controllers

import (
	"context"

	"github.com/buemura/event-driven-commerce/product-svc/internal/application/usecases"
	"github.com/buemura/event-driven-commerce/product-svc/internal/domain/product"
	"github.com/buemura/event-driven-commerce/product-svc/internal/infra/database"
	"github.com/buemura/event-driven-commerce/product-svc/internal/infra/grpc/server"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductController struct {
	server.UnimplementedProductServiceServer
}

func (c ProductController) GetManyProducts(
	ctx context.Context,
	in *server.GetManyProductsRequest,
) (*server.GetManyProductsResponse, error) {
	var page, items int = 1, 10
	if in.Page != 0 {
		page = int(in.Page)
	}
	if in.Items != 0 {
		items = int(in.Items)
	}

	repo := database.NewPgxProductRepository(database.Conn)
	usecase := usecases.NewGetManyProductUsecase(repo)

	res, err := usecase.Execute(&product.GetManyProductsIn{
		Page:  page,
		Items: items,
	})
	if err != nil {
		return nil, err
	}

	var productList []*server.ProductResponse
	for _, p := range res.ProductList {
		productList = append(productList, &server.ProductResponse{
			Id:       int32(p.ID),
			Name:     p.Name,
			Price:    int64(p.Price),
			Quantity: int32(p.Quantity),
			ImageUrl: p.ImageUrl,
		})
	}

	return &server.GetManyProductsResponse{
		ProductList: productList,
		Meta: &server.PaginationMeta{
			Page:       int32(res.Meta.Page),
			Items:      int32(res.Meta.Items),
			TotalPages: int32(res.Meta.TotalPages),
			TotalItems: int32(res.Meta.TotalItems),
		},
	}, nil
}

func (c ProductController) GetProduct(
	ctx context.Context,
	in *server.GetProductRequest,
) (*server.ProductResponse, error) {
	if in.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "missing id parameter")
	}

	repo := database.NewPgxProductRepository(database.Conn)
	usecase := usecases.NewGetProductUsecase(repo)

	prod, err := usecase.Execute(int(in.Id))
	if err != nil {
		return nil, err
	}

	return &server.ProductResponse{
		Id:       int32(prod.ID),
		Name:     prod.Name,
		Price:    int64(prod.Price),
		Quantity: int32(prod.Quantity),
		ImageUrl: prod.ImageUrl,
	}, nil
}
