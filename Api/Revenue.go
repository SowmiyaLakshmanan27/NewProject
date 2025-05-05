package api

import (
	dbconnection "NewProject/db"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type RevenueReqStruct struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type RevenueDetails struct {
	ProductWiseRevenueArr []ProductWiseRevenue  `json:"ProductWiseRevenueArr"`
	CategoryWiseRevenue   []CategoryWiseRevenue `json:"categoryWiseRevenue"`
	RegionWiseRevenue     []RegionWiseRevenue   `json:"regionWiseRevenue"`
	TotalRevenue          float64               `json:"totalrevenue"`
}

type RevenueResponseStruct struct {
	RevenueDetailsRec RevenueDetails `json:"revenueDetailsRec"`
	Status            string         `json:"status"`
	Msg               string         `json:"msg"`
}

type ProductWiseRevenue struct {
	ProductName        string  `json:"productName"`
	ProductWiseRevenue float64 `json:"productWiseRevenue"`
}
type CategoryWiseRevenue struct {
	Category            string  `json:"category"`
	CategoryWiseRevenue float64 `json:"categoryWiseRevenue"`
}
type RegionWiseRevenue struct {
	Region            string  `json:"region"`
	RegionWiseRevenue float64 `json:"regionWiseRevenue"`
}

func GetRevenueDetails(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
	log.Println("GetRevenueDetails (+)..")

	var lRespRec RevenueResponseStruct

	if strings.EqualFold(r.Method, http.MethodPost) {

		log.Println("GetRevenueDetails +" + r.Method)
		var lRevenueReqRec RevenueReqStruct
		lRespRec.Status = "S"
		lBody, lErr := ioutil.ReadAll(r.Body)
		if lErr != nil {
			log.Println("GetRevenueDetails:001 " + lErr.Error())
			lRespRec.Status = "E"
			lRespRec.Msg = "GetRevenueDetails:001 " + lErr.Error()
		} else {
			lErr = json.Unmarshal(lBody, &lRevenueReqRec)
			if lErr != nil {
				log.Println("GetRevenueDetails:002 " + lErr.Error())
				lRespRec.Status = "E"
				lRespRec.Msg = "GetRevenueDetails:002 " + lErr.Error()
			} else {
				lRespRec.RevenueDetailsRec.TotalRevenue, lErr = FetchTotalRevenue(lRevenueReqRec)
				if lErr != nil {
					log.Println("GetRevenueDetails:003 " + lErr.Error())
					lRespRec.Status = "E"
					lRespRec.Msg = "GetRevenueDetails:003 " + lErr.Error()
				} else {
					lRespRec.RevenueDetailsRec.CategoryWiseRevenue, lErr = FetchCategoryWiseRevenue(lRevenueReqRec)
					if lErr != nil {
						log.Println("GetRevenueDetails:003 " + lErr.Error())
						lRespRec.Status = "E"
						lRespRec.Msg = "GetRevenueDetails:003 " + lErr.Error()
					} else {
						lRespRec.RevenueDetailsRec.RegionWiseRevenue, lErr = FetchRegionWiseRevenue(lRevenueReqRec)
						if lErr != nil {
							log.Println("GetRevenueDetails:003 " + lErr.Error())
							lRespRec.Status = "E"
							lRespRec.Msg = "GetRevenueDetails:003 " + lErr.Error()
						} else {
							lRespRec.RevenueDetailsRec.ProductWiseRevenueArr, lErr = FetchProductWiseRevenue(lRevenueReqRec)
							if lErr != nil {
								log.Println("GetRevenueDetails:003 " + lErr.Error())
								lRespRec.Status = "E"
								lRespRec.Msg = "GetRevenueDetails:003 " + lErr.Error()
							}
						}
					}
				}
			}

		}
	} else {
		lRespRec.Status = "E"
		lRespRec.Msg = "Invalid Methods"
	}
	lDatas, err := json.Marshal(lRespRec)
	if err != nil {
		fmt.Fprintf(w, "Error taking data"+err.Error())
	} else {
		fmt.Fprintf(w, string(lDatas))
	}
}

func FetchTotalRevenue(RevenueReqRec RevenueReqStruct) (float64, error) {
	log.Println("FetchTotalRevenue +")
	var lTotalRevenue float64
	//Establish DB connection
	lDb, lErr := dbconnection.LocalDbConnect()
	if lErr != nil {
		return lTotalRevenue, fmt.Errorf("FetchTotalRevenue:001 " + lErr.Error())
	} else {
		defer lDb.Close()
		// discount - 10%

		lCoreString :=
			`SELECT 
			SUM((od.Quantity * UnitPrice) - Discount + od.ShippingCost ) AS TotalRevenue
			FROM order_Details od
			WHERE Date BETWEEN ? AND ?;`
		lRows, lErr := lDb.Query(lCoreString, RevenueReqRec.StartDate, RevenueReqRec.EndDate)
		if lErr != nil {
			return lTotalRevenue, fmt.Errorf("FetchTotalRevenue:002 " + lErr.Error())
		} else {
			for lRows.Next() {
				lErr := lRows.Scan(&lTotalRevenue)
				if lErr != nil {
					return lTotalRevenue, fmt.Errorf("FetchTotalRevenue:003 " + lErr.Error())
				}
			}
		}
	}
	return lTotalRevenue, nil
}

func FetchCategoryWiseRevenue(RevenueReqRec RevenueReqStruct) (lCategoryWiseRevenueArr []CategoryWiseRevenue, lErr error) {
	log.Println("FetchCategoryWiseRevenue +")
	var lCategoryWiseRevenue CategoryWiseRevenue
	//Establish DB connection
	lDb, lErr := dbconnection.LocalDbConnect()
	if lErr != nil {
		return lCategoryWiseRevenueArr, fmt.Errorf("FetchCategoryWiseRevenue:001 " + lErr.Error())
	} else {
		defer lDb.Close()
		// discount - 10%

		lCoreString :=
			`SELECT 
			Category,SUM((od.Quantity * UnitPrice) - Discount + od.ShippingCost ) AS TotalRevenue
			FROM order_Details od
			WHERE Date BETWEEN ? AND ?
			group by Category;`
		lRows, lErr := lDb.Query(lCoreString, RevenueReqRec.StartDate, RevenueReqRec.EndDate)
		if lErr != nil {
			return lCategoryWiseRevenueArr, fmt.Errorf("FetchCategoryWiseRevenue:002 " + lErr.Error())
		} else {
			for lRows.Next() {
				lErr := lRows.Scan(&lCategoryWiseRevenue.Category, &lCategoryWiseRevenue.CategoryWiseRevenue)
				if lErr != nil {
					return lCategoryWiseRevenueArr, fmt.Errorf("FetchCategoryWiseRevenue:003 " + lErr.Error())
				} else {
					lCategoryWiseRevenueArr = append(lCategoryWiseRevenueArr, lCategoryWiseRevenue)
				}
			}
		}
	}
	return lCategoryWiseRevenueArr, nil
}

func FetchRegionWiseRevenue(RevenueReqRec RevenueReqStruct) (RegionWiseRevenueArr []RegionWiseRevenue, lErr error) {
	log.Println("FetchRegionWiseRevenue +")
	var RegionWiseRevenue RegionWiseRevenue
	//Establish DB connection
	lDb, lErr := dbconnection.LocalDbConnect()
	if lErr != nil {
		return RegionWiseRevenueArr, fmt.Errorf("FetchRegionWiseRevenue:001 " + lErr.Error())
	} else {
		defer lDb.Close()
		// discount - 10%

		lCoreString :=
			`SELECT Region,
			SUM((od.Quantity * UnitPrice) - Discount + od.ShippingCost ) AS TotalRevenue
			FROM order_Details od 
			WHERE Date BETWEEN ? AND ?
			group by Region ;`
		lRows, lErr := lDb.Query(lCoreString, RevenueReqRec.StartDate, RevenueReqRec.EndDate)
		if lErr != nil {
			return RegionWiseRevenueArr, fmt.Errorf("FetchRegionWiseRevenue:002 " + lErr.Error())
		} else {
			for lRows.Next() {
				lErr := lRows.Scan(&RegionWiseRevenue.Region, &RegionWiseRevenue.RegionWiseRevenue)
				if lErr != nil {
					return RegionWiseRevenueArr, fmt.Errorf("FetchRegionWiseRevenue:003 " + lErr.Error())
				} else {
					RegionWiseRevenueArr = append(RegionWiseRevenueArr, RegionWiseRevenue)
				}
			}
		}
	}
	return RegionWiseRevenueArr, nil
}

func FetchProductWiseRevenue(RevenueReqRec RevenueReqStruct) (ProductWiseRevenueArr []ProductWiseRevenue, lErr error) {
	log.Println("FetchProductWiseRevenue +")
	var ProductWiseRevenue ProductWiseRevenue
	//Establish DB connection
	lDb, lErr := dbconnection.LocalDbConnect()
	if lErr != nil {
		return ProductWiseRevenueArr, fmt.Errorf("FetchProductWiseRevenue:001 " + lErr.Error())
	} else {
		defer lDb.Close()
		// discount - 10%

		lCoreString :=
			`SELECT ProductName,
			SUM((od.Quantity * UnitPrice) - Discount + od.ShippingCost ) AS TotalRevenue
			FROM order_Details od 
			WHERE Date BETWEEN ? AND ?
			group by ProductName ;`
		lRows, lErr := lDb.Query(lCoreString, RevenueReqRec.StartDate, RevenueReqRec.EndDate)
		if lErr != nil {
			return ProductWiseRevenueArr, fmt.Errorf("FetchProductWiseRevenue:002 " + lErr.Error())
		} else {
			for lRows.Next() {
				lErr := lRows.Scan(&ProductWiseRevenue.ProductName, &ProductWiseRevenue.ProductWiseRevenue)
				if lErr != nil {
					return ProductWiseRevenueArr, fmt.Errorf("FetchProductWiseRevenue:003 " + lErr.Error())
				} else {
					ProductWiseRevenueArr = append(ProductWiseRevenueArr, ProductWiseRevenue)
				}
			}
		}
	}
	return ProductWiseRevenueArr, nil
}
