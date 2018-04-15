package cmd

import (
	"testing"
	"github.com/daviddengcn/go-assert"
	"path"
	"github.com/spf13/viper"
)

// Check file exists
func TestExists(t *testing.T) {
	fileExists, err := exists(path.Join("fixtures", "config.yml"))
	assert.Equals(t, "error", err, nil)
	assert.IsTrue(t, "file exists", fileExists)
}

// Check file does not exist
func TestExistsNot(t *testing.T) {
	fileExists, err := exists(path.Join("fixtures", "missing.yml"))
	assert.Equals(t, "error", err, nil)
	assert.IsFalse(t, "file does not exist", fileExists)
}

// Check file but whole directory does not exist
func TestExistsNotDirectory(t *testing.T) {
	fileExists, err := exists(path.Join("fixtures", "missing", "missing.yml"))
	assert.Equals(t, "error", err, nil)
	assert.IsFalse(t, "file does not exist", fileExists)
}

func testReadConfig(t *testing.T, configName string) {
	readConfig(configName)

	assert.IsTrue(t, "BoolTrue1", viper.GetBool("BoolTrue1"))
	assert.IsTrue(t, "BoolTrue2", viper.GetBool("BoolTrue2"))
	assert.IsTrue(t, "BoolTrue3", viper.GetBool("BoolTrue3"))
	assert.IsTrue(t, "BoolTrue4", viper.GetBool("BoolTrue4"))
	assert.IsTrue(t, "BoolTrue5", viper.GetBool("BoolTrue5"))

	assert.IsFalse(t, "BoolFalse1", viper.GetBool("BoolFalse1"))
	assert.IsFalse(t, "BoolFalse2", viper.GetBool("BoolFalse2"))
	assert.IsFalse(t, "BoolFalse3", viper.GetBool("BoolFalse3"))
	assert.IsFalse(t, "BoolFalse4", viper.GetBool("BoolFalse4"))
	assert.IsFalse(t, "BoolFalse5", viper.GetBool("BoolFalse5"))

	//viper.GetString()
	//viper.GetStringSlice()
	//viper.GetStringMapString()
}

// Read in yaml config into viper
func TestReadConfigYaml(t *testing.T) {
	testReadConfig(t, path.Join("fixtures", "config.yml"))
}

// Read in toml config into viper
func TestReadConfigToml(t *testing.T) {
	testReadConfig(t, path.Join("fixtures", "config.toml"))
}

// Read in json config into viper
func TestReadConfigJson(t *testing.T) {
	testReadConfig(t, path.Join("fixtures", "config.json"))
}
