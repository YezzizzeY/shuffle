package shuffle

import (
	"fmt"
	"math/big"
	"testing"
)

func TestProduct(t *testing.T) {

	com := Common{EC.G,EC.H}
	// prover generate  C
	rx,ry,rz := big.NewInt(11), big.NewInt(22),big.NewInt(33)
	x, y, z := big.NewInt(2), big.NewInt(4),big.NewInt(6)

	X, Y, Z := Com(com,x,rx), Com(com,y,ry), Com(com,z,rz)
	P_product := Product(com, X, x, rx, Y, y, ry, Z, rz)
	fmt.Println("Verify product: ", VerifyProduct(com, P_product))

	P_product_value := ProductValue(com, X,x,rx,Y,y,ry,z)
	fmt.Println("Verify product value: ", VerifyProduct(com, P_product_value))

	P_Multi_product := MultiProduct(com, []ECPoint{X,Y,Z}, []*big.Int{x,y,z}, []*big.Int{rx,ry,rz})
	fmt.Println("Verify multiproduct product value: ", VerifyMultiProduct(com, P_Multi_product, big.NewInt(48)))
}
