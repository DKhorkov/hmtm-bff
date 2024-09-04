package loadenv_test

import (
	"hmtmbff/pkg/loadenv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	t.Run("key exists", func(t *testing.T) {
		t.Setenv("TEST_KEY", "GRAPHQL_PORT")
		want := "GRAPHQL_PORT"
		got := loadenv.GetEnv("TEST_KEY", "")
		assert.Equal(t, want, got, "\nGetEnv should return a value from an environment variable.")
	})

	t.Run("key does not exist", func(t *testing.T) {
		want := ""
		got := loadenv.GetEnv("NON_EXISTENT_KEY", want)
		assert.Equal(t, want, got, "\nGetEnv should return the default value if the key does not exist.")
	})
}

func TestGetEnvAsInt(t *testing.T) {
	t.Run("env var exists and is valid integer", func(t *testing.T) {
		t.Setenv("TEST_INT", "8081")
		want := 8081
		got := loadenv.GetEnvAsInt("TEST_INT", 8080)
		assert.Equal(t, want, got, "\nGetEnvAsInt should return an integer value from an environment "+
			"variable.")
	})

	t.Run("env var exists but is invalid integer", func(t *testing.T) {
		t.Setenv("TEST_INVALID_INT", "abc")
		want := 8080
		got := loadenv.GetEnvAsInt("TEST_INVALID_INT", 8080)
		assert.Equal(t, want, got, "\nGetEnvAsInt should return the default value if the environment"+
			" variable is not a valid integer.")
	})

	t.Run("env var does not exist", func(t *testing.T) {
		want := 8080
		got := loadenv.GetEnvAsInt("NON_EXISTENT_KEY", 8080)
		assert.Equal(t, want, got, "\nGetEnvAsInt should return the default value if the environment"+
			" variable does not exist")
	})
}

func TestGetEnvAsSlice(t *testing.T) {
	t.Run("env var exists but is invalid slice", func(t *testing.T) {
		t.Setenv("TEST_INVALID_SLICE", "fs")
		want := []string{"1", "2"}
		got := loadenv.GetEnvAsSlice("TEST_INVALID_SLICE", []string{"1", "2"}, ",")
		assert.Equal(t, want, got, "GetEnvAsSlice should return the default value if the environment "+
			"variable is not a valid slice.")
	})

	t.Run("env var exists but is invalid slice with different separator", func(t *testing.T) {
		t.Setenv("TEST_INVALID_SLICE", "fs")
		want := []string{"1", "2"}
		got := loadenv.GetEnvAsSlice("TEST_INVALID_SLICE", []string{"1", "2"}, "|")
		assert.Equal(t, want, got, "GetEnvAsSlice should return the default value if the environment "+
			"variable is not a valid slice with different separator.")
	})

	t.Run("env var exists and is valid slice", func(t *testing.T) {
		t.Setenv("TEST_VALID_SLICE", "fs,a,ass")
		want := []string{"fs", "a", "ass"}
		got := loadenv.GetEnvAsSlice("TEST_VALID_SLICE", []string{"1", "2"}, ",")
		assert.Equal(t, want, got, "GetEnvAsSlice should return the slice if the environment "+
			"variable is a valid slice.")
	})

	t.Run("env var exists and is valid slice with different separator", func(t *testing.T) {
		t.Setenv("TEST_VALID_SLICE", "fs|a|ass")
		want := []string{"fs", "a", "ass"}
		got := loadenv.GetEnvAsSlice("TEST_VALID_SLICE", []string{"1", "2"}, "|")
		assert.Equal(t, want, got, "GetEnvAsSlice should return the slice if the environment "+
			"variable is a valid slice with different separator.")
	})
}

func TestGetEnvAsBool(t *testing.T) {
	t.Run("env var exists and is true", func(t *testing.T) {
		t.Setenv("TEST_BOOL_TRUE", "true")
		got := loadenv.GetEnvAsBool("TEST_BOOL_TRUE", false)
		assert.Equal(t, true, got, "\nGetEnvAsBool should return true if the environment"+
			" variable is 'true'.")
	})

	t.Run("env var exists and is false", func(t *testing.T) {
		t.Setenv("TEST_BOOL_FALSE", "false")
		got := loadenv.GetEnvAsBool("TEST_BOOL_FALSE", false)
		assert.Equal(t, false, got, "\nGetEnvAsBool should return false if the environment"+
			" variable is 'false'.")
	})

	t.Run("env var exists but is invalid boolean", func(t *testing.T) {
		t.Setenv("TEST_BOOL_INVALID", "invalid")
		got := loadenv.GetEnvAsBool("TEST_BOOL_INVALID", false)
		assert.Equal(t, false, got, "\nGetEnvAsBool should return the default value if the"+
			" environment variable is not a valid boolean.")
	})

	t.Run("env var does not exist", func(t *testing.T) {
		got := loadenv.GetEnvAsBool("TEST_BOOL_NOT_EXIST", true)
		assert.Equal(t, true, got, "\nGetEnvAsBool should return the default value if the"+
			" environment variable does not exist.")
	})
}

func TestIsStringIsValidSlice(t *testing.T) {
	t.Run("slice is valid", func(t *testing.T) {
		got := loadenv.IsStringIsValidSlice("Uno,Dos,Tres", ",")
		assert.True(t, got, "\nTestIsStringIsValidSlice should return 'True' for valid slice")
	})

	t.Run("slice is not valid", func(t *testing.T) {
		got := loadenv.IsStringIsValidSlice("Uno,", ",")
		assert.False(t, got, "\nTestIsStringIsValidSlice should return return 'False' for not "+
			"valid slice")
	})

	t.Run("slice is valid with whitespaces", func(t *testing.T) {
		got := loadenv.IsStringIsValidSlice("Uno, Dos, Tres", ",")
		assert.True(t, got, "\nTestIsStringIsValidSlice should return 'True' for valid slice with "+
			"whitespaces between separated values")
	})
}
