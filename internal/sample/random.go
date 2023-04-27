package sample

import (
	"math"
	"math/rand"
	"time"

	"github.com/ryanreadbooks/go-grpc-example/pb"
)

func init() {
	rand.Seed(time.Now().Unix())
}

var (
	cameraBrands            = []string{"Sony", "Leica", "Canon", "Nikon"}
	cameraSpecs             = []string{"12MP+20MP", "10MP+20MP", "20MP+45MP", "10MP+15MP"}
	screenResolutions       = []string{"1920x1080", "2560x1440", "1280x720", "1334x750"}
	manufacturers           = []string{"Snapdragon", "Nvidia", "Intel", "Samsung", "MediaTek"}
	operatingSystems        = []string{"Android", "iOS", "Windows Mobile", "ColorOS", "MIUI", "HarmonyOS", "OriginOS"}
	operatingSystemVersions = []string{"Stable", "Dev", "Beta", "Insider"}
	cellPhoneBrands         = []string{"Apple", "Samsung", "Huawei", "Xiaomi", "OPPO", "VIVO", "Honor", "Pixel"}
)

func randStringFromSlice(candidates []string) string {
	return candidates[rand.Intn(len(candidates))]
}

func randomCellphoneBrand() string {
	return randStringFromSlice(cellPhoneBrands)
}

func randomCameraBrand() string {
	return randStringFromSlice(cameraBrands)
}

func randomCameraSpec() string {
	return randStringFromSlice(cameraSpecs)
}

func randomFloat64(min, max float64) float64 {
	f := rand.Float64()
	return f*(max-min) + min
}

func randomScreenResolution() string {
	return randStringFromSlice(screenResolutions)
}

func randomScreenSize() float64 {
	return randomFloat64(4.4, 7.0)
}

func randomManufacturer() string {
	return randStringFromSlice(manufacturers)
}

func randomPowerOfTwo(max int) int32 {
	return (int32)(math.Pow(2, float64(rand.Intn(max+1))))
}

func randomInt32(min, max int32) int32 {
	return rand.Int31n(max-min) + min
}

func randomMemoryUnit() pb.Unit {
	idx := rand.Intn(3)
	switch idx {
	case 0:
		return pb.Unit_UnitMB
	case 1:
		return pb.Unit_UnitGB
	case 2:
		return pb.Unit_UnitTB
	default:
		return pb.Unit_UnitGB
	}
}

func randomStorageType() pb.StorageType {
	idx := rand.Intn(2)
	if idx == 0 {
		return pb.StorageType_SSD
	}
	return pb.StorageType_HDD
}

func randomDDRType() pb.DDRType {
	idx := rand.Intn(3)
	switch idx {
	case 0:
		return pb.DDRType_DDR3
	case 1:
		return pb.DDRType_DDR4
	case 2:
		return pb.DDRType_DDR5
	default:
		return pb.DDRType_DDR4
	}
}

func randomOperatingSystem() string {
	return randStringFromSlice(operatingSystems)
}

func randomOperatingSystemVersion() string {
	return randStringFromSlice(operatingSystemVersions)
}
