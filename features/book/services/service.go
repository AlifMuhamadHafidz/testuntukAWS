package services

import (
	"api/features/book"
	"api/helper"
	"errors"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
)

type bookSrv struct {
	data book.BookData
	vld  *validator.Validate
}

func New(d book.BookData) book.BookService {
	return &bookSrv{
		data: d,
		vld:  validator.New(),
	}
}

func (bs *bookSrv) Add(token interface{}, newBook book.Core) (book.Core, error) {
	userID := helper.ExtractToken(token)
	if userID <= 0 {
		return book.Core{}, errors.New("user not found")
	}

	err := bs.vld.Struct(newBook)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			log.Println(err)
		}
		return book.Core{}, errors.New("input buku tidak sesuai dengan arahan")
	}

	res, err := bs.data.Add(uint(userID), newBook)
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "buku not found"
		} else {
			msg = "terjadi kesalahan pada server"
		}
		return book.Core{}, errors.New(msg)
	}
	res.UserID = uint(userID)

	return res, nil

}
func (bs *bookSrv) Update(token interface{}, bookID uint, updatedData book.Core) (book.Core, error) {
	id := helper.ExtractToken(token)

	if id <= 0 {
		return book.Core{}, errors.New("data not found")
	}

	res, err := bs.data.Update(uint(id), bookID, updatedData)

	if err != nil {
		msg := ""

		if strings.Contains(err.Error(), "not found") {
			msg = "book not found"
		} else if strings.Contains(err.Error(), "unauthorized") {
			msg = "unauthorized request"
		} else {
			msg = "there is a problem with server"
		}
		return book.Core{}, errors.New(msg)
	}

	res.ID = bookID
	res.UserID = uint(id)

	return res, nil

}

func (bs *bookSrv) Delete(token interface{}, bookID uint) error {
	id := helper.ExtractToken(token)

	if id <= 0 {
		return errors.New("data not found")
	}

	err := bs.data.Delete(uint(id), bookID)

	if err != nil {
		log.Println("delete query error", err.Error())
		return err
	}

	return nil
}
