package processor

// Holds the result after processing the hashes for the file
type Result struct {
	File       string
	MD5        string
	SHA1       string
	SHA256     string
	SHA512     string
	Blake2b256 string
	Blake2b512 string
	Sha3224    string
	Sha3256    string
	Sha3384    string
	Sha3512    string
	Bytes      int64
}
