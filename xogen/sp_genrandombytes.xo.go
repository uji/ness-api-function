// Package xogen contains the types for schema 'public'.
package xogen

// Code generated by xo. DO NOT EDIT.

// GenRandomBytes calls the stored procedure 'public.gen_random_bytes(integer) bytea' on db.
func GenRandomBytes(db XODB, v0 int) ([]byte, error) {
	var err error

	// sql query
	const sqlstr = `SELECT public.gen_random_bytes($1)`

	// run query
	var ret []byte
	XOLog(sqlstr, v0)
	err = db.QueryRow(sqlstr, v0).Scan(&ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
