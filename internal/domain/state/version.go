package state

type Version int

const Version1 Version = 1

func NewVersion(value int) (Version, error) {
	version := Version(value)
	if version != Version1 {
		return 0, ErrUnsupportedVersion
	}

	return version, nil
}
