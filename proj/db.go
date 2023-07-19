package proj

import (
	"context"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/db"
)

func UseAccountDbClient(ctx context.Context, driver gira.DbDao, config *gira.Config) error {
	if client, err := db.NewConfigDbClient(ctx, gira.ACCOUNTDB_NAME, *config.Db[gira.ACCOUNTDB_NAME]); err != nil {
		return err
	} else if err := driver.UseClient(client); err != nil {
		return err
	} else {
		return nil
	}
}

func UseGameDbClient(ctx context.Context, driver gira.DbDao, config *gira.Config) error {
	if client, err := db.NewConfigDbClient(ctx, gira.GAMEDB_NAME, *config.Db[gira.GAMEDB_NAME]); err != nil {
		return err
	} else if err := driver.UseClient(client); err != nil {
		return err
	} else {
		return nil
	}
}

func UseResourceDbClient(ctx context.Context, driver gira.DbDao, config *gira.Config) error {
	if client, err := db.NewConfigDbClient(ctx, gira.RESOURCEDB_NAME, *config.Db[gira.RESOURCEDB_NAME]); err != nil {
		return err
	} else if err := driver.UseClient(client); err != nil {
		return err
	} else {
		return nil
	}
}

func NewResourceDbClient(ctx context.Context, config *gira.Config) (gira.DbClient, error) {
	if client, err := db.NewConfigDbClient(ctx, gira.RESOURCEDB_NAME, *config.Db[gira.RESOURCEDB_NAME]); err != nil {
		return nil, err
	} else {
		return client, nil
	}
}
