package utils

import (
    "io"
    "crypto/rand"
)

var StdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

func NewPassword(length int) string {
    return rand_char(length, StdChars)
}

func rand_char(length int, chars []byte) string {
    new_pword := make([]byte, length)
    random_data := make([]byte, length)
    clen := byte(len(chars))
    if _, err := io.ReadFull(rand.Reader, random_data); err != nil {
        panic(err)
    }
    for i, c := range random_data {
        new_pword[i] = chars[c%clen]
    }
    return string(new_pword)
}
