package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
)

// 对称加密的加密方法

// encryptAES 使用AES-CBC模式加密数据

func EncryptAES(plainText, key []byte) (string, error) {
	// 生成一个16字节的随机IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("failed to generate IV: %v", err)
	}

	// 创建AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %v", err)
	}

	// 使用CBC模式加密
	cipherText := make([]byte, len(plainText)+aes.BlockSize)
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, plainText)

	// 将IV和密文组合在一起
	combinedData := append(iv, cipherText...)

	// 将组合后的数据编码为Base64格式
	encodedData := base64.StdEncoding.EncodeToString(combinedData)
	return encodedData, nil
}

// 对称加密的解密方法

// decryptAES 使用AES-CBC模式解密数据

func DecryptAES(encryptedData, key []byte) (string, error) {

	// 解码Base64数据
	decodedData, err := base64.StdEncoding.DecodeString(string(encryptedData))
	decodeKey, err := base64.StdEncoding.DecodeString(string(key))
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 data: %v", err)
	}

	// 分离IV和密文
	if len(decodedData) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := decodedData[:aes.BlockSize]
	cipherText := decodedData[aes.BlockSize:]

	// 创建AES cipher
	block, err := aes.NewCipher(decodeKey)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %v", err)
	}

	// 使用CBC模式解密
	if len(cipherText)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)

	// 去除PKCS7填充
	plainText, err := unpad(cipherText)
	if err != nil {
		return "", fmt.Errorf("failed to unpad plaintext: %v", err)
	}

	return string(plainText), nil
}

// unpad 去除PKCS7填充

func unpad(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return nil, errors.New("padding is incorrect")
	}

	padding := int(src[length-1])
	if padding > length {
		return nil, errors.New("padding is incorrect")
	}

	for i := length - padding; i < length; i++ {
		if int(src[i]) != padding {
			return nil, errors.New("padding is incorrect")
		}
	}

	return src[:length-padding], nil
}

// EncryptRSA 使用 RSA 公钥加密数据
func EncryptRSA(data, publicKeyPEM string) (string, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return "", errors.New("failed to parse PEM block containing the public key")
	}

	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse public key: %v", err)
	}

	hashed := sha256.New()
	encryptedBytes, err := rsa.EncryptOAEP(
		hashed, rand.Reader, pub, []byte(data), nil,
	)
	if err != nil {
		return "", fmt.Errorf("encryption failed: %v", err)
	}

	return base64.URLEncoding.EncodeToString(encryptedBytes), nil
}

// 非对称加密的解密方法

// DecryptRSA 使用 RSA 私钥解密数据
func DecryptRSA(encryptedData, privateKeyPEM string) (string, error) {
	cipherText, err := base64.URLEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 data: %v", err)
	}

	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return "", errors.New("failed to parse PEM block containing the private key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}

	hashed := sha256.New()
	decryptedBytes, err := rsa.DecryptOAEP(hashed, rand.Reader, priv, cipherText, nil)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %v", err)
	}

	return string(decryptedBytes), nil
}

// 对称加密秘钥生成方法

// GenerateAESKey 生成一个随机的AES密钥，并返回Base64编码的字符串
func GenerateAESKey(keySize int) (string, error) {
	// AES密钥的长度必须是16, 24, 或32字节
	if keySize != 16 && keySize != 24 && keySize != 32 {
		return "", fmt.Errorf("invalid key size: %d. Valid sizes are 16, 24, or 32 bytes", keySize)
	}

	// 生成随机密钥
	key := make([]byte, keySize)
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return "", fmt.Errorf("failed to generate random key: %v", err)
	}

	// 将密钥编码为Base64格式
	encodedKey := base64.StdEncoding.EncodeToString(key)
	return encodedKey, nil
}

// GenerateRSAKeys 生成 RSA 私钥和公钥，并以 PEM 格式返回
func GenerateRSAKeys(bits int) (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		},
	)

	publicKeyBytes := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	publicKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicKeyBytes,
		},
	)

	return string(privateKeyPEM), string(publicKeyPEM), nil
}
