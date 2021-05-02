package shuffle

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

type Statement struct {
	C_array []ECPoint
	C_value ECPoint
	CA_array []ECPoint
}

type Multi_proof struct {
	// statement
	Statement
	// message step 2
	CA0 ECPoint
	CB_array []ECPoint
	Ek_array []ECPoint
	x *big.Int
	// calculated message
	a,r,b,s,tao *big.Int
	X_array []*big.Int
}

// Common means two points in a struct, vary from Encrypt value to Ek
type Common struct {
	g1 ECPoint
	g2 ECPoint
}


func Encryptb(com Common, bk *big.Int, r *big.Int) ECPoint {
	g2 := com.g2.Mult(bk).Add(com.g2.Mult(r))
	return g2
}

func Encrypt1(com Common, r *big.Int) ECPoint {
	return com.g2.Mult(r)
}
func GenCom(com Common, v *big.Int, r *big.Int) ECPoint{
	return com.g1.Mult(v).Add(com.g2.Mult(r))
}


func Multi_exp(m int, st Statement, a_array []*big.Int, r_array []*big.Int, rou *big.Int) Multi_proof {

	// a_array index from 0
	// r_array index from 0
	// C_array CA_array index from 0
	if len(a_array)!= m{fmt.Println("vector a length != m")}
	if len(r_array)!= m{fmt.Println("vector r length != m")}
	com := Common{EC.G,EC.H}

	// calculate vector Ci * ai
	// Encrypt rou
	ciai := Encrypt1(com,rou)
	for i:=0;i<m;i++{
		ciai = ciai.Add(st.C_array[i].Mult(a_array[i]))
	}

	// calculate C and statements C1...Cm mult a1...am mult Epk(1;rou)
	if !st.C_value.Equal(ciai){fmt.Println("C != eciai")}
	for i:=0;i<m;i++ {
		if !st.CA_array[i].Equal(GenCom(com, a_array[i], r_array[i])) {
			fmt.Println("nunber", i, "CA comm != GenCom(a,r)")
		}
	}
	// select a0 and r0
	a0, err := rand.Int(rand.Reader, EC.N)
	check(err)

	// form matrix C1*a0 ... C1*am
	//			   ...		 ...
	//			   Cm*a0 ... Cm*am
	Cmatrix := [][]ECPoint{}
	for t:=0;t<m;t++{
		heng := []ECPoint{}
		heng = append(heng, st.C_array[t].Mult(a0))
		for i:=0;i<m;i++{
			heng = append(heng,st.C_array[t].Mult(a_array[i]))
		}
		Cmatrix = append(Cmatrix, heng)
	}

	// select random number r0, b0...b2m-1, s0...s2m-1, tao0...tao2m-1
	r0, err := rand.Int(rand.Reader, EC.N)
	check(err)

	// b_array, s_array, tao_array index: 0~2m-1 total 2m numbers
	b_array := []*big.Int{}
	for k:=0;k<=(2*m-1);k++{
		r, err := rand.Int(rand.Reader, EC.N)
		check(err)
		b_array = append(b_array, r)
	}
	b_array[m] = big.NewInt(0)

	s_array := []*big.Int{}
	for k:=0;k<=(2*m-1);k++{
		r, err := rand.Int(rand.Reader, EC.N)
		check(err)
		s_array = append(s_array, r)
	}
	s_array[m] = big.NewInt(0)

	tao_array := []*big.Int{}
	for k:=0;k<=(2*m-1);k++{
		r, err := rand.Int(rand.Reader, EC.N)
		check(err)
		tao_array = append(tao_array, r)
	}
	tao_array[m] = rou

	// Com(a0,r0)
	CA0 := GenCom(com, a0, r0)

	// initiate CB1...CBm Ek1...Ekm
	CB_array := []ECPoint{}
	Ek_array := []ECPoint{}

	// form CBk and Ek
	for k:=0;k<=(2*m-1);k++ {
		CB_array = append(CB_array, GenCom(com,b_array[k],s_array[k]))
	}

	// a_array1 put a0 in the front, index from 0~m, total m+1
	a_array1 := []*big.Int{}
	a_array1 = append(a_array1,a0)
	a_array1 = append(a_array1,a_array...)

	// C_array1 put 1 in the front, C[0] is not useful C[i] represents Ci in the matrix, total m
	C_array1 := []ECPoint{}
	C_array1 = append(C_array1, EC.G.Mult(big.NewInt(0)))
	C_array1 = append(C_array1, st.C_array...)

	// calculate every Ek and append it to Ek_array, the order is E0, E1, ... E2m-1, should be totaly 2m numbers
	for k:=0;k<=(2*m-1);k++{
		if k<m {
			temp := Encryptb(com, b_array[k],tao_array[k])
			for i:=m-k;i<=m;i++ {
				temp = temp.Add(C_array1[i].Mult(a_array1[k-m+i]))
			}
			Ek_array = append(Ek_array,temp)
		}
		if k==m{
			temp := Encryptb(com, b_array[k],tao_array[k])
			for i:=1;i<=m;i++{
				temp = temp.Add(C_array1[i].Mult(a_array1[k-m+i]))
			}
			Ek_array = append(Ek_array,temp)
		}
		if k>m{
			temp := Encryptb(com, b_array[k],tao_array[k])
			for i:=1;i<=2*m-k;i++{
				temp = temp.Add(C_array1[i].Mult(a_array1[k-m+i]))
			}
			Ek_array = append(Ek_array,temp)
		}
	}

	// ----------------------------- select challenge x, step2----------------------------------------------------
	//TODO this place, x should be midified
	var EkString string
	for i:=0;i<=2*m-1;i++{
		EkString = EkString + Ek_array[i].X.String() + Ek_array[i].Y.String()
	}
	GHString := com.g1.X.String()+com.g1.Y.String()+com.g2.X.String()+com.g2.Y.String()
	c := sha256.Sum256([]byte(EkString+GHString))
	intc := new(big.Int).SetBytes(c[:])
	x := intc

	// x_array: index from 0, total m, x^1 - x^m
	X_array := []*big.Int{}
	for i:=1;i<=m;i++{
		X_array = append(X_array,new(big.Int).Exp(x,big.NewInt(int64(i)),nil))
	}

	// culculate a = a0+ i:1~m, aixi
	a := a0
	for i:=0;i<m;i++{
		a = new(big.Int).Add(a,new(big.Int).Mul(a_array[i],X_array[i]))
	}

	// culculate r = r0+ i:1~m, rixi
	r := r0
	for i:=0;i<m;i++{
		r = new(big.Int).Add(r,new(big.Int).Mul(r_array[i],X_array[i]))
	}

	// culculate b = b0+ bkx^k
	b := b_array[0]
	for k:=1;k<=2*m-1;k++{
		b = new(big.Int).Add(b, new(big.Int).Mul(b_array[k],new(big.Int).Exp(x,big.NewInt(int64(k)),nil)))
	}

	// culculate s = s0+ skx^k
	s := s_array[0]
	for k:=1;k<=2*m-1;k++{
		s = new(big.Int).Add(s, new(big.Int).Mul(s_array[k],new(big.Int).Exp(x,big.NewInt(int64(k)),nil)))
	}

	// culculate tao = tao0+ taokx^k
	tao := tao_array[0]
	for k:=1;k<=2*m-1;k++{
		tao = new(big.Int).Add(tao, new(big.Int).Mul(tao_array[k],new(big.Int).Exp(x,big.NewInt(int64(k)),nil)))
	}

	// now a_array, r, b, s, tao are all generated the next thing is to verify them
	result := Multi_proof{}
	result.Statement = st
	result.CA0 = CA0
	result.CB_array = CB_array
	result.Ek_array = Ek_array
	result.x = x
	result.a = a
	result.r = r
	result.b = b
	result.s = s
	result.tao = tao
	result.X_array = X_array
	fmt.Println("Multi_exp generated")
	return result
}

