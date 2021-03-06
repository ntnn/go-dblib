// SPDX-FileCopyrightText: 2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by "gen_type BigTime time.Time"; DO NOT EDIT.

package integration

import (
	"database/sql"

	"testing"

	"time"
)

// DoTestBigTime tests the handling of the BigTime.
func DoTestBigTime(t *testing.T) {
	TestForEachDB("TestBigTime", t, testBigTime)
	//
}

func testBigTime(t *testing.T, db *sql.DB, tableName string) {
	pass := make([]interface{}, len(samplesBigTime))
	mySamples := make([]time.Time, len(samplesBigTime))

	for i, sample := range samplesBigTime {

		mySample := sample

		pass[i] = mySample
		mySamples[i] = mySample
	}

	rows, teardownFn, err := SetupTableInsert(db, tableName, "bigtime", pass...)
	if err != nil {
		t.Errorf("Error preparing table: %v", err)
		return
	}
	defer rows.Close()
	defer teardownFn()

	i := 0
	var recv time.Time
	for rows.Next() {
		if err := rows.Scan(&recv); err != nil {
			t.Errorf("Scan failed on %dth scan: %v", i, err)
			continue
		}

		if recv != mySamples[i] {

			t.Errorf("Received value does not match passed parameter")
			t.Errorf("Expected: %v", mySamples[i])
			t.Errorf("Received: %v", recv)
		}

		i++
	}

	if err := rows.Err(); err != nil {
		t.Errorf("Error preparing rows: %v", err)
	}

	if i != len(pass) {
		t.Errorf("Only read %d values from database, expected to read %d", i, len(pass))
	}
}
