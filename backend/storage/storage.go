package storage

import "github.com/declanl482/go-book-tracker-app/backend/types"

type Storage interface {
	CreateUser(user *types.User) (*types.User, error)
	GetUserByEmail(email string) (*types.User, error)
	GetUser(id int) (*types.User, error)
	UpdateUser(user *types.User) (*types.User, error)
	IsEmailTaken(email string) (bool, error)
	DeleteUser(user *types.User) error

	CreateBook(book *types.Book) (*types.Book, error)
	GetBooks(id int) (*[]types.Book, error)
	GetBook(id int) (*types.Book, error)
	UpdateBook(book *types.Book) (*types.Book, error)
	DeleteBook(book *types.Book) error
}
