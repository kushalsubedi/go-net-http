package student

import (
	"encoding/json"
	"errors"
	"fmt"

	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/kushalsubedi/students-api/internal/types"
	"github.com/kushalsubedi/students-api/internal/utils/response"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		if err := validator.New().Struct(student); err != nil {
			if validateErrs, ok := err.(validator.ValidationErrors); ok {
				response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			} else {
				response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			}
		}

		w.Write([]byte("Welcome to students api "))
		response.WriteJson(w, http.StatusCreated, map[string]string{"success": "OK"})
	}
}
