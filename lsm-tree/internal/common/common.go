package common

type Serializable interface {
	Serialize() []byte
	Deserialize([]byte) Serializable
}
