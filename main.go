package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

func HomepageHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Tech Company listing API with Golang"})
}

func main() {
	CallGenQR(GenQRCodeReqDTO{})
	router := gin.Default()
	router.GET("/", HomepageHandler)
	router.GET("/companies", GetCompaniesHandler)
	router.POST("/company", NewCompanyHandler)
	router.POST("/gen-transaction", NewGenTransactionHandler)
	router.PUT("/company/:id", UpdateCompanyHandler)
	router.DELETE("/company/:id", DeleteCompanyHandler)
	router.Run()
}

//
type GenQRCodeReqDTO struct {
	CommonReqDTO
	GenQRCodeReqDataDTO GenQRCodeReqDataDTO `json:"data"`
}
type GenQRCodeReqDataDTO struct {
	MerchantId    string `json:"merchantId"`
	TerminalId    string `json:"terminalId"`
	OrderId       string `json:"orderId"`
	AccountNo     string `json:"accountNo"`
	Description   string `json:"description"`
	BillNumber    string `json:"billNumber"`
	TransactionId string `json:"transactionId"`
	Amount        int64  `json:"amount"`
}

type CommonReqDTO struct {
	RequestId   string `json:"requestId"`
	RequestTime string `json:"requestTime"`
	Signature   string `json:"signature"`
}

type CommonRespDTO struct {
	RequestId       string `json:"requestId"`
	ResponseCode    string `json:"responseCode"`
	ResponseTime    string `json:"responseTime"`
	ResponseMessage string `json:"responseMessage"`
}

func NewUpRouter() *gin.Engine {
	router := gin.Default()

	return router
}

// add header
func buildEsbHeader() map[string]string {
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["X-IBM-CLIENT-SECRET"] = "WWCOiqcuSo8i0m5Gh7rLXaztbzqxxcpY7iqphLmG"

	return header
}

type GenQRCodeResqDTO struct {
	CommonRespDTO
	GenQRCodeResqDataDTO GenQRCodeResqDataDTO `json:"data"`
}
type GenQRCodeResqDataDTO struct {
	ImageStr      string `json:"imageStr"`
	TransactionId string `json:"transactionId"`
}

func BuildSignature(secretkey, plainText, algorithm string) string {
	var (
		signature string
	)

	switch algorithm {
	case "SHA256":
		sum := sha256.Sum256([]byte(plainText))
		signature = fmt.Sprintf("%x", sum)
	case "HMAC_SHA256":
		mac := hmac.New(sha256.New, []byte(secretkey))
		mac.Write([]byte(plainText))
		expectedMAC := mac.Sum(nil)
		signature = hex.EncodeToString(expectedMAC)
	default:
		mac := hmac.New(sha256.New, []byte(secretkey))
		mac.Write([]byte(plainText))
		expectedMAC := mac.Sum(nil)
		signature = hex.EncodeToString(expectedMAC)
	}
	return signature
}
func plainTextGen(reqDTO GenQRCodeReqDTO) string {
	plainText := reqDTO.GenQRCodeReqDataDTO.MerchantId +
		reqDTO.GenQRCodeReqDataDTO.TerminalId +
		reqDTO.GenQRCodeReqDataDTO.OrderId +
		reqDTO.GenQRCodeReqDataDTO.AccountNo +
		strconv.Itoa(int(reqDTO.GenQRCodeReqDataDTO.Amount)) +
		reqDTO.GenQRCodeReqDataDTO.BillNumber +
		reqDTO.GenQRCodeReqDataDTO.Description
	return plainText
}
func CallGenQR(req GenQRCodeReqDTO) (GenQRCodeResqDTO, error) {
	var (
		resp = GenQRCodeResqDTO{}
	)
	plainText := plainTextGen(req)
	signature := BuildSignature("qr-code-hdbank", plainText, os.Getenv("ALGORITHM_GEN"))
	req.CommonReqDTO.Signature = signature

	clientHTTP := NewClientHttp()
	reqC := ClientHttpRequest{
		Body:   req,
		Method: "POST",
		Url:    "https://7ods3u4e25.execute-api.ap-southeast-1.amazonaws.com/dev/gen-transaction",
		Header: buildEsbHeader(),
	}
	ctx := context.Background()
	httpResp, err := clientHTTP.Post(ctx, reqC)
	if err != nil {
		return resp, err
	}
	bodyByte, err := io.ReadAll(httpResp.Body)
	// bodyByte, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return resp, err
	}
	err = json.Unmarshal(bodyByte, &resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
	// r := NewUpRouter()

	// jsonValue, err := json.Marshal(genReq)
	// if err != nil {
	// 	panic(err)
	// }
	// req, err := http.NewRequest("POST", "https://7ods3u4e25.execute-api.ap-southeast-1.amazonaws.com/dev/gen-transaction", bytes.NewBuffer(jsonValue))
	// if err != nil {
	// 	panic(err)
	// }
	// w := httptest.NewRecorder()
	// r.ServeHTTP(w, req)
	// fmt.Println(w)
}

