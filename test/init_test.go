//go:build e2e
// +build e2e

package test

import (
	"context"
	"testing"

	"github.com/cryptellation/exchanges/configs"
	"github.com/cryptellation/exchanges/pkg/client"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

func TestEndToEndSuite(t *testing.T) {
	suite.Run(t, new(EndToEndSuite))
}

type EndToEndSuite struct {
	suite.Suite
	client client.Client
}

func (suite *EndToEndSuite) SetupSuite() {
	client, err := client.New(
		viper.GetString(configs.EnvTemporalAddress),
	)
	suite.Require().NoError(err)
	suite.client = client
}

func (suite *EndToEndSuite) TearDownSuite() {
	suite.client.Close(context.Background())
}
