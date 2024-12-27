package main

import (
	"advancedGo/configs"
	"advancedGo/internal/order"
	"advancedGo/internal/product"
	"advancedGo/internal/user"
	"advancedGo/pkg/jwt"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Структура для тестового окружения
type testEnv struct {
	db *gorm.DB
}

// Инициализация тестовой БД
func setupTestDB(t *testing.T) *testEnv {
	err := godotenv.Load("../.env.test")
	if err != nil {
		t.Fatal("Error loading .env file:", err)
	}

	// Используем отдельную тестовую БД
	testDSN := os.Getenv("TEST_DSN")
	db, err := gorm.Open(postgres.Open(testDSN), &gorm.Config{})
	if err != nil {
		t.Fatal("Failed to connect to test database:", err)
	}

	// Миграция схемы для тестовой БД
	err = db.AutoMigrate(&product.Product{}, &user.User{}, &order.Order{})
	if err != nil {
		t.Fatal("Failed to migrate test database:", err)
	}

	return &testEnv{
		db: db,
	}
}

// Подготовка тестовых данных
func (env *testEnv) prepareTestData(t *testing.T) (*user.User, []product.Product) {
	// Создаем тестового пользователя
	testUser := &user.User{
		Phone:            "+79999999999",
		SessionId:        "test_session",
		VerificationCode: "1234",
		IsVerified:       true,
	}
	if err := env.db.Create(testUser).Error; err != nil {
		t.Fatal("Failed to create test user:", err)
	}

	// Создаем тестовые продукты
	testProducts := []product.Product{
		{
			Name:        "СУП",
			Description: "Борщ с мясом",
		},
		{
			Name:        "Пицца",
			Description: "Пеперони",
		},
	}

	for i := range testProducts {
		if err := env.db.Create(&testProducts[i]).Error; err != nil {
			t.Fatal("Failed to create test product:", err)
		}
	}

	return testUser, testProducts
}

func TestCreateOrder(t *testing.T) {
	// Подготовка тестового окружения
	env := setupTestDB(t)

	// Очистка данных после теста
	defer env.cleanupTestData(t)

	// Подготовка тестовых данных
	testUser, testProducts := env.prepareTestData(t)

	// Создаем конфиг для тестового сервера
	config := &configs.Config{
		DB: configs.DBConfig{
			DSN: os.Getenv("TEST_DSN"),
		},
		Auth: configs.AuthConfig{
			Token: os.Getenv("TEST_TOKEN"),
		},
	}

	// Создаем тестовый сервер
	ts := httptest.NewServer(App(config))
	defer ts.Close()

	// Подготовка данных для запроса создания заказа
	orderRequest := order.OrderCreateRequest{
		ProductIds: []uint{testProducts[0].ID, testProducts[1].ID},
	}
	requestBody, err := json.Marshal(orderRequest)
	if err != nil {
		t.Fatal("Failed to marshal request:", err)
	}

	// Создаем HTTP запрос
	req, err := http.NewRequest("POST", ts.URL+"/order", bytes.NewReader(requestBody))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}

	// Добавляем токен авторизации
	req.Header.Set("Authorization", "Bearer "+createTestToken(testUser))
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("Failed to send request:", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	// Проверяем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Failed to read response body:", err)
	}

	var createdOrder order.Order
	if err := json.Unmarshal(body, &createdOrder); err != nil {
		t.Fatal("Failed to unmarshal response:", err)
	}

	// Проверяем создание заказа в БД
	var dbOrder order.Order
	if err := env.db.Preload("Products").First(&dbOrder, createdOrder.ID).Error; err != nil {
		t.Fatal("Failed to fetch created order from DB:", err)
	}

	// Проверяем количество продуктов в заказе
	if len(dbOrder.Products) != 2 {
		t.Errorf("Expected 2 products in order, got %d", len(dbOrder.Products))
	}
}

// Создаем тестовый токен
func createTestToken(user *user.User) string {
	token, err := jwt.NewJWT(os.Getenv("TEST_TOKEN")).Create(jwt.JWTData{
		Phone: user.Phone,
	})
	if err != nil {
		panic(err)
	}
	return token
}

// Очистка тестовых данных
func (env *testEnv) cleanupTestData(t *testing.T) {
	tables := []string{"order_products", "orders", "products", "users"}
	for _, table := range tables {
		if err := env.db.Exec("DELETE FROM " + table).Error; err != nil {
			t.Errorf("Failed to cleanup table %s: %v", table, err)
		}
	}
}
