package database

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/gocraft/dbr/v2"

	"github.com/guilhermealegre/go-clean-arch-core-lib/pagination"

	"github.com/gocraft/dbr/v2/dialect"
	"github.com/guilhermealegre/go-clean-arch-core-lib/database/session"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/errors"
)

type Database struct {
	Params SearchParams
	DatabaseSearch
	SearchQueryBuilder
	engine
}

func New(client Client, log logFunction, maxResults int) IDatabase {
	return &Database{
		DatabaseSearch: DatabaseSearch{
			logger:             log,
			client:             client,
			maxPageResultLimit: maxResults,
		},
	}
}

func (d *Database) IsMaster(b bool) IDatabase {
	d.isMaster = b
	return d
}

func (d *Database) UrlValues(values url.Values) IDatabase {
	d.values = values
	return d
}

func (d *Database) Context(ctx context.Context) IDatabase {
	d.context = ctx
	return d
}

func (d *Database) Distinct() IDatabase {
	d.distinct = true
	return d
}

func (d *Database) Sort(column ...string) IDatabase {
	d.sortColumns = column
	return d
}

func (d *Database) SortDir(direction OrderDirection) IDatabase {
	d.sortDirection = direction
	return d
}

func (d *Database) Searchable(column ...string) IDatabase {
	d.searchableColumns = column
	return d
}

func (d *Database) Group(column string) IDatabase {
	d.groupColumn = column
	return d
}

func (d *Database) Filter(filter ...FilterDatabase) IDatabase {
	d.filters = filter
	return d
}

func (d *Database) Bind(obj interface{}) IDatabase {
	d.bindObject = obj
	return d
}

func (d *Database) QueryBuilder(obj SearchQueryBuilder) IDatabase {
	d.SearchQueryBuilder = obj
	return d
}

func (d *Database) SearchWithPagination() (data interface{}, result *SearchResult, pag *pagination.Pagination, err error) {
	//In the search, the limits and offset are miscalculated,
	var page, length int
	var isPage bool
	var isLength bool

	isLength = true
	isPage = true

	if value := d.context.Value(LimitParameter); value != nil {
		if v, ok := value.(int); ok {
			length = v
		}
	}

	if length == 0 {
		length, _ = strconv.Atoi(d.values.Get(LimitParameter))
		if length == 0 {
			isLength = false
			length = d.maxPageResultLimit
		}
	}

	if value := d.context.Value(PageParameter); value != nil {
		if v, ok := value.(int); ok {
			page = v
		}
	}

	if page == 0 {
		page, _ = strconv.Atoi(d.values.Get(PageParameter))
		if page == 0 {
			page = 1
			isPage = false
		}
	}

	if isPage || isLength {
		d.values.Set(LengthParameter, strconv.Itoa(length))
		d.values.Set(StartParameter, strconv.Itoa((page-1)*length))
	}

	data, result, err = d.Search()
	if err != nil {
		return nil, nil, nil, err
	}

	pag = &pagination.Pagination{
		Page:    uint64(page),
		Limit:   uint64(length),
		Records: uint64(result.RecordsTotal),
		Path:    d.path,
	}

	pag.Prepare()

	return data, result, pag, nil
}

func (d *Database) Search() (data interface{}, _ *SearchResult, err error) {
	result := &SearchResult{}

	if err = d.prepare(); err != nil {
		return nil, nil, err
	}

	// calculate countTotal
	countTotal, err := d.doCountTotal()
	if err != nil {
		return nil, nil, d.logger(err)
	}

	if len(strings.Split(d.groupColumn, ",")) > 1 {
		d.builder.GroupBy(d.groupColumn)
	}

	// apply parameters
	if len(d.sortColumns) > 0 {
		if OrderDirection(d.Params.OrderDir) == OrderDirectionDesc || (d.Params.OrderDir == "" && d.sortDirection == OrderDirectionDesc) {
			d.builder.OrderDir(d.sortColumns[d.Params.OrderColumns], false)
		} else {
			d.builder.OrderDir(d.sortColumns[d.Params.OrderColumns], true)
		}
	}

	// pagination
	if d.Params.Length > 0 {
		d.builder.Limit(uint64(d.Params.Length)).Offset(uint64(d.Params.Start))
	}

	var countFiltered int
	countFiltered, err = d.builder.LoadContext(d.context, d.bindObject)
	if err != nil {
		return nil, nil, d.logger(err)
	}

	// parse null string from database
	if _, err = json.Marshal(d.bindObject); err != nil {
		_ = d.logger(err)
	}

	// facets
	facets, err := d.Facets()
	if err != nil {
		return nil, nil, err
	}

	// result
	result.Facets = facets
	result.RecordsFiltered = countFiltered
	result.RecordsTotal = countTotal

	if countFiltered > 0 {
		data = d.bindObject
	}

	return data, result, err
}

