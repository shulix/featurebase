// Copyright 2021 Molecula Corp. All rights reserved.

package planner

import (
	"context"
	"fmt"
	"time"

	pilosa "github.com/featurebasedb/featurebase/v3"
	"github.com/featurebasedb/featurebase/v3/sql3/parser"
	"github.com/featurebasedb/featurebase/v3/sql3/planner/types"
)

// PlanOpFeatureBaseTables wraps a []*IndexInfo that is returned from
// schemaAPI.Schema().
type PlanOpFeatureBaseTables struct {
	indexInfo []*pilosa.IndexInfo
	warnings  []string
}

func NewPlanOpFeatureBaseTables(indexInfo []*pilosa.IndexInfo) *PlanOpFeatureBaseTables {
	return &PlanOpFeatureBaseTables{
		indexInfo: indexInfo,
		warnings:  make([]string, 0),
	}
}

func (p *PlanOpFeatureBaseTables) Plan() map[string]interface{} {
	result := make(map[string]interface{})
	result["_op"] = fmt.Sprintf("%T", p)
	result["_schema"] = p.Schema().Plan()
	return result
}

func (p *PlanOpFeatureBaseTables) String() string {
	return ""
}

func (p *PlanOpFeatureBaseTables) AddWarning(warning string) {
	p.warnings = append(p.warnings, warning)
}

func (p *PlanOpFeatureBaseTables) Warnings() []string {
	return p.warnings
}

func (p *PlanOpFeatureBaseTables) Schema() types.Schema {
	return types.Schema{
		&types.PlannerColumn{
			RelationName: "fb_tables",
			ColumnName:   "_id",
			Type:         parser.NewDataTypeString(),
		},
		&types.PlannerColumn{
			RelationName: "fb_tables",
			ColumnName:   "name",
			Type:         parser.NewDataTypeString(),
		},
		&types.PlannerColumn{
			RelationName: "fb_tables",
			ColumnName:   "owner",
			Type:         parser.NewDataTypeString(),
		},
		&types.PlannerColumn{
			RelationName: "fb_tables",
			ColumnName:   "updated_by",
			Type:         parser.NewDataTypeString(),
		},
		&types.PlannerColumn{
			RelationName: "fb_tables",
			ColumnName:   "created_at",
			Type:         parser.NewDataTypeTimestamp(),
		},
		&types.PlannerColumn{
			RelationName: "fb_tables",
			ColumnName:   "updated_at",
			Type:         parser.NewDataTypeTimestamp(),
		},
		&types.PlannerColumn{
			RelationName: "fb_tables",
			ColumnName:   "keys",
			Type:         parser.NewDataTypeBool(),
		},
		&types.PlannerColumn{
			RelationName: "fb_tables",
			ColumnName:   "description",
			Type:         parser.NewDataTypeString(),
		},
	}
}

func (p *PlanOpFeatureBaseTables) Children() []types.PlanOperator {
	return []types.PlanOperator{}
}

func (p *PlanOpFeatureBaseTables) Iterator(ctx context.Context, row types.Row) (types.RowIterator, error) {
	return &showTablesRowIter{
		indexInfo: p.indexInfo,
	}, nil
}

func (p *PlanOpFeatureBaseTables) WithChildren(children ...types.PlanOperator) (types.PlanOperator, error) {
	return NewPlanOpFeatureBaseTables(p.indexInfo), nil
}

type showTablesRowIter struct {
	indexInfo []*pilosa.IndexInfo
	rowIndex  int
}

var _ types.RowIterator = (*showTablesRowIter)(nil)

func (i *showTablesRowIter) Next(ctx context.Context) (types.Row, error) {
	if i.rowIndex < len(i.indexInfo) {
		createdAt := time.Unix(0, i.indexInfo[i.rowIndex].CreatedAt)
		updatedAt := time.Unix(0, i.indexInfo[i.rowIndex].UpdatedAt)
		row := []interface{}{
			i.indexInfo[i.rowIndex].Name,
			i.indexInfo[i.rowIndex].Name,
			i.indexInfo[i.rowIndex].Owner,
			i.indexInfo[i.rowIndex].LastUpdateUser,
			createdAt.Format(time.RFC3339),
			updatedAt.Format(time.RFC3339),
			i.indexInfo[i.rowIndex].Options.Keys,
			i.indexInfo[i.rowIndex].Options.Description,
		}
		i.rowIndex += 1
		return row, nil
	}
	return nil, types.ErrNoMoreRows
}
