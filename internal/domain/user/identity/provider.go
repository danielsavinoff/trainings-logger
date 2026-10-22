package identity

type Provider string

const ProviderDropbox Provider = "Dropbox"

func NewProvider(value string) (Provider, error) {
	provider := Provider(value)
	if provider != ProviderDropbox {
		return "", ErrInvalidProvider
	}

	return provider, nil
}
