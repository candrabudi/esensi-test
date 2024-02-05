package factory

import (
	"esensi-test/database"
	"esensi-test/internal/repository"
)

type Factory struct {
	UserRepository        repository.User
	UserSessionRepository repository.UserSession
	ItemRepository        repository.Item
	CustomerRepository    repository.Customer
	InvoiceRepository     repository.Invoice
	RedisRepository       repository.RedisRepository
}

func NewFactory() *Factory {
	// Check db connection
	db := database.GetConnection()
	rdb := database.GetRedisClient()
	return &Factory{
		// Pass the db connection to repository package for database query calling
		repository.NewUserRepository(db),
		repository.NewUserSessionRepository(db),
		repository.NewItemRepository(db),
		repository.NewCustomerRepository(db),
		repository.NewInvoiceRepository(db),
		repository.NewRedisRepository(rdb),
	}
}
