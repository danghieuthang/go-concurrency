package concurrency

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// Version is a wrapper around sql.NullString to handle nullable string fields in a database.
type Version sql.NullString

// NewVersion generates a new Version with a unique string value using ksuid.New().String().
func NewVersion() Version {
	return Version{
		Valid:  true,
		String: ksuid.New().String(),
	}
}

// Scan implements the Scanner interface. It is used to convert the SQL value into the Version type.
func (v *Version) Scan(value interface{}) error {
	return (*sql.NullString)(v).Scan(value)
}

// Value implements the driver Valuer interface. It is used to convert the Version type into a SQL value.
func (v Version) Value() (driver.Value, error) {
	if !v.Valid {
		return nil, nil
	}
	return v.String, nil
}

// UnmarshalJSON is used to convert the JSON value into the Version type.
func (v *Version) UnmarshalJSON(bytes []byte) error {
	if string(bytes) == "null" {
		v.Valid = false
		return nil
	}
	err := json.Unmarshal(bytes, &v.String)
	if err == nil {
		v.Valid = true
	}
	return err
}

// MarshalJSON is used to convert the Version type into a JSON value.
func (v *Version) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	}
	return json.Marshal(nil)
}

// CreateClauses returns a slice of clause.Interface to customize the SQL clauses generated by GORM when creating a record.
func (v *Version) CreateClauses(field *schema.Field) []clause.Interface {
	return []clause.Interface{VersionCreateClause{Field: field}}
}

// VersionCreateClause is used to provide custom behavior for creating records with a Version field.
type VersionCreateClause struct {
	Field *schema.Field
}

func (v VersionCreateClause) Name() string {
	return ""
}

func (v VersionCreateClause) Build(builder clause.Builder) {

}

func (v VersionCreateClause) MergeClause(c *clause.Clause) {

}

// ModifyStatement is used to modify the SQL statement for creating a record.
func (v VersionCreateClause) ModifyStatement(stmt *gorm.Statement) {
	if stmt.SQL.Len() == 0 {
		nv := NewVersion()
		stmt.AddClause(clause.Set{{Column: clause.Column{Name: v.Field.DBName}, Value: nv.String}})
		stmt.SetColumn(v.Field.DBName, nv.String, true)
	}
}

// UpdateClauses returns a slice of clause.Interface to customize the SQL clauses generated by GORM when updating a record.
func (v *Version) UpdateClauses(field *schema.Field) []clause.Interface {
	return []clause.Interface{VersionUpdateClause{Field: field}}
}

// VersionUpdateClause is used to provide custom behavior for updating records with a Version field.
type VersionUpdateClause struct {
	Field *schema.Field
}

func (v VersionUpdateClause) Name() string {
	return ""
}

func (v VersionUpdateClause) Build(builder clause.Builder) {

}

func (v VersionUpdateClause) MergeClause(c *clause.Clause) {

}

// ModifyStatement is used to modify the SQL statement for updating a record.
func (v VersionUpdateClause) ModifyStatement(stmt *gorm.Statement) {
	if stmt.SQL.Len() == 0 {
		if _, ok := stmt.Clauses["concurrency_query"]; !ok {
			if cv, zero := v.Field.ValueOf(stmt.Context, stmt.ReflectValue); !zero {
				if cvv, ok := cv.(Version); ok && cvv.Valid {
					stmt.AddClause(clause.Where{Exprs: []clause.Expression{
						clause.Eq{Column: clause.Column{Table: clause.CurrentTable, Name: v.Field.DBName}, Value: cvv.String},
					}})
				}
			}
			stmt.Clauses["concurrency_query"] = clause.Clause{}
		}

		nv := NewVersion()
		stmt.SetColumn(v.Field.DBName, nv.String, true)
	}
}
