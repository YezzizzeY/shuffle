package shuffle

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

func ComZero(com Common, m int, z *big.Int) []ECPoint{
	re := []ECPoint{}
	for i:=0;i<m;i++ {
		re = append(re, com.g1.Mult(z).Neg())
	}
	return re
}

type Shu_proof struct {
	C_array []ECPoint
	C_pie []ECPoint
	CA_array []ECPoint
	x *big.Int
	CB_array []ECPoint
	y,z *big.Int
	kco KCO
	pro []ProductProof
	mul Multi_proof
}

// attention: pi should contain zero
func Shuffle(com Common, m int, C_array []ECPoint, pi []int, rou_array []*big.Int) Shu_proof {

	C_pie := []ECPoint{}
	for i:=0;i<m;i++{
		C_pie = append(C_pie, C_array[pi[i]].Add(Encrypt1(com, rou_array[i])))
	}

	// initiate vector r
	r_array := []*big.Int{}
	for i:=0;i<m;i++{
		r, err := rand.Int(rand.Reader, EC.N)
		check(err)
		r_array = append(r_array, r)
	}

	// initiate vector a
	a_array := []*big.Int{}
	for i:=0;i<m;i++{
		a_array = append(a_array, big.NewInt(int64(pi[i])))
	}

	// calculate vector CA
	CA_array := []ECPoint{}
	for i:=0;i<m;i++{
		CA_array = append(CA_array, Com(com, a_array[i], r_array[i]))
	}

	//TODO: challenge x is supposed to be modified
	var CAString string
	for i:=0;i<m;i++{
		CAString = CAString + CA_array[i].X.String() + CA_array[i].Y.String()
	}
	GHString := com.g1.X.String()+com.g1.Y.String()+com.g2.X.String()+com.g2.Y.String()
	c := sha256.Sum256([]byte(CAString+GHString))
	intc := new(big.Int).SetBytes(c[:])
	x := intc

	// initiate vector s
	s_array := []*big.Int{}
	for i:=0;i<m;i++{
		s, err := rand.Int(rand.Reader, EC.N)
		check(err)
		s_array = append(s_array, s)
	}

	// calculate vector x
	x_array := []*big.Int{}
	for i:=0;i<m;i++{
		x_array = append(x_array, new(big.Int).Exp(x, big.NewInt(int64(i+1)),nil))
	}

	// initiate vector b
	b_array := []*big.Int{}
	for i:=0;i<m;i++{
		b_array = append(b_array, x_array[pi[i]])
	}

	// initiate vector CB
	CB_array := []ECPoint{}
	for i:=0;i<m;i++{
		CB_array = append(CB_array, Com(com, b_array[i], s_array[i]))
	}

	// TODO: challenges y,z should be modified
	var CBString string
	for i:=0;i<m;i++{
		CBString = CBString + CB_array[i].X.String() + CB_array[i].Y.String()
	}
	y_s := sha256.Sum256([]byte(CAString+CBString+GHString))
	z_S := sha256.Sum256([]byte(CAString+CBString+GHString))
	y := new(big.Int).SetBytes(y_s[:])
	z := new(big.Int).SetBytes(z_S[:])
	// z_array
	z_array := []*big.Int{}
	for i:=0;i<m;i++{
		z_array = append(z_array, new(big.Int).Neg(z))
	}

	// calculate vector CD
	CD_array := []ECPoint{}
	for i:=0;i<m;i++{
		CD_array = append(CD_array, CA_array[i].Mult(y).Add(CB_array[i]))
	}

	// calculate vector d
	d_array := []*big.Int{}
	for i:=0;i<m;i++{
		d_array = append(d_array, new(big.Int).Add(new(big.Int).Mul(y,new(big.Int).Add(big.NewInt(int64(1)),a_array[i])),b_array[i]))
	}

	// calculate vector t
	t_array := []*big.Int{}
	for i:=0;i<m;i++{
		t_array = append(t_array, new(big.Int).Add(new(big.Int).Mul(y,r_array[i]),s_array[i]))
	}

	// calculate rou
	// rou here is -rou, it will be calculated in multi_exp
	rou := big.NewInt(0)
	for i:=0;i<m;i++{
		rou = new(big.Int).Add(rou,new(big.Int).Mul(rou_array[i], b_array[i]))
	}
	rou = new(big.Int).Neg(rou)

	// Product Argument: include Know_comments_openings argument and Open_products argument
	// form vector CD*C_-z
	// first calculate vector d-z and generate kco proof
	d_z_array := []*big.Int{}
	for i:=0;i<m;i++{
		d_z_array = append(d_z_array, new(big.Int).Add(d_array[i], z_array[i]))
	}

	Cd_Cz_array := []ECPoint{}
	for i:=0;i<m;i++{
		Cd_Cz_array = append(Cd_Cz_array, Com(com, d_z_array[i], t_array[i]))
	}

	kco := KCO_proof(com, d_z_array, t_array)

	// generate product proof
	prod := MultiProduct(com, Cd_Cz_array, d_z_array, t_array)


	// Multi-exponentiation Argument
	// first calculate C_value
	C_value := EC.G.Mult(big.NewInt(0))
	for i:=0;i<m;i++{
		n := C_array[i].Mult(x_array[i])
		C_value = C_value.Add(n)
	}

	// second form statement, C'[i] = C_array[pi]*E(1,rou[i])
	r_pie_array := []*big.Int{}
	for i:=0;i<m;i++{
		r_pie_array = append(r_pie_array, new(big.Int).Mul(rou_array[i], r_array[i]))
	}

	st := Statement{
		C_array: C_pie,
		C_value: C_value,
		CA_array: CB_array,
	}

	mulp := Multi_exp(m, st, b_array, s_array, rou)

	shu := Shu_proof{
		C_array: C_array,
		C_pie: C_pie,
		CA_array: CA_array,
		x: x,
		CB_array: CB_array,
		y :y,
		z: z,
		kco: kco,
		pro: prod,
		mul: mulp,
	}
	return shu
}

func VerifyShuffle(shu Shu_proof) bool {
	com := Common{EC.G,EC.H}
	// first Verify KCO
	if !VerifyKCO(com, shu.kco){
		return false
	}

	// calculaete i=m : y*i+x^i-z
	m := len(shu.C_array)
	v := big.NewInt(1)
	for i:=0;i<m;i++{
		// y*i
		yi := new(big.Int).Mul(shu.y, big.NewInt(int64(i+1)))
		// x^i
		xi := new(big.Int).Exp(shu.x, big.NewInt(int64(i+1)), nil)
		// y*i + x^i
		yx := new(big.Int).Add(yi,xi)
		// final
		f := new(big.Int).Sub(yx, shu.z)
		v = new(big.Int).Mul(v, f)
	}
	if !VerifyMultiProduct(com, shu.pro, v){
		return false
	}

	if !VerifyMulExp(shu.mul){
		return false
	}
	fmt.Println("Verify Shuffle Succeed")
	return true
}