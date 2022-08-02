package main

import (
	"context"
	"io"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/entity"
)

type table struct {
	adapter *Database
	name    string
	schema  sql.Schema
	indexes []sql.Index
	entity  entity.EntityStoreClient
}

var (
	_ sql.Table = (*table)(nil)
)

func newTable(ctx *sql.Context, d *Database, tblName string, client entity.EntityStoreClient) (sql.Table, error) {
	t := &table{
		adapter: d,
		name:    tblName,
		entity:  client,
	}

	t.schema = append(t.schema, &sql.Column{
		Name:     "payload",
		Type:     sql.JSON,
		Nullable: true,
		Source:   tblName,
	})

	return t, nil
}

func (t *table) Name() string {
	return t.name
}

func (t *table) String() string {
	return t.name
}

func (t *table) Schema() sql.Schema {
	return t.schema
}

func (t *table) Partitions(*sql.Context) (sql.PartitionIter, error) {
	return &noopPartitionIter{}, nil
}

func (t *table) PartitionRows(ctx *sql.Context, _ sql.Partition) (sql.RowIter, error) {
	res, err := t.entity.ListFolder(context.Background(), &entity.ListFolderRequest{
		Path:        "devenv/dev-dashboards/datasource-postgres",
		WithPayload: true,
	})
	if err != nil {
		return nil, err
	}

	currentIdx := 0
	return &rowIter{table: t, res: res, currentIdx: &currentIdx}, nil

}

type noopPartitionIter struct {
	done bool
}

func (i *noopPartitionIter) Next(*sql.Context) (sql.Partition, error) {
	if !i.done {
		i.done = true
		return noopParition, nil
	}
	return nil, io.EOF
}

func (i *noopPartitionIter) Close(*sql.Context) error {
	return nil
}

var noopParition = partition(nil)

type partition []byte

func (p partition) Key() []byte {
	return p
}

type rowIter struct {
	table      *table
	res        *entity.FolderListing
	currentIdx *int
}

func (i rowIter) Next(ctx *sql.Context) (sql.Row, error) {
	values := make([]interface{}, len(i.table.schema))

	if *i.currentIdx >= len(i.res.GetItems()) {
		return nil, io.EOF
	}
	values[0] = i.res.GetItems()[*i.currentIdx].Payload
	*i.currentIdx += 1
	return sql.Row(values), nil
}

func (i rowIter) Close(_ *sql.Context) error {
	return nil
}
