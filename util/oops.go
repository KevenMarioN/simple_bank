package util

import (
	"fmt"
	"log"

	"github.com/lib/pq"
)

func HttpError(err error) (result error) {

	if pqErr, ok := err.(*pq.Error); ok {
		result = errorsPQ(pqErr)
	}
	if result == nil {
		return err
	}
	return result
}

func errorsPQ(pqErr *pq.Error) error {
	log.Println(pqErr.Code)
	switch pqErr.Code {
	case "23503":
		return fmt.Errorf("Dado indicado não é uma referência válida")
	case "23505":
		return fmt.Errorf("Registro duplicado")
	}
	return nil
}
