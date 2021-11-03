package hbarbullnft

import (
	"fmt"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

func CreateAndMintSingle(client *hedera.Client, nftStorageKey string, photoPath string,
	name string, description string, creator string, category string) (string, error) {
	var supply int64 = 1
	cid, err := UploadNft(photoPath, name, description, creator, category, supply, nftStorageKey)
	if err != nil {
		return "", err
	}
	tokenSymbol := fmt.Sprintf("IPFS://%s", cid)
	nftMetadata := fmt.Sprintf("https://cloudflare-ipfs.com/ipfs/%s", cid)
	tokenId, err := CreateToken(client, name, tokenSymbol, supply)
	if err != nil {
		return "", err
	}
	nftId, err := MintToken(client, tokenId, nftMetadata)
	if err != nil {
		return "", err
	}
	return nftId, nil
}
