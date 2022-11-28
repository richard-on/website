package test

import (
	"github.com/joho/godotenv"
	"github.com/richard-on/website/config"
	"testing"
)

func TestInit(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Log(err)
	}

	err = config.Init()
	if err != nil {
		t.Fatal(err)
	}
}
