package database

import (
	"context"
	"net/url"

	"github.com/guilhermealegre/go-clean-arch-core-lib/pagination"
)

type IDatabase interface {
	IsMaster(bool) IDatabase
	UrlValues(values url.Values) IDatabase
	Context(ctx context.Context) IDatabase
	Distinct() IDatabase
	Sort(column ...string) IDatabase
	SortDir(direction OrderDirection) IDatabase
	Searchable(column ...string) IDatabase
	Group(column string) IDatabase
	Filter(filter ...FilterDatabase) IDatabase
	Bind(obj interface{}) IDatabase
	QueryBuilder(obj SearchQueryBuilder) IDatabase
	SearchWithPagination() (data interface{}, result *SearchResult, pag *pagination.Pagination, err error)
	Search() (data interface{}, result *SearchResult, err error)
	Count() (count int64, err error)
	Facets() (FacetDatabaseMap, error)
	Path(path string) IDatabase
}
