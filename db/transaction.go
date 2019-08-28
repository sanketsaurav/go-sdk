package db

import (
	"database/sql"
	"fmt"
	"github.com/blend/go-sdk/ex"
	"strconv"
	"sync"
)

// ErrInnerTxExists is an error that occurs when you try to Begin multiple transactions inside of the same parent
// transaction. You can have multiple inner save points, but you must nest them one inside the other. A single
// parent can have no more than one inner transaction
const ErrInnerTxExists ex.Class = "db: inner transaction already exists"


// Tx is a database transaction
type Tx struct {
	*sql.Tx

	savePoint     int
	nested        *Tx
	resolved      bool
	resolutionMux sync.Mutex
}

// Begin starts a new nested transaction
func (tx *Tx) Begin() (*Tx, error) {
	tx.resolutionMux.Lock()
	defer tx.resolutionMux.Unlock()
	if tx.nested != nil && !tx.nested.resolved {
		return nil, ex.New(ErrInnerTxExists)
	}
	if tx.savePoint > 0 && tx.resolved {
		return nil, sql.ErrTxDone
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
	tx.resolutionMux.Lock()
	defer tx.resolutionMux.Unlock()
	if tx.resolved {
		return sql.ErrTxDone
	}
	tx.markResolved()

	if tx.savePoint > 0 {
		_, err := tx.Exec(fmt.Sprintf("ROLLBACK TO SAVEPOINT PT%s", strconv.Itoa(tx.savePoint)))
		return err
	}

	return tx.Tx.Rollback()
}

func (tx *Tx) markResolved() {
	if tx.nested != nil {
		tx.nested.resolutionMux.Lock()
		tx.nested.markResolved()
		tx.nested.resolutionMux.Unlock()
	}
	tx.resolved = true
}

// Commit commits the transaction. If this is an inner transaction it releases the last save point.
func (tx *Tx) Commit() error {
	tx.resolutionMux.Lock()
	defer tx.resolutionMux.Unlock()

	// If the client is trying to re-commit a transaction thats resolved, error
	if tx.resolved {
		return sql.ErrTxDone
	}

	return tx.innerCommit()
}

func (tx *Tx) innerCommit() error {
	if tx.nested != nil && !tx.nested.resolved {
		tx.resolved = true
		err := tx.nested.innerCommit()
		if err != nil {
			return err
		}
	}

	if tx.savePoint > 0 && !tx.resolved {
		tx.resolved = true
		_, err := tx.Exec(fmt.Sprintf("RELEASE SAVEPOINT PT%s", strconv.Itoa(tx.savePoint)))
		return err
	}

	// This is a "root" transaction commit it
	return tx.Tx.Commit()
}