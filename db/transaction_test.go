package db

import (
	"database/sql"
	"github.com/blend/go-sdk/assert"
	"testing"
)

func TestTransactionInnerRollback(t *testing.T) {
	assert := assert.New(t)

	tx, err := defaultDB().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	assert.Nil(IgnoreExecResult(defaultDB().Invoke(OptTx(tx)).Exec("CREATE TABLE IF NOT EXISTS unique_obj (id int not null primary key, name varchar)")))
	assert.Nil(defaultDB().Invoke(OptTx(tx)).Create(&uniqueObj{ID: 1, Name: "one"}))
	var verify uniqueObj
	assert.Nil(secondArgErr(defaultDB().Invoke(OptTx(tx)).Get(&verify, 1)))

	inner, err := tx.Begin()
	assert.Nil(err)
	defaultDB().Invoke(OptTx(inner)).Create(&uniqueObj{ID: 2, Name: "two"})

	// Available on both
	found, err := defaultDB().Invoke(OptTx(tx)).Query("Select 1 from unique_obj where id = $1", 2).Any()
	assert.Nil(err)
	assert.True(found)

	found, err = defaultDB().Invoke(OptTx(inner)).Query("Select 1 from unique_obj where id = $1", 2).Any()
	assert.Nil(err)
	assert.True(found)


	illegal, err := tx.Begin()
	assert.NotNil(err)
	assert.Nil(illegal)
	assert.Equal(ErrInnerTxExists, err.Error())

	err = inner.Rollback()
	assert.Nil(err)

	// Not available on tx
	found, err = defaultDB().Invoke(OptTx(tx)).Query("Select 1 from unique_obj where id = $1", 2).Any()
	assert.Nil(err)
	assert.False(found)

	// Cant roll back resolved tx
	err = inner.Rollback()
	assert.Equal(sql.ErrTxDone, err)
}

func TestTransactionOuterRollback(t *testing.T) {
	assert := assert.New(t)

	tx, err := defaultDB().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	assert.Nil(IgnoreExecResult(defaultDB().Invoke(OptTx(tx)).Exec("CREATE TABLE IF NOT EXISTS unique_obj (id int not null primary key, name varchar)")))

	outer, err := tx.Begin()
	assert.Nil(err)

	assert.Nil(defaultDB().Invoke(OptTx(outer)).Create(&uniqueObj{ID: 1, Name: "one"}))
	var verify uniqueObj
	assert.Nil(secondArgErr(defaultDB().Invoke(OptTx(outer)).Get(&verify, 1)))

	inner, err := outer.Begin()
	assert.Nil(err)
	defaultDB().Invoke(OptTx(inner)).Create(&uniqueObj{ID: 2, Name: "two"})

	// Available on all
	found, err := defaultDB().Invoke(OptTx(inner)).Query("Select 1 from unique_obj where id = $1", 2).Any()
	assert.Nil(err)
	assert.True(found)

	found, err = defaultDB().Invoke(OptTx(outer)).Query("Select 1 from unique_obj where id = $1", 2).Any()
	assert.Nil(err)
	assert.True(found)

	found, err = defaultDB().Invoke(OptTx(tx)).Query("Select 1 from unique_obj where id = $1", 2).Any()
	assert.Nil(err)
	assert.True(found)

	err = outer.Rollback()
	assert.Nil(err)

	// Outer can't be committed after rollback
	err = outer.Commit()
	assert.Equal(sql.ErrTxDone, err)

	// Not available on tx
	found, err = defaultDB().Invoke(OptTx(tx)).Query("Select 1 from unique_obj where id = $1", 2).Any()
	assert.Nil(err)
	assert.False(found)

	// Inner can't be rolled back
	err = inner.Rollback()
	assert.Equal(sql.ErrTxDone, err)

	// Inner can't be committed back
	err = inner.Commit()
	assert.Equal(sql.ErrTxDone, err)
}