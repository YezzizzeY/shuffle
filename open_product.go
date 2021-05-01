package shuffle

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type ProductProof struct {
	z *big.Int
	X, Y, Z ECPoint
	arefa, beita, derta ECPoint
	c, z1, z2, z3, z4, z5 *big.Int
}

// implement of part of <Hyrax: Doubly-efficient zkSnarks without trusted setup> 2018
func Product(com Common, X ECPoint,x *big.Int, rx *big.Int, Y ECPoint, y *big.Int, ry *big.Int, Z ECPoint, rz *big.Int) ProductProof {
	// first Prover picks b1... b5
	b1, err := rand.Int(rand.Reader, EC.N)
	check(err)
	b2, err := rand.Int(rand.Reader, EC.N)
	check(err)
	b3, err := rand.Int(rand.Reader, EC.N)
	check(err)
	b4, err := rand.Int(rand.Reader, EC.N)
	check(err)
	b5, err := rand.Int(rand.Reader, EC.N)
	check(err)

	arefa := Com(com, b1, b2)
	beita := Com(com, b3, b4)
	derta := Com(Common{X,com.g2}, b3, b5)

	//TODO: challenge c shuold be hash()
	c, err := rand.Int(rand.Reader, EC.N)
	check(err)

	// calculate z1...z5
	z1 := new(big.Int).Add(b1, new(big.Int).Mul(c,x))
	z2 := new(big.Int).Add(b2, new(big.Int).Mul(c,rx))
	z3 := new(big.Int).Add(b3, new(big.Int).Mul(c,y))
	z4 := new(big.Int).Add(b4, new(big.Int).Mul(c,ry))
	z5 := new(big.Int).Add(b5, new(big.Int).Mul(c,new(big.Int).Sub(rz, new(big.Int).Mul(rx,y))))

	P := ProductProof{}
	P.X, P.Y, P.Z = X, Y, Z
	P.arefa, P.beita, P.derta = arefa, beita, derta
	P.z1, P.z2, P.z3, P.z4, P.z5 = z1,z2,z3,z4,z5
	P.c = c
	P.z = new(big.Int).Mul(x,y)
	return P
}

func Com(com Common, x *big.Int, r *big.Int) ECPoint {
	return com.g1.Mult(x).Add(com.g2.Mult(r))
}

func VerifyProduct(com Common, P ProductProof) bool {
	if !P.arefa.Add(P.X.Mult(P.c)).Equal(Com(com, P.z1, P.z2)){
		fmt.Println("arefa verify failed")
		return false
	}
	if !P.beita.Add(P.Y.Mult(P.c)).Equal(Com(com, P.z3, P.z4)){
		fmt.Println("beita verify failed")
		return false
	}
	if !P.derta.Add(P.Z.Mult(P.c)).Equal(Com(Common{P.X,com.g2}, P.z3, P.z5)){
		fmt.Println("derta verify failed")
		return false
	}
	return true
}

func ProductValue(com Common, X ECPoint, x *big.Int, rx *big.Int, Y ECPoint, y *big.Int, ry *big.Int, z *big.Int) ProductProof {
	rz, err := rand.Int(rand.Reader, EC.N)
	check(err)
	Z := Com(com, z, rz)
	return Product(com, X, x, rx, Y, y, ry, Z, rz)
}

// extension CA...CN, z such that a*b*...*n = z
// the length of Q v r must be the same
func MultiProduct(com Common, Q []ECPoint, v []*big.Int, r []*big.Int) []ProductProof{
	m := len(Q)

	re := []ProductProof{}
	for i:=0;i<m-1;i++{
		p := ProductValue(com, Q[i], v[i], r[i], Q[i+1], v[i+1], r[i+1], new(big.Int).Mul(v[i],v[i+1]))
		re = append(re, p)
		if i==m-2{
			q := ProductValue(com, Q[0], v[0], r[0], Q[m-1], v[m-1], r[m-1], new(big.Int).Mul(v[0],v[m-1]))
			re = append(re, q)
		}
	}
	fmt.Println("multiproduct generated")
	return re
}

func VerifyMultiProduct(com Common, P_array []ProductProof, z *big.Int) bool {
	t := big.NewInt(1)
	for i:=0;i<len(P_array);i++{
		if !VerifyProduct(com, P_array[i]) {
			fmt.Println("Verify MultiProduct step1 failed")
			return false
		}
		t = new(big.Int).Mul(t, P_array[i].z)
	}
	fmt.Println("Verify MultiProduct step1 Succeed")
	if sqrt(*t).Cmp(z)!=0{
		fmt.Println("Verify MultiProduct step2 Failed")
		return false
	}
	fmt.Println("Verify MultiProduct Succeed")
	return true
}

func sqrt(n big.Int) *big.Int {
	var  a, b, m, m2 big.Int


	a.SetInt64(int64(1))
	b.Set(&n)

	for {
		m.Add(&a, &b).Div(&m, big.NewInt(2))

		if m.Cmp(&a) == 0 || m.Cmp(&b) == 0 {
			break
		}

		m2.Mul(&m, &m)
		if m2.Cmp(&n) > 0 {
			b.Set(&m)
		} else {
			a.Set(&m)
		}
	}

	return &m
}
