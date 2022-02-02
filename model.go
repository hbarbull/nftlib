package nftlib

type NftMetadata struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Creator     string                 `json:"creator"`
	Category    string                 `json:"category"`
	Supply      int64                  `json:"supply"`
	Properties  map[string]interface{} `json:"properties"`
	Image       string                 `json:"image"`
}

type NftMetadataWRoyalties struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Creator     string                 `json:"creator"`
	Category    string                 `json:"category"`
	Supply      int64                  `json:"supply"`
	Properties  map[string]interface{} `json:"properties"`
	Royalties   Royalties              `json:"royalties"`
	Image       string                 `json:"image"`
}

type OneCenterMetadata struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Creator     string                 `json:"creator"`
	Category    string                 `json:"category"`
	Supply      int64                  `json:"supply"`
	Properties  map[string]interface{} `json:"properties"`
	Royalties   Royalties              `json:"royalties"`
	Image       string                 `json:"image"`
}

type OneCenterPropertyAddons struct {
	AddOns []string `json:"addons"`
}

type Royalties struct {
	Numerator   int64 `json:"numerator"`
	Denominator int64 `json:"denominator"`
	FallBackFee int64 `json:"fallbackFee"`
}

type NftPhoto struct {
	Photo string `json:"photo"`
}

type NftStorageOkResponse struct {
	Ok    string                `json:"ok"`
	Value NftStorageResponseCID `json:"value"`
}

type NftStorageResponseCID struct {
	Cid string `json:"cid"`
}
