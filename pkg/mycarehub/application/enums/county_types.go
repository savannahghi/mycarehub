package enums

import (
	"fmt"
	"io"
	"strconv"
)

// CountyType defines various Counties available
type CountyType string

const (
	// CountyTypeMombasa is a county in country type KENYA
	CountyTypeMombasa CountyType = "Mombasa"
	// CountyTypeKwale is a county in country type KENYA
	CountyTypeKwale CountyType = "Kwale"
	// CountyTypeKilifi is a county in country type KENYA
	CountyTypeKilifi CountyType = "Kilifi"
	// CountyTypeTanaRiver is a county in country type KENYA
	CountyTypeTanaRiver CountyType = "Tana_River"
	// CountyTypeLamu is a county in country type KENYA
	CountyTypeLamu CountyType = "Lamu"
	// CountyTypeTaitaTaveta is a county in country type KENYA
	CountyTypeTaitaTaveta CountyType = "Taita_Taveta"
	// CountyTypeGarissa is a county in country type KENYA
	CountyTypeGarissa CountyType = "Garissa"
	// CountyTypeWajir is a county in country type KENYA
	CountyTypeWajir CountyType = "Wajir"
	// CountyTypeMandera is a county in country type KENYA
	CountyTypeMandera CountyType = "Mandera"
	// CountyTypeMarsabit is a county in country type KENYA
	CountyTypeMarsabit CountyType = "Marsabit"
	// CountyTypeIsiolo is a county in country type KENYA
	CountyTypeIsiolo CountyType = "Isiolo"
	// CountyTypeMeru is a county in country type KENYA
	CountyTypeMeru CountyType = "Meru"
	// CountyTypeTharakaNithi is a county in country type KENYA
	CountyTypeTharakaNithi CountyType = "Tharaka_Nithi"
	// CountyTypeEmbu is a county in country type KENYA
	CountyTypeEmbu CountyType = "Embu"
	// CountyTypeKitui is a county in country type KENYA
	CountyTypeKitui CountyType = "Kitui"
	// CountyTypeMachakos is a county in country type KENYA
	CountyTypeMachakos CountyType = "Machakos"
	// CountyTypeMakueni      is a county in country type KENYA
	CountyTypeMakueni CountyType = "Makueni"
	// CountyTypeNyandarua is a county in country type KENYA
	CountyTypeNyandarua CountyType = "Nyandarua"
	// CountyTypeNyeri is a county in country type KENYA
	CountyTypeNyeri CountyType = "Nyeri"
	// CountyTypeKirinyaga is a county in country type KENYA
	CountyTypeKirinyaga CountyType = "Kirinyaga"
	// CountyTypeMuranga is a county in country type KENYA
	CountyTypeMuranga CountyType = "Muranga"
	// CountyTypeKiambu is a county in country type KENYA
	CountyTypeKiambu CountyType = "Kiambu"
	// CountyTypeTurkana is a county in country type KENYA
	CountyTypeTurkana CountyType = "Turkana"
	// CountyTypeWestPokot is a county in country type KENYA
	CountyTypeWestPokot CountyType = "West_Pokot"
	// CountyTypeSamburu is a county in country type KENYA
	CountyTypeSamburu CountyType = "Samburu"
	// CountyTypeTransNzoia is a county in country type KENYA
	CountyTypeTransNzoia CountyType = "Trans_Nzoia"
	// CountyTypeUasinGishu is a county in country type KENYA
	CountyTypeUasinGishu CountyType = "Uasin_Gishu"
	// CountyTypeElgeyoMarakwet is a county in country type KENYA.
	CountyTypeElgeyoMarakwet CountyType = "Elgeyo_Marakwet"
	// CountyTypeNandi is a county in country type KENYA
	CountyTypeNandi CountyType = "Nandi"
	// CountyTypeBaringo is a county in country type KENYA
	CountyTypeBaringo CountyType = "Baringo"
	// CountyTypeLaikipia is a county in country type KENYA
	CountyTypeLaikipia CountyType = "Laikipia"
	// CountyTypeNakuru is a county in country type KENYA
	CountyTypeNakuru CountyType = "Nakuru"
	// CountyTypeNarok is a county in country type KENYA
	CountyTypeNarok CountyType = "Narok"
	// CountyTypeKajiado is a county in country type KENYA
	CountyTypeKajiado CountyType = "Kajiado"
	// CountyTypeKericho is a county in country type KENYA
	CountyTypeKericho CountyType = "Kericho"
	// CountyTypeBomet        is a county in country type KENYA
	CountyTypeBomet CountyType = "Bomet"
	// CountyTypeKakamega is a county in country type KENYA
	CountyTypeKakamega CountyType = "Kakamega"
	// CountyTypeVihiga is a county in country type KENYA
	CountyTypeVihiga CountyType = "Vihiga"
	// CountyTypeBungoma is a county in country type KENYA
	CountyTypeBungoma CountyType = "Bungoma"
	// CountyTypeBusia is a county in country type KENYA
	CountyTypeBusia CountyType = "Busia"
	// CountyTypeSiaya is a county in country type KENYA
	CountyTypeSiaya CountyType = "Siaya"
	// CountyTypeKisumu is a county in country type KENYA
	CountyTypeKisumu CountyType = "Kisumu"
	// CountyTypeHomaBay is a county in country type KENYA
	CountyTypeHomaBay CountyType = "Homa_Bay"
	// CountyTypeMigori is a county in country type KENYA
	CountyTypeMigori CountyType = "Migori"
	// CountyTypeKisii is a county in country type KENYA
	CountyTypeKisii CountyType = "Kisii"
	// CountyTypeNyamira is a county in country type KENYA
	CountyTypeNyamira CountyType = "Nyamira"
	// CountyTypeNairobi is a county in country type KENYA
	CountyTypeNairobi CountyType = "Nairobi"

	// Other counties
)

