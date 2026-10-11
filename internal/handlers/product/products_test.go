package product

import (
	"bytes"
	"errors"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/mock"

	"go-pet-shop/internal/handlers/product/mocks"
)

// Get Product - Ready
func TestGetAllProducts_Success(t *testing.T) {
	// Мокаем storage — он вернёт один продукт.
	productsMock := mocks.NewProducts(t)
	productsMock.On("GetAllProducts", mock.Anything).Return([]models.Product{
		{ID: 1, Name: "Dog Food"},
	}, nil)

	// Создаем HTTP-запрос GET /products
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), productsMock)

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
	productsMock := mocks.NewProducts(t)
	productsMock.On("GetAllProducts", mock.Anything).Return(nil, errors.New("DB error"))

	// Создаем запрос
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), productsMock)
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
	productsMock := mocks.NewProducts(t)
	productsMock.On("CreateProduct", mock.Anything, mock.Anything).Return(1, nil)

	req := httptest.NewRequest(
		http.MethodPost,
		"/products",
		bytes.NewBufferString(`{"name": "Dog Food", "price": 10.5, "stock": 100}`),
	)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), productsMock)
	handler.CreateProduct(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Product created successfully") {
		t.Errorf("unexpected response body: %s", w.Body.String())
	}
}

func TestCreateProduct_BadRequest(t *testing.T) {
	productsMock := mocks.NewProducts(t)

	req := httptest.NewRequest(
		http.MethodPost,
		"/products",
		bytes.NewBufferString(`{"name":}`),
	)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), productsMock)
	handler.CreateProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Invalid JSON payload") {
		t.Errorf("unexpected response body: %s", w.Body.String())
	}
}

func TestCreateProduct_Fail(t *testing.T) {
	productsMock := mocks.NewProducts(t)
	productsMock.On("CreateProduct", mock.Anything, mock.Anything).Return(0, errors.New("DB error"))

	req := httptest.NewRequest(
		http.MethodPost,
		"/products",
		bytes.NewBufferString(`{"name": "Dog Food", "price": 10.5, "stock": 100}`),
	)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), productsMock)
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
	productsMock := mocks.NewProducts(t)
	productsMock.On("UpdateProduct", mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(
		http.MethodPut,
		"/products/1",
		bytes.NewBufferString(`{"id": 1, "name": "Dog Food", "price": 10.5, "stock": 100}`),
	)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), productsMock)

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
	productsMock := mocks.NewProducts(t)

	req := httptest.NewRequest(
		http.MethodPut,
		"/products/1",
		bytes.NewBufferString(`{"id": 1, "name":}`),
	)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), productsMock)

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
	productsMock := mocks.NewProducts(t)
	productsMock.On("UpdateProduct", mock.Anything, mock.Anything).Return(errors.New("DB error"))

	req := httptest.NewRequest(
		http.MethodPut,
		"/products/1",
		bytes.NewBufferString(`{"id": 1, "name": "Dog Food", "price": 10.5, "stock": 100}`),
	)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), productsMock)

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
	productsMock := mocks.NewProducts(t)
	productsMock.On("DeleteProduct", mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), productsMock)

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
	productsMock := mocks.NewProducts(t)

	req := httptest.NewRequest(http.MethodDelete, "/products", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), productsMock)
	handler.DeleteProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Product ID is required") {
		t.Errorf("unexpected response body: %s", w.Body.String())
	}
}

func TestDeleteProduct_Fail(t *testing.T) {
	productsMock := mocks.NewProducts(t)
	productsMock.On("DeleteProduct", mock.Anything, mock.Anything).Return(errors.New("DB error"))

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), productsMock)

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
