package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type SubjectServer struct {
	St      *SubjectTable
	Session *SessionTable
}

type SubjectServerData struct {
	Session         string `json:"session"`
	VerifactionData string `json:"veri_data"`
}

func (ss *SubjectServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") == "application/json" {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var subjData SubjectServerData
		err = json.Unmarshal(data, &subjData)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Printf("%v", subjData)
		if r.Method == "POST" {

			//the verifacte this /
			ok, err := ss.VerifactionDataFromApple(subjData.VerifactionData, subjData.Session)
			if err != nil {
				log.Printf("%v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if !ok {
				log.Printf("%s", "ExpiresDate is old .")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		} else if r.Method == "GET" {
			date, err := ss.QueryExpiresDate(subjData.Session)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.Write([]byte(date))
			return
		}

	}
	w.WriteHeader(http.StatusBadRequest)

}

type VerifacteData struct {
	ExcludeOldData bool   `json:"exclude-old-transactions"`
	Password       string `json:"password"`
	VeriData       string `json:"receipt-data"`
}

type VerifacteDataResp struct {
	Infos       []Info        `json:"latest_receipt_info"`
	RenewalInfo []RenewalInfo `json:"pending_renewal_info"`
}

type Info struct {
	ProductId             string `json:"product_id"` //EasyPaperTracker_Subscription
	TransactionId         string `json:"transaction_id"`
	OriginalTransactionId string `json:"original_transaction_id"`
	PurchaseDate          string `json:"purchase_date"`
	OriginalPurchaseDate  string `json:"original_purchase_date"`
	ExpiresDate           string `json:"expires_date"`
	InAppOwnershipType    string `json:"in_app_ownership_type"`
}
type RenewalInfo struct {
	AutoRenewalStatus string `json:"auto_renew_status"`
}

func (ss *SubjectServer) VerifactionDataFromApple(veriData string, session string) (bool, error) {

	arg := VerifacteData{
		ExcludeOldData: true,
		Password:       "91716a5c1381490384d07db9b5d80a12",
		VeriData:       veriData,
	}
	data, err := json.Marshal(arg)
	if err != nil {
		return false, err
	}
	buffer := bytes.NewBuffer(data)
	request, err := http.NewRequest("POST", "https://sandbox.itunes.apple.com/verifyReceipt", buffer)
	if err != nil {
		return false, err
	}
	log.Printf("GO Veri..")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return false, err
	}
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	var vresp VerifacteDataResp
	err = json.Unmarshal(respData, &vresp)
	if err != nil {
		return false, err
	}
	var subjectionInfo Info

	for _, info := range vresp.Infos {
		if info.ProductId == "EasyPaperTracker_Subscription" {
			log.Printf(" ExpiresDate %s", info.ExpiresDate)
			log.Printf(" ExpiresDate Fixed %s", fixedTime(info.ExpiresDate))
			log.Printf("CurrentDate %s", time.Now().UTC().String())
			if info.ExpiresDate < time.Now().UTC().String() {
				return false, nil
			}
			subjectionInfo = info
			break
		}
	}
	//query uid
	row := ss.Session.db.QueryRow("select uid from sessionTable where session=? ", session)
	var uid int
	err = row.Scan(&uid)
	if err != nil {
		return false, err
	}
	exist, err := ss.St.TransactionIdisExist(subjectionInfo.TransactionId)
	if err != nil {
		return false, err
	}
	if exist {
		return true, nil
	}
	err = ss.St.InsertLatestReceiptInfo(uid, subjectionInfo.ProductId,
		subjectionInfo.TransactionId, subjectionInfo.OriginalTransactionId,
		fixedTime(subjectionInfo.PurchaseDate), fixedTime(subjectionInfo.OriginalPurchaseDate),
		fixedTime(subjectionInfo.ExpiresDate), subjectionInfo.InAppOwnershipType, "TEST", vresp.RenewalInfo[len(vresp.RenewalInfo)-1].AutoRenewalStatus)
	if err != nil {
		return false, err
	}
	return true, nil

}

func fixedTime(GMT string) string {
	strs := strings.Split(GMT, " ")
	var ret string
	for _, sub := range strs[:len(strs)-1] {
		ret += sub + " "
	}
	return ret
}

func (ss *SubjectServer) QueryExpiresDate(session string) (string, error) {

	row := ss.Session.db.QueryRow("select uid from sessionTable where session=? ", session)
	var uid int
	err := row.Scan(&uid)
	if err != nil {
		return "", err
	}
	qrow := ss.St.db.QueryRow("select ifnull (max(expires_date),\"NO SUBJECT\") from subjectTable where uid= ?", uid)
	var date string
	err = qrow.Scan(&date)
	if err != nil {
		return "", err
	}
	return date, nil
}
