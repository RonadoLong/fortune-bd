package helper

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/gommon/log"
)

const (
	pri_key = `-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQDIsspLnpwUaVDlx3E3a2nz29oi8/qtXoL4ZJMCuAGZbdzDnwzx
NDIhayYcQaoOnI83y+nNbpj68994qpSnh3wj+VFxRD182StB0EmBVxOKC+Qspb0V
Qezu0niEcf3oCKp7AYzf8C2hND+s3kdyXyomLZQn36ewtDh79xp6cCU05wIDAQAB
AoGAfMnWSLCFIZfeKhEZTzkldu/zMRp8ekGys5ltYxpgPDL4OlXxqSQoK2lBF/6o
K0+jKTFL3WTwD9GE2LVPmt7+CxI0ss57zFN51q0VMAoWNHticGtLFnUh705fjCoX
iVNewdGpcrF4uFsrj9N/POPOaeTB+jL/S2erhx3XRXKgiLkCQQDeiDjn3grrHrNs
BC7uaJIYe48NKl3JAJxHi61fNnhoIxn0sOyoYqyr5ugCxkKiyvy34V86Bwh+oiao
5Sr2smUDAkEA5uHw8RkCJWVH9E6wKkUXFrqFHyAllPpKqaP67dczwB2zyUwGguHR
W7x4WBvOOIDGRDaIDuRIv8FzTsimHXfxTQJAZRdlIpBQTXdo8s0DtPJ0TAL1fXmd
mU5ZsHbXj8Vi9YvcorgtCmGpJ36CL6B5bRLhs3cCl43SYhSvk1JoLiHkmQJAMz+F
rs6BRnGzxgvNWKSbWmUudVk6XlYsSnlmknKJPySYqp7gdx7OzNEJ2WzamnojCDMe
gkezyjSTdrJdBP+BpQJAeOw2JgLvF0mm65hVg+u+ipzOb9ZZeM+CCggu6LAIFVHx
ybXR0k/DXOWAW9C31r5G8D5eOLlKbfj9vrAHhYhfxw==
-----END RSA PRIVATE KEY-----`

	pub_key = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDIsspLnpwUaVDlx3E3a2nz29oi
8/qtXoL4ZJMCuAGZbdzDnwzxNDIhayYcQaoOnI83y+nNbpj68994qpSnh3wj+VFx
RD182StB0EmBVxOKC+Qspb0VQezu0niEcf3oCKp7AYzf8C2hND+s3kdyXyomLZQn
36ewtDh79xp6cCU05wIDAQAB
-----END PUBLIC KEY-----`
)

func GetRsaPriKey() *rsa.PrivateKey {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pri_key))
	if err != nil {
		log.Printf("ParseRSAPrivateKeyFromPEM err %v",err)
	}
	return key
}

func GetRsaPublicKey() *rsa.PublicKey {
	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pri_key))
	if err != nil {
		log.Printf("ParseRSAPrivateKeyFromPEM err %v",err)
	}
	return key
}