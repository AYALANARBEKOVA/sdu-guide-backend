package services

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"sdu-guide/internal/repositories"
	"sdu-guide/internal/structures"
	"sdu-guide/internal/utils"
)

type File struct {
	repo *repositories.Repository
}

func newFileService(repo *repositories.Repository) *File {
	return &File{repo: repo}
}

func (m *File) StoreFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	defer file.Close()

	hash := utils.RandomString(6)
	fileName := hash + "_" + header.Filename

	dir := "./uploads"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	filePath := filepath.Join(dir, fileName)

	out, err := os.Create(filePath)
	if err != nil {

		return hash, err
	}
	defer out.Close()

	if _, err = io.Copy(out, file); err != nil {

		return hash, err
	}

	xlsx := structures.File{
		Name:     header.Filename,
		Hash:     hash,
		Path:     filePath,
		FileType: filepath.Ext(header.Filename),
	}

	if err := m.repo.CreateFile(xlsx); err != nil {
		return hash, err
	}

	return hash, nil
}

func (m *File) GetFilebyHash(hash string) (structures.File, error) {
	return m.repo.FileRepo.GetFile(hash)
}
