package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Gurveer1510/student-api/internal/storage"
	"github.com/Gurveer1510/student-api/internal/types"
	"github.com/Gurveer1510/student-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("Creating a Student")

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty payload")))
			return
		}
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Request validation

		if err := validator.New().Struct(student); err != nil {
			validationErrs := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validationErrs))
			return
		}

		id, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err)
			return
		}
		slog.Info("Student Added successfully", slog.String("name", student.Name))
		response.WriteJSON(w, http.StatusCreated, map[string]int64{"id": id})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("getting a student", slog.String("id", id))

		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		student, err := storage.GetStudentById(idInt)
		if err != nil {
			slog.Error("Error getting user", slog.String("id", id))
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJSON(w, http.StatusOK, student)

	}
}
