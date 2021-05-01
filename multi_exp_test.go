package shuffle

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"
)

func TestMulti_exp(t *testing.T) {

	com := Common{EC.G,EC.H}
	// prover generate  C
	r1,r2,r3 := big.NewInt(11), big.NewInt(22),big.NewInt(33)
	a1, a2, a3 := big.NewInt(2), big.NewInt(1),big.NewInt(3)

	// it's ok whatever C1 and C2 is
	C1, C2, C3 := GenCom(com,big.NewInt(1),r1), GenCom(com,big.NewInt(2),r2), GenCom(com,big.NewInt(3),r3)
	C_array := []ECPoint{C1,C2,C3}
	a_array := []*big.Int{a1, a2, a3}
	rou, err := rand.Int(rand.Reader, EC.N)
	check(err)
	C_value := Encrypt1(com, rou).Add(C1.Mult(a1)).Add(C2.Mult(a2)).Add(C3.Mult(a3))
	CA_array := []ECPoint{GenCom(com,a1,r1),GenCom(com,a2,r2),GenCom(com,a3,r3)}
	r_array := []*big.Int{r1,r2,r3}
	prove := Multi_exp(3,Statement{C_array,C_value,CA_array},a_array,r_array,rou)

	fmt.Println("VerifyMulExp",VerifyMulExp(prove))
}
