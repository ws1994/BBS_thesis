package main

import (
	"fmt"
	"pbc"
)

func main() {
	// In a real application, generate this once and publish it
	params := pbc.GenerateA(160, 512)

	pairing := params.NewPairing()

	// Initialize group elements. pbc automatically handles garbage collection.
	g := pairing.NewG1()
	s := pairing.NewZr()
	pks := pairing.NewG1()
	a := pairing.NewZr()
	pkuLeft := pairing.NewG1()
	pkuRight := pairing.NewG1()
	h1t := pairing.NewG1()
	timeKey := pairing.NewG1()
	r := pairing.NewZr()
	rasg := pairing.NewG1()
	rg := pairing.NewG1()
	K := pairing.NewGT()
	K2Temp := pairing.NewGT()
	K2 := pairing.NewGT()
	check1 := pairing.NewGT()
	check2 := pairing.NewGT()

	// check3 := pairing.NewGT()
	// check4 := pairing.NewGT()
	// tempcheck := pairing.NewG1()

	//*****
	// Generate server keys
	//*****
	// secret key
	s.Rand()
	fmt.Printf("s = %s\n", s)
	// public key: pks = sG
	g.Rand()
	fmt.Printf("G = %s\n", g)
	pks = pks.MulZn(g, s)
	fmt.Printf("sG = %s\n", pks)

	//*****
	// Generate User keys
	//*****
	// secret key
	a.Rand()
	fmt.Printf("a = %s\n", a)
	// public key; pkuLeft=ag; pkuRight=asg
	pkuLeft = pkuLeft.MulZn(g, a)
	fmt.Printf("aG = %s\n", pkuLeft)
	pkuRight = pkuRight.MulZn(pks, a)
	fmt.Printf("asG = %s\n", pkuRight)

	//*****
	// time-bound key update
	//*****
	h1t = h1t.SetFromHash([]byte("20220610"))
	fmt.Printf("H1(t) = %s\n", h1t)
	timeKey = timeKey.MulZn(h1t, s)
	fmt.Printf("sH1(t) = %s\n", timeKey)

	//*****
	// encryption
	//*****
	// step1: check
	check1 = check1.Pair(pkuLeft, pks)
	check2 = check2.Pair(g, pkuRight)
	check := check1.Equals(check2)
	if check {
		fmt.Printf("pass check!\n")
	} else {
		fmt.Printf("check failed!\n")
	}
	// step2: random r
	r.Rand()
	fmt.Printf("r = %s\n", r)
	rg = rg.MulZn(g, r)
	fmt.Printf("rg = %s\n", rg)
	rasg = rasg.MulZn(pkuRight, r)
	fmt.Printf("rasg = %s\n", rasg)

	// step3: compute encryption key
	K = K.Pair(rasg, h1t)
	fmt.Printf("K = %s\n", K)

	// check3 = check3.Pair(h1t, g)
	// fmt.Printf("e(g,h1t) = %s\n", check3)
	// check4 = check4.PowZn(check3, r)
	// fmt.Printf("()^r = %s\n", check4)
	// check4 = check4.PowZn(check4, a)
	// fmt.Printf("()^ra = %s\n", check4)
	// check4 = check4.PowZn(check4, s)
	// fmt.Printf("()^ras = %s\n", check4)

	// step4: encrypt message
	n := K.BytesLen()
	fmt.Printf("length of H2(k): %d \n", n)
	h2k := K.Bytes()
	fmt.Printf("H2(K) is %s \n", h2k)

	message := []byte("the solution is abcd the solution is abcd the solution is abcd the solution is abcd")

	fmt.Printf("message is: %s \n", message)
	mLength := len(message)
	fmt.Printf("length of message: %d \n", mLength)

	ciphertext := make([]byte, mLength)
	for i := 0; i < mLength; i++ {
		ciphertext[i] = message[i] ^ h2k[i%n]
	}
	fmt.Printf("ciphertext is: %s \n", ciphertext)
	ciphtLength := len(ciphertext)
	fmt.Printf("length of ciphertext is: %d \n", ciphtLength)

	//*****
	// decryption
	//*****
	// calculate k'
	K2Temp = K2Temp.Pair(rg, timeKey)
	K2 = K2.PowZn(K2Temp, a)
	fmt.Printf("K2 = %s\n", K2)

	// recover message
	message2 := make([]byte, mLength)

	n = K2.BytesLen()
	h2k2 := K2.Bytes()
	fmt.Printf("length of H2(k2): %d \n", n)
	fmt.Printf("H2(K2) is %s \n", h2k2)

	for i := 0; i < mLength; i++ {
		message2[i] = ciphertext[i] ^ h2k2[i%n]
	}
	mLength = len(message2)
	fmt.Printf("length of message2 is: %d \n", mLength)
	fmt.Printf("message2 is: %s \n", message2)
}
