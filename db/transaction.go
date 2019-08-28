package db

import (
	"database/sql"
	"fmt"
	"github.com/blend/go-sdk/ex"
	"strconv"
)

// ErrInnerTxExists is an error that occurs when you try to Begin multiple transactions inside of the same parent
// transaction. You can have multiple inner save points, but you must nest them one inside the other. A single
// parent can have no more than one inner transaction
const ErrInnerTxExists ex.Class = "db: inner transaction already exists"

// Tx is a database transaction
type Tx struct {
	*sql.Tx

	savePoint int
	nested    *Tx
	resolved  bool
}

// Begin starts a new nested transaction
func (tx *Tx) Begin() (*Tx, error) {
	if tx.nested != nil && !tx.nested.resolved {
		return nil, ex.New(ErrInnerTxExists)
	}
	tx.nested = &Tx{
		Tx:        tx.Tx,
		savePoint: tx.savePoint + 1,
	}

	_, err := tx.Exec(fmt.Sprintf("SAVEPOINT PT%s", strconv.Itoa(tx.nested.savePoint)))

	if err != nil {
		return nil, err
	}

	return tx.nested, nil
}

// Rollback rolls back the transaction. If this is an inner transaction it rolls back to the last save point
func (tx *Tx) Rollback() error {
	tx.resolved = true

	if tx.savePoint > 0 {
		_, err := tx.Exec(fmt.Sprintf("ROLLBACK TO SAVEPOINT PT%s", strconv.Itoa(tx.savePoint)))
		return err
	}

	return tx.Tx.Rollback()
}

// Commit commits the transaction. If this is an inner transaction it releases the last save point.
func (tx *Tx) Commit() error {
	if tx.nested != nil && !tx.nested.resolved {
		err := tx.nested.Commit()
		if err != nil {
			return err
		}
	}

	tx.resolved = true

	if tx.savePoint > 0 {
		_, err := tx.Exec(fmt.Sprintf("RELEASE SAVEPOINT PT%s", strconv.Itoa(tx.savePoint)))

		return err
	}

	return tx.Tx.Commit()
}