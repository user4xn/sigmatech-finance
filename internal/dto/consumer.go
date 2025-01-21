package dto

import "mime/multipart"

type (
	ResponseMyProfile struct {
		NIK        string `json:"nik"`
		FullName   string `json:"full_name"`
		LegalName  string `json:"legal_name"`
		BirthDate  string `json:"birth_date"`
		BirthPlace string `json:"birth_place"`
		Salary     int    `json:"salary"`
		KTPURL     string `json:"ktp_url"`
		SelfieURL  string `json:"selfie_url"`
	}

	PayloadMyProfile struct {
		NIK        string                `form:"nik" binding:"required"`
		FullName   string                `form:"full_name" binding:"required"`
		LegalName  string                `form:"legal_name" binding:"required"`
		BirthDate  string                `form:"birth_date" binding:"required"`
		BirthPlace string                `form:"birth_place" binding:"required"`
		Salary     int                   `form:"salary" binding:"required"`
		FileKTP    *multipart.FileHeader `form:"file_ktp"`
		FileSelfie *multipart.FileHeader `form:"file_selfie"`
		KTPURL     string
		SelfieURL  string
	}
)
