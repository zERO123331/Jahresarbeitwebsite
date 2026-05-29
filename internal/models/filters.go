package models

import "Jahresarbeitwebsite/internal/validator"

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.CheckFilterErrors(f.Page > 0, "page", "must be greater than zero")
	v.CheckFilterErrors(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.CheckFilterErrors(f.PageSize > 0, "page_size", "must be greater than zero")
	v.CheckFilterErrors(f.PageSize <= 100, "page_size", "must be a maximum of 100")

	v.CheckFilterErrors(validator.PermittedValue(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}
