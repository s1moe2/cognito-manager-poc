package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type newClientPayload struct {
	Name string `json:"name"`
}

type newClientResponse struct {
	PoolID       string `json:"poolID"`
	PoolClientID string `json:"poolClientID"`
}

func (app *application) newClientHandler(w http.ResponseWriter, r *http.Request) {
	var payload newClientPayload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		app.respondError(w, AppError{
			Message:    "failed to decode newClientPayload",
			error:      err,
			statusCode: http.StatusInternalServerError,
		})
		return
	}

	pool, err := app.createCognitoUserPool(r.Context(), payload.Name)

	if err != nil {
		app.respondError(w, AppError{
			Message:    "failed to create pool",
			error:      err,
			statusCode: http.StatusInternalServerError,
		})
		return
	}

	client, err := app.createCognitoUserPoolClient(r.Context(), fmt.Sprintf("%s_client", *pool.UserPool.Name), *pool.UserPool.Id)
	if err != nil {
		delErr := app.deleteCognitoUserPool(r.Context(), *pool.UserPool.Id)
		if delErr != nil {
			log.Println(delErr)
			app.respondError(w, AppError{
				Message:    "failed to create pool client; failed to rollback and delete pool",
				error:      err,
				statusCode: http.StatusInternalServerError,
			})
			return
		}
		app.respondError(w, AppError{
			Message:    "failed to create pool client",
			error:      err,
			statusCode: http.StatusInternalServerError,
		})
		return
	}

	app.respondJSON(w, newClientResponse{
		PoolID:       *pool.UserPool.Id,
		PoolClientID: *client.UserPoolClient.ClientId,
	}, http.StatusOK)
}

func (app *application) deleteClientHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	poolID := params.ByName("id")
	if len(poolID) == 0 {
		app.respondError(w, AppError{
			Message:    "bad client id",
			statusCode: http.StatusBadRequest,
		})
		return
	}

	err := app.deleteCognitoUserPool(r.Context(), poolID)
	if err != nil {
		log.Println(err)
		app.respondError(w, AppError{
			Message:    "failed to delete pool",
			error:      err,
			statusCode: http.StatusInternalServerError,
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
