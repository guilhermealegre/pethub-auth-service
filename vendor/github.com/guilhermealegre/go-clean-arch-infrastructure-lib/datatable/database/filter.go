package database

import (
	"context"
	"fmt"

	"github.com/guilhermealegre/go-clean-arch-core-lib/database/session"
)

func (f *FilterDatabase) getFacetData(client session.ISession, ctx context.Context, bindObject string) (FacetDatabaseList, error) {
	var err error
	facet := f.loadCachedFacet(bindObject)
	if len(facet) == 0 {
		facet, err = f.loadFacetFromDb(client, ctx)
		if err != nil {
			return nil, err
		}

		f.saveCacheFacet(facet, bindObject)
	}

	return facet, nil
}

func (f *FilterDatabase) loadCachedFacet(bindObject string) FacetDatabaseList {
	if _, ok := cachedFacetDatabase[f.FilterName]; ok {
		if _, ok := cachedFacetDatabase[f.FilterName][bindObject]; ok {
			return cachedFacetDatabase[f.FilterName][bindObject]
		}
	}
	return nil
}

func (f *FilterDatabase) saveCacheFacet(facet FacetDatabaseList, bindObject string) {
	if cachedFacetDatabase == nil {
		cachedFacetDatabase = make(map[string]map[string]FacetDatabaseList)
	}

	if cachedFacetDatabase[f.FilterName] == nil {
		cachedFacetDatabase[f.FilterName] = make(map[string]FacetDatabaseList)
	}

	if cachedFacetDatabase[f.FilterName][bindObject] == nil {
		cachedFacetDatabase[f.FilterName][bindObject] = FacetDatabaseList{}
	}

	cachedFacetDatabase[f.FilterName][bindObject] = facet
}

func (f *FilterDatabase) loadFacetFromDb(client session.ISession, ctx context.Context) (FacetDatabaseList, error) {
	facet := FacetDatabaseList{}
	builder := client.
		Select(f.ShowFields...).
		From(f.Table)

	if len(f.Join) > 0 {
		for _, joinTable := range f.Join {
			primaryKey := fmt.Sprintf("id_%s", joinTable.Table)
			foreignKey := fmt.Sprintf("fk_%s", joinTable.Table)
			if joinTable.PrimaryKey != "" {
				primaryKey = joinTable.PrimaryKey
			}
			if joinTable.ForeignKey != "" {
				foreignKey = joinTable.ForeignKey
			}
			if joinTable.LeftJoin {
				builder.LeftJoin(joinTable.Table, fmt.Sprintf("%s = %s", primaryKey, foreignKey))
			} else {
				builder.Join(joinTable.Table, fmt.Sprintf("%s = %s", primaryKey, foreignKey))
			}
		}
	}

	if len(f.Exceptions) > 0 {
		for _, exception := range f.Exceptions {
			builder.Where(exception)
		}
	}

	if f.GroupValue != "" {
		builder.GroupBy(f.GroupValue)
	}

	if f.OrderValue != "" {
		builder.OrderBy(f.OrderValue)
	}

	if _, err := builder.OrderBy("name").LoadContext(ctx, &facet); err != nil {
		return nil, err
	}

	return facet, nil
}
