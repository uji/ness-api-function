// Package xogen contains the types for schema 'public'.
package xogen

// Code generated by xo. DO NOT EDIT.

// Decrypt calls the stored procedure 'public.decrypt(bytea, bytea, text) bytea' on db.
func Decrypt(db XODB, v0 []byte, v1 []byte, v2 string) ([]byte, error) {
	var err error

	// sql query
	const sqlstr = `SELECT public.decrypt($1, $2, $3)`

	// run query
	var ret []byte
	XOLog(sqlstr, v0, v1, v2)
	err = db.QueryRow(sqlstr, v0, v1, v2).Scan(&ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