// KenyanCounties represent all contyTypes of country type KENYA
var KenyanCounties = []CountyType{
	CountyTypeMombasa,
	CountyTypeKwale,
	CountyTypeKilifi,
	CountyTypeTanaRiver,
	CountyTypeLamu,
	CountyTypeTaitaTaveta,
	CountyTypeGarissa,
	CountyTypeWajir,
	CountyTypeMandera,
	CountyTypeMarsabit,
	CountyTypeIsiolo,
	CountyTypeMeru,
	CountyTypeTharakaNithi,
	CountyTypeEmbu,
	CountyTypeKitui,
	CountyTypeMachakos,
	CountyTypeMakueni,
	CountyTypeNyandarua,
	CountyTypeNyeri,
	CountyTypeKirinyaga,
	CountyTypeMuranga,
	CountyTypeKiambu,
	CountyTypeTurkana,
	CountyTypeWestPokot,
	CountyTypeSamburu,
	CountyTypeTransNzoia,
	CountyTypeUasinGishu,
	CountyTypeElgeyoMarakwet,
	CountyTypeNandi,
	CountyTypeBaringo,
	CountyTypeLaikipia,
	CountyTypeNakuru,
	CountyTypeNarok,
	CountyTypeKajiado,
	CountyTypeKericho,
	CountyTypeBomet,
	CountyTypeKakamega,
	CountyTypeVihiga,
	CountyTypeBungoma,
	CountyTypeBusia,
	CountyTypeSiaya,
	CountyTypeKisumu,
	CountyTypeHomaBay,
	CountyTypeMigori,
	CountyTypeKisii,
	CountyTypeNyamira,
	CountyTypeNairobi,
}

// IsValid returns true if a county type is valid
func (e CountyType) IsValid() bool {
	switch e {
	case CountyTypeMombasa,
		CountyTypeKwale,
		CountyTypeKilifi,
		CountyTypeTanaRiver,
		CountyTypeLamu,
		CountyTypeTaitaTaveta,
		CountyTypeGarissa,
		CountyTypeWajir,
		CountyTypeMandera,
		CountyTypeMarsabit,
		CountyTypeIsiolo,
		CountyTypeMeru,
		CountyTypeTharakaNithi,
		CountyTypeEmbu,
		CountyTypeKitui,
		CountyTypeMachakos,
		CountyTypeMakueni,
		CountyTypeNyandarua,
		CountyTypeNyeri,
		CountyTypeKirinyaga,
		CountyTypeMuranga,
		CountyTypeKiambu,
		CountyTypeTurkana,
		CountyTypeWestPokot,
		CountyTypeSamburu,
		CountyTypeTransNzoia,
		CountyTypeUasinGishu,
		CountyTypeElgeyoMarakwet,
		CountyTypeNandi,
		CountyTypeBaringo,
		CountyTypeLaikipia,
		CountyTypeNakuru,
		CountyTypeNarok,
		CountyTypeKajiado,
		CountyTypeKericho,
		CountyTypeBomet,
		CountyTypeKakamega,
		CountyTypeVihiga,
		CountyTypeBungoma,
		CountyTypeBusia,
		CountyTypeSiaya,
		CountyTypeKisumu,
		CountyTypeHomaBay,
		CountyTypeMigori,
		CountyTypeKisii,
		CountyTypeNyamira,
		CountyTypeNairobi:
		return true
	}
	return false
}

// String converts county type to string.
func (e CountyType) String() string {
	return string(e)
}

// UnmarshalGQL converts the supplied value to a county type.
func (e *CountyType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = CountyType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid CountyType", str)
	}
	return nil
}

// MarshalGQL writes the county type to the supplied writer
func (e CountyType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
