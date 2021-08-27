package main

import (
	"encoding/json"
	"fmt"
	stdlog "log"

	"net/http"

	"github.com/aih/billtitles"
	"github.com/rs/zerolog/log"
)

/*
NOTE: I tried to create a service with echo and then with gin.
In both cases, importing the billtitles package (and in particular, something related to opening the SQLite database)
led to segfaults. Serving with plain net/http works.
*/

const serverPort = ":3333"

type message struct {
	Message string `json:"message"`
}

//----------
// Handlers
//----------

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func getBill(w http.ResponseWriter, r *http.Request) {
	var db = billtitles.GetDb("")

	var bill billtitles.Bill
	billString := r.URL.Query().Get("bill")
	if billString == "" {
		noBill := "No bill specified"

		json.NewEncoder(w).Encode(&message{noBill})
		return
	}
	billnumber := billtitles.BillnumberRegexCompiled.ReplaceAllString(billString, "$1$2$3")
	db.Find(&bill, "Billnumberversion = ?", billnumber)
	fmt.Println("{}", bill)

	json.NewEncoder(w).Encode(bill)
}
func getRelatedBills(w http.ResponseWriter, r *http.Request) {
	var db = billtitles.GetDb("")

	billString := r.URL.Query().Get("bill")
	if billString == "" {
		noBill := "No bill specified"

		json.NewEncoder(w).Encode(&message{noBill})
		return
	}
	billnumber := billtitles.BillnumberRegexCompiled.ReplaceAllString(billString, "$1$2$3")
	bills, _, err := billtitles.GetBillsWithSameTitleDb(db, billnumber)
	mytitles := billtitles.GetTitlesByBillnumberDb(db, bills[0].Billnumber)
	log.Info().Msgf("Found %d titles related to sample bill: %+v", len(mytitles), mytitles[0].Title)
	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode(&message{err.Error()})
		return

	}

	json.NewEncoder(w).Encode(bills)
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/bills", getBill)
	http.HandleFunc("/related", getRelatedBills)
	fmt.Println("***********************")
	fmt.Println("Server started on port:", serverPort)
	fmt.Println("***********************")
	stdlog.Fatal(http.ListenAndServe(serverPort, nil))
}

func main() {
	handleRequests()
}
