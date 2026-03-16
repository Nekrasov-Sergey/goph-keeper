package types

type SecretPayload struct {
	Name     string
	Type     SecretType
	Data     []byte
	Metadata *string
}
