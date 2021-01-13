// SPDX-FileCopyrightText: 2021 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by "stringer -type=DynamicStatusType"; DO NOT EDIT.

package tds

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TDS_DYNAMIC_UNUSED-0]
	_ = x[TDS_DYNAMIC_HASARGS-1]
	_ = x[TDS_DYNAMIC_SUPPRESS_FMT-2]
	_ = x[TDS_DYNAMIC_BATCH_PARAMS-4]
	_ = x[TDS_DYNAMIC_SUPPRESS_PARAMFMT-8]
}

const (
	_DynamicStatusType_name_0 = "TDS_DYNAMIC_UNUSEDTDS_DYNAMIC_HASARGSTDS_DYNAMIC_SUPPRESS_FMT"
	_DynamicStatusType_name_1 = "TDS_DYNAMIC_BATCH_PARAMS"
	_DynamicStatusType_name_2 = "TDS_DYNAMIC_SUPPRESS_PARAMFMT"
)

var (
	_DynamicStatusType_index_0 = [...]uint8{0, 18, 37, 61}
)

func (i DynamicStatusType) String() string {
	switch {
	case i <= 2:
		return _DynamicStatusType_name_0[_DynamicStatusType_index_0[i]:_DynamicStatusType_index_0[i+1]]
	case i == 4:
		return _DynamicStatusType_name_1
	case i == 8:
		return _DynamicStatusType_name_2
	default:
		return "DynamicStatusType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
