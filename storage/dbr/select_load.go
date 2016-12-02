package dbr

import (
	"database/sql"
	"reflect"

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/util/errors"
)

// Rows executes a query and returns the row result.
func (b *SelectBuilder) Rows(hss ...SelectHook) (*sql.Rows, error) {

	SelectHooks(hss).Apply(b)

	sqlStr, args, err := b.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "[store] SelectBuilder.Rows.ToSql")
	}
	rows, err := b.Querier.Query(sqlStr, args...)

	return rows, errors.Wrap(err, "[store] SelectBuilder.Rows.QueryContext")
}

// Unvetted thots:
// Given a query and given a structure (field list), there's 2 sets of fields.
// Take the intersection. We can fill those in. great.
// For fields in the structure that aren't in the query, we'll let that slide if db:"-"
// For fields in the structure that aren't in the query but without db:"-", return error
// For fields in the query that aren't in the structure, we'll ignore them.

// LoadStructs executes the SelectBuilder and loads the resulting data into a
// slice of structs dest must be a pointer to a slice of pointers to structs.
// Returns the number of items found (which is not necessarily the # of items
// set). Slow because of the massive use of reflection.
func (b *SelectBuilder) LoadStructs(dest interface{}) (int, error) {
	//
	// Validate the dest, and extract the reflection values we need.
	//

	// This must be a pointer to a slice
	valueOfDest := reflect.ValueOf(dest)
	kindOfDest := valueOfDest.Kind()

	if kindOfDest != reflect.Ptr {
		return 0, errors.NewNotValidf("[dbr] invalid type passed to LoadStructs. Need a pointer to a slice")
	}

	// This must a slice
	valueOfDest = reflect.Indirect(valueOfDest)
	kindOfDest = valueOfDest.Kind()

	if kindOfDest != reflect.Slice {
		return 0, errors.NewNotValidf("[dbr] invalid type passed to LoadStructs. Need a pointer to a slice")
	}

	// The slice elements must be pointers to structures
	recordType := valueOfDest.Type().Elem()
	if recordType.Kind() != reflect.Ptr {
		return 0, errors.NewNotValidf("[dbr] Elements need to be pointers to structures")
	}

	recordType = recordType.Elem()
	if recordType.Kind() != reflect.Struct {
		return 0, errors.NewNotValidf("[dbr] Elements need to be pointers to structures")
	}

	//
	// Get full SQL
	//
	tSQL, tArg, err := b.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "[dbr] select.load_structs.tosql")
	}

	fullSql, err := Preprocess(tSQL, tArg)
	if err != nil {
		return 0, errors.Wrap(err, "[dbr] select.load_all.interpolate")
	}

	numberOfRowsReturned := 0

	if b.Logger.IsInfo() {
		defer log.WhenDone(b.Logger).Info("dbr.SelectBuilder.LoadStructs.QueryContext.timing", log.String("sql", fullSql))
	}

	// Run the query:
	rows, err := b.Query(fullSql)
	if err != nil {
		return 0, errors.Wrap(err, "[dbr] select.load_all.query")
	}
	defer rows.Close()

	// Get the columns returned
	columns, err := rows.Columns()
	if err != nil {
		return numberOfRowsReturned, errors.Wrap(err, "[dbr] select.load_one.rows.Columns")
	}

	// Create a map of this result set to the struct fields
	fieldMap, err := calculateFieldMap(recordType, columns, false)
	if err != nil {
		return numberOfRowsReturned, errors.Wrap(err, "[dbr] select.load_all.calculateFieldMap")
	}

	// Build a 'holder', which is an []interface{}. Each value will be the set to address of the field corresponding to our newly made records:
	holder := make([]interface{}, len(fieldMap))

	// Iterate over rows and scan their data into the structs
	sliceValue := valueOfDest
	for rows.Next() {
		// Create a new record to store our row:
		pointerToNewRecord := reflect.New(recordType)
		newRecord := reflect.Indirect(pointerToNewRecord)

		// Prepare the holder for this record
		scannable, err := prepareHolderFor(newRecord, fieldMap, holder)
		if err != nil {
			return numberOfRowsReturned, errors.Wrap(err, "[dbr] select.load_all.holderFor")
		}

		// Load up our new structure with the row's values
		err = rows.Scan(scannable...)
		if err != nil {
			return numberOfRowsReturned, errors.Wrap(err, "[dbr] select.load_all.scan")
		}

		// Append our new record to the slice:
		sliceValue = reflect.Append(sliceValue, pointerToNewRecord)

		numberOfRowsReturned++
	}
	valueOfDest.Set(sliceValue)

	// Check for errors at the end. Supposedly these are error that can happen during iteration.
	if err = rows.Err(); err != nil {
		return numberOfRowsReturned, errors.Wrap(err, "[dbr] select.load_all.rows_err")
	}

	return numberOfRowsReturned, nil
}

