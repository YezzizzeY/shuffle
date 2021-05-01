package shuffle

import (
	"fmt"
	"math/big"
	"testing"
)

func TestShuffle(t *testing.T) {


	com := Common{EC.G,EC.H}

	rx,ry,rz := big.NewInt(11), big.NewInt(22), big.NewInt(33)
	x, y, z := big.NewInt(2), big.NewInt(3), big.NewInt(4)

	X, Y, Z := Com(com,x,rx), Com(com,y,ry), Com(com,z,rz)
	C_array := []ECPoint{X, Y, Z}
	pi := []int{2,1,0}
	rou_array := []*big.Int{big.NewInt(2),big.NewInt(4),big.NewInt(3)}

	if VerifyShuffle(Shuffle(com, 3, C_array, pi, rou_array)){
		fmt.Println("shuffle completed!")
	}
}
