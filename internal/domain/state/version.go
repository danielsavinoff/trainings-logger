package state

type Version int

const Version4 Version = 4

func NewVersion(value int) (Version, error) {
	version := Version(value)
	if version != Version4 {
		return 0, ErrUnsupportedVersion
	}

	return version, nil
}
