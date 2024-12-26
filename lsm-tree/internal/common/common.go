package common

type Serializable interface {
	New() Serializable
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}
