package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenQR(t *testing.T) {
	genReq := GenQRCodeReqDTO{
		CommonReqDTO: CommonReqDTO{
			RequestId:   "",
			RequestTime: "",
			Signature:   "d1270a72bc42276e28f613d575211fc7b6628a252f8ba85e9acb43a748dcc9e4",
		},
		GenQRCodeReqDataDTO: GenQRCodeReqDataDTO{
			MerchantId:    "4313564511",
			TerminalId:    "462578343222",
			OrderId:       "06ca91fe9d-80dc-4819-2222-db21b2b01177",
			AccountNo:     "045704050000026",
			Description:   "this is a message description 230000",
			BillNumber:    "230000",
			TransactionId: "",
			Amount:        230000,
		},
	}
	resp, err := CallGenQR(genReq)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "00", resp.CommonRespDTO.ResponseCode, "")
}
