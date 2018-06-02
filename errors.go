package fsv

type errNumber uint8

const (
	no_INVALID_FLAG uint = 1 + iota
	no_OCCUPIED_PATH
	no_MISSING_TARGETDIR
	no_MISSING_REC_FLAG
	no_MISSING_OS_SUPPORT
	no_UNKNOWN_ERR
)

var (
	INVALID_FLAG       = Error{no_INVALID_FLAG, _PATH_EMPTY, _FLAG_EMPTY}
	OCCUPIED_PATH      = Error{no_OCCUPIED_PATH, _PATH_EMPTY, _FLAG_EMPTY}
	MISSING_TARGETDIR  = Error{no_MISSING_TARGETDIR, _PATH_EMPTY, _FLAG_EMPTY}
	MISSING_REC_FLAG   = Error{no_MISSING_REC_FLAG, _PATH_EMPTY, _FLAG_EMPTY}
	MISSING_OS_SUPPORT = Error{no_MISSING_OS_SUPPORT, _PATH_EMPTY, _FLAG_EMPTY}
	UNKNOWN_ERR        = Error{no_UNKNOWN_ERR, _PATH_EMPTY, _FLAG_EMPTY}
)

const (
	_FLAG_EMPTY rune = 0
	_PATH_EMPTY Path = ""
)

type Error struct {
	Id   uint
	Path Path
	Flag rune
}

func (e Error) Error() string {
	switch e.Id {
	case no_INVALID_FLAG:
		return "Invalid flag: " + string(e.Flag)

	case no_OCCUPIED_PATH:
		return "Occupied path: " + string(e.Path)

	case no_MISSING_TARGETDIR:
		return "Inexistent target directory: " + string(e.Path)

	case no_MISSING_REC_FLAG:
		return "Copying/Moving dir requires recursive flag."

	case no_MISSING_OS_SUPPORT:
		return "Operating system does not support this operation."

	case no_UNKNOWN_ERR:
		return "Unkown error."

	default:
		panic("Tried to call Error() on unidentifiable error :(")
	}

}

func (proto Error) new(path Path, flag rune) Error {
	return Error{
		proto.Id,
		path,
		flag,
	}
}

// TODO: errors for PathList.Each, to which the returned errors will be appended

func (proto Error) IsTypeOf(e error) bool {
	switch {
	case e == nil && error(proto) == nil:
		return true
	case e != nil && error(proto) == nil:
		return false
	case e == nil && error(proto) != nil:
		return false
	}

	fsve, ok := e.(Error)
	if !ok {
		return false
	}

	if fsve.Id == proto.Id {
		return true
	} else {
		return false
	}
}
