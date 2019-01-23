package constants

//TODO: get cat, pit, and drip file method signatures directly from the ABI
func biteMethod() string               { return GetSolidityMethodSignature(CatABI(), "Bite") }
func catFileChopLumpMethod() string    { return "file(bytes32,bytes32,uint256)" }
func catFileFlipMethod() string        { return GetSolidityMethodSignature(CatABI(), "file") }
func catFilePitVowMethod() string      { return "file(bytes32,address)" }
func dealMethod() string               { return GetSolidityMethodSignature(FlipperABI(), "deal") }
func dentMethod() string               { return GetSolidityMethodSignature(FlipperABI(), "dent") }
func dripDripMethod() string           { return GetSolidityMethodSignature(DripABI(), "drip") }
func dripFileIlkMethod() string        { return "file(bytes32,bytes32,uint256)" }
func dripFileRepoMethod() string       { return GetSolidityMethodSignature(DripABI(), "file") }
func dripFileVowMethod() string        { return "file(bytes32,bytes32)" }
func flapKickMethod() string           { return GetSolidityMethodSignature(FlapperABI(), "Kick") }
func flipKickMethod() string           { return GetSolidityMethodSignature(FlipperABI(), "Kick") }
func flopKickMethod() string           { return GetSolidityMethodSignature(FlopperABI(), "Kick") }
func frobMethod() string               { return GetSolidityMethodSignature(PitABI(), "Frob") }
func logValueMethod() string           { return GetSolidityMethodSignature(MedianizerABI(), "LogValue") }
func pitFileDebtCeilingMethod() string { return "file(bytes32,uint256)" }
func pitFileIlkMethod() string         { return "file(bytes32,bytes32,uint256)" }
func tendMethod() string               { return GetSolidityMethodSignature(FlipperABI(), "tend") }
func vatFluxMethod() string            { return GetSolidityMethodSignature(VatABI(), "flux") }
func vatFoldMethod() string            { return GetSolidityMethodSignature(VatABI(), "fold") }
func vatGrabMethod() string            { return GetSolidityMethodSignature(VatABI(), "grab") }
func vatHealMethod() string            { return GetSolidityMethodSignature(VatABI(), "heal") }
func vatInitMethod() string            { return GetSolidityMethodSignature(VatABI(), "init") }
func vatMoveMethod() string            { return GetSolidityMethodSignature(VatABI(), "move") }
func vatSlipMethod() string            { return GetSolidityMethodSignature(VatABI(), "slip") }
func vatTollMethod() string            { return GetSolidityMethodSignature(VatABI(), "toll") }
func vatTuneMethod() string            { return GetSolidityMethodSignature(VatABI(), "tune") }
func vowFlogMethod() string            { return GetSolidityMethodSignature(VowABI(), "flog") }
