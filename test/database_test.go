package test

import (
	"testing"
	"twittir-go/internal/database"

	"github.com/stretchr/testify/assert"
)

func TestConnectDB(t *testing.T) {
	err := database.ConnectDB()

	// Assert that no error occurs
	assert.NoError(t, err, "Database connection should not return an error")

	// Assert that the global DB variable is not nil
	db := database.GetDB()
	assert.NotNil(t, db, "Database connection should not be nil")
}
