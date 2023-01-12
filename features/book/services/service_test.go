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
}
