//go:build unit
// +build unit

package configs

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

func TestViperSuite(t *testing.T) {
	suite.Run(t, new(ViperSuite))
}

type ViperSuite struct {
	suite.Suite
}

func (suite *ViperSuite) TestDBDSN() {
	// Test the default value of the database DSN
	suite.Equal(DefaultDBDSN, viper.GetString(EnvDBDSN))

	// Set environment variable for the database DSN
	os.Setenv(strings.ToUpper(EnvDBDSN), "test")

	// Test the overridden value of the database DSN
	suite.Equal("test", viper.GetString(EnvDBDSN))
}
