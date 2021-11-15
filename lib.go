package nftlib

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

func CreateAndMintSeries(client *hedera.Client, nftStorageKey string, photoPath string,
	name string, description string, creator string, category string, supply int64) ([]string, error) {
	var nfts []string
	cid, err := UploadNft(photoPath, name, description, creator, category, supply, nftStorageKey)
	if err != nil {
		return nfts, err
	}
	tokenSymbol := fmt.Sprintf("IPFS://%s", cid)
	nftMetadata := fmt.Sprintf("https://cloudflare-ipfs.com/ipfs/%s", cid)
	tokenId, err := CreateToken(client, name, tokenSymbol, supply)
	if err != nil {
		return nfts, err
	}
	var i int64
	for i = 1; i <= supply; i++ {
		nftId, err := MintToken(client, tokenId, nftMetadata)
		nfts = append(nfts, nftId)
		if err != nil {
			return nfts, err
		}
	}
	return nfts, nil
}
