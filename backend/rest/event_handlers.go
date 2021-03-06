package rest

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	R "github.com/go-pkgz/rest"
	"github.com/ts-dmitry/cronpad/backend/repository"
	"net/http"
)

type eventHandlers struct {
	service EventService
}

type EventService interface {
	Create(record repository.Event, userID string) (string, error)
	Update(event repository.Event, userID string) (string, error)
	Delete(eventID string, userID string) error
}

func (t *eventHandlers) create(writer http.ResponseWriter, request *http.Request) {
	user, err := GetUserInfo(request)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusForbidden, err, "user should be logged in", ErrInternal)
		return
	}

	var event repository.Event
	err = json.NewDecoder(request.Body).Decode(&event)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't parse json", ErrInternal)
		return
	}

	event.PrepareReceived()

	id, err := t.service.Create(event, user.ID)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't create event", ErrInternal)
		return
	}

	render.Status(request, http.StatusCreated)
	render.JSON(writer, request, R.JSON{"id": id})
}

func (t *eventHandlers) update(writer http.ResponseWriter, request *http.Request) {
	user, err := GetUserInfo(request)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusForbidden, err, "user should be logged in", ErrInternal)
		return
	}

	id := chi.URLParam(request, "id")
	if len(id) == 0 {
		SendErrorJSON(writer, request, http.StatusBadRequest, errors.New("id can't be empty"), "", ErrInternal)
		return
	}

	var event repository.Event
	err = json.NewDecoder(request.Body).Decode(&event)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't parse json", ErrInternal)
		return
	}

	event.PrepareReceived()
	event.ID = id

	id, err = t.service.Update(event, user.ID)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusBadRequest, err, "can't update event", ErrInternal)
		return
	}

	render.Status(request, http.StatusCreated)
	render.JSON(writer, request, R.JSON{"id": id})
}

func (t *eventHandlers) delete(writer http.ResponseWriter, request *http.Request) {
	user, err := GetUserInfo(request)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusForbidden, err, "user should be logged in", ErrInternal)
		return
	}

	id := chi.URLParam(request, "id")

	err = t.service.Delete(id, user.ID)
	if err != nil {
		SendErrorJSON(writer, request, http.StatusInternalServerError, err, "can't delete event", ErrInternal)
		return
	}

	render.Status(request, http.StatusOK)
	render.JSON(writer, request, R.JSON{"id": id})
}
