package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"github.com/go-ini/ini"
	"log"
)

var aesKey []byte

// 初始化：从 config.ini 读取 AES 密钥
func init() {
	// 加载 config.ini 文件
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalf("无法加载 config.ini: %v", err)
	}
	// 读取 [security] 节下的 aes_secret_key
	keyStr := cfg.Section("security").Key("aes_secret_key").String()
	if keyStr == "" {
		log.Fatal("config.ini 中未配置 aes_secret_key")
	}
	// 转换为字节数组（十六进制字符串需解码）
	aesKey, err = hex.DecodeString(keyStr) // 如果密钥是十六进制格式
	// 若密钥是原始字符串（非十六进制），则直接转换：
	// aesKey = []byte(keyStr)

	// 验证密钥长度
	if len(aesKey) != 16 && len(aesKey) != 24 && len(aesKey) != 32 {
		log.Fatalf("AES 密钥长度必须是 16/24/32 字节，当前长度: %d", len(aesKey))
	}
}

// AESEncrypt 加密
func AESEncrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	plaintextBytes := []byte(plaintext)
	padding := blockSize - len(plaintextBytes)%blockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	plaintextBytes = append(plaintextBytes, padtext...)

	ciphertext := make([]byte, blockSize+len(plaintextBytes))
	iv := ciphertext[:blockSize]
	if _, err := rand.Read(iv); err != nil { // 需要导入 "crypto/rand"
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[blockSize:], plaintextBytes)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// AESDecrypt 解密
func AESDecrypt(ciphertext string) (string, error) {
	ciphertextBytes, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	iv := ciphertextBytes[:blockSize]
	ciphertextBytes = ciphertextBytes[blockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertextBytes, ciphertextBytes)

	padding := int(ciphertextBytes[len(ciphertextBytes)-1])
	return string(ciphertextBytes[:len(ciphertextBytes)-padding]), nil
}
