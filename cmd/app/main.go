package main

import (
	"github.com/sirupsen/logrus"

	"social-network-otus/internal/app"
)

func init() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
}

func main() {
	a, err := app.NewApp()
	if err != nil {
		logrus.Fatal(err)
	}

	err = a.Run()
	if err != nil {
		logrus.Fatal(err)
	}
}