// LoadStruct executes the SelectBuilder and loads the resulting data into a
// struct dest must be a pointer to a struct Returns ErrNotFound behaviour. Slow
// because of the massive use of reflection.
func (b *SelectBuilder) LoadStruct(dest interface{}) error {
	//
	// Validate the dest, and extract the reflection values we need.
	//
	valueOfDest := reflect.ValueOf(dest)
	indirectOfDest := reflect.Indirect(valueOfDest)
	kindOfDest := valueOfDest.Kind()

	if kindOfDest != reflect.Ptr || indirectOfDest.Kind() != reflect.Struct {
		return errors.NewNotValidf("[dbr] you need to pass in the address of a struct")
	}

	recordType := indirectOfDest.Type()

	//
	// Get full SQL
	//
	tSQL, tArg, err := b.ToSql()
	if err != nil {
		return errors.Wrap(err, "[dbr] select.load_struct.tosql")
	}

	fullSql, err := Preprocess(tSQL, tArg)
	if err != nil {
		return err
	}

	if b.Logger.IsInfo() {
		defer log.WhenDone(b.Logger).Info("dbr.SelectBuilder.LoadStruct.ExecContext.timing", log.String("sql", fullSql))
	}

	// Run the query:
	rows, err := b.Query(fullSql)
	if err != nil {
		return errors.Wrap(err, "[dbr] select.load_one.query")
	}
	defer rows.Close()

	// Get the columns of this result set
	columns, err := rows.Columns()
	if err != nil {
		return errors.Wrap(err, "[dbr] select.load_one.rows.Columns")
	}

	// Create a map of this result set to the struct columns
	fieldMap, err := calculateFieldMap(recordType, columns, false)
	if err != nil {
		return errors.Wrap(err, "[dbr] select.load_one.calculateFieldMap")
	}

	// Build a 'holder', which is an []interface{}. Each value will be the set to address of the field corresponding to our newly made records:
	holder := make([]interface{}, len(fieldMap))

	if rows.Next() {
		// Build a 'holder', which is an []interface{}. Each value will be the address of the field corresponding to our newly made record:
		scannable, err := prepareHolderFor(indirectOfDest, fieldMap, holder)
		if err != nil {
			return errors.Wrap(err, "[dbr] select.load_one.holderFor")
		}

		// Load up our new structure with the row's values
		err = rows.Scan(scannable...)
		if err != nil {
			return errors.Wrap(err, "[dbr] select.load_one.scan")
		}
		return nil
	}

	if err := rows.Err(); err != nil {
		return errors.Wrap(err, "[dbr] select.load_one.rows_err")
	}

	return errors.NewNotFoundf("[dbr] Entry not found")
}

