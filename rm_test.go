package fsv_test

import (
	"fsv"
)

// case01: target is file => no error
// case02: target is symlink => no error
// case03: target doesn't exist => no error
// case04: target is dir & no r flag => error
// case05: target is dir & r flag => no error

var testLocRm = testDir.AppendStr("Rm")
var rmTests = []testStruct{
	testFsvErr{

		_name: "case01",
		modify: func() error {
			path := testLocRm.AppendStr("case01/file")
			return path.Rm()
		},
		_expect: nil,
	}, testFsvErr{

		_name: "case02",
		modify: func() error {
			path := testLocRm.AppendStr("case02/symlink")
			return path.Rm()
		},
		_expect: nil,
	}, testFsvErr{

		_name: "case03",
		modify: func() error {
			path := testLocRm.AppendStr("case03/none")
			return path.Rm()
		},
		_expect: nil,
	}, testFsvErr{

		_name: "case04",
		modify: func() error {
			path := testLocRm.AppendStr("case04/dir")
			return path.Rm()
		},
		_expect: fsv.MISSING_REC_FLAG,
	}, testFsvErr{

		_name: "case05",
		modify: func() error {
			path := testLocRm.AppendStr("case05/dir")
			return path.Rm('r')
		},
		_expect: nil,
	},
}
