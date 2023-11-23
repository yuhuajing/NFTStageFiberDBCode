package traderjoe

import (
	"crypto/ecdsa"
	"io/ioutil"
	"log"
	"main/app/types"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
)

func Signature(c *fiber.Ctx) error {
	payload := &types.SigValue{}
	// if err := handler.ParseBody(c, payload); err != nil {
	// 	return err
	// }
	if err := c.BodyParser(payload); err != nil {
		return c.Status(400).JSON(types.ErrorResponse{
			Error: err.Error(),
		})
	}
	scaddress := common.HexToAddress(payload.Scaddress)
	sender := common.HexToAddress(payload.Sender)
	cosigner := common.HexToAddress("0x6593B47be3F4Bd1154c2faFb8Ad4aC4EFddD618f") //(os.Getenv("Signer_Addr"))
	qty := payload.Qty
	chainId := payload.ChainId
	timestamp := uint64(time.Now().Unix())
	requestId := utils.UUIDv4()
	packed := Abiencode(scaddress, sender, cosigner, qty, chainId, requestId, timestamp)
	res := hexutil.Encode(packed)[:2] + hexutil.Encode(packed)[66:]
	///auth/login

	//fmt.Println(fmt.Sprintf("scaddress:%s sender=%s, cosigner=%s, qty=%d, chainId=%d,timestamp=%d,rid=%s", scaddress, sender, cosigner, qty, chainId, timestamp, requestId))
	return c.SendString(res)
	//return handler.SendResponse(c, res)
}

type MyData struct {
	FieldSc        common.Address `json:"field_sc"`
	FieldSender    common.Address `json:"field_sender"`
	FieldCosigner  common.Address `json:"field_cosigner"`
	FieldQty       uint32         `json:"field_qty"`
	FieldChainid   uint32         `json:"field_chainid"`
	FieldRequestid string         `json:"field_requestid"`
	FieldTimestamp uint64         `json:"field_timestamp"`
	FieldSig       []byte         `json:"field_sig"`
}

func Abiencode(wallet, sender, cosigner common.Address, qty, chainID uint32, requestId string, timestamp uint64) []byte {
	// keystorePath := os.Getenv("KeystorePath")
	// keystorePass := os.Getenv("KeystorePass")
	// privateKey := GetPrivatefromKeystore(keystorePath, keystorePass)
	privateKey, _ := crypto.HexToECDSA("797391c7bd2e156e52329ceb6471496798e0c125ef35c4c3393329bd2a64f3f5")
	sig, err := PersonalSign(wallet, sender, cosigner, qty, chainID, requestId, timestamp, privateKey)
	if err != nil {
		log.Fatal(err)
	}
	//	fmt.Println(sig)

	structThing, _ := abi.NewType("tuple", "struct thing", []abi.ArgumentMarshaling{
		{Name: "field_sc", Type: "address"},
		{Name: "field_sender", Type: "address"},
		{Name: "field_cosigner", Type: "address"},
		{Name: "field_qty", Type: "uint32"},
		{Name: "field_chainid", Type: "uint32"},
		{Name: "field_requestid", Type: "string"},
		{Name: "field_timestamp", Type: "uint64"},
		{Name: "field_sig", Type: "bytes"},
	})

	args := abi.Arguments{
		{Type: structThing, Name: "param_one"},
	}

	record := MyData{
		wallet,
		sender,
		cosigner,
		qty,
		chainID,
		requestId,
		timestamp,
		sig,
	}

	packed, err := args.Pack(&record)
	if err != nil {
		log.Fatalf("error in pack: %v", err)
	}
	return packed

	//fmt.Printf("signature: %s\nabi encoded: %s\n", sig, hexutil.Encode(packed)[:2]+hexutil.Encode(packed)[66:])
}

func GetPrivatefromKeystore(ksfile string, pass string) *ecdsa.PrivateKey {
	keyjson, err := ioutil.ReadFile(ksfile)
	key, err := keystore.DecryptKey(keyjson, pass)
	if err != nil {
		panic(err)
	}
	return key.PrivateKey
}

func PersonalSign(wallet, sender, cosigner common.Address, qty, chainID uint32, requestId string, timestamp uint64, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	mes := solsha3.SoliditySHA3(
		// types
		//not uint256
		[]string{"address", "address", "address", "uint32", "uint32", "string", "uint64"},

		// values
		[]interface{}{
			wallet,
			sender,
			cosigner,
			qty,
			chainID,
			requestId,
			timestamp,
		},
	)
	hasht := solsha3.SoliditySHA3WithPrefix(mes)
	signatureBytes, err := crypto.Sign(hasht, privateKey)

	if err != nil {
		return nil, err
	}
	signatureBytes[64] += 27
	return signatureBytes, nil
}
