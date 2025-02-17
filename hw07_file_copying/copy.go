package main

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrCopyToTheSameFile     = errors.New("copy to the same file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFiledToMovePointer    = errors.New("filed to move pointer")
	ErrCopyAndRemove         = errors.New("error copy, unable to remove file")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// проверяем исходный файл
	fromFileInfo, err := os.Stat(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}

	// это обычный файл
	if !fromFileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	// смещение не больше размера файла
	if offset > fromFileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	// проверяем копируемый файл
	if toFileInfo, err := os.Stat(toPath); err == nil {
		// если файл существует,
		// проверим что он отличается от исходного файла
		if os.SameFile(fromFileInfo, toFileInfo) {
			return ErrCopyToTheSameFile
		}
	}

	if limit < 0 {
		limit = 0
	}

	// определяем сколько байт копировать с учетом смещения
	// учитываем лимит если указан, иначе копируем до конца файла
	size := fromFileInfo.Size() - offset
	if limit > 0 && limit < size {
		size = limit
	}

	// открываем исходный файл
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		// обработка самой игнорируемой ошибки, просто пишем в лог
		if err := f.Close(); err != nil {
			log.Println("error closing file:", err)
		}
	}(fromFile)

	// создадим или откроем копируемый файл
	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		// обработка самой игнорируемой ошибки, просто пишем в лог
		if err := f.Close(); err != nil {
			log.Println("error closing file:", err)
		}
	}(toFile)

	// смещаем указатель чтения файла
	if _, err := fromFile.Seek(offset, io.SeekStart); err != nil {
		return ErrFiledToMovePointer
	}

	// обертка для прогресс-бара
	bar := pb.Full.Start64(size)
	barReader := bar.NewProxyReader(fromFile)
	defer bar.Finish()

	// копируем данные
	_, err = io.CopyN(toFile, barReader, size)
	if err != nil && !errors.Is(err, io.EOF) {
		if errRemove := os.Remove(toPath); errRemove != nil {
			return ErrCopyAndRemove
		}
		return err
	}

	return nil
}
