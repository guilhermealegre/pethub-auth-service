package database

import "strings"

func (s *engine) addTableByField(field string) *engine {
	if s.tablesActivatedByFilters == nil {
		s.tablesActivatedByFilters = make(map[string]interface{})
	}
	t := strings.Split(field, ".")
	if _, ok := s.tablesActivatedByFilters[t[0]]; !ok {
		s.tablesActivatedByFilters[t[0]] = t[1]
	}

	return s
}

func (s *engine) getDependentTables() map[string]interface{} {
	tmp := make(map[string]interface{}, len(s.tablesActivatedByFilters))
	for table := range s.tablesActivatedByFilters {
		if table == s.primaryTable {
			continue
		}

		if _, ok := tmp[table]; !ok {
			tmp[table] = nil
		}
	}

	return tmp
}
