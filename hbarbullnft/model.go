package hbarbullnft

type NftMetadata struct {
	Name        string         `json:"name"`
	Description NftDescription `json:"description"`
	Creator     string         `json:"creator"`
	Category    string         `json:"category"`
	Supply      int64          `json:"supply"`
	Image       NftImage       `json:"image"`
}

type NftDescription struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type NftImage NftDescription

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
