package utility

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const ENV_KEY_PORT = "FRED_PORT_F9mF"

// TestHelloWorld
func TestEnvironmentVariableSet(t *testing.T) {
	port := Env(ENV_KEY_PORT, "8080")

	assert.Equal(t, "8080", port, "The port value has to be equal to the defined default of 8080 because port env var doesnt exist.")
}

func TestEnvironmentVariableNotSet(t *testing.T) {
	os.Setenv(ENV_KEY_PORT, "80")
	defer os.Unsetenv(ENV_KEY_PORT)

	port := Env(ENV_KEY_PORT, "8080")

	assert.Equal(t, "80", port, "The port value has to be read from the env var PORT and be 80.")
}
