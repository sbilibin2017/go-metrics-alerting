package key

type KeyEngineInterface interface {
	Encode(key *Key) (string, error)
	Decode(key string) (*Key, error)
}
