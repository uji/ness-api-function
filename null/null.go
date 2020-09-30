package null

import "database/sql"

type (
	Bool   = sql.NullBool
	Int32  = sql.NullInt32
	Int64  = sql.NullInt64
	Float  = sql.NullFloat64
	String = sql.NullString
	Time   = sql.NullTime
)
