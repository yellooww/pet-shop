package product

import (
	"bytes"
	"context"
	"errors"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
)

// Get Product - Ready
func TestGetAllProducts_Success(t *testing.T) {
	// Мокаем storage — он вернёт один продукт.
	mock := &ProductsMock{
		GetAllProductsFunc: func(ctx context.Context) ([]models.Product, error) {
			return []models.Product{
				{ID: 1, Name: "Dog Food"},
			}, nil
		},
	}

	// Создаем HTTP-запрос GET /products
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), mock)

	// Вызываем метод GetAllProducts, который является http.HandlerFunc
	handler.GetAllProducts(w, req)

	// Проверяем HTTP-код
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Dog Food") {
		t.Errorf("unexpected response body: %s", w.Body.String())
	}
}
func TestGetAllProducts_Error(t *testing.T) {
	// Мокаем storage — он будет возвращать ошибку
	mock := &ProductsMock{
		GetAllProductsFunc: func(ctx context.Context) ([]models.Product, error) {
			return nil, errors.New("DB error")
		},
	}

	// Создаем запрос
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.GetAllProducts(w, req)

	// Ожидаем HTTP 500
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Failed to retrieve products") {
		t.Errorf("unexpected response body: %s", w.Body.String())
	}
}

// =======================
// Create Product
// =======================

func TestCreateProduct_Success(t *testing.T) {
	mock := &ProductsMock{
		CreateProductFunc: func(ctx context.Context, product models.Product) (int, error) {
			return 1, nil
		},
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/products",
		bytes.NewBufferString(`{"name": "Dog Food", "price": 10.5, "quantity": 100}`),
	)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.CreateProduct(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Product created successfully") {
		t.Errorf("unexpected response body: %s", w.Body.String())
	}
}

func TestCreateProduct_BadRequest(t *testing.T) {
	mock := &ProductsMock{
		CreateProductFunc: func(ctx context.Context, product models.Product) (int, error) {
			return 1, nil
		},
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/products",
		bytes.NewBufferString(`{"name":}`),
	)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.CreateProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Invalid JSON payload") {
		t.Errorf("unexpected response body: %s", w.Body.String())
	}
}

func TestCreateProduct_Fail(t *testing.T) {
	mock := &ProductsMock{
		CreateProductFunc: func(ctx context.Context, product models.Product) (int, error) {
			return 0, errors.New("DB error")
		},
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/products",
		bytes.NewBufferString(`{"name": "Dog Food", "price": 10.5, "quantity": 100}`),
	)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.CreateProduct(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Failed to create product") {
		t.Errorf("unexpected response body: %s", w.Body.String())
	}
}

// =======================
// Update Product
// =======================

func TestUpdateProduct_Success(t *testing.T) {
	mock := &ProductsMock{
		UpdateProductFunc: func(ctx context.Context, product models.Product) error {
			return nil
		},
	}

	req := httptest.NewRequest(
		http.MethodPut,
		"/products/1",
		bytes.NewBufferString(`{"id": 1, "name": "Dog Food", "price": 10.5, "quantity": 100}`),
	)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)

	router := chi.NewRouter()
	router.Put("/products/{id}", handler.UpdateProduct)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Product updated successfully") {
		t.Errorf("unexpected response body: %s", w.Body.String())
	}
}

func TestUpdateProduct_BadRequest(t *testing.T) {
	mock := &ProductsMock{
		UpdateProductFunc: func(ctx context.Context, product models.Product) error {
			return nil
		},
	}

	req := httptest.NewRequest(
		http.MethodPut,
		"/products/1",
		bytes.NewBufferString(`{"id": 1, "name":}`),
	)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)

	router := chi.NewRouter()
	router.Put("/products/{id}", handler.UpdateProduct)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Invalid JSON payload") {
		t.Errorf("unexpected response body: %s", w.Body.String())
	}
}

func TestUpdateProduct_Fail(t *testing.T) {
	mock := &ProductsMock{
		UpdateProductFunc: func(ctx context.Context, product models.Product) error {
			return errors.New("DB error")
		},
	}

	req := httptest.NewRequest(
		http.MethodPut,
		"/products/1",
		bytes.NewBufferString(`{"id": 1, "name": "Dog Food", "price": 10.5, "quantity": 100}`),
	)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)

	router := chi.NewRouter()
	router.Put("/products/{id}", handler.UpdateProduct)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Failed to update product") {
		t.Errorf("unexpected response body: %s", w.Body.String())
	}
}

// =======================
// Delete Product
// =======================

func TestDeleteProduct_Success(t *testing.T) {
	mock := &ProductsMock{
		DeleteProductFunc: func(ctx context.Context, id int) error {
			return nil
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)

	router := chi.NewRouter()
	router.Delete("/products/{id}", handler.DeleteProduct)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Product deleted successfully") {
		t.Errorf("unexpected response body: %s", w.Body.String())
	}
}

func TestDeleteProduct_BadRequest(t *testing.T) {
	mock := &ProductsMock{
		DeleteProductFunc: func(ctx context.Context, id int) error {
			return nil
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/products", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.DeleteProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Product ID is required") {
		t.Errorf("unexpected response body: %s", w.Body.String())
	}
}

func TestDeleteProduct_Fail(t *testing.T) {
	mock := &ProductsMock{
		DeleteProductFunc: func(ctx context.Context, id int) error {
			return errors.New("DB error")
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)

	router := chi.NewRouter()
	router.Delete("/products/{id}", handler.DeleteProduct)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Failed to delete product") {
		t.Errorf("unexpected response body: %s", w.Body.String())
	}
}
