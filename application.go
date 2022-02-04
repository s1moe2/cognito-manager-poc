package main

import (
	"context"
	"encoding/json"
	"errors"
	awsconf "github.com/aws/aws-sdk-go-v2/config"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

type application struct {
	cognitoClient *cip.Client
}

func initApplication() *application {
	awsConf, err := awsconf.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return &application{
		cognitoClient: cip.NewFromConfig(awsConf),
	}
}

func (app *application) serve() {
	server := &http.Server{
		Addr:         ":4004",
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Server started on %s", server.Addr)
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodPost, "/clients", app.newClientHandler)
	router.HandlerFunc(http.MethodDelete, "/clients/:id", app.deleteClientHandler)

	return router
}

type AppError struct {
	Message    string `json:"message"`
	error      error
	statusCode int
}

func (app *application) respondError(w http.ResponseWriter, e AppError) {
	if e.error != nil {
		log.Println(e.error)
	}

	jsonData, err := json.Marshal(e)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.statusCode)
	w.Write(jsonData)
}

func (app *application) respondJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	response := struct {
		Data interface{} `json:"data"`
	}{
		Data: data,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		app.respondError(w, AppError{
			Message:    "internal error",
			error:      err,
			statusCode: http.StatusInternalServerError,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonData)
}
