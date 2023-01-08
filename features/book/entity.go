package book

import "github.com/labstack/echo/v4"

type Core struct {
	ID          uint
	Judul       string `validate:"required"`
	TahunTerbit int    `validate:"required"`
	Penulis     string `validate:"required"`
	UserID      uint
	Pemilik     string
}

type BookHandler interface {
	Add() echo.HandlerFunc
	Update() echo.HandlerFunc
	// Delete() echo.HandlerFunc
	// MyBook() echo.HandlerFunc
}

type BookService interface {
	Add(token interface{}, newBook Core) (Core, error)
	Update(token interface{}, bookID uint, updatedData Core) (Core, error)
	// Delete(token interface{}, bookID int) error
	// MyBook(token interface{}) ([]Core, error)
}

type BookData interface {
	Add(userID uint, newBook Core) (Core, error)
	Update(userID uint, bookID uint, updatedData Core) (Core, error)
	Delete(userID uint, bookID uint) error
	// MyBook(userID int) ([]Core, error)
}
