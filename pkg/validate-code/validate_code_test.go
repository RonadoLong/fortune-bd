package validate_code

import (
	"log"
	"strings"
	"testing"
)

func BenchmarkGenValidateCode(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			got := GenValidateCode(6)
			log.Println(got)
		}
	})
}

func TestGenValidateCode(t *testing.T) {
	//i := 0
	if strings.HasPrefix("19888", "0") {
		t.Log("ffff")
	}
	//for i<1000 {
	//	a := GenValidateCode(6)
	//	t.Log(a)
	//	//if strings.HasPrefix(a, "0"){
	//	//
	//	//}
	//	i++
	//}
}