// LoadValues executes the SelectBuilder and loads the resulting data into a
// slice of primitive values Returns ErrNotFound behaviour if no value was
// found, and it was therefore not set. Slow because of the massive use of
// reflection.
func (b *SelectBuilder) LoadValues(dest interface{}) (int, error) {
	// Validate the dest and reflection values we need

	// This must be a pointer to a slice
	valueOfDest := reflect.ValueOf(dest)
	kindOfDest := valueOfDest.Kind()

	if kindOfDest != reflect.Ptr {
		return 0, errors.NewNotValidf("[dbr] invalid type passed to LoadValues. Need a pointer to a slice")
	}

	// This must a slice
	valueOfDest = reflect.Indirect(valueOfDest)
	kindOfDest = valueOfDest.Kind()

	if kindOfDest != reflect.Slice {
		return 0, errors.NewNotValidf("[dbr] invalid type passed to LoadValues. Need a pointer to a slice")
	}

	recordType := valueOfDest.Type().Elem()

	recordTypeIsPtr := recordType.Kind() == reflect.Ptr
	if recordTypeIsPtr {
		reflect.ValueOf(dest)
	}

	//
	// Get full SQL
	//
	tSQL, tArg, err := b.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "[dbr] select.load_values.tosql")
	}

	fullSql, err := Preprocess(tSQL, tArg)
	if err != nil {
		return 0, err
	}

	numberOfRowsReturned := 0

	if b.Logger.IsInfo() {
		defer log.WhenDone(b.Logger).Info("dbr.SelectBuilder.LoadValues.QueryContext.timing", log.String("sql", fullSql))
	}

	// Run the query:
	rows, err := b.Query(fullSql)
	if err != nil {
		return numberOfRowsReturned, errors.Wrap(err, "[dbr] select.load_all_values.query")
	}
	defer rows.Close()

	sliceValue := valueOfDest
	for rows.Next() {
		// Create a new value to store our row:
		pointerToNewValue := reflect.New(recordType)
		newValue := reflect.Indirect(pointerToNewValue)

		err = rows.Scan(pointerToNewValue.Interface())
		if err != nil {
			return numberOfRowsReturned, errors.Wrap(err, "[dbr] select.load_all_values.scan")
		}

		// Append our new value to the slice:
		sliceValue = reflect.Append(sliceValue, newValue)

		numberOfRowsReturned++
	}
	valueOfDest.Set(sliceValue)

	if err := rows.Err(); err != nil {
		return numberOfRowsReturned, errors.Wrap(err, "[dbr] select.load_all_values.rows_err")
	}

	return numberOfRowsReturned, nil
}

// LoadValue executes the SelectBuilder and loads the resulting data into a
// primitive value Returns ErrNotFound if no value was found, and it was
// therefore not set. Slow because of the massive use of reflection.
func (b *SelectBuilder) LoadValue(dest interface{}) error {
	// Validate the dest
	valueOfDest := reflect.ValueOf(dest)
	kindOfDest := valueOfDest.Kind()

	if kindOfDest != reflect.Ptr {
		return errors.NewNotValidf("[dbr] Destination must be a pointer")
	}

	//
	// Get full SQL
	//
	tSQL, tArg, err := b.ToSql()
	if err != nil {
		return errors.Wrap(err, "[dbr] select.load_value.tosql")
	}

	fullSql, err := Preprocess(tSQL, tArg)
	if err != nil {
		return err
	}

	if b.Logger.IsInfo() {
		defer log.WhenDone(b.Logger).Info("dbr.SelectBuilder.LoadValue.QueryContext.timing", log.String("sql", fullSql))
	}

	// Run the query:
	rows, err := b.Query(fullSql)
	if err != nil {
		return errors.Wrap(err, "[dbr] select.load_value.query")
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(dest)
		if err != nil {
			return errors.Wrap(err, "[dbr] select.load_value.scan")
		}
		return nil
	}

	if err := rows.Err(); err != nil {
		return errors.Wrap(err, "[dbr] select.load_value.rows_err")
	}

	return errors.NewNotFoundf("[dbr] Entry not found")
}
