package api

import (
	dbconnection "NewProject/db"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

type GetOrderDetails struct {
	OrderId        string  `json:"orderId"`
	ProductId      string  `json:"productId"`
	CustomerId     string  `json:"customerId"`
	ProductName    string  `json:"productName"`
	Category       string  `json:"category"`
	Region         string  `json:"region"`
	Date           string  `json:"date"`
	Quantity       int     `json:"quantity"`
	UnitPrice      float64 `json:"unitPrice"`
	Discount       float64 `json:"discount"`
	ShippingCost   float64 `json:"shippingCost"`
	PaymentMethods string  `json:"paymentMethods"`
	Name           string  `json:"name"`
	EmailId        string  `json:"emailId"`
	Address        string  `json:"address"`
	User           string  `json:"user"`
}
type OrderDetailsResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

func CSVRefreshHandler(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	var lOrderRespRec OrderDetailsResponse

	if strings.EqualFold(r.Method, http.MethodPost) {
		lTwoDArrayString, lErr := ReadCSV(r, "FileName")
		if lErr != nil {
			lOrderRespRec.Status = "E"
			lOrderRespRec.Msg = "CSVRefreshHandler:001 " + lErr.Error()
		} else {
			lOrderDetails, lErr := TwoDConversion(lTwoDArrayString)
			if lErr != nil {
				lOrderRespRec.Status = "E"
				lOrderRespRec.Msg = "CSVRefreshHandler:002 " + lErr.Error()
			} else {
				lErr = InsertOrders(lOrderDetails)
				log.Println(lOrderDetails, "lOrderDetails")
				if lErr != nil {
					lOrderRespRec.Status = "E"
					lOrderRespRec.Msg = "CSVRefreshHandler:003 " + lErr.Error()
				} else {
					lOrderRespRec.Status = "S"
					lOrderRespRec.Msg = "Successfully Inserted"
				}
			}
		}

	} else {
		lOrderRespRec.Status = "E"
		lOrderRespRec.Msg = "Invalid Methods"
	}
	lDatas, lErr := json.Marshal(lOrderRespRec)
	if lErr != nil {
		fmt.Fprintf(w, "Error taking data"+lErr.Error())
	} else {
		fmt.Fprint(w, string(lDatas))
	}
}
func GetFileDetails(r *http.Request, formName string) (*strings.Reader, string, *multipart.FileHeader, error) {
	log.Println("GetFileDetails(+)")

	fileStr := ""
	var file *strings.Reader

	// Attempt to retrieve the file data and header using r.FormFile(formName)
	fileBody, header, lErr := r.FormFile(formName)

	if lErr != nil {
		// If an error occurs during retrieval, return an empty file, empty fileStr, the header, and the error
		return file, fileStr, header, fmt.Errorf("GetFileDetails:001" + lErr.Error())
	} else {
		// If the file data is successfully retrieved, read its content into fileStr
		datas, _ := io.ReadAll(fileBody)
		fileStr = string(datas)

		// Create a strings.Reader (file) from fileStr to facilitate further use of the file's content
		file = strings.NewReader(fileStr)

		log.Println("GetFileDetails(-)")
		// Log the end of the function
		return file, fileStr, header, nil
	}
}

// ReadCSV reads the contents of a CSV file from a strings.Reader and returns the data as a 2D slice of strings.

// Step 1: Initialize a 2D slice to store the CSV data
// Step 2: Create a CSV reader for the input string reader
// Step 3: Read the CSV file row by row
// Step 4: Check for the end of the file
// Step 5: Append each row to the 2D slice
// Step 6: Return the 2D slice containing the CSV data and no error
func ReadCSV(r *http.Request, pFile string) ([][]string, error) {
	log.Println("ReadCSV (+)")
	var lRecord [][]string
	lFile, _, _, lErr := GetFileDetails(r, pFile)
	if lErr != nil {
		return lRecord, fmt.Errorf("ReadCSV:001" + lErr.Error())
	} else {
		// Step 2: Create a CSV reader for the input string reader
		lRows := csv.NewReader(lFile)

		// Step 3: Read the CSV file row by row
		for {
			// Step 4: Read a row from the CSV
			lRecordRow, lErr := lRows.Read()

			// Step 4: Check for the end of the file
			if lErr == io.EOF {
				break // Exit the loop when we reach the end of the file
			} else {
				// Step 5: Append the read row to the 2D slice
				lRecord = append(lRecord, lRecordRow)
			}
		}
	}
	log.Println("ReadCSV (-)")
	// Step 6: Return the 2D slice containing the CSV data and no error
	return lRecord, nil
}
func TwoDConversion(data [][]string) (orders []GetOrderDetails, lErr error) {
	log.Println("TwoDConversion (+)")
	for _, row := range data {
		if len(row) < 15 {
			fmt.Println("Skipping row due to insufficient data:", row)
			continue
		}

		quantity, _ := strconv.Atoi(row[7])
		unitPrice, _ := strconv.ParseFloat(row[8], 64)
		discount, _ := strconv.ParseFloat(row[9], 64)
		shippingCost, _ := strconv.ParseFloat(row[10], 64)

		order := GetOrderDetails{
			OrderId:        row[0],
			ProductId:      row[1],
			CustomerId:     row[2],
			ProductName:    row[3],
			Category:       row[4],
			Region:         row[5],
			Date:           row[6],
			Quantity:       quantity,
			UnitPrice:      unitPrice,
			Discount:       discount,
			ShippingCost:   shippingCost,
			PaymentMethods: row[11],
			Name:           row[12],
			EmailId:        row[13],
			Address:        row[14],
			User:           "Sowmiya Lakshmanan",
		}

		orders = append(orders, order)
	}
	log.Println("TwoDConversion (-)")
	return orders, nil
}

func InsertOrders(lOrderDetails []GetOrderDetails) error {
	log.Println("InsertOrders (+)")
	lDb, lErr := dbconnection.LocalDbConnect()
	if lErr != nil {
		return fmt.Errorf("InsertOrders:001" + lErr.Error())
	} else {
		defer lDb.Close()

		for i := 0; i < len(lOrderDetails); i++ {
			lCoreString := `INSERT INTO order_Details
		(OrderId,ProductId, CustomerId, ProductName,Category, Region,Date,Quantity,UnitPrice,Discount,ShippingCost, PaymentMethods, Name,EmailId,Address, Created_By, Created_Date, Updated_By, Updated_date)
		VALUES(?, ?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,now(), ?,now(), ?, now());`
			_, lErr := lDb.Exec(lCoreString, lOrderDetails[i].OrderId, lOrderDetails[i].ProductId, lOrderDetails[i].CustomerId, lOrderDetails[i].ProductName, lOrderDetails[i].Category, lOrderDetails[i].Region, lOrderDetails[i].Date, lOrderDetails[i].Quantity, lOrderDetails[i].UnitPrice, lOrderDetails[i].Discount, lOrderDetails[i].ShippingCost, lOrderDetails[i].PaymentMethods, lOrderDetails[i].Name, lOrderDetails[i].EmailId, lOrderDetails[i].Address, lOrderDetails[i].User, lOrderDetails[i].User)
			if lErr != nil {
				return fmt.Errorf("InsertOrders:002" + lErr.Error())
			}
		}
	}
	log.Println("InsertOrders-")
	return nil
}
