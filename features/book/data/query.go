package data

import (
	"api/features/book"
	"errors"
	"log"

	"gorm.io/gorm"
)

type bookData struct {
	db *gorm.DB
}

func New(db *gorm.DB) book.BookData {
	return &bookData{
		db: db,
	}
}

func (bd *bookData) Add(userID uint, newBook book.Core) (book.Core, error) {
	cnv := CoreToData(newBook)
	cnv.UserID = uint(userID)
	err := bd.db.Create(&cnv).Error
	if err != nil {
		return book.Core{}, err
	}

	newBook.ID = cnv.ID

	return newBook, nil
}
func (bd *bookData) Update(userID uint, bookID uint, updatedData book.Core) (book.Core, error) {
	getID := Books{}
	err := bd.db.Where("id = ?", bookID).First(&getID).Error

	if err != nil {
		log.Println("get user book error", err.Error())
		return book.Core{}, err
	}

	if getID.UserID != userID {
		log.Println("tidak memiliki akses")
		return book.Core{}, errors.New("tidak memiliki akses")
	}

	cnv := CoreToData(updatedData)
	qry := bd.db.Where("id = ?", bookID).Updates(&cnv)
	if qry.RowsAffected <= 0 {
		log.Println("update book query error : data not found")
		return book.Core{}, errors.New("not found")
	}

	if err := qry.Error; err != nil {
		log.Println("update book query error :", err.Error())
		return book.Core{}, err
	}

	return ToCore(cnv), nil
}

func (bd *bookData) Delete(userID uint, bookID uint) error {
	getID := Books{}
	err := bd.db.Where("id = ? ", bookID).First(&getID).Error

	if err != nil {
		log.Println("get user book error", err.Error())
		return errors.New("failed to get user book data")
	}

	if getID.UserID != userID {
		log.Println("tidak memiliki akses")
		return errors.New("tidak memiliki akses")
	}

	qry := bd.db.Delete(&Books{}, bookID)

	affRows := qry.RowsAffected

	if affRows <= 0 {
		log.Println("no rows affected")
		return errors.New("failed to delete user book, data not found")
	}
	return nil

}

// func (bd *bookData) MyBook(userID int) ([]book.Core, error) {
// 	return nil, nil
// }
