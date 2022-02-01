package nftlib

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func GetPhoto(photoPath string) (NftPhoto, error) {
	fh, err := os.Open(photoPath)
	if err != nil {
		return NftPhoto{}, err
	}
	data, err := ioutil.ReadAll(fh)
	if err != nil {
		return NftPhoto{}, err
	}
	mimeType := http.DetectContentType(data)
	var output string
	switch mimeType {
	case "image/jpeg":
		output = "data:image/jpeg;base64,"
	case "image/png":
		output = "data:image/png;base64,"
	default:
		output = fmt.Sprintf("data:%s;base64,", mimeType)
	}
	output = output + base64.StdEncoding.EncodeToString(data)
	return NftPhoto{Photo: output}, nil
}

func NewMetadata(name string, description string, creator string, category string, supply int64, image string) NftMetadata {
	return NftMetadata{
		Name:        name,
		Description: description,
		Creator:     creator,
		Category:    category,
		Supply:      supply,
		Image:       image,
	}
}

func NewOneCenterMetadata(name string, description string, creator string,
	category string, supply int64, image string, addons []string, royalties_numerator int64,
	royalties_denominator int64, royalties_fallback int64) OneCenterMetadata {
	return OneCenterMetadata{
		Name:        name,
		Description: description,
		Creator:     creator,
		Category:    category,
		Supply:      supply,
		Properties: OneCenterNftProperties{
			Type: "object",
			Description: OneCenterDescription{
				AddOns: addons,
			},
		},
		Royalties: Royalties{
			Numerator:   royalties_numerator,
			Denominator: royalties_denominator,
			FallBackFee: royalties_fallback,
		},
		Image: image,
	}
}

func Ping() {
	fmt.Println("Pong")
}

func UploadMetadata(imageMetaData NftMetadata, nftStorageKey string) (string, error) {
	return Upload(imageMetaData, nftStorageKey)
}

func UploadImage(photoPath string, nftStorageKey string) (string, error) {
	nftPhoto, err := GetPhoto(photoPath)
	if err != nil {
		return "", err
	}
	return Upload(nftPhoto, nftStorageKey)
}

func UploadNft(photoPath string, name string, description string,
	creator string, category string, supply int64, nftStorageKey string) (string, error) {
	nftPhoto, err := GetPhoto(photoPath)
	if err != nil {
		return "", err
	}
	cid1, err := Upload(nftPhoto, nftStorageKey)
	if err != nil {
		return "", err
	}
	imageUrl := fmt.Sprintf("https://cloudflare-ipfs.com/ipfs/%s", cid1)
	imageMetaData := NewMetadata(name, description, creator, category, supply, imageUrl)
	cid2, err := Upload(imageMetaData, nftStorageKey)
	if err != nil {
		return "", err
	}
	fmt.Println("Created the following two uploads: " + cid1 + " " + cid2)
	return cid2, nil
}

func UploadRaw(data []byte, nftStorageKey string) (string, error) {
	client := &http.Client{
		Timeout: time.Second * 600,
	}
	req, err := http.NewRequest("POST", "https://api.nft.storage/upload", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+nftStorageKey)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		var response NftStorageOkResponse
		json.Unmarshal(body, &response)
		return response.Value.Cid, nil

	} else {
		return "", errors.New(string(body))
	}

}

func UploadPhotoRaw(photoPath string, nftStorageKey string) (string, error) {
	fh, err := os.Open(photoPath)
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(fh)
	if err != nil {
		return "", err
	}
	return UploadRaw(data, nftStorageKey)

}

func Upload(data interface{}, nftStorageKey string) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return UploadRaw(jsonBytes, nftStorageKey)
}
