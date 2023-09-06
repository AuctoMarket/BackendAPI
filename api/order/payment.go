package order

import (
	"BackendAPI/data"
	"BackendAPI/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

func CreatePaymentRequest(amount float64, orderId string, paymentType string, isGuest bool) (data.PaymentRequestResponseData, *utils.ErrorHandler) {
	var response data.PaymentRequestResponseData
	var webhookResource string
	currency := "SGD"
	hitpayBaseUrl, envHitpayExists := os.LookupEnv("HITPAY_BASE_URL")
	auctoBaseUrl, envAuctoExists := os.LookupEnv("AUCTO_BASE_URL")
	apiBaseUrl, envApiExists := os.LookupEnv("API_BASE_URL")
	hitpayApiKey, envApiKeyExists := os.LookupEnv("HITPAY_API_KEY")
	redirectResource := "/orders/" + orderId + "/payment-complete"
	if isGuest {
		webhookResource = "/api/v1/orders/" + orderId + "/payment-complete/guest"
	} else {
		webhookResource = "/api/v1/orders/" + orderId + "/payment-complete"
	}

	if !envHitpayExists || !envApiExists || !envAuctoExists || !envApiKeyExists {
		errResp := utils.InternalServerError(nil)
		return response, errResp
	}

	//Create Request Body
	var requestBody data.PaymentRequestData = data.PaymentRequestData{
		Amount:         amount,
		Currency:       currency,
		RedirectUrl:    auctoBaseUrl + redirectResource,
		Webhook:        apiBaseUrl + webhookResource,
		PaymentMethods: []string{paymentType}}

	requestBodyJSON, err := json.Marshal(requestBody)

	req, err := http.NewRequest(http.MethodPost, hitpayBaseUrl, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in creating POST request")
		return response, errResp
	}

	req.Header.Add("X-BUSINESS-API-KEY", hitpayApiKey)
	req.Header.Add("accept", "application/json")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		errResp := utils.InternalServerError(err)
		utils.LogError(err, "Error in sending POST request")
		return response, errResp
	}

	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&response)

	return response, nil
}
