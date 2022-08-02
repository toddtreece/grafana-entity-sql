package main

import (
	"context"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/entity"
)

type Database struct {
	tableNames []string
	client     entity.EntityStoreClient
}

func NewDatabase(entityClient entity.EntityStoreClient) *Database {
	return &Database{tableNames: []string{"dashboards"}, client: entityClient}
}

func (d *Database) Name() string {
	return "grafana"
}

func (d *Database) GetTableInsensitive(ctx *sql.Context, tblName string) (sql.Table, bool, error) {
	t, err := newTable(ctx, d, tblName, d.client)
	return t, true, err
}

func (d *Database) GetTableNames(ctx *sql.Context) ([]string, error) {
	return d.tableNames, nil
}

type tokenAuth struct {
	token string
}

func (t tokenAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + t.token,
	}, nil
}

func (tokenAuth) RequireTransportSecurity() bool {
	return false
}
