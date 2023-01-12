package services

import (
	"api/features/book"
	"api/helper"
	"api/mocks"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	repo := mocks.NewBookData(t)

	t.Run("berhasil tambah buku", func(t *testing.T) {
		inputBook := book.Core{Judul: "One Piece", TahunTerbit: 1997, Penulis: "Eichiro Oda"}
		resBook := book.Core{ID: uint(1), Judul: "One Piece", TahunTerbit: 1997, Penulis: "Eichiro Oda"}
		repo.On("Add", uint(1), inputBook).Return(resBook, nil).Once()

		srv := New(repo)
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Add(pToken, inputBook)
		assert.Nil(t, err)
		assert.Equal(t, resBook.ID, res.ID)
		repo.AssertExpectations(t)

	})

	t.Run("masalah di server", func(t *testing.T) {
		inputBook := book.Core{Judul: "One Piece", TahunTerbit: 1997, Penulis: "Eichiro Oda"}
		repo.On("Add", uint(1), inputBook).Return(book.Core{}, errors.New("terdapat masalah pada server")).Once()
		srv := New(repo)
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Add(pToken, inputBook)
		assert.NotNil(t, err)
		assert.Equal(t, uint(0), res.ID)
		assert.ErrorContains(t, err, "server")
		repo.AssertExpectations(t)
	})

	t.Run("user tidak ditemukan", func(t *testing.T) {
		inputBook := book.Core{Judul: "One Piece", TahunTerbit: 1997, Penulis: "Eichiro Oda"}
		repo.On("Add", uint(1), inputBook).Return(book.Core{}, errors.New("not found")).Once()
		srv := New(repo)
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Add(pToken, inputBook)
		assert.NotNil(t, err)
		assert.Equal(t, uint(0), res.ID)
		assert.ErrorContains(t, err, "not found")
		repo.AssertExpectations(t)
	})

	t.Run("jwt tidak valid", func(t *testing.T) {
		inputBook := book.Core{Judul: "One Piece", TahunTerbit: 1997, Penulis: "Eichiro Oda"}
		srv := New(repo)

		_, token := helper.GenerateJWT(1)
		res, err := srv.Add(token, inputBook)
		assert.NotNil(t, err)
		assert.Equal(t, uint(0), res.ID)
		assert.ErrorContains(t, err, "not found")
	})
}

func TestUpdate(t *testing.T) {
	repo := mocks.NewBookData(t)

	t.Run("suskes update data", func(t *testing.T) {
		inputBook := book.Core{Judul: "Naruto", TahunTerbit: 1999, Penulis: "Masashi Kishimoto"}
		resBook := book.Core{ID: uint(1), Judul: "Naruto", TahunTerbit: 1999, Penulis: "Masashi Kishimoto"}
		repo.On("Update", uint(1), uint(1), inputBook).Return(resBook, nil).Once()

		srv := New(repo)
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Update(pToken, uint(1), inputBook)
		assert.Nil(t, err)
		assert.Equal(t, resBook.ID, res.ID)
		assert.Equal(t, inputBook.Judul, res.Judul)
		assert.Equal(t, inputBook.TahunTerbit, res.TahunTerbit)
		assert.Equal(t, inputBook.Penulis, res.Penulis)
		repo.AssertExpectations(t)

	})

	t.Run("jwt tidak valid", func(t *testing.T) {
		inputBook := book.Core{Judul: "One Piece", TahunTerbit: 1997, Penulis: "Eichiro Oda"}
		srv := New(repo)

		_, token := helper.GenerateJWT(0)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Update(pToken, 1, inputBook)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
		assert.Equal(t, uint(0), res.ID)
	})

	t.Run("data tidak ditemukan", func(t *testing.T) {
		inputBook := book.Core{Judul: "One Piece", TahunTerbit: 1997, Penulis: "Eichiro Oda"}
		repo.On("Update", uint(2), uint(2), inputBook).Return(book.Core{}, errors.New("data not found")).Once()

		srv := New(repo)
		_, token := helper.GenerateJWT(2)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Update(pToken, 2, inputBook)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
		assert.Equal(t, uint(0), res.ID)
		repo.AssertExpectations(t)
	})

	t.Run("masalah di server", func(t *testing.T) {
		inputBook := book.Core{Judul: "One Piece", TahunTerbit: 1997, Penulis: "Eichiro Oda"}
		repo.On("Update", uint(1), uint(1), inputBook).Return(book.Core{}, errors.New("terdapat masalah pada server")).Once()

		srv := New(repo)
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Update(pToken, 1, inputBook)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		assert.Equal(t, uint(0), res.ID)
		repo.AssertExpectations(t)
	})

}

func TestDelete(t *testing.T) {
	repo := mocks.NewBookData(t)

	t.Run("suskes hapus buku", func(t *testing.T) {
		repo.On("Delete", uint(1), uint(1)).Return(nil).Once()

		srv := New(repo)
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		err := srv.Delete(pToken, 1)
		assert.Nil(t, err)
		repo.AssertExpectations(t)

	})

	t.Run("jwt tidak valid", func(t *testing.T) {
		srv := New(repo)

		_, token := helper.GenerateJWT(0)
		err := srv.Delete(token, 1)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
	})

	t.Run("data tidak ditemukan", func(t *testing.T) {
		repo.On("Delete", uint(2), uint(2)).Return(errors.New("data not found")).Once()

		srv := New(repo)
		_, token := helper.GenerateJWT(2)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		err := srv.Delete(pToken, 2)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
		repo.AssertExpectations(t)
	})

}
