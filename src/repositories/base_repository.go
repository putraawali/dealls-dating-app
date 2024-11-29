package repositories

import "github.com/sarulabs/di"

type Repositories struct {
	User        UserRepository
	Swipe       SwipeRepository
	Transaction TransactionRepository
}

// Initiate repository layer, accept dependency injection as parameter and return *Repositories
func NewRepository(di di.Container) *Repositories {
	return &Repositories{
		User:        NewUserRepository(di),
		Swipe:       NewSwipeRepository(di),
		Transaction: NewTransactionRepository(di),
	}
}
