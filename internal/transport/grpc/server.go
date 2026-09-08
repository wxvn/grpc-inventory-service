package grpc

import (
	"context"
	"strconv"

	"github.com/wxvn/grpc-inventory-service/internal/domain"
	pb "github.com/wxvn/grpc-inventory-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	pb.UnimplementedInventoryServiceServer
	inventory Inventory
}

type Inventory interface {
	CreateProduct(ctx context.Context, p domain.Product) (domain.Product, error)
	GetProduct(ctx context.Context, id int) (domain.Product, error)
	GetProducts(ctx context.Context, page, page_size int32) ([]domain.Product, error)
}

func Register(gRPCServer *grpc.Server, inventory Inventory) {
	pb.RegisterInventoryServiceServer(
		gRPCServer,
		&serverAPI{
			inventory: inventory,
		},
	)
}

func (s *serverAPI) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	if req.Price <= 0 {
		return nil, status.Error(codes.InvalidArgument, "price must be greater than 0")
	}

	if req.Stock < 0 {
		return nil, status.Error(codes.InvalidArgument, "stock cannot be negative")
	}

	p := domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	product, err := s.inventory.CreateProduct(ctx, p)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed create product")
	}

	return &pb.CreateProductResponse{
		Product: &pb.Product{
			Id:          strconv.Itoa(product.ID),
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
		},
	}, nil
}

func (s *serverAPI) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	idStr, err := strconv.Atoi(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "id")
	}

	product, err := s.inventory.GetProduct(ctx, idStr)
	if err != nil {
		return nil, status.Error(codes.NotFound, "failed get product")
	}

	return &pb.GetProductResponse{
		Product: &pb.Product{
			Id:          strconv.Itoa(product.ID),
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
		},
	}, nil
}

func (s *serverAPI) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	products, err := s.inventory.GetProducts(ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, status.Error(codes.NotFound, "failed get products")
	}

	result := make([]*pb.Product, 0, len(products))

	for _, product := range products {
		result = append(result, &pb.Product{
			Id:          strconv.Itoa(product.ID),
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
		})
	}

	return &pb.ListProductsResponse{
		Products: result,
	}, nil
}
