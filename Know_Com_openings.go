package shuffle

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type KCO struct {
	C_array          []ECPoint
	c0               ECPoint
	e                *big.Int
	z_array			 []*big.Int
	s                *big.Int
}

func ComMul(com Common, x_array []*big.Int, r0 *big.Int) ECPoint {
	res := com.g2.Mult(r0)
	for i:=0;i<len(x_array);i++{
		res = res.Add(com.g1.Mult(x_array[i]))
	}
	return res
}

// argument of knowledge of commitment openings
// <linear algebra with Sub-linear Zero-Knowledge Arguments> 2009 A
func KCO_proof(com Common, x_array []*big.Int, r_array []*big.Int) KCO {
	C_array := []ECPoint{}
	m := len(x_array)
	for i:=0;i<m;i++{
		C_array = append(C_array, Com(com, x_array[i], r_array[i]))
	}
	r0, err := rand.Int(rand.Reader, EC.N)
	check(err)
	// x0 is array in the article
	x0, err := rand.Int(rand.Reader, EC.N)
	check(err)
	c0 := Com(com, x0, r0)
	x1_array := []*big.Int{x0}
	x1_array = append(x1_array,x_array...)
	r1_array := []*big.Int{r0}
	r1_array = append(r1_array,r_array...)
	c_array := []ECPoint{c0}
	c_array = append(c_array,C_array...)

	//TODO: e should be hash()
	e, err := rand.Int(rand.Reader, EC.N)
	check(err)

	z_array := []*big.Int{}
	for i:=0;i<=m;i++{
		z_array = append(z_array, new(big.Int).Mul(new(big.Int).Exp(e, big.NewInt(int64(i)),nil),x1_array[i]))
	}

	s := big.NewInt(0)
	for i:=0;i<=m;i++{
		s = new(big.Int).Add(s, new(big.Int).Mul(new(big.Int).Exp(e, big.NewInt(int64(i)),nil),r1_array[i]))
	}

	k := KCO{
		C_array: c_array,
		e: e,
		z_array: z_array,
		s:s,
	}
	fmt.Println("KCO generated")
	return k
}

func VerifyKCO(com Common, kco KCO) bool{
	v_left := EC.G.Mult(big.NewInt(0))
	for i:=0;i<len(kco.C_array);i++{
		v_left = v_left.Add(kco.C_array[i].Mult(new(big.Int).Exp(kco.e, big.NewInt(int64(i)) ,nil)))
	}

	v_right := ComMul(com, kco.z_array, kco.s)

	if !v_left.Equal(v_right){
		fmt.Println("Verify KCO failed")
		return false
	}
	fmt.Println("Verify KCO Succeed")
	return true
}
