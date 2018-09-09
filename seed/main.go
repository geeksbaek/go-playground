package main

import (
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "errors"
    "fmt"
    "io"

    "github.com/geeksbaek/seed"
)

func main() {
    cipherKey := []byte("0123456789012345")
    msg := "A quick brown fox jumped over the lazy dog."

    encrypted, err := encrypt(cipherKey, msg)
    if err != nil {
        panic(err)
    }

    fmt.Printf("CIPHER KEY: %s\n", cipherKey)
    fmt.Printf("ENCRYPTED: %s\n", encrypted)

    decrypted, err := decrypt(cipherKey, encrypted)
    if err != nil {
        panic(err)
    }

    fmt.Printf("DECRYPTED: %s\n", decrypted)

    if msg != decrypted {
        panic("do not match msg and decrypted")
    }

    // CIPHER KEY: 0123456789012345
    // ENCRYPTED: 9VzqUQJh1JWmboAw_tfzzbHdaI8_53NHhBTFoNFPiPn4fqe_G44K0xQpYRyqRWAIp9ao-6OnTkJCh08=
    // DECRYPTED: A quick brown fox jumped over the lazy dog.
}

func encrypt(key []byte, message string) (encmess string, err error) {
    plainText := []byte(message)

    block, err := seed.NewCipher(key)
    if err != nil {
        return
    }

    // IV needs to be unique, but doesn't have to be secure.
    // It's common to put it at the beginning of the ciphertext.
    cipherText := make([]byte, seed.BlockSize+len(plainText))
    iv := cipherText[:seed.BlockSize]
    if _, err = io.ReadFull(rand.Reader, iv); err != nil {
        return
    }

    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(cipherText[seed.BlockSize:], plainText)

    // returns to base64 encoded string
    encmess = base64.URLEncoding.EncodeToString(cipherText)
    return
}

func decrypt(key []byte, securemess string) (decodedmess string, err error) {
    cipherText, err := base64.URLEncoding.DecodeString(securemess)
    if err != nil {
        return
    }

    block, err := seed.NewCipher(key)
    if err != nil {
        return
    }

    if len(cipherText) < seed.BlockSize {
        err = errors.New("Ciphertext block size is too short")
        return
    }

    // IV needs to be unique, but doesn't have to be secure.
    // It's common to put it at the beginning of the ciphertext.
    iv := cipherText[:seed.BlockSize]
    cipherText = cipherText[seed.BlockSize:]

    stream := cipher.NewCFBDecrypter(block, iv)
    // XORKeyStream can work in-place if the two arguments are the same.
    stream.XORKeyStream(cipherText, cipherText)

    decodedmess = string(cipherText)
    return
}
