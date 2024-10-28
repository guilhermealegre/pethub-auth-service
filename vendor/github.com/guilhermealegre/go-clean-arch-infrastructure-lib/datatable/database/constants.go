package database

const (
	Inner joinType = iota
	Left
	Right
	Full

	MaxPageResultLimit      = 48
	StartParameter          = "start"
	LengthParameter         = "length"
	SearchParameter         = "search"
	MasterParameter         = "master"
	OrderColumnParameter    = "order[column]"
	OrderDirectionParameter = "order[dir]"
	DrawParameter           = "draw"
	LimitParameter          = "limit" //Only with SearchWithPagination
	PageParameter           = "page"  //Only with SearchWithPagination

	OrderDirectionAsc  OrderDirection = "asc"
	OrderDirectionDesc OrderDirection = "desc"
)
