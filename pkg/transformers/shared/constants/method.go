package constants

var (
	//TODO: get cat, pit, and drip file method signatures directly from the ABI
	biteMethod               = GetSolidityMethodSignature(CatABI, "Bite")
	catFileChopLumpMethod    = "file(bytes32,bytes32,uint256)"
	catFileFlipMethod        = GetSolidityMethodSignature(CatABI, "file")
	catFilePitVowMethod      = "file(bytes32,address)"
	dealMethod               = GetSolidityMethodSignature(FlipperABI, "deal")
	dentMethod               = GetSolidityMethodSignature(FlipperABI, "dent")
	dripDripMethod           = GetSolidityMethodSignature(DripABI, "drip")
	dripFileIlkMethod        = "file(bytes32,bytes32,uint256)"
	dripFileRepoMethod       = GetSolidityMethodSignature(DripABI, "file")
	dripFileVowMethod        = "file(bytes32,bytes32)"
	flapKickMethod           = GetSolidityMethodSignature(FlapperABI, "Kick")
	flipKickMethod           = GetSolidityMethodSignature(FlipperABI, "Kick")
	flopKickMethod           = GetSolidityMethodSignature(FlopperABI, "Kick")
	frobMethod               = GetSolidityMethodSignature(PitABI, "Frob")
	logValueMethod           = GetSolidityMethodSignature(MedianizerABI, "LogValue")
	pitFileDebtCeilingMethod = "file(bytes32,uint256)"
	pitFileIlkMethod         = "file(bytes32,bytes32,uint256)"
	tendMethod               = GetSolidityMethodSignature(FlipperABI, "tend")
	vatFluxMethod            = GetSolidityMethodSignature(VatABI, "flux")
	vatFoldMethod            = GetSolidityMethodSignature(VatABI, "fold")
	vatGrabMethod            = GetSolidityMethodSignature(VatABI, "grab")
	vatHealMethod            = GetSolidityMethodSignature(VatABI, "heal")
	vatInitMethod            = GetSolidityMethodSignature(VatABI, "init")
	vatMoveMethod            = GetSolidityMethodSignature(VatABI, "move")
	vatSlipMethod            = GetSolidityMethodSignature(VatABI, "slip")
	vatTollMethod            = GetSolidityMethodSignature(VatABI, "toll")
	vatTuneMethod            = GetSolidityMethodSignature(VatABI, "tune")
	vowFlogMethod            = GetSolidityMethodSignature(VowABI, "flog")
)