func VerifyMulExp(mul Multi_proof) bool {

	com := Common{EC.G,EC.H}
	m := len(mul.Statement.CA_array)
	// first CB=Com_ck(0;0)
	if !mul.CB_array[m].Equal(GenCom(com,big.NewInt(0),big.NewInt(0))){
		fmt.Println("Verify Multi-exponentiation Failed!  CB[m] is wrong")
	}
	fmt.Println("Verify Multi-exponentiation step 1 Succeed")
	if !mul.Ek_array[m].Equal(mul.Statement.C_value){
		fmt.Println("Verify Multi-exponentiation Failed! Em!=C")
	}
	fmt.Println("Verify Multi-exponentiation step 2 Succeed")
	// CA0CA^x and com_ck(a,r)
	ver3 := mul.CA0
	for i:=0;i<m;i++{
		ver3 = ver3.Add( mul.CA_array[i].Mult(mul.X_array[i]))
	}
	if !ver3.Equal(GenCom(com,mul.a,mul.r)){
		fmt.Println("Verify Multi-exponentiation Failed! verify 3 false")
	}
	fmt.Println("Verify Multi-exponentiation step 3 Succeed")
	// Verify CB, also verify step 4
	ver4 := mul.CB_array[0]
	for k := 1; k <= 2*m-1; k++ {
		ver4 = ver4.Add(mul.CB_array[k].Mult(new(big.Int).Exp(mul.x,big.NewInt(int64(k)),nil)))
	}
	if !ver4.Equal(GenCom(com,mul.b,mul.s)){
		fmt.Println("Verify Multi-exponentiation Failed! verify 4 false")
	}
	fmt.Println("Verify Multi-exponentiation step 4 Succeed")
	// verify 5
	// first calculate verify5_left
	ver5 := mul.Ek_array[0]
	for k:=1;k<=2*m-1;k++ {
		ver5 = ver5.Add(mul.Ek_array[k].Mult(new(big.Int).Exp(mul.x,big.NewInt(int64(k)),nil)))
	}

	// next calculate verify5_right
	ver5_right := Encryptb(com,mul.b,mul.tao)
	for i:=1; i<=m; i++ {
		ver5_right = ver5_right.Add(mul.C_array[i-1].Mult(new(big.Int).Mul(new(big.Int).Exp(mul.x,new(big.Int).Sub(big.NewInt(int64(m)),big.NewInt(int64(i))),nil),mul.a)))
	}
	if !ver5.Equal(ver5_right){
		fmt.Println("Verify Multi-exponentiation Failed! ver5 false")
	}
	fmt.Println("Verify Multi-exponentiation step 5 Succeed")
	fmt.Println("Verify Multi-exponentiation Succeed")
	return true
}