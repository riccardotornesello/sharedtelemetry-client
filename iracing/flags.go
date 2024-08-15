package iracing

import "sharedtelemetry/client/common"

type IRSDKFlagBitmap uint32

const (
	IRSDKFlagCheckered     IRSDKFlagBitmap = 0x00000001
	IRSDKFlagWhite         IRSDKFlagBitmap = 0x00000002
	IRSDKFlagGreen         IRSDKFlagBitmap = 0x00000004
	IRSDKFlagYellow        IRSDKFlagBitmap = 0x00000008
	IRSDKFlagRed           IRSDKFlagBitmap = 0x00000010
	IRSDKFlagBlue          IRSDKFlagBitmap = 0x00000020
	IRSDKFlagDebris        IRSDKFlagBitmap = 0x00000040
	IRSDKFlagCrossed       IRSDKFlagBitmap = 0x00000080
	IRSDKFlagYellowWaving  IRSDKFlagBitmap = 0x00000100
	IRSDKFlagOneLapToGreen IRSDKFlagBitmap = 0x00000200
	IRSDKFlagGreenHeld     IRSDKFlagBitmap = 0x00000400
	IRSDKFlagTtenToGo      IRSDKFlagBitmap = 0x00000800
	IRSDKFlagFiveToGo      IRSDKFlagBitmap = 0x00001000
	IRSDKFlagRandomWaving  IRSDKFlagBitmap = 0x00002000
	IRSDKFlagCaution       IRSDKFlagBitmap = 0x00004000
	IRSDKFlagCautionWaving IRSDKFlagBitmap = 0x00008000
	IRSDKFlagBlack         IRSDKFlagBitmap = 0x00010000
	IRSDKFlagDisqualify    IRSDKFlagBitmap = 0x00020000
	IRSDKFlagServicible    IRSDKFlagBitmap = 0x00040000
	IRSDKFlagFurled        IRSDKFlagBitmap = 0x00080000
	IRSDKFlagRepair        IRSDKFlagBitmap = 0x00100000
	IRSDKFlagStartHidden   IRSDKFlagBitmap = 0x10000000
	IRSDKFlagStartReady    IRSDKFlagBitmap = 0x20000000
	IRSDKFlagStartSet      IRSDKFlagBitmap = 0x40000000
	IRSDKFlagStartGo       IRSDKFlagBitmap = 0x80000000
)

var IRSDKFlagNames = map[IRSDKFlagBitmap]common.Flag{
	IRSDKFlagCheckered:     common.FlagCheckered,
	IRSDKFlagWhite:         common.FlagWhite,
	IRSDKFlagGreen:         common.FlagGreen,
	IRSDKFlagYellow:        common.FlagYellow,
	IRSDKFlagRed:           common.FlagRed,
	IRSDKFlagBlue:          common.FlagBlue,
	IRSDKFlagDebris:        common.FlagDebris,
	IRSDKFlagCrossed:       common.FlagCrossed,
	IRSDKFlagYellowWaving:  common.FlagYellowWaving,
	IRSDKFlagOneLapToGreen: common.FlagOneLapToGreen,
	IRSDKFlagGreenHeld:     common.FlagGreenHeld,
	IRSDKFlagTtenToGo:      common.FlagTtenToGo,
	IRSDKFlagFiveToGo:      common.FlagFiveToGo,
	IRSDKFlagRandomWaving:  common.FlagRandomWaving,
	IRSDKFlagCaution:       common.FlagCaution,
	IRSDKFlagCautionWaving: common.FlagCautionWaving,
	IRSDKFlagBlack:         common.FlagBlack,
	IRSDKFlagDisqualify:    common.FlagDisqualify,
	IRSDKFlagServicible:    common.FlagServicible,
	IRSDKFlagFurled:        common.FlagFurled,
	IRSDKFlagRepair:        common.FlagRepair,
	IRSDKFlagStartHidden:   common.FlagStartHidden,
	IRSDKFlagStartReady:    common.FlagStartReady,
	IRSDKFlagStartSet:      common.FlagStartSet,
	IRSDKFlagStartGo:       common.FlagStartGo,
}

func fetchFlags(flags uint32) []common.Flag {
	var extractedFlags []common.Flag

	for flag, name := range IRSDKFlagNames {
		if flags&uint32(flag) != 0 {
			extractedFlags = append(extractedFlags, name)
		}
	}

	return extractedFlags
}
