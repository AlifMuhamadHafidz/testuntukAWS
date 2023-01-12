package services

import (
	"api/features/user"
	"api/helper"
	"api/mocks"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegister(t *testing.T) {
	repo := mocks.NewUserData(t)

	t.Run("Berhasil Register", func(t *testing.T) {
		inputData := user.Core{Nama: "alif", Email: "alif@be14.com", Alamat: "bangka", HP: "088", Password: "alif123"}
		resData := user.Core{ID: uint(1), Nama: "alif", Email: "alif@be14.com", Alamat: "bangka", HP: "088"}
		repo.On("Register", mock.Anything).Return(resData, nil).Once()
		srv := New(repo)
		res, err := srv.Register(inputData)
		assert.Nil(t, err)
		assert.Equal(t, resData.ID, res.ID)
		assert.Equal(t, resData.Nama, res.Nama)
		assert.Equal(t, resData.Alamat, res.Alamat)
		assert.Equal(t, resData.HP, res.HP)
		repo.AssertExpectations(t)
	})

	t.Run("masalah di server", func(t *testing.T) {
		inputData := user.Core{Nama: "alif", Email: "alif@be14.com", Alamat: "bangka", HP: "088", Password: "alif123"}
		resData := user.Core{ID: uint(1), Nama: "alif", Email: "alif@be14.com", Alamat: "bangka", HP: "088"}
		repo.On("Register", mock.Anything).Return(resData, errors.New("terdapat masalah pada server")).Once()
		srv := New(repo)
		res, err := srv.Register(inputData)
		assert.NotNil(t, err)
		assert.Equal(t, uint(0), res.ID)
		assert.ErrorContains(t, err, "server")
		repo.AssertExpectations(t)
	})

	t.Run("data sudah terdaftar", func(t *testing.T) {
		inputData := user.Core{Nama: "alif", Email: "alif@be14.com", Alamat: "bangka", HP: "088", Password: "alif123"}
		// resData := user.Core{ID: uint(1), Nama: "alif", Email: "alif@be14.com", Alamat: "bangka", HP: "088"}
		repo.On("Register", mock.Anything).Return(user.Core{}, errors.New("duplicated")).Once()
		srv := New(repo)
		res, err := srv.Register(inputData)
		assert.NotNil(t, err)
		assert.Equal(t, uint(0), res.ID)
		assert.ErrorContains(t, err, "sudah terdaftar")
		repo.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	repo := mocks.NewUserData(t) // mock data

	t.Run("Berhasil login", func(t *testing.T) {
		// input dan respond untuk mock data
		inputEmail := "alif@be14.com"
		// res dari data akan mengembalik password yang sudah di hash
		hashed, _ := helper.GeneratePassword("be1422")
		resData := user.Core{ID: uint(1), Nama: "alif", Email: "alif@be14.com", HP: "088888", Password: hashed}

		repo.On("Login", inputEmail).Return(resData, nil) // simulasi method login pada layer data

		srv := New(repo)
		token, res, err := srv.Login(inputEmail, "be1422")
		assert.Nil(t, err)
		assert.NotEmpty(t, token)
		assert.Equal(t, resData.ID, res.ID)
		repo.AssertExpectations(t)
	})

	t.Run("Tidak ditemukan", func(t *testing.T) {
		inputEmail := "alif@be14.com"
		repo.On("Login", inputEmail).Return(user.Core{}, errors.New("data not found"))

		srv := New(repo)
		token, res, err := srv.Login(inputEmail, "be1422")
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "tidak ditemukan")
		assert.Empty(t, token)
		assert.Equal(t, uint(0), res.ID)
		repo.AssertExpectations(t)
	})

	t.Run("Salah password", func(t *testing.T) {
		inputEmail := "alif@be14.com"
		hashed, _ := helper.GeneratePassword("be1422")
		resData := user.Core{ID: uint(1), Nama: "alif", Email: "alif@be14.com", HP: "088888", Password: hashed}
		repo.On("Login", inputEmail).Return(resData, nil)

		srv := New(repo)
		token, res, err := srv.Login(inputEmail, "be1423")
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "password tidak sesuai")
		assert.Empty(t, token)
		assert.Equal(t, uint(0), res.ID)
		repo.AssertExpectations(t)
	})

	t.Run("masalah di server", func(t *testing.T) {
		inputEmail := "alif@be14.com"
		hashed, _ := helper.GeneratePassword("be1422")
		resData := user.Core{ID: uint(1), Nama: "alif", Email: "alif@be14.com", HP: "088888", Password: hashed}
		repo.On("Login", inputEmail).Return(resData, errors.New("terdapat masalah pada server")).Once()

		srv := New(repo)
		token, res, err := srv.Login(inputEmail, "be1423")
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		assert.Empty(t, token)
		assert.Equal(t, uint(0), res.ID)
		repo.AssertExpectations(t)
	})

}

func TestProfile(t *testing.T) {
	repo := mocks.NewUserData(t)

	t.Run("Sukses lihat profile", func(t *testing.T) {
		resData := user.Core{ID: uint(1), Nama: "alif", Email: "alif@be14.com", HP: "088888"}

		repo.On("Profile", uint(1)).Return(resData, nil).Once()

		srv := New(repo)

		_, token := helper.GenerateJWT(1)

		pToken := token.(*jwt.Token)
		pToken.Valid = true

		res, err := srv.Profile(pToken)
		assert.Nil(t, err)
		assert.Equal(t, resData.ID, res.ID)
		repo.AssertExpectations(t)
	})

	t.Run("jwt tidak valid", func(t *testing.T) {
		srv := New(repo)

		_, token := helper.GenerateJWT(1)

		res, err := srv.Profile(token)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "tidak ditemukan")
		assert.Equal(t, uint(0), res.ID)
	})

	t.Run("data tidak ditemukan", func(t *testing.T) {
		repo.On("Profile", uint(4)).Return(user.Core{}, errors.New("data not found")).Once()

		srv := New(repo)

		_, token := helper.GenerateJWT(4)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Profile(pToken)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "tidak ditemukan")
		assert.Equal(t, uint(0), res.ID)
		repo.AssertExpectations(t)
	})

	t.Run("masalah di server", func(t *testing.T) {
		repo.On("Profile", mock.Anything).Return(user.Core{}, errors.New("terdapat masalah pada server")).Once()
		srv := New(repo)

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Profile(pToken)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		assert.Equal(t, uint(0), res.ID)
		repo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	repo := mocks.NewUserData(t)

	t.Run("suskes update data", func(t *testing.T) {
		input := user.Core{Nama: "alip", Email: "alip@be14.com", HP: "08888"}
		hashed, _ := helper.GeneratePassword("be1422")
		resData := user.Core{ID: uint(1), Nama: "alip", Email: "alip@be14.com", HP: "08888", Password: hashed}
		repo.On("Update", uint(1), input).Return(resData, nil).Once()

		srv := New(repo)
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Update(pToken, input)
		assert.Nil(t, err)
		assert.Equal(t, resData.ID, res.ID)
		assert.Equal(t, input.Nama, res.Nama)
		assert.Equal(t, input.Email, res.Email)
		assert.Equal(t, input.HP, res.HP)
		repo.AssertExpectations(t)

	})

	t.Run("jwt tidak valid", func(t *testing.T) {
		input := user.Core{Nama: "alif", Email: "alif@be14.com", HP: "088"}
		srv := New(repo)

		_, token := helper.GenerateJWT(0)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Update(pToken, input)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
		assert.Equal(t, uint(0), res.ID)
	})

	t.Run("data tidak ditemukan", func(t *testing.T) {
		input := user.Core{Nama: "alif", Email: "alif@be14.com", HP: "088"}
		repo.On("Update", uint(2), input).Return(user.Core{}, errors.New("data not found")).Once()

		srv := New(repo)
		_, token := helper.GenerateJWT(2)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Update(pToken, input)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "tidak ditemukan")
		assert.Equal(t, uint(0), res.ID)
		repo.AssertExpectations(t)
	})

	t.Run("masalah di server", func(t *testing.T) {
		input := user.Core{Nama: "alif", Email: "alif@be14.com", HP: "088"}
		repo.On("Update", uint(1), input).Return(user.Core{}, errors.New("terdapat masalah pada server")).Once()

		srv := New(repo)
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Update(pToken, input)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		assert.Equal(t, uint(0), res.ID)
		repo.AssertExpectations(t)
	})

}

func TestDeactive(t *testing.T) {
	repo := mocks.NewUserData(t)

	t.Run("suskes hapus profile", func(t *testing.T) {
		repo.On("Deactive", uint(1)).Return(nil).Once()

		srv := New(repo)
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		err := srv.Deactive(pToken)
		assert.Nil(t, err)
		repo.AssertExpectations(t)

	})

	t.Run("jwt tidak valid", func(t *testing.T) {
		srv := New(repo)

		_, token := helper.GenerateJWT(1)
		err := srv.Deactive(token)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
	})

	t.Run("data tidak ditemukan", func(t *testing.T) {
		repo.On("Deactive", uint(2)).Return(errors.New("data not found")).Once()

		srv := New(repo)
		_, token := helper.GenerateJWT(2)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		err := srv.Deactive(pToken)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "tidak ditemukan")
		repo.AssertExpectations(t)
	})

	t.Run("masalah di server", func(t *testing.T) {
		repo.On("Deactive", mock.Anything).Return(errors.New("terdapat masalah pada server")).Once()

		srv := New(repo)
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		err := srv.Deactive(pToken)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		repo.AssertExpectations(t)
	})
}
