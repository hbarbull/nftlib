package hbarbullnft

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
		return NftPhoto{}, errors.New("Not a supported mimetype: " + mimeType)
	}
	output = output + base64.StdEncoding.EncodeToString(data)
	return NftPhoto{Photo: output}, nil
}

func NewMetadata(name string, description string, creator string, category string, supply int64, image string) NftMetadata {
	return NftMetadata{
		Name:        name,
		Description: NftDescription{Type: "string", Description: description},
		Creator:     creator,
		Category:    category,
		Supply:      supply,
		Image:       NftImage{Type: "string", Description: image},
	}
}

/*
func UploadImage(photoPath string, nftStorageKey string) (string, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	nftPhoto, err := GetPhoto(photoPath)
	if err != nil {
		return "", err
	}
	jsonBytes, err := json.Marshal(nftPhoto)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", "https://api.nft.storage/upload", bytes.NewBuffer(jsonBytes))
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

func UploadMetadata(imageMetaData NftMetadata, nftStorageKey string) (string, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	jsonBytes, err := json.Marshal(imageMetaData)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", "https://api.nft.storage/upload", bytes.NewBuffer(jsonBytes))
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
*/

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

func Upload(data interface{}, nftStorageKey string) (string, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", "https://api.nft.storage/upload", bytes.NewBuffer(jsonBytes))
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
