// Code generated by "stringer -type=SliceViewSignals"; DO NOT EDIT.

package giv

import (
	"fmt"
	"strconv"
)

const _SliceViewSignals_name = "SliceViewDoubleClickedSliceViewSignalsN"

var _SliceViewSignals_index = [...]uint8{0, 22, 39}

func (i SliceViewSignals) String() string {
	if i < 0 || i >= SliceViewSignals(len(_SliceViewSignals_index)-1) {
		return "SliceViewSignals(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _SliceViewSignals_name[_SliceViewSignals_index[i]:_SliceViewSignals_index[i+1]]
}

func (i *SliceViewSignals) FromString(s string) error {
	for j := 0; j < len(_SliceViewSignals_index)-1; j++ {
		if s == _SliceViewSignals_name[_SliceViewSignals_index[j]:_SliceViewSignals_index[j+1]] {
			*i = SliceViewSignals(j)
			return nil
		}
	}
	return fmt.Errorf("String %v is not a valid option for type SliceViewSignals", s)
}
