package nftlib

import (
	"errors"
	"strconv"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

func GetTestNetClient(acctId string, acctKey string) (*hedera.Client, error) {
	accountId, err := hedera.AccountIDFromString(acctId)
	if err != nil {
		return nil, err
	}

	privateKey, err := hedera.PrivateKeyFromString(acctKey)
	if err != nil {
		return nil, err
	}

	client := hedera.ClientForTestnet()
	client.SetOperator(accountId, privateKey)
	return client, nil
}

func GetMainNetClient(acctId string, acctKey string) (*hedera.Client, error) {
	accountId, err := hedera.AccountIDFromString(acctId)
	if err != nil {
		return nil, err
	}
	privateKey, err := hedera.PrivateKeyFromString(acctKey)
	if err != nil {
		return nil, err
	}

	client := hedera.ClientForMainnet()
	client.SetOperator(accountId, privateKey)
	return client, nil
}

func GetAccountBallance(client *hedera.Client, account_id string) (string, error) {
	accountId, err := hedera.AccountIDFromString(account_id)
	if err != nil {
		return "", err
	}
	query := hedera.NewAccountBalanceQuery().SetAccountID(accountId)
	accountBalance, err := query.Execute(client)
	if err != nil {
		return "", err
	}
	return strconv.FormatFloat(accountBalance.Hbars.As(hedera.HbarUnits.Hbar), 'f', 4, 64), nil
}

func CreateTokenWithRoyalty(client *hedera.Client, tokenName string, tokenSymbol string, maxSupply int64, r_num int64, r_denom int64, r_fallback int64) (string, error) {
	royalities := []hedera.Fee{hedera.NewCustomRoyaltyFee().
		SetFeeCollectorAccountID(client.GetOperatorAccountID()).
		SetDenominator(r_denom).
		SetNumerator(r_num).
		SetFallbackFee(
			hedera.NewCustomFixedFee().
				SetFeeCollectorAccountID(client.GetOperatorAccountID()).
				SetAmount(r_fallback),
		),
	}
	tokenCreateTx, err := hedera.NewTokenCreateTransaction().
		SetTokenName(tokenName).
		SetTokenSymbol(tokenSymbol).
		SetTokenType(hedera.TokenTypeNonFungibleUnique).
		SetSupplyKey(client.GetOperatorPublicKey()).
		SetSupplyType(hedera.TokenSupplyTypeFinite).
		SetDecimals(0).
		SetInitialSupply(0).
		SetMaxSupply(maxSupply).
		SetAdminKey(client.GetOperatorPublicKey()).
		SetTreasuryAccountID(client.GetOperatorAccountID()).
		SetCustomFees(royalities).
		Execute(client)
	if err != nil {
		return "", err
	}
	rcpt, err := tokenCreateTx.GetReceipt(client)
	if err != nil {
		return "", err
	}
	newTokenId := rcpt.TokenID
	return newTokenId.String(), nil
}

func CreateToken(client *hedera.Client, tokenName string, tokenSymbol string, maxSupply int64) (string, error) {
	tokenCreateTx, err := hedera.NewTokenCreateTransaction().
		SetTokenName(tokenName).
		SetTokenSymbol(tokenSymbol).
		SetTokenType(hedera.TokenTypeNonFungibleUnique).
		SetSupplyKey(client.GetOperatorPublicKey()).
		SetSupplyType(hedera.TokenSupplyTypeFinite).
		SetDecimals(0).
		SetInitialSupply(0).
		SetMaxSupply(maxSupply).
		SetAdminKey(client.GetOperatorPublicKey()).
		SetTreasuryAccountID(client.GetOperatorAccountID()).
		Execute(client)
	if err != nil {
		return "", err
	}
	rcpt, err := tokenCreateTx.GetReceipt(client)
	if err != nil {
		return "", err
	}
	newTokenId := rcpt.TokenID
	return newTokenId.String(), nil
}

func MintToken(client *hedera.Client, tokenIdStr string, metadata string) (string, error) {
	tokenId, err := hedera.TokenIDFromString(tokenIdStr)
	if err != nil {
		return "", err
	}
	tokenMintTx, err := hedera.NewTokenMintTransaction().
		SetTokenID(tokenId).
		SetMetadata([]byte(metadata)).
		FreezeWith(client)
	if err != nil {
		return "", err
	}
	signOp, err := tokenMintTx.SignWithOperator(client)
	if err != nil {
		return "", err
	}
	signTx, err := signOp.Execute(client)
	if err != nil {
		return "", err
	}
	signRcpt, err := signTx.GetReceipt(client)
	if err != nil {
		return "", err
	}
	serialNumbers := signRcpt.SerialNumbers
	if len(serialNumbers) == 1 {
		serNum := serialNumbers[0]
		nft := hedera.NftID{
			TokenID:      tokenId,
			SerialNumber: serNum,
		}
		return nft.String(), nil
	} else {
		return "", errors.New("error no serial numbers were found")
	}
}

func TransferNft(client *hedera.Client, nftIdStr string, toAcctStr string) error {
	nftId, err := hedera.NftIDFromString(nftIdStr)
	if err != nil {
		return err
	}
	toAcct, err := hedera.AccountIDFromString(toAcctStr)
	if err != nil {
		return err
	}
	nftTransferTx, err := hedera.NewTransferTransaction().
		AddNftTransfer(nftId, client.GetOperatorAccountID(), toAcct).
		Execute(client)
	if err != nil {
		return err
	}
	nftTransferRcpt, err := nftTransferTx.GetReceipt(client)
	if err != nil {
		return err
	}
	if nftTransferRcpt.Status == hedera.StatusSuccess {
		return nil
	} else {
		return errors.New("tranfer error")
	}
}
