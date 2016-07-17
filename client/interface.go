package client

type Client interface {
	// Store accepts an id and a payload in bytes
	Store(id, payload []byte) (aesKey []byte, err error)

	// Retrieve accepts an id and an AES key, and requests the original
	// (decrypted) bytes
	Retrieve(id, aesKey []byte) (payload []byte, err error)
}
