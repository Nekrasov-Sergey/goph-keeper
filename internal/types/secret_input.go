package types

type SecretInput struct {
	Name     string
	Type     SecretType
	Data     []byte
	Metadata *string
}
