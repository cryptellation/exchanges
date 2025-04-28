//go:build integration
// +build integration

package sql

import (
	"context"
	"testing"

	"github.com/cryptellation/dbmigrator"
	"github.com/cryptellation/exchanges/configs"
	"github.com/cryptellation/exchanges/configs/sql/down"
	"github.com/cryptellation/exchanges/configs/sql/up"
	"github.com/cryptellation/exchanges/svc/db"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

func TestExchangesSuite(t *testing.T) {
	suite.Run(t, new(ExchangesSuite))
}

type ExchangesSuite struct {
	db.ExchangesSuite
}

func (suite *ExchangesSuite) SetupSuite() {
	act, err := New(context.Background(), viper.GetString(configs.EnvDBDSN))
	suite.Require().NoError(err)

	mig, err := dbmigrator.NewMigrator(context.Background(), act.db, up.Migrations, down.Migrations, nil)
	suite.Require().NoError(err)
	suite.Require().NoError(mig.MigrateToLatest(context.Background()))

	suite.DB = act
}

func (suite *ExchangesSuite) SetupTest() {
	db := suite.DB.(*Activities)
	suite.Require().NoError(db.Reset(context.Background()))
}
