package config

import (
	"context"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var FireBaseApp *firebase.App

func IntializeFireBaseApp() {
	opt := option.WithCredentialsFile("hostego-firebase-secret.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)

	if err != nil {
		log.Fatal("Error Initializing Firebase", err)
	}
	FireBaseApp = app
	log.Println("FireBase initiliazed successfully !")
}
