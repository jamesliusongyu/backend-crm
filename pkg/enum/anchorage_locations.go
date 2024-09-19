package enum

type AnchorageLocation string

const (
	CHANGI_BARGE_TEMPORARY_HOLDING_ANCHORAGE     AnchorageLocation = "Changi Barge Temporary Holding Anchorage (ACBTH)"
	CHANGI_GENERAL_PURPOSES_ANCHORAGE            AnchorageLocation = "Changi General Purposes Anchorage (ACGP)"
	MAN_OF_WAR_ANCHORAGE                         AnchorageLocation = "Man-of-War Anchorage (AMOW)"
	EASTERN_BUNKERING_A_ANCHORAGE                AnchorageLocation = "Eastern Bunkering A Anchorage (AEBA)"
	EASTERN_BUNKERING_B_ANCHORAGE                AnchorageLocation = "Eastern Bunkering B Anchorage (AEBB)"
	EASTERN_PETROLEUM_C_ANCHORAGE                AnchorageLocation = "Eastern Petroleum C Anchorage (AEPBC)"
	SMALL_CRAFT_B_ANCHORAGE                      AnchorageLocation = "Small Craft B Anchorage (ASCB)"
	SMALL_CRAFT_A_ANCHORAGE                      AnchorageLocation = "Small Craft A Anchorage (ASCA)"
	EASTERN_PETROLEUM_B_ANCHORAGE                AnchorageLocation = "Eastern Petroleum B Anchorage (AEPBB)"
	EASTERN_SPECIAL_PURPOSES_A_ANCHORAGE         AnchorageLocation = "Eastern Special Purposes A Anchorage (AESPA)"
	EASTERN_BUNKERING_C_ANCHORAGE                AnchorageLocation = "Eastern Bunkering C Anchorage (AEBC)"
	EASTERN_HOLDING_A_ANCHORAGE                  AnchorageLocation = "Eastern Holding A Anchorage (AEHA)"
	EASTERN_PETROLEUM_A_ANCHORAGE                AnchorageLocation = "Eastern Petroleum A Anchorage (AEPA)"
	EASTERN_ANCHORAGE                            AnchorageLocation = "Eastern Anchorage (AEW)"
	EASTERN_HOLDING_B_ANCHORAGE                  AnchorageLocation = "Eastern Holding B Anchorage (AEHB)"
	EASTERN_HOLDING_C_ANCHORAGE                  AnchorageLocation = "Eastern Holding C Anchorage (AEHC)"
	WESTERN_QUARANTINE_AND_IMMIGRATION_ANCHORAGE AnchorageLocation = "Western Quarantine and Immigration Anchorage (AWQI)"
	WESTERN_ANCHORAGE                            AnchorageLocation = "Western Anchorage (AWW)"
	WESTERN_PETROLEUM_A_ANCHORAGE                AnchorageLocation = "Western Petroleum A Anchorage (AWPA)"
	WESTERN_HOLDING_ANCHORAGE                    AnchorageLocation = "Western Holding Anchorage (AWH)"
	WESTERN_PETROLEUM_B_ANCHORAGE                AnchorageLocation = "Western Petroleum B Anchorage (AWPB)"
	RAFFLES_RESERVED_ANCHORAGE                   AnchorageLocation = "Raffles Reserved Anchorage (ARAFR)"
	RAFFLES_PETROLEUM_ANCHORAGE                  AnchorageLocation = "Raffles Petroleum Anchorage (ARP)"
	SELAT_PAUH_ANCHORAGE                         AnchorageLocation = "Selat Pauh Anchorage (ASPLU)"
	SELAT_PAUH_PETROLEUM_ANCHORAGE               AnchorageLocation = "Selat Pauh Petroleum Anchorage (ASPP)"
	SUDONG_PETROLEUM_HOLDING_ANCHORAGE           AnchorageLocation = "Sudong Petroleum Holding Anchorage (ASPH)"
	SUDONG_EXPLOSIVE                             AnchorageLocation = "Sudong Explosive (ASUEX)"
	SUDONG_SPECIAL_PURPOSE                       AnchorageLocation = "Sudong Special Purpose (ASSPU)"
	SUDONG_HOLDING_ANCHORAGE                     AnchorageLocation = "Sudong Holding Anchorage (ASH)"
	VERY_LARGE_CRUDE_CARRIER_ANCHORAGE           AnchorageLocation = "Very Large Crude Carrier Anchorage (AVLCC)"
)

var AnchorageLocations = []AnchorageLocation{
	CHANGI_BARGE_TEMPORARY_HOLDING_ANCHORAGE,
	CHANGI_GENERAL_PURPOSES_ANCHORAGE,
	MAN_OF_WAR_ANCHORAGE,
	EASTERN_BUNKERING_A_ANCHORAGE,
	EASTERN_BUNKERING_B_ANCHORAGE,
	EASTERN_PETROLEUM_C_ANCHORAGE,
	SMALL_CRAFT_B_ANCHORAGE,
	SMALL_CRAFT_A_ANCHORAGE,
	EASTERN_PETROLEUM_B_ANCHORAGE,
	EASTERN_SPECIAL_PURPOSES_A_ANCHORAGE,
	EASTERN_BUNKERING_C_ANCHORAGE,
	EASTERN_HOLDING_A_ANCHORAGE,
	EASTERN_PETROLEUM_A_ANCHORAGE,
	EASTERN_ANCHORAGE,
	EASTERN_HOLDING_B_ANCHORAGE,
	EASTERN_HOLDING_C_ANCHORAGE,
	WESTERN_QUARANTINE_AND_IMMIGRATION_ANCHORAGE,
	WESTERN_ANCHORAGE,
	WESTERN_PETROLEUM_A_ANCHORAGE,
	WESTERN_HOLDING_ANCHORAGE,
	WESTERN_PETROLEUM_B_ANCHORAGE,
	RAFFLES_RESERVED_ANCHORAGE,
	RAFFLES_PETROLEUM_ANCHORAGE,
	SELAT_PAUH_ANCHORAGE,
	SELAT_PAUH_PETROLEUM_ANCHORAGE,
	SUDONG_PETROLEUM_HOLDING_ANCHORAGE,
	SUDONG_EXPLOSIVE,
	SUDONG_SPECIAL_PURPOSE,
	SUDONG_HOLDING_ANCHORAGE,
	VERY_LARGE_CRUDE_CARRIER_ANCHORAGE,
}

// GetAnchorageLocations returns a slice of all anchorage locations
func GetAnchorageLocations() []string {
	locations := make([]string, len(AnchorageLocations))
	for i, location := range AnchorageLocations {
		locations[i] = string(location)
	}
	return locations
}
