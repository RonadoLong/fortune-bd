package gzip

import "testing"

func Test_Decompress_Success(t *testing.T) {
	buf, _ := GZipCompress("huobi")

	result, _ := GZipDecompress(buf)

	expected := "huobi"
	if result != expected {
		t.Errorf("expected: %s, actual: %s", expected, result)
	}
}
