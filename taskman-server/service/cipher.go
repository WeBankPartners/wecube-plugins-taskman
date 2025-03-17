package service

import (
	"github.com/WeBankPartners/go-common-lib/cipher"
	"strings"
)

var (
	// DEFALT_CIPHER CMDB默认加密
	DEFALT_CIPHER = "CIPHER_A"
	// DEFALT_CIPHER_B taskman敏感数据加密
	DEFALT_CIPHER_B = "CIPHER_B"
	// DEFALT_CIPHER_C taskman密码加密
	DEFALT_CIPHER_C = "CIPHER_C"
	CIPHER_MAP      = map[string]string{"CIPHER_A": "{cipher_a}", "CIPHER_B": "{cipher_b}", "CIPHER_C": "{cipher_c}"}
)

func AesEnPasswordByGuid(guid, seed, password, cipherStr string) (string, error) {
	if seed == "" {
		return password, nil
	}
	for _, _cipher := range CIPHER_MAP {
		if strings.HasPrefix(password, _cipher) {
			return password, nil
		}
	}
	if cipherStr == "" {
		cipherStr = DEFALT_CIPHER
	}
	md5sum := cipher.Md5Encode(guid + seed)
	enPassword, err := cipher.AesEncode(md5sum[0:16], password)
	if err != nil {
		return "", err
	}
	return CIPHER_MAP[cipherStr] + enPassword, nil
}

func AesDePasswordByGuid(guid, seed, password string) (string, error) {
	var cipherStr string
	for _, _cipher := range CIPHER_MAP {
		if strings.HasPrefix(password, _cipher) {
			cipherStr = _cipher
			break
		}
	}
	if cipherStr == "" {
		return password, nil
	}
	password = password[len(cipherStr):]
	md5sum := cipher.Md5Encode(guid + seed)
	dePassword, err := cipher.AesDecode(md5sum[0:16], password)
	if err != nil {
		return "", err
	}
	return dePassword, nil
}
