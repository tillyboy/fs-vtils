package fsv

import (
	"os"
	"os/user"
	"syscall"
)

// Path is the main type of this package. It provides additional type safety
// over the usage of strings as paths and anything revolving around paths and
// files (inodes) can be defined as a method on the new type.
type Path string

// -------------------------- general information --------------------------- //

// Exists checks wether a file (including directories, links etc.) exists at p.
func (p Path) Exists() bool {
	_, err := os.Stat(string(p))
	return err == nil || os.IsExist(err)
}

// Info returns the os.FileInfo of the file located at p.
func (p Path) Info() (os.FileInfo, error) {
	return os.Lstat(string(p))
}

// IsFile returns wether the file sitting at p is a regular file.
// (i.e. not a link, directory etc.)
func (p Path) IsFile() (bool, error) {
	info, err := p.Info()
	if err != nil {
		return false, err
	}
	return info.Mode().IsRegular(), nil
}

// IsSymlink returns true when the inode located at p is a symlink.
func (p Path) IsSymlink() (bool, error) {
	info, err := p.Info()
	if err != nil {
		return false, err
	}
	return (info.Mode()&os.ModeSymlink != 0), nil
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
	info, err := p.Info()
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

// Ls tries to list the all paths of entries of the directory residing at p.
func (p Path) Ls() (PathList, error) {
	d, err := os.Open(string(p))
	defer closeOrPanic(d)
	if err != nil {
		return nil, err
	}

	es, err := d.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	var ps PathList = make([]Path, len(es))
	for i, e := range es {
		ps[i] = p.ExtendStr(e)
	}

	return ps, nil
}

// ----------------------------- based on name ------------------------------ //

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

// ------------------------------ mode & owner ------------------------------ //

// Mode returns the os.FileMode of the file located at p.
func (p Path) Mode() (os.FileMode, error) {
	info, err := p.Info()
	if err != nil {
		return 0, err
	}
	return info.Mode(), nil
}

// Chmod changes the permissions of the file located at p.
func (p Path) Chmod(mode os.FileMode) error {
	f, err := os.Open(string(p))
	if err != nil {
		return err
	}
	defer closeOrPanic(f)

	return f.Chmod(mode)
}

// Owner returns the *user.User who owns the file located at p.
func (p Path) Owner(os.FileMode) (*user.User, error) {
	info, err := os.Lstat(string(p))
	if err != nil {
		return &user.User{}, err
	}

	sysInfo := info.Sys().(*syscall.Stat_t)
	return user.LookupId(string(sysInfo.Uid))
}

// Chown changes the owner of the file located at p.
func (p Path) Chown(uid, gid int) error {
	f, err := os.Open(string(p))
	if err != nil {
		return err
	}
	defer closeOrPanic(f)

	return f.Chown(uid, gid)
}

// -------------------------------- contents -------------------------------- //

// Size returns the size of a file in bytes.
func (p Path) Size() (int64, error) {
	info, err := p.Info()
	if err != nil {
		return 0, err
	} else if !info.Mode().IsRegular() {
		return 0, FILE_OPERATION.new(p, _FLAG_EMPTY)
	}

	return info.Size(), nil
}

// CountRunes returns the length the file contents converted to a rune-slice.
// This involves opening the file and converting its contents to an intermittent
// string. Performance might be suboptimal.
func (p Path) CountRunes() (int, error) {
	bytes, err := p.ReadBytes()
	if err != nil {
		return 0, err
	}
	return len([]rune(string(bytes))), nil
}

// CountLines returns the amount of Newline-Characters ('\n') found in a file.
// This involves reading the file and ranging over its contents. Performance
// might be suboptimal.
func (p Path) CountLines() (int, error) {
	contents, err := p.ReadString()
	if err != nil {
		return 0, err
	}
	c := 0

	for r := range []rune(contents) {
		if r == '\n' {
			c++
		}
	}

	// -1 to remove trailing newline: TODO: do or don't?
	return c - 1, nil
}
