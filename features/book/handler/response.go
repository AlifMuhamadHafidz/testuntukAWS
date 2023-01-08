package handler

import "api/features/book"

type BookResponse struct {
	ID          uint   `json:"id"`
	Judul       string `json:"judul"`
	TahunTerbit int    `json:"tahun_terbit"`
	Penulis     string `json:"penulis"`
	Pemilik     string `json:"pemilik"`
}

type AddBookResponse struct {
	Judul       string `json:"judul"`
	TahunTerbit int    `json:"tahun_terbit"`
	Penulis     string `json:"penulis"`
}

func ToResponse(data book.Core) BookResponse {
	return BookResponse{
		ID:          data.ID,
		Judul:       data.Judul,
		TahunTerbit: data.TahunTerbit,
		Penulis:     data.Penulis,
		Pemilik:     data.Pemilik,
	}
}
