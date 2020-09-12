// Package xogen contains the types for schema 'public'.
package xogen

// Code generated by xo. DO NOT EDIT.

// GenSalt calls the stored procedure 'public.gen_salt(text, text, integer) text' on db.
func GenSalt(db XODB, v0 string, v1 string, v2 int) (string, error) {
	var err error

	// sql query
	const sqlstr = `SELECT public.gen_salt($1, $2, $3)`

	// run query
	var ret string
	XOLog(sqlstr, v0, v1, v2)
	err = db.QueryRow(sqlstr, v0, v1, v2).Scan(&ret)
	if err != nil {
		return "", err
	}

	return ret, nil
}