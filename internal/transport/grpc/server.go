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
	UpdateStockProduct(ctx context.Context, id int, quantity int32) (domain.Product, error)
	DeleteProduct(ctx context.Context, id int) error
	CheckAvailability(ctx context.Context, id int, quantity int32) (bool, int32, error)
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

func (s *serverAPI) UpdateStock(ctx context.Context, req *pb.UpdateStockRequest) (*pb.UpdateStockResponse, error) {
	idStr, err := strconv.Atoi(req.ProductId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "id")
	}

	updateProduct, err := s.inventory.UpdateStockProduct(ctx, idStr, req.Quantity)
	if err != nil {
		return nil, status.Error(codes.NotFound, "failed update products")
	}

	return &pb.UpdateStockResponse{
		Product: &pb.Product{
			Id:          strconv.Itoa(updateProduct.ID),
			Name:        updateProduct.Name,
			Description: updateProduct.Description,
			Price:       updateProduct.Price,
			Stock:       updateProduct.Stock,
		},
	}, nil
}

func (s *serverAPI) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	id, err := strconv.Atoi(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	err = s.inventory.DeleteProduct(ctx, id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "failed to delete product")
	}

	return &pb.DeleteProductResponse{
		Success: true,
	}, nil
}

func (s *serverAPI) CheckAvailability(ctx context.Context, req *pb.CheckAvailabilityRequest) (*pb.CheckAvailabilityResponse, error) {
	if req.Quantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "quantity must be greater than zero")
	}

	id, err := strconv.Atoi(req.ProductId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	available, currentStock, err := s.inventory.CheckAvailability(ctx, id, req.Quantity)
	if err != nil {
		return nil, status.Error(codes.NotFound, "failed to check availability")
	}

	return &pb.CheckAvailabilityResponse{
		Available:    available,
		CurrentStock: currentStock,
	}, nil
}
