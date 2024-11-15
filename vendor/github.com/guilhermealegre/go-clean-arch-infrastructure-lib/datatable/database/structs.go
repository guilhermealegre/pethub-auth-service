package database

import (
	"context"
	"net/url"

	"github.com/gocraft/dbr/v2"

	"github.com/guilhermealegre/go-clean-arch-core-lib/database/session"
)

type DatabaseSearch struct {
	values             url.Values
	sortColumns        []string
	sortDirection      OrderDirection
	searchableColumns  []string
	groupColumn        string
	filters            FilterDatabaseList
	bindObject         interface{}
	logger             logFunction
	client             Client
	builder            *dbr.SelectBuilder
	maxPageResultLimit int
	context            context.Context
	distinct           bool
	path               string
}

type Client struct {
	Reader session.ISession
	Writer session.ISession
}

type joinType uint8
type OrderDirection string
type logFunction func(err error) error

type SearchResult struct {
	Facets          interface{} `json:"facets"`
	RecordsTotal    int64       `json:"recordsTotal"`
	RecordsFiltered int         `json:"recordsFiltered"`
}

type FilterDatabaseList []FilterDatabase
type FilterDatabase struct {
	FilterName        string                 // [required] filter name
	Table             interface{}            // [required] table name
	ShowFields        []any                  // [required] fields to show on the filter list
	FilterField       string                 // [required] filter by column name
	Join              FilterJoinDatabaseList // [optional] join with other table
	NullLabel         string                 // [optional] set null label default value as None
	NullValue         bool                   // [optional] set filter default value as null
	Exceptions        []interface{}          // [optional] exceptions, e.g: ignore specific record
	StringId          bool                   // [optional] when filter is made by a string
	ArrayId           bool                   // [optional] when filter is made by an array
	DisabledFacetLoad bool
	ExtraFacet        ExtraFacetDatabase // [optional] Add Extra facet in case you need to add more specific filters
	SpecialFilter     SpecialFiltersDatabaseList
	GroupValue        string // [optional] Case you need to group by value
	OrderValue        string // [optional] Case you need to order by for specific value
	SkipDbFacetLoad   bool   //[optional] Case you only want static (extra) facets without db load
}

type FilterJoinDatabaseList []FilterJoinDatabase
type FilterJoinDatabase struct {
	LeftJoin   bool
	Table      string
	ForeignKey string
	PrimaryKey string
}

type ExtraFacetDatabase struct {
	Active bool // define if you need extra facets
	Facets FacetDatabaseList
}

type FacetDatabaseList []FacetDatabase
type FacetDatabaseMap map[string]FacetDatabaseList
type FacetDatabase struct {
	Id   interface{} `json:"id"`
	Key  string      `json:"key,omitempty"`
	Name string      `json:"name,omitempty"`
}

type SpecialFiltersDatabaseList []SpecialFiltersDatabase
type SpecialFiltersDatabase struct {
	Value     interface{} // Value to compare with the value given in the input
	Condition string      // Condition to filter the value given in the input
}

type SearchQueryBuilder struct {
	Columns []interface{}
	From    SearchFromBuilder
	Join    []SearchJoinBuilder
	Where   []SearchWhereBuilder
	Order   []SearchOrderBuilder
	GroupBy []string
}

type SearchFromBuilder struct {
	Builder dbr.Builder
	Alias   string
}

type SearchJoinBuilder struct {
	JoinType         joinType
	Table            interface{}
	On               interface{}
	Alias            string
	RequiredForCount bool
}

type SearchWhereBuilder struct {
	Query interface{}
	Value []interface{}
}

type SearchOrderBuilder struct {
	Field     string
	OrderDesc bool
}

var cachedFacetDatabase map[string]map[string]FacetDatabaseList

type SearchParams struct {
	Start        int
	Length       int
	Search       string
	OrderColumns int
	OrderDir     string
}

type engine struct {
	isMaster                 bool
	primaryTable             string
	tablesActivatedByFilters map[string]interface{}
}
