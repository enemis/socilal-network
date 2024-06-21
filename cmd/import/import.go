package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"social-network-otus/internal/app"
)

func init() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		lvl = "debug"
	}
	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.DebugLevel
	}
	logrus.SetLevel(ll)
}

func main() {
	_, err := app.NewSeeder()
	if err != nil {
		logrus.Fatal(err.Error())
	}
}
