package helpers

import (
	"crypto/md5"
	"fmt"
	"io"
)

func Gravatar(email string) string {
	hash := md5.New()
	io.WriteString(hash, email)
	return fmt.Sprintf("http://www.gravatar.com/avatar/%x", hash.Sum(nil))
}
