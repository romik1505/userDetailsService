package repository

import (
	"os"
	"testing"
)

var (
	personsRepo *PersonsRepo
)

func TestMain(m *testing.M) {
	personsRepo = NewPersonsRepo()
	exit := m.Run()
	os.Exit(exit)
}