func NewGenTransactionHandler(c *gin.Context) {
	var genReq GenQRCodeReqDTO
	err := c.ShouldBindJSON(&genReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// newCompany.ID = xid.New().String()
	// companies = append(companies, newCompany)
	// c.JSON(http.StatusCreated, newCompany)
}

// common
type EsbHeader struct {
	XClientID string `json:"xClientId"`
	XAPIKey   string `json:"xAPIKey"`
}
type Common struct {
	ServiceVersion   string `json:"serviceVersion"`
	MessageId        string `json:"messageId"`
	TransactionId    string `json:"transactionId"`
	MessageTimestamp string `json:"messageTimestamp"`
}
type Client struct {
	SourceAppID  string     `json:"sourceAppID"`
	TargetAppIDs []string   `json:"targetAppIDs"`
	UserDetail   UserDetail `json:"userDetail"`
}
type UserDetail struct {
	UserID       string `json:"userID"`
	UserPassword string `json:"userPassword"`
}
type Header struct {
	Common Common `json:"common"`
	Client Client `json:"client"`
}

// request
type VerifyAccountNoReq struct {
	InquireAcctInfoReq InquireAcctInfoReq `json:"inquireAcctInfoReq"`
}

type InquireAcctInfoReq struct {
	Header  Header  `json:"header"`
	BodyReq BodyReq `json:"bodyReq"`
}

type BodyReq struct {
	FunctionCode string `json:"functionCode"`
	AcctNo       string `json:"acctNo"`
}

// response
type VerifyAccountNoResq struct {
	InquireAcctInfoRes InquireAcctInfoRes `json:"inquireAcctInfoRes"`
}

type InquireAcctInfoRes struct {
	Header         Header         `json:"header"`
	BodyRes        BodyResq       `json:"bodyRes"`
	ResponseStatus ResponseStatus `json:"responseStatus"`
}

type ResponseStatus struct {
	Status                 string      `json:"status"`
	GlobalErrorCode        string      `json:"globalErrorCode"`
	GlobalErrorDescription string      `json:"globalErrorDescription"`
	ErrorInfo              []ErrorInfo `json:"errorInfo"`
}
type ErrorInfo struct {
	SourceAppID string `json:"sourceAppID"`
	ErrorCode   string `json:"errorCode"`
	ErrorDesc   string `json:"errorDesc"`
}
type BodyResq struct {
	Acct         AcctResp `json:"acct"`
	FunctionCode string   `json:"functionCode"`
}

type AcctResp struct {
	AcctNo         string `json:"acctNo"`
	AcctCcy        string `json:"acctCcy"`
	AcctDesc       string `json:"acctDesc"`
	AcctOpenDate   string `json:"acctOpenDate"`
	OpeningBal     string `json:"openingBal"`
	CalcBal        string `json:"calcBal"`
	ActualBal      string `json:"actualBal"`
	LedgerBal      string `json:"ledgerBal"`
	OutstandingAmt string `json:"outstandingAmt"`
	LoanAvailAmt   string `json:"loanAvailAmt"`
	HoldAmt        string `json:"holdAmt"`
	AvailableAmt   string `json:"availableAmt"`
	IntRate        string `json:"intRate"`
	AcctStatus     string `json:"acctStatus"`
	AcctBranch     string `json:"acctBranch"`
	AcctBranchName string `json:"acctBranchName"`
	ProductCode    string `json:"productCode"`
	ProductName    string `json:"productName"`
	GlobalId       string `json:"globalId"`
	ClientShort    string `json:"clientShort"`
}

type Company struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	CEO     string `json:"ceo"`
	Revenue string `json:"revenue"`
}

var companies = []Company{
	{ID: "1", Name: "Dell", CEO: "Michael Dell", Revenue: "92.2 billion"},
	{ID: "2", Name: "Netflix", CEO: "Reed Hastings", Revenue: "20.2 billion"},
	{ID: "3", Name: "Microsoft", CEO: "Satya Nadella", Revenue: "320 million"},
}

func NewCompanyHandler(c *gin.Context) {
	var newCompany Company
	if err := c.ShouldBindJSON(&newCompany); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	newCompany.ID = xid.New().String()
	companies = append(companies, newCompany)
	c.JSON(http.StatusCreated, newCompany)
}
func GetCompaniesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, companies)
}
func UpdateCompanyHandler(c *gin.Context) {
	id := c.Param("id")
	var company Company
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	index := -1
	for i := 0; i < len(companies); i++ {
		if companies[i].ID == id {
			index = 1
		}
	}
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Company not found",
		})
		return
	}
	companies[index] = company
	c.JSON(http.StatusOK, company)
}
func DeleteCompanyHandler(c *gin.Context) {
	id := c.Param("id")
	index := -1
	for i := 0; i < len(companies); i++ {
		if companies[i].ID == id {
			index = 1
		}
	}
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Company not found",
		})
		return
	}
	companies = append(companies[:index], companies[index+1:]...)
	c.JSON(http.StatusOK, gin.H{
		"message": "Company has been deleted",
	})
}
