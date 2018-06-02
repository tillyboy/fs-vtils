package fsv

import (
	"os"
)

// Path is the main type of this package. It provides additional type safety
// over the usage of strings as paths and anything revolving around paths and
// files (inodes) can be defined as a method on the new type.
type Path string

// -------------------------- information methods --------------------------- //

// TODO: Grep(re string) (lineNo int, startByte int, stopByte int, line string, match string)

// Exists checks wether a file (including directories, links etc.) exists at p.
func (p Path) Exists() bool {
	_, err := os.Stat(string(p))
	return err == nil || os.IsExist(err)
}

// Info returns the os.FileInfo of the file located at p.
func (p Path) Info() (os.FileInfo, error) {
	return os.Lstat(string(p))
}

// Mode returns the os.FileInfo of the file located at p.
func (p Path) Mode() (os.FileMode, error) {
	nfo, err := p.Info()
	if err != nil {
		return 0, err
	}
	return nfo.Mode(), nil
}

// IsFile returns wether the file sitting at p is a regular file.
// (i.e. not a link, directory etc.)
func (p Path) IsFile() (bool, error) {
	nfo, err := p.Info()
	if err != nil {
		return false, err
	}
	return nfo.Mode().IsRegular(), nil
}

// IsSymlink returns true when the inode located at p is a symlink.
func (p Path) IsSymlink() (bool, error) {
	nfo, err := p.Info()
	if err != nil {
		return false, err
	}
	return (nfo.Mode()&os.ModeSymlink != 0), nil
}

// Follow tries to read the path that a symlink residing at p points to.
func (p Path) Follow() (Path, error) {
	target, err := os.Readlink(string(p))
	return Path(target), err
}

// Target tries to recursviely follow a symlink until a non-symlink is found.
func (p Path) Target() (Path, error) {
	isLn, err := p.IsSymlink()
	if err != nil {
		return Path(""), err
	}

	for isLn {
		p, err = p.Follow()
		if err != nil {
			return Path(""), err
		}

		isLn, err = p.IsSymlink()
		if err != nil {
			return Path(""), err
		}
	}

	return p, nil
}

// IsDir returns true when the inode located at p is a directory.
func (p Path) IsDir() (bool, error) {
	nfo, err := p.Info()
	if err != nil {
		return false, err
	}
	return nfo.IsDir(), nil
}

// Ls tries to list the all paths of entries of the directory residing at p.
func (p Path) Ls() (PathList, error) {
	d, err := os.Open(string(p))
	defer d.Close()
	if err != nil {
		return nil, err
	}

	es, err := d.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	var ps PathList = make([]Path, len(es))
	for i, e := range es {
		ps[i] = p.AppendStr(e)
	}

	return ps, nil
}

// IsOsFile checks the basename of p to determine wether the file located at p
// is generated by the operating system or not.
// Filenames that lead to a true return value:
//    - .DS_Store
//    - __MACOSX
//    - desktop.ini
//    - Thumbs.db
//    - thumbs.db
func (p Path) IsOsFile() bool {
	name := string(p.Base())

	for _, n := range []string{
		".DS_Store",
		"___MACOSX",
		"desktop.ini",
		"Thumbs.db",
		"thumbs.db",
	} {
		if name == n {
			return true
		}
	}

	return false

}

// IsHidden determines wether a file is considered hidden by name.
func (p Path) IsHidden() bool {
	// only for unix so far => TODO: other OS's
	return string(p.Base())[0] == '.'
}

// IsVisible is the negation of IsHidden.
func (p Path) IsVisible() bool {
	return !p.IsHidden()
}
