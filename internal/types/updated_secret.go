package types

type UpdatedSecret struct {
	ID       int64
	Name     *string
	Data     []byte
	Metadata *string
}
