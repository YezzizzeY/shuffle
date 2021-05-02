package shuffle

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"
)

func TestShuffle(t *testing.T) {


	com := Common{EC.G,EC.H}

	rx,ry,rz,rw := big.NewInt(11), big.NewInt(22), big.NewInt(33),big.NewInt(44)
	x, y, z, w := big.NewInt(2), big.NewInt(3), big.NewInt(2), big.NewInt(4)

	X, Y, Z, W := Com(com,x,rx), Com(com,y,ry), Com(com,z,rz), Com(com,w,rw)
	C_array := []ECPoint{X, Y, Z, W}
	fmt.Println("four initial commitments: ")
	for i:=0;i<4;i++{
		fmt.Println(C_array[i])
	}
	pi := []int{2,1,3,0}
	fmt.Println(" Rearranged sequence: ", pi)
	r_array := []*big.Int{}
	for i:=0;i<4;i++{
		r, err := rand.Int(rand.Reader, EC.N)
		check(err)
		r_array = append(r_array, r)
	}

	fmt.Println("four random numbers: ",)
	for i:=0;i<4;i++{
		fmt.Println(r_array[i])
	}
	shuff := Shuffle(com, 4, C_array, pi, r_array)
	fmt.Println("shuffle proof:")
	fmt.Println("previous commitments:")
	for i:=0;i<4;i++{
		fmt.Println(shuff.C_array[i])
	}
	fmt.Println("shuffled commitments:")
	for i:=0;i<4;i++{
		fmt.Println(shuff.C_pie[i])
	}
	fmt.Println("commitments for queue a:", )
	for i:=0;i<4;i++{
		fmt.Println(shuff.CA_array[i])
	}
	fmt.Println("challenges x,y,z: ", shuff.x, shuff.y, shuff.z)
	fmt.Println("commitments for queue b=x^a:", )
	for i:=0;i<4;i++{
		fmt.Println(shuff.CB_array[i])
	}

	if VerifyShuffle(shuff){
		fmt.Println("shuffle completed!")
	}
}