func (d *Database) Facets() (_ FacetDatabaseMap, err error) {
	result := FacetDatabaseMap{}
	name := reflect.TypeOf(d.bindObject).Elem().Name()

	for _, filter := range d.filters {
		if !filter.DisabledFacetLoad {
			facet := make(FacetDatabaseList, 0)
			if !filter.SkipDbFacetLoad {
				facet, err = filter.getFacetData(d.getClient(), d.context, name)
				if err != nil {
					return nil, d.logger(err)
				}
			}

			if filter.NullValue {
				aux := FacetDatabaseList{}
				label := "None"
				if filter.NullLabel != "" {
					label = filter.NullLabel
				}

				aux = append(aux, FacetDatabase{Id: 0, Key: label, Name: label})
				facet = append(aux, facet...)
			}

			if filter.ExtraFacet.Active {
				facet = append(filter.ExtraFacet.Facets, facet...)
			}

			result[filter.FilterName] = facet
		}
	}

	return result, nil
}

func (d *Database) Path(path string) IDatabase {
	d.path = path
	return d
}

func (d *Database) doCountTotal() (int64, error) {
	var countTotal int64
	var counterBuilder = d.getClient().Select("*")

	groupColumn := d.groupColumn

	if groupColumn == "" {
		groupColumn = "*"
	} else {
		counterBuilder.Group = d.builder.Group
		if len(strings.Split(groupColumn, ",")) > 1 {
			groupColumn = "*"
		} else {
			d.engine.addTableByField(groupColumn)
			groupColumn = "DISTINCT (" + groupColumn + ")"
			counterBuilder.Group = nil
		}
	}

	counterBuilder.Column = []interface{}{
		"COUNT(" + groupColumn + ") AS countTotal",
	}

	d.prepareCountBuilder(counterBuilder)

	counterBuilder.WhereCond = d.builder.WhereCond

	_, err := counterBuilder.LoadContext(d.context, &countTotal)
	if err != nil {
		return 0, d.logger(err)
	}

	return countTotal, nil
}

func (d *Database) Count() (int64, error) {
	if err := d.prepare(); err != nil {
		return 0, err
	}

	count, err := d.doCountTotal()
	if err != nil {
		return 0, d.logger(err)
	}

	return count, nil
}

func (d *Database) prepare() error {
	if !d.validateSearchableColumns() {
		return errors.ErrorDatatableSearchableColumn()
	}
	if !d.validateFilters() {
		return errors.ErrorDatatableFields()
	}
	if !d.validateGroupColumn() {
		return errors.ErrorDatatableGroup()
	}
	d.initializeEngine().
		prepareContext().
		prepareParams().
		prepareBuilder().
		prepareSearch().
		prepareFilters()
	return nil
}

func (d *Database) validateSearchableColumns() bool {
	for _, column := range d.searchableColumns {
		if len(strings.Split(column, ".")) == 0 {
			return false
		}
	}

	return true
}

func (d *Database) validateFilters() bool {
	for _, filter := range d.filters {
		if len(strings.Split(filter.FilterField, ".")) == 0 {
			return false
		}
	}

	return true
}

func (d *Database) validateGroupColumn() bool {
	return len(strings.Split(d.groupColumn, ".")) != 0
}

func (d *Database) initializeEngine() *Database {
	d.engine.primaryTable = d.From.Alias

	return d
}

func (d *Database) prepareContext() *Database {
	if d.context == nil {
		d.context = context.Background()
	}

	return d
}

