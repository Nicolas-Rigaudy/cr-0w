package contract

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
)

// ZeroSHA256 is the prev_sha256 of the first validation of a mission.
var ZeroSHA256 = strings.Repeat("0", 64)

// Canonical encodes v as canonical JSON (req_037): object keys sorted, separators `,` and
// `:`, UTF-8 without ASCII or HTML escaping. It matches Python's
// json.dumps(v, sort_keys=True, separators=(",", ":"), ensure_ascii=False), except that Go
// escapes U+2028 and U+2029.
func Canonical(v any) ([]byte, error) {
	raw, err := encode(v)
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var tree any
	if err := dec.Decode(&tree); err != nil {
		return nil, err
	}
	return encode(tree) // maps encode with sorted keys
}

func encode(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return bytes.TrimSuffix(buf.Bytes(), []byte("\n")), nil
}

// SHA256Hex is the lowercase hex SHA-256 of b.
func SHA256Hex(b []byte) string {
	s := sha256.Sum256(b)
	return hex.EncodeToString(s[:])
}

// PayloadSHA256 is the payload_sha256 of one publish entry: the SHA-256 of its canonical JSON.
func PayloadSHA256(p Publish) (string, error) {
	b, err := Canonical(p)
	if err != nil {
		return "", err
	}
	return SHA256Hex(b), nil
}
