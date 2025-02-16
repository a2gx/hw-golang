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
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	log.Printf("COPY from=%q to=%q offset=%d limit=%d", fromPath, toPath, offset, limit)

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
	sizeToCopy := fromFileInfo.Size() - offset
	if limit > 0 && limit < sizeToCopy {
		sizeToCopy = limit
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

	// обертка для прогресс-бара
	bar := pb.Full.Start64(sizeToCopy)
	barReader := bar.NewProxyReader(fromFile)
	defer bar.Finish()

	// копируем данные
	_, err = io.CopyN(toFile, barReader, sizeToCopy)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	return nil
}