func (d *Database) prepareBuilder() *Database {
	d.builder = d.getClient().Select().
		From(d.From.Builder)

	if d.distinct {
		d.builder.Distinct()
	}

	d.builder.Column = d.Columns

	for _, join := range d.Join {
		switch join.JoinType {
		case Inner:
			d.builder.Join(join.Table, join.On)
		case Left:
			d.builder.LeftJoin(join.Table, join.On)
		case Right:
			d.builder.RightJoin(join.Table, join.On)
		case Full:
			d.builder.FullJoin(join.Table, join.On)
		}
	}

	for _, condition := range d.Where {
		if condition.Value == nil {
			d.builder.Where(condition.Query)
		} else {
			d.builder.Where(condition.Query, condition.Value...)
		}
	}

	for _, order := range d.Order {
		d.builder.OrderDir(order.Field, !order.OrderDesc)
	}

	if len(d.GroupBy) > 0 {
		d.builder.GroupBy(d.GroupBy...)
	}

	return d
}

func (d *Database) prepareCountBuilder(counterBuilder *dbr.SelectStmt) *dbr.SelectStmt {
	counterBuilder.Table = d.From.Builder

	dependentTables := d.engine.getDependentTables()
	for _, join := range d.Join {
		if join.RequiredForCount || len(dependentTables) > 0 {
			switch join.JoinType {
			case Inner:
				counterBuilder.Join(join.Table, join.On)
			case Left:
				counterBuilder.LeftJoin(join.Table, join.On)
			case Right:
				counterBuilder.RightJoin(join.Table, join.On)
			case Full:
				counterBuilder.FullJoin(join.Table, join.On)
			}

			delete(dependentTables, join.Alias)
		}
	}

	return counterBuilder
}

func (d *Database) prepareParams() *Database {
	d.Params.Start = getIntParameter(StartParameter, d.values, d.context)
	d.Params.Length = getIntParameter(LengthParameter, d.values, d.context)
	d.Params.Search = strings.Replace(html.UnescapeString(getStringParameter(SearchParameter, d.values, d.context)), "?", "??", -1)
	d.Params.OrderColumns = getIntParameter(OrderColumnParameter, d.values, d.context)
	d.Params.OrderDir = getStringParameter(OrderDirectionParameter, d.values, d.context)

	// set minimum to show up
	if d.Params.Length < 0 {
		d.Params.Length = d.maxPageResultLimit
	}

	return d
}

func (d *Database) prepareSearch() *Database {
	// search
	var whereStr string

	if len(d.Params.Search) > 0 {
		for k, v := range d.searchableColumns {
			d.engine.addTableByField(v)
			if k != 0 {
				whereStr += " OR "
			}
			if v == "created_at" || v == "updated_at" {
				whereStr += " to_char(" + v + ", 'yyyy-MM-dd HH:min:ss') LIKE " + dialect.PostgreSQL.EncodeString(d.Params.Search+"%") //like
			} else {
				whereStr += v + " ILIKE " + dialect.PostgreSQL.EncodeString("%"+strings.Replace(d.Params.Search, " ", "%", -1)+"%") //like
			}
		}

		d.builder.Where(whereStr)
	}

	return d
}

func (d *Database) prepareFilters() *Database {
	// filters
	var value string
	for _, filter := range d.filters {
		value = d.values.Get(filter.FilterName)
		if value != "" {
			filterVal, _ := strconv.Atoi(value)

			var foundSpecialFilter bool

			if len(filter.SpecialFilter) > 0 {
				for _, specialFilter := range filter.SpecialFilter {
					if filterVal == specialFilter.Value {
						d.builder.Where(specialFilter.Condition)
						foundSpecialFilter = true
						break
					}
				}
			}

			if !foundSpecialFilter {
				d.engine.addTableByField(filter.FilterField)
				if filter.StringId {
					d.builder.Where(fmt.Sprintf("%s = ?", filter.FilterField), value)
				} else if filter.ArrayId {
					value = strings.NewReplacer("[", "", "]", "").Replace(value)
					values := strings.Split(value, ",")
					d.builder.Where(fmt.Sprintf("%s IN ?", filter.FilterField), values)
				} else {
					filterVal, err := strconv.Atoi(value)
					if err == nil {
						if filterVal > 0 {
							d.builder.Where(fmt.Sprintf("%s = ?", filter.FilterField), filterVal)
						} else {
							d.builder.Where(fmt.Sprintf("%s IS NULL", filter.FilterField))
						}
					} else {
						boolVal, err := strconv.ParseBool(value)
						if err == nil {
							d.builder.Where(fmt.Sprintf("%s = ?", filter.FilterField), boolVal)
						}
					}
				}
			}
		}
	}

	return d
}

func (d *Database) getClient() session.ISession {
	if d.isMaster {
		return d.client.Writer
	}

	return d.client.Reader
}
