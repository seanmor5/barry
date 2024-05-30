package mercury

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const apiBaseURL = "https://api.mercury.com/api/v1/"

type Config struct {
	APIKey string
}

type AccountResponse struct {
	Accounts []Account `json:"accounts"`
}

type Account struct {
	ID                string  `json:"id"`
	AccountNumber     string  `json:"accountNumber"`
	RoutingNumber     string  `json:"routingNumber"`
	Name              string  `json:"name"`
	Status            string  `json:"status"`
	Type              string  `json:"type"`
	CreatedAt         string  `json:"createdAt"`
	AvailableBalance  float64 `json:"availableBalance"`
	CurrentBalance    float64 `json:"currentBalance"`
	Kind              string  `json:"kind"`
	LegalBusinessName string  `json:"legalBusinessName"`
	DashboardLink     string  `json:"dashboardLink"`
}

type TransactionResponse struct {
	Total        int           `json:"total"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	Amount                     float64               `json:"amount"`
	BankDescription            *string               `json:"bankDescription"`
	CounterpartyId             string                `json:"counterpartyId"`
	CounterpartyName           string                `json:"counterpartyName"`
	CounterpartyNickname       *string               `json:"counterpartyNickname"`
	CreatedAt                  time.Time             `json:"createdAt"`
	DashboardLink              string                `json:"dashboardLink"`
	Details                    *TransactionDetails   `json:"details"`
	EstimatedDeliveryDate      time.Time             `json:"estimatedDeliveryDate"`
	FailedAt                   *time.Time            `json:"failedAt"`
	ID                         string                `json:"id"`
	Kind                       string                `json:"kind"`
	Note                       *string               `json:"note"`
	ExternalMemo               *string               `json:"externalMemo"`
	PostedAt                   *time.Time            `json:"postedAt"`
	ReasonForFailure           *string               `json:"reasonForFailure"`
	Status                     string                `json:"status"`
	FeeId                      *string               `json:"feeId"`
	CurrencyExchangeInfo       *CurrencyExchangeInfo `json:"currencyExchangeInfo"`
	CompliantWithReceiptPolicy *bool                 `json:"compliantWithReceiptPolicy"`
	HasGeneratedReceipt        *bool                 `json:"hasGeneratedReceipt"`
	CreditAccountPeriodId      *string               `json:"creditAccountPeriodId"`
	MercuryCategory            *string               `json:"mercuryCategory"`
	GeneralLedgerCodeName      *string               `json:"generalLedgerCodeName"`
	Attachments                []Attachment          `json:"attachments"`
}

type TransactionDetails struct {
	Address                      *Address                      `json:"address"`
	DomesticWireRoutingInfo      *DomesticWireRoutingInfo      `json:"domesticWireRoutingInfo"`
	ElectronicRoutingInfo        *ElectronicRoutingInfo        `json:"electronicRoutingInfo"`
	InternationalWireRoutingInfo *InternationalWireRoutingInfo `json:"internationalWireRoutingInfo"`
	DebitCardInfo                *CardInfo                     `json:"debitCardInfo"`
	CreditCardInfo               *CardInfo                     `json:"creditCardInfo"`
}

type Address struct {
	Address1   string  `json:"address1"`
	Address2   *string `json:"address2"`
	City       string  `json:"city"`
	State      *string `json:"state"`
	PostalCode string  `json:"postalCode"`
}

type DomesticWireRoutingInfo struct {
	BankName      *string  `json:"bankName"`
	AccountNumber string   `json:"accountNumber"`
	RoutingNumber string   `json:"routingNumber"`
	Address       *Address `json:"address"`
}

type ElectronicRoutingInfo struct {
	AccountNumber string  `json:"accountNumber"`
	RoutingNumber string  `json:"routingnumber"`
	BankName      *string `json:"bankName"`
}

type InternationalWireRoutingInfo struct {
	IBAN              string             `json:"iban"`
	SwiftCode         string             `json:"swiftCode"`
	CorrespondentInfo *CorrespondentInfo `json:"correspondentInfo"`
	BankDetails       *BankDetails       `json:"bankDetails"`
	Address           *Address           `json:"address"`
	PhoneNumber       *string            `json:"phoneNumber"`
	CountrySpecific   *CountrySpecific   `json:"countrySpecific"`
}

type CorrespondentInfo struct {
	RoutingNumber *string `json:"routingNumber"`
	SwiftCode     *string `json:"swiftCode"`
	BankName      *string `json:"bankName"`
}

type BankDetails struct {
	BankName  string `json:"bankName"`
	CityState string `json:"cityState"`
	Country   string `json:"country"`
}

type CountrySpecific struct {
	CountrySpecificDataCanada      *CountrySpecificDataCanada      `json:"countrySpecificDataCanada"`
	CountrySpecificDataAustralia   *CountrySpecificDataAustralia   `json:"countrySpecificDataAustralia"`
	CountrySpecificDataIndia       *CountrySpecificDataIndia       `json:"countrySpecificDataIndia"`
	CountrySpecificDataRussia      *CountrySpecificDataRussia      `json:"countrySpecificDataRussia"`
	CountrySpecificDataPhilippines *CountrySpecificDataPhilippines `json:"countrySpecificDataPhilippines"`
	CountrySpecificDataSouthAfrica *CountrySpecificDataSouthAfrica `json:"countrySpecificDataSouthAfrica"`
}

type CountrySpecificDataCanada struct {
	BankCode      string `json:"bankCode"`
	TransitNumber string `json:"transitNumber"`
}

type CountrySpecificDataAustralia struct {
	BSBCode string `json:"bsbCode"`
}

type CountrySpecificDataIndia struct {
	IFSCCode string `json:"ifscCode"`
}

type CountrySpecificDataRussia struct {
	INN string `json:"inn"`
}

type CountrySpecificDataPhilippines struct {
	RoutingNumber string `json:"routingNumber"`
}

type CountrySpecificDataSouthAfrica struct {
	BranchCode string `json:"branchCode"`
}

type CardInfo struct {
	ID string `json:"id"`
}

type CurrencyExchangeInfo struct {
	ConvertedFromCurrency string  `json:"convertedFromCurrency"`
	ConvertedToCurrency   string  `json:"convertedToCurrency"`
	ConvertedFromAmount   float64 `json:"convertedFromAmount"`
	ConvertedToAmount     float64 `json:"convertedToAmount"`
	FeeAmount             float64 `json:"feeAmount"`
	FeePercentage         float64 `json:"feePercentage"`
	ExchangeRate          float64 `json:"exchangeRate"`
	FeeTransactionId      string  `json:"feeTransactionId"`
}

type Attachment struct {
	FileName       string `json:"fileName"`
	URL            string `json:"url"`
	AttachmentType string `json:"attachmentType"`
}

type ListTransactionsParams struct {
	Limit  *int32
	Offset *int32
	Status *string
	Start  *string
	End    *string
	Search *string
}

func ListAccounts(config Config) (*AccountResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/accounts", apiBaseURL), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.APIKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch accounts: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var accountResponse AccountResponse
	if err := json.Unmarshal(body, &accountResponse); err != nil {
		return nil, err
	}

	return &accountResponse, nil
}

func ListTransactions(config Config, accountID string, params ListTransactionsParams) (*TransactionResponse, error) {
	apiUrl := fmt.Sprintf("%s/account/%s/transactions", apiBaseURL, accountID)

	query := url.Values{}
	if params.Limit != nil {
		query.Set("limit", fmt.Sprintf("%d", *params.Limit))
	}
	if params.Offset != nil {
		query.Set("offset", fmt.Sprintf("%d", *params.Offset))
	}
	if params.Status != nil {
		query.Set("status", *params.Status)
	}
	if params.Start != nil {
		query.Set("start", *params.Start)
	}
	if params.End != nil {
		query.Set("end", *params.End)
	}
	if params.Search != nil {
		query.Set("search", *params.Search)
	}

	if len(query) > 0 {
		apiUrl += "?" + query.Encode()
	}

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.APIKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch transactions: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var transactionResponse TransactionResponse
	if err := json.Unmarshal(body, &transactionResponse); err != nil {
		return nil, err
	}

	return &transactionResponse, nil
}
