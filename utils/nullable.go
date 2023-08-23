package utils

import "database/sql"

/*
If string is null, returns a null string else a valid string
*/
func NewNullableString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
