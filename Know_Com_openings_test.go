package shuffle

import (
	"fmt"
	"math/big"
	"testing"
)

func TestKCO_proof(t *testing.T) {

	com := Common{EC.G,EC.H}

	rx,ry,rz := big.NewInt(11), big.NewInt(22),big.NewInt(33)
	x, y, z := big.NewInt(2), big.NewInt(3),big.NewInt(6)

	x_array := []*big.Int{x,y,z}
	r_array := []*big.Int{rx,ry,rz}
	proof := KCO_proof(com, x_array, r_array)

	fmt.Println("Verify KCO: ", VerifyKCO(com, proof))
}
