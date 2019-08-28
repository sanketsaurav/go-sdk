package db

import (
	"github.com/blend/go-sdk/assert"
	"testing"
)

func TestTransactionSingleCheckpointRollback(t *testing.T) {
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

	err = inner.Rollback()
	assert.Nil(err)

	// Not available on tx
	found, err = defaultDB().Invoke(OptTx(tx)).Query("Select 1 from unique_obj where id = $1", 2).Any()
	assert.Nil(err)
	assert.False(found)
}
