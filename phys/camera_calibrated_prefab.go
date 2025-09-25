package phys


var IntrinsicsFireflyDLComputar16mm = CameraIntrinsics{
	Width:  1440,
	Height: 1080,
	Fx:     4736.71083258,
	Fy:     4743.03975965,
	Cx:     770.21866744,
	Cy:     483.50827668,
	K1:     -0.09261328,
	K2:     -1.43023836,
	P1:     -0.00215911,
	P2:     -0.00187976,
	K3:     188.33757455,
	// K4..K6 are zero for this 5-parameter fit.
}

var IntrinsicsFireflyDLComputar12mm = CameraIntrinsics{
	Width:  1440,
	Height: 1080,
	Fx:     3613.49651386,
	Fy:     3617.43390846,
	Cx:     837.17440873,
	Cy:     412.78087519,
	K1:     -0.18060152,
	K2:     2.48103332,
	P1:     -0.00571920,
	P2:     0.00121639,
	K3:     -29.84726761,
	// K4..K6 are zero for this 5-parameter fit.
}

var IntrinsicsFireflyDLGeneric6mm = CameraIntrinsics{
	Width:  1440,
	Height: 1080,
	Fx:     1804.17453167,
	Fy:     1804.69144616,
	Cx:     756.49974101,
	Cy:     481.63486915,
	K1:     -0.50722235,
	K2:     0.44907698,
	P1:     0.00151234,
	P2:     -0.00094105,
	K3:     -0.72605770,
	// K4..K6 are zero for this 5-parameter fit.
}
