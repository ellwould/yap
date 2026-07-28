/*
    YAP (Yet Another PBX)
    Copyright (C) 2025 Elliot Michael Keavney

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as published
    by the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see https://github.com/ellwould/yap/blob/main/LICENSE
*/

package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/ellwould/csvcell"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/sony/sonyflake"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

//----------------------------------------------------------------------------------------------------

// Constant for directory path that contains the files yap-start.html and yap-end.html
const dirHTML string = "/etc/yap/html-css"

// Constant for fileStartHTML file
const fileStartHTML string = "yap-start.html"

// Constant for fileEndHTML file
const fileEndHTML string = "yap-end.html"

//----------------------------------------------------------------------------------------------------

// Function for webpage wallpaper
func wallpaper(w http.ResponseWriter, wallpaperCSS string) {
	fmt.Fprintf(w, "<div class=\"wallpaper "+wallpaperCSS+"\"></div>")
	fmt.Fprintf(w, "<div class=\"wallpaper wallpaper2 "+wallpaperCSS+"\"></div>")
	fmt.Fprintf(w, "<div class=\"wallpaper wallpaper3 "+wallpaperCSS+"\"></div>")
}

// Function to retrive HTTP email header
func emailHeaderHTTP(r *http.Request) (email string) {
	email = r.Header.Get("X-Email")
	return email
}

// Function to convert a string to a float64
func stringToFloat64(valueString string) float64 {
	valueFloat64, err := strconv.ParseFloat(valueString, 64)
	if err != nil {
		log.Fatal(err)
	}
	return valueFloat64
}

// Function to generate a unique ID using Sonyflake
func genID() (uniqueID string) {
	// Sonyflake custom setting
	sonyFlakeSetting := sonyflake.Settings{

		// Start time is set to YAP epoch (GitHub initial commit)
		StartTime: time.Date(2025, 10, 20, 02, 37, 0, 0, time.UTC),
	}

	// Generate Sonyflake
	generatedSonyFlake := sonyflake.NewSonyflake(sonyFlakeSetting)
	if generatedSonyFlake == nil {
		panic("sonyflake not generated")
	}

	// Generate a new ID
	uniqueIdentifier, err := generatedSonyFlake.NextID()
	if err != nil {
		panic("Unique identifier not generated")
	}
	uniqueID = strconv.FormatUint(uniqueIdentifier, 10)
	return uniqueID
}

// Function to generate random passowrds
func genPassword(passwordLength int) string {
	rand.Seed(time.Now().UnixNano())
	character := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_$"
	randomPassword := make([]byte, passwordLength)
	for i := 0; i < passwordLength; i++ {
		randomPassword[i] = character[rand.Intn(len(character))]
	}
	return string(randomPassword)
}

// Function for error message
func errorBox(w http.ResponseWriter, errorType string, headerCSS string, buttonCSS string) {
	fmt.Fprintf(w, "<div class=\"div-error-box\">")
	fmt.Fprintf(w, "  <h1 class=\""+headerCSS+"\">")
	if errorType == "email_error" {
		fmt.Fprintf(w, "    User Account Not Found<br>")
		fmt.Fprintf(w, "    <a href=\"/oauth2/sign_out?rd=https://github.com/logout\" class=\"button-general button-header "+buttonCSS+"\">Logout</a>")
	} else if errorType == "account_type_error" {
		fmt.Fprintf(w, "    Account Type Forbidden<br>")
		fmt.Fprintf(w, "    <a href=\"/yap\" class=\"button-general button-header "+buttonCSS+"\">Main Menu</a>")
	} else if errorType == "url_error" {
		fmt.Fprintf(w, "    The URL is Empty in /etc/yap.env<br>")
		fmt.Fprintf(w, "    <a href=\"/yap\" class=\"button-general button-header "+buttonCSS+"\">Main Menu</a>")
	} else if errorType == "client_id_error" {
		fmt.Fprintf(w, "    The Client ID is Empty in /etc/yap.env<br>")
		fmt.Fprintf(w, "    <a href=\"/yap\" class=\"button-general button-header "+buttonCSS+"\">Main Menu</a>")
	} else if errorType == "client_secret_error" {
		fmt.Fprintf(w, "    The Client Secret is Empty in /etc/yap.env<br>")
		fmt.Fprintf(w, "    <a href=\"/yap\" class=\"button-general button-header "+buttonCSS+"\">Main Menu</a>")
	} else if errorType == "refresh_token_error" {
		fmt.Fprintf(w, "    The Refresh Token is Empty in /etc/yap.env<br>")
		fmt.Fprintf(w, "    <a href=\"/yap\" class=\"button-general button-header "+buttonCSS+"\">Main Menu</a>")
	} else if errorType == "currency_code_error" {
		fmt.Fprintf(w, "    The Currency Code is Empty in /etc/yap.env<br>")
		fmt.Fprintf(w, "    <a href=\"/yap\" class=\"button-general button-header "+buttonCSS+"\">Main Menu</a>")
	} else {
		fmt.Fprintf(w, "    Unknown Error<br>")
		fmt.Fprintf(w, "    <a href=\"/yap\" class=\"button-general button-header "+buttonCSS+"\">Main Menu</a>")
	}
	fmt.Fprintf(w, "</h1>")
	fmt.Fprintf(w, "</div>")
}

// Function for the header
func header(w http.ResponseWriter, pageTitle string, headerCSS string, extraButtonName string, extraButtonURL string) {
	fmt.Fprintf(w, "<div class=\"div-header\">")
	fmt.Fprintf(w, "  <h1 class=\""+headerCSS+"\">")
	fmt.Fprintf(w, "    ⊛ YAP [Yet Another PBX] ⊛")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "    <a href=\"/oauth2/sign_out?rd=https://github.com/logout\" class=\"button-general button-header button-logout\">Logout</a>")
	fmt.Fprintf(w, "    <a href=\"/yap\" class=\"button-general button-header button-home\">Home</a>")
	fmt.Fprintf(w, "    <a href=\"https://github.com/ellwould/yap/blob/main/LICENSE\" target=\"_blank\" class=\"button-general button-header button-license\">License</a>")
	if extraButtonName != "" {
		fmt.Fprintf(w, "    <a href=\""+extraButtonURL+"\" class=\"button-general button-header button-extra\">"+extraButtonName+"</a>")
	}
	fmt.Fprintf(w, "    <br>")
	fmt.Fprintf(w, "    "+pageTitle)
	fmt.Fprintf(w, "</h1>")
	fmt.Fprintf(w, "</div>")
}

// Function for the footer
func footer(w http.ResponseWriter, headerCSS string, buttonCSS string) {
	fmt.Fprintf(w, "<div class=\"div-footer\">")
	fmt.Fprintf(w, "  <h1 class=\""+headerCSS+"\">")
	fmt.Fprintf(w, "    <a href=\"https://github.com/ellwould/yap\" target=\"_blank\" class=\"button-general button-footer "+buttonCSS+"\">YAP Source Code</a>")
	fmt.Fprintf(w, "    <a href=\"https://github.com/ellwould\" target=\"_blank\" class=\"button-general button-footer "+buttonCSS+"\">Other Software</a>")
	fmt.Fprintf(w, "  </h1>")
	fmt.Fprintf(w, "</div>")
}

// Function to format ISO date format
func formatDate(date string) string {
	const (
		iso    = "2006-01-02"
		layout = "2/January/2006"
	)

	parse, _ := time.Parse(iso, date)
	return string(parse.Format(layout))
}

// Function to format ISO datetime format
func formatDateTime(dateTime string) string {
	const (
		iso    = "2006-01-02 15:04:05"
		layout = "2/January/2006 15:04:05 PM"
	)

	parse, _ := time.Parse(iso, dateTime)
	return string(parse.Format(layout))
}

// Get current date
func currentDate() string {
	date := time.Now().UTC()
	result := date.Format("2006-01-02")
	return string(result)
}

// Function to trim string
func trimString(unTrimmedString string, characterToTrim string) (result string) {
	result = strings.Replace(unTrimmedString, characterToTrim, "", -1)
	return result
}

//----------------------------------------------------------------------------------------------------

// Pure embedded HTML Go functions

// Function to create HTML date input
func dateHTML(w http.ResponseWriter, inputValue string, labelMessage string) {
	fmt.Fprintf(w, "  <label for=\""+inputValue+"\"><b>Enter/Select "+labelMessage+":</b>")
	fmt.Fprintf(w, "  </label><br>")
	fmt.Fprintf(w, "  <input type=\"date\" id=\""+inputValue+"\" name=\""+inputValue+"\">")
	fmt.Fprintf(w, "<br>")
}

// Function to create HTML text input
func inputHTML(w http.ResponseWriter, inputValue string, labelMessage string) {
	fmt.Fprintf(w, "  <label for=\""+inputValue+"\"><b>Enter "+labelMessage+":</b>")
	fmt.Fprintf(w, "  </label><br>")
	fmt.Fprintf(w, "  <input type=\"text\" id=\""+inputValue+"\" name=\""+inputValue+"\">")
	fmt.Fprintf(w, "<br>")
}

// Function to create read only HTML text input
func inputReadOnlyHTML(w http.ResponseWriter, inputValue string, labelMessage string, readOnlyData string) {
	fmt.Fprintf(w, "  <label for=\""+inputValue+"\"><b>"+labelMessage+":</b>")
	fmt.Fprintf(w, "  </label><br>")
	fmt.Fprintf(w, "  <input style=\"text-align: center;\" type=\"text\" id=\""+inputValue+"\" name=\""+inputValue+"\" value=\""+readOnlyData+"\" readonly>")
	fmt.Fprintf(w, "<br>")
}

// Function to put slice elements into HTML select options
func selectSingleHTML(w http.ResponseWriter, selectValue string, labelMessage string, optionValue []string) {
	fmt.Fprintf(w, "  <label for=\""+selectValue+"\"><b>Select "+labelMessage+":</b>")
	fmt.Fprintf(w, "  </label><br>")
	fmt.Fprintf(w, "  <select id=\""+selectValue+"\" name=\""+selectValue+"\">")
	for _, value := range optionValue {
		fmt.Fprintf(w, "<option value=\""+string(value)+"\">&nbsp "+string(value)+"</option>")
	}
	fmt.Fprintf(w, "  </select>")
}

// Function to put nested slice elements into HTML select options via a loop; shows the option value in the option tags and prepends with ID: and Name:
func selectDoubleHTML(w http.ResponseWriter, selectValue string, labelMessage string, optionValue [][]string) {
	fmt.Fprintf(w, "  <label for=\""+selectValue+"\"><b>Select "+labelMessage+" (Cannot Be Empty):</b>")
	fmt.Fprintf(w, "  </label><br>")
	fmt.Fprintf(w, "  <select id=\""+selectValue+"\" name=\""+selectValue+"\">")
	fmt.Fprintf(w, "<option value></option>")
	for _, value := range optionValue {
		fmt.Fprintf(w, "<option value=\""+string(value[0:][0])+"\">&nbsp "+labelMessage+" ID: "+string(value[0:][0])+" | "+labelMessage+" Name: "+string(value[1:][0])+"</option>")
	}
	fmt.Fprintf(w, "  </select>")
}

// Function to put nested slice elements into HTML select options via a loop; shows one element in the option tags, does not show the option value in the option tags
func selectDoubleHiddenHTML(w http.ResponseWriter, selectValue string, labelMessage string, optionValue [][]string) {
	fmt.Fprintf(w, "  <label for=\""+selectValue+"\"><b>Select "+labelMessage+":</b>")
	fmt.Fprintf(w, "  </label><br>")
	fmt.Fprintf(w, "  <select id=\""+selectValue+"\" name=\""+selectValue+"\">")
	fmt.Fprintf(w, "<option value></option>")
	for _, value := range optionValue {
		fmt.Fprintf(w, "<option value=\""+string(value[0:][0])+"\">&nbsp "+string(value[1:][0])+"</option>")
	}
	fmt.Fprintf(w, "  </select>")
}

// Function to create warning message
func messageHTML(w http.ResponseWriter, message string, messageType string) {
	if messageType == "warning" || messageType == "success" {
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "  <table>")
		fmt.Fprintf(w, "    <tr>")
		fmt.Fprintf(w, "      <th class=\"table-"+messageType+"-message\">"+message+"</th>")
		fmt.Fprintf(w, "    </tr>")
		fmt.Fprintf(w, "  </table>")
	} else {
		panic("Message type must be either warning or success")
	}
}

// Embedded JavaScript and associated HTML functions
type jsFunctionParameter struct {
	funcNameJS   string
	inputID      string
	tableID      string
	divID        string
	placeholder  string
	columnNumber int
	data         string
	buttonCSS    string
	fileName     string
	pathURL      string
}

// JavaScript toggle function
func toggleDivJS(w http.ResponseWriter, parameter jsFunctionParameter) {
	fmt.Fprintf(w, "<script>")
	fmt.Fprintf(w, "  function "+parameter.funcNameJS+"() {")
	fmt.Fprintf(w, "  var x = document.getElementById(\""+parameter.divID+"\");")
	fmt.Fprintf(w, "  if (x.style.display === \"none\") {")
	fmt.Fprintf(w, "    x.style.display = \"table\";")
	fmt.Fprintf(w, "  } else {")
	fmt.Fprintf(w, "    x.style.display = \"none\";")
	fmt.Fprintf(w, "  }")
	fmt.Fprintf(w, "}")
	fmt.Fprintf(w, "</script>")
}

// Input for filtering a HTML table
func inputTableHTML(w http.ResponseWriter, parameter jsFunctionParameter) {
	fmt.Fprintf(w, "<input type=\"text\" id=\""+parameter.inputID+"\" onkeyup=\""+parameter.funcNameJS+"()\" placeholder=\"Filter Via "+parameter.placeholder+"...\" title=\""+parameter.placeholder+"\">")
}

// JavaScript filter HTML table function
func filterTableJS(w http.ResponseWriter, parameter jsFunctionParameter) {
	if parameter.columnNumber > 11 {
		panic("Table column number cannot exceed 11")
	} else if parameter.columnNumber < 0 {
		panic("Table column number cannot be a negative number")
	} else {
		fmt.Fprintf(w, "<script>")
		fmt.Fprintf(w, "function "+parameter.funcNameJS+"() {")
		fmt.Fprintf(w, "  var input, filter, table, tr, td, i, txtValue;")
		fmt.Fprintf(w, "  input = document.getElementById(\""+parameter.inputID+"\");")
		fmt.Fprintf(w, "  filter = input.value.toUpperCase();")
		fmt.Fprintf(w, "  table = document.getElementById(\""+parameter.tableID+"\");")
		fmt.Fprintf(w, "  tr = table.getElementsByTagName(\"tr\");")
		fmt.Fprintf(w, "  for (i = 0; i < tr.length; i++) {")
		fmt.Fprintf(w, "    td = tr[i].getElementsByTagName(\"td\")["+strconv.Itoa(parameter.columnNumber)+"];")
		fmt.Fprintf(w, "    if (td) {")
		fmt.Fprintf(w, "      txtValue = td.textContent || td.innerText;")
		fmt.Fprintf(w, "      if (txtValue.toUpperCase().indexOf(filter) > -1) {")
		fmt.Fprintf(w, "        tr[i].style.display = \"\";")
		fmt.Fprintf(w, "      } else {")
		fmt.Fprintf(w, "        tr[i].style.display = \"none\";")
		fmt.Fprintf(w, "      }")
		fmt.Fprintf(w, "    }")
		fmt.Fprintf(w, "  }")
		fmt.Fprintf(w, "}")
		fmt.Fprintf(w, "</script>")
	}
}

// JavaScript to copy data to users clipboard
func copyButtonJS(w http.ResponseWriter, parameter jsFunctionParameter) {
	dataTrim := strings.Replace(parameter.data, "-", "", -1)
	fmt.Fprintf(w, "<div class=\"button-data-space\"></div><button class=\"button-general "+parameter.buttonCSS+"\" onclick=cp"+dataTrim+"()>&nbsp Copy &#10697 &nbsp</button><br>")
	fmt.Fprintf(w, "<script>")
	fmt.Fprintf(w, "  function cp"+dataTrim+"() {")
	fmt.Fprintf(w, "    navigator.clipboard.writeText('"+parameter.data+"');")
	fmt.Fprintf(w, "  }")
	fmt.Fprintf(w, "</script>")
}

// HTML button to call the JavaScript exportCSVJS function
func exportCSVButtonHTML(w http.ResponseWriter, parameter jsFunctionParameter) {
	fmt.Fprintf(w, "<button class=\"button-general "+parameter.buttonCSS+"\" onclick=\"exportTable"+parameter.funcNameJS+"ToCSV()\">Export to CSV</button><br>")
}

// JavaScript to download or view a HTML table as a CSV file
func exportCSVJS(w http.ResponseWriter, parameter jsFunctionParameter) {
	fmt.Fprintf(w, "<script>")
	fmt.Fprintf(w, "  function exportTable"+parameter.funcNameJS+"ToCSV() {")
	fmt.Fprintf(w, "    const table = document.getElementById('"+parameter.tableID+"');")
	fmt.Fprintf(w, "    let csvData = '';")
	fmt.Fprintf(w, "    for (let i = 0; i < table.rows.length; i++) {")
	fmt.Fprintf(w, "      let row = table.rows[i];")
	fmt.Fprintf(w, "      let rowData = [];")
	fmt.Fprintf(w, "      for (let j = 0; j < row.cells.length; j++) {")
	fmt.Fprintf(w, "        let cell = row.cells[j];")
	fmt.Fprintf(w, "        rowData.push(cell.textContent.replace(/,/g, ''));")
	fmt.Fprintf(w, "      }")
	fmt.Fprintf(w, "      csvData += rowData.join(',') + '\\r\\n';")
	fmt.Fprintf(w, "    }")
	fmt.Fprintf(w, "    const blob = new Blob([csvData], { type: 'text/csv' });")
	fmt.Fprintf(w, "    const url = URL.createObjectURL(blob);")
	fmt.Fprintf(w, "    const csv = document.createElement('a');")
	fmt.Fprintf(w, "    csv.href = url;")
	fmt.Fprintf(w, "    csv.download = '"+parameter.fileName+".csv';")
	fmt.Fprintf(w, "    document.body.append(csv);")
	fmt.Fprintf(w, "    csv.click();")
	fmt.Fprintf(w, "    document.body.remove(csv);")
	fmt.Fprintf(w, "    window.open('/"+parameter.pathURL+"', '_blank');")
	fmt.Fprintf(w, "  }")
	fmt.Fprintf(w, "</script>")
}

//----------------------------------------------------------------------------------------------------

// Database functions

// Database struct
type databaseFunctionParameter struct {
	connection          *sql.DB
	database            string
	table               string
	column              string
	columnWhere         string
	columnWhereValue    string
	columnWhereAnd      string
	columnWhereValueAnd string
	countMinusOne       bool
}

// General struct for list/add/edit/delete functions
type generalFunctionParameter struct {
	userID                  string
	userTypeID              string
	userCustomerID          string
	userPBXID               string
	defaultExtLimit         string
	currencySymbol          string
	yapAdminUKVATRegistered string
}

// InvoicePBXExt struct
type invoicePBXExtFunctionParameter struct {
	customerID        string
	pbxID             string
	serviceProduct    string
	tag               string
	sellPrice         string
	salesTaxRate      string
	salesTaxStatus    string
	billItemOnce      string
	itemOnHold        string
	contractLength    string
	contractStartDate string
}

// Struct for accounting software API
type accountingSoftwareParameter struct {
	url            string
	accessToken    string
	customerID     string
	clientID       string
	clientSecret   string
	refreshToken   string
	currencyCode   string
	httpStatusCode int
}

// Function to insert NULL value into database column
func nullSQL(value string) sql.NullString {
	if len(value) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: value,
		Valid:  true,
	}
}

// Function to select value from a table with the WHERE clause
func selectWhere(dbSelectWhere databaseFunctionParameter) string {
	var selectWhere string
	selectWhereQuery, err := dbSelectWhere.connection.Query(`SELECT
								   `+dbSelectWhere.column+`
								 FROM
								   `+dbSelectWhere.database+`.`+dbSelectWhere.table+`
								 WHERE
								   `+dbSelectWhere.columnWhere+` = ?;`, dbSelectWhere.columnWhereValue)

	if err != nil {
		panic(err)
	}
	for selectWhereQuery.Next() {
		err := selectWhereQuery.Scan(&selectWhere)
		if err != nil {
			panic(err)
		}
	}
	return selectWhere
}

// Function to select value from a table with the WHERE clause & the AND operator
func selectWhereAnd(dbSelectWhereAnd databaseFunctionParameter) string {
	var selectWhereAnd string
	selectWhereAndQuery, err := dbSelectWhereAnd.connection.Query(`SELECT
                                                    	                 `+dbSelectWhereAnd.column+`
                                                                       FROM
                                                                         `+dbSelectWhereAnd.database+`.`+dbSelectWhereAnd.table+`
                                                                       WHERE
                                                                         `+dbSelectWhereAnd.columnWhere+` = ?`+`
                                                                       AND
                                                                         `+dbSelectWhereAnd.columnWhereAnd+` = ?;`, dbSelectWhereAnd.columnWhereValue, dbSelectWhereAnd.columnWhereValueAnd)

	if err != nil {
		panic(err)
	}
	for selectWhereAndQuery.Next() {
		err := selectWhereAndQuery.Scan(&selectWhereAnd)
		if err != nil {
			panic(err)
		}
	}
	return selectWhereAnd
}

// Function to count rows in a table
func totalTableCount(dbTotalTableCount databaseFunctionParameter) string {
	if dbTotalTableCount.countMinusOne == true {
		var countMinusOne string
		countMinusOneQuery, err := dbTotalTableCount.connection.Query(`SELECT
									       COUNT(*) -1
									       FROM
									         ` + dbTotalTableCount.database + `.` + dbTotalTableCount.table + `;`)
		if err != nil {
			panic(err)
		}
		for countMinusOneQuery.Next() {
			err = countMinusOneQuery.Scan(&countMinusOne)
			// Error
			if err != nil {
				panic(err)
			}
		}
		return countMinusOne
	} else {
		var count string
		countQuery, err := dbTotalTableCount.connection.Query(`SELECT
								       COUNT(*)
								       FROM
								         ` + dbTotalTableCount.database + `.` + dbTotalTableCount.table + `;`)
		if err != nil {
			panic(err)
		}
		for countQuery.Next() {
			err = countQuery.Scan(&count)
			// Error
			if err != nil {
				panic(err)
			}
		}
		return count
	}
}

// Function to counts rows in a table with the WHERE clause
func totalTableCountWhere(dbTotalTableCountWhere databaseFunctionParameter) string {
	if dbTotalTableCountWhere.countMinusOne == true {
		var countMinusOne string
		countMinusOneQuery, err := dbTotalTableCountWhere.connection.Query(`SELECT
                                                                    COUNT(*) -1
                                                                    FROM
                                                                      `+dbTotalTableCountWhere.database+`.`+dbTotalTableCountWhere.table+`
                                                                    WHERE
                                                                      `+dbTotalTableCountWhere.columnWhere+` = ?;`, dbTotalTableCountWhere.columnWhereValue)
		//Error
		if err != nil {
			panic(err)
		}
		for countMinusOneQuery.Next() {
			err = countMinusOneQuery.Scan(&countMinusOne)
			// Error
			if err != nil {
				panic(err)
			}
		}
		return countMinusOne
	} else {
		var count string
		countQuery, err := dbTotalTableCountWhere.connection.Query(`SELECT
								    COUNT(*)
								    FROM
								      `+dbTotalTableCountWhere.database+`.`+dbTotalTableCountWhere.table+`
								    WHERE
								      `+dbTotalTableCountWhere.columnWhere+` = ?;`, dbTotalTableCountWhere.columnWhereValue)
		//Error
		if err != nil {
			panic(err)
		}
		for countQuery.Next() {
			err = countQuery.Scan(&count)
			// Error
			if err != nil {
				panic(err)
			}
		}
		return count
	}
}

// Function to counts rows in a table with the WHERE clause & the AND operator
func totalTableCountWhereAnd(dbTotalTableCountWhereAnd databaseFunctionParameter) string {
	var count string
	countQuery, err := dbTotalTableCountWhereAnd.connection.Query(`SELECT
								       COUNT(*)
								       FROM
								         `+dbTotalTableCountWhereAnd.database+`.`+dbTotalTableCountWhereAnd.table+`
								       WHERE
								         `+dbTotalTableCountWhereAnd.columnWhere+` = ?`+`
								       AND
								         `+dbTotalTableCountWhereAnd.columnWhereAnd+` = ?;`, dbTotalTableCountWhereAnd.columnWhereValue, dbTotalTableCountWhereAnd.columnWhereValueAnd)
	//Error
	if err != nil {
		panic(err)
	}
	for countQuery.Next() {
		err = countQuery.Scan(&count)
		// Error
		if err != nil {
			panic(err)
		}
	}
	return count
}

func userAccountData(dbUserAccountData databaseFunctionParameter, data string) string {

	var dbSelectWhere databaseFunctionParameter
	dbSelectWhere.connection = dbUserAccountData.connection
	dbSelectWhere.database = dbUserAccountData.database
	dbSelectWhere.table = "view___account_detail"
	if data == "id" {
		dbSelectWhere.column = "user_account_id"
	} else if data == "type_id" {
		dbSelectWhere.column = "user_account_type_id"
	} else if data == "customer_id" {
		dbSelectWhere.column = "customer_id"
	} else if data == "customer_name" {
		dbSelectWhere.column = "customer_name"
	} else if data == "pbx_id" {
		dbSelectWhere.column = "pbx_id"
	} else if data == "pbx_name" {
		dbSelectWhere.column = "pbx_name"
	} else {
		panic("The function userAccountData can only accept the following arguments: id, type_id, customer_id, customer_name, pbx_id or pbx_name")
	}
	dbSelectWhere.columnWhere = "user_account_email"
	dbSelectWhere.columnWhereValue = dbUserAccountData.columnWhereValue

	return selectWhere(dbSelectWhere)
}

// Function to retrive account type name(s) and ID(s) or just the account type ID(s) from the user_account_type table
func userAccountTypeSlice(dbDetail databaseFunctionParameter) ([][]string, []string) {
	// Get account type ID and name from the database and append to slice
	var userAccountTypeIDNameList [][]string
	var userAccountTypeIDList []string

	var userAccountTypeID string
	var userAccountTypeName string

	userAccountTypeIDNameSQL, err := dbDetail.connection.Query(`SELECT
                                                                      id,
                                                                      type
                                                                    FROM
                                                                      yap.user_account_type;`)

	// Error
	if err != nil {
		panic(err)
	}

	for userAccountTypeIDNameSQL.Next() {

		err = userAccountTypeIDNameSQL.Scan(
			&userAccountTypeID,
			&userAccountTypeName,
		)

		// Error
		if err != nil {
			panic(err)
		}

		var userAccountTypeIDAndName []string
		userAccountTypeIDAndName = append([]string{userAccountTypeID}, []string{userAccountTypeName}...)
		userAccountTypeIDNameList = append(userAccountTypeIDNameList, userAccountTypeIDAndName)
		userAccountTypeIDList = append(userAccountTypeIDList, userAccountTypeID)
	}
	return userAccountTypeIDNameList, userAccountTypeIDList
}

// Function to retrive customer name(s) and ID(s) or just the customer ID(s) from the view___customer_detail view
func customerSlice(dbDetail databaseFunctionParameter) ([][]string, []string) {
	// Get customer ID and name from the database and append to slice
	var customerIDNameList [][]string
	var customerIDList []string

	var customerID string
	var customerName string

	customerIDNameSQL, err := dbDetail.connection.Query(`SELECT
                                                               customer_id,
                                                               customer_name
                                                             FROM
                                                               yap.view___customer_detail;`)

	// Error
	if err != nil {
		panic(err)
	}

	for customerIDNameSQL.Next() {

		err = customerIDNameSQL.Scan(
			&customerID,
			&customerName,
		)

		// Error
		if err != nil {
			panic(err)
		}

		var customerIDAndName []string
		if customerID != "1" {
			customerIDAndName = append([]string{customerID}, []string{customerName}...)
			customerIDNameList = append(customerIDNameList, customerIDAndName)
			customerIDList = append(customerIDList, customerID)
		}

	}
	return customerIDNameList, customerIDList
}

// Function to retrive PBX name(s) and ID(s) or just the PBX ID(s) from the view___pbx_detail view
func pbxSlice(dbDetail databaseFunctionParameter) ([][]string, []string) {
	// Get PBX name and ID from the database and append to slice
	var pbxIDNameList [][]string
	var pbxIDList []string

	var pbxID string
	var pbxName string

	pbxIDNameSQL, err := dbDetail.connection.Query(`SELECT
						          pbx_id,
                                                          pbx_name
                                                        FROM
                                                          yap.view___pbx_detail;`)

	// Error
	if err != nil {
		panic(err)
	}

	for pbxIDNameSQL.Next() {

		err = pbxIDNameSQL.Scan(
			&pbxID,
			&pbxName,
		)

		// Error
		if err != nil {
			panic(err)
		}

		var pbxIDAndName []string
		if pbxID != "1" {
			pbxIDAndName = append([]string{pbxID}, []string{pbxName}...)
			pbxIDNameList = append(pbxIDNameList, pbxIDAndName)
			pbxIDList = append(pbxIDList, pbxID)
		}

	}
	return pbxIDNameList, pbxIDList
}

// Function to retrive PBX name(s) and ID(s) or just the PBX ID(s) from the view___pbx_detail view based on customer ID
func pbxWhereSlice(dbDetail databaseFunctionParameter) ([][]string, []string) {
	// Get PBX name and ID from the database and append to slice
	var pbxIDNameList [][]string
	var pbxIDList []string

	var pbxID string
	var pbxName string

	pbxIDNameSQL, err := dbDetail.connection.Query(`SELECT
                                                          pbx_id,
                                                          pbx_name
                                                        FROM
                                                          yap.view___pbx_detail
                                                        WHERE
                                                          `+dbDetail.columnWhere+` = ?;`, dbDetail.columnWhereValue)

	// Error
	if err != nil {
		panic(err)
	}

	for pbxIDNameSQL.Next() {

		err = pbxIDNameSQL.Scan(
			&pbxID,
			&pbxName,
		)

		// Error
		if err != nil {
			panic(err)
		}

		var pbxIDAndName []string
		if pbxID != "1" {
			pbxIDAndName = append([]string{pbxID}, []string{pbxName}...)
			pbxIDNameList = append(pbxIDNameList, pbxIDAndName)
			pbxIDList = append(pbxIDList, pbxID)
		}

	}
	return pbxIDNameList, pbxIDList
}

// Function to retrive values from a single column from inside a database table, the data is appended to a slice
func singleColumnSlice(dbDetail databaseFunctionParameter) []string {
	// Get values from the database and append to a slice
	var singleColumnList []string
	var singleColumn string

	singleColumnSQL, err := dbDetail.connection.Query(`SELECT
                                                                ` + dbDetail.column + `
                                                              FROM
                                                                yap.` + dbDetail.table + `;`)

	// Error
	if err != nil {
		panic(err)
	}

	for singleColumnSQL.Next() {

		err = singleColumnSQL.Scan(
			&singleColumn,
		)

		// Error
		if err != nil {
			panic(err)
		}

		singleColumnList = append(singleColumnList, singleColumn)
	}
	return singleColumnList
}

// Function to add ext and PBX setup/rental/cease charges
func invoicePBXExtAdd(dbDetail databaseFunctionParameter, invoicePBXExt invoicePBXExtFunctionParameter) {
	// Convert string values to a float64 to use the math package to round to the nearest two decimal places
	sellPriceFloat64 := stringToFloat64(invoicePBXExt.sellPrice)

	dbDetail.connection.Query(`INSERT 
                                             INTO
                                           invoice_item (
                                             customer_id,
                                             pbx_id,
                                             service_product_name,
                                             tag,
                                             sell_price,
                                             sales_tax_rate,
                                             sales_tax_status,
                                             bill_item_once,
                                             item_on_hold,
                                             contract_length,
                                             contract_start_date
                                           )
                                           VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`,
		invoicePBXExt.customerID,
		invoicePBXExt.pbxID,
		invoicePBXExt.serviceProduct,
		nullSQL(invoicePBXExt.tag),
		math.Round(sellPriceFloat64*100)/100,
		invoicePBXExt.salesTaxRate,
		invoicePBXExt.salesTaxStatus,
		invoicePBXExt.billItemOnce,
		invoicePBXExt.itemOnHold,
		nullSQL(invoicePBXExt.contractLength),
		nullSQL(invoicePBXExt.contractStartDate),
	)
}

//----------------------------------------------------------------------------------------------------

// Functions to send invoices to accounting software via API

// Function to get access token via the accounting software API
func accessToken(accountingSoftware accountingSoftwareParameter) string {

	// Set data to POST
	data := url.Values{}
	data.Set("client_id", accountingSoftware.clientID)
	data.Set("client_secret", accountingSoftware.clientSecret)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", accountingSoftware.refreshToken)

	// Send POST request
	request, error := http.NewRequest("POST", accountingSoftware.url+`/token_endpoint`, strings.NewReader(data.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Accept", "application/json")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}

	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	jsonReturned := string(bodyBytes)
	// Parse JSON
	token := gjson.Get(jsonReturned, "access_token")
	// Trim " out of string
	tokenTrimmed := trimString(token.Raw, "\"")
	// Return the access token
	return string(tokenTrimmed)
}

// Function for invoice item JSON (this function is only accessed by the postInvoice)
func invoiceItemJSON(dbDetail databaseFunctionParameter) string {

	var (
		serviceProductName        string
		serviceProductType        string
		invoiceItemTag            string
		invoiceItemSellPrice      string
		invoiceItemSalesTaxRate   string
		invoiceItemSalesTaxStatus string
	)

	itemJSONSQL, err := dbDetail.connection.Query(`SELECT
                                                         service_product_name,
                                                         service_product_type,
                                                         invoice_item_tag,
                                                         invoice_item_sell_price,
                                                         invoice_item_sales_tax_rate,
                                                         invoice_item_sales_tax_status
                                                       FROM
                                                         yap.view___invoice_item;`)

	// Error
	if err != nil {
		panic(err)

	}

	var item strings.Builder

	for itemJSONSQL.Next() {

		err = itemJSONSQL.Scan(
			&serviceProductName,
			&serviceProductType,
			&invoiceItemTag,
			&invoiceItemSellPrice,
			&invoiceItemSalesTaxRate,
			&invoiceItemSalesTaxStatus,
		)

		// Error
		if err != nil {
			panic(err)
		}

		serviceProductName = serviceProductName + " (" + invoiceItemTag + ")"

		item.WriteString(`,{"description":"` + serviceProductName + `",
                                            "item_type":"` + serviceProductType + `", 
                                            "price":"` + invoiceItemSellPrice + `",
                                            "quantity":"1",
                                            "sales_tax_rate":"` + invoiceItemSalesTaxRate + `",
                                            "sales_tax_status":"` + invoiceItemSalesTaxStatus + `"
                                           }`)
	}
	return item.String()
}

// Function to post invoice via the accounting software API (this function is only accessed by the sendCustomerInvoice function)
func postInvoice(dbDetail databaseFunctionParameter, accountingSoftware accountingSoftwareParameter) int {
	item := invoiceItemJSON(dbDetail)

	dbDetail.table = "view___customer_detail"
	dbDetail.column = "customer_uk_based"
	dbDetail.columnWhere = "customer_id"
	dbDetail.columnWhereValue = accountingSoftware.customerID
	customerUKBased := selectWhere(dbDetail)

	// Determine the ecStatus from the SQL select statment
	var ecStatus string

	if customerUKBased == "yes" {
		ecStatus = "UK"
	} else if customerUKBased == "no" {
		ecStatus = "Reverse Charge"
	} else {
		ecStatus = ""
	}

	// Set invoiceDate to current date
	invoiceDate := currentDate()

	var dataJSON = []byte(`{
                                "invoice":
                                  {
                                  "contact":"` + accountingSoftware.customerID + `",
                                  "dated_on":"` + invoiceDate + `",
                                  "status":"Scheduled To Email",
                                  "currency":"` + accountingSoftware.currencyCode + `",
                                  "ec_status":"` + ecStatus + `",
                                  "send_thank_you_emails":true,
                                  "send_reminder_emails":true,
                                  "send_new_invoice_emails":true,
                                  "payment_terms_in_days":0,     
                                       "payment_methods": {
                                            "paypal": true,
                                            "stripe": true
                                        },
                                       "invoice_items":[
                                       {
                                       "description":"To see in depth details please use the portal",
                                       "item_type":"Comment", 
                                       "quantity":"0"
                                       }
                                       ` + item + `
                                ]
                        }
                        }
                        `)

	request, error := http.NewRequest("POST", accountingSoftware.url+`/invoices`, bytes.NewBuffer(dataJSON))
	request.Header.Add("Authorization", "Bearer "+accountingSoftware.accessToken)
	request.Header.Add("Content-Type", "application/json; charset=UTF-8")
	request.Header.Add("Accept", "application/json")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}

	defer response.Body.Close()

	// Return the HTTP status code
	httpStatusCode := response.StatusCode
	return httpStatusCode
}

// function to send invoice delete all invoice items that are set to bill once and update all invoice items that are on hold to being off hold
func sendCustomerInvoice(dbDetail databaseFunctionParameter, accountingSoftware accountingSoftwareParameter) {

	// Get an access token from accounting software API
	accountingSoftware.accessToken = accessToken(accountingSoftware)

	var invoiceItemCustomerID string

	customerIDSQL, err := dbDetail.connection.Query(`SELECT DISTINCT
                     	                                   customer_id
                                                         FROM
                                                           view___invoice_item;`)

	// Error
	if err != nil {
		panic(err)

	}

	for customerIDSQL.Next() {

		err = customerIDSQL.Scan(
			&invoiceItemCustomerID,
		)

		// Error
		if err != nil {
			panic(err)
		}

		accountingSoftware.customerID = invoiceItemCustomerID
		accountingSoftware.httpStatusCode = postInvoice(dbDetail, accountingSoftware)
	}

	// The HTTP status code response code determines if the invoices were sent succesfully
	if accountingSoftware.httpStatusCode == 201 {
		// Delete all invoice items that are set to bill once
		dbDetail.connection.Query("DELETE FROM `invoice_item` WHERE `bill_item_once` = 'yes';")

		// Update all invoice items that are on hold to being off hold
		dbDetail.connection.Query("UPDATE `invoice_item` SET `item_on_hold` = 'no' WHERE `item_on_hold` = 'yes';")
	} else if accountingSoftware.httpStatusCode == 400 {

	} else if accountingSoftware.httpStatusCode == 404 {

	}
}

//----------------------------------------------------------------------------------------------------

// Function to return slice of pbxLimitList
func pbxLimitSlice() []string {
	pbxLimitList := []string{"", "1", "2", "3", "4", "5", "10", "25", "50", "75", "100", "150", "200", "250", "500", "750", "1000", "1500", "2000", "2500", "5000"}
	return pbxLimitList
}

// Function to return slice of extLimitList
func extLimitSlice() []string {
	extLimitList := []string{"", "1", "2", "3", "4", "5", "10", "25", "50", "75", "100", "150", "200", "250", "500", "750", "1000", "1500", "2000", "2500", "5000"}
	return extLimitList
}

// Function to return slice of yesList
func yesSlice() []string {
	yesList := []string{"", "yes"}
	return yesList
}

// Function to return slice of yesNoList
func yesNoSlice() []string {
	yesNoList := []string{"", "yes", "no"}
	return yesNoList
}

// Function to return nested slice of invoiceColumnList
func invoiceColumnSlice() [][]string {
	invoiceColumnValueName := [][]string{
		{"address_line_1", "Invoice Address Line One"},
		{"address_line_2", "Invoice Address Line Two"},
		{"city_town_village", "Invoice City Town Village"},
		{"county_state_region", "Invoice County/State/Region"},
		{"postcode_zip_code", "Invoice Postcode/Zip Code"},
		{"country", "Invoice Country"},
		{"contact_email", "Invoice Contact Email"},
		{"contact_number", "Invoice Contact Number"},
	}
	return invoiceColumnValueName
}

// Function to return nested slice of siteColumnList
func siteColumnSlice() [][]string {
	siteColumnValueName := [][]string{
		{"address_line_1", "Site Address Line One"},
		{"address_line_2", "Site Address Line Two"},
		{"city_town_village", "Site City Town Village"},
		{"county_state_region", "Site County/State/Region"},
		{"postcode_zip_code", "Site Postcode/Zip Code"},
		{"country", "Site Country"},
		{"contact_email", "Site Contact Email"},
		{"contact_number", "Site Contact Number"},
	}
	return siteColumnValueName
}

// Function to return slice of codecAllowedSlice
func codecAllowedSlice() ([][]string, []string) {
	codecAllowedValueName := [][]string{
		{"alaw", "A-LAW"},
		{"ulaw", "U-lAW"},
	}

	var codecAllowedValue []string
	for _, value := range codecAllowedValueName {
		codecAllowedValue = append(codecAllowedValue, value[0])
	}

	return codecAllowedValueName, codecAllowedValue
}

// Function to return nested slice of dtmfModeSlice
func dtmfModeSlice() ([][]string, []string) {
	dtmfModeValueName := [][]string{
		{"rfc4733", "RFC 4733"},
		{"inband", "Inbound"},
		{"info", "Info"},
		{"auto", "Auto"},
		{"auto_info", "Auto Info"},
	}

	var dtmfModeValue []string
	for _, value := range dtmfModeValueName {
		dtmfModeValue = append(dtmfModeValue, value[0])
	}

	dtmfModeValue = append(dtmfModeValue, "")
	return dtmfModeValueName, dtmfModeValue
}

// Function to return nested slice of mediaEncryptionSlice
func mediaEncryptionSlice() ([][]string, []string) {
	mediaEncryptionValueName := [][]string{
		{"sdes", "SDES (Session Description Protocol Security Descriptions)"},
		{"no", "no"},
	}

	var mediaEncryptionValue []string
	for _, value := range mediaEncryptionValueName {
		mediaEncryptionValue = append(mediaEncryptionValue, value[0])
	}

	mediaEncryptionValue = append(mediaEncryptionValue, "")
	return mediaEncryptionValueName, mediaEncryptionValue
}

// Function to return nested slice of directMediaMethodSlice
func directMediaMethodSlice() ([][]string, []string) {
	directMediaMethodValueName := [][]string{
		{"invite", "Invite"},
		{"reinvite", "Reinvite"},
		{"update", "Update"},
	}

	var directMediaMethodValue []string
	for _, value := range directMediaMethodValueName {
		directMediaMethodValue = append(directMediaMethodValue, value[0])
	}

	directMediaMethodValue = append(directMediaMethodValue, "")
	return directMediaMethodValueName, directMediaMethodValue
}

// Function to return nested slice of callerIDPrivacySlice
func callerIDPrivacySlice() ([][]string, []string) {
	callerIDPrivacyValueName := [][]string{
		{"allowed_not_screened", "Allowed Not Screened"},
		{"allowed_passed_screen", "Allowed Passed Screen"},
		{"allowed_failed_screen", "Allowed Failed Screen"},
		{"allowed", "Allowed"},
		{"prohib_not_screened", "Prohib Not Screened"},
		{"prohib_passed_screen", "Prohib Passed Screen"},
		{"prohib_failed_screen", "Prohib Failed Screen"},
		{"prohib", "Prohib"},
		{"unavailable", "Unavailable"},
	}

	var callerIDPrivacyValue []string
	for _, value := range callerIDPrivacyValueName {
		callerIDPrivacyValue = append(callerIDPrivacyValue, value[0])
	}

	callerIDPrivacyValue = append(callerIDPrivacyValue, "")
	return callerIDPrivacyValueName, callerIDPrivacyValue
}

// Function to return slice of serviceProductTypeSlice
func serviceProductTypeSlice() []string {
	serviceProductTypeList := []string{"", "Services", "Products"}
	return serviceProductTypeList
}

//----------------------------------------------------------------------------------------------------

// Constants for validation messages
const validationMessageEmail string = " value must be a valid email address with a maxamium of 30 characters used"
const validationMessagePhoneNumber string = " value must be a valid phone number in e.164 format with a maxamium of 16 characters must be used"
const validationMessageNumber string = " value must be a number"
const validationMessageAlphaNum string = " value must be 1 to 30 characters and must only contain characters: a-z, A-Z or numbers"
const validationMessageAlphaNumEmpty string = " value can be empty or must be a maxamium of 30 characters and must only contain characters: a-z, A-Z or numbers"
const validationMessagePrice string = " value must be a decimal number with maxamium of 8 numbers"
const validationMessageTax string = " value must be a decimal number with maxamium of 6 numbers"
const validationMessageBoolean string = " value must be yes or no"
const validationMessageBooleanEmpty string = " value must be yes, no or empty"
const validationMessageDate string = " value must be in the date format DD-MM-YYYY"
const validationMessageIPAddress string = " value must be a valid IP address or empty"
const validationMessageFilePath string = " value must be a valid absloute path or empty"

const validationMessageAlreadyExist string = " already exists"
const validationMessageDoesNotExist string = " does not exist"
const validationMessageCreated string = " created"
const validationMessageNotCreated string = " not created"
const validationMessageDeleted string = " deleted"
const validationMessageNotDeleted string = " not deleted"
const validationMessageInvalidOption string = "Invalid option selected for "

// user-account page specfic HTML messages
const validationMessageAccountFirstName string = "A first name" + validationMessageAlphaNumEmpty
const validationMessageAccountLastName string = "A last name" + validationMessageAlphaNumEmpty
const validationMessageAccountType string = validationMessageInvalidOption + "account type"
const validationMessageAccountID string = "User account ID" + validationMessageDoesNotExist
const validationMessageAccountEmail string = "Account email" + validationMessageEmail
const validationMessageAccountCreated string = "User account" + validationMessageCreated

const validationMessageAccountDeleted string = "User account(s)" + validationMessageDeleted
const validationMessageAccountIndividualDeleted string = "User account" + validationMessageDeleted

const validationMessageAccountAlreadyExist string = "User account(s)" + validationMessageAlreadyExist

const validationMessageAccountYAPAdmin string = "Must be a YAP Admin (100) account with account ID 1"
const validationMessageAccountIDOne string = "User account with ID 1 cannot be used"
const validationMessageAccountColumn string = validationMessageInvalidOption + "account column"

const validationMessageAccountPBXDoesNotExist string = "User account(s) with PBX ID" + validationMessageDoesNotExist
const validationMessageAccountMultipleDeleted string = "User account(s)" + validationMessageDeleted
const validationMessageAccountMultipleNotDeleted string = "User account(s)" + validationMessageNotDeleted
const validationMessageAccountCustomerDoesNotExist string = "User account(s) with customer ID" + validationMessageDoesNotExist

// customer page specific HTML messages
const validationMessageCustomerID string = "A customer ID" + validationMessageAlphaNum
const validationMessageCustomerName string = "A customer name" + validationMessageAlphaNum
const validationMessageCustomerUKBased string = "UK based" + validationMessageBoolean
const validationMessageCustomerResellingMinutes string = "Reselling minutes" + validationMessageBoolean
const validationMessageCustomerConsumerType string = validationMessageInvalidOption + "consumer type"
const validationMessageCustomerUKVATRegistered string = "UK VAT registered" + validationMessageBoolean
const validationMessageCustomerUKVATRegisteredEmpty string = "When UK VAT registered is set to yes the UK VAT number cannot be empty"
const validationMessageCustomerUKVATNumber string = "UK VAT number" + validationMessageAlphaNumEmpty
const validationMessageCustomerPBXLimit string = validationMessageInvalidOption + "PBX limit"
const validationMessageCustomerPBXSalesTaxRate string = validationMessageInvalidOption + "PBX sales tax rate"
const validationMessageCustomerPBXSalesTaxStatus string = validationMessageInvalidOption + "PBX sales tax status"
const validationMessageCustomerExtSalesTaxRate string = validationMessageInvalidOption + "ext sales tax rate"
const validationMessageCustomerExtSalesTaxStatus string = validationMessageInvalidOption + "ext sales tax status"

const validationMessageCustomerSetupPrice string = "Setup price" + validationMessagePrice
const validationMessageCustomerRentalPrice string = "Rental price" + validationMessagePrice
const validationMessageCustomerCeasePrice string = "Cease price" + validationMessagePrice

const validationMessageCustomerCreated string = "Customer" + validationMessageCreated
const validationMessageCustomerNotCreated string = "Customer" + validationMessageNotCreated
const validationMessageCustomerlAlreadyExist string = "Customer" + validationMessageAlreadyExist
const validationMessageCustomerDoesNotExist string = "Customer" + validationMessageDoesNotExist
const validationMessageCustomerDeleted string = "Customer" + validationMessageDeleted
const validationMessageCustomerNotDeleted string = "Customer" + validationMessageNotDeleted

const validationMessageCustomerColumn string = validationMessageInvalidOption + "customer column"
const validationMessageCustomerEmail string = "A customer email" + validationMessageEmail
const validationMessageCustomerPhoneNumber string = "A customer phone number" + validationMessagePhoneNumber

// PBX page specific HTML messages
const validationMessagePBXName string = "A PBX name" + validationMessageAlphaNum
const validationMessagePBXExtLimit string = validationMessageInvalidOption + "ext limit"
const validationMessagePBXCreated string = "PBX" + validationMessageCreated
const validationMessagePBXNotCreated string = "PBX" + validationMessageNotCreated

const validationMessagePBXColumn string = validationMessageInvalidOption + "PBX column"
const validationMessagePBXDeleted string = "PBX" + validationMessageDeleted
const validationMessagePBXNotDeleted string = "PBX" + validationMessageNotDeleted
const validationMessagePBXSiteEmail string = "A PBX site email" + validationMessageEmail
const validationMessagePBXSitePhoneNumber string = "A PBX site phone number" + validationMessagePhoneNumber
const validationMessagePBXMaxPBX string = "Max amount of PBXs allowed for the customer"
const validationMessagePBXIDOne string = "PBX with ID 1 cannot be used"

// Ext page specific HTML messages
const validationMessageExt string = "Ext" + validationMessageAlphaNum
const validationMessageExtCodecAllowed string = validationMessageInvalidOption + "codec allowed"
const validationMessageExtDTMFMode string = validationMessageInvalidOption + "DTMF mode"
const validationMessageExtCallGroup string = "Call group" + validationMessageAlphaNumEmpty
const validationMessageExtPickupGroup string = "Pickup Group" + validationMessageAlphaNumEmpty

const validationMessageExtMediaEncryption string = validationMessageInvalidOption + "media encryption"
const validationMessageExtICESupport string = "Ice support" + validationMessageBooleanEmpty
const validationMessageExtDirectMedia string = "Direct media" + validationMessageBooleanEmpty
const validationMessageExtDirectMediaMethod string = validationMessageInvalidOption + "direct media method"
const validationMessageExtRewriteContact string = "Rewrite contact" + validationMessageBooleanEmpty
const validationMessageExtRTPSymmetric string = "RTP symmetric" + validationMessageBooleanEmpty
const validationMessageExtForceRPort string = "Force RPort" + validationMessageBooleanEmpty
const validationMessageExtRestrictExt string = "Restrict to ext" + validationMessageIPAddress

const validationMessageExtAllowTransfer string = "Allow Transfer" + validationMessageBooleanEmpty
const validationMessageExtCallerID string = "Caller ID" + validationMessageAlphaNumEmpty
const validationMessageExtCallerIDPrivacy string = validationMessageInvalidOption + "caller ID privacy"
const validationMessageExtContactUser string = "SIP header - contact user" + validationMessageAlphaNumEmpty
const validationMessageExtFromUser string = "SIP header - from user" + validationMessageAlphaNumEmpty
const validationMessageExtFromDomain string = "SIP header - from domain" + validationMessageAlphaNumEmpty
const validationMessageExtStirShaken string = validationMessageInvalidOption + "stir/shaken"
const validationMessageExtStirShakenProfile string = "Stir/shaken profile" + validationMessageFilePath

const validationMessageExtAlreadyExist string = "Ext" + validationMessageAlreadyExist
const validationMessageExtCreated string = "Ext" + validationMessageCreated
const validationMessageExtNotCreated string = "Ext" + validationMessageNotCreated
const validationMessageExtMaxExt string = "Max amount of extensions allowed for the PBX"
const validationMessageExtColumn string = validationMessageInvalidOption + "ext column"
const validationMessageExtDeleted string = "Ext" + validationMessageDeleted
const validationMessageExtNotDeleted string = "Ext" + validationMessageNotDeleted + " or" + validationMessageDoesNotExist
const validationMessageExtDoesNotExist string = "Ext" + validationMessageDoesNotExist

// invoice page specific HTML messages
const validationMessageInvoiceServiceProduct string = validationMessageInvalidOption + "service/product"
const validationMessageInvoiceServiceProductTag string = "Service/product tag" + validationMessageAlphaNumEmpty
const validationMessageInvoiceItemPrice string = "Item price" + validationMessagePrice
const validationMessageInvoiceSalesTaxRate string = validationMessageInvalidOption + "sales tax rate"
const validationMessageInvoiceSalesTaxStatus string = validationMessageInvalidOption + "sales tax status"
const validationMessageInvoiceBillItemOnce string = validationMessageInvalidOption + "bill item once"
const validationMessageInvoiceItemOnHold string = validationMessageInvalidOption + "item on hold"
const validationMessageInvoiceContractStartDate string = "Contract start date" + validationMessageDate
const validationMessageInvoiceContractStartDateEmpty string = "If contract length is not empty then contract start date must have a value"

const validationMessageInvoice string = validationMessageInvalidOption + "invoice"
const validationMessageInvoiceDoesNotExist string = "Invoice" + validationMessageDoesNotExist
const validationMessageInvoiceDeleted string = "Invoice" + validationMessageDeleted
const validationMessageInvoiceNotDeleted string = "Invoice" + validationMessageNotDeleted
const validationMessageInvoiceID string = "Invoice ID" + validationMessageNumber

// service-product page specific HTML messages
const validationMessageServiceProductName string = "Service/product name" + validationMessageAlphaNum
const validationMessageServiceProductType string = validationMessageInvalidOption + "service/product type"
const validationMessageServiceProductYAP string = "Service/product supplier name cannot be ⊛ YAP (Yet Another PBX) ⊛"
const validationMessageServiceProductSupplierName string = validationMessageInvalidOption + "service/product supplier name"
const validationMessageServiceProductCreated string = "Service/product" + validationMessageCreated
const validationMessageServiceProductNotCreated string = "Service/product" + validationMessageNotCreated
const validationMessageServiceProductAlreadyExist string = "Service/product" + validationMessageAlreadyExist

const validationMessageSupplierYAP string = "Supplier name cannot be ⊛ YAP (Yet Another PBX) ⊛"
const validationMessageSupplierAlreadyExist string = "Supplier name" + validationMessageAlreadyExist
const validationMessageSupplierCreated string = "Supplier" + validationMessageCreated
const validationMessageSupplierNotCreated string = "Supplier" + validationMessageNotCreated
const validationMessageSupplierName string = "Supplier name" + validationMessageAlphaNum

const validationMessageSalesTaxRateAlreadyExist string = "Sales tax rate" + validationMessageAlreadyExist
const validationMessageSalesTaxRateCreated string = "Sales tax rate" + validationMessageCreated
const validationMessageSalesTaxRateNotCreated string = "Sales tax rate" + validationMessageNotCreated
const validationMessageSalesTaxRate string = "Sales tax rate" + validationMessageTax

const validationMessageServiceProductDoesNotExist string = "Service/product" + validationMessageDoesNotExist
const validationMessageServiceProductDeleted string = "Service/product" + validationMessageDeleted
const validationMessageServiceProductNotDeleted string = "Service/product" + validationMessageNotDeleted

const validationMessageSupplierDoesNotExist string = "Supplier" + validationMessageDoesNotExist
const validationMessageSupplierDeleted string = "Supplier" + validationMessageDeleted
const validationMessageSupplierNotDeleted string = "Supplier" + validationMessageNotDeleted

const validationMessageSalesTaxRateDoesNotExist string = "Sales tax rate" + validationMessageDoesNotExist
const validationMessageSalesTaxRateDeleted string = "Sales tax rate" + validationMessageDeleted
const validationMessageSalesTaxRateNotDeleted string = "Sales tax rate" + validationMessageNotDeleted

const validationMessageServiceProductID string = "Service/product ID" + validationMessageNumber
const validationMessageServiceProductColumn string = validationMessageInvalidOption + "service/product column"

const validationMessageSupplierExistingValue string = "Supplier" + validationMessageAlphaNum
const validationMessageSupplierColumn string = validationMessageInvalidOption + "supplier column"

const validationMessageSalesTaxRateColumn string = validationMessageInvalidOption + "sales tax rate column"
const validationMessageSalesTaxRateTax string = "Sales tax rate" + validationMessageTax

// General/multi-page HTML messsages
const validationMessageCustomer string = validationMessageInvalidOption + "customer"
const validationMessageEmailAlreadyExist string = "Email" + validationMessageAlreadyExist
const validationMessageEmailDoesNotExist string = "Email" + validationMessageDoesNotExist
const validationMessagePBX string = validationMessageInvalidOption + "PBX"
const validationMessagePBXAlreadyExist string = "PBX" + validationMessageAlreadyExist
const validationMessagePBXDoesNotExist string = "PBX" + validationMessageDoesNotExist
const validationMessageGenericInvalidOption string = "Invalid option selected"
const validationMessageGenericAlphaNumEmpty string = "The" + validationMessageAlphaNumEmpty
const validationMessageGenericAlphaNum string = "The" + validationMessageAlphaNum
const validationMessageGenericPrice string = "Price" + validationMessagePrice
const validationMessageConfirmation string = "Confirmation must be yes"
const validationMessageContractLength string = validationMessageInvalidOption + "contract length"
const validationMessageInvalid string = "Invalid"

const validationMessageAddresslineOne string = "Address line one" + validationMessageAlphaNumEmpty
const validationMessageAddresslineTwo string = "Address line two" + validationMessageAlphaNumEmpty
const validationMessageCityTownVillage string = "City/town/village" + validationMessageAlphaNumEmpty
const validationMessageCountyStateRegion string = "County/state/region" + validationMessageAlphaNumEmpty
const validationMessagePostcodeZipCode string = "Postcode zip code" + validationMessageAlphaNumEmpty
const validationMessageCountry string = "Country" + validationMessageAlphaNumEmpty
const validationMessageCustomerSiteEmail string = "A customer site email" + validationMessageEmail
const validationMessageCustomerSitePhoneNumber string = "A customer site phone number" + validationMessagePhoneNumber
const validationMessageInvoiceEmail string = "A customer invoice email" + validationMessageEmail
const validationMessageInvoicePhoneNumber string = "A customer invoice phone number" + validationMessagePhoneNumber

// Function to validate user input utlising the Go validator version 10 package
func validateInput(value string, valueType string) (validation bool) {
	validateInput := validator.New()
	// Conditional statments are used for each type of validation needed
	if valueType == "email" {
		validateInputEmailErr := validateInput.Var(value, "email,max=200")
		if validateInputEmailErr != nil {
			validation = false
			return
		} else {
			validation = true
			return
		}
	} else if valueType == "phoneNumber" {
		validateInputPhoneNumberErr := validateInput.Var(value, "e164,max=16")
		if validateInputPhoneNumberErr != nil {
			validation = false
			return
		} else {
			validation = true
			return
		}
	} else if valueType == "number" {
		validateInputNumberErr := validateInput.Var(value, "number,min=1,max=30")
		if validateInputNumberErr != nil {
			validation = false
			return
		} else {
			validation = true
			return
		}
	} else if valueType == "alphaNum" {
		validateInputAlphaNumErr := validateInput.Var(value, "alphanumspace,min=1,max=30")
		validateInputSymbolErr := validateInput.Var(value, "excludes=`!\"£$%^&*()-_=+{}[];:@'#~\\.<>/?")
		if validateInputAlphaNumErr != nil || validateInputSymbolErr != nil {
			validation = false
			return
		} else {
			validation = true
			return
		}
	} else if valueType == "alphaNumEmpty" {
		validateInputAlphaNumEmptyErr := validateInput.Var(value, "ascii,max=30")
		validateInputSymbolErr := validateInput.Var(value, "excludes=`!\"£$%^&*()-_=+{}[];:@'#~\\.<>/?")
		if validateInputAlphaNumEmptyErr != nil || validateInputSymbolErr != nil {
			validation = false
			return
		} else {
			validation = true
			return
		}
	} else if valueType == "price" {
		validateInputNumberErr := validateInput.Var(value, "numeric,max=9")
		validateInputDecimalErr := validateInput.Var(value, "contains=.")
		if validateInputNumberErr != nil || validateInputDecimalErr != nil {
			validation = false
			return
		} else {
			validation = true
			return
		}
	} else if valueType == "tax" {
		validateInputNumberErr := validateInput.Var(value, "numeric,max=6")
		validateInputDecimalErr := validateInput.Var(value, "contains=.")
		if validateInputNumberErr != nil || validateInputDecimalErr != nil {
			validation = false
			return
		} else {
			validation = true
			return
		}
	} else if valueType == "date" {
		validateInputDateErr := validateInput.Var(value, "omitempty,datetime=2006-01-02")
		if validateInputDateErr != nil {
			validation = false
			return
		} else {
			validation = true
			return
		}
	} else if valueType == "ipAddress" {
		validateInputIPAddressErr := validateInput.Var(value, "omitempty,ip_addr")
		if validateInputIPAddressErr != nil {
			validation = false
			return
		} else {
			validation = true
			return
		}
	} else if valueType == "filePath" {
		validateInputDirErr := validateInput.Var(value, "omitempty,dir")
		if validateInputDirErr != nil {
			validation = false
			return
		} else {
			validation = true
			return
		}
	} else if valueType == "extension" {
		validateInputExtensionErr := validateInput.Var(value, "ascii,max=30")
		validateInputSymbolErr := validateInput.Var(value, "excludes=`!\"£$%^&*()_=+{}[];:@'#~\\.<>/?")
		if validateInputExtensionErr != nil || validateInputSymbolErr != nil {
			validation = false
			return
		} else {
			validation = true
			return
		}
	} else {
		panic("The validateInput function can only take the following arguments: email, phoneNumber, number, alphaNum, alphaNumEmpty, price, tax, data, ipAddress, filePath or extension")
	}
}

//----------------------------------------------------------------------------------------------------

// Main menu page functions

// YAP account function
func mainMenuYapAccount(w http.ResponseWriter, dbYapAccount databaseFunctionParameter) {

	var dbTotalTableCount databaseFunctionParameter
	dbTotalTableCount.connection = dbYapAccount.connection
	dbTotalTableCount.database = dbYapAccount.database

	var dbTotalTableCountWhere databaseFunctionParameter
	dbTotalTableCountWhere.connection = dbYapAccount.connection
	dbTotalTableCountWhere.database = dbYapAccount.database
	dbTotalTableCountWhere.table = "user_account"
	dbTotalTableCountWhere.columnWhere = "user_account_type_id"

	fmt.Fprintf(w, "</table>")
	fmt.Fprintf(w, "<table id=\"table\" class=\"table-main-menu\">")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>Total Customers</th>")
	fmt.Fprintf(w, "    <th>Total PBXs</th>")
	fmt.Fprintf(w, "    <th>Total Exts</th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "  <tr>")
	dbTotalTableCount.table = "view___customer_detail"
	dbTotalTableCount.countMinusOne = true
	fmt.Fprintf(w, "    <td>"+totalTableCount(dbTotalTableCount)+"</td>")
	dbTotalTableCount.table = "view___pbx_detail"
	dbTotalTableCount.countMinusOne = true
	fmt.Fprintf(w, "    <td>"+totalTableCount(dbTotalTableCount)+"</td>")
	dbTotalTableCount.table = "view___sip_extension_detail"
	dbTotalTableCount.countMinusOne = false
	fmt.Fprintf(w, "    <td>"+totalTableCount(dbTotalTableCount)+"</td>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "</table>")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<table id=\"table\" class=\"table-main-menu\">")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>Total YAP<br>Admin<br>Accounts<br>(Type ID: 100)</th>")
	fmt.Fprintf(w, "    <th>Total Customer<br>Admin<br>Accounts<br>(Type ID: 200)</th>")
	fmt.Fprintf(w, "    <th>Total Customer<br>Regular<br>Accounts<br>(Type ID: 201)</th>")
	fmt.Fprintf(w, "    <th>Total PBX<br>Admin<br>Accounts<br>(Type ID: 300)</th>")
	fmt.Fprintf(w, "    <th>Total PBX<br>Regular<br>Accounts<br>(Type ID: 301)</th>")
	fmt.Fprintf(w, "    <th>Total PBX<br>Read Only<br>Accounts<br>(Type ID: 302)</th>")
	fmt.Fprintf(w, "    <th>Total Customer<br>Invoice<br>Accounts<br>(Type ID: 400)</th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "  <tr>")
	dbTotalTableCountWhere.columnWhereValue = "100"
	fmt.Fprintf(w, "    <td>"+totalTableCountWhere(dbTotalTableCountWhere)+"</td>")
	dbTotalTableCountWhere.columnWhereValue = "200"
	fmt.Fprintf(w, "    <td>"+totalTableCountWhere(dbTotalTableCountWhere)+"</td>")
	dbTotalTableCountWhere.columnWhereValue = "201"
	fmt.Fprintf(w, "    <td>"+totalTableCountWhere(dbTotalTableCountWhere)+"</td>")
	dbTotalTableCountWhere.columnWhereValue = "300"
	fmt.Fprintf(w, "    <td>"+totalTableCountWhere(dbTotalTableCountWhere)+"</td>")
	dbTotalTableCountWhere.columnWhereValue = "301"
	fmt.Fprintf(w, "    <td>"+totalTableCountWhere(dbTotalTableCountWhere)+"</td>")
	dbTotalTableCountWhere.columnWhereValue = "302"
	fmt.Fprintf(w, "    <td>"+totalTableCountWhere(dbTotalTableCountWhere)+"</td>")
	dbTotalTableCountWhere.columnWhereValue = "400"
	fmt.Fprintf(w, "    <td>"+totalTableCountWhere(dbTotalTableCountWhere)+"</td>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "</table>")

}

// Customer account function
func mainMenuCustomerAccount(w http.ResponseWriter, dbCustomerAccount databaseFunctionParameter) {

	result, err := dbCustomerAccount.connection.Query(`SELECT
					                  customer_name,
					                  customer_id,
					                  customer_site_address_line_1,
					                  customer_site_address_line_2,
					                  customer_site_city_town_village,
					                  customer_site_county_state_region,
					                  customer_site_postcode_zip_code,
					                  customer_site_country,
					                  customer_site_contact_email,
					                  customer_site_contact_number,
					                  customer_invoice_address_line_1,
					                  customer_invoice_address_line_2,
					                  customer_invoice_city_town_village,
					                  customer_invoice_county_state_region,
					                  customer_invoice_postcode_zip_code,
					                  customer_invoice_country,
					                  customer_invoice_contact_email,
					                  customer_invoice_contact_number,
					                  pbx_id
					                FROM
					                  yap.view___account_detail
					                WHERE
					                  user_account_email = ?;`, dbCustomerAccount.columnWhereValue)

	// Error
	if err != nil {
		panic(err)
	}

	for result.Next() {
		var (
			customerName                     string
			customerID                       string
			customerSiteAddressLine1         string
			customerSiteAddressLine2         string
			customerSiteCityTownVillage      string
			customerSiteCountyStateRegion    string
			customerSitePostcodeZipCode      string
			customerSiteCountry              string
			customerSiteContactEmail         string
			customerSiteContactNumber        string
			customerInvoiceAddressLine1      string
			customerInvoiceAddressLine2      string
			customerInvoiceCityTownVillage   string
			customerInvoiceCountyStateRegion string
			customerInvoicePostcodeZipCode   string
			customerInvoiceCountry           string
			customerInvoiceContactEmail      string
			customerInvoiceContactNumber     string
			pbxID                            string
		)

		err = result.Scan(
			&customerName,
			&customerID,
			&customerSiteAddressLine1,
			&customerSiteAddressLine2,
			&customerSiteCityTownVillage,
			&customerSiteCountyStateRegion,
			&customerSitePostcodeZipCode,
			&customerSiteCountry,
			&customerSiteContactEmail,
			&customerSiteContactNumber,
			&customerInvoiceAddressLine1,
			&customerInvoiceAddressLine2,
			&customerInvoiceCityTownVillage,
			&customerInvoiceCountyStateRegion,
			&customerInvoicePostcodeZipCode,
			&customerInvoiceCountry,
			&customerInvoiceContactEmail,
			&customerInvoiceContactNumber,
			&pbxID,
		)

		// Error
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, "<table id=\"table\" class=\"table-main-menu\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>Customer Name and ID</th>")
		fmt.Fprintf(w, "    <th>Customers Total PBXs</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td>Customer Name: "+customerName+"<br><br>Customer ID: "+customerID+"</td>")
		var dbTotalTableCountWhere databaseFunctionParameter
		dbTotalTableCountWhere.connection = dbCustomerAccount.connection
		dbTotalTableCountWhere.database = dbCustomerAccount.database
		dbTotalTableCountWhere.table = "pbx"
		dbTotalTableCountWhere.columnWhere = "id"
		dbTotalTableCountWhere.columnWhereValue = pbxID
		fmt.Fprintf(w, "    <td>"+totalTableCountWhere(dbTotalTableCountWhere)+"</td>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<table id=\"table\" class=\"table-main-menu\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>Customer Site Address</th>")
		fmt.Fprintf(w, "    <th>Customer Site Email</th>")
		fmt.Fprintf(w, "    <th>Customer Site Phone Number</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td style=\"text-align: left;\">"+customerSiteAddressLine1+"<br>"+customerSiteAddressLine2+"<br>"+customerSiteCityTownVillage+"<br>"+customerSiteCountyStateRegion+"<br><br>"+customerSitePostcodeZipCode+"<br><br>"+customerSiteCountry+"</td>")
		fmt.Fprintf(w, "    <td><a href=\"mailto:"+customerSiteContactEmail+"\">"+customerSiteContactEmail+"</a></td>")
		fmt.Fprintf(w, "    <td><a href=\"tel:"+customerSiteContactNumber+"\">"+customerSiteContactNumber+"</a></td>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<table id=\"table\" class=\"table-main-menu\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>Customer Invoice Address</th>")
		fmt.Fprintf(w, "    <th>Customer Invoice Email</th>")
		fmt.Fprintf(w, "    <th>Customer Invoice Phone Number</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td style=\"text-align: left;\">"+customerInvoiceAddressLine1+"<br>"+customerInvoiceAddressLine2+"<br>"+customerInvoiceCityTownVillage+"<br>"+customerInvoiceCountyStateRegion+"<br><br>"+customerInvoicePostcodeZipCode+"<br><br>"+customerInvoiceCountry+"</td>")
		fmt.Fprintf(w, "    <td><a href=\"mailto:"+customerInvoiceContactEmail+"\">"+customerInvoiceContactEmail+"</a></td>")
		fmt.Fprintf(w, "    <td><a href=\"tel:"+customerInvoiceContactNumber+"\">"+customerInvoiceContactNumber+"</a></td>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
	}
}

// PBX account function
func mainMenuPBXAccount(w http.ResponseWriter, dbPBXAccount databaseFunctionParameter) {

	result, err := dbPBXAccount.connection.Query(`SELECT
					                pbx_name,
					                pbx_id,
					                pbx_site_address_line_1,
					                pbx_site_address_line_2,
					                pbx_site_city_town_village,
					                pbx_site_county_state_region,
					                pbx_site_postcode_zip_code,
					                pbx_site_country,
					                pbx_site_contact_email,
					                pbx_site_contact_number
						      FROM
						        yap.view___account_detail
					  	      WHERE
					  	        user_account_email = ?;`, dbPBXAccount.columnWhereValue)

	// Error
	if err != nil {
		panic(err)
	}

	for result.Next() {

		var (
			pbxName                  string
			pbxID                    string
			pbxSiteAddressLine1      string
			pbxSiteAddressLine2      string
			pbxSiteCityTownVillage   string
			pbxSiteCountyStateRegion string
			pbxSitePostcodeZipCode   string
			pbxSiteCountry           string
			pbxSiteContactEmail      string
			pbxSiteContactNumber     string
		)

		err = result.Scan(
			&pbxName,
			&pbxID,
			&pbxSiteAddressLine1,
			&pbxSiteAddressLine2,
			&pbxSiteCityTownVillage,
			&pbxSiteCountyStateRegion,
			&pbxSitePostcodeZipCode,
			&pbxSiteCountry,
			&pbxSiteContactEmail,
			&pbxSiteContactNumber,
		)

		// Error
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, "<table id=\"table\" class=\"table-main-menu\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>PBX Name and ID</th>")
		fmt.Fprintf(w, "    <th>Total Exts in PBX</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td>PBX Name: "+pbxName+"<br><br>PBX ID: "+pbxID+"</td>")
		var dbTotalTableCountWhere databaseFunctionParameter
		dbTotalTableCountWhere.connection = dbPBXAccount.connection
		dbTotalTableCountWhere.database = dbPBXAccount.database
		dbTotalTableCountWhere.table = "view___sip_extension_detail"
		dbTotalTableCountWhere.columnWhere = "pbx_id"
		dbTotalTableCountWhere.columnWhereValue = pbxID
		fmt.Fprintf(w, "    <td>"+totalTableCountWhere(dbTotalTableCountWhere)+"</td>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<table id=\"table\" class=\"table-main-menu\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>PBX Site Address</th>")
		fmt.Fprintf(w, "    <th>PBX Site Email</th>")
		fmt.Fprintf(w, "    <th>PBX Site Phone Number</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td style=\"text-align: left;\">"+pbxSiteAddressLine1+"&nbsp<br>"+pbxSiteAddressLine2+"<br>"+pbxSiteCityTownVillage+"<br>"+pbxSiteCountyStateRegion+"<br><br>"+pbxSitePostcodeZipCode+"<br><br>"+pbxSiteCountry+"</td>")
		fmt.Fprintf(w, "    <td><a href=\"mailto:"+pbxSiteContactEmail+"\">"+pbxSiteContactEmail+"</a></td>")
		fmt.Fprintf(w, "    <td><a href=\"tel:"+pbxSiteContactNumber+"\">"+pbxSiteContactNumber+"</a></td>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")

	}

}

// Function for main menu page user information
func mainMenuUserInformation(w http.ResponseWriter, dbUserInformation databaseFunctionParameter, genDetail generalFunctionParameter) {

	result, err := dbUserInformation.connection.Query(`SELECT
							     user_account_id,
					                     user_account_first_name,
					                     user_account_last_name,
					                     user_account_email,
					                     user_account_type,
					                     user_account_date_time_added,
					                     user_account_type_permission
					                   FROM
					                     yap.view___account_detail
					                   WHERE
					                     user_account_email = ?;`, dbUserInformation.columnWhereValue)

	// Error
	if err != nil {
		panic(err)
	}

	for result.Next() {
		var (
			userAccountID             string
			userAccountFirstName      string
			userAccountLastName       string
			userAccountEmail          string
			userAccountType           string
			userAccountDateTimeAdded  string
			userAccountTypePermission string
		)

		err = result.Scan(
			&userAccountID,
			&userAccountFirstName,
			&userAccountLastName,
			&userAccountEmail,
			&userAccountType,
			&userAccountDateTimeAdded,
			&userAccountTypePermission,
		)

		// Error
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, "<div>")
		fmt.Fprintf(w, "<table id=\"table\" class=\"table-main-menu\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table id=\"table\" class=\"table-main-menu\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <th>Account ID</th>")
		fmt.Fprintf(w, "          <th>Name</th>")
		fmt.Fprintf(w, "          <th>Email</th>")
		fmt.Fprintf(w, "          <th>Account Type</th>")
		fmt.Fprintf(w, "          <th>Account Created</th>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>"+userAccountID+"</td>")
		fmt.Fprintf(w, "          <td>"+userAccountFirstName+"<br>"+userAccountLastName+"</td>")
		fmt.Fprintf(w, "          <td>"+userAccountEmail+"</td>")
		fmt.Fprintf(w, "          <td>"+userAccountType+"</td>")
		fmt.Fprintf(w, "          <td>"+formatDateTime(userAccountDateTimeAdded)+"</td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><button onclick=\"toggleAccountDetail() \"class=\"button-general\">&nbsp Show/Hide More Account Details &nbsp</button></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		// Account detail tables
		fmt.Fprintf(w, "</div>")
		fmt.Fprintf(w, "<div id=\"account-detail-div\" style=\"display:none\">")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<table id=\"table\" class=\"table-main-menu\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>Account Type Permissions - "+userAccountType+"</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td style=\"text-align: left\">"+userAccountTypePermission+"</td>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
	}

	var dbDetail databaseFunctionParameter
	dbDetail.connection = dbUserInformation.connection
	dbDetail.database = dbUserInformation.database
	dbDetail.columnWhereValue = dbUserInformation.columnWhereValue

	if genDetail.userTypeID == "100" {
		fmt.Fprintf(w, "<br>")
		mainMenuYapAccount(w, dbDetail)
	} else if genDetail.userTypeID == "200" || genDetail.userTypeID == "201" || genDetail.userTypeID == "400" {
		fmt.Fprintf(w, "<br>")
		mainMenuCustomerAccount(w, dbDetail)
	} else if genDetail.userTypeID == "300" || genDetail.userTypeID == "301" || genDetail.userTypeID == "302" {
		fmt.Fprintf(w, "<br>")
		mainMenuPBXAccount(w, dbDetail)
	} else {
	}
	fmt.Fprintf(w, "</div>")
	var toggleDivJSArgument jsFunctionParameter
	toggleDivJSArgument.funcNameJS = "toggleAccountDetail"
	toggleDivJSArgument.divID = "account-detail-div"
	toggleDivJS(w, toggleDivJSArgument)
}

type mainMenuParameter struct {
	writeHTTP  http.ResponseWriter
	buttonName string
	hyperlink  string
	headerCSS  string
	buttonCSS  string
}

// Function for main menu page buttons (hyperlinks)
func mainMenuButton(mainMenu mainMenuParameter) {
	fmt.Fprintf(mainMenu.writeHTTP, "&nbsp")
	fmt.Fprintf(mainMenu.writeHTTP, "<h2 class=\""+mainMenu.headerCSS+"\">")
	fmt.Fprintf(mainMenu.writeHTTP, "<a href=\""+mainMenu.hyperlink+"\" class=\"button-general button-main-menu "+mainMenu.buttonCSS+"\"><p>"+mainMenu.buttonName+"</p></a>")
	fmt.Fprintf(mainMenu.writeHTTP, "</h2>")
	fmt.Fprintf(mainMenu.writeHTTP, "&nbsp")
}

//----------------------------------------------------------------------------------------------------

// User account page functions

// Function uses different CSS class for table td depending on tag content
func userAccountListTdHTML(w http.ResponseWriter, userAccountID string, userAccountTypeID string, tagContent string) {
	if userAccountID == "1" && userAccountTypeID == "100" {
		fmt.Fprintf(w, "          <td class=\"td-yap-admin-id-1\"><b>"+tagContent+"</b></td>")
	} else if userAccountTypeID == "100" {
		fmt.Fprintf(w, "          <td class=\"td-yap-admin\"><b>"+tagContent+"</b></td>")
	} else if userAccountTypeID == "200" || userAccountTypeID == "201" {
		fmt.Fprintf(w, "	  <td class=\"td-customer\"><b>"+tagContent+"</b></td>")
	} else if userAccountTypeID == "300" || userAccountTypeID == "301" || userAccountTypeID == "302" {
		fmt.Fprintf(w, "          <td class=\"td-pbx\"><b>"+tagContent+"</b></td>")
	} else if userAccountTypeID == "400" {
		fmt.Fprintf(w, "          <td class=\"td-invoice\"><b>"+tagContent+"</b></td>")
	}
}

// Function to list user accounts
func userAccountList(w http.ResponseWriter, dbDetail databaseFunctionParameter, genDetail generalFunctionParameter) {

	// Only account type ID 100, 200, 201, 300, 301, 302, 400 should be able to use this function
	if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" || genDetail.userTypeID == "300" || genDetail.userTypeID == "301" || genDetail.userTypeID == "302" || genDetail.userTypeID == "400" {

		var (
			userAccountID            string
			userAccountFirstName     string
			userAccountLastName      string
			userAccountEmail         string
			userAccountType          string
			userAccountTypeID        string
			userAccountDateTimeAdded string
			customerID               string
			customerName             string
			pbxID                    string
			pbxName                  string
		)

		ownUserAccountSQL, err := dbDetail.connection.Query(`SELECT
							       user_account_id,
							       user_account_first_name,
							       user_account_last_name,
							       user_account_email,
							       user_account_type,
							       user_account_date_time_added,
							       customer_id,
							       pbx_id
							     FROM
							       yap.view___account_detail
							     WHERE
							       user_account_email = ?;`, dbDetail.columnWhereValue)

		// Error
		if err != nil {
			panic(err)
		}

		for ownUserAccountSQL.Next() {

			err = ownUserAccountSQL.Scan(
				&userAccountID,
				&userAccountFirstName,
				&userAccountLastName,
				&userAccountEmail,
				&userAccountType,
				&userAccountDateTimeAdded,
				&customerID,
				&pbxID,
			)

			// Error
			if err != nil {
				panic(err)
			}

			var dbTotalTableCount databaseFunctionParameter
			dbTotalTableCount.connection = dbDetail.connection
			dbTotalTableCount.database = dbDetail.database

			var dbTotalTableCountWhere databaseFunctionParameter
			dbTotalTableCountWhere.connection = dbDetail.connection
			dbTotalTableCountWhere.database = dbDetail.database
			dbTotalTableCountWhere.table = "user_account"
			dbTotalTableCountWhere.columnWhere = "user_account_type_id"

			if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" || genDetail.userTypeID == "300" || genDetail.userTypeID == "301" {
				fmt.Fprintf(w, "<table id=\"table\" class=\"table-user-account\">")
				fmt.Fprintf(w, "  <tr>")
				if genDetail.userTypeID == "100" {
					fmt.Fprintf(w, "    <th>Total YAP<br>Admin<br>Accounts<br>(Type ID: 100)</th>")
				}
				if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" {
					fmt.Fprintf(w, "    <th>Total Customer<br>Admin<br>Accounts<br>(Type ID: 200)</th>")
				}
				if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
					fmt.Fprintf(w, "    <th>Total Customer<br>Regular<br>Accounts<br>(Type ID: 201)</th>")
				}
				if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" || genDetail.userTypeID == "300" {
					fmt.Fprintf(w, "    <th>Total PBX<br>Admin<br>Accounts<br>(Type ID: 300)</th>")
				}
				if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" || genDetail.userTypeID == "300" || genDetail.userTypeID == "301" {
					fmt.Fprintf(w, "    <th>Total PBX<br>Regular<br>Accounts<br>(Type ID: 301)</th>")
					fmt.Fprintf(w, "    <th>Total PBX<br>Read Only<br>Accounts<br>(Type ID: 302)</th>")
				}
				if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" {
					fmt.Fprintf(w, "    <th>Total Customer<br>Invoice<br>Accounts<br>(Type ID: 400)</th>")
				}
				fmt.Fprintf(w, "  </tr>")
				fmt.Fprintf(w, "  <tr>")
				if genDetail.userTypeID == "100" {
					dbTotalTableCountWhere.columnWhereValue = "100"
					fmt.Fprintf(w, "    <td>"+totalTableCountWhere(dbTotalTableCountWhere)+"</td>")
					dbTotalTableCountWhere.columnWhereValue = "200"
					fmt.Fprintf(w, "    <td>"+totalTableCountWhere(dbTotalTableCountWhere)+"</td>")
					dbTotalTableCountWhere.columnWhereValue = "201"
					fmt.Fprintf(w, "    <td>"+totalTableCountWhere(dbTotalTableCountWhere)+"</td>")
					dbTotalTableCountWhere.columnWhereValue = "300"
					fmt.Fprintf(w, "    <td>"+totalTableCountWhere(dbTotalTableCountWhere)+"</td>")
					dbTotalTableCountWhere.columnWhereValue = "301"
					fmt.Fprintf(w, "    <td>"+totalTableCountWhere(dbTotalTableCountWhere)+"</td>")
					dbTotalTableCountWhere.columnWhereValue = "302"
					fmt.Fprintf(w, "    <td>"+totalTableCountWhere(dbTotalTableCountWhere)+"</td>")
					dbTotalTableCountWhere.columnWhereValue = "400"
					fmt.Fprintf(w, "    <td>"+totalTableCountWhere(dbTotalTableCountWhere)+"</td>")
				} else if genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
					dbTotalTableCountWhere.columnWhereAnd = "customer_id"
					dbTotalTableCountWhere.columnWhereValueAnd = customerID
					if genDetail.userTypeID == "200" {
						dbTotalTableCountWhere.columnWhereValue = "200"
						fmt.Fprintf(w, "    <td>"+totalTableCountWhereAnd(dbTotalTableCountWhere)+"</td>")
						dbTotalTableCountWhere.columnWhereValue = "201"
						fmt.Fprintf(w, "    <td>"+totalTableCountWhereAnd(dbTotalTableCountWhere)+"</td>")
					} else if genDetail.userTypeID == "201" {
						dbTotalTableCountWhere.columnWhereValue = "201"
						fmt.Fprintf(w, "    <td>"+totalTableCountWhereAnd(dbTotalTableCountWhere)+"</td>")
					}
					dbTotalTableCountWhere.columnWhereValue = "300"
					fmt.Fprintf(w, "    <td>"+totalTableCountWhereAnd(dbTotalTableCountWhere)+"</td>")
					dbTotalTableCountWhere.columnWhereValue = "301"
					fmt.Fprintf(w, "    <td>"+totalTableCountWhereAnd(dbTotalTableCountWhere)+"</td>")
					dbTotalTableCountWhere.columnWhereValue = "302"
					fmt.Fprintf(w, "    <td>"+totalTableCountWhereAnd(dbTotalTableCountWhere)+"</td>")
					if genDetail.userTypeID == "200" {
						dbTotalTableCountWhere.columnWhereValue = "400"
						fmt.Fprintf(w, "    <td>"+totalTableCountWhereAnd(dbTotalTableCountWhere)+"</td>")
					}
				} else if genDetail.userTypeID == "300" || genDetail.userTypeID == "301" {
					dbTotalTableCountWhere.columnWhereAnd = "pbx_id"
					dbTotalTableCountWhere.columnWhereValueAnd = pbxID
					if genDetail.userTypeID == "300" {
						dbTotalTableCountWhere.columnWhereValue = "300"
						fmt.Fprintf(w, "    <td>"+totalTableCountWhereAnd(dbTotalTableCountWhere)+"</td>")
						dbTotalTableCountWhere.columnWhereValue = "301"
						fmt.Fprintf(w, "    <td>"+totalTableCountWhereAnd(dbTotalTableCountWhere)+"</td>")
					} else if genDetail.userTypeID == "301" {
						dbTotalTableCountWhere.columnWhereValue = "301"
						fmt.Fprintf(w, "    <td>"+totalTableCountWhereAnd(dbTotalTableCountWhere)+"</td>")
					}
					dbTotalTableCountWhere.columnWhereValue = "302"
					fmt.Fprintf(w, "    <td>"+totalTableCountWhereAnd(dbTotalTableCountWhere)+"</td>")
				}
				fmt.Fprintf(w, "  </tr>")
				fmt.Fprintf(w, "</table>")
				fmt.Fprintf(w, "<br>")

			}
			fmt.Fprintf(w, "<table id=\"table\" class=\"table-user-account\">")
			fmt.Fprintf(w, "  <tr>")
			fmt.Fprintf(w, "    <th class=\"table-title\";>Own User Account Details:</th>")
			fmt.Fprintf(w, "  </tr>")
			fmt.Fprintf(w, "  <tr>")
			fmt.Fprintf(w, "    <th>")
			fmt.Fprintf(w, "      <table id=\"table\" class=\"table-user-account\">")
			fmt.Fprintf(w, "        <tr>")
			fmt.Fprintf(w, "          <th>Account ID</th>")
			fmt.Fprintf(w, "          <th>Name</th>")
			fmt.Fprintf(w, "          <th>Email</th>")
			fmt.Fprintf(w, "          <th>Account Type</th>")
			fmt.Fprintf(w, "          <th>Account Created</th>")
			fmt.Fprintf(w, "        </tr>")
			fmt.Fprintf(w, "        <tr>")
			fmt.Fprintf(w, "          <td>"+userAccountID+"</td>")
			fmt.Fprintf(w, "          <td>"+userAccountFirstName+"<br>"+userAccountLastName+"</td>")
			fmt.Fprintf(w, "          <td>"+userAccountEmail+"</td>")
			fmt.Fprintf(w, "          <td>"+userAccountType+"</td>")
			fmt.Fprintf(w, "          <td>"+formatDateTime(userAccountDateTimeAdded)+"</td>")
			fmt.Fprintf(w, "        </tr>")
			fmt.Fprintf(w, "      </table>")
			fmt.Fprintf(w, "    </th>")
			fmt.Fprintf(w, "  </tr>")
			if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" || genDetail.userTypeID == "300" || genDetail.userTypeID == "301" {
				fmt.Fprintf(w, "  <tr>")
				fmt.Fprintf(w, "    <th><button onclick=\"toggleOtherAccount() \"class=\"button-general button-user-account\">&nbsp Show/Hide Other Accounts &nbsp</button></th>")
				fmt.Fprintf(w, "  </tr>")
			}
			fmt.Fprintf(w, "</table>")
		}

		if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" || genDetail.userTypeID == "300" || genDetail.userTypeID == "301" {

			userCustomerID := userAccountData(dbDetail, "customer_id")
			userCustomerName := userAccountData(dbDetail, "customer_name")
			userPBXID := userAccountData(dbDetail, "pbx_id")
			userPBXName := userAccountData(dbDetail, "pbx_name")

			fmt.Fprintf(w, "<div id=\"other-account-div\" style=\"display:none\">")
			fmt.Fprintf(w, "<br>")
			fmt.Fprintf(w, "<table id=\"table\" class=\"table-user-account\">")
			fmt.Fprintf(w, "  <tr>")
			if genDetail.userTypeID == "100" {
				fmt.Fprintf(w, "    <th class=\"table-title\";>All User Account Details on the YAP Server:</th>")
			} else if genDetail.userTypeID == "200" {
				fmt.Fprintf(w, "    <th class=\"table-title\";>User Account Details for the Customer<br>"+userCustomerName+"<br>(Customer ID: "+userCustomerID+")</th>")
			} else if genDetail.userTypeID == "201" {
				fmt.Fprintf(w, "    <th class=\"table-title\";>PBX User Account Details for the Customer<br>"+userCustomerName+"<br>(Customer ID: "+userCustomerID+")</th>")
			} else if genDetail.userTypeID == "300" || genDetail.userTypeID == "301" {
				fmt.Fprintf(w, "    <th class=\"table-title\";>PBX User Account Details Within the PBX<br>"+userPBXName+"<br>(PBX ID: "+userPBXID+")</th>")
			}
			fmt.Fprintf(w, "  </tr>")
			fmt.Fprintf(w, "  <tr>")
			fmt.Fprintf(w, "    <th>")
			fmt.Fprintf(w, "    <br>")
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			var inputTableHTMLArgument jsFunctionParameter
			inputTableHTMLArgument.inputID = "other-account-input-id"
			inputTableHTMLArgument.funcNameJS = "otherAccountSearchID"
			inputTableHTMLArgument.placeholder = "ID"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "other-account-input-name"
			inputTableHTMLArgument.funcNameJS = "otherAccountSearchName"
			inputTableHTMLArgument.placeholder = "Name"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "other-account-input-email"
			inputTableHTMLArgument.funcNameJS = "otherAccountSearchEmail"
			inputTableHTMLArgument.placeholder = "Email"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "other-account-input-type"
			inputTableHTMLArgument.funcNameJS = "otherAccountSearchType"
			inputTableHTMLArgument.placeholder = "Account Type"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			fmt.Fprintf(w, "    <br>")
			fmt.Fprintf(w, "    <br>")
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "other-account-input-date-time"
			inputTableHTMLArgument.funcNameJS = "otherAccountSearchDateTime"
			inputTableHTMLArgument.placeholder = "Date & Time Created"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
				inputTableHTMLArgument.inputID = "other-account-input-pbx-id"
				inputTableHTMLArgument.funcNameJS = "otherAccountSearchPBXID"
				inputTableHTMLArgument.placeholder = "PBX ID"
				inputTableHTML(w, inputTableHTMLArgument)
				fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
				inputTableHTMLArgument.inputID = "other-account-input-pbx-name"
				inputTableHTMLArgument.funcNameJS = "otherAccountSearchPBXName"
				inputTableHTMLArgument.placeholder = "PBX Name"
				inputTableHTML(w, inputTableHTMLArgument)
			}
			if genDetail.userTypeID == "100" {
				fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
				inputTableHTMLArgument.inputID = "other-account-input-customer-id"
				inputTableHTMLArgument.funcNameJS = "otherAccountSearchCustomerID"
				inputTableHTMLArgument.placeholder = "Customer ID"
				inputTableHTML(w, inputTableHTMLArgument)
				fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
				fmt.Fprintf(w, "    <br>")
				fmt.Fprintf(w, "    <br>")
				inputTableHTMLArgument.inputID = "other-account-input-customer-name"
				inputTableHTMLArgument.funcNameJS = "otherAccountSearchCustomerName"
				inputTableHTMLArgument.placeholder = "Customer Name"
				inputTableHTML(w, inputTableHTMLArgument)
			}
			fmt.Fprintf(w, "    <br>")
			fmt.Fprintf(w, "    <br>")
			fmt.Fprintf(w, "    </th>")
			fmt.Fprintf(w, "  </tr>")
			fmt.Fprintf(w, "  <tr>")
			fmt.Fprintf(w, "    <th>")
			var exportCSVButtonHTMLArgument jsFunctionParameter
			exportCSVButtonHTMLArgument.funcNameJS = "OtherAccount"
			exportCSVButtonHTMLArgument.buttonCSS = "button-user-account"
			exportCSVButtonHTML(w, exportCSVButtonHTMLArgument)
			fmt.Fprintf(w, "    </th>")
			fmt.Fprintf(w, "  </tr>")
			fmt.Fprintf(w, "  <tr>")
			fmt.Fprintf(w, "    <th>")
			fmt.Fprintf(w, "      <table id=\"other-account-table\" class=\"table-user-account\">")
			fmt.Fprintf(w, "        <tr>")
			fmt.Fprintf(w, "          <th>Account ID</th>")
			fmt.Fprintf(w, "          <th>Name</th>")
			fmt.Fprintf(w, "          <th>Email</th>")
			fmt.Fprintf(w, "          <th>Account Type</th>")
			fmt.Fprintf(w, "          <th>Account Created</th>")
			if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
				fmt.Fprintf(w, "          <th>PBX ID</th>")
				fmt.Fprintf(w, "          <th>PBX Name</th>")
			}
			if genDetail.userTypeID == "100" {
				fmt.Fprintf(w, "          <th>Customer ID</th>")
				fmt.Fprintf(w, "          <th>Customer Name</th>")
			}
			fmt.Fprintf(w, "        </tr>")

			var whereClause string

			if genDetail.userTypeID == "100" {
				whereClause = "WHERE customer_id != ? AND pbx_id != ?;"
				userCustomerID = "0"
				userPBXID = "0"
			} else if genDetail.userTypeID == "200" {
				whereClause = "WHERE customer_id = ? AND pbx_id != ?;"
				userPBXID = "0"
			} else if genDetail.userTypeID == "201" {
				whereClause = "WHERE customer_id = ? AND pbx_id != ? AND user_account_type_id != 200;"
				userPBXID = "0"
			} else if genDetail.userTypeID == "300" {
				whereClause = "WHERE customer_id = ? AND pbx_id = ?;"
			} else if genDetail.userTypeID == "301" {
				whereClause = "WHERE customer_id = ? AND pbx_id = ? AND user_account_type_id != 300;"
			}

			otherUserAccountSQL, err := dbDetail.connection.Query(`SELECT
									 user_account_id,
						     			 user_account_first_name,
						     			 user_account_last_name,  
						     			 user_account_email,                                                   
						     			 user_account_type,  
						     			 user_account_type_id,
						     			 user_account_date_time_added, 
						     			 customer_name,
						     			 customer_id,
						     			 pbx_name,
						     			 pbx_id						     
								       FROM
								         yap.view___account_detail
								       `+whereClause, userCustomerID, userPBXID)
			// Error
			if err != nil {
				panic(err)
			}

			for otherUserAccountSQL.Next() {

				err = otherUserAccountSQL.Scan(
					&userAccountID,
					&userAccountFirstName,
					&userAccountLastName,
					&userAccountEmail,
					&userAccountType,
					&userAccountTypeID,
					&userAccountDateTimeAdded,
					&customerName,
					&customerID,
					&pbxName,
					&pbxID,
				)

				// Error
				if err != nil {
					panic(err)
				}

				fmt.Fprintf(w, "        <tr>")
				userAccountListTdHTML(w, userAccountID, userAccountTypeID, userAccountID)
				userAccountListTdHTML(w, userAccountID, userAccountTypeID, userAccountFirstName+" "+userAccountLastName)
				userAccountListTdHTML(w, userAccountID, userAccountTypeID, userAccountEmail)
				userAccountListTdHTML(w, userAccountID, userAccountTypeID, userAccountType)
				userAccountListTdHTML(w, userAccountID, userAccountTypeID, formatDateTime(userAccountDateTimeAdded))

				if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
					if pbxID != "1" {
						userAccountListTdHTML(w, userAccountID, userAccountTypeID, pbxID)
					} else {
						userAccountListTdHTML(w, userAccountID, userAccountTypeID, "-")
					}
					if pbxName != "system" {
						userAccountListTdHTML(w, userAccountID, userAccountTypeID, pbxName)
					} else {
						userAccountListTdHTML(w, userAccountID, userAccountTypeID, "-")
					}
				}
				if genDetail.userTypeID == "100" {
					if customerID != "1" {
						userAccountListTdHTML(w, userAccountID, userAccountTypeID, customerID)
					} else {
						userAccountListTdHTML(w, userAccountID, userAccountTypeID, "-")
					}
					if customerName != "system" {
						userAccountListTdHTML(w, userAccountID, userAccountTypeID, customerName)
					} else {
						userAccountListTdHTML(w, userAccountID, userAccountTypeID, "-")
					}
				}
				fmt.Fprintf(w, "        </tr>")
			}

			fmt.Fprintf(w, "      </table>")

			var filterTableJSArgument jsFunctionParameter
			filterTableJSArgument.tableID = "other-account-table"
			// JS filter function for account ID in the other account table
			filterTableJSArgument.funcNameJS = "otherAccountSearchID"
			filterTableJSArgument.inputID = "other-account-input-id"
			filterTableJSArgument.columnNumber = 0
			filterTableJS(w, filterTableJSArgument)
			// JS filter function for name in the other account table
			filterTableJSArgument.funcNameJS = "otherAccountSearchName"
			filterTableJSArgument.inputID = "other-account-input-name"
			filterTableJSArgument.columnNumber = 1
			filterTableJS(w, filterTableJSArgument)
			// JS filter function for email in the other account table
			filterTableJSArgument.funcNameJS = "otherAccountSearchEmail"
			filterTableJSArgument.inputID = "other-account-input-email"
			filterTableJSArgument.columnNumber = 2
			filterTableJS(w, filterTableJSArgument)
			// JS filter function for type in the other account table
			filterTableJSArgument.funcNameJS = "otherAccountSearchType"
			filterTableJSArgument.inputID = "other-account-input-type"
			filterTableJSArgument.columnNumber = 3
			filterTableJS(w, filterTableJSArgument)
			// JS filter function for date and time in the other account table
			filterTableJSArgument.funcNameJS = "otherAccountSearchDateTime"
			filterTableJSArgument.inputID = "other-account-input-date-time"
			filterTableJSArgument.columnNumber = 4
			filterTableJS(w, filterTableJSArgument)
			if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
				// JS filter function for PBX ID in the other account table
				filterTableJSArgument.funcNameJS = "otherAccountSearchPBXID"
				filterTableJSArgument.inputID = "other-account-input-pbx-id"
				filterTableJSArgument.columnNumber = 5
				filterTableJS(w, filterTableJSArgument)
				// JS filter function for PBX name in the other account table
				filterTableJSArgument.funcNameJS = "otherAccountSearchPBXName"
				filterTableJSArgument.inputID = "other-account-input-pbx-name"
				filterTableJSArgument.columnNumber = 6
				filterTableJS(w, filterTableJSArgument)
			}
			if genDetail.userTypeID == "100" {
				// JS filter function for the customer ID in the other account table
				filterTableJSArgument.funcNameJS = "otherAccountSearchCustomerID"
				filterTableJSArgument.inputID = "other-account-input-customer-id"
				filterTableJSArgument.columnNumber = 7
				filterTableJS(w, filterTableJSArgument)
				// JS filter function for the customer name in the other account table
				filterTableJSArgument.funcNameJS = "otherAccountSearchCustomerName"
				filterTableJSArgument.inputID = "other-account-input-customer-name"
				filterTableJSArgument.columnNumber = 8
				filterTableJS(w, filterTableJSArgument)
			}
			var exportCSVJSArgument jsFunctionParameter
			exportCSVJSArgument.funcNameJS = "OtherAccount"
			exportCSVJSArgument.tableID = "other-account-table"
			exportCSVJSArgument.fileName = "YAP_user_account_details"
			exportCSVJSArgument.pathURL = "user-account"
			exportCSVJS(w, exportCSVJSArgument)
			fmt.Fprintf(w, "    </th>")
			fmt.Fprintf(w, "  </tr>")
			fmt.Fprintf(w, "</table>")
			fmt.Fprintf(w, "</div>")
			var toggleDivJSArgument jsFunctionParameter
			toggleDivJSArgument.funcNameJS = "toggleOtherAccount"
			toggleDivJSArgument.divID = "other-account-div"
			toggleDivJS(w, toggleDivJSArgument)
		}
	} else {
		panic("userAccountList function should only be called with account type ID 100, 200, 201, 300, 301, 302, 400")
	}
}

// Function to add new user account
func userAccountAdd(w http.ResponseWriter, r *http.Request, dbDetail databaseFunctionParameter, genDetail generalFunctionParameter) {

	// Only account type ID 100 should be able to use this function
	if genDetail.userTypeID == "100" {

		fmt.Fprintf(w, "<form method=\"POST\" action=\"/user-account\">")
		fmt.Fprintf(w, "<table class=\"table-add\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>Add a New User Account</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_account_input_first_name", "First Name (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_account_input_last_name", "Last Name (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_account_input_email", "Email Address (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		userAccountTypeIDNameList, _ := userAccountTypeSlice(dbDetail)
		selectDoubleHTML(w, "add_account_select_account_type", "Account Type", userAccountTypeIDNameList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		pbxIDNameList, _ := pbxSlice(dbDetail)
		selectDoubleHTML(w, "add_account_select_pbx_id", "PBX", pbxIDNameList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		customerIDNameList, _ := customerSlice(dbDetail)
		selectDoubleHTML(w, "add_account_select_customer_id", "Customer", customerIDNameList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><input type=\"submit\" value=\"Create Account\"></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</form>")

		addAccountInputFirstName := r.FormValue("add_account_input_first_name")
		addAccountInputLastName := r.FormValue("add_account_input_last_name")
		addAccountInputEmail := r.FormValue("add_account_input_email")
		addAccountSelectAccountType := r.FormValue("add_account_select_account_type")
		addAccountSelectPBXID := r.FormValue("add_account_select_pbx_id")
		addAccountSelectCustomerID := r.FormValue("add_account_select_customer_id")

		// Validate the first name string
		validateFirstName := validateInput(addAccountInputFirstName, "alphaNum")

		// Validate the last name string
		validateLastName := validateInput(addAccountInputLastName, "alphaNum")

		// Validate the email address string
		validateEmail := validateInput(addAccountInputEmail, "email")

		// Check user type ID is contained in the slice
		_, userAccountTypeIDList := userAccountTypeSlice(dbDetail)
		validateUserAccountTypeID := slices.Contains(userAccountTypeIDList, addAccountSelectAccountType)

		// Check PBX ID is contained in the slice
		_, pbxIDList := pbxSlice(dbDetail)
		pbxIDList = append(pbxIDList, "")
		validatePBXID := slices.Contains(pbxIDList, addAccountSelectPBXID)

		// Check customer ID is contained in the slice
		_, customerIDList := customerSlice(dbDetail)
		customerIDList = append(customerIDList, "")
		validateCustomerID := slices.Contains(customerIDList, addAccountSelectCustomerID)

		if addAccountInputFirstName == "" && addAccountInputLastName == "" && addAccountInputEmail == "" && addAccountSelectAccountType == "" && addAccountSelectPBXID == "" && addAccountSelectCustomerID == "" {
			// Do Nothing
		} else if validateFirstName == false {
			messageHTML(w, validationMessageAccountFirstName, "warning")
		} else if validateLastName == false {
			messageHTML(w, validationMessageAccountLastName, "warning")
		} else if validateEmail == false {
			messageHTML(w, validationMessageAccountEmail, "warning")
		} else if validateUserAccountTypeID == false {
			messageHTML(w, validationMessageAccountType, "warning")
		} else if validatePBXID == false {
			messageHTML(w, validationMessagePBX, "warning")
		} else if validateCustomerID == false {
			messageHTML(w, validationMessageCustomer, "warning")
		} else if addAccountSelectAccountType == "100" && genDetail.userID != "1" {
			messageHTML(w, validationMessageAccountYAPAdmin, "warning")
		} else if addAccountSelectAccountType == "200" && addAccountSelectCustomerID == "" {
			messageHTML(w, validationMessageCustomer, "warning")
		} else if addAccountSelectAccountType == "201" && addAccountSelectCustomerID == "" {
			messageHTML(w, validationMessageCustomer, "warning")
		} else if addAccountSelectAccountType == "400" && addAccountSelectCustomerID == "" {
			messageHTML(w, validationMessageCustomer, "warning")
		} else if addAccountSelectAccountType == "300" && addAccountSelectPBXID == "" {
			messageHTML(w, validationMessagePBX, "warning")
		} else if addAccountSelectAccountType == "301" && addAccountSelectPBXID == "" {
			messageHTML(w, validationMessagePBX, "warning")
		} else if addAccountSelectAccountType == "302" && addAccountSelectPBXID == "" {
			messageHTML(w, validationMessagePBX, "warning")
		} else if validateFirstName == true && validateLastName == true && validateEmail == true && validateUserAccountTypeID == true && validatePBXID == true && validateCustomerID == true {
			dbDetail.table = "view___account_detail"
			dbDetail.column = "user_account_email"
			dbDetail.columnWhere = "user_account_email"
			dbDetail.columnWhereValue = addAccountInputEmail

			checkAccountEmailExist := selectWhere(dbDetail)

			if checkAccountEmailExist == addAccountInputEmail {
				messageHTML(w, validationMessageEmailAlreadyExist, "warning")
			} else {
				// The database is designed not to allow two or more records with the same email address
				// The conditional statements are mostly used to inform the user with messages in HTML
				if addAccountSelectAccountType == "100" {
					addAccountSelectCustomerID = "1"
					addAccountSelectPBXID = "1"
					checkAccountType100Created := selectWhere(dbDetail)
					if checkAccountType100Created == addAccountInputEmail {
						messageHTML(w, validationMessageEmailAlreadyExist, "warning")
					} else {
						messageHTML(w, validationMessageAccountCreated, "success")
					}
				} else if addAccountSelectAccountType == "200" || addAccountSelectAccountType == "201" || addAccountSelectAccountType == "400" {
					addAccountSelectPBXID = "1"
					if addAccountSelectAccountType == "200" {
						checkAccountType200Created := selectWhere(dbDetail)
						if checkAccountType200Created == addAccountInputEmail {
							messageHTML(w, validationMessageEmailAlreadyExist, "warning")
						} else {
							messageHTML(w, validationMessageAccountCreated, "success")
						}
					} else if addAccountSelectAccountType == "201" {
						checkAccountType201Created := selectWhere(dbDetail)
						if checkAccountType201Created == addAccountInputEmail {
							messageHTML(w, validationMessageEmailAlreadyExist, "warning")
						} else {
							messageHTML(w, validationMessageAccountCreated, "success")
						}
					} else if addAccountSelectAccountType == "400" {
						checkAccountType400Created := selectWhere(dbDetail)
						if checkAccountType400Created == addAccountInputEmail {
							messageHTML(w, validationMessageEmailAlreadyExist, "warning")
						} else {
							messageHTML(w, validationMessageAccountCreated, "success")
						}
					}
				} else if addAccountSelectAccountType == "300" || addAccountSelectAccountType == "301" || addAccountSelectAccountType == "302" {
					if addAccountSelectAccountType == "300" {
						checkAccountType300Created := selectWhere(dbDetail)
						if checkAccountType300Created == addAccountInputEmail {
							messageHTML(w, validationMessageEmailAlreadyExist, "warning")
						} else {
							messageHTML(w, validationMessageAccountCreated, "success")
						}
					} else if addAccountSelectAccountType == "301" {
						checkAccountType301Created := selectWhere(dbDetail)
						if checkAccountType301Created == addAccountInputEmail {
							messageHTML(w, validationMessageEmailAlreadyExist, "warning")
						} else {
							messageHTML(w, validationMessageAccountCreated, "success")
						}
					} else if addAccountSelectAccountType == "302" {
						checkAccountType302Created := selectWhere(dbDetail)
						if checkAccountType302Created == addAccountInputEmail {
							messageHTML(w, validationMessageEmailAlreadyExist, "warning")
						} else {
							messageHTML(w, validationMessageAccountCreated, "success")
						}
					}
					dbDetail.table = "view___pbx_detail"
					dbDetail.column = "customer_id"
					dbDetail.columnWhere = "pbx_id"
					dbDetail.columnWhereValue = addAccountSelectPBXID

					addAccountSelectCustomerID = selectWhere(dbDetail)
				}

				dbDetail.connection.Query(`INSERT 
        	                   INTO
	       		       user_account (
			           email,
			           first_name,
			           last_name,
			           user_account_type_id,
			           customer_id,
			           pbx_id
			       )
			       VALUES(?, ?, ?, ?, ?, ?);`,
					addAccountInputEmail,
					nullSQL(addAccountInputFirstName),
					nullSQL(addAccountInputLastName),
					addAccountSelectAccountType,
					addAccountSelectCustomerID,
					addAccountSelectPBXID)
			}
		}
	} else {
		panic("userAccountAdd function should only be called with account type ID 100")
	}
}

// User account edit function
func userAccountEdit(w http.ResponseWriter, r *http.Request, dbDetail databaseFunctionParameter, genDetail generalFunctionParameter) {

	// Only account type ID 100 should be able to use this function
	if genDetail.userTypeID == "100" {

		// List of first name and last name column names from the user account table
		accountColumnList := [][]string{
			{"first_name", "First Name"},
			{"last_name", "Last Name"},
		}

		fmt.Fprintf(w, "<form method=\"POST\" action=\"/user-account\">")
		fmt.Fprintf(w, "<table class=\"table-user-account\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>Edit User Account Details</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td style=\"text-align: left;\">")
		fmt.Fprintf(w, "      <b>Acceptable Values for Columns</b><br><br>")
		fmt.Fprintf(w, "      <b>First Name:</b> text<br>")
		fmt.Fprintf(w, "      <b>Last Name:</b> text<br>")
		fmt.Fprintf(w, "    </td>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "edit_account_input_account_id", "Account ID (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "edit_account_input_email", "Account Email (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		selectDoubleHiddenHTML(w, "edit_account_select_column", "Column to Edit (Cannot Be Empty)", accountColumnList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "edit_account_input_new_value", "New Value")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><input type=\"submit\" value=\"Update Account\"></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</form>")

		editAccountInputAccountID := r.FormValue("edit_account_input_account_id")
		editAccountInputEmail := r.FormValue("edit_account_input_email")
		editAccountSelectColumn := r.FormValue("edit_account_select_column")
		editAccountInputNewValue := r.FormValue("edit_account_input_new_value")

		dbDetail.table = "view___account_detail"

		// Validate account ID from account ID List
		dbDetail.column = "user_account_id"
		accountIDList := singleColumnSlice(dbDetail)
		validateAccountID := slices.Contains(accountIDList, editAccountInputAccountID)

		// Validate email from email List
		dbDetail.column = "user_account_email"
		accountEmailList := singleColumnSlice(dbDetail)
		validateAccountEmail := slices.Contains(accountEmailList, editAccountInputEmail)

		if editAccountInputAccountID == "" && editAccountInputEmail == "" && editAccountSelectColumn == "" && editAccountInputNewValue == "" {
			// Do Nothing
		} else if validateAccountID == false {
			messageHTML(w, validationMessageAccountID, "warning")
		} else if editAccountInputAccountID == "1" {
			messageHTML(w, validationMessageAccountIDOne, "warning")
		} else if validateAccountEmail == false {
			messageHTML(w, validationMessageEmailDoesNotExist, "warning")
		} else if editAccountSelectColumn == "" {
			messageHTML(w, validationMessageAccountColumn, "warning")
		} else if editAccountSelectColumn == "first_name" || editAccountSelectColumn == "last_name" {
			// Validate editAccountInputNewValue is a string
			validateNewValue := validateInput(editAccountInputNewValue, "alphaNumEmpty")
			if validateNewValue == true {
				dbDetail.connection.Query("UPDATE user_account SET "+editAccountSelectColumn+" = ? WHERE id = ? AND email = ?;", editAccountInputNewValue, editAccountInputAccountID, editAccountInputEmail)
			} else {
				messageHTML(w, validationMessageAccountColumn, "warning")
			}
		} else {
			messageHTML(w, validationMessageAccountColumn, "warning")
		}
	} else {
		panic("userAccountEdit function should only be called with account type ID 100")
	}
}

func userAccountDelete(w http.ResponseWriter, r *http.Request, dbDetail databaseFunctionParameter, genDetail generalFunctionParameter) {

	// Only account type ID 100 should be able to use this function
	if genDetail.userTypeID == "100" {

		// Get account type ID and email from the database and append to slice
		var userAccountIDList []string
		var userAccountID string

		var userAccountEmailList []string
		var userAccountEmail string

		userAccountIDEmailSQL, err := dbDetail.connection.Query(`SELECT
	                                                             user_account_id,
	                                                             user_account_email
	                                                         FROM
	                                                             yap.view___account_detail;`)

		// Error
		if err != nil {
			panic(err)
		}

		for userAccountIDEmailSQL.Next() {

			err = userAccountIDEmailSQL.Scan(
				&userAccountID,
				&userAccountEmail,
			)

			// Error
			if err != nil {
				panic(err)
			}

			userAccountIDList = append(userAccountIDList, userAccountID)
			userAccountEmailList = append(userAccountEmailList, userAccountEmail)
		}

		// Delete individual user account code
		fmt.Fprintf(w, "<form method=\"POST\" action=\"/user-account\">")
		fmt.Fprintf(w, "<table class=\"table-delete\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>Delete Individual User Account</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "delete_individual_user_account_input_account_id", "Account ID (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "delete_individual_user_account_input_account_email", "Account Email (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		userAccountTypeIDNameList, _ := userAccountTypeSlice(dbDetail)
		selectDoubleHTML(w, "delete_individual_user_account_select_account_type", "Account Type", userAccountTypeIDNameList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><input type=\"submit\" value=\"Delete Account\"></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</form>")

		deleteIndividualUserAccountInputAccountID := r.FormValue("delete_individual_user_account_input_account_id")
		deleteIndividualUserAccountInputEmail := r.FormValue("delete_individual_user_account_input_account_email")
		deleteIndividualUserAccountSelectAccountType := r.FormValue("delete_individual_user_account_select_account_type")

		// Check user account ID is contained in the slice
		validateIndividualUserAccountID := slices.Contains(userAccountIDList, deleteIndividualUserAccountInputAccountID)

		// Check user email is contained in the slice
		validateIndividualUserAccountEmail := slices.Contains(userAccountEmailList, deleteIndividualUserAccountInputEmail)

		// Check user type ID is contained in the slice
		_, userAccountTypeIDList := userAccountTypeSlice(dbDetail)
		validateIndividualUserAccountTypeID := slices.Contains(userAccountTypeIDList, deleteIndividualUserAccountSelectAccountType)

		if deleteIndividualUserAccountInputAccountID == "" && deleteIndividualUserAccountInputEmail == "" && deleteIndividualUserAccountSelectAccountType == "" {
			// Do nothing
		} else if validateIndividualUserAccountID == false || deleteIndividualUserAccountInputAccountID == "" {
			messageHTML(w, validationMessageAccountID, "warning")
		} else if validateIndividualUserAccountEmail == false || deleteIndividualUserAccountInputEmail == "" {
			messageHTML(w, validationMessageEmailDoesNotExist, "warning")
		} else if validateIndividualUserAccountTypeID == false || deleteIndividualUserAccountSelectAccountType == "" {
			messageHTML(w, validationMessageAccountType, "warning")
		} else if deleteIndividualUserAccountInputAccountID == "1" {
			messageHTML(w, validationMessageAccountIDOne, "warning")
		} else if deleteIndividualUserAccountSelectAccountType == "100" && genDetail.userID != "1" {
			messageHTML(w, validationMessageAccountYAPAdmin, "warning")
		} else if validateIndividualUserAccountID == true && validateIndividualUserAccountEmail == true && validateIndividualUserAccountTypeID == true {
			dbDetail.table = "view___account_detail"
			dbDetail.column = "user_account_id"
			dbDetail.columnWhere = "user_account_id"
			dbDetail.columnWhereValue = deleteIndividualUserAccountInputAccountID

			dbDetail.connection.Query(`DELETE FROM user_account WHERE id = ? AND user_account_type_id = ?;`, deleteIndividualUserAccountInputAccountID, deleteIndividualUserAccountSelectAccountType)

			checkIndividualUserAccountDeleted := selectWhere(dbDetail)

			if checkIndividualUserAccountDeleted == "" {
				messageHTML(w, validationMessageAccountIndividualDeleted, "success")
			} else {
				messageHTML(w, validationMessageAccountType, "warning")
			}
		} else {
			messageHTML(w, validationMessageInvalid, "warning")
		}

		// Delete all user accounts for a PBX code
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<form method=\"POST\" action=\"/user-account\">")
		fmt.Fprintf(w, "<table class=\"table-delete\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>Delete All User Accounts for a PBX</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		pbxIDNameList, _ := pbxSlice(dbDetail)
		selectDoubleHTML(w, "delete_pbx_user_account_select_account_pbx_id", "PBX", pbxIDNameList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		confirmList := yesSlice()
		selectSingleHTML(w, "delete_pbx_user_account_select_confirm", "yes to Confirm (Cannot Be Empty)", confirmList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><input type=\"submit\" value=\"Delete PBX Accounts\"></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</form>")

		deletePBXUserAccountSelectPBXID := r.FormValue("delete_pbx_user_account_select_account_pbx_id")
		deletePBXUserAccountSelectConfirm := r.FormValue("delete_pbx_user_account_select_confirm")

		// Validate PBX List
		_, pbxIDList := pbxSlice(dbDetail)
		validatePBXID := slices.Contains(pbxIDList, deletePBXUserAccountSelectPBXID)

		if deletePBXUserAccountSelectPBXID == "" && deletePBXUserAccountSelectConfirm == "" {
			// Do Nothing
		} else if validatePBXID == false && deletePBXUserAccountSelectConfirm == "yes" {
			messageHTML(w, validationMessagePBX, "warning")
		} else if validatePBXID == true && deletePBXUserAccountSelectConfirm != "yes" {
			messageHTML(w, validationMessageConfirmation, "warning")
		} else if deletePBXUserAccountSelectPBXID == "1" {
			messageHTML(w, validationMessagePBX, "warning")
		} else if validatePBXID == true && deletePBXUserAccountSelectConfirm == "yes" {

			dbDetail.table = "view___account_detail"
			dbDetail.column = "user_account_id"
			dbDetail.columnWhere = "pbx_id"
			dbDetail.columnWhereValue = deletePBXUserAccountSelectPBXID

			checkPBXUserAccountExist := selectWhere(dbDetail)

			if checkPBXUserAccountExist == "" {
				messageHTML(w, validationMessageAccountPBXDoesNotExist, "warning")
			} else {

				dbDetail.connection.Query(`DELETE FROM user_account WHERE pbx_id = ?;`, deletePBXUserAccountSelectPBXID)

				checkPBXUserAccountDeleted := selectWhere(dbDetail)

				if checkPBXUserAccountDeleted == "" {
					messageHTML(w, validationMessageAccountMultipleDeleted, "success")
				} else {
					messageHTML(w, validationMessageAccountMultipleNotDeleted, "warning")
				}
			}
		} else {
			messageHTML(w, validationMessageInvalid, "warning")
		}

		// Delete all user accounts for a Customer
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<form method=\"POST\" action=\"/user-account\">")
		fmt.Fprintf(w, "<table class=\"table-delete\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>Delete All User Accounts for a Customer</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		customerIDNameList, _ := customerSlice(dbDetail)
		selectDoubleHTML(w, "delete_customer_user_account_select_account_customer_id", "Customer", customerIDNameList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		selectSingleHTML(w, "delete_customer_user_account_select_confirm", "yes to Confirm (Cannot Be Empty)", confirmList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><input type=\"submit\" value=\"Delete Customer Accounts\"></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</form>")

		deleteCustomerUserAccountSelectCustomerID := r.FormValue("delete_customer_user_account_select_account_customer_id")
		deleteCustomerUserAccountSelectConfirm := r.FormValue("delete_customer_user_account_select_confirm")

		// Validate Customer List
		_, customerIDList := customerSlice(dbDetail)
		validateCustomerID := slices.Contains(customerIDList, deleteCustomerUserAccountSelectCustomerID)

		if deleteCustomerUserAccountSelectCustomerID == "" && deleteCustomerUserAccountSelectConfirm == "" {
			// Do Nothing
		} else if validateCustomerID == false && deleteCustomerUserAccountSelectConfirm == "yes" {
			messageHTML(w, validationMessageCustomer, "warning")
		} else if validateCustomerID == true && deleteCustomerUserAccountSelectConfirm != "yes" {
			messageHTML(w, validationMessageConfirmation, "warning")
		} else if deleteCustomerUserAccountSelectCustomerID == "1" {
			messageHTML(w, validationMessageAccountIDOne, "warning")
		} else if validateCustomerID == true && deleteCustomerUserAccountSelectConfirm == "yes" {

			dbDetail.table = "view___account_detail"
			dbDetail.column = "user_account_id"
			dbDetail.columnWhere = "customer_id"
			dbDetail.columnWhereValue = deleteCustomerUserAccountSelectCustomerID

			checkCustomerUserAccountExist := selectWhere(dbDetail)

			if checkCustomerUserAccountExist == "" {
				messageHTML(w, validationMessageAccountCustomerDoesNotExist, "warning")
			} else {

				dbDetail.connection.Query(`DELETE FROM user_account WHERE customer_id = ?;`, deleteCustomerUserAccountSelectCustomerID)

				checkCustomerUserAccountDeleted := selectWhere(dbDetail)

				if checkCustomerUserAccountDeleted == "" {
					messageHTML(w, validationMessageAccountMultipleDeleted, "success")
				} else {
					messageHTML(w, validationMessageAccountMultipleNotDeleted, "warning")
				}
			}

		} else {
			messageHTML(w, validationMessageInvalid, "warning")
		}
	} else {
		panic("userAccountDelete function should only be called with account type ID 100")
	}
}

//----------------------------------------------------------------------------------------------------

// Customer page functions
func customerList(w http.ResponseWriter, dbDetail databaseFunctionParameter, genDetail generalFunctionParameter) {

	// Only account type ID 100, 200, 201, 400 should be able to use this function
	if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" || genDetail.userTypeID == "400" {

		var (
			customerName                     string
			customerID                       string
			customerDateTimeAdded            string
			customerUKBased                  string
			customerConsumerType             string
			customerUKVATRegistered          string
			customerUKVATNumber              string
			customerResellingMinutes         string
			customerPBXLimit                 string
			customerPBXSalesTaxRate          string
			customerPBXSalesTaxStatus        string
			customerExtSalesTaxRate          string
			customerExtSalesTaxStatus        string
			customerPBXSetupPrice            string
			customerPBXRentalPrice           string
			customerPBXCeasePrice            string
			customerPBXContractLength        string
			customerExtSetupPrice            string
			customerExtRentalPrice           string
			customerExtCeasePrice            string
			customerExtContractLength        string
			customerSiteAddressLine1         string
			customerSiteAddressLine2         string
			customerSiteCityTownVillage      string
			customerSiteCountyStateRegion    string
			customerSitePostcodeZipCode      string
			customerSiteCountry              string
			customerSiteContactEmail         string
			customerSiteContactNumber        string
			customerInvoiceAddressLine1      string
			customerInvoiceAddressLine2      string
			customerInvoiceCityTownVillage   string
			customerInvoiceCountyStateRegion string
			customerInvoicePostcodeZipCode   string
			customerInvoiceCountry           string
			customerInvoiceContactEmail      string
			customerInvoiceContactNumber     string
		)

		var dbTableCountUserCustomer databaseFunctionParameter
		dbTableCountUserCustomer.connection = dbDetail.connection
		dbTableCountUserCustomer.database = dbDetail.database
		dbTableCountUserCustomer.table = "view___customer_detail"

		if genDetail.userTypeID == "100" {
			fmt.Fprintf(w, "<table id=\"table\" class=\"table-customer\">")
			fmt.Fprintf(w, "  <tr>")
			fmt.Fprintf(w, "    <th>")
			fmt.Fprintf(w, "      <table id=\"table\" class=\"table-customer\">")
			fmt.Fprintf(w, "        <tr>")
			fmt.Fprintf(w, "          <th>Total Customers on the YAP Server</th>")
			fmt.Fprintf(w, "        </tr>")
			fmt.Fprintf(w, "        <tr>")
			dbTableCountUserCustomer.countMinusOne = true
			fmt.Fprintf(w, "          <td>"+totalTableCount(dbTableCountUserCustomer)+"</td>")
			fmt.Fprintf(w, "        </tr>")
			fmt.Fprintf(w, "      </table>")
			fmt.Fprintf(w, "    </th>")
			fmt.Fprintf(w, "  </tr>")
			fmt.Fprintf(w, "  <tr>")
			fmt.Fprintf(w, "    <th><button onclick=\"toggleCustomer() \"class=\"button-general button-customer\">&nbsp Show/Hide Customers &nbsp</button></th>")
			fmt.Fprintf(w, "  </tr>")
			fmt.Fprintf(w, "</table>")
		}

		if genDetail.userTypeID == "100" {
			fmt.Fprintf(w, "<div id=\"customer-div\" style=\"display:none\">")
			fmt.Fprintf(w, "<br>")
		} else {
			fmt.Fprintf(w, "<div id=\"customer-div\">")
		}
		fmt.Fprintf(w, "<table id=\"table\" class=\"table-customer\">")
		fmt.Fprintf(w, "  <tr>")
		if genDetail.userTypeID == "100" {
			fmt.Fprintf(w, "    <th class=\"table-title\";>All Customer Contact Details on the YAP Server:</th>")
		} else {
			fmt.Fprintf(w, "    <th class=\"table-title\";>Customer Contact Details</th>")
		}
		fmt.Fprintf(w, "  </tr>")
		if genDetail.userTypeID == "100" {
			fmt.Fprintf(w, "  <tr>")
			fmt.Fprintf(w, "    <th>")
			fmt.Fprintf(w, "    <br>")
			var inputTableHTMLArgument jsFunctionParameter
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "customer-contact-input-customer-id"
			inputTableHTMLArgument.funcNameJS = "customerContactSearchCustomerID"
			inputTableHTMLArgument.placeholder = "Customer ID"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "customer-contact-input-customer-name"
			inputTableHTMLArgument.funcNameJS = "customerContactSearchCustomerName"
			inputTableHTMLArgument.placeholder = "Customer Name"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "customer-contact-input-site-address"
			inputTableHTMLArgument.funcNameJS = "customerContactSearchSiteAddress"
			inputTableHTMLArgument.placeholder = "Customer Site Address"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "customer-contact-input-site-email"
			inputTableHTMLArgument.funcNameJS = "customerContactSearchSiteEmail"
			inputTableHTMLArgument.placeholder = "Customer Site Email Address"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			fmt.Fprintf(w, "    <br>")
			fmt.Fprintf(w, "    <br>")
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "customer-contact-input-site-phone"
			inputTableHTMLArgument.funcNameJS = "customerContactSearchSitePhone"
			inputTableHTMLArgument.placeholder = "Customer Site Phone Number"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "customer-contact-input-invoice-address"
			inputTableHTMLArgument.funcNameJS = "customerContactSearchInvoiceAddress"
			inputTableHTMLArgument.placeholder = "Customer Invoice Address"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "customer-contact-input-invoice-email"
			inputTableHTMLArgument.funcNameJS = "customerContactSearchInvoiceEmail"
			inputTableHTMLArgument.placeholder = "Customer Invoice Email Address"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "customer-contact-input-invoice-phone"
			inputTableHTMLArgument.funcNameJS = "customerContactSearchInvoicePhone"
			inputTableHTMLArgument.placeholder = "Customer Invoice Phone Number"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			fmt.Fprintf(w, "    <br>")
			fmt.Fprintf(w, "    <br>")
			fmt.Fprintf(w, "    </th>")
			fmt.Fprintf(w, "  </tr>")
		}
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		var exportCSVButtonHTMLArgument jsFunctionParameter
		exportCSVButtonHTMLArgument.funcNameJS = "CustomerContact"
		exportCSVButtonHTMLArgument.buttonCSS = "button-customer"
		exportCSVButtonHTML(w, exportCSVButtonHTMLArgument)
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table id=\"customer-contact-table\" class=\"table-customer\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <th>Customer ID</th>")
		fmt.Fprintf(w, "          <th>Customer Name</th>")
		fmt.Fprintf(w, "          <th>Site Address</th>")
		fmt.Fprintf(w, "          <th>Site Email Address</th>")
		fmt.Fprintf(w, "          <th>Site Phone Number</th>")
		fmt.Fprintf(w, "          <th>Invoice Address</th>")
		fmt.Fprintf(w, "          <th>Invoice Email Address</th>")
		fmt.Fprintf(w, "          <th>Invoice Phone Number</th>")
		fmt.Fprintf(w, "        </tr>")

		var whereClause string

		if genDetail.userTypeID == "100" {
			whereClause = "WHERE customer_id != ?;"
			genDetail.userCustomerID = "1"
		} else if genDetail.userTypeID == "200" || genDetail.userTypeID == "201" || genDetail.userTypeID == "400" {
			whereClause = "WHERE customer_id = ?;"
		}

		customerContactSQL, err := dbDetail.connection.Query(`SELECT
							customer_id,
							customer_name,
							customer_site_address_line_1,
					                customer_site_address_line_2,
					                customer_site_city_town_village,
					                customer_site_county_state_region,
					                customer_site_postcode_zip_code,
					                customer_site_country,
					                customer_site_contact_email,
					                customer_site_contact_number,
					                customer_invoice_address_line_1,
					                customer_invoice_address_line_2,
					                customer_invoice_city_town_village,
					                customer_invoice_county_state_region,
					                customer_invoice_postcode_zip_code,
					                customer_invoice_country,
					                customer_invoice_contact_email,
					                customer_invoice_contact_number
					              FROM
					  	        yap.view___customer_detail
						      `+whereClause, genDetail.userCustomerID)

		// Error
		if err != nil {
			panic(err)

		}

		for customerContactSQL.Next() {

			err = customerContactSQL.Scan(
				&customerID,
				&customerName,
				&customerSiteAddressLine1,
				&customerSiteAddressLine2,
				&customerSiteCityTownVillage,
				&customerSiteCountyStateRegion,
				&customerSitePostcodeZipCode,
				&customerSiteCountry,
				&customerSiteContactEmail,
				&customerSiteContactNumber,
				&customerInvoiceAddressLine1,
				&customerInvoiceAddressLine2,
				&customerInvoiceCityTownVillage,
				&customerInvoiceCountyStateRegion,
				&customerInvoicePostcodeZipCode,
				&customerInvoiceCountry,
				&customerInvoiceContactEmail,
				&customerInvoiceContactNumber,
			)

			// Error
			if err != nil {
				panic(err)
			}
			fmt.Fprintf(w, "        <tr>")
			fmt.Fprintf(w, "          <td>"+customerID+"</td>")
			fmt.Fprintf(w, "          <td>"+customerName+"</td>")
			fmt.Fprintf(w, "          <td style=\"text-align: left;\">"+customerSiteAddressLine1+"&nbsp<br>"+customerSiteAddressLine2+"&nbsp<br>"+customerSiteCityTownVillage+"&nbsp<br>"+customerSiteCountyStateRegion+"&nbsp<br><br>"+customerSitePostcodeZipCode+"&nbsp<br><br>"+customerSiteCountry+"&nbsp</td>")
			fmt.Fprintf(w, "          <td><a href=\"mailto:"+customerSiteContactEmail+"\">"+customerSiteContactEmail+"</a></td>")
			fmt.Fprintf(w, "          <td><a href=\"tel:"+customerSiteContactNumber+"\">"+customerSiteContactNumber+"</a></td>")
			fmt.Fprintf(w, "          <td style=\"text-align: left;\">"+customerInvoiceAddressLine1+"&nbsp<br>"+customerInvoiceAddressLine2+"&nbsp<br>"+customerInvoiceCityTownVillage+"&nbsp<br>"+customerInvoiceCountyStateRegion+"&nbsp<br><br>"+customerInvoicePostcodeZipCode+"&nbsp<br><br>"+customerInvoiceCountry+"&nbsp</td>")
			fmt.Fprintf(w, "          <td><a href=\"mailto:"+customerInvoiceContactEmail+"\">"+customerInvoiceContactEmail+"</a></td>")
			fmt.Fprintf(w, "          <td><a href=\"tel:"+customerInvoiceContactNumber+"\">"+customerInvoiceContactNumber+"</a></td>")
			fmt.Fprintf(w, "        </tr>")
		}

		fmt.Fprintf(w, "      </table>")
		if genDetail.userTypeID == "100" {
			var filterTableJSArgument jsFunctionParameter
			filterTableJSArgument.tableID = "customer-contact-table"
			// JS filter function for the customer ID in the customer contact table
			filterTableJSArgument.funcNameJS = "customerContactSearchCustomerID"
			filterTableJSArgument.inputID = "customer-contact-input-customer-id"
			filterTableJSArgument.columnNumber = 0
			filterTableJS(w, filterTableJSArgument)
			// JS filter function for the customer name in the customer contact table
			filterTableJSArgument.funcNameJS = "customerContactSearchCustomerName"
			filterTableJSArgument.inputID = "customer-contact-input-customer-name"
			filterTableJSArgument.columnNumber = 1
			filterTableJS(w, filterTableJSArgument)
			// JS filter function for site address in the customer contact table
			filterTableJSArgument.funcNameJS = "customerContactSearchSiteAddress"
			filterTableJSArgument.inputID = "customer-contact-input-site-address"
			filterTableJSArgument.columnNumber = 2
			filterTableJS(w, filterTableJSArgument)
			// JS filter function for site email in the customer contact table
			filterTableJSArgument.funcNameJS = "customerContactSearchSiteEmail"
			filterTableJSArgument.inputID = "customer-contact-input-site-email"
			filterTableJSArgument.columnNumber = 3
			filterTableJS(w, filterTableJSArgument)
			// JS filter function for site phone in the customer contact table
			filterTableJSArgument.funcNameJS = "customerContactSearchSitePhone"
			filterTableJSArgument.inputID = "customer-contact-input-site-phone"
			filterTableJSArgument.columnNumber = 4
			filterTableJS(w, filterTableJSArgument)
			// JS filter function for invoice address in the customer contact table
			filterTableJSArgument.funcNameJS = "customerContactSearchInvoiceAddress"
			filterTableJSArgument.inputID = "customer-contact-input-invoice-address"
			filterTableJSArgument.columnNumber = 5
			filterTableJS(w, filterTableJSArgument)
			// JS filter function for invoice email in the customer contact table
			filterTableJSArgument.funcNameJS = "customerContactSearchInvoiceEmail"
			filterTableJSArgument.inputID = "customer-contact-input-invoice-email"
			filterTableJSArgument.columnNumber = 6
			filterTableJS(w, filterTableJSArgument)
			// JS filter function for invoice phone in the customer contact table
			filterTableJSArgument.funcNameJS = "customerContactSearchInvoicePhone"
			filterTableJSArgument.inputID = "customer-contact-input-invoice-phone"
			filterTableJSArgument.columnNumber = 7
			filterTableJS(w, filterTableJSArgument)
		}
		var exportCSVJSArgument jsFunctionParameter
		exportCSVJSArgument.funcNameJS = "CustomerContact"
		exportCSVJSArgument.tableID = "customer-contact-table"
		exportCSVJSArgument.fileName = "YAP_customer_contact_details"
		exportCSVJSArgument.pathURL = "customer"
		exportCSVJS(w, exportCSVJSArgument)
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")

		// Customer Miscellaneous Information Table
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<table id=\"table\" class=\"table-customer\">")
		fmt.Fprintf(w, "  <tr>")
		if genDetail.userTypeID == "100" {
			fmt.Fprintf(w, "    <th class=\"table-title\";>All Customer Miscellaneous Information:</th>")
		} else {
			fmt.Fprintf(w, "    <th class=\"table-title\";>Miscellaneous Information for Customer</th>")
		}
		fmt.Fprintf(w, "  </tr>")
		if genDetail.userTypeID == "100" {
			fmt.Fprintf(w, "  <tr>")
			fmt.Fprintf(w, "    <th>")
			fmt.Fprintf(w, "    <br>")
			var inputTableHTMLArgument jsFunctionParameter
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "customer-miscellaneous-input-customer-id"
			inputTableHTMLArgument.funcNameJS = "customerMiscellaneousSearchCustomerID"
			inputTableHTMLArgument.placeholder = "Customer ID"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "customer-miscellaneous-input-customer-name"
			inputTableHTMLArgument.funcNameJS = "customerMiscellaneousSearchCustomerName"
			inputTableHTMLArgument.placeholder = "Customer Name"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "customer-miscellaneous-input-date-time"
			inputTableHTMLArgument.funcNameJS = "customerMiscellaneousSearchDateTime"
			inputTableHTMLArgument.placeholder = "Date & Time Created"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			fmt.Fprintf(w, "    <br>")
			fmt.Fprintf(w, "    <br>")
			fmt.Fprintf(w, "    </th>")
			fmt.Fprintf(w, "  </tr>")
		}
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		exportCSVButtonHTMLArgument.funcNameJS = "CustomerMiscellaneous"
		exportCSVButtonHTMLArgument.buttonCSS = "button-customer"
		exportCSVButtonHTML(w, exportCSVButtonHTMLArgument)
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table id=\"customer-miscellaneous-table\" class=\"table-customer\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <th>Customer ID</th>")
		fmt.Fprintf(w, "          <th>Customer Name</th>")
		fmt.Fprintf(w, "          <th>Customer Added</th>")
		fmt.Fprintf(w, "          <th>Customer Information</th>")
		fmt.Fprintf(w, "          <th>Customer Pricing/Contract Length</th>")
		fmt.Fprintf(w, "        </tr>")

		customerMiscellaneousSQL, err := dbDetail.connection.Query(`SELECT
							customer_id,
							customer_name,
							customer_date_time_added,
							customer_uk_based,
							customer_consumer_type,
							customer_uk_vat_registered,
							customer_uk_vat_number,
							customer_reselling_minutes,
							customer_pbx_limit,
							customer_pbx_sales_tax_rate,
							customer_pbx_sales_tax_status,
                        				customer_ext_sales_tax_rate,
                        				customer_ext_sales_tax_status,
							customer_pbx_setup_price,
							customer_pbx_rental_price,
							customer_pbx_cease_price,
							customer_pbx_contract_length,
							customer_ext_setup_price,
							customer_ext_rental_price,
							customer_ext_cease_price,
							customer_ext_contract_length					              
					              FROM
					  	        yap.view___customer_detail
						      `+whereClause, genDetail.userCustomerID)

		// Error
		if err != nil {
			panic(err)

		}

		for customerMiscellaneousSQL.Next() {

			err = customerMiscellaneousSQL.Scan(
				&customerID,
				&customerName,
				&customerDateTimeAdded,
				&customerUKBased,
				&customerConsumerType,
				&customerUKVATRegistered,
				&customerUKVATNumber,
				&customerResellingMinutes,
				&customerPBXLimit,
				&customerPBXSalesTaxRate,
				&customerPBXSalesTaxStatus,
				&customerExtSalesTaxRate,
				&customerExtSalesTaxStatus,
				&customerPBXSetupPrice,
				&customerPBXRentalPrice,
				&customerPBXCeasePrice,
				&customerPBXContractLength,
				&customerExtSetupPrice,
				&customerExtRentalPrice,
				&customerExtCeasePrice,
				&customerExtContractLength,
			)

			// Error
			if err != nil {
				panic(err)
			}
			fmt.Fprintf(w, "        <tr>")
			fmt.Fprintf(w, "          <td>"+customerID+"</td>")
			fmt.Fprintf(w, "          <td>"+customerName+"</td>")
			fmt.Fprintf(w, "          <td>"+formatDateTime(customerDateTimeAdded)+"</td>")
			fmt.Fprintf(w, "          <td style=\"text-align: left; vertical-align: top;\">")
			fmt.Fprintf(w, "            <b>UK Based:</b> "+customerUKBased+"<br>")
			fmt.Fprintf(w, "            <b>Consumer Type:</b> "+customerConsumerType+"<br>")
			fmt.Fprintf(w, "            <b>UK VAT Registered:</b> "+customerUKVATRegistered+"<br>")
			fmt.Fprintf(w, "            <b>UK VAT Number:</b> "+customerUKVATNumber+"<br>")
			fmt.Fprintf(w, "            <b>Reselling Minutes:</b> "+customerResellingMinutes+"<br>")
			fmt.Fprintf(w, "            <b>PBX Limit:</b> "+customerPBXLimit+"<br>")
			fmt.Fprintf(w, "            <b>Ext Default Limit:</b> "+genDetail.defaultExtLimit+"<br>")
			fmt.Fprintf(w, "            <b>PBX Sales Tax Rate:</b> "+customerPBXSalesTaxRate+"&#37<br>")
			fmt.Fprintf(w, "            <b>PBX Sales Tax Status:</b> "+customerPBXSalesTaxStatus+"<br>")
			fmt.Fprintf(w, "            <b>Ext Sales Tax Rate:</b> "+customerExtSalesTaxRate+"&#37<br>")
			fmt.Fprintf(w, "            <b>Ext Sales Tax Status:</b> "+customerExtSalesTaxStatus)
			fmt.Fprintf(w, "          </td>")
			fmt.Fprintf(w, "          <td style=\"text-align: left; vertical-align: top;\">")
			fmt.Fprintf(w, "            <b>PBX Setup Price:</b> "+genDetail.currencySymbol+customerPBXSetupPrice+"<br>")
			fmt.Fprintf(w, "            <b>PBX Rental Price:</b> "+genDetail.currencySymbol+customerPBXRentalPrice+"<br>")
			fmt.Fprintf(w, "            <b>PBX Cease Price:</b> "+genDetail.currencySymbol+customerPBXCeasePrice+"<br>")
			fmt.Fprintf(w, "            <b>PBX Contract Length:</b> "+customerPBXContractLength+"<br>")
			fmt.Fprintf(w, "            <b>SIP Ext Setup Price:</b> "+genDetail.currencySymbol+customerExtSetupPrice+"<br>")
			fmt.Fprintf(w, "            <b>SIP Ext Rental Price:</b> "+genDetail.currencySymbol+customerExtRentalPrice+"<br>")
			fmt.Fprintf(w, "            <b>SIP Ext Cease Price:</b> "+genDetail.currencySymbol+customerExtCeasePrice+"<br>")
			fmt.Fprintf(w, "            <b>SIP Ext Contract Length:</b> "+customerExtContractLength+"<br>")
			fmt.Fprintf(w, "          </td>")
			fmt.Fprintf(w, "        </tr>")
		}

		fmt.Fprintf(w, "      </table>")
		if genDetail.userTypeID == "100" {
			var filterTableJSArgument jsFunctionParameter
			filterTableJSArgument.tableID = "customer-miscellaneous-table"
			// Call JS filter function for the customer name in the customer miscellaneous table
			filterTableJSArgument.funcNameJS = "customerMiscellaneousSearchCustomerID"
			filterTableJSArgument.inputID = "customer-miscellaneous-input-customer-id"
			filterTableJSArgument.columnNumber = 0
			filterTableJS(w, filterTableJSArgument)
			// Call JS filter function for the customer ID in the customer miscellaneous table
			filterTableJSArgument.funcNameJS = "customerMiscellaneousSearchCustomerName"
			filterTableJSArgument.inputID = "customer-miscellaneous-input-customer-name"
			filterTableJSArgument.columnNumber = 1
			filterTableJS(w, filterTableJSArgument)
			// Call JS filter function for date and time in the customer miscellaneous table
			filterTableJSArgument.funcNameJS = "customerMiscellaneousSearchDateTime"
			filterTableJSArgument.inputID = "customer-miscellaneous-input-date-time"
			filterTableJSArgument.columnNumber = 2
			filterTableJS(w, filterTableJSArgument)
		}
		exportCSVJSArgument.funcNameJS = "CustomerMiscellaneous"
		exportCSVJSArgument.tableID = "customer-miscellaneous-table"
		exportCSVJSArgument.fileName = "YAP_customer_miscellaneous_details"
		exportCSVJSArgument.pathURL = "customer"
		exportCSVJS(w, exportCSVJSArgument)
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</div>")
		if genDetail.userTypeID == "100" {
			var toggleDivJSArgument jsFunctionParameter
			toggleDivJSArgument.funcNameJS = "toggleCustomer"
			toggleDivJSArgument.divID = "customer-div"
			toggleDivJS(w, toggleDivJSArgument)
		}
	} else {
		panic("customerList function should only be called with account type ID 100, 200, 201, 400")
	}
}

// Add customer function
func customerAdd(w http.ResponseWriter, r *http.Request, dbDetail databaseFunctionParameter, genDetail generalFunctionParameter) {

	// Only account type ID 100 should be able to use this function
	if genDetail.userTypeID == "100" {

		fmt.Fprintf(w, "<form method=\"POST\" action=\"/customer\">")
		fmt.Fprintf(w, "<table class=\"table-add\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>Add a New Customer<br>(Typing GEN in the Customer ID Box Will Automatically Generate a Random Customer ID)</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_customer_input_customer_id", "Customer ID<br>(Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_customer_input_customer_name", "Customer Name<br>(Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		ukBasedList := yesNoSlice()
		selectSingleHTML(w, "add_customer_select_uk_based", "Customer UK Based<br>(Cannot Be Empty)", ukBasedList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		resellingMinutesList := yesNoSlice()
		selectSingleHTML(w, "add_customer_select_reselling_minutes", "Customer Reselling Minutes<br>(Cannot Be Empty)", resellingMinutesList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		dbDetail.table = "consumer_type_lookup"
		dbDetail.column = "consumer_type"
		consumerTypeList := singleColumnSlice(dbDetail)
		consumerTypeList = append([]string{""}, consumerTypeList...)
		selectSingleHTML(w, "add_customer_select_consumer_type", "Consumer Type<br>(Cannot Be Empty)", consumerTypeList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		ukVATRegisteredList := yesNoSlice()
		selectSingleHTML(w, "add_customer_select_uk_vat_registered", "UK VAT Registered<br>(Cannot Be Empty)", ukVATRegisteredList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_customer_input_uk_vat_number", "UK VAT Number<br>(Cannot Be Empty if UK VAT Registered yes)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		pbxLimitList := pbxLimitSlice()
		selectSingleHTML(w, "add_customer_select_pbx_limit", "PBX Limt<br>(Cannot Be Empty)", pbxLimitList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		dbDetail.table = "sales_tax_rate_lookup"
		dbDetail.column = "sales_tax_rate"
		salesTaxRateList := singleColumnSlice(dbDetail)
		salesTaxRateList = append([]string{""}, salesTaxRateList...)
		selectSingleHTML(w, "add_customer_select_pbx_sales_tax_rate", "PBX Sales Tax Rate &#37<br>(Cannot Be Empty)", salesTaxRateList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		salesTaxStatusList := []string{"", "TAXABLE", "EXEMPT"}
		selectSingleHTML(w, "add_customer_select_pbx_sales_tax_status", "PBX Sales Tax Status<br>(Cannot Be Empty)", salesTaxStatusList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		selectSingleHTML(w, "add_customer_select_ext_sales_tax_rate", "Ext Sales Tax Rate &#37<br>(Cannot Be Empty)", salesTaxRateList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		selectSingleHTML(w, "add_customer_select_ext_sales_tax_status", "Ext Sales Tax Status<br>(Cannot Be Empty)", salesTaxStatusList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td style=\"border: none;\">")
		fmt.Fprintf(w, "            <br>")
		fmt.Fprintf(w, "            <br>")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_customer_input_pbx_setup_price", "PBX Setup Price (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_customer_input_pbx_rental_price", "PBX Rental Price (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_customer_input_pbx_cease_price", "PBX Cease Price (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		dbDetail.table = "contract_length_lookup"
		dbDetail.column = "contract_length"
		contractLengthList := singleColumnSlice(dbDetail)
		contractLengthList = append([]string{""}, contractLengthList...)
		selectSingleHTML(w, "add_customer_select_pbx_contract_length", "PBX Contract Length", contractLengthList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_customer_input_ext_setup_price", "EXT Setup Price (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_customer_input_ext_rental_price", "EXT Rental Price (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_customer_input_ext_cease_price", "EXT Cease Price (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		selectSingleHTML(w, "add_customer_select_ext_contract_length", "EXT Contract Length", contractLengthList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td style=\"border: none;\">")
		fmt.Fprintf(w, "            <br>")
		fmt.Fprintf(w, "            <br>")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_customer_input_site_address_line_1", "Site Address Line One")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_customer_input_site_address_line_2", "Site Address Line Two")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_customer_input_site_city_town_village", "Site City/Town/Village")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_customer_input_site_county_state_region", "Site County/State/Region")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_customer_input_site_postcode_zip_code", "Site Postcode/Zip Code")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_customer_input_site_country", "Site Country")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_customer_input_site_contact_email", "Site Email (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_customer_input_site_contact_number", "Site Phone (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td style=\"border: none;\">")
		fmt.Fprintf(w, "            <br>")
		fmt.Fprintf(w, "            <br>")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_customer_input_invoice_address_line_1", "Invoice Address Line One")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_customer_input_invoice_address_line_2", "Invoice Address Line Two")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_customer_input_invoice_city_town_village", "Invoice City/Town/Village")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_customer_input_invoice_county_state_region", "Invoice County/State/Region")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_customer_input_invoice_postcode_zip_code", "Invoice Postcode/Zip Code")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_customer_input_invoice_country", "Invoice Country")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_customer_input_invoice_contact_email", "Invoice Email (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_customer_input_invoice_contact_number", "Invoice Phone (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><input type=\"submit\" value=\"Create Customer\"></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</form>")

		addCustomerInputCustomerID := r.FormValue("add_customer_input_customer_id")
		addCustomerInputCustomerName := r.FormValue("add_customer_input_customer_name")
		addCustomerSelectUKBased := r.FormValue("add_customer_select_uk_based")
		addCustomerSelectResellingMinutes := r.FormValue("add_customer_select_reselling_minutes")
		addCustomerSelectConsumerType := r.FormValue("add_customer_select_consumer_type")
		addCustomerSelectUKVATRegistered := r.FormValue("add_customer_select_uk_vat_registered")
		addCustomerInputUKVATNumber := r.FormValue("add_customer_input_uk_vat_number")
		addCustomerSelectPBXLimit := r.FormValue("add_customer_select_pbx_limit")
		addCustomerSelectPBXSalesTaxRate := r.FormValue("add_customer_select_pbx_sales_tax_rate")
		addCustomerSelectPBXSalesTaxStatus := r.FormValue("add_customer_select_pbx_sales_tax_status")
		addCustomerSelectExtSalesTaxRate := r.FormValue("add_customer_select_ext_sales_tax_rate")
		addCustomerSelectExtSalesTaxStatus := r.FormValue("add_customer_select_ext_sales_tax_status")

		addCustomerInputPBXSetupPrice := r.FormValue("add_customer_input_pbx_setup_price")
		addCustomerInputPBXRentalPrice := r.FormValue("add_customer_input_pbx_rental_price")
		addCustomerInputPBXCeasePrice := r.FormValue("add_customer_input_pbx_cease_price")
		addCustomerSelectPBXContractLength := r.FormValue("add_customer_select_pbx_contract_length")
		addCustomerInputExtSetupPrice := r.FormValue("add_customer_input_ext_setup_price")
		addCustomerInputExtRentalPrice := r.FormValue("add_customer_input_ext_rental_price")
		addCustomerInputExtCeasePrice := r.FormValue("add_customer_input_ext_cease_price")
		addCustomerSelectExtContractLength := r.FormValue("add_customer_select_ext_contract_length")

		addCustomerInputSiteAddressLine1 := r.FormValue("add_customer_input_site_address_line_1")
		addCustomerInputSiteAddressLine2 := r.FormValue("add_customer_input_site_address_line_2")
		addCustomerInputSiteCityTownVillage := r.FormValue("add_customer_input_site_city_town_village")
		addCustomerInputSiteCountyStateRegion := r.FormValue("add_customer_input_site_county_state_region")
		addCustomerInputSitePostcodeZipCode := r.FormValue("add_customer_input_site_postcode_zip_code")
		addCustomerInputSiteCountry := r.FormValue("add_customer_input_site_country")
		addCustomerInputSiteContactEmail := r.FormValue("add_customer_input_site_contact_email")
		addCustomerInputSiteContactNumber := r.FormValue("add_customer_input_site_contact_number")

		addCustomerInputInvoiceAddressLine1 := r.FormValue("add_customer_input_invoice_address_line_1")
		addCustomerInputInvoiceAddressLine2 := r.FormValue("add_customer_input_invoice_address_line_2")
		addCustomerInputInvoiceCityTownVillage := r.FormValue("add_customer_input_invoice_city_town_village")
		addCustomerInputInvoiceCountyStateRegion := r.FormValue("add_customer_input_invoice_county_state_region")
		addCustomerInputInvoicePostcodeZipCode := r.FormValue("add_customer_input_invoice_postcode_zip_code")
		addCustomerInputInvoiceCountry := r.FormValue("add_customer_input_invoice_country")
		addCustomerInputInvoiceContactEmail := r.FormValue("add_customer_input_invoice_contact_email")
		addCustomerInputInvoiceContactNumber := r.FormValue("add_customer_input_invoice_contact_number")

		// Validate the customer ID
		validateCustomerID := validateInput(addCustomerInputCustomerID, "alphaNum")
		// Validate the customer name
		validateCustomerName := validateInput(addCustomerInputCustomerName, "alphaNum")
		// Validate UK based
		validateUKBased := slices.Contains(ukBasedList, addCustomerSelectUKBased)
		// Validate reselling minutes
		validateResellingMinutes := slices.Contains(resellingMinutesList, addCustomerSelectResellingMinutes)
		// Validate consumer type
		validateConsumerType := slices.Contains(consumerTypeList, addCustomerSelectConsumerType)
		// Validate UK VAT registered status
		validateUKVATRegistered := slices.Contains(ukVATRegisteredList, addCustomerSelectUKVATRegistered)
		// Validate UK VAT number
		validateUKVATNumber := validateInput(addCustomerInputUKVATNumber, "alphaNumEmpty")
		// Validate PBX limit
		validatePBXLimit := slices.Contains(pbxLimitList, addCustomerSelectPBXLimit)
		// Validate PBX sales tax rate
		validatePBXSalesTaxRate := slices.Contains(salesTaxRateList, addCustomerSelectPBXSalesTaxRate)
		// Validate PBX sales tax status
		validatePBXSalesTaxStatus := slices.Contains(salesTaxStatusList, addCustomerSelectPBXSalesTaxStatus)
		// Validate ext sales tax rate
		validateExtSalesTaxRate := slices.Contains(salesTaxRateList, addCustomerSelectExtSalesTaxRate)
		// Validate ext sales tax status
		validateExtSalesTaxStatus := slices.Contains(salesTaxStatusList, addCustomerSelectExtSalesTaxStatus)

		// Validate the PBX setup price
		validatePBXSetupPrice := validateInput(addCustomerInputPBXSetupPrice, "price")
		// Validate the PBX rental price
		validatePBXRentalPrice := validateInput(addCustomerInputPBXRentalPrice, "price")
		// Validate the PBX cease price
		validatePBXCeasePrice := validateInput(addCustomerInputPBXCeasePrice, "price")
		// Validate the PBX contract length
		validatePBXContractLength := slices.Contains(contractLengthList, addCustomerSelectPBXContractLength)
		// Validate the ext setup price
		validateExtSetupPrice := validateInput(addCustomerInputExtSetupPrice, "price")
		// Validate the ext rental price
		validateExtRentalPrice := validateInput(addCustomerInputExtRentalPrice, "price")
		// Validate the ext cease price
		validateExtCeasePrice := validateInput(addCustomerInputExtCeasePrice, "price")
		// Validate the ext contract length
		validateExtContractLength := slices.Contains(contractLengthList, addCustomerSelectExtContractLength)

		// Validate site address line one
		validateSiteAddressLine1 := validateInput(addCustomerInputSiteAddressLine1, "alphaNumEmpty")
		// Validate site address line two
		validateSiteAddressLine2 := validateInput(addCustomerInputSiteAddressLine2, "alphaNumEmpty")
		// Validate site city/town/village
		validateSiteCityTownVillage := validateInput(addCustomerInputSiteCityTownVillage, "alphaNumEmpty")
		// Validate site county/state/region
		validateSiteCountyStateRegion := validateInput(addCustomerInputSiteCountyStateRegion, "alphaNumEmpty")
		// Validate site postcode/zip code
		validateSitePostcodeZipCode := validateInput(addCustomerInputSitePostcodeZipCode, "alphaNumEmpty")
		// Validate site country
		validateSiteCountry := validateInput(addCustomerInputSiteCountry, "alphaNumEmpty")
		// Validate site contact emial
		validateSiteContactEmail := validateInput(addCustomerInputSiteContactEmail, "email")
		// Validate Site contact phone number
		validateSiteContactNumber := validateInput(addCustomerInputSiteContactNumber, "phoneNumber")

		// Validate invoice address line one
		validateInvoiceAddressLine1 := validateInput(addCustomerInputInvoiceAddressLine1, "alphaNumEmpty")
		// Validate invoice address line two
		validateInvoiceAddressLine2 := validateInput(addCustomerInputInvoiceAddressLine2, "alphaNumEmpty")
		// Validate invoice city/town/village
		validateInvoiceCityTownVillage := validateInput(addCustomerInputInvoiceCityTownVillage, "alphaNumEmpty")
		// Validate invoice county/state/region
		validateInvoiceCountyStateRegion := validateInput(addCustomerInputInvoiceCountyStateRegion, "alphaNumEmpty")
		// Validate invoice postcode/zip code
		validateInvoicePostcodeZipCode := validateInput(addCustomerInputInvoicePostcodeZipCode, "alphaNumEmpty")
		// Validate invoice country
		validateInvoiceCountry := validateInput(addCustomerInputInvoiceCountry, "alphaNumEmpty")
		// Validate invoice contact emial
		validateInvoiceContactEmail := validateInput(addCustomerInputInvoiceContactEmail, "email")
		// Validate invoice contact phone number
		validateInvoiceContactNumber := validateInput(addCustomerInputInvoiceContactNumber, "phoneNumber")

		if addCustomerInputCustomerID == "" && addCustomerInputCustomerName == "" && addCustomerSelectUKBased == "" && addCustomerSelectResellingMinutes == "" && addCustomerSelectConsumerType == "" && addCustomerSelectUKVATRegistered == "" && addCustomerSelectPBXLimit == "" && addCustomerSelectPBXSalesTaxRate == "" && addCustomerSelectPBXSalesTaxStatus == "" && addCustomerSelectExtSalesTaxRate == "" && addCustomerSelectExtSalesTaxStatus == "" && addCustomerInputPBXSetupPrice == "" && addCustomerInputPBXRentalPrice == "" && addCustomerInputPBXCeasePrice == "" && addCustomerInputExtSetupPrice == "" && addCustomerInputExtRentalPrice == "" && addCustomerInputExtCeasePrice == "" && addCustomerInputSiteContactEmail == "" && addCustomerInputSiteContactNumber == "" && addCustomerInputInvoiceContactEmail == "" && addCustomerInputInvoiceContactNumber == "" {
			// Do Nothing
		} else if validateCustomerID == false {
			messageHTML(w, validationMessageCustomerID, "warning")
		} else if validateCustomerName == false {
			messageHTML(w, validationMessageCustomerName, "warning")
		} else if validateUKBased == false || addCustomerSelectUKBased == "" {
			messageHTML(w, validationMessageCustomerUKBased, "warning")
		} else if validateResellingMinutes == false || addCustomerSelectResellingMinutes == "" {
			messageHTML(w, validationMessageCustomerResellingMinutes, "warning")
		} else if validateConsumerType == false || addCustomerSelectConsumerType == "" {
			messageHTML(w, validationMessageCustomerConsumerType, "warning")
		} else if validateUKVATRegistered == false || addCustomerSelectUKVATRegistered == "" {
			messageHTML(w, validationMessageCustomerUKVATRegistered, "warning")
		} else if addCustomerSelectUKVATRegistered == "yes" && addCustomerInputUKVATNumber == "" {
			messageHTML(w, validationMessageCustomerUKVATRegisteredEmpty, "warning")
		} else if validateUKVATNumber == false {
			messageHTML(w, validationMessageCustomerUKVATNumber, "warning")
		} else if validatePBXLimit == false || addCustomerSelectPBXLimit == "" {
			messageHTML(w, validationMessageCustomerPBXLimit, "warning")
		} else if validatePBXSalesTaxRate == false || addCustomerSelectPBXSalesTaxRate == "" {
			messageHTML(w, validationMessageCustomerPBXSalesTaxRate, "warning")
		} else if validatePBXSalesTaxStatus == false || addCustomerSelectPBXSalesTaxStatus == "" {
			messageHTML(w, validationMessageCustomerPBXSalesTaxStatus, "warning")
		} else if validateExtSalesTaxRate == false || addCustomerSelectExtSalesTaxRate == "" {
			messageHTML(w, validationMessageCustomerExtSalesTaxRate, "warning")
		} else if validateExtSalesTaxStatus == false || addCustomerSelectExtSalesTaxStatus == "" {
			messageHTML(w, validationMessageCustomerExtSalesTaxStatus, "warning")
		} else if validatePBXSetupPrice == false {
			messageHTML(w, validationMessageCustomerSetupPrice, "warning")
		} else if validatePBXRentalPrice == false {
			messageHTML(w, validationMessageCustomerRentalPrice, "warning")
		} else if validatePBXCeasePrice == false {
			messageHTML(w, validationMessageCustomerCeasePrice, "warning")
		} else if validatePBXContractLength == false {
			messageHTML(w, validationMessageContractLength, "warning")
		} else if validateExtSetupPrice == false {
			messageHTML(w, validationMessageCustomerSetupPrice, "warning")
		} else if validateExtRentalPrice == false {
			messageHTML(w, validationMessageCustomerRentalPrice, "warning")
		} else if validateExtCeasePrice == false {
			messageHTML(w, validationMessageCustomerCeasePrice, "warning")
		} else if validateExtContractLength == false {
			messageHTML(w, validationMessageContractLength, "warning")
		} else if validateSiteAddressLine1 == false {
			messageHTML(w, validationMessageAddresslineOne, "warning")
		} else if validateSiteAddressLine2 == false {
			messageHTML(w, validationMessageAddresslineTwo, "warning")
		} else if validateSiteCityTownVillage == false {
			messageHTML(w, validationMessageCityTownVillage, "warning")
		} else if validateSiteCountyStateRegion == false {
			messageHTML(w, validationMessageCountyStateRegion, "warning")
		} else if validateSitePostcodeZipCode == false {
			messageHTML(w, validationMessagePostcodeZipCode, "warning")
		} else if validateSiteCountry == false {
			messageHTML(w, validationMessageCountry, "warning")
		} else if validateSiteContactEmail == false {
			messageHTML(w, validationMessageCustomerSiteEmail, "warning")
		} else if validateSiteContactNumber == false {
			messageHTML(w, validationMessageCustomerSitePhoneNumber, "warning")
		} else if validateInvoiceAddressLine1 == false {
			messageHTML(w, validationMessageAddresslineOne, "warning")
		} else if validateInvoiceAddressLine2 == false {
			messageHTML(w, validationMessageAddresslineTwo, "warning")
		} else if validateInvoiceCityTownVillage == false {
			messageHTML(w, validationMessageCityTownVillage, "warning")
		} else if validateInvoiceCountyStateRegion == false {
			messageHTML(w, validationMessageCountyStateRegion, "warning")
		} else if validateInvoicePostcodeZipCode == false {
			messageHTML(w, validationMessagePostcodeZipCode, "warning")
		} else if validateInvoiceCountry == false {
			messageHTML(w, validationMessageCountry, "warning")
		} else if validateInvoiceContactEmail == false {
			messageHTML(w, validationMessageInvoiceEmail, "warning")
		} else if validateInvoiceContactNumber == false {
			messageHTML(w, validationMessageInvoicePhoneNumber, "warning")
		} else {
			if addCustomerInputCustomerID == "Gen" || addCustomerInputCustomerID == "GEN" || addCustomerInputCustomerID == "gen" || addCustomerInputCustomerID == "Generate" || addCustomerInputCustomerID == "GENERATE" || addCustomerInputCustomerID == "generate" {
				addCustomerInputCustomerID = genID()
			}

			dbDetail.table = "view___customer_detail"
			dbDetail.column = "customer_id"
			dbDetail.columnWhere = "customer_id"
			dbDetail.columnWhereValue = addCustomerInputCustomerID

			checkCustomerExist := selectWhere(dbDetail)

			if checkCustomerExist == addCustomerInputCustomerID {
				messageHTML(w, validationMessageCustomerlAlreadyExist, "warning")
			} else {

				// Convert string values to a float64 to use the math package to round to the nearest two decimal places
				addCustomerInputPBXSetupPriceFloat64 := stringToFloat64(addCustomerInputPBXSetupPrice)
				addCustomerInputPBXRentalPriceFloat64 := stringToFloat64(addCustomerInputPBXRentalPrice)
				addCustomerInputPBXCeasePriceFloat64 := stringToFloat64(addCustomerInputPBXCeasePrice)
				addCustomerInputExtSetupPriceFloat64 := stringToFloat64(addCustomerInputExtSetupPrice)
				addCustomerInputExtRentalPriceFloat64 := stringToFloat64(addCustomerInputExtRentalPrice)
				addCustomerInputExtCeasePriceFloat64 := stringToFloat64(addCustomerInputExtCeasePrice)

				dbDetail.connection.Query(`INSERT 
                                   INTO
                                 customer (
                                   id,
                                   name,
                                   uk_based,
                                   reselling_minutes,
                                   consumer_type,
                                   uk_vat_registered,
                                   uk_vat_number,
                                   pbx_limit,
                                   pbx_sales_tax_rate,
                                   pbx_sales_tax_status,
                                   ext_sales_tax_rate,
                                   ext_sales_tax_status,
                                   pbx_setup_price,
                                   pbx_rental_price,
                                   pbx_cease_price,
                                   pbx_contract_length,  
                                   ext_setup_price,
                                   ext_rental_price,
                                   ext_cease_price,
                                   ext_contract_length
                                 )
                                 VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`,
					addCustomerInputCustomerID,
					addCustomerInputCustomerName,
					addCustomerSelectUKBased,
					addCustomerSelectResellingMinutes,
					addCustomerSelectConsumerType,
					addCustomerSelectUKVATRegistered,
					nullSQL(addCustomerInputUKVATNumber),
					addCustomerSelectPBXLimit,
					addCustomerSelectPBXSalesTaxRate,
					addCustomerSelectPBXSalesTaxStatus,
					addCustomerSelectExtSalesTaxRate,
					addCustomerSelectExtSalesTaxStatus,
					math.Round(addCustomerInputPBXSetupPriceFloat64*100)/100,
					math.Round(addCustomerInputPBXRentalPriceFloat64*100)/100,
					math.Round(addCustomerInputPBXCeasePriceFloat64*100)/100,
					nullSQL(addCustomerSelectPBXContractLength),
					math.Round(addCustomerInputExtSetupPriceFloat64*100)/100,
					math.Round(addCustomerInputExtRentalPriceFloat64*100)/100,
					math.Round(addCustomerInputExtCeasePriceFloat64*100)/100,
					nullSQL(addCustomerSelectExtContractLength))

				dbDetail.connection.Query(`INSERT 
                                   INTO
                                 customer_site_address (
                                   id,
                                   address_line_1,
                                   address_line_2,
                                   city_town_village,
                                   county_state_region,
                                   postcode_zip_code,
                                   country,
                                   contact_email,
                                   contact_number 
                                 )
                                 VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?);`,
					addCustomerInputCustomerID,
					nullSQL(addCustomerInputSiteAddressLine1),
					nullSQL(addCustomerInputSiteAddressLine2),
					nullSQL(addCustomerInputSiteCityTownVillage),
					nullSQL(addCustomerInputSiteCountyStateRegion),
					nullSQL(addCustomerInputSitePostcodeZipCode),
					nullSQL(addCustomerInputSiteCountry),
					nullSQL(addCustomerInputSiteContactEmail),
					nullSQL(addCustomerInputSiteContactNumber))

				dbDetail.connection.Query(`INSERT 
                                   INTO
                                 customer_invoice_address (
                                   id,
                                   address_line_1,
                                   address_line_2,
                                   city_town_village,
                                   county_state_region,
                                   postcode_zip_code,
                                   country,
                                   contact_email,
                                   contact_number 
                                 )
                                 VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?);`,
					addCustomerInputCustomerID,
					nullSQL(addCustomerInputInvoiceAddressLine1),
					nullSQL(addCustomerInputInvoiceAddressLine2),
					nullSQL(addCustomerInputInvoiceCityTownVillage),
					nullSQL(addCustomerInputInvoiceCountyStateRegion),
					nullSQL(addCustomerInputInvoicePostcodeZipCode),
					nullSQL(addCustomerInputInvoiceCountry),
					nullSQL(addCustomerInputInvoiceContactEmail),
					nullSQL(addCustomerInputInvoiceContactNumber))

				checkCustomerCreated := selectWhere(dbDetail)

				if checkCustomerCreated == addCustomerInputCustomerID {
					messageHTML(w, validationMessageCustomerCreated, "success")
				} else {
					messageHTML(w, validationMessageCustomerNotCreated, "success")
				}
			}
		}
	} else {
		panic("customerAdd function should only be called with account type ID 100")
	}
}

// Customer edit function
func customerEdit(w http.ResponseWriter, r *http.Request, dbDetail databaseFunctionParameter, genDetail generalFunctionParameter) {

	// Only account type ID 100 should be able to use this function
	if genDetail.userTypeID == "100" {

		// List of all column names from the customer table
		customerColumnList := [][]string{
			{"name", "Name"},
			{"uk_based", "UK Based"},
			{"reselling_minutes", "Rselling Minutes"},
			{"consumer_type", "Consumer Type"},
			{"uk_vat_registered", "UK VAT Registered"},
			{"uk_vat_number", "UK VAT Number"},
			{"pbx_limit", "PBX Limit"},
			{"pbx_sales_tax_rate", "PBX Sales Tax Rate"},
			{"pbx_sales_tax_status", "PBX Sales Tax Status"},
			{"ext_sales_tax_rate", "EXT Sales Tax Rate"},
			{"ext_sales_tax_status", "EXT Sales Tax Status"},
			{"pbx_setup_price", "PBX Setup Price"},
			{"pbx_rental_price", "PBX Rental Price"},
			{"pbx_cease_price", "PBX Cease Price"},
			{"pbx_contract_length", "PBX Contract Length"},
			{"ext_setup_price", "SIP EXT Setup Price"},
			{"ext_rental_price", "SIP EXT Rental Price"},
			{"ext_cease_price", "SIP EXT Cease Price"},
			{"ext_contract_length", "SIP EXT Contract Length"},
		}

		fmt.Fprintf(w, "<form method=\"POST\" action=\"/customer\">")
		fmt.Fprintf(w, "<table class=\"table-customer\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>Edit Customer Details</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td style=\"text-align: left;\">")
		fmt.Fprintf(w, "      <b><u>Acceptable Values for Columns</u></b><br><br>")
		fmt.Fprintf(w, "      <b>Name:</b> text<br>")
		fmt.Fprintf(w, "      <b>UK Based:</b> yes, no<br>")
		fmt.Fprintf(w, "      <b>Resell Minutes:</b> yes, no<br>")
		dbDetail.table = "consumer_type_lookup"
		dbDetail.column = "consumer_type"
		consumerTypeList := singleColumnSlice(dbDetail)
		fmt.Fprintf(w, "      <b>Consumer Type: </b>")
		fmt.Fprintf(w, strings.Join(consumerTypeList, ", "))
		fmt.Fprintf(w, "      <br>")
		fmt.Fprintf(w, "      <b>UK VAT Registered:</b> yes, no<br>")
		fmt.Fprintf(w, "      <b>UK VAT Number:</b> text, EMPTY<br>")
		fmt.Fprintf(w, "      <b>PBX Limit:</b> 1, 2, 3, 4, 5, 10, 25, 50, 75, 100, 150, 200, 250, 500, 750, 1000, 1500, 2000, 2500, 5000<br>")
		dbDetail.table = "sales_tax_rate_lookup"
		dbDetail.column = "sales_tax_rate"
		salesTaxRateList := singleColumnSlice(dbDetail)
		fmt.Fprintf(w, "      <b>PBX Sales Tax Rate &#37:</b> ")
		fmt.Fprintf(w, strings.Join(salesTaxRateList, ", "))
		fmt.Fprintf(w, "      <br>")
		fmt.Fprintf(w, "      <b>PBX Sales Tax Status:</b> TAXABLE, EXEMPT<br>")
		fmt.Fprintf(w, "      <b>EXT Sales Tax Rate &#37:</b> ")
		fmt.Fprintf(w, strings.Join(salesTaxRateList, ", "))
		fmt.Fprintf(w, "      <br>")
		fmt.Fprintf(w, "      <b>EXT Sales Tax Status:</b> TAXABLE, EXEMPT<br>")
		fmt.Fprintf(w, "      <b>PBX Setup Price:</b> decimal number<br>")
		fmt.Fprintf(w, "      <b>PBX Rental Price:</b> decimal number<br>")
		fmt.Fprintf(w, "      <b>PBX Cease Price:</b> decimal number<br> ")
		dbDetail.table = "contract_length_lookup"
		dbDetail.column = "contract_length"
		contractLengthList := singleColumnSlice(dbDetail)
		fmt.Fprintf(w, "      <b>PBX Contract Length: </b>EMPTY, ")
		fmt.Fprintf(w, strings.Join(contractLengthList, ", "))
		fmt.Fprintf(w, "      <br>")
		fmt.Fprintf(w, "      <b>EXT Setup Price:</b> decimal number<br>")
		fmt.Fprintf(w, "      <b>EXT Rental Price:</b> decimal number<br>")
		fmt.Fprintf(w, "      <b>EXT Cease Pricce:</b> decimal number<br> ")
		fmt.Fprintf(w, "      <b>EXT Contract Length: </b>EMPTY, ")
		fmt.Fprintf(w, strings.Join(contractLengthList, ", "))
		fmt.Fprintf(w, "    </td>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		customerIDNameList, _ := customerSlice(dbDetail)
		selectDoubleHTML(w, "edit_customer_select_customer_id", "Customer", customerIDNameList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		selectDoubleHiddenHTML(w, "edit_customer_select_column", "Column to Edit (Cannot Be Empty)", customerColumnList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "edit_customer_input_new_value", "New Value")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><input type=\"submit\" value=\"Update Customer\"></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</form>")

		editCustomerSelectCustomerID := r.FormValue("edit_customer_select_customer_id")
		editCustomerSelectColumn := r.FormValue("edit_customer_select_column")
		editCustomerInputNewValue := r.FormValue("edit_customer_input_new_value")

		// Validate Customer List
		_, customerIDList := customerSlice(dbDetail)
		validateCustomerID := slices.Contains(customerIDList, editCustomerSelectCustomerID)

		if editCustomerSelectCustomerID == "" && editCustomerSelectColumn == "" && editCustomerInputNewValue == "" {
			// Do Nothing
		} else if validateCustomerID == false {
			messageHTML(w, validationMessageCustomer, "warning")
		} else if editCustomerSelectColumn == "" {
			messageHTML(w, validationMessageCustomerColumn, "warning")
		} else if editCustomerSelectColumn == "name" || editCustomerSelectColumn == "uk_based" || editCustomerSelectColumn == "reselling_minutes" || editCustomerSelectColumn == "consumer_type" || editCustomerSelectColumn == "uk_vat_registered" || editCustomerSelectColumn == "uk_vat_number" || editCustomerSelectColumn == "pbx_contract_length" || editCustomerSelectColumn == "ext_contract_length" {
			// Validate editCustomerInputNewValue is a string
			validateNewValue := validateInput(editCustomerInputNewValue, "alphaNumEmpty")
			if validateNewValue == true {
				dbDetail.connection.Query("UPDATE customer SET "+editCustomerSelectColumn+" = ? WHERE id = ?;", editCustomerInputNewValue, editCustomerSelectCustomerID)
			} else {
				messageHTML(w, validationMessageGenericAlphaNumEmpty, "warning")
			}
		} else if editCustomerSelectColumn == "pbx_setup_price" || editCustomerSelectColumn == "pbx_rental_price" || editCustomerSelectColumn == "pbx_cease_price" || editCustomerSelectColumn == "ext_setup_price" || editCustomerSelectColumn == "ext_rental_price" || editCustomerSelectColumn == "ext_cease_price" {
			// Validate editCustomerSelectColumn is a decimal
			validateNewValue := validateInput(editCustomerInputNewValue, "price")
			if validateNewValue == true {
				// Convert string values to a float64 to use the math package to round to the nearest two decimal places
				editCustomerInputNewValueFloat64 := stringToFloat64(editCustomerInputNewValue)
				dbDetail.connection.Query("UPDATE customer SET "+editCustomerSelectColumn+" = ? WHERE id = ?;", math.Round(editCustomerInputNewValueFloat64*100)/100, editCustomerSelectCustomerID)
			} else {
				messageHTML(w, validationMessageGenericPrice, "warning")
			}
		} else if editCustomerSelectColumn == "pbx_limit" {
			pbxLimitList := pbxLimitSlice()
			// Validate editCustomerInputNewValue is in the pbxLimitList slice
			validateNewValue := slices.Contains(pbxLimitList, editCustomerInputNewValue)
			if validateNewValue == true {
				dbDetail.connection.Query("UPDATE customer SET "+editCustomerSelectColumn+" = ? WHERE id = ?;", editCustomerInputNewValue, editCustomerSelectCustomerID)
			} else {
				messageHTML(w, validationMessagePBX, "warning")
			}
		} else if editCustomerSelectColumn == "pbx_sales_tax_rate" || editCustomerSelectColumn == "ext_sales_tax_rate" {
			salesTaxRateList := singleColumnSlice(dbDetail)
			// Validate editCustomerSelectColumn is in the salesTaxRateList Slice
			validateNewValue := slices.Contains(salesTaxRateList, editCustomerInputNewValue)
			if validateNewValue == true {
				dbDetail.connection.Query("UPDATE customer SET "+editCustomerSelectColumn+" = ? WHERE id = ?;", editCustomerInputNewValue, editCustomerSelectCustomerID)
			} else {
				messageHTML(w, validationMessageGenericInvalidOption, "warning")
			}
		} else if editCustomerSelectColumn == "pbx_sales_tax_status" || editCustomerSelectColumn == "ext_sales_tax_status" {
			salesTaxStatusList := []string{"TAXABLE", "EXEMPT"}
			// Validate editCustomerSelectColumn is in the salesTaxStatusList Slice
			validateNewValue := slices.Contains(salesTaxStatusList, editCustomerInputNewValue)
			if validateNewValue == true {
				dbDetail.connection.Query("UPDATE customer SET "+editCustomerSelectColumn+" = ? WHERE id = ?;", editCustomerInputNewValue, editCustomerSelectCustomerID)
			} else {
				messageHTML(w, validationMessageGenericInvalidOption, "warning")
			}
		} else {
			messageHTML(w, validationMessageCustomerColumn, "warning")
		}

		// Customer site address edit code
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<form method=\"POST\" action=\"/customer\">")
		fmt.Fprintf(w, "<table class=\"table-customer\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>Edit Customer Site Details</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td style=\"text-align: left;\">")
		fmt.Fprintf(w, "      <b><u>Acceptable Values for Columns</u></b><br><br>")
		fmt.Fprintf(w, "      <b>Site Address Line One:</b> text<br>")
		fmt.Fprintf(w, "      <b>Site Address Line Two:</b> text<br>")
		fmt.Fprintf(w, "      <b>Site City Town Village:</b> text<br>")
		fmt.Fprintf(w, "      <b>Site County State Region:</b> text<br>")
		fmt.Fprintf(w, "      <b>Site Postcode Zip Code:</b> text<br>")
		fmt.Fprintf(w, "      <b>Site Country:</b> text<br>")
		fmt.Fprintf(w, "      <b>Site Contact Email:</b> valid email address<br>")
		fmt.Fprintf(w, "      <b>Site Contact Number:</b> phone number in e.164 format<br>")
		fmt.Fprintf(w, "    </td>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		selectDoubleHTML(w, "edit_customer_site_select_customer_id", "Customer", customerIDNameList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		siteColumnList := siteColumnSlice()
		selectDoubleHiddenHTML(w, "edit_customer_site_select_column", "Column to Edit (Cannot Be Empty)", siteColumnList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "edit_customer_site_input_new_value", "New Value")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><input type=\"submit\" value=\"Update Customer Site\"></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</form>")

		editCustomerSiteSelectCustomerID := r.FormValue("edit_customer_site_select_customer_id")
		editCustomerSiteSelectColumn := r.FormValue("edit_customer_site_select_column")
		editCustomerSiteInputNewValue := r.FormValue("edit_customer_site_input_new_value")

		// Validate Customer List
		_, customerSiteIDList := customerSlice(dbDetail)
		validateSiteCustomerID := slices.Contains(customerSiteIDList, editCustomerSiteSelectCustomerID)

		if editCustomerSiteSelectCustomerID == "" && editCustomerSiteSelectColumn == "" && editCustomerSiteInputNewValue == "" {
			// Do Nothing
		} else if validateSiteCustomerID == false {
			messageHTML(w, validationMessageCustomer, "warning")
		} else if editCustomerSiteSelectColumn == "" {
			messageHTML(w, validationMessageCustomerColumn, "warning")
		} else if editCustomerSiteSelectColumn == "address_line_1" || editCustomerSiteSelectColumn == "address_line_2" || editCustomerSiteSelectColumn == "city_town_village" || editCustomerSiteSelectColumn == "county_state_region" || editCustomerSiteSelectColumn == "postcode_zip_code" || editCustomerSiteSelectColumn == "country" {
			// Validate editCustomerSiteInputNewValue is a string
			validateNewValue := validateInput(editCustomerSiteInputNewValue, "alphaNumEmpty")
			if validateNewValue == true {
				dbDetail.connection.Query("UPDATE customer_site_address SET "+editCustomerSiteSelectColumn+" = ? WHERE id = ?;", editCustomerSiteInputNewValue, editCustomerSiteSelectCustomerID)
			} else {
				messageHTML(w, validationMessageGenericAlphaNumEmpty, "warning")
			}
		} else if editCustomerSiteSelectColumn == "contact_email" {
			// Validate editCustomerSiteInputNewValue is a email
			validateNewValue := validateInput(editCustomerSiteInputNewValue, "email")
			if validateNewValue == true {
				dbDetail.connection.Query("UPDATE customer_site_address SET "+editCustomerSiteSelectColumn+" = ? WHERE id = ?;", editCustomerSiteInputNewValue, editCustomerSiteSelectCustomerID)
			} else {
				messageHTML(w, validationMessageCustomerEmail, "warning")
			}
		} else if editCustomerSiteSelectColumn == "contact_number" {
			// Validate editCustomerSiteInputNewValue is a phone number
			validateNewValue := validateInput(editCustomerSiteInputNewValue, "phoneNumber")
			if validateNewValue == true {
				dbDetail.connection.Query("UPDATE customer_site_address SET "+editCustomerSiteSelectColumn+" = ? WHERE id = ?;", editCustomerSiteInputNewValue, editCustomerSiteSelectCustomerID)
			} else {
				messageHTML(w, validationMessageCustomerPhoneNumber, "warning")
			}
		} else {
			messageHTML(w, validationMessageCustomerColumn, "warning")
		}

		// Customer invoice address edit code
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<form method=\"POST\" action=\"/customer\">")
		fmt.Fprintf(w, "<table class=\"table-customer\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>Edit Customer Invoice Details</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td style=\"text-align: left;\">")
		fmt.Fprintf(w, "      <b><u>Acceptable Values for Columns</u></b><br><br>")
		fmt.Fprintf(w, "      <b>Invoice Address Line One:</b> text<br>")
		fmt.Fprintf(w, "      <b>Invoice Address Line Two:</b> text<br>")
		fmt.Fprintf(w, "      <b>Invoice City Town Village:</b> text<br>")
		fmt.Fprintf(w, "      <b>Invoice County State Region:</b> text<br>")
		fmt.Fprintf(w, "      <b>Invoice Postcode Zip Code:</b> text<br>")
		fmt.Fprintf(w, "      <b>Invoice Country:</b> text<br>")
		fmt.Fprintf(w, "      <b>Invoice Contact Email:</b> valid email address<br>")
		fmt.Fprintf(w, "      <b>Invoice Contact Number:</b> phone number in e.164 format<br>")
		fmt.Fprintf(w, "    </td>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		selectDoubleHTML(w, "edit_customer_invoice_select_customer_id", "Customer", customerIDNameList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		invoiceColumnList := invoiceColumnSlice()
		selectDoubleHiddenHTML(w, "edit_customer_invoice_select_column", "Column to Edit (Cannot Be Empty)", invoiceColumnList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "edit_customer_invoice_input_new_value", "New Value")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><input type=\"submit\" value=\"Update Customer Invoice\"></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</form>")

		editCustomerInvoiceSelectCustomerID := r.FormValue("edit_customer_invoice_select_customer_id")
		editCustomerInvoiceSelectColumn := r.FormValue("edit_customer_invoice_select_column")
		editCustomerInvoiceInputNewValue := r.FormValue("edit_customer_invoice_input_new_value")

		// Validate Customer List
		_, customerInvoiceIDList := customerSlice(dbDetail)
		validateInvoiceCustomerID := slices.Contains(customerInvoiceIDList, editCustomerInvoiceSelectCustomerID)

		if editCustomerInvoiceSelectCustomerID == "" && editCustomerInvoiceSelectColumn == "" && editCustomerInvoiceInputNewValue == "" {
			// Do Nothing
		} else if validateInvoiceCustomerID == false {
			messageHTML(w, validationMessageCustomer, "warning")
		} else if editCustomerInvoiceSelectColumn == "" {
			messageHTML(w, validationMessageCustomerColumn, "warning")
		} else if editCustomerInvoiceSelectColumn == "address_line_1" || editCustomerInvoiceSelectColumn == "address_line_2" || editCustomerInvoiceSelectColumn == "city_town_village" || editCustomerInvoiceSelectColumn == "county_state_region" || editCustomerInvoiceSelectColumn == "postcode_zip_code" || editCustomerInvoiceSelectColumn == "country" {
			// Validate editCustomerInvoiceInputNewValue is a string
			validateNewValue := validateInput(editCustomerInvoiceInputNewValue, "alphaNumEmpty")
			if validateNewValue == true {
				dbDetail.connection.Query("UPDATE customer_invoice_address SET "+editCustomerInvoiceSelectColumn+" = ? WHERE id = ?;", editCustomerInvoiceInputNewValue, editCustomerInvoiceSelectCustomerID)
			} else {
				messageHTML(w, validationMessageGenericAlphaNumEmpty, "warning")
			}
		} else if editCustomerInvoiceSelectColumn == "contact_email" {
			// Validate editCustomerInvoiceInputNewValue is a email
			validateNewValue := validateInput(editCustomerInvoiceInputNewValue, "email")
			if validateNewValue == true {
				dbDetail.connection.Query("UPDATE customer_invoice_address SET "+editCustomerInvoiceSelectColumn+" = ? WHERE id = ?;", editCustomerInvoiceInputNewValue, editCustomerInvoiceSelectCustomerID)
			} else {
				messageHTML(w, validationMessageInvoiceEmail, "warning")
			}
		} else if editCustomerInvoiceSelectColumn == "contact_number" {
			// Validate editCustomerInvoiceInputNewValue is a phone number
			validateNewValue := validateInput(editCustomerInvoiceInputNewValue, "phoneNumber")
			if validateNewValue == true {
				dbDetail.connection.Query("UPDATE customer_invoice_address SET "+editCustomerInvoiceSelectColumn+" = ? WHERE id = ?;", editCustomerInvoiceInputNewValue, editCustomerInvoiceSelectCustomerID)
			} else {
				messageHTML(w, validationMessageInvoicePhoneNumber, "warning")
			}
		} else {
			messageHTML(w, validationMessageCustomerColumn, "warning")
		}
	} else {
		panic("customerEdit function should only be called with account type ID 100")
	}
}

// Customer delete function
func customerDelete(w http.ResponseWriter, r *http.Request, dbDetail databaseFunctionParameter, genDetail generalFunctionParameter) {

	// Only account type ID 100 should be able to use this function
	if genDetail.userTypeID == "100" {

		// Delete a Customer
		fmt.Fprintf(w, "<form method=\"POST\" action=\"/customer\">")
		fmt.Fprintf(w, "<table class=\"table-delete\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>Delete a Customer<br>(This Will Delete All User Accounts, PBXs and Exts Part of the Customer)</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		customerIDNameList, _ := customerSlice(dbDetail)
		selectDoubleHTML(w, "delete_customer_select_customer_id", "Customer", customerIDNameList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		confirmList := yesSlice()
		selectSingleHTML(w, "delete_customer_select_confirm", "yes to Confirm (Cannot Be Empty)", confirmList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><input type=\"submit\" value=\"Delete Customer\"></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</form>")

		deleteCustomerSelectCustomerID := r.FormValue("delete_customer_select_customer_id")
		deleteCustomerSelectConfirm := r.FormValue("delete_customer_select_confirm")

		// Validate Customer List
		_, customerIDList := customerSlice(dbDetail)
		validateCustomerID := slices.Contains(customerIDList, deleteCustomerSelectCustomerID)

		if deleteCustomerSelectCustomerID == "" && deleteCustomerSelectConfirm == "" {
			// Do Nothing
		} else if validateCustomerID == false && deleteCustomerSelectConfirm == "yes" {
			messageHTML(w, validationMessageCustomer, "warning")
		} else if validateCustomerID == true && deleteCustomerSelectConfirm != "yes" {
			messageHTML(w, validationMessageConfirmation, "warning")
		} else if deleteCustomerSelectCustomerID == "1" {
			messageHTML(w, validationMessageCustomer, "warning")
		} else if validateCustomerID == true && deleteCustomerSelectConfirm == "yes" {

			dbDetail.table = "view___customer_detail"
			dbDetail.column = "customer_id"
			dbDetail.columnWhere = "customer_id"
			dbDetail.columnWhereValue = deleteCustomerSelectCustomerID

			checkCustomerExist := selectWhere(dbDetail)

			if checkCustomerExist == "" {
				messageHTML(w, validationMessageCustomerDoesNotExist, "warning")
			} else {

				dbDetail.connection.Query(`DELETE FROM customer WHERE id = ?;`, deleteCustomerSelectCustomerID)

				checkCustomerDeleted := selectWhere(dbDetail)

				if checkCustomerDeleted == "" {
					messageHTML(w, validationMessageCustomerDeleted, "success")
				} else {
					messageHTML(w, validationMessageCustomerNotDeleted, "warning")
				}
			}

		} else {
			messageHTML(w, validationMessageInvalid, "warning")
		}
	} else {
		panic("customerDelete function should only be called with account type ID 100")
	}
}

//----------------------------------------------------------------------------------------------------

// PBX page functions
func pbxList(w http.ResponseWriter, dbDetail databaseFunctionParameter, genDetail generalFunctionParameter) {

	// Only account type ID 100, 200, 201, 300, 301, 302 should be able to use this function
	if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" || genDetail.userTypeID == "300" || genDetail.userTypeID == "301" || genDetail.userTypeID == "302" {

		var (
			pbxID                    string
			pbxName                  string
			pbxDateTimeAdded         string
			pbxSIPExtensionLimit     string
			pbxSiteAddressLine1      string
			pbxSiteAddressLine2      string
			pbxSiteCityTownVillage   string
			pbxSiteCountyStateRegion string
			pbxSitePostcodeZipCode   string
			pbxSiteCountry           string
			pbxSiteContactEmail      string
			pbxSiteContactNumber     string
			customerID               string
			customerName             string
		)

		var dbTableCountUserPBX databaseFunctionParameter
		dbTableCountUserPBX.connection = dbDetail.connection
		dbTableCountUserPBX.database = dbDetail.database
		dbTableCountUserPBX.table = "pbx"

		if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
			fmt.Fprintf(w, "<table id=\"table\" class=\"table-pbx\">")
			fmt.Fprintf(w, "  <tr>")
			fmt.Fprintf(w, "    <th>")
			fmt.Fprintf(w, "      <table id=\"table\" class=\"table-pbx\">")
			fmt.Fprintf(w, "        <tr>")
			if genDetail.userTypeID == "100" {
				fmt.Fprintf(w, "          <th>Total PBXs on the YAP Server</th>")
			} else if genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
				fmt.Fprintf(w, "          <th>Customers Total PBXs</th>")
			}
			fmt.Fprintf(w, "        </tr>")
			fmt.Fprintf(w, "        <tr>")
			if genDetail.userTypeID == "100" {
				dbTableCountUserPBX.countMinusOne = true
				fmt.Fprintf(w, "    <td>"+totalTableCount(dbTableCountUserPBX)+"</td>")
			} else if genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
				dbTableCountUserPBX.columnWhere = "customer_id"
				dbTableCountUserPBX.columnWhereValue = genDetail.userCustomerID
				fmt.Fprintf(w, "    <td>"+totalTableCountWhere(dbTableCountUserPBX)+"</td>")
			}
			fmt.Fprintf(w, "        </tr>")
			fmt.Fprintf(w, "      </table>")
			fmt.Fprintf(w, "    </th>")
			fmt.Fprintf(w, "  </tr>")
			fmt.Fprintf(w, "  <tr>")
			fmt.Fprintf(w, "    <th><button onclick=\"togglePBX() \"class=\"button-general button-pbx\">&nbsp Show/Hide PBXs &nbsp</button></th>")
			fmt.Fprintf(w, "  </tr>")
			fmt.Fprintf(w, "</table>")
		}

		if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
			fmt.Fprintf(w, "<div id=\"pbx-div\" style=\"display:none\">")
			fmt.Fprintf(w, "<br>")
		} else {
			fmt.Fprintf(w, "<div id=\"pbx-div\">")
		}
		fmt.Fprintf(w, "<table id=\"table\" class=\"table-pbx\">")
		fmt.Fprintf(w, "  <tr>")
		if genDetail.userTypeID == "100" {
			fmt.Fprintf(w, "    <th class=\"table-title\";>All PBX Contact Details on the YAP Server:</th>")
		} else {
			fmt.Fprintf(w, "    <th class=\"table-title\";>PBX Contact Details</th>")
		}
		fmt.Fprintf(w, "  </tr>")
		if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
			fmt.Fprintf(w, "  <tr>")
			fmt.Fprintf(w, "    <th>")
			fmt.Fprintf(w, "    <br>")
			var inputTableHTMLArgument jsFunctionParameter
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "pbx-contact-input-pbx-id"
			inputTableHTMLArgument.funcNameJS = "pbxContactSearchPBXID"
			inputTableHTMLArgument.placeholder = "PBX ID"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "pbx-contact-input-pbx-name"
			inputTableHTMLArgument.funcNameJS = "pbxContactSearchPBXName"
			inputTableHTMLArgument.placeholder = "PBX Name"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "pbx-contact-input-site-address"
			inputTableHTMLArgument.funcNameJS = "pbxContactSearchSiteAddress"
			inputTableHTMLArgument.placeholder = "PBX Site Address"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "pbx-contact-input-site-email"
			inputTableHTMLArgument.funcNameJS = "pbxContactSearchSiteEmail"
			inputTableHTMLArgument.placeholder = "PBX Site Email Address"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			fmt.Fprintf(w, "    <br>")
			fmt.Fprintf(w, "    <br>")
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "pbx-contact-input-site-phone"
			inputTableHTMLArgument.funcNameJS = "pbxContactSearchSitePhone"
			inputTableHTMLArgument.placeholder = "PBX Site Phone Number"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			if genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
				fmt.Fprintf(w, "    <br>")
				fmt.Fprintf(w, "    <br>")
			}
			if genDetail.userTypeID == "100" {
				inputTableHTMLArgument.inputID = "pbx-contact-input-customer-id"
				inputTableHTMLArgument.funcNameJS = "pbxContactSearchCustomerID"
				inputTableHTMLArgument.placeholder = "Customer ID"
				inputTableHTML(w, inputTableHTMLArgument)
				fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
				inputTableHTMLArgument.inputID = "pbx-contact-input-customer-name"
				inputTableHTMLArgument.funcNameJS = "pbxContactSearchCustomerName"
				inputTableHTMLArgument.placeholder = "Customer Name"
				inputTableHTML(w, inputTableHTMLArgument)
				fmt.Fprintf(w, "    <br>")
				fmt.Fprintf(w, "    <br>")
			}
			fmt.Fprintf(w, "    </th>")
			fmt.Fprintf(w, "  </tr>")
		}
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		var exportCSVButtonHTMLArgument jsFunctionParameter
		exportCSVButtonHTMLArgument.funcNameJS = "PBXContact"
		exportCSVButtonHTMLArgument.buttonCSS = "button-pbx"
		exportCSVButtonHTML(w, exportCSVButtonHTMLArgument)
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table id=\"pbx-contact-table\" class=\"table-pbx\">")
		fmt.Fprintf(w, "        <tr>")
		if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
			fmt.Fprintf(w, "          <th>PBX ID</th>")
			fmt.Fprintf(w, "          <th>PBX Name</th>")
		}
		fmt.Fprintf(w, "          <th>Site Address</th>")
		fmt.Fprintf(w, "          <th>Site Email Address</th>")
		fmt.Fprintf(w, "          <th>Site Phone Number</th>")
		if genDetail.userTypeID == "100" {
			fmt.Fprintf(w, "          <th>Customer ID</th>")
			fmt.Fprintf(w, "          <th>Customer Name</th>")
		}
		fmt.Fprintf(w, "        </tr>")

		var whereClause string
		var userWhereID string

		if genDetail.userTypeID == "100" {
			whereClause = "WHERE pbx_id != ?;"
			userWhereID = "1"
		} else if genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
			whereClause = "WHERE customer_id = ?;"
			userWhereID = genDetail.userCustomerID
		} else if genDetail.userTypeID == "300" || genDetail.userTypeID == "301" || genDetail.userTypeID == "302" {
			whereClause = "WHERE pbx_id = ?;"
			userWhereID = genDetail.userPBXID
		}

		pbxContactSQL, err := dbDetail.connection.Query(`SELECT
							pbx_id,
							pbx_name,
							pbx_site_address_line_1,
					                pbx_site_address_line_2,
					                pbx_site_city_town_village,
					                pbx_site_county_state_region,
					                pbx_site_postcode_zip_code,
					                pbx_site_country,
					                pbx_site_contact_email,
					                pbx_site_contact_number,
					                customer_id,
					                customer_name
					              FROM
					  	        yap.view___pbx_detail
						      `+whereClause, userWhereID)

		// Error
		if err != nil {
			panic(err)

		}

		for pbxContactSQL.Next() {

			err = pbxContactSQL.Scan(
				&pbxID,
				&pbxName,
				&pbxSiteAddressLine1,
				&pbxSiteAddressLine2,
				&pbxSiteCityTownVillage,
				&pbxSiteCountyStateRegion,
				&pbxSitePostcodeZipCode,
				&pbxSiteCountry,
				&pbxSiteContactEmail,
				&pbxSiteContactNumber,
				&customerID,
				&customerName,
			)

			// Error
			if err != nil {
				panic(err)
			}
			fmt.Fprintf(w, "        <tr>")
			if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
				fmt.Fprintf(w, "          <td>"+pbxID+"</td>")
				fmt.Fprintf(w, "          <td>"+pbxName+"</td>")
			}
			fmt.Fprintf(w, "          <td style=\"text-align: left;\">"+pbxSiteAddressLine1+"&nbsp<br>"+pbxSiteAddressLine2+"&nbsp<br>"+pbxSiteCityTownVillage+"&nbsp<br>"+pbxSiteCountyStateRegion+"&nbsp<br><br>"+pbxSitePostcodeZipCode+"&nbsp<br><br>"+pbxSiteCountry+"&nbsp</td>")
			fmt.Fprintf(w, "          <td><a href=\"mailto:"+pbxSiteContactEmail+"\">"+pbxSiteContactEmail+"</a></td>")
			fmt.Fprintf(w, "          <td><a href=\"tel:"+pbxSiteContactNumber+"\">"+pbxSiteContactNumber+"</a></td>")
			if genDetail.userTypeID == "100" {
				fmt.Fprintf(w, "          <td>"+customerID+"</td>")
				fmt.Fprintf(w, "          <td>"+customerName+"</td>")
			}
			fmt.Fprintf(w, "        </tr>")
		}

		fmt.Fprintf(w, "      </table>")
		if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
			var filterTableJSArgument jsFunctionParameter
			filterTableJSArgument.tableID = "pbx-contact-table"
			// Call JS filter function for PBX ID in the PBX contact table
			filterTableJSArgument.funcNameJS = "pbxContactSearchPBXID"
			filterTableJSArgument.inputID = "pbx-contact-input-pbx-id"
			filterTableJSArgument.columnNumber = 0
			filterTableJS(w, filterTableJSArgument)
			// Call JS filter function for PBX name in the PBX contact table
			filterTableJSArgument.funcNameJS = "pbxContactSearchPBXName"
			filterTableJSArgument.inputID = "pbx-contact-input-pbx-name"
			filterTableJSArgument.columnNumber = 1
			filterTableJS(w, filterTableJSArgument)
			// Call JS filter function for site address in the PBX contact table
			filterTableJSArgument.funcNameJS = "pbxContactSearchSiteAddress"
			filterTableJSArgument.inputID = "pbx-contact-input-site-address"
			filterTableJSArgument.columnNumber = 2
			filterTableJS(w, filterTableJSArgument)
			// Call JS filter function for site email in the PBX contact table
			filterTableJSArgument.funcNameJS = "pbxContactSearchSiteEmail"
			filterTableJSArgument.inputID = "pbx-contact-input-site-email"
			filterTableJSArgument.columnNumber = 3
			filterTableJS(w, filterTableJSArgument)
			// Call JS filter function for site phone in the PBX contact table
			filterTableJSArgument.funcNameJS = "pbxContactSearchSitePhone"
			filterTableJSArgument.inputID = "pbx-contact-input-site-phone"
			filterTableJSArgument.columnNumber = 4
			filterTableJS(w, filterTableJSArgument)
			if genDetail.userTypeID == "100" {
				// Call JS filter function for the customer ID in the PBX contact table
				filterTableJSArgument.funcNameJS = "pbxContactSearchCustomerID"
				filterTableJSArgument.inputID = "pbx-contact-input-customer-id"
				filterTableJSArgument.columnNumber = 5
				filterTableJS(w, filterTableJSArgument)
				// Call JS filter function for the customer name in the PBX contact table
				filterTableJSArgument.funcNameJS = "pbxContactSearchCustomerName"
				filterTableJSArgument.inputID = "pbx-contact-input-customer-name"
				filterTableJSArgument.columnNumber = 6
				filterTableJS(w, filterTableJSArgument)
			}
		}
		var exportCSVJSArgument jsFunctionParameter
		exportCSVJSArgument.funcNameJS = "PBXContact"
		exportCSVJSArgument.tableID = "pbx-contact-table"
		exportCSVJSArgument.fileName = "YAP_pbx_contact_details"
		exportCSVJSArgument.pathURL = "pbx"
		exportCSVJS(w, exportCSVJSArgument)
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")

		// Customer resource table
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<table id=\"table\" class=\"table-pbx\">")
		fmt.Fprintf(w, "  <tr>")
		if genDetail.userTypeID == "100" {
			fmt.Fprintf(w, "    <th class=\"table-title\";>All PBX Resource Limits on the YAP Server:</th>")
		} else {
			fmt.Fprintf(w, "    <th class=\"table-title\";>PBX Resource Limits</th>")
		}
		fmt.Fprintf(w, "  </tr>")
		if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
			fmt.Fprintf(w, "  <tr>")
			fmt.Fprintf(w, "    <th>")
			fmt.Fprintf(w, "    <br>")
			var inputTableHTMLArgument jsFunctionParameter
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "pbx-resource-input-pbx-id"
			inputTableHTMLArgument.funcNameJS = "pbxResourceSearchPBXID"
			inputTableHTMLArgument.placeholder = "PBX ID"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "pbx-resource-input-pbx-name"
			inputTableHTMLArgument.funcNameJS = "pbxResourceSearchPBXName"
			inputTableHTMLArgument.placeholder = "PBX Name"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "pbx-resource-input-date-time"
			inputTableHTMLArgument.funcNameJS = "pbxResourceSearchDateTime"
			inputTableHTMLArgument.placeholder = "Date & Time Created"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			if genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
				fmt.Fprintf(w, "    <br>")
				fmt.Fprintf(w, "    <br>")
			}
			if genDetail.userTypeID == "100" {
				inputTableHTMLArgument.inputID = "pbx-resource-input-customer-id"
				inputTableHTMLArgument.funcNameJS = "pbxResourceSearchCustomerID"
				inputTableHTMLArgument.placeholder = "Customer ID"
				inputTableHTML(w, inputTableHTMLArgument)
				fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
				fmt.Fprintf(w, "    <br>")
				fmt.Fprintf(w, "    <br>")
				fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
				inputTableHTMLArgument.inputID = "pbx-resource-input-customer-name"
				inputTableHTMLArgument.funcNameJS = "pbxResourceSearchCustomerName"
				inputTableHTMLArgument.placeholder = "Customer Name"
				inputTableHTML(w, inputTableHTMLArgument)
				fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
				fmt.Fprintf(w, "    <br>")
				fmt.Fprintf(w, "    <br>")
			}
			fmt.Fprintf(w, "    </th>")
			fmt.Fprintf(w, "  </tr>")
		}
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		exportCSVButtonHTMLArgument.funcNameJS = "PBXResource"
		exportCSVButtonHTMLArgument.buttonCSS = "button-pbx"
		exportCSVButtonHTML(w, exportCSVButtonHTMLArgument)
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table id=\"pbx-resource-table\" class=\"table-pbx\">")
		fmt.Fprintf(w, "        <tr>")
		if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
			fmt.Fprintf(w, "          <th>PBX ID</th>")
			fmt.Fprintf(w, "          <th>PBX Name</th>")
		}
		fmt.Fprintf(w, "          <th>Date & Time Added</th>")
		fmt.Fprintf(w, "          <th>Ext Limit for PBX</th>")
		if genDetail.userTypeID == "100" {
			fmt.Fprintf(w, "          <th>Customer ID</th>")
			fmt.Fprintf(w, "          <th>Customer Name</th>")
		}
		fmt.Fprintf(w, "        </tr>")

		pbxResourceSQL, err := dbDetail.connection.Query(`SELECT
							pbx_id,
							pbx_name,
							pbx_date_time_added,
							pbx_sip_extension_limit,
							customer_id,
							customer_name
					              FROM
					  	        yap.view___pbx_detail
						      `+whereClause, userWhereID)

		// Error
		if err != nil {
			panic(err)

		}

		for pbxResourceSQL.Next() {

			err = pbxResourceSQL.Scan(
				&pbxID,
				&pbxName,
				&pbxDateTimeAdded,
				&pbxSIPExtensionLimit,
				&customerID,
				&customerName,
			)

			// Error
			if err != nil {
				panic(err)
			}
			fmt.Fprintf(w, "        <tr>")
			if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
				fmt.Fprintf(w, "          <td>"+pbxID+"</td>")
				fmt.Fprintf(w, "          <td>"+pbxName+"</td>")
			}
			fmt.Fprintf(w, "          <td>"+formatDateTime(pbxDateTimeAdded)+"</td>")
			fmt.Fprintf(w, "          <td>"+pbxSIPExtensionLimit+"</td>")
			if genDetail.userTypeID == "100" {
				fmt.Fprintf(w, "          <td>"+customerID+"</td>")
				fmt.Fprintf(w, "          <td>"+customerName+"</td>")
			}
			fmt.Fprintf(w, "        </tr>")
		}

		fmt.Fprintf(w, "      </table>")
		if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
			var filterTableJSArgument jsFunctionParameter
			filterTableJSArgument.tableID = "pbx-resource-table"
			// Call JS filter function for PBX ID in the PBX resource table
			filterTableJSArgument.funcNameJS = "pbxResourceSearchPBXID"
			filterTableJSArgument.inputID = "pbx-resource-input-pbx-id"
			filterTableJSArgument.columnNumber = 0
			filterTableJS(w, filterTableJSArgument)
			// Call JS filter function for PBX name in the PBX resource table
			filterTableJSArgument.funcNameJS = "pbxResourceSearchPBXName"
			filterTableJSArgument.inputID = "pbx-resource-input-pbx-name"
			filterTableJSArgument.columnNumber = 1
			filterTableJS(w, filterTableJSArgument)
			// Call JS filter function for date and time in the PBX resource table
			filterTableJSArgument.funcNameJS = "pbxResourceSearchDateTime"
			filterTableJSArgument.inputID = "pbx-resource-input-date-time"
			filterTableJSArgument.columnNumber = 2
			filterTableJS(w, filterTableJSArgument)
			if genDetail.userTypeID == "100" {
				// Call JS filter function for the customer ID in the PBX resource table
				filterTableJSArgument.funcNameJS = "pbxResourceSearchCustomerID"
				filterTableJSArgument.inputID = "pbx-resource-input-customer-id"
				filterTableJSArgument.columnNumber = 4
				filterTableJS(w, filterTableJSArgument)
				// Call JS filter function for the customer name in the PBX resource table
				filterTableJSArgument.funcNameJS = "pbxResourceSearchCustomerName"
				filterTableJSArgument.inputID = "pbx-resource-input-customer-name"
				filterTableJSArgument.columnNumber = 5
				filterTableJS(w, filterTableJSArgument)
			}
		}
		exportCSVJSArgument.funcNameJS = "PBXResource"
		exportCSVJSArgument.tableID = "pbx-resource-table"
		exportCSVJSArgument.fileName = "YAP_pbx_resource_details"
		exportCSVJSArgument.pathURL = "pbx"
		exportCSVJS(w, exportCSVJSArgument)
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</div>")
		if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
			var toggleDivJSArgument jsFunctionParameter
			toggleDivJSArgument.funcNameJS = "togglePBX"
			toggleDivJSArgument.divID = "pbx-div"
			toggleDivJS(w, toggleDivJSArgument)
		}
	} else {
		panic("pbxList function should only be called with account type ID 100, 200, 201, 300, 301, 302")
	}
}

// Add PBX function
func pbxAdd(w http.ResponseWriter, r *http.Request, dbDetail databaseFunctionParameter, genDetail generalFunctionParameter) {

	// Only account type ID 100, 200, 201 should be able to use this function
	if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {

		fmt.Fprintf(w, "<form method=\"POST\" action=\"/pbx\">")
		fmt.Fprintf(w, "<table class=\"table-add\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>Add a New PBX</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		if genDetail.userTypeID == "100" {
			customerIDNameList, _ := customerSlice(dbDetail)
			selectDoubleHTML(w, "add_pbx_select_customer_id", "Customer", customerIDNameList)
		} else {
			// Do Nothing
		}
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_pbx_input_pbx_name", "PBX Name (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		extLimitList := extLimitSlice()
		if genDetail.userTypeID == "100" {
			selectSingleHTML(w, "add_pbx_select_ext_limit", "Ext Limt (Cannot Be Empty)", extLimitList)
		} else {
			inputReadOnlyHTML(w, "add_pbx_input_default_ext_limit", "Default Ext Limit for PBX", genDetail.defaultExtLimit)
		}
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		if genDetail.userTypeID == "100" {
			inputReadOnlyHTML(w, "add_pbx_input_default_ext_limit", "Default Ext Limit for PBX", genDetail.defaultExtLimit)
		} else {
			// Do Nothing
		}
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td style=\"border: none;\">")
		fmt.Fprintf(w, "            <br>")
		fmt.Fprintf(w, "            <br>")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_pbx_input_site_address_line_1", "Site Address Line One")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_pbx_input_site_address_line_2", "Site Address Line Two")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_pbx_input_site_city_town_village", "Site City/Town/Village")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_pbx_input_site_county_state_region", "Site County/State/Region")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_pbx_input_site_postcode_zip_code", "Site Postcode/Zip Code")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_pbx_input_site_country", "Site Country")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_pbx_input_site_contact_email", "Site Email (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_pbx_input_site_contact_number", "Site Phone (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><input type=\"submit\" value=\"Create PBX\"></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</form>")

		addPBXSelectCustomerID := r.FormValue("add_pbx_select_customer_id")
		addPBXInputPBXName := r.FormValue("add_pbx_input_pbx_name")
		addPBXSelectExtLimit := r.FormValue("add_pbx_select_ext_limit")

		addPBXInputSiteAddressLine1 := r.FormValue("add_pbx_input_site_address_line_1")
		addPBXInputSiteAddressLine2 := r.FormValue("add_pbx_input_site_address_line_2")
		addPBXInputSiteCityTownVillage := r.FormValue("add_pbx_input_site_city_town_village")
		addPBXInputSiteCountyStateRegion := r.FormValue("add_pbx_input_site_county_state_region")
		addPBXInputSitePostcodeZipCode := r.FormValue("add_pbx_input_site_postcode_zip_code")
		addPBXInputSiteCountry := r.FormValue("add_pbx_input_site_country")
		addPBXInputSiteContactEmail := r.FormValue("add_pbx_input_site_contact_email")
		addPBXInputSiteContactNumber := r.FormValue("add_pbx_input_site_contact_number")

		if genDetail.userTypeID != "100" {
			addPBXSelectExtLimit = genDetail.defaultExtLimit
		}

		// Check customer ID is contained in the slice
		_, customerIDList := customerSlice(dbDetail)
		customerIDList = append(customerIDList, "")
		validateCustomerID := slices.Contains(customerIDList, addPBXSelectCustomerID)

		// Validate the PBX name
		validatePBXName := validateInput(addPBXInputPBXName, "alphaNum")

		// Validate ext limit
		validateExtLimit := slices.Contains(extLimitList, addPBXSelectExtLimit)

		// Validate site address line one
		validateSiteAddressLine1 := validateInput(addPBXInputSiteAddressLine1, "alphaNumEmpty")

		// Validate site address line two
		validateSiteAddressLine2 := validateInput(addPBXInputSiteAddressLine2, "alphaNumEmpty")

		// Validate site city/town/village
		validateSiteCityTownVillage := validateInput(addPBXInputSiteCityTownVillage, "alphaNumEmpty")

		// Validate site county/state/region
		validateSiteCountyStateRegion := validateInput(addPBXInputSiteCountyStateRegion, "alphaNumEmpty")

		// Validate site postcode/zip code
		validateSitePostcodeZipCode := validateInput(addPBXInputSitePostcodeZipCode, "alphaNumEmpty")

		// Validate site country
		validateSiteCountry := validateInput(addPBXInputSiteCountry, "alphaNumEmpty")

		// Validate site contact emial
		validateSiteContactEmail := validateInput(addPBXInputSiteContactEmail, "email")

		// Validate Site contact phone number
		validateSiteContactNumber := validateInput(addPBXInputSiteContactNumber, "phoneNumber")

		if genDetail.userTypeID != "100" {
			addPBXSelectCustomerID = genDetail.userCustomerID
		}

		if addPBXSelectCustomerID == "" && addPBXInputPBXName == "" && addPBXSelectExtLimit == "" {
			// Do Nothing
		} else if genDetail.userTypeID != "100" && addPBXInputPBXName == "" {
			// Do Nothing
		} else if validateCustomerID == false || addPBXSelectCustomerID == "" {
			messageHTML(w, validationMessageCustomer, "warning")
		} else if validatePBXName == false {
			messageHTML(w, validationMessagePBXName, "warning")
		} else if validateExtLimit == false || addPBXSelectExtLimit == "" {
			messageHTML(w, validationMessagePBXExtLimit, "warning")
		} else if validateSiteAddressLine1 == false {
			messageHTML(w, validationMessageAddresslineOne, "warning")
		} else if validateSiteAddressLine2 == false {
			messageHTML(w, validationMessageAddresslineTwo, "warning")
		} else if validateSiteCityTownVillage == false {
			messageHTML(w, validationMessageCityTownVillage, "warning")
		} else if validateSiteCountyStateRegion == false {
			messageHTML(w, validationMessageCountyStateRegion, "warning")
		} else if validateSitePostcodeZipCode == false {
			messageHTML(w, validationMessagePostcodeZipCode, "warning")
		} else if validateSiteCountry == false {
			messageHTML(w, validationMessageCountry, "warning")
		} else if validateSiteContactEmail == false {
			messageHTML(w, validationMessagePBXSiteEmail, "warning")
		} else if validateSiteContactNumber == false {
			messageHTML(w, validationMessagePBXSitePhoneNumber, "warning")
		} else {

			// Used to compare the max allowed PBXs to the number of PBXs that already exist
			var pbxMaxLimit string
			dbDetail.table = "view___customer_detail"
			dbDetail.column = "customer_pbx_limit"
			dbDetail.columnWhere = "customer_id"
			dbDetail.columnWhereValue = addPBXSelectCustomerID
			pbxMaxLimit = selectWhere(dbDetail)

			var pbxCount string
			dbDetail.table = "view___pbx_detail"
			dbDetail.column = "customer_id"
			dbDetail.columnWhere = "customer_id"
			dbDetail.countMinusOne = false
			dbDetail.columnWhereValue = addPBXSelectCustomerID
			pbxCount = totalTableCountWhere(dbDetail)

			if pbxCount >= pbxMaxLimit {
				messageHTML(w, validationMessagePBXMaxPBX, "warning")
			} else {

				pbxID := genID()

				dbDetail.connection.Query(`INSERT 
					     INTO
					   pbx (
					     id,
                                             name,
                                             customer_id,
                                             sip_extension_limit
	                                   )
                                           VALUES(?, ?, ?, ?);`,
					pbxID,
					addPBXInputPBXName,
					addPBXSelectCustomerID,
					addPBXSelectExtLimit)

				dbDetail.connection.Query(`INSERT 
                         	             INTO
                                           pbx_site_address (
                                           id,
                                           address_line_1,
                                           address_line_2,
                                           city_town_village,
                                           county_state_region,
                                           postcode_zip_code,
                                           country,
                                           contact_email,
                                           contact_number 
                                           )
                                           VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?);`,
					pbxID,
					nullSQL(addPBXInputSiteAddressLine1),
					nullSQL(addPBXInputSiteAddressLine2),
					nullSQL(addPBXInputSiteCityTownVillage),
					nullSQL(addPBXInputSiteCountyStateRegion),
					nullSQL(addPBXInputSitePostcodeZipCode),
					nullSQL(addPBXInputSiteCountry),
					nullSQL(addPBXInputSiteContactEmail),
					nullSQL(addPBXInputSiteContactNumber))

				dbDetail.table = "view___pbx_detail"
				dbDetail.column = "pbx_id"
				dbDetail.columnWhere = "pbx_id"
				dbDetail.columnWhereValue = pbxID

				checkPBXCreated := selectWhere(dbDetail)

				if checkPBXCreated == pbxID {
					messageHTML(w, validationMessagePBXCreated, "success")

					dbDetail.table = "view___customer_detail"
					dbDetail.columnWhere = "customer_id"
					dbDetail.columnWhereValue = addPBXSelectCustomerID

					// Get PBX setup price
					dbDetail.column = "customer_pbx_setup_price"
					setupPrice := selectWhere(dbDetail)

					// Get PBX sales tax rate
					dbDetail.column = "customer_pbx_sales_tax_rate"
					salesTaxRate := selectWhere(dbDetail)

					// Get PBX sales tax status
					dbDetail.column = "customer_pbx_sales_tax_status"
					salesTaxStatus := selectWhere(dbDetail)

					// Get PBX contract length
					dbDetail.column = "customer_pbx_contract_length"
					contractLength := selectWhere(dbDetail)

					var invoicePBXExt invoicePBXExtFunctionParameter

					invoicePBXExt.customerID = addPBXSelectCustomerID
					invoicePBXExt.serviceProduct = "⊛ YAP PBX Setup ⊛"
					invoicePBXExt.tag = pbxID
					invoicePBXExt.pbxID = pbxID
					invoicePBXExt.sellPrice = setupPrice
					invoicePBXExt.salesTaxRate = salesTaxRate
					invoicePBXExt.salesTaxStatus = salesTaxStatus
					invoicePBXExt.billItemOnce = "yes"
					invoicePBXExt.itemOnHold = "no"
					invoicePBXExt.contractLength = contractLength
					invoicePBXExt.contractStartDate = currentDate()

					// Add PBX setup to invoice
					invoicePBXExtAdd(dbDetail, invoicePBXExt)

					// Get PBX setup price
					dbDetail.column = "customer_pbx_rental_price"
					rentalPrice := selectWhere(dbDetail)

					invoicePBXExt.customerID = addPBXSelectCustomerID
					invoicePBXExt.serviceProduct = "⊛ YAP PBX Rental ⊛"
					invoicePBXExt.tag = pbxID
					invoicePBXExt.pbxID = pbxID
					invoicePBXExt.sellPrice = rentalPrice
					invoicePBXExt.salesTaxRate = salesTaxRate
					invoicePBXExt.salesTaxStatus = salesTaxStatus
					invoicePBXExt.billItemOnce = "no"
					invoicePBXExt.itemOnHold = "yes"
					invoicePBXExt.contractLength = contractLength
					invoicePBXExt.contractStartDate = currentDate()

					// Add PBX rental to invoice
					invoicePBXExtAdd(dbDetail, invoicePBXExt)

				} else {
					messageHTML(w, validationMessagePBXNotCreated, "warning")
				}
			}
		}
	} else {
		panic("pbxAdd function should only be called with account type ID 100, 200, 201")
	}
}

// PBX edit function
func pbxEdit(w http.ResponseWriter, r *http.Request, dbDetail databaseFunctionParameter, genDetail generalFunctionParameter) {

	// Only account type ID 100, 200, 201 should be able to use this function
	if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {

		// List of all column names from the PBX table
		pbxColumnList := [][]string{
			{"name", "PBX Name"},
		}

		// List of all column names from the PBX table
		pbxColumnExtraList := [][]string{
			{"sip_extension_limit", "Ext Limit"},
		}

		pbxColumnExtraList = append(pbxColumnList, pbxColumnExtraList...)

		fmt.Fprintf(w, "<form method=\"POST\" action=\"/pbx\">")
		fmt.Fprintf(w, "<table class=\"table-pbx\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>Edit PBX Details</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td style=\"text-align: left;\">")
		fmt.Fprintf(w, "      <b><u>Acceptable Values for Columns</u></b><br><br>")
		fmt.Fprintf(w, "      <b>PBX Name:</b> text<br>")
		if genDetail.userTypeID == "100" {
			fmt.Fprintf(w, "      <b>Ext Limit:</b> 1, 2, 3, 4, 5, 10, 25, 50, 75, 100, 150, 200, 250, 500, 750, 1000, 1500, 2000, 2500, 5000<br>")
		}
		fmt.Fprintf(w, "    </td>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		pbxIDNameList, _ := pbxSlice(dbDetail)
		dbDetail.columnWhere = "customer_id"
		dbDetail.columnWhereValue = genDetail.userCustomerID
		pbxWhereIDNameList, _ := pbxWhereSlice(dbDetail)
		if genDetail.userTypeID == "100" {
			selectDoubleHTML(w, "edit_pbx_select_pbx_id", "PBX", pbxIDNameList)
		} else if genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
			selectDoubleHTML(w, "edit_pbx_select_pbx_id", "PBX", pbxWhereIDNameList)
		}
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		if genDetail.userTypeID == "100" {
			selectDoubleHiddenHTML(w, "edit_pbx_select_column", "Column to Edit", pbxColumnExtraList)
		} else if genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
			selectDoubleHiddenHTML(w, "edit_pbx_select_column", "Column to Edit", pbxColumnList)
		}
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "edit_pbx_input_new_value", "New Value")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><input type=\"submit\" value=\"Update PBX\"></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</form>")

		editPBXSelectPBXID := r.FormValue("edit_pbx_select_pbx_id")
		editPBXSelectColumn := r.FormValue("edit_pbx_select_column")
		editPBXInputNewValue := r.FormValue("edit_pbx_input_new_value")

		// Validate the PBX ID
		_, pbxIDList := pbxSlice(dbDetail)
		validatePBXID := slices.Contains(pbxIDList, editPBXSelectPBXID)
		_, pbxWhereIDList := pbxWhereSlice(dbDetail)
		validatePBXWhereID := slices.Contains(pbxWhereIDList, editPBXSelectPBXID)

		if editPBXSelectPBXID == "" && editPBXSelectColumn == "" && editPBXInputNewValue == "" {
			// Do Nothing
		} else if genDetail.userTypeID == "100" && validatePBXID == false {
			messageHTML(w, validationMessagePBX, "warning")
		} else if genDetail.userTypeID == "200" && validatePBXWhereID == false {
			messageHTML(w, validationMessagePBX, "warning")
		} else if genDetail.userTypeID == "201" && validatePBXWhereID == false {
			messageHTML(w, validationMessagePBX, "warning")
		} else if editPBXSelectColumn == "" {
			messageHTML(w, validationMessagePBXColumn, "warning")
		} else if editPBXSelectColumn == "name" {
			// Validate editPBXInputNewValue is a string
			validateNewValue := validateInput(editPBXInputNewValue, "alphaNumEmpty")
			if genDetail.userTypeID == "100" {
				if validateNewValue == true {
					dbDetail.connection.Query("UPDATE pbx SET "+editPBXSelectColumn+" = ? WHERE id = ?;", editPBXInputNewValue, editPBXSelectPBXID)
				} else {
					messageHTML(w, validationMessageGenericAlphaNumEmpty, "warning")
				}
			} else if genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
				if validateNewValue == true {
					dbDetail.connection.Query("UPDATE pbx SET "+editPBXSelectColumn+" = ? WHERE id = ? AND customer_id = ?;", editPBXInputNewValue, editPBXSelectPBXID, genDetail.userCustomerID)
				} else {
					messageHTML(w, validationMessageGenericAlphaNumEmpty, "warning")
				}
			} else {
				messageHTML(w, validationMessagePBXColumn, "warning")
			}
		} else if genDetail.userTypeID == "100" && editPBXSelectColumn == "sip_extension_limit" {
			extLimitList := extLimitSlice()
			// Validate editCustomerSelectColumn is in the extLimitList Slice
			validateNewValue := slices.Contains(extLimitList, editPBXInputNewValue)
			if validateNewValue == true {
				dbDetail.connection.Query("UPDATE pbx SET "+editPBXSelectColumn+" = ? WHERE id = ?;", editPBXInputNewValue, editPBXSelectPBXID)
			} else {
				messageHTML(w, validationMessagePBXExtLimit, "warning")
			}
		} else {
			messageHTML(w, validationMessagePBXColumn, "warning")
		}

		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<form method=\"POST\" action=\"/pbx\">")
		fmt.Fprintf(w, "<table class=\"table-pbx\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>Edit PBX Site Details</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td style=\"text-align: left;\">")
		fmt.Fprintf(w, "      <b><u>Acceptable Values for Columns</u></b><br><br>")
		fmt.Fprintf(w, "      <b>Site Address Line One:</b> text<br>")
		fmt.Fprintf(w, "      <b>Site Address Line Two:</b> text<br>")
		fmt.Fprintf(w, "      <b>Site City Town Village:</b> text<br>")
		fmt.Fprintf(w, "      <b>Site County State Region:</b> text<br>")
		fmt.Fprintf(w, "      <b>Site Postcode Zip Code:</b> text<br>")
		fmt.Fprintf(w, "      <b>Site Country:</b> text<br>")
		fmt.Fprintf(w, "      <b>Site Contact Email:</b> valid email address<br>")
		fmt.Fprintf(w, "      <b>Site Contact Number:</b> phone number in e.164 format<br>")
		fmt.Fprintf(w, "    </td>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		if genDetail.userTypeID == "100" {
			selectDoubleHTML(w, "edit_pbx_site_select_pbx_id", "PBX", pbxIDNameList)
		} else if genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
			selectDoubleHTML(w, "edit_pbx_site_select_pbx_id", "PBX", pbxWhereIDNameList)
		}
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		siteColumnList := siteColumnSlice()
		selectDoubleHiddenHTML(w, "edit_pbx_site_select_column", "Column to Edit (Cannot Be Empty)", siteColumnList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "edit_pbx_site_input_new_value", "New Value")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><input type=\"submit\" value=\"Update PBX Site\"></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</form>")

		editPBXSiteSelectPBXID := r.FormValue("edit_pbx_site_select_pbx_id")
		editPBXSiteSelectColumn := r.FormValue("edit_pbx_site_select_column")
		editPBXSiteInputNewValue := r.FormValue("edit_pbx_site_input_new_value")

		// Validate the PBX ID
		_, pbxSiteIDList := pbxSlice(dbDetail)
		validatePBXSiteID := slices.Contains(pbxSiteIDList, editPBXSiteSelectPBXID)
		_, pbxSiteWhereIDList := pbxWhereSlice(dbDetail)
		validatePBXSiteWhereID := slices.Contains(pbxSiteWhereIDList, editPBXSiteSelectPBXID)

		if editPBXSiteSelectPBXID == "" && editPBXSiteSelectColumn == "" && editPBXSiteInputNewValue == "" {
			// Do Nothing
		} else if genDetail.userTypeID == "100" && validatePBXSiteID == false {
			messageHTML(w, validationMessagePBX, "warning")
		} else if genDetail.userTypeID == "200" && validatePBXSiteWhereID == false {
			messageHTML(w, validationMessagePBX, "warning")
		} else if genDetail.userTypeID == "201" && validatePBXSiteWhereID == false {
			messageHTML(w, validationMessagePBX, "warning")
		} else if editPBXSiteSelectColumn == "" {
			messageHTML(w, validationMessagePBXColumn, "warning")
		} else if editPBXSiteSelectColumn == "address_line_1" || editPBXSiteSelectColumn == "address_line_2" || editPBXSiteSelectColumn == "city_town_village" || editPBXSiteSelectColumn == "county_state_region" || editPBXSiteSelectColumn == "postcode_zip_code" || editPBXSiteSelectColumn == "country" {
			// Validate editPBXSiteInputNewValue is a string
			validateNewValue := validateInput(editPBXSiteInputNewValue, "alphaNumEmpty")
			if validateNewValue == true {
				dbDetail.connection.Query("UPDATE pbx_site_address SET "+editPBXSiteSelectColumn+" = ? WHERE id = ?;", editPBXSiteInputNewValue, editPBXSiteSelectPBXID)
			} else {
				messageHTML(w, validationMessageCustomerColumn, "warning")
			}
		} else if editPBXSiteSelectColumn == "contact_email" {
			// Validate editPBXSiteInputNewValue is a email
			validateNewValue := validateInput(editPBXSiteInputNewValue, "email")
			if validateNewValue == true {
				dbDetail.connection.Query("UPDATE pbx_site_address SET "+editPBXSiteSelectColumn+" = ? WHERE id = ?;", editPBXSiteInputNewValue, editPBXSiteSelectPBXID)
			} else {
				messageHTML(w, validationMessagePBXSiteEmail, "warning")
			}
		} else if editPBXSiteSelectColumn == "contact_number" {
			// Validate editPBXSiteInputNewValue is a phone number
			validateNewValue := validateInput(editPBXSiteInputNewValue, "phoneNumber")
			if validateNewValue == true {
				dbDetail.connection.Query("UPDATE pbx_site_address SET "+editPBXSiteSelectColumn+" = ? WHERE id = ?;", editPBXSiteInputNewValue, editPBXSiteSelectPBXID)
			} else {
				messageHTML(w, validationMessagePBXSitePhoneNumber, "warning")
			}
		} else {
			messageHTML(w, validationMessagePBXColumn, "warning")
		}
	} else {
		panic("pbxEdit function should only be called with account type ID 100, 200, 201")
	}
}

// PBX delete function
func pbxDelete(w http.ResponseWriter, r *http.Request, dbDetail databaseFunctionParameter, genDetail generalFunctionParameter) {

	// Only account type ID 100, 200 should be able to use this function
	if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" {

		// Delete a Customer
		fmt.Fprintf(w, "<form method=\"POST\" action=\"/pbx\">")
		fmt.Fprintf(w, "<table class=\"table-delete\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>Delete a PBX<br>(This Will Delete All User Accounts and Exts Part of the Customer)</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		pbxIDNameList, _ := pbxSlice(dbDetail)
		dbDetail.columnWhere = "customer_id"
		dbDetail.columnWhereValue = genDetail.userCustomerID
		pbxWhereIDNameList, _ := pbxWhereSlice(dbDetail)
		if genDetail.userTypeID == "100" {
			selectDoubleHTML(w, "delete_pbx_select_pbx_id", "PBX", pbxIDNameList)
		} else if genDetail.userTypeID == "200" {
			selectDoubleHTML(w, "delete_pbx_select_pbx_id", "PBX", pbxWhereIDNameList)
		}
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		confirmList := yesSlice()
		selectSingleHTML(w, "delete_pbx_select_confirm", "yes to Confirm (Cannot Be Empty)", confirmList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><input type=\"submit\" value=\"Delete PBX\"></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</form>")

		deletePBXSelectPBXID := r.FormValue("delete_pbx_select_pbx_id")
		deletePBXSelectConfirm := r.FormValue("delete_pbx_select_confirm")

		// Validate the PBX ID
		_, pbxIDList := pbxSlice(dbDetail)
		validatePBXID := slices.Contains(pbxIDList, deletePBXSelectPBXID)
		_, pbxWhereIDList := pbxWhereSlice(dbDetail)
		validatePBXWhereID := slices.Contains(pbxWhereIDList, deletePBXSelectPBXID)

		// Check PBX Exist
		dbDetail.table = "view___pbx_detail"
		dbDetail.column = "pbx_id"
		dbDetail.columnWhere = "pbx_id"

		// Variables for adding cease charge to invoice; they are safe because they are not the input from the user
		var invoicePBXExt invoicePBXExtFunctionParameter
		invoicePBXExt.serviceProduct = "⊛ YAP PBX Cease ⊛"
		invoicePBXExt.billItemOnce = "yes"
		invoicePBXExt.itemOnHold = "no"
		invoicePBXExt.contractLength = ""
		invoicePBXExt.contractStartDate = ""

		if deletePBXSelectPBXID == "" && deletePBXSelectConfirm == "" {
			// Do Nothing
		} else if genDetail.userTypeID == "100" && validatePBXID == false {
			messageHTML(w, validationMessagePBX, "warning")
		} else if genDetail.userTypeID == "200" && validatePBXWhereID == false {
			messageHTML(w, validationMessagePBX, "warning")
		} else if validatePBXID == false && deletePBXSelectConfirm == "yes" {
			messageHTML(w, validationMessagePBX, "warning")
		} else if validatePBXID == true && deletePBXSelectConfirm != "yes" {
			messageHTML(w, validationMessageConfirmation, "warning")
		} else if deletePBXSelectPBXID == "1" {
			messageHTML(w, validationMessagePBXIDOne, "warning")
		} else if genDetail.userTypeID == "100" && validatePBXID == true && deletePBXSelectConfirm == "yes" {

			dbDetail.columnWhereValue = deletePBXSelectPBXID
			checkPBXExist := selectWhere(dbDetail)

			if checkPBXExist == "" {
				messageHTML(w, validationMessagePBXDoesNotExist, "warning")
			} else {
				// Get customer ID based on the PBX ID
				dbDetail.column = "customer_id"
				userCustomerID := selectWhere(dbDetail)

				// Check how many exts the PBX has
				dbDetail.table = "view___sip_extension_detail"
				dbDetail.countMinusOne = false
				dbDetail.columnWhere = "pbx_id"
				dbDetail.columnWhereValue = deletePBXSelectPBXID
				extTotal := totalTableCountWhere(dbDetail)

				dbDetail.connection.Query(`DELETE FROM pbx WHERE id = ?;`, deletePBXSelectPBXID)

				// Check PBX deleted
				dbDetail.column = "pbx_id"
				checkPBXDeleted := selectWhere(dbDetail)

				if checkPBXDeleted == "" {
					messageHTML(w, validationMessagePBXDeleted, "success")

					// Variables for customer table columns
					dbDetail.table = "view___pbx_detail"
					dbDetail.table = "view___customer_detail"
					dbDetail.column = "pbx_id"
					dbDetail.columnWhere = "customer_id"
					dbDetail.columnWhereValue = userCustomerID

					// Add cease charge for exts deleted
					extTotalFloat64 := stringToFloat64(extTotal)

					if extTotalFloat64 >= 1 {

						// Get ext cease price
						dbDetail.column = "customer_ext_cease_price"
						extCeasePrice := selectWhere(dbDetail)

						// Get ext sales tax rate
						dbDetail.column = "customer_ext_sales_tax_rate"
						salesTaxRate := selectWhere(dbDetail)

						// Get ext sales tax status
						dbDetail.column = "customer_ext_sales_tax_status"
						salesTaxStatus := selectWhere(dbDetail)

						invoicePBXExt.serviceProduct = "⊛ YAP Extension Cease ⊛"
						invoicePBXExt.customerID = userCustomerID
						invoicePBXExt.pbxID = deletePBXSelectPBXID
						invoicePBXExt.tag = "Ext Cease X" + extTotal + ": " + deletePBXSelectPBXID
						invoicePBXExt.contractStartDate = currentDate()
						extCeasePriceFloat64 := stringToFloat64(extCeasePrice)
						extTotalExtCeasePrice := extTotalFloat64 * extCeasePriceFloat64
						invoicePBXExt.sellPrice = strconv.FormatFloat(extTotalExtCeasePrice, 'f', -1, 64)
						invoicePBXExt.salesTaxRate = salesTaxRate
						invoicePBXExt.salesTaxStatus = salesTaxStatus

						// Add PBX cease to invoice
						invoicePBXExtAdd(dbDetail, invoicePBXExt)

						// Delete PBX rental record from invoice_item table
						dbDetail.connection.Query(`DELETE FROM invoice_item WHERE pbx_id = ? AND service_product_name = ?;`, deletePBXSelectPBXID, "⊛ YAP Extension Rental ⊛")

					}

					// Get PBX cease price
					dbDetail.column = "customer_pbx_cease_price"
					pbxCeasePrice := selectWhere(dbDetail)

					// Get PBX sales tax rate
					dbDetail.column = "customer_pbx_sales_tax_rate"
					salesTaxRate := selectWhere(dbDetail)

					// Get PBX sales tax status
					dbDetail.column = "customer_pbx_sales_tax_status"
					salesTaxStatus := selectWhere(dbDetail)

					// These variables have to be here because they need validation first!
					invoicePBXExt.serviceProduct = "⊛ YAP PBX Cease ⊛"
					invoicePBXExt.customerID = userCustomerID
					invoicePBXExt.pbxID = deletePBXSelectPBXID
					invoicePBXExt.tag = deletePBXSelectPBXID
					invoicePBXExt.contractStartDate = currentDate()
					invoicePBXExt.sellPrice = pbxCeasePrice
					invoicePBXExt.salesTaxRate = salesTaxRate
					invoicePBXExt.salesTaxStatus = salesTaxStatus

					// Add PBX cease to invoice
					invoicePBXExtAdd(dbDetail, invoicePBXExt)

					// Delete PBX rental record from invoice_item table
					dbDetail.connection.Query(`DELETE FROM invoice_item WHERE pbx_id = ? AND service_product_name = ?;`, deletePBXSelectPBXID, "⊛ YAP PBX Rental ⊛")

				} else {
					messageHTML(w, validationMessagePBXNotDeleted, "warning")
				}
			}

		} else if genDetail.userTypeID == "200" && validatePBXWhereID == true && deletePBXSelectConfirm == "yes" {

			dbDetail.columnWhereValue = deletePBXSelectPBXID
			checkPBXExist := selectWhere(dbDetail)

			if checkPBXExist == "" {
				messageHTML(w, validationMessagePBXDoesNotExist, "warning")
			} else {
				// Check how many exts the PBX has
				dbDetail.table = "view___sip_extension_detail"
				dbDetail.countMinusOne = false
				dbDetail.columnWhere = "pbx_id"
				dbDetail.columnWhereValue = deletePBXSelectPBXID
				extTotal := totalTableCountWhere(dbDetail)

				dbDetail.connection.Query(`DELETE FROM pbx WHERE id = ?;`, deletePBXSelectPBXID)

				// Check PBX deleted
				dbDetail.column = "pbx_id"
				checkPBXDeleted := selectWhere(dbDetail)

				if checkPBXDeleted == "" {
					messageHTML(w, validationMessagePBXDeleted, "success")

					// Variables for customer table columns
					dbDetail.table = "view___pbx_detail"
					dbDetail.table = "view___customer_detail"
					dbDetail.column = "pbx_id"
					dbDetail.columnWhere = "customer_id"
					dbDetail.columnWhereValue = genDetail.userCustomerID

					// Add cease charge for exts deleted
					extTotalFloat64 := stringToFloat64(extTotal)

					if extTotalFloat64 >= 1 {

						// Get ext cease price
						dbDetail.column = "customer_ext_cease_price"
						extCeasePrice := selectWhere(dbDetail)

						// Get ext sales tax rate
						dbDetail.column = "customer_ext_sales_tax_rate"
						salesTaxRate := selectWhere(dbDetail)

						// Get ext sales tax status
						dbDetail.column = "customer_ext_sales_tax_status"
						salesTaxStatus := selectWhere(dbDetail)

						invoicePBXExt.serviceProduct = "⊛ YAP Extension Cease ⊛"
						invoicePBXExt.customerID = genDetail.userCustomerID
						invoicePBXExt.pbxID = deletePBXSelectPBXID
						invoicePBXExt.tag = "Ext Cease X" + extTotal + ": " + deletePBXSelectPBXID
						invoicePBXExt.contractStartDate = currentDate()
						extCeasePriceFloat64 := stringToFloat64(extCeasePrice)
						extTotalExtCeasePrice := extTotalFloat64 * extCeasePriceFloat64
						invoicePBXExt.sellPrice = strconv.FormatFloat(extTotalExtCeasePrice, 'f', -1, 64)
						invoicePBXExt.salesTaxRate = salesTaxRate
						invoicePBXExt.salesTaxStatus = salesTaxStatus

						// Add PBX cease to invoice
						invoicePBXExtAdd(dbDetail, invoicePBXExt)

						// Delete PBX rental record from invoice_item table
						dbDetail.connection.Query(`DELETE FROM invoice_item WHERE pbx_id = ? AND service_product_name = ?;`, deletePBXSelectPBXID, "⊛ YAP Extension Rental ⊛")

					}

					// Get PBX cease price
					dbDetail.column = "customer_pbx_cease_price"
					pbxCeasePrice := selectWhere(dbDetail)

					// Get PBX sales tax rate
					dbDetail.column = "customer_pbx_sales_tax_rate"
					salesTaxRate := selectWhere(dbDetail)

					// Get PBX sales tax status
					dbDetail.column = "customer_pbx_sales_tax_status"
					salesTaxStatus := selectWhere(dbDetail)

					// These variables have to be here because they need validation first!
					invoicePBXExt.serviceProduct = "⊛ YAP PBX Cease ⊛"
					invoicePBXExt.customerID = genDetail.userCustomerID
					invoicePBXExt.pbxID = deletePBXSelectPBXID
					invoicePBXExt.tag = deletePBXSelectPBXID
					invoicePBXExt.contractStartDate = currentDate()
					invoicePBXExt.sellPrice = pbxCeasePrice
					invoicePBXExt.salesTaxRate = salesTaxRate
					invoicePBXExt.salesTaxStatus = salesTaxStatus

					// Add PBX cease to invoice
					invoicePBXExtAdd(dbDetail, invoicePBXExt)

					// Delete PBX rental record from invoice_item table
					dbDetail.connection.Query(`DELETE FROM invoice_item WHERE pbx_id = ? AND service_product_name = ?;`, deletePBXSelectPBXID, "⊛ YAP PBX Rental ⊛")

				} else {
					messageHTML(w, validationMessagePBXNotDeleted, "warning")
				}
			}

		} else {
			messageHTML(w, validationMessageInvalid, "warning")
		}
	} else {
		panic("pbxDelete function should only be called with account type ID 100, 200")
	}
}

//----------------------------------------------------------------------------------------------------

// SIP extension page functions
func extList(w http.ResponseWriter, dbDetail databaseFunctionParameter, genDetail generalFunctionParameter) {

	// Only account type ID 100, 200, 201, 300, 301, 302 should be able to use this function
	if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" || genDetail.userTypeID == "300" || genDetail.userTypeID == "301" || genDetail.userTypeID == "302" {

		var (
			sipUsername          string
			sipPassword          string
			inboundContext       string
			codecAllowed         string
			dtmfMode             string
			namedCallGroup       string
			namedPickupGroup     string
			mediaEncryption      string
			iceSupport           string
			directMedia          string
			directMediaMethod    string
			rewriteContact       string
			rtpSymmetric         string
			forceRPort           string
			ipAddressAllowed     string
			allowTransfer        string
			callerID             string
			callerIDPrivacy      string
			contactSIPHeaderUser string
			fromSIPHeaderUser    string
			fromSIPHeaderDomain  string
			stirShaken           string
			stirShakenProfile    string
			registered           string
			pbxID                string
			pbxName              string
			customerID           string
			customerName         string
		)

		// Registered table
		var (
			uri       string
			userAgent string
		)

		var dbTableCountUserExt databaseFunctionParameter
		dbTableCountUserExt.connection = dbDetail.connection
		dbTableCountUserExt.database = dbDetail.database
		dbTableCountUserExt.table = "view___sip_extension_detail"
		dbTableCountUserExt.columnWhere = "sip_username"

		fmt.Fprintf(w, "<table id=\"table\" class=\"table-ext\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table id=\"table\" class=\"table-ext\">")
		fmt.Fprintf(w, "        <tr>")
		if genDetail.userTypeID == "100" {
			fmt.Fprintf(w, "          <th>Total Extensions on the YAP Server</th>")
		} else if genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
			fmt.Fprintf(w, "          <th>Total Extensions for the Customer</th>")
		} else if genDetail.userTypeID == "300" || genDetail.userTypeID == "301" || genDetail.userTypeID == "302" {
			fmt.Fprintf(w, "          <th>Total Extensions Within the PBX</th>")
		}
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		if genDetail.userTypeID == "100" {
			dbTableCountUserExt.countMinusOne = false
			fmt.Fprintf(w, "          <td>"+totalTableCount(dbTableCountUserExt)+"</td>")
		} else if genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
			var dbTableCountUserExtWhere databaseFunctionParameter
			dbTableCountUserExtWhere.connection = dbDetail.connection
			dbTableCountUserExtWhere.database = dbDetail.database
			dbTableCountUserExtWhere.table = "view___sip_extension_detail"
			dbTableCountUserExtWhere.columnWhere = "customer_id"
			dbTableCountUserExtWhere.columnWhereValue = genDetail.userCustomerID
			fmt.Fprintf(w, "    <td>"+totalTableCountWhere(dbTableCountUserExtWhere)+"</td>")
		} else if genDetail.userTypeID == "300" || genDetail.userTypeID == "301" || genDetail.userTypeID == "302" {
			var dbTableCountUserExtWhere databaseFunctionParameter
			dbTableCountUserExtWhere.connection = dbDetail.connection
			dbTableCountUserExtWhere.database = dbDetail.database
			dbTableCountUserExtWhere.table = "view___sip_extension_detail"
			dbTableCountUserExtWhere.columnWhere = "pbx_id"
			dbTableCountUserExtWhere.columnWhereValue = genDetail.userPBXID
			fmt.Fprintf(w, "    <td>"+totalTableCountWhere(dbTableCountUserExtWhere)+"</td>")
		}
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><button onclick=\"toggleExt() \"class=\"button-general button-ext\">&nbsp Show/Hide Extension &nbsp</button></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")

		fmt.Fprintf(w, "<div id=\"ext-div\" style=\"display:none\">")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<table id=\"table\" class=\"table-ext\">")
		fmt.Fprintf(w, "  <tr>")
		if genDetail.userTypeID == "100" {
			fmt.Fprintf(w, "    <th class=\"table-title\";>All Extension Details on the YAP Server:</th>")
		} else if genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
			fmt.Fprintf(w, "    <th class=\"table-title\";>All Extension Details for the Customer:</th>")
		} else if genDetail.userTypeID == "300" || genDetail.userTypeID == "301" || genDetail.userTypeID == "302" {
			fmt.Fprintf(w, "    <th class=\"table-title\";>All Extension Details Within the PBX:</th>")
		}
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		var inputTableHTMLArgument jsFunctionParameter
		inputTableHTMLArgument.inputID = "ext-detail-input-sip-username"
		inputTableHTMLArgument.funcNameJS = "extDetailSearchSIPUsername"
		inputTableHTMLArgument.placeholder = "SIP Username/PBX ID"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "ext-detail-input-option"
		inputTableHTMLArgument.funcNameJS = "extDetailSearchOption"
		inputTableHTMLArgument.placeholder = "Options"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
			inputTableHTMLArgument.inputID = "ext-detail-input-pbx-name"
			inputTableHTMLArgument.funcNameJS = "extDetailSearchPBXName"
			inputTableHTMLArgument.placeholder = "PBX Name"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		}
		if genDetail.userTypeID == "100" {
			inputTableHTMLArgument.inputID = "ext-detail-input-customer-id"
			inputTableHTMLArgument.funcNameJS = "extDetailSearchCustomerID"
			inputTableHTMLArgument.placeholder = "Customer ID"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			fmt.Fprintf(w, "    <br>")
			fmt.Fprintf(w, "    <br>")
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "ext-detail-input-customer-name"
			inputTableHTMLArgument.funcNameJS = "extDetailSearchCustomerName"
			inputTableHTMLArgument.placeholder = "Customer Name"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		}
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table id=\"ext-detail-table\" class=\"table-ext\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <th>SIP Username</th>")
		fmt.Fprintf(w, "          <th>SIP Password</th>")
		fmt.Fprintf(w, "          <th>Registered</th>")
		fmt.Fprintf(w, "          <th>Options</th>")
		if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
			fmt.Fprintf(w, "          <th>PBX Name</th>")
		}
		if genDetail.userTypeID == "100" {
			fmt.Fprintf(w, "          <th>Customer ID</th>")
			fmt.Fprintf(w, "          <th>Customer Name</th>")
		}
		fmt.Fprintf(w, "        </tr>")

		var whereClause string
		var userWhereID string

		if genDetail.userTypeID == "100" {
			whereClause = "WHERE customer_id != ?;"
			userWhereID = "1"
		} else if genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
			whereClause = "WHERE customer_id = ?;"
			userWhereID = genDetail.userCustomerID
		} else if genDetail.userTypeID == "300" || genDetail.userTypeID == "301" || genDetail.userTypeID == "302" {
			whereClause = "WHERE pbx_id = ?;"
			userWhereID = genDetail.userPBXID
		}

		extDetailSQL, err := dbDetail.connection.Query(`SELECT
							sip_username,
							sip_password,
							inbound_context,
							codec_allowed,
							dtmf_mode,
							named_call_group,
							named_pickup_group,
							media_encryption,
							ice_support,
							direct_media,
							direct_media_method,
							rewrite_contact,
							rtp_symmetric,
							force_rport,
							ip_address_allowed,
							allow_transfer,
							caller_id,
							caller_id_privacy,
							contact_sip_header_user,
							from_sip_header_user,
							from_sip_header_domain,
							stir_shaken,
							stir_shaken_profile,  
					                registered,
					                pbx_id,
					                pbx_name,
					                customer_id,
					                customer_name
					              FROM
					                 yap.view___sip_extension_detail
						      `+whereClause, userWhereID)

		// Error
		if err != nil {
			panic(err)

		}

		for extDetailSQL.Next() {

			err = extDetailSQL.Scan(
				&sipUsername,
				&sipPassword,
				&inboundContext,
				&codecAllowed,
				&dtmfMode,
				&namedCallGroup,
				&namedPickupGroup,
				&mediaEncryption,
				&iceSupport,
				&directMedia,
				&directMediaMethod,
				&rewriteContact,
				&rtpSymmetric,
				&forceRPort,
				&ipAddressAllowed,
				&allowTransfer,
				&callerID,
				&callerIDPrivacy,
				&contactSIPHeaderUser,
				&fromSIPHeaderUser,
				&fromSIPHeaderDomain,
				&stirShaken,
				&stirShakenProfile,
				&registered,
				&pbxID,
				&pbxName,
				&customerID,
				&customerName,
			)

			// Error
			if err != nil {
				panic(err)
			}

			fmt.Fprintf(w, "        <tr>")
			fmt.Fprintf(w, "          <td>"+sipUsername)
			var copyButtonJSArgument jsFunctionParameter
			copyButtonJSArgument.buttonCSS = "button-ext"
			copyButtonJSArgument.data = sipUsername
			copyButtonJS(w, copyButtonJSArgument)
			fmt.Fprintf(w, "	  </td>")
			fmt.Fprintf(w, "          <td>"+sipPassword)
			copyButtonJSArgument.data = sipPassword
			copyButtonJS(w, copyButtonJSArgument)
			fmt.Fprintf(w, "          </td>")
			if registered == "1" {
				fmt.Fprintf(w, "          <td>&#128994</td>")
			} else {
				fmt.Fprintf(w, "          <td>&#128308</td>")
			}
			fmt.Fprintf(w, "          <td style=\"text-align: left;\">")
			fmt.Fprintf(w, "          <b>Inbound Context:</b> "+inboundContext+"<br>")
			fmt.Fprintf(w, "          <hr class=\"line-table\"></h>")
			fmt.Fprintf(w, "          <b>Codec Allowed:</b> "+codecAllowed+"<br>")
			fmt.Fprintf(w, "          <b>DTMF Mode:</b> "+dtmfMode+"<br>")
			fmt.Fprintf(w, "          <b>Call Group:</b> "+namedCallGroup+"<br>")
			fmt.Fprintf(w, "          <b>Pickup Group:</b> "+namedPickupGroup+"<br>")
			fmt.Fprintf(w, "          <b>Media Encryption:</b> "+mediaEncryption+"<br>")
			fmt.Fprintf(w, "          <b>ICE Support:</b> "+iceSupport+"<br>")
			fmt.Fprintf(w, "          <b>Direct Media:</b> "+directMedia+"<br>")
			fmt.Fprintf(w, "          <b>Direct Media Method:</b> "+directMediaMethod+"<br>")
			fmt.Fprintf(w, "          <b>Rewrite Contact:</b> "+rewriteContact+"<br>")
			fmt.Fprintf(w, "          <b>RTP Symmetric:</b> "+rtpSymmetric+"<br>")
			fmt.Fprintf(w, "          <b>Force Rport:</b> "+forceRPort+"<br>")
			fmt.Fprintf(w, "          <b>IP Address Allowed:</b> "+ipAddressAllowed+"<br>")
			fmt.Fprintf(w, "          <b>Allow Transfer:</b> "+allowTransfer+"<br>")
			fmt.Fprintf(w, "          <b>Caller ID:</b> "+callerID+"<br>")
			fmt.Fprintf(w, "          <b>Caller ID Privacy:</b> "+callerIDPrivacy+"<br>")
			fmt.Fprintf(w, "          <b>Contact SIP Header User:</b> "+contactSIPHeaderUser+"<br>")
			fmt.Fprintf(w, "          <b>From SIP Header User:</b> "+fromSIPHeaderUser+"<br>")
			fmt.Fprintf(w, "          <b>From SIP Header Domain:</b> "+fromSIPHeaderDomain+"<br>")
			fmt.Fprintf(w, "          <b>STIR/SHAKEN Enabled:</b> "+stirShaken+"<br>")
			fmt.Fprintf(w, "          <b>STIR/SHAKEN Profile:</b> "+stirShakenProfile+"<br>")
			fmt.Fprintf(w, "          </td>")
			if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
				fmt.Fprintf(w, "          <td>"+pbxName+"</td>")
			}
			if genDetail.userTypeID == "100" {
				fmt.Fprintf(w, "          <td>"+customerID+"</td>")
				fmt.Fprintf(w, "          <td>"+customerName+"</td>")
			}
			fmt.Fprintf(w, "        </tr>")
		}

		fmt.Fprintf(w, "      </table>")
		var filterTableJSArgument jsFunctionParameter
		filterTableJSArgument.tableID = "ext-detail-table"
		// Call JS filter function for SIP username in the SIP extension detail table
		filterTableJSArgument.funcNameJS = "extDetailSearchSIPUsername"
		filterTableJSArgument.inputID = "ext-detail-input-sip-username"
		filterTableJSArgument.columnNumber = 0
		filterTableJS(w, filterTableJSArgument)
		// Call JS filter function for options in the SIP extenson detail table
		filterTableJSArgument.funcNameJS = "extDetailSearchOption"
		filterTableJSArgument.inputID = "ext-detail-input-option"
		filterTableJSArgument.columnNumber = 3
		filterTableJS(w, filterTableJSArgument)
		if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
			// Call JS filter function for PBX name in the SIP extension detail table
			filterTableJSArgument.funcNameJS = "extDetailSearchPBXName"
			filterTableJSArgument.inputID = "ext-detail-input-pbx-name"
			filterTableJSArgument.columnNumber = 4
			filterTableJS(w, filterTableJSArgument)
		}
		if genDetail.userTypeID == "100" {
			// Call JS filter function for the customer ID in the SIP extension detail table
			filterTableJSArgument.funcNameJS = "extDetailSearchCustomerID"
			filterTableJSArgument.inputID = "ext-detail-input-customer-id"
			filterTableJSArgument.columnNumber = 5
			filterTableJS(w, filterTableJSArgument)
			// Call JS filter function for the customer name in the SIP extension detail table
			filterTableJSArgument.funcNameJS = "extDetailSearchCustomerName"
			filterTableJSArgument.inputID = "ext-detail-input-customer-name"
			filterTableJSArgument.columnNumber = 6
			filterTableJS(w, filterTableJSArgument)
		}
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "<br>")

		// Registered SIP extension Table
		fmt.Fprintf(w, "<table id=\"table\" class=\"table-ext\">")
		fmt.Fprintf(w, "  <tr>")
		if genDetail.userTypeID == "100" {
			fmt.Fprintf(w, "    <th class=\"table-title\";>All Extensions Registered on the Server:</th>")
		} else if genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
			fmt.Fprintf(w, "    <th class=\"table-title\";>All Extensions Registered for the Customer:</th>")
		} else if genDetail.userTypeID == "300" || genDetail.userTypeID == "301" || genDetail.userTypeID == "302" {
			fmt.Fprintf(w, "    <th class=\"table-title\";>All Extensions Registered Within the PBX:</th>")
		}
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "ext-reg-input-sip-username"
		inputTableHTMLArgument.funcNameJS = "extRegSearchSIPUsername"
		inputTableHTMLArgument.placeholder = "SIP Username/PBX ID"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "ext-reg-input-uri"
		inputTableHTMLArgument.funcNameJS = "extRegSearchURI"
		inputTableHTMLArgument.placeholder = "URI"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "ext-reg-input-user-agent"
		inputTableHTMLArgument.funcNameJS = "extRegSearchUserAgent"
		inputTableHTMLArgument.placeholder = "User Agent"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
			inputTableHTMLArgument.inputID = "ext-reg-input-pbx-name"
			inputTableHTMLArgument.funcNameJS = "extRegSearchPBXName"
			inputTableHTMLArgument.placeholder = "PBX Name"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		}
		if genDetail.userTypeID == "100" {
			fmt.Fprintf(w, "    <br>")
			fmt.Fprintf(w, "    <br>")
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "ext-reg-input-customer-id"
			inputTableHTMLArgument.funcNameJS = "extRegSearchCustomerID"
			inputTableHTMLArgument.placeholder = "Customer ID"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "ext-reg-input-customer-name"
			inputTableHTMLArgument.funcNameJS = "extRegSearchCustomerName"
			inputTableHTMLArgument.placeholder = "Customer Name"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		}
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table id=\"ext-reg-table\" class=\"table-ext\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <th>SIP Username</th>")
		fmt.Fprintf(w, "          <th>URI</th>")
		fmt.Fprintf(w, "          <th>User Agent/SIP Client</th>")
		if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
			fmt.Fprintf(w, "          <th>PBX Name</th>")
		}
		if genDetail.userTypeID == "100" {
			fmt.Fprintf(w, "          <th>Customer ID</th>")
			fmt.Fprintf(w, "          <th>Customer Name</th>")
		}
		fmt.Fprintf(w, "        </tr>")

		extRegSQL, err := dbDetail.connection.Query(`SELECT
							sip_username,
							uri,
							user_agent,
					                pbx_name,
					                customer_id,
					                customer_name
					              FROM
					  	        yap.view___sip_extension_registered
					  	      `+whereClause, userWhereID)

		// Error
		if err != nil {
			panic(err)
		}

		for extRegSQL.Next() {

			err = extRegSQL.Scan(
				&sipUsername,
				&uri,
				&userAgent,
				&pbxName,
				&customerID,
				&customerName,
			)

			// Error
			if err != nil {
				panic(err)
			}
			fmt.Fprintf(w, "        <tr>")
			fmt.Fprintf(w, "          <td>"+sipUsername+"</td>")
			fmt.Fprintf(w, "          <td>"+uri+"</td>")
			fmt.Fprintf(w, "          <td>"+userAgent+"</td>")
			fmt.Fprintf(w, "          <td>"+pbxName+"</td>")
			fmt.Fprintf(w, "          <td>"+customerID+"</td>")
			fmt.Fprintf(w, "          <td>"+customerName+"</td>")
			fmt.Fprintf(w, "        </tr>")
		}

		fmt.Fprintf(w, "      </table>")
		filterTableJSArgument.tableID = "ext-reg-table"
		// Call JS filter function for SIP username in the SIP extension registration (reg) table
		filterTableJSArgument.funcNameJS = "extRegSearchSIPUsername"
		filterTableJSArgument.inputID = "ext-reg-input-sip-username"
		filterTableJSArgument.columnNumber = 0
		filterTableJS(w, filterTableJSArgument)
		// Call JS filter function for URI in the SIP extension registration (reg) table
		filterTableJSArgument.funcNameJS = "extRegSearchURI"
		filterTableJSArgument.inputID = "ext-reg-input-uri"
		filterTableJSArgument.columnNumber = 1
		filterTableJS(w, filterTableJSArgument)
		// Call JS filter function for user agent in the SIP extension registration (reg) table
		filterTableJSArgument.funcNameJS = "extRegSearchUserAgent"
		filterTableJSArgument.inputID = "ext-reg-input-user-agent"
		filterTableJSArgument.columnNumber = 2
		filterTableJS(w, filterTableJSArgument)
		if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
			// Call JS filter function for PBX name in the SIP extension registration (reg) table
			filterTableJSArgument.funcNameJS = "extRegSearchPBXName"
			filterTableJSArgument.inputID = "ext-reg-input-pbx-name"
			filterTableJSArgument.columnNumber = 3
			filterTableJS(w, filterTableJSArgument)
		}
		if genDetail.userTypeID == "100" {
			// Call JS filter function for the customer ID in the SIP extension registration (reg) table
			filterTableJSArgument.funcNameJS = "extRegSearchCustomerID"
			filterTableJSArgument.inputID = "ext-reg-input-customer-id"
			filterTableJSArgument.columnNumber = 4
			filterTableJS(w, filterTableJSArgument)
			// Call JS filter function for the customer name in the SIP extension registration (reg) table
			filterTableJSArgument.funcNameJS = "extRegSearchCustomerName"
			filterTableJSArgument.inputID = "ext-reg-input-customer-name"
			filterTableJSArgument.columnNumber = 5
			filterTableJS(w, filterTableJSArgument)
		}
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</div>")
		var toggleDivJSArgument jsFunctionParameter
		toggleDivJSArgument.funcNameJS = "toggleExt"
		toggleDivJSArgument.divID = "ext-div"
		toggleDivJS(w, toggleDivJSArgument)

	} else {
		panic("extList function should only be called with account type ID 100, 200, 201, 300, 301, 302")
	}
}

// Add ext function
func extAdd(w http.ResponseWriter, r *http.Request, dbDetail databaseFunctionParameter, genDetail generalFunctionParameter) {

	// Only account type ID 100, 200, 201, 300, 301 should be able to use this function
	if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" || genDetail.userTypeID == "300" || genDetail.userTypeID == "301" {

		fmt.Fprintf(w, "<form method=\"POST\" action=\"/extension\">")
		fmt.Fprintf(w, "<table class=\"table-add\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>Add a New Extension</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		pbxIDNameList, _ := pbxSlice(dbDetail)
		dbDetail.columnWhere = "customer_id"
		dbDetail.columnWhereValue = genDetail.userCustomerID
		pbxWhereIDNameList, _ := pbxWhereSlice(dbDetail)
		if genDetail.userTypeID == "100" {
			selectDoubleHTML(w, "add_ext_select_pbx_id", "PBX", pbxIDNameList)
		} else if genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
			selectDoubleHTML(w, "add_ext_select_pbx_id", "PBX", pbxWhereIDNameList)
		} else if genDetail.userTypeID == "300" || genDetail.userTypeID == "301" {
			inputReadOnlyHTML(w, "add_ext_input_pbx_id", "PBX ID is", genDetail.userPBXID)
		}
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_ext_input_sip_ext", "Extension (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		codecAllowedValueNameList, _ := codecAllowedSlice()
		selectDoubleHiddenHTML(w, "add_ext_select_codec_allowed", "Codec Allowed (Cannot Be Empty)", codecAllowedValueNameList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		dtmfModeValueNameList, _ := dtmfModeSlice()
		selectDoubleHiddenHTML(w, "add_ext_select_dtmf_mode", "DTMF Mode", dtmfModeValueNameList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_ext_input_call_group", "Call Group")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_ext_input_pickup_group", "Pickup Group")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td style=\"border: none;\">")
		fmt.Fprintf(w, "            <br>")
		fmt.Fprintf(w, "            <br>")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		mediaEncryptionValueNameList, _ := mediaEncryptionSlice()
		selectDoubleHiddenHTML(w, "add_ext_select_media_encryption", "Media Encryption (Recommended)", mediaEncryptionValueNameList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		iceSupportValueList := yesNoSlice()
		selectSingleHTML(w, "add_ext_select_ice_support", "ICE Support", iceSupportValueList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		directMediaValueList := yesNoSlice()
		selectSingleHTML(w, "add_ext_select_direct_media", "Direct Media", directMediaValueList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		directMediaMethodValueNameList, _ := directMediaMethodSlice()
		selectDoubleHiddenHTML(w, "add_ext_select_direct_media_method", "Direct Media Method", directMediaMethodValueNameList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		rewriteContactValueList := yesNoSlice()
		selectSingleHTML(w, "add_ext_select_rewrite_contact", "Rewrite Contact", rewriteContactValueList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		rtpSymmetricValueList := yesNoSlice()
		selectSingleHTML(w, "add_ext_select_rtp_symmetric", "RTP Symmetric", rtpSymmetricValueList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		forceRPortValueList := yesNoSlice()
		selectSingleHTML(w, "add_ext_select_force_rport", "Force RPort", forceRPortValueList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_ext_input_ip_address", "IP Address to Restrict Ext to")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		allowTransferValueList := yesNoSlice()
		callerIDPrivacyValueNameList, _ := callerIDPrivacySlice()
		stirShakenValueList := yesNoSlice()
		if genDetail.userTypeID == "100" {
			fmt.Fprintf(w, "        <tr>")
			fmt.Fprintf(w, "          <td style=\"border: none;\">")
			fmt.Fprintf(w, "            <br>")
			fmt.Fprintf(w, "            <br>")
			fmt.Fprintf(w, "          </td>")
			fmt.Fprintf(w, "        </tr>")
			fmt.Fprintf(w, "        <tr>")
			fmt.Fprintf(w, "          <td>")
			selectSingleHTML(w, "add_ext_select_allow_transfer", "Allow Transfer", allowTransferValueList)
			fmt.Fprintf(w, "          </td>")
			fmt.Fprintf(w, "          <td>")
			inputHTML(w, "add_ext_input_caller_id", "Caller ID")
			fmt.Fprintf(w, "          </td>")
			fmt.Fprintf(w, "          <td>")
			selectDoubleHiddenHTML(w, "add_ext_select_caller_id_privacy", "Caller ID Privacy", callerIDPrivacyValueNameList)
			fmt.Fprintf(w, "          </td>")
			fmt.Fprintf(w, "          <td>")
			inputHTML(w, "add_ext_input_contact_user", "SIP Header - Contact User")
			fmt.Fprintf(w, "          </td>")
			fmt.Fprintf(w, "        </tr>")
			fmt.Fprintf(w, "        <tr>")
			fmt.Fprintf(w, "          <td>")
			inputHTML(w, "add_ext_input_from_user", "SIP Header - From User")
			fmt.Fprintf(w, "          </td>")
			fmt.Fprintf(w, "          <td>")
			inputHTML(w, "add_ext_input_from_domain", "SIP header - From Domain")
			fmt.Fprintf(w, "          </td>")
			fmt.Fprintf(w, "          <td>")
			selectSingleHTML(w, "add_ext_select_stir_shaken", "Stir Shaken", stirShakenValueList)
			fmt.Fprintf(w, "          </td>")
			fmt.Fprintf(w, "          <td>")
			inputHTML(w, "add_ext_input_stir_shaken_profile", "Stir Shaken Profile")
			fmt.Fprintf(w, "          </td>")
			fmt.Fprintf(w, "        </tr>")
		}
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><input type=\"submit\" value=\"Create Extension\"></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</form>")

		addExtSelectPBXID := r.FormValue("add_ext_select_pbx_id")
		addExtInputExt := r.FormValue("add_ext_input_sip_ext")
		addExtSelectCodecAllowed := r.FormValue("add_ext_select_codec_allowed")
		addExtSelectDTMFMode := r.FormValue("add_ext_select_dtmf_mode")
		addExtInputCallGroup := r.FormValue("add_ext_input_call_group")
		addExtInputPickupGroup := r.FormValue("add_ext_input_pickup_group")

		addExtSelectMediaEncryption := r.FormValue("add_ext_select_media_encryption")
		addExtSelectICESupport := r.FormValue("add_ext_select_ice_support")
		addExtSelectDirectMedia := r.FormValue("add_ext_select_direct_media")
		addExtSelectDirectMediaMethod := r.FormValue("add_ext_select_direct_media_method")
		addExtSelectRewriteContact := r.FormValue("add_ext_select_rewrite_contact")
		addExtSelectRTPSymmetric := r.FormValue("add_ext_select_rtp_symmetric")
		addExtSelectForceRPort := r.FormValue("add_ext_select_force_rport")
		addExtInputIPAddress := r.FormValue("add_ext_input_ip_address")

		addExtSelectAllowTransfer := r.FormValue("add_ext_select_allow_transfer")
		addExtInputCallerID := r.FormValue("add_ext_input_caller_id")
		addExtSelectCallerIDPrivacy := r.FormValue("add_ext_select_caller_id_privacy")
		addExtInputContactUser := r.FormValue("add_ext_input_contact_user")
		addExtInputFromUser := r.FormValue("add_ext_input_from_user")
		addExtInputFromDomain := r.FormValue("add_ext_input_from_domain")
		addExtSelectStirShaken := r.FormValue("add_ext_select_stir_shaken")
		addExtInputStirShakenProfile := r.FormValue("add_ext_input_stir_shaken_profile")

		// Validate the PBX ID
		_, pbxIDList := pbxSlice(dbDetail)
		validatePBXID := slices.Contains(pbxIDList, addExtSelectPBXID)
		_, pbxWhereIDList := pbxWhereSlice(dbDetail)
		validatePBXWhereID := slices.Contains(pbxWhereIDList, addExtSelectPBXID)
		// Validate Ext
		validateExt := validateInput(addExtInputExt, "alphaNum")
		// Validate Codec
		_, codecAllowedValueList := codecAllowedSlice()
		validateCodecAllowed := slices.Contains(codecAllowedValueList, addExtSelectCodecAllowed)
		// Validate DTMF Mode
		_, dtmfModeValueList := dtmfModeSlice()
		validateDTMFMode := slices.Contains(dtmfModeValueList, addExtSelectDTMFMode)
		// Validate Call Group
		validateCallGroup := validateInput(addExtInputCallGroup, "alphaNumEmpty")
		// Validate Pickup Group
		validatePickupGroup := validateInput(addExtInputPickupGroup, "alphaNumEmpty")

		// Validate Media Encryption
		_, mediaEncryptionValueList := mediaEncryptionSlice()
		validateMediaEncryption := slices.Contains(mediaEncryptionValueList, addExtSelectMediaEncryption)
		// Validate ICE Support
		validateICESupport := slices.Contains(iceSupportValueList, addExtSelectICESupport)
		// Validate Direct Media
		validateDirectMedia := slices.Contains(directMediaValueList, addExtSelectDirectMedia)
		// Validate Direct Media Method
		_, directMediaMethodValueList := directMediaMethodSlice()
		validateDirectMediaMethod := slices.Contains(directMediaMethodValueList, addExtSelectDirectMediaMethod)
		// Validate Rewrite Contact
		validateRewriteContact := slices.Contains(rewriteContactValueList, addExtSelectRewriteContact)
		// Validate RTP Symmetric
		validateRTPSymmetric := slices.Contains(rtpSymmetricValueList, addExtSelectRTPSymmetric)
		// Validate Force RPort
		validateForceRPort := slices.Contains(forceRPortValueList, addExtSelectForceRPort)
		// Validate IP Address
		validateIPAddress := validateInput(addExtInputIPAddress, "ipAddress")

		// Validate Allow Transfer
		validateAllowTransfer := slices.Contains(allowTransferValueList, addExtSelectAllowTransfer)
		// Validate Direct Media Method
		validateCallerID := validateInput(addExtInputCallerID, "alphaNumEmpty")
		// Validate Rewrite Contact
		_, callerIDPrivacyValueList := callerIDPrivacySlice()
		validateCallerIDPrivacy := slices.Contains(callerIDPrivacyValueList, addExtSelectCallerIDPrivacy)
		// Validate Contact User
		validateContactUser := validateInput(addExtInputContactUser, "alphaNumEmpty")
		// Validate From User
		validateFromUser := validateInput(addExtInputFromUser, "alphaNumEmpty")
		// Validate From Domain
		validateFromDomain := validateInput(addExtInputFromDomain, "alphaNumEmpty")
		// Validate Stir/Shaken
		validateStirShaken := slices.Contains(stirShakenValueList, addExtSelectStirShaken)
		// Validate Stir/Shaken Profile
		validateStirShakenProfile := validateInput(addExtInputStirShakenProfile, "alphaNumEmpty")

		if genDetail.userTypeID == "300" || genDetail.userTypeID == "301" {
			addExtSelectPBXID = genDetail.userPBXID
		}

		if addExtSelectPBXID == "" && addExtInputExt == "" && addExtSelectCodecAllowed == "" {
			// Do Nothing
		} else if genDetail.userTypeID == "300" && addExtInputExt == "" && addExtSelectCodecAllowed == "" {
			// Do Nothing
		} else if genDetail.userTypeID == "301" && addExtInputExt == "" && addExtSelectCodecAllowed == "" {
			// Do Nothing
		} else if genDetail.userTypeID == "100" && validatePBXID == false {
			messageHTML(w, validationMessagePBX, "warning")
		} else if genDetail.userTypeID == "200" && validatePBXWhereID == false {
			messageHTML(w, validationMessagePBX, "warning")
		} else if genDetail.userTypeID == "201" && validatePBXWhereID == false {
			messageHTML(w, validationMessagePBX, "warning")
		} else if validateExt == false {
			messageHTML(w, validationMessageExt, "warning")
		} else if validateCodecAllowed == false || addExtSelectCodecAllowed == "" {
			messageHTML(w, validationMessageExtCodecAllowed, "warning")
		} else if validateDTMFMode == false {
			messageHTML(w, validationMessageExtDTMFMode, "warning")
		} else if validateCallGroup == false {
			messageHTML(w, validationMessageExtCallGroup, "warning")
		} else if validatePickupGroup == false {
			messageHTML(w, validationMessageExtPickupGroup, "warning")
		} else if validateMediaEncryption == false {
			messageHTML(w, validationMessageExtMediaEncryption, "warning")
		} else if validateICESupport == false {
			messageHTML(w, validationMessageExtICESupport, "warning")
		} else if validateDirectMedia == false {
			messageHTML(w, validationMessageExtDirectMedia, "warning")
		} else if validateDirectMediaMethod == false {
			messageHTML(w, validationMessageExtDirectMediaMethod, "warning")
		} else if validateRewriteContact == false {
			messageHTML(w, validationMessageExtRewriteContact, "warning")
		} else if validateRTPSymmetric == false {
			messageHTML(w, validationMessageExtRTPSymmetric, "warning")
		} else if validateForceRPort == false {
			messageHTML(w, validationMessageExtForceRPort, "warning")
		} else if validateIPAddress == false {
			messageHTML(w, validationMessageExtRestrictExt, "warning")
		} else if validateAllowTransfer == false {
			messageHTML(w, validationMessageExtAllowTransfer, "warning")
		} else if validateCallerID == false {
			messageHTML(w, validationMessageExtCallerID, "warning")
		} else if validateCallerIDPrivacy == false {
			messageHTML(w, validationMessageExtCallerIDPrivacy, "warning")
		} else if validateContactUser == false {
			messageHTML(w, validationMessageExtContactUser, "warning")
		} else if validateFromUser == false {
			messageHTML(w, validationMessageExtFromUser, "warning")
		} else if validateFromDomain == false {
			messageHTML(w, validationMessageExtFromDomain, "warning")
		} else if validateStirShaken == false {
			messageHTML(w, validationMessageExtStirShaken, "warning")
		} else if validateStirShakenProfile == false {
			messageHTML(w, validationMessageExtStirShakenProfile, "warning")
		} else {

			// Used to compare the max allowed extensions to the number of extensions that already exist
			var extMaxLimit string
			dbDetail.table = "view___pbx_detail"
			dbDetail.column = "pbx_sip_extension_limit"
			dbDetail.columnWhere = "pbx_id"
			dbDetail.columnWhereValue = addExtSelectPBXID
			extMaxLimit = selectWhere(dbDetail)

			var extCount string
			dbDetail.table = "view___sip_extension_detail"
			dbDetail.column = "sip_username"
			dbDetail.columnWhere = "pbx_id"
			dbDetail.countMinusOne = false
			dbDetail.columnWhereValue = addExtSelectPBXID
			extCount = totalTableCountWhere(dbDetail)

			if extCount >= extMaxLimit {
				messageHTML(w, validationMessageExtMaxExt, "warning")
			} else {

				extPBXID := addExtInputExt + "-" + addExtSelectPBXID

				var callGroupPBXID string
				if addExtInputCallGroup != "" {
					callGroupPBXID = addExtInputCallGroup + "-" + addExtSelectPBXID
				}

				var pickupGroupPBXID string
				if addExtInputPickupGroup != "" {
					pickupGroupPBXID = addExtInputPickupGroup + "-" + addExtSelectPBXID
				}

				dbDetail.table = "view___sip_extension_detail"
				dbDetail.column = "sip_username"
				dbDetail.columnWhere = "sip_username"
				dbDetail.columnWhereValue = extPBXID

				checkExtExist := selectWhere(dbDetail)

				if checkExtExist == extPBXID {
					messageHTML(w, validationMessageExtAlreadyExist, "warning")
				} else {

					if genDetail.userTypeID != "100" {
						addExtSelectAllowTransfer = ""
						addExtInputCallerID = ""
						addExtSelectCallerIDPrivacy = ""
						addExtInputContactUser = ""
						addExtInputFromUser = ""
						addExtInputFromDomain = ""
						addExtSelectStirShaken = ""
						addExtInputStirShakenProfile = ""
					}

					dbDetail.connection.Query(`INSERT 
                                   INTO
                                 ps_endpoints (
                                   id,
                                   aors,
                                   auth,
                                   context,
                                   disallow,
                                   allow,
                                   dtmf_mode,
                                   named_call_group,
                                   named_pickup_group,
                                   media_encryption,
                                   ice_support,
                                   direct_media,
                                   direct_media_method,
                                   rewrite_contact,
                                   rtp_symmetric,
                                   force_rport,
                                   allow_transfer,
                                   callerid,
                                   callerid_privacy,
                                   contact_user,
                                   from_user,
                                   from_domain,
                                   stir_shaken,
                                   stir_shaken_profile,
                                   pbx_id
                                 )
                                 VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`,
						extPBXID,
						extPBXID,
						extPBXID,
						"inbound_"+addExtSelectPBXID,
						"all",
						addExtSelectCodecAllowed,
						nullSQL(addExtSelectDTMFMode),
						nullSQL(callGroupPBXID),
						nullSQL(pickupGroupPBXID),
						nullSQL(addExtSelectMediaEncryption),
						nullSQL(addExtSelectICESupport),
						nullSQL(addExtSelectDirectMedia),
						nullSQL(addExtSelectDirectMediaMethod),
						nullSQL(addExtSelectRewriteContact),
						nullSQL(addExtSelectRTPSymmetric),
						nullSQL(addExtSelectForceRPort),
						nullSQL(addExtSelectAllowTransfer),
						nullSQL(addExtInputCallerID),
						nullSQL(addExtSelectCallerIDPrivacy),
						nullSQL(addExtInputContactUser),
						nullSQL(addExtInputFromUser),
						nullSQL(addExtInputFromDomain),
						nullSQL(addExtSelectStirShaken),
						nullSQL(addExtInputStirShakenProfile),
						addExtSelectPBXID)

					dbDetail.connection.Query(`INSERT 
                                   INTO
                                 ps_aors (
                                   id,
                                   max_contacts,
                                   pbx_id
                                 )
                                 VALUES(?, ?, ?);`,
						extPBXID,
						"2",
						addExtSelectPBXID)

					password := genPassword(20)

					dbDetail.connection.Query(`INSERT 
                                   INTO
                                 ps_auths (
                                   id,
                                   auth_type,
                                   username,
                                   password,
                                   pbx_id
                                 )
                                 VALUES(?, ?, ?, ?, ?);`,
						extPBXID,
						"userpass",
						extPBXID,
						password,
						addExtSelectPBXID)

					checkExtCreated := selectWhere(dbDetail)

					if checkExtCreated == extPBXID {
						messageHTML(w, validationMessageExtCreated, "success")

						// If userTypeID is 100 then get userCustomerID based on the addExtSelectPBXID
						if genDetail.userTypeID == "100" {
							dbDetail.table = "view___pbx_detail"
							dbDetail.column = "customer_id"
							dbDetail.columnWhere = "pbx_id"
							dbDetail.columnWhereValue = addExtSelectPBXID

							genDetail.userCustomerID = selectWhere(dbDetail)
						}

						dbDetail.table = "view___customer_detail"
						dbDetail.columnWhere = "customer_id"
						dbDetail.columnWhereValue = genDetail.userCustomerID

						// Get ext setup price
						dbDetail.column = "customer_ext_setup_price"
						setupPrice := selectWhere(dbDetail)

						// Get ext sales tax rate
						dbDetail.column = "customer_ext_sales_tax_rate"
						salesTaxRate := selectWhere(dbDetail)

						// Get ext sales tax status
						dbDetail.column = "customer_ext_sales_tax_status"
						salesTaxStatus := selectWhere(dbDetail)

						// Get ext contract length
						dbDetail.column = "customer_ext_contract_length"
						contractLength := selectWhere(dbDetail)

						var invoicePBXExt invoicePBXExtFunctionParameter

						invoicePBXExt.customerID = genDetail.userCustomerID
						invoicePBXExt.pbxID = addExtSelectPBXID
						invoicePBXExt.serviceProduct = "⊛ YAP Extension Setup ⊛"
						invoicePBXExt.tag = extPBXID
						invoicePBXExt.sellPrice = setupPrice
						invoicePBXExt.salesTaxRate = salesTaxRate
						invoicePBXExt.salesTaxStatus = salesTaxStatus
						invoicePBXExt.billItemOnce = "yes"
						invoicePBXExt.itemOnHold = "no"
						invoicePBXExt.contractLength = contractLength
						invoicePBXExt.contractStartDate = currentDate()

						// Add extension setup to invoice
						invoicePBXExtAdd(dbDetail, invoicePBXExt)

						// Get ext setup price
						dbDetail.column = "customer_ext_rental_price"
						rentalPrice := selectWhere(dbDetail)

						invoicePBXExt.customerID = genDetail.userCustomerID
						invoicePBXExt.pbxID = addExtSelectPBXID
						invoicePBXExt.serviceProduct = "⊛ YAP Extension Rental ⊛"
						invoicePBXExt.tag = extPBXID
						invoicePBXExt.sellPrice = rentalPrice
						invoicePBXExt.salesTaxRate = salesTaxRate
						invoicePBXExt.salesTaxStatus = salesTaxStatus
						invoicePBXExt.billItemOnce = "no"
						invoicePBXExt.itemOnHold = "yes"
						invoicePBXExt.contractLength = contractLength
						invoicePBXExt.contractStartDate = currentDate()

						// Add extension rental to invoice
						invoicePBXExtAdd(dbDetail, invoicePBXExt)

					} else {
						messageHTML(w, validationMessageExtNotCreated, "success")
					}
				}
			}
		}
	} else {
		panic("extAdd function should only be called with account type ID 100, 200, 201, 300, 301")
	}
}

// Ext edit function
func extEdit(w http.ResponseWriter, r *http.Request, dbDetail databaseFunctionParameter, genDetail generalFunctionParameter) {

	// Only account type ID 100, 200, 201, 300, 301 should be able to use this function
	if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" || genDetail.userTypeID == "300" || genDetail.userTypeID == "301" {

		// List of all column names to edit in the ps_endpoints table
		extColumnList := [][]string{
			{"password", "Generate New Ext Password"},
			{"allow", "Codec"},
			{"dtmf_mode", "DTMF Mode"},
			{"named_call_group", "Call Group"},
			{"named_pickup_group", "Pickup Group"},
			{"media_encryption", "Media Encryption"},
			{"ice_support", "ICE Support"},
			{"direct_media", "Direct Media"},
			{"direct_media_method", "Direct Media Method"},
			{"rewrite_contact", "Rewrite Contact"},
			{"rtp_symmetric", "RTP Symmetric"},
			{"force_rport", "Force RPort"},
			{"permit", "Restrict to IP Address"},
		}

		extExtraColumnList := [][]string{
			{"allow_transfer", "Allow Transfer"},
			{"callerid", "Caller ID"},
			{"callerid_privacy", "Caller ID Privacy"},
			{"contact_user", "SIP Header - Contact User"},
			{"from_user", "SIP Header - From User"},
			{"from_domain", "SIP Header - From Domain"},
			{"stir_shaken", "Stir Shaken"},
			{"stir_shaken_profile", "Stir Shaken Profile"},
		}

		extExtraColumnList = append(extColumnList, extExtraColumnList...)

		fmt.Fprintf(w, "<form method=\"POST\" action=\"/extension\">")
		fmt.Fprintf(w, "<table class=\"table-ext\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>Edit Extension Details</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td style=\"text-align: left;\">")
		fmt.Fprintf(w, "      <b><u>Acceptable Values for Columns</u></b><br><br>")
		fmt.Fprintf(w, "      <b>Generate New Ext Password:</b> yes<br>")
		fmt.Fprintf(w, "      <b>Codec:</b> alaw, ulaw<br>")
		fmt.Fprintf(w, "      <b>Call Group:</b> text<br>")
		fmt.Fprintf(w, "      <b>Pickup Group:</b> text<br>")
		fmt.Fprintf(w, "      <b>Media Encryption:</b> sdes, no<br>")
		fmt.Fprintf(w, "      <b>ICE Support:</b> yes, no<br>")
		fmt.Fprintf(w, "      <b>Direct Media:</b> yes, no<br>")
		fmt.Fprintf(w, "      <b>Direct Media Method:</b> invite, reinvite, update<br>")
		fmt.Fprintf(w, "      <b>Rewrite Contact:</b> yes, no<br>")
		fmt.Fprintf(w, "      <b>RTP Symmetric:</b> yes, no<br>")
		fmt.Fprintf(w, "      <b>Force RPort:</b> yes, no<br>")
		fmt.Fprintf(w, "      <b>Restrict to IP Address:</b> valid IP address<br>")
		if genDetail.userTypeID == "100" {
			fmt.Fprintf(w, "      <b>Allow Transfer:</b> yes, no<br>")
			fmt.Fprintf(w, "      <b>Caller ID:</b> text<br>")
			fmt.Fprintf(w, "      <b>Caller ID Privacy:</b> allowed_not_screened, allowed_passed_screen, allowed_failed_screen, allowed, prohib_not_screened, prohib_passed_screen, prohib_failed_screen, prohib, unavailable<br>")
			fmt.Fprintf(w, "      <b>SIP Header - Contact User:</b> text<br>")
			fmt.Fprintf(w, "      <b>SIP Header - From User:</b> text<br>")
			fmt.Fprintf(w, "      <b>SIP Header - From Domain:</b> text<br>")
			fmt.Fprintf(w, "      <b>Stir Shaken:</b> yes, no<br>")
			fmt.Fprintf(w, "      <b>Stir Shaken Profile:</b> valid file path<br>")
		}
		fmt.Fprintf(w, "    </td>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "edit_ext_input_ext", "Ext (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		if genDetail.userTypeID == "100" {
			selectDoubleHiddenHTML(w, "edit_ext_select_column", "Column to Edit (Cannot Be Empty)", extExtraColumnList)
		} else if genDetail.userTypeID == "200" || genDetail.userTypeID == "201" || genDetail.userTypeID == "300" || genDetail.userTypeID == "301" {
			selectDoubleHiddenHTML(w, "edit_ext_select_column", "Column to Edit (Cannot Be Empty)", extColumnList)
		}
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "edit_ext_input_new_value", "New Value")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><input type=\"submit\" value=\"Update Extension\"></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</form>")

		editExtInputExt := r.FormValue("edit_ext_input_ext")
		editExtSelectColumn := r.FormValue("edit_ext_select_column")
		editExtInputNewValue := r.FormValue("edit_ext_input_new_value")

		// Validate Ext
		validateExt := validateInput(editExtInputExt, "extension")

		var pbxID string

		// If the userTypeID is 300 or 301 then set the pbxID to the user account userPBXID
		if genDetail.userTypeID == "300" || genDetail.userTypeID == "301" {
			pbxID = genDetail.userPBXID
		}

		// If the userTypeID is 200 or 201 then get the extensions pbxID and set it
		dbDetail.table = "view___sip_extension_detail"
		dbDetail.column = "pbx_id"
		dbDetail.columnWhere = "sip_username"
		dbDetail.columnWhereAnd = "customer_id"

		if editExtInputExt == "" && editExtSelectColumn == "" && editExtInputNewValue == "" {
			// Do Nothing
		} else if validateExt == false || editExtInputExt == "" {
			messageHTML(w, validationMessageExt, "warning")
		} else if editExtSelectColumn == "" {
			messageHTML(w, validationMessageExtColumn, "warning")
		} else if editExtSelectColumn == "password" {
			if editExtInputNewValue == "y" || editExtInputNewValue == "Y" || editExtInputNewValue == "yes" || editExtInputNewValue == "Yes" || editExtInputNewValue == "YES" {
				newExtPassword := genPassword(20)
				if genDetail.userTypeID == "100" {
					dbDetail.connection.Query("UPDATE ps_auths SET "+editExtSelectColumn+" = ? WHERE id = ?;", newExtPassword, editExtInputExt)
				} else if genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
					dbDetail.columnWhereValue = editExtInputExt
					dbDetail.columnWhereValueAnd = genDetail.userCustomerID
					pbxID = selectWhereAnd(dbDetail)
					if pbxID == "" {
						// Do Nothing
					} else {
						dbDetail.connection.Query("UPDATE ps_auths SET "+editExtSelectColumn+" = ? WHERE id = ? AND pbx_id = ?;", newExtPassword, editExtInputExt, pbxID)
					}
				} else if genDetail.userTypeID == "300" || genDetail.userTypeID == "301" {
					dbDetail.connection.Query("UPDATE ps_auths SET "+editExtSelectColumn+" = ? WHERE id = ? AND pbx_id = ?;", newExtPassword, editExtInputExt, pbxID)
				}
			} else {
				messageHTML(w, validationMessageGenericAlphaNumEmpty, "warning")
			}
		} else if editExtSelectColumn == "allow" {
			// Validate editExtInputNewValue is a string and not empty
			validateNewValue := validateInput(editExtInputNewValue, "alphaNum")
			if validateNewValue == true {
				if genDetail.userTypeID == "100" {
					dbDetail.connection.Query("UPDATE ps_endpoints SET "+editExtSelectColumn+" = ? WHERE id = ?;", editExtInputNewValue, editExtInputExt)
				} else if genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
					dbDetail.columnWhereValue = editExtInputExt
					dbDetail.columnWhereValueAnd = genDetail.userCustomerID
					pbxID = selectWhereAnd(dbDetail)
					if pbxID == "" {
						// Do Nothing
					} else {
						dbDetail.connection.Query("UPDATE ps_endpoints SET "+editExtSelectColumn+" = ? WHERE id = ? AND pbx_id = ?;", editExtInputNewValue, editExtInputExt, pbxID)
					}
				} else if genDetail.userTypeID == "300" || genDetail.userTypeID == "301" {
					dbDetail.connection.Query("UPDATE ps_endpoints SET "+editExtSelectColumn+" = ? WHERE id = ? AND pbx_id = ?;", editExtInputNewValue, editExtInputExt, pbxID)
				}
			} else {
				messageHTML(w, validationMessageGenericAlphaNumEmpty, "warning")
			}
		} else if editExtSelectColumn == "editExtSelectColumn" || editExtSelectColumn == "dtmf_mode" || editExtSelectColumn == "named_call_group" || editExtSelectColumn == "named_pickup_group" || editExtSelectColumn == "media_encryption" || editExtSelectColumn == "ice_support" || editExtSelectColumn == "direct_media" || editExtSelectColumn == "direct_media_method" || editExtSelectColumn == "rewrite_contact" || editExtSelectColumn == "rtp_symmetric" || editExtSelectColumn == "force_rport" {
			// Validate editExtInputNewValue is a string
			validateNewValue := validateInput(editExtInputNewValue, "alphaNumEmpty")
			if validateNewValue == true {
				if genDetail.userTypeID == "100" {
					dbDetail.connection.Query("UPDATE ps_endpoints SET "+editExtSelectColumn+" = ? WHERE id = ?;", editExtInputNewValue, editExtInputExt)
				} else if genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
					dbDetail.columnWhereValue = editExtInputExt
					dbDetail.columnWhereValueAnd = genDetail.userCustomerID
					pbxID = selectWhereAnd(dbDetail)
					if pbxID == "" {
						// Do Nothing
					} else {
						dbDetail.connection.Query("UPDATE ps_endpoints SET "+editExtSelectColumn+" = ? WHERE id = ? AND pbx_id = ?;", editExtInputNewValue, editExtInputExt, pbxID)
					}
				} else if genDetail.userTypeID == "300" || genDetail.userTypeID == "301" {
					dbDetail.connection.Query("UPDATE ps_endpoints SET "+editExtSelectColumn+" = ? WHERE id = ? AND pbx_id = ?;", editExtInputNewValue, editExtInputExt, pbxID)
				}
			} else {
				messageHTML(w, validationMessageGenericAlphaNumEmpty, "warning")
			}
		} else if editExtSelectColumn == "permit" {
			// Validate editExtInputNewValue is an IP Address
			validateNewValue := validateInput(editExtInputNewValue, "ipAddress")
			if validateNewValue == true {
				if genDetail.userTypeID == "100" {
					dbDetail.connection.Query("UPDATE ps_endpoints SET "+editExtSelectColumn+" = ? WHERE id = ?;", editExtInputNewValue, editExtInputExt)
				} else if genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
					dbDetail.columnWhereValue = editExtInputExt
					dbDetail.columnWhereValueAnd = genDetail.userCustomerID
					pbxID = selectWhereAnd(dbDetail)
					if pbxID == "" {
						// Do Nothing
					} else {
						dbDetail.connection.Query("UPDATE ps_endpoints SET "+editExtSelectColumn+" = ? WHERE id = ? AND pbx_id = ?;", editExtInputNewValue, editExtInputExt, pbxID)
					}
				} else if genDetail.userTypeID == "300" || genDetail.userTypeID == "301" {
					dbDetail.connection.Query("UPDATE ps_endpoints SET "+editExtSelectColumn+" = ? WHERE id = ? AND pbx_id = ?;", editExtInputNewValue, editExtInputExt, pbxID)
				}
			} else {
				messageHTML(w, validationMessageGenericAlphaNumEmpty, "warning")
			}
		} else if genDetail.userTypeID == "100" && editExtSelectColumn == "allow_transfer" || editExtSelectColumn == "callerid" || editExtSelectColumn == "callerid_privacy" || editExtSelectColumn == "contact_user" || editExtSelectColumn == "from_user" || editExtSelectColumn == "from_domain" || editExtSelectColumn == "stir_shaken" {
			// Validate editExtInputNewValue is a string
			validateNewValue := validateInput(editExtInputNewValue, "alphaNumEmpty")
			if validateNewValue == true {
				if genDetail.userTypeID == "100" {
					dbDetail.connection.Query("UPDATE ps_endpoints SET "+editExtSelectColumn+" = ? WHERE id = ?;", editExtInputNewValue, editExtInputExt)
				}
			} else {
				messageHTML(w, validationMessageGenericAlphaNumEmpty, "warning")
			}
		} else if genDetail.userTypeID == "100" && editExtSelectColumn == "stir_shaken_profile" {
			// Validate editExtInputNewValue is a string
			validateNewValue := validateInput(editExtInputNewValue, "alphaNumEmpty")
			if validateNewValue == true {
				if genDetail.userTypeID == "100" {
					dbDetail.connection.Query("UPDATE ps_endpoints SET "+editExtSelectColumn+" = ? WHERE id = ?;", editExtInputNewValue, editExtInputExt)
				}
			} else {
				messageHTML(w, validationMessageGenericAlphaNumEmpty, "warning")
			}
		} else {
			messageHTML(w, validationMessageExtColumn, "warning")
		}
	} else {
		panic("extEdit function should only be called with account type ID 100, 200, 201, 300, 301")
	}
}

// Ext delete function
func extDelete(w http.ResponseWriter, r *http.Request, dbDetail databaseFunctionParameter, genDetail generalFunctionParameter) {

	// Only account type ID 100, 200, 201, 300 should be able to use this function
	if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "201" || genDetail.userTypeID == "300" {

		// Delete an ext
		fmt.Fprintf(w, "<form method=\"POST\" action=\"/extension\">")
		fmt.Fprintf(w, "<table class=\"table-delete\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>Delete an Extension</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "delete_ext_input_ext", "Ext (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		confirmList := yesSlice()
		selectSingleHTML(w, "delete_ext_select_confirm", "yes to Confirm (Cannot Be Empty)", confirmList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><input type=\"submit\" value=\"Delete Extension\"></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</form>")

		deleteExtInputExt := r.FormValue("delete_ext_input_ext")
		deleteExtSelectConfirm := r.FormValue("delete_ext_select_confirm")

		// Validate Ext
		validateExt := validateInput(deleteExtInputExt, "extension")

		if deleteExtInputExt == "" && deleteExtSelectConfirm == "" {
			// Do Nothing
		} else if validateExt == false && deleteExtSelectConfirm == "yes" {
			messageHTML(w, validationMessageExt, "warning")
		} else if validateExt == true && deleteExtSelectConfirm != "yes" {
			messageHTML(w, validationMessageConfirmation, "warning")
		} else if validateExt == true && deleteExtSelectConfirm == "yes" {

			var pbxID string

			// These variables have to be here because they need validation first!
			// Variables for adding cease charge to invoice; they are safe because they are not the input from the user
			var invoicePBXExt invoicePBXExtFunctionParameter
			invoicePBXExt.serviceProduct = "⊛ YAP Extension Cease ⊛"
			invoicePBXExt.billItemOnce = "yes"
			invoicePBXExt.itemOnHold = "no"
			invoicePBXExt.contractLength = ""
			invoicePBXExt.contractStartDate = ""

			if genDetail.userTypeID == "100" {
				// Get customer ID based on the ext
				dbDetail.table = "view___sip_extension_detail"
				dbDetail.column = "customer_id"
				dbDetail.columnWhere = "sip_username"
				dbDetail.columnWhereValue = deleteExtInputExt
				customerID := selectWhere(dbDetail)

				// Get PBX ID based on the ext
				dbDetail.column = "pbx_id"
				pbxID = selectWhere(dbDetail)

				dbDetail.connection.Query(`DELETE FROM ps_endpoints WHERE id = ?;`, deleteExtInputExt)
				dbDetail.table = "view___sip_extension_detail"
				dbDetail.column = "sip_username"
				dbDetail.columnWhere = "sip_username"
				dbDetail.columnWhereValue = deleteExtInputExt
				checkExtDeleted := selectWhere(dbDetail)
				if checkExtDeleted == "" {
					messageHTML(w, validationMessageExtDeleted, "success")

					dbDetail.table = "view___customer_detail"
					dbDetail.columnWhere = "customer_id"
					dbDetail.columnWhereValue = customerID

					// Get Ext cease price
					dbDetail.column = "customer_ext_cease_price"
					extCeasePrice := selectWhere(dbDetail)

					// Get Ext sales tax rate
					dbDetail.column = "customer_ext_sales_tax_rate"
					salesTaxRate := selectWhere(dbDetail)

					// Get Ext sales tax status
					dbDetail.column = "customer_ext_sales_tax_status"
					salesTaxStatus := selectWhere(dbDetail)

					invoicePBXExt.customerID = customerID
					invoicePBXExt.pbxID = pbxID
					invoicePBXExt.tag = deleteExtInputExt
					invoicePBXExt.sellPrice = extCeasePrice
					invoicePBXExt.salesTaxRate = salesTaxRate
					invoicePBXExt.salesTaxStatus = salesTaxStatus

					// Add Ext cease to invoice
					invoicePBXExtAdd(dbDetail, invoicePBXExt)

					// Delete Ext rental record from invoice_item table
					dbDetail.connection.Query(`DELETE FROM invoice_item WHERE tag = ? AND customer_id = ? AND service_product_name = ?;`, deleteExtInputExt, customerID, "⊛ YAP Extension Rental ⊛")

				} else {
					messageHTML(w, validationMessageExtNotDeleted, "warning")
				}
			} else if genDetail.userTypeID == "200" || genDetail.userTypeID == "201" {
				// If the userTypeID is 200 then get the extensions pbxID and set it
				dbDetail.table = "view___sip_extension_detail"
				dbDetail.column = "pbx_id"
				dbDetail.columnWhere = "sip_username"
				dbDetail.columnWhereValue = deleteExtInputExt
				dbDetail.columnWhereAnd = "customer_id"
				dbDetail.columnWhereValueAnd = genDetail.userCustomerID
				pbxID = selectWhereAnd(dbDetail)
				if pbxID == "" {
					messageHTML(w, validationMessageExtDoesNotExist, "warning")
				} else {
					dbDetail.connection.Query("DELETE FROM ps_endpoints WHERE id = ? AND pbx_id = ?;", deleteExtInputExt, pbxID)
					dbDetail.table = "view___sip_extension_detail"
					dbDetail.column = "sip_username"
					dbDetail.columnWhere = "sip_username"
					dbDetail.columnWhereValue = deleteExtInputExt
					checkExtDeleted := selectWhere(dbDetail)
					if checkExtDeleted == "" {

						messageHTML(w, validationMessageExtDeleted, "success")

						dbDetail.table = "view___customer_detail"
						dbDetail.columnWhere = "customer_id"
						dbDetail.columnWhereValue = genDetail.userCustomerID

						// Get Ext cease price
						dbDetail.column = "customer_ext_cease_price"
						extCeasePrice := selectWhere(dbDetail)

						// Get Ext sales tax rate
						dbDetail.column = "customer_ext_sales_tax_rate"
						salesTaxRate := selectWhere(dbDetail)

						// Get Ext sales tax status
						dbDetail.column = "customer_ext_sales_tax_status"
						salesTaxStatus := selectWhere(dbDetail)

						invoicePBXExt.customerID = genDetail.userCustomerID
						invoicePBXExt.pbxID = pbxID
						invoicePBXExt.tag = deleteExtInputExt
						invoicePBXExt.sellPrice = extCeasePrice
						invoicePBXExt.salesTaxRate = salesTaxRate
						invoicePBXExt.salesTaxStatus = salesTaxStatus

						// Add Ext cease to invoice
						invoicePBXExtAdd(dbDetail, invoicePBXExt)

						// Delete Ext rental record from invoice_item table
						dbDetail.connection.Query(`DELETE FROM invoice_item WHERE tag = ? AND customer_id = ? AND service_product_name = ?;`, deleteExtInputExt, genDetail.userCustomerID, "⊛ YAP Extension Rental ⊛")

					} else {
						messageHTML(w, validationMessageExtNotDeleted, "warning")
					}
				}
			} else if genDetail.userTypeID == "300" {
				pbxID = genDetail.userPBXID
				dbDetail.connection.Query("DELETE FROM ps_endpoints WHERE id = ? AND pbx_id = ?;", deleteExtInputExt, pbxID)
				dbDetail.table = "view___sip_extension_detail"
				dbDetail.column = "sip_username"
				dbDetail.columnWhere = "sip_username"
				dbDetail.columnWhereValue = deleteExtInputExt
				checkExtDeleted := selectWhere(dbDetail)
				if checkExtDeleted == "" {

					messageHTML(w, validationMessageExtDeleted, "success")

					dbDetail.table = "view___customer_detail"
					dbDetail.columnWhere = "customer_id"
					dbDetail.columnWhereValue = genDetail.userCustomerID

					// Get Ext cease price
					dbDetail.column = "customer_ext_cease_price"
					extCeasePrice := selectWhere(dbDetail)

					// Get Ext sales tax rate
					dbDetail.column = "customer_ext_sales_tax_rate"
					salesTaxRate := selectWhere(dbDetail)

					// Get Ext sales tax status
					dbDetail.column = "customer_ext_sales_tax_status"
					salesTaxStatus := selectWhere(dbDetail)

					invoicePBXExt.customerID = genDetail.userCustomerID
					invoicePBXExt.pbxID = pbxID
					invoicePBXExt.tag = deleteExtInputExt
					invoicePBXExt.sellPrice = extCeasePrice
					invoicePBXExt.salesTaxRate = salesTaxRate
					invoicePBXExt.salesTaxStatus = salesTaxStatus

					// Add Ext cease to invoice
					invoicePBXExtAdd(dbDetail, invoicePBXExt)

					// Delete Ext rental record from invoice_item table
					dbDetail.connection.Query(`DELETE FROM invoice_item WHERE tag = ? AND customer_id = ? AND service_product_name = ?;`, deleteExtInputExt, genDetail.userCustomerID, "⊛ YAP Extension Rental ⊛")
				} else {
					messageHTML(w, validationMessageExtNotDeleted, "warning")
				}
			}
		}
	} else {
		panic("extDelete function should only be called with account type ID 100, 200, 201, 300")
	}
}

//----------------------------------------------------------------------------------------------------

// Invoice page functions
func invoiceList(w http.ResponseWriter, dbDetail databaseFunctionParameter, genDetail generalFunctionParameter) {

	if genDetail.userTypeID == "100" || genDetail.userTypeID == "200" || genDetail.userTypeID == "400" {

		var (
			invoiceItemID                        string
			customerName                         string
			customerID                           string
			customerUKBased                      string
			customerResellingMinutes             string
			customerUKVATRegistered              string
			customerUKVATNumber                  string
			invoiceItemTag                       string
			invoiceItemSellPrice                 string
			invoiceItemDateTimeAdded             string
			invoiceItemSalesTaxRate              string
			invoiceItemSalesTaxStatus            string
			invoiceBillItemOnce                  string
			invoiceItemOnHold                    string
			invoiceItemContractLength            string
			invoiceItemContractStartDate         string
			serviceProductName                   string
			serviceProductType                   string
			serviceProductSupplierName           string
			serviceProductSupplierContractLength string
		)

		var dbTableCountInvoice databaseFunctionParameter
		dbTableCountInvoice.connection = dbDetail.connection
		dbTableCountInvoice.database = dbDetail.database
		dbTableCountInvoice.table = "view___invoice_item"

		if genDetail.userTypeID == "100" {
			fmt.Fprintf(w, "<table id=\"table\" class=\"table-invoice\">")
			fmt.Fprintf(w, "  <tr>")
			fmt.Fprintf(w, "    <th>")
			fmt.Fprintf(w, "      <table id=\"table\" class=\"table-invoice\">")
			fmt.Fprintf(w, "        <tr>")
			fmt.Fprintf(w, "          <th>Total Invoice Services/Products</th>")
			fmt.Fprintf(w, "        </tr>")
			fmt.Fprintf(w, "        <tr>")
			dbTableCountInvoice.countMinusOne = false
			fmt.Fprintf(w, "          <td>"+totalTableCount(dbTableCountInvoice)+"</td>")
			fmt.Fprintf(w, "        </tr>")
			fmt.Fprintf(w, "      </table>")
			fmt.Fprintf(w, "    </th>")
			fmt.Fprintf(w, "  </tr>")
			fmt.Fprintf(w, "  <tr>")
			fmt.Fprintf(w, "    <th><button onclick=\"toggleInvoice() \"class=\"button-general button-invoice\">&nbsp Show/Hide Invoice &nbsp</button></th>")
			fmt.Fprintf(w, "  </tr>")
			fmt.Fprintf(w, "</table>")
		}

		if genDetail.userTypeID == "100" {
			fmt.Fprintf(w, "<div id=\"invoice-div\" style=\"display:none\">")
			fmt.Fprintf(w, "<br>")
		} else {
			fmt.Fprintf(w, "<div id=\"invoice-div\">")
		}
		fmt.Fprintf(w, "<table id=\"table\" class=\"table-invoice\">")
		fmt.Fprintf(w, "  <tr>")
		if genDetail.userTypeID == "100" {
			fmt.Fprintf(w, "    <th class=\"table-title\";>All Customer Service/Product Invoice Items on the YAP Server:</th>")
		} else {
			fmt.Fprintf(w, "    <th class=\"table-title\";>Customer Invoice Services/Products</th>")
		}
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "    <br>")

		var inputTableHTMLArgument jsFunctionParameter
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "invoice-input-item-id"
		inputTableHTMLArgument.funcNameJS = "invoiceSearchItemID"
		inputTableHTMLArgument.placeholder = "Item ID"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "invoice-input-name-information"
		inputTableHTMLArgument.funcNameJS = "invoiceSearchNameInformation"
		inputTableHTMLArgument.placeholder = "Service/Product Name & Information"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "invoice-input-sale-price"
		inputTableHTMLArgument.funcNameJS = "invoiceSearchSalePrice"
		inputTableHTMLArgument.placeholder = "Service/Product Sale Price"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		if genDetail.userTypeID == "100" {
			inputTableHTMLArgument.inputID = "invoice-input-detail"
			inputTableHTMLArgument.funcNameJS = "invoiceSearchDetail"
			inputTableHTMLArgument.placeholder = "Invoice Item Details"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			fmt.Fprintf(w, "    <br><br>")
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "invoice-input-customer-id"
			inputTableHTMLArgument.funcNameJS = "invoiceSearchCustomerID"
			inputTableHTMLArgument.placeholder = "Customer ID"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "invoice-input-customer-name"
			inputTableHTMLArgument.funcNameJS = "invoiceSearchCustomerName"
			inputTableHTMLArgument.placeholder = "Customer Name"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		}

		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		var exportCSVButtonHTMLArgument jsFunctionParameter
		exportCSVButtonHTMLArgument.funcNameJS = "Invoice"
		exportCSVButtonHTMLArgument.buttonCSS = "button-invoice"
		exportCSVButtonHTML(w, exportCSVButtonHTMLArgument)
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table id=\"invoice-table\" class=\"table-invoice\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <th>Item ID</th>")
		fmt.Fprintf(w, "          <th>Name & Information</th>")
		fmt.Fprintf(w, "          <th>Sale Price</th>")
		fmt.Fprintf(w, "          <th>Details</th>")
		if genDetail.userTypeID == "100" {
			fmt.Fprintf(w, "          <th>Customer ID</th>")
			fmt.Fprintf(w, "          <th>Customer Name</th>")
		}
		fmt.Fprintf(w, "        </tr>")

		var whereClause string

		if genDetail.userTypeID == "100" {
			whereClause = "WHERE customer_id != ?;"
			genDetail.userCustomerID = "1"
		} else if genDetail.userTypeID == "200" || genDetail.userTypeID == "400" {
			whereClause = "WHERE customer_id = ?;"
		}

		invoiceSQL, err := dbDetail.connection.Query(`SELECT
			                     		invoice_item_id,
			                                customer_id,
                                                        customer_name,
                                                        customer_uk_based,
                                                        customer_reselling_minutes,
                                                        customer_uk_vat_registered,
                                                        customer_uk_vat_number,
							invoice_item_tag,
							invoice_item_sell_price,
							invoice_item_date_time_added,
							invoice_item_sales_tax_rate,
							invoice_item_sales_tax_status,
							invoice_bill_item_once,
							invoice_item_on_hold,
							invoice_item_contract_length,
							invoice_item_contract_start_date,
							service_product_name,
							service_product_type,
							service_product_supplier_name,
							service_product_supplier_contract_length
					              FROM
					  	        yap.view___invoice_item
						      `+whereClause, genDetail.userCustomerID)

		// Error
		if err != nil {
			panic(err)

		}

		for invoiceSQL.Next() {

			err = invoiceSQL.Scan(
				&invoiceItemID,
				&customerID,
				&customerName,
				&customerUKBased,
				&customerResellingMinutes,
				&customerUKVATRegistered,
				&customerUKVATNumber,
				&invoiceItemTag,
				&invoiceItemSellPrice,
				&invoiceItemDateTimeAdded,
				&invoiceItemSalesTaxRate,
				&invoiceItemSalesTaxStatus,
				&invoiceBillItemOnce,
				&invoiceItemOnHold,
				&invoiceItemContractLength,
				&invoiceItemContractStartDate,
				&serviceProductName,
				&serviceProductType,
				&serviceProductSupplierName,
				&serviceProductSupplierContractLength,
			)

			// Error
			if err != nil {
				panic(err)
			}

			if genDetail.userTypeID != "100" && invoiceItemOnHold == "yes" {
				// Do Nothing
			} else {
				fmt.Fprintf(w, "        <tr>")
				fmt.Fprintf(w, "          <td>"+invoiceItemID+"</td>")
				fmt.Fprintf(w, "          <td style=\"text-align: left;\">")
				fmt.Fprintf(w, "            "+serviceProductName+"<br><br>")
				fmt.Fprintf(w, "            <b>Service/Product Tag:</b> "+invoiceItemTag+"<br>")
				fmt.Fprintf(w, "          </td>")
				fmt.Fprintf(w, "          <td style=\"text-align: left;\">")
				//If the YAP Admin is UK VAT registered and the service/product is taxable
				if genDetail.yapAdminUKVATRegistered == "yes" && invoiceItemSalesTaxStatus == "TAXABLE" {
					fmt.Fprintf(w, "            <b>Price (exVAT):</b> "+genDetail.currencySymbol+invoiceItemSellPrice+"<br>")
					// Convert UK sales VAT rate to float64
					invoiceItemSalesTaxRateFloat64 := stringToFloat64(invoiceItemSalesTaxRate)
					fmt.Fprintf(w, "            <b>VAT Rate:</b> "+strconv.FormatFloat(invoiceItemSalesTaxRateFloat64, 'f', -1, 64)+"&#37;<br>")
					// Convert item sell price to float64
					invoiceItemSellPriceExVATFloat64 := stringToFloat64(invoiceItemSellPrice)
					var invoiceItemSellPriceIncVATFloat64 float64 = invoiceItemSellPriceExVATFloat64 * (invoiceItemSalesTaxRateFloat64/100 + 1)
					var invoiceItemSellVATFloat64 float64 = invoiceItemSellPriceIncVATFloat64 - invoiceItemSellPriceExVATFloat64
					fmt.Fprintf(w, "            <b>VAT:</b> "+genDetail.currencySymbol+strconv.FormatFloat(invoiceItemSellVATFloat64, 'f', 2, 64)+"<br>")
					fmt.Fprintf(w, "            <b>Total Price (incVAT):</b> "+genDetail.currencySymbol+strconv.FormatFloat(invoiceItemSellPriceIncVATFloat64, 'f', 2, 64)+"<br>")
					//If the YAP Admin is UK VAT registered and the service/product is exempt
				} else if genDetail.yapAdminUKVATRegistered == "yes" && invoiceItemSalesTaxStatus == "EXEMPT" {
					fmt.Fprintf(w, "            <b>Price (exVAT):</b> "+genDetail.currencySymbol+invoiceItemSellPrice+"<br>")
					fmt.Fprintf(w, "            <b>VAT Rate:</b> Exempt<br>")
					fmt.Fprintf(w, "            <b>VAT:</b> "+genDetail.currencySymbol+"0.00</b><br>")
					fmt.Fprintf(w, "            <b>Total Price (incVAT):</b> "+genDetail.currencySymbol+invoiceItemSellPrice)
					//If the YAP Admin is not UK VAT registered
				} else if genDetail.yapAdminUKVATRegistered == "no" {
					fmt.Fprintf(w, "            "+genDetail.currencySymbol+invoiceItemSellPrice+"<br>")
				}
				fmt.Fprintf(w, "          </td>")
				fmt.Fprintf(w, "          <td style=\"text-align: left; vertical-align: top;\">")
				fmt.Fprintf(w, "            <b><u>Item Details</u></b><br><br>")
				fmt.Fprintf(w, " 	    <b>Item Added Date & Time: </b>"+formatDateTime(invoiceItemDateTimeAdded)+"<br>")
				fmt.Fprintf(w, "            <b>Item Type: </b>"+serviceProductType+"<br>")
				if invoiceItemContractLength == "" {
					// Do Nothing
				} else {
					// The date is automatically added for YAP PBX/Ext charges
					fmt.Fprintf(w, "            <b>Contract Start Date: </b>"+formatDate(invoiceItemContractStartDate)+"<br>")
					fmt.Fprintf(w, "            <b>Contract Length: </b>"+invoiceItemContractLength+"<br>")
				}
				if genDetail.userTypeID == "100" {
					fmt.Fprintf(w, "            <b>Sale VAT Status: </b>"+invoiceItemSalesTaxStatus+"<br>")
					fmt.Fprintf(w, "            <b>Bill Item Once: </b>"+invoiceBillItemOnce+"<br>")
					fmt.Fprintf(w, "            <b>Item on Hold: </b>"+invoiceItemOnHold+"<br>")
					fmt.Fprintf(w, "            <b>Item Type: </b>"+serviceProductType+"<br>")
					fmt.Fprintf(w, "            <hr class=\"line-table\"></h>")
					fmt.Fprintf(w, "            <b><u>Customer Details</u></b><br><br>")
					fmt.Fprintf(w, "            <b>Reselling Minutes: </b>"+customerResellingMinutes+"<br>")
					fmt.Fprintf(w, "            <b>UK Based: </b>"+customerUKBased+"<br>")
					fmt.Fprintf(w, "            <b>UK VAT Registered: </b>"+customerUKVATRegistered+"<br>")
					fmt.Fprintf(w, "            <b>UK VAT Number: </b>"+customerUKVATNumber+"<br>")
					fmt.Fprintf(w, "            <hr class=\"line-table\"></h>")
					fmt.Fprintf(w, "            <b><u>Supplier Details</u></b><br><br>")
					fmt.Fprintf(w, "            <b>Supplier Name: </b>"+serviceProductSupplierName+"<br>")
					fmt.Fprintf(w, "            <b>Supplier Contract Length: </b>"+serviceProductSupplierContractLength)
				}
				fmt.Fprintf(w, "          </td>")
				if genDetail.userTypeID == "100" {
					fmt.Fprintf(w, "          <td>"+customerID+"</td>")
					fmt.Fprintf(w, "          <td>"+customerName+"</td>")
				}
				fmt.Fprintf(w, "        </tr>")
			}
		}

		fmt.Fprintf(w, "      </table>")
		var filterTableJSArgument jsFunctionParameter
		filterTableJSArgument.tableID = "invoice-table"

		filterTableJSArgument.funcNameJS = "invoiceSearchItemID"
		filterTableJSArgument.inputID = "invoice-input-item-id"
		filterTableJSArgument.columnNumber = 0
		filterTableJS(w, filterTableJSArgument)

		filterTableJSArgument.funcNameJS = "invoiceSearchNameInformation"
		filterTableJSArgument.inputID = "invoice-input-name-information"
		filterTableJSArgument.columnNumber = 1
		filterTableJS(w, filterTableJSArgument)

		filterTableJSArgument.funcNameJS = "invoiceSearchSalePrice"
		filterTableJSArgument.inputID = "invoice-input-sale-price"
		filterTableJSArgument.columnNumber = 2
		filterTableJS(w, filterTableJSArgument)

		if genDetail.userTypeID == "100" {
			filterTableJSArgument.funcNameJS = "invoiceSearchDetail"
			filterTableJSArgument.inputID = "invoice-input-detail"
			filterTableJSArgument.columnNumber = 3
			filterTableJS(w, filterTableJSArgument)

			filterTableJSArgument.funcNameJS = "invoiceSearchCustomerID"
			filterTableJSArgument.inputID = "invoice-input-customer-id"
			filterTableJSArgument.columnNumber = 4
			filterTableJS(w, filterTableJSArgument)

			filterTableJSArgument.funcNameJS = "invoiceSearchCustomerName"
			filterTableJSArgument.inputID = "invoice-input-customer-name"
			filterTableJSArgument.columnNumber = 5
			filterTableJS(w, filterTableJSArgument)
		}
		var exportCSVJSArgument jsFunctionParameter
		exportCSVJSArgument.funcNameJS = "Invoice"
		exportCSVJSArgument.tableID = "invoice-table"
		exportCSVJSArgument.fileName = "YAP_invoice_item_details"
		exportCSVJSArgument.pathURL = "invoice"
		exportCSVJS(w, exportCSVJSArgument)
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</div>")
		var toggleDivJSArgument jsFunctionParameter
		toggleDivJSArgument.funcNameJS = "toggleInvoice"
		toggleDivJSArgument.divID = "invoice-div"
		toggleDivJS(w, toggleDivJSArgument)

	} else {
		panic("invoiceList function should only be called with account type ID 100, 200, 400")
	}
}

// Invoice add function
func invoiceAdd(w http.ResponseWriter, r *http.Request, dbDetail databaseFunctionParameter, genDetail generalFunctionParameter) {

	// Only account type ID 100 should be able to use this function
	if genDetail.userTypeID == "100" {

		fmt.Fprintf(w, "<form method=\"POST\" action=\"/invoice\">")
		fmt.Fprintf(w, "<table class=\"table-add\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>Add a New Invoice Item<br>(YAP PBX & Ext Invoices Cannot Be Created Manually)</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		customerIDNameList, _ := customerSlice(dbDetail)
		selectDoubleHTML(w, "add_invoice_select_customer_id", "Customer", customerIDNameList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		dbDetail.table = "service_product"
		dbDetail.column = "name"
		serviceProductList := singleColumnSlice(dbDetail)
		serviceProductList = append([]string{""}, serviceProductList...)
		serviceProductList = serviceProductList[:len(serviceProductList)-6]
		selectSingleHTML(w, "add_invoice_select_service_product", "Service/Product (Cannot Be Empty)", serviceProductList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_invoice_input_tag", "Item Tag (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_invoice_input_price", "Item Price (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		dbDetail.table = "sales_tax_rate_lookup"
		dbDetail.column = "sales_tax_rate"
		salesTaxRateList := singleColumnSlice(dbDetail)
		salesTaxRateList = append([]string{""}, salesTaxRateList...)
		selectSingleHTML(w, "add_invoice_select_sales_tax_rate", "Sales Tax Rate &#37; (Cannot Be Empty)", salesTaxRateList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		salesTaxStatusList := []string{"", "TAXABLE", "EXEMPT"}
		selectSingleHTML(w, "add_invoice_select_sales_tax_status", "Sales Tax Status (Cannot Be Empty)", salesTaxStatusList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		billItemOnceList := yesNoSlice()
		selectSingleHTML(w, "add_invoice_select_bill_item_once", "Bill Item Once (Cannot Be Empty)", billItemOnceList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		itemOnHoldList := yesNoSlice()
		selectSingleHTML(w, "add_invoice_select_item_on_hold", "Item On Hold (Cannot Be Empty)", itemOnHoldList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		dbDetail.table = "contract_length_lookup"
		dbDetail.column = "contract_length"
		contractLengthList := singleColumnSlice(dbDetail)
		contractLengthList = append([]string{""}, contractLengthList...)
		selectSingleHTML(w, "add_invoice_select_contract_length", "Contract Length", contractLengthList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		dateHTML(w, "add_invoice_input_contract_start_date", "Contract Start Date")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><input type=\"submit\" value=\"Create Invoice Item\"></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</form>")

		addInvoiceSelectCustomerID := r.FormValue("add_invoice_select_customer_id")
		addInvoiceSelectServiceProduct := r.FormValue("add_invoice_select_service_product")
		addInvoiceInputTag := r.FormValue("add_invoice_input_tag")
		addInvoiceInputSellPrice := r.FormValue("add_invoice_input_price")
		addInvoiceSelectSalesTaxRate := r.FormValue("add_invoice_select_sales_tax_rate")
		addInvoiceSelectSalesTaxStatus := r.FormValue("add_invoice_select_sales_tax_status")
		addInvoiceSelectBillItemOnce := r.FormValue("add_invoice_select_bill_item_once")
		addInvoiceSelectItemOnHold := r.FormValue("add_invoice_select_item_on_hold")
		addInvoiceSelectContractLength := r.FormValue("add_invoice_select_contract_length")
		addInvoiceInputContractStartDate := r.FormValue("add_invoice_input_contract_start_date")

		// Check customer ID is contained in the slice
		_, customerIDList := customerSlice(dbDetail)
		customerIDList = append(customerIDList, "")
		validateCustomerID := slices.Contains(customerIDList, addInvoiceSelectCustomerID)

		// Validate service/product is contained in the slice
		validateServiceProduct := slices.Contains(serviceProductList, addInvoiceSelectServiceProduct)

		// Validate the good/service tag
		validateTag := validateInput(addInvoiceInputTag, "alphaNumEmpty")

		// Validate the invoice item price
		validateSellPrice := validateInput(addInvoiceInputSellPrice, "price")

		// Check sales tax rate is contained in the slice
		validateSalesTaxRate := slices.Contains(salesTaxRateList, addInvoiceSelectSalesTaxRate)

		// Check sales tax status is contained in the slice
		validateSalesTaxStatus := slices.Contains(salesTaxStatusList, addInvoiceSelectSalesTaxStatus)

		// Check bill item once is contained in the slice
		validateBillItemOnce := slices.Contains(billItemOnceList, addInvoiceSelectBillItemOnce)

		// Check item on hold is contained in the slice
		validateItemOnHold := slices.Contains(itemOnHoldList, addInvoiceSelectItemOnHold)

		// Check contract length is contained in the slice
		validateContractLength := slices.Contains(contractLengthList, addInvoiceSelectContractLength)

		// Validate the contract start date
		validateStartDate := validateInput(addInvoiceInputContractStartDate, "date")

		if addInvoiceSelectCustomerID == "" && addInvoiceSelectServiceProduct == "" && addInvoiceInputTag == "" && addInvoiceInputSellPrice == "" && addInvoiceSelectSalesTaxRate == "" && addInvoiceSelectSalesTaxStatus == "" && addInvoiceSelectBillItemOnce == "" && addInvoiceSelectItemOnHold == "" && addInvoiceSelectContractLength == "" && addInvoiceInputContractStartDate == "" {
			// Do Nothing
		} else if validateCustomerID == false || addInvoiceSelectCustomerID == "" {
			messageHTML(w, validationMessageCustomer, "warning")
		} else if validateServiceProduct == false || addInvoiceSelectServiceProduct == "" {
			messageHTML(w, validationMessageInvoiceServiceProduct, "warning")
		} else if validateTag == false {
			messageHTML(w, validationMessageInvoiceServiceProductTag, "warning")
		} else if validateSellPrice == false || addInvoiceInputSellPrice == "" {
			messageHTML(w, validationMessageInvoiceItemPrice, "warning")
		} else if validateSalesTaxRate == false || addInvoiceSelectSalesTaxRate == "" {
			messageHTML(w, validationMessageInvoiceSalesTaxRate, "warning")
		} else if validateSalesTaxStatus == false || addInvoiceSelectSalesTaxStatus == "" {
			messageHTML(w, validationMessageInvoiceSalesTaxStatus, "warning")
		} else if validateBillItemOnce == false || addInvoiceSelectBillItemOnce == "" {
			messageHTML(w, validationMessageInvoiceBillItemOnce, "warning")
		} else if validateItemOnHold == false || addInvoiceSelectItemOnHold == "" {
			messageHTML(w, validationMessageInvoiceItemOnHold, "warning")
		} else if validateContractLength == false {
			messageHTML(w, validationMessageContractLength, "warning")
		} else if addInvoiceSelectContractLength != "" && addInvoiceInputContractStartDate == "" {
			messageHTML(w, validationMessageInvoiceContractStartDateEmpty, "warning")
		} else if validateStartDate == false {
			messageHTML(w, validationMessageInvoiceContractStartDate, "warning")
		} else {

			// Convert string values to a float64 to use the math package to round to the nearest two decimal places
			addInvoiceInputSellPriceFloat64 := stringToFloat64(addInvoiceInputSellPrice)

			dbDetail.connection.Query(`INSERT 
        	              		     INTO
        	              		   invoice_item (
					     customer_id,
					     pbx_id,
					     service_product_name,
					     tag,
					     sell_price,
					     sales_tax_rate,
					     sales_tax_status,
					     bill_item_once,
					     item_on_hold,
					     contract_length,
					     contract_start_date
					   )
					   VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`,
				addInvoiceSelectCustomerID,
				"1",
				addInvoiceSelectServiceProduct,
				nullSQL(addInvoiceInputTag),
				math.Round(addInvoiceInputSellPriceFloat64*100)/100,
				addInvoiceSelectSalesTaxRate,
				addInvoiceSelectSalesTaxStatus,
				addInvoiceSelectBillItemOnce,
				addInvoiceSelectItemOnHold,
				nullSQL(addInvoiceSelectContractLength),
				nullSQL(addInvoiceInputContractStartDate),
			)

		}
	} else {
		panic("invoiceAdd function should only be called with account type ID 100")
	}
}

// Invoice delete function
func invoiceDelete(w http.ResponseWriter, r *http.Request, dbDetail databaseFunctionParameter, genDetail generalFunctionParameter) {

	// Only account type ID 100 should be able to use this function
	if genDetail.userTypeID == "100" {

		// Delete a invoice
		fmt.Fprintf(w, "<form method=\"POST\" action=\"/invoice\">")
		fmt.Fprintf(w, "<table class=\"table-delete\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>Delete an Invoice<br>(YAP PBX and Ext Invoices Cannot Be Deleted Manually)</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "delete_invoice_input_invoice_id", "Invoice ID (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		confirmList := yesSlice()
		selectSingleHTML(w, "delete_invoice_select_confirm", "yes to Confirm (Cannot Be Empty)", confirmList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><input type=\"submit\" value=\"Delete Invoice Item\"></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</form>")

		deleteInvoiceInputInvoiceID := r.FormValue("delete_invoice_input_invoice_id")
		deleteInvoiceSelectConfirm := r.FormValue("delete_invoice_select_confirm")

		// Validate Ext
		validateInvoiceID := validateInput(deleteInvoiceInputInvoiceID, "number")

		if deleteInvoiceInputInvoiceID == "" && deleteInvoiceSelectConfirm == "" {
			// Do Nothing
		} else if validateInvoiceID == false && deleteInvoiceSelectConfirm == "yes" {
			messageHTML(w, validationMessageInvoiceID, "warning")
		} else if validateInvoiceID == true && deleteInvoiceSelectConfirm != "yes" {
			messageHTML(w, validationMessageInvoice, "warning")
		} else if validateInvoiceID == true && deleteInvoiceSelectConfirm == "yes" {

			dbDetail.table = "view___invoice_item"
			dbDetail.column = "invoice_item_id"
			dbDetail.columnWhere = "invoice_item_id"
			dbDetail.columnWhereValue = deleteInvoiceInputInvoiceID

			checkInvoiceExist := selectWhere(dbDetail)

			if checkInvoiceExist == "" {
				messageHTML(w, validationMessageInvoiceDoesNotExist, "warning")
			} else {

				dbDetail.connection.Query(`DELETE FROM invoice_item WHERE id = ? AND pbx_id = ?;`, deleteInvoiceInputInvoiceID, "1")

				checkInvoiceDeleted := selectWhere(dbDetail)

				if checkInvoiceDeleted == "" {
					messageHTML(w, validationMessageInvoiceDeleted, "success")
				} else {
					messageHTML(w, validationMessageInvoiceNotDeleted, "warning")
				}
			}
		} else {
			messageHTML(w, validationMessageInvalid, "warning")
		}
	} else {
		panic("invoiceDelete function should only be called with account type ID 100")
	}
}

//----------------------------------------------------------------------------------------------------

// Service/Products page functions

// Service/products list function
func serviceProductList(w http.ResponseWriter, dbDetail databaseFunctionParameter, genDetail generalFunctionParameter) {

	if genDetail.userTypeID == "100" {

		// service_product table columns
		var (
			serviceProductID                     string
			serviceProductName                   string
			serviceProductType                   string
			serviceProductSupplierName           string
			serviceProductSupplierContractLength string
			serviceProductDateTimeAdded          string
		)

		// sales_tax_rate table column
		var (
			salesTaxRate string
		)

		// supplier table column
		var (
			supplierName          string
			supplierDateTimeAdded string
		)

		var dbTableCountServiceProduct databaseFunctionParameter
		dbTableCountServiceProduct.connection = dbDetail.connection
		dbTableCountServiceProduct.database = dbDetail.database
		dbTableCountServiceProduct.table = "view___service_product"

		fmt.Fprintf(w, "<table id=\"table\" class=\"table-service-product\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table id=\"table\" class=\"table-service-product\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <th>Total Services/Products Available</th>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		dbTableCountServiceProduct.countMinusOne = false
		fmt.Fprintf(w, "          <td>"+totalTableCount(dbTableCountServiceProduct)+"</td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><button onclick=\"toggleServiceProduct() \"class=\"button-general button-service-product\">&nbsp Show/Hide Services/Products &nbsp</button></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")

		fmt.Fprintf(w, "<div id=\"service-product-div\" style=\"display:none\">")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<table id=\"table\" class=\"table-service-product\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>Service/Product Information:</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "    <br>")
		var inputTableHTMLArgument jsFunctionParameter
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "service-product-input-id"
		inputTableHTMLArgument.funcNameJS = "serviceProductSearchID"
		inputTableHTMLArgument.placeholder = "Service/Product ID"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "service-product-input-name"
		inputTableHTMLArgument.funcNameJS = "serviceProductSearchName"
		inputTableHTMLArgument.placeholder = "Service/Product Name"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "service-product-input-type"
		inputTableHTMLArgument.funcNameJS = "serviceProductSearchType"
		inputTableHTMLArgument.placeholder = "Service/Product Type"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "service-product-input-supplier-name"
		inputTableHTMLArgument.funcNameJS = "serviceProductSearchSupplierName"
		inputTableHTMLArgument.placeholder = "Service/Product Supplier Name"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "service-product-input-supplier-contract-length"
		inputTableHTMLArgument.funcNameJS = "serviceProductSearchSupplierContractLength"
		inputTableHTMLArgument.placeholder = "Supplier Contract Length"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "service-product-input-date-time"
		inputTableHTMLArgument.funcNameJS = "serviceProductSearchDateTime"
		inputTableHTMLArgument.placeholder = "Date & Time Added"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		var exportCSVButtonHTMLArgument jsFunctionParameter
		exportCSVButtonHTMLArgument.funcNameJS = "ServiceProductInfo"
		exportCSVButtonHTMLArgument.buttonCSS = "button-service-product"
		exportCSVButtonHTML(w, exportCSVButtonHTMLArgument)
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table id=\"service-product-info-table\" class=\"table-service-product\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <th>Service/Product ID</th>")
		fmt.Fprintf(w, "          <th>Service/Product Name</th>")
		fmt.Fprintf(w, "          <th>Service/Product Type</th>")
		fmt.Fprintf(w, "          <th>Supplier Name</th>")
		fmt.Fprintf(w, "          <th>Supplier Contract Length</th>")
		fmt.Fprintf(w, "          <th>Date & Time Added</th>")
		fmt.Fprintf(w, "        </tr>")
		serviceProductInfoSQL, err := dbDetail.connection.Query(`SELECT
									   service_product_id,
									   service_product_name,
									   service_product_type,
									   service_product_supplier_name,
									   service_product_supplier_contract_length,
									   service_product_date_time_added
									 FROM
									   yap.view___service_product;`)

		// Error
		if err != nil {
			panic(err)

		}

		for serviceProductInfoSQL.Next() {

			err = serviceProductInfoSQL.Scan(
				&serviceProductID,
				&serviceProductName,
				&serviceProductType,
				&serviceProductSupplierName,
				&serviceProductSupplierContractLength,
				&serviceProductDateTimeAdded,
			)

			// Error
			if err != nil {
				panic(err)
			}

			fmt.Fprintf(w, "        <tr>")
			fmt.Fprintf(w, "          <td>"+serviceProductID+"</td>")
			fmt.Fprintf(w, "          <td>"+serviceProductName+"</td>")
			fmt.Fprintf(w, "          <td>"+serviceProductType+"</td>")
			fmt.Fprintf(w, "          <td>"+serviceProductSupplierName+"</td>")
			fmt.Fprintf(w, "          <td>"+serviceProductSupplierContractLength+"</td>")
			fmt.Fprintf(w, "          <td>"+formatDateTime(serviceProductDateTimeAdded)+"</td>")
			fmt.Fprintf(w, "        </tr>")
		}

		fmt.Fprintf(w, "      </table>")
		var filterTableJSArgument jsFunctionParameter
		filterTableJSArgument.tableID = "service-product-info-table"

		filterTableJSArgument.funcNameJS = "serviceProductSearchID"
		filterTableJSArgument.inputID = "service-product-input-id"
		filterTableJSArgument.columnNumber = 0
		filterTableJS(w, filterTableJSArgument)

		filterTableJSArgument.funcNameJS = "serviceProductSearchName"
		filterTableJSArgument.inputID = "service-product-input-name"
		filterTableJSArgument.columnNumber = 1
		filterTableJS(w, filterTableJSArgument)

		filterTableJSArgument.funcNameJS = "serviceProductSearchType"
		filterTableJSArgument.inputID = "service-product-input-type"
		filterTableJSArgument.columnNumber = 2
		filterTableJS(w, filterTableJSArgument)

		filterTableJSArgument.funcNameJS = "serviceProductSearchSupplierName"
		filterTableJSArgument.inputID = "service-product-input-supplier-name"
		filterTableJSArgument.columnNumber = 3
		filterTableJS(w, filterTableJSArgument)

		filterTableJSArgument.funcNameJS = "serviceProductSearchSupplierContractLength"
		filterTableJSArgument.inputID = "service-product-input-supplier-contract-length"
		filterTableJSArgument.columnNumber = 4
		filterTableJS(w, filterTableJSArgument)

		filterTableJSArgument.funcNameJS = "serviceProductSearchDateTime"
		filterTableJSArgument.inputID = "service-product-input-date-time"
		filterTableJSArgument.columnNumber = 5
		filterTableJS(w, filterTableJSArgument)

		var exportCSVJSArgument jsFunctionParameter
		exportCSVJSArgument.funcNameJS = "ServiceProductInfo"
		exportCSVJSArgument.tableID = "service-product-info-table"
		exportCSVJSArgument.fileName = "YAP_service_product_information_details"
		exportCSVJSArgument.pathURL = "service-product"
		exportCSVJS(w, exportCSVJSArgument)

		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")

		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<table id=\"table\" class=\"table-service-product\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>Service/Product Suppliers Available:</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td>")
		fmt.Fprintf(w, "      <table id=\"service-product-supplier-table\" class=\"table-service-product-supplier\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <th>Supplier Name</th>")
		fmt.Fprintf(w, "          <th>Date & Time Added</th>")
		fmt.Fprintf(w, "        </tr>")

		supplierSQL, err := dbDetail.connection.Query(`SELECT
                                                                 name,
                                                                 date_time_added
                                                               FROM
                                                                 yap.supplier;`)

		// Error
		if err != nil {
			panic(err)

		}

		for supplierSQL.Next() {

			err = supplierSQL.Scan(
				&supplierName,
				&supplierDateTimeAdded,
			)

			// Error
			if err != nil {
				panic(err)
			}

			fmt.Fprintf(w, "        <tr>")
			fmt.Fprintf(w, "          <td>"+supplierName+"</td>")
			fmt.Fprintf(w, "          <td>"+formatDateTime(supplierDateTimeAdded)+"</td>")
			fmt.Fprintf(w, "        </tr>")
		}

		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </td>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")

		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<table id=\"table\" class=\"table-service-product\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>All Sales Tax Rates Used:</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td>")
		fmt.Fprintf(w, "      <table id=\"sales-tax-rate-table\" class=\"table-sales-tax-rate\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <th>Sales Tax Rates</th>")
		fmt.Fprintf(w, "        </tr>")

		salesTaxRateSQL, err := dbDetail.connection.Query(`SELECT
                                                                     sales_tax_rate
                                                                   FROM
                                                                     yap.sales_tax_rate_lookup;`)

		// Error
		if err != nil {
			panic(err)

		}

		for salesTaxRateSQL.Next() {

			err = salesTaxRateSQL.Scan(
				&salesTaxRate,
			)

			// Error
			if err != nil {
				panic(err)
			}

			fmt.Fprintf(w, "        <tr>")
			fmt.Fprintf(w, "          <td>"+salesTaxRate+"&#37</td>")
			fmt.Fprintf(w, "        </tr>")
		}

		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </td>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</div>")
		var toggleDivJSArgument jsFunctionParameter
		toggleDivJSArgument.funcNameJS = "toggleServiceProduct"
		toggleDivJSArgument.divID = "service-product-div"
		toggleDivJS(w, toggleDivJSArgument)

	} else {
		panic("productServiceList function shoud only be called with account type ID 100")
	}
}

// Service/product Add function
func serviceProductAdd(w http.ResponseWriter, r *http.Request, dbDetail databaseFunctionParameter, genDetail generalFunctionParameter) {

	// Only account type ID 100 should be able to use this function
	if genDetail.userTypeID == "100" {

		fmt.Fprintf(w, "<form method=\"POST\" action=\"/service-product\">")
		fmt.Fprintf(w, "<table class=\"table-add\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>Add a New Service/Product<br>(If the Supplier Does Not Exist, Add a New One Below First)</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_service_product_input_name", "Service/Product Name<br>(Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		dbDetail.table = "service_product_type_lookup"
		dbDetail.column = "service_product_type"
		serviceProductTypeList := singleColumnSlice(dbDetail)
		serviceProductTypeList = append([]string{""}, serviceProductTypeList...)
		selectSingleHTML(w, "add_service_product_select_type", "Service/Product Type<br>(Cannot Be Empty)", serviceProductTypeList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		dbDetail.table = "supplier"
		dbDetail.column = "name"
		supplierNameList := singleColumnSlice(dbDetail)
		supplierNameList = append([]string{""}, supplierNameList...)
		supplierNameList = supplierNameList[:len(supplierNameList)-1]
		selectSingleHTML(w, "add_service_product_select_supplier_name", "Supplier Name<br>(Cannot Be Empty)", supplierNameList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		dbDetail.table = "contract_length_lookup"
		dbDetail.column = "contract_length"
		contractLengthList := singleColumnSlice(dbDetail)
		contractLengthList = append([]string{""}, contractLengthList...)
		selectSingleHTML(w, "add_service_product_select_contract_length", "Supplier Contract Length<br>(Can Be Empty)", contractLengthList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><input type=\"submit\" value=\"Create Service/Product\"></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</form>")

		addServiceProductInputName := r.FormValue("add_service_product_input_name")
		addServiceProductSelectType := r.FormValue("add_service_product_select_type")
		addServiceProductSelectSupplierName := r.FormValue("add_service_product_select_supplier_name")
		addServiceProductSelectContractLength := r.FormValue("add_service_product_select_contract_length")

		// Validate the name
		validateServiceProductName := validateInput(addServiceProductInputName, "alphaNum")
		// Validate the type
		validateServiceProductType := slices.Contains(serviceProductTypeList, addServiceProductSelectType)
		// Validate the supplier name
		validateServiceProductSupplierName := slices.Contains(supplierNameList, addServiceProductSelectSupplierName)
		// Validate the contract length
		validateServiceProductContractLength := slices.Contains(contractLengthList, addServiceProductSelectContractLength)

		if addServiceProductInputName == "" && addServiceProductSelectType == "" && addServiceProductSelectSupplierName == "" && addServiceProductSelectContractLength == "" {
			// Do Nothing
		} else if validateServiceProductName == false || addServiceProductInputName == "" {
			messageHTML(w, validationMessageServiceProductName, "warning")
		} else if validateServiceProductType == false || addServiceProductSelectType == "" {
			messageHTML(w, validationMessageServiceProductType, "warning")
		} else if addServiceProductSelectSupplierName == "⊛ YAP (Yet Another PBX) ⊛" {
			messageHTML(w, validationMessageServiceProductYAP, "warning")
		} else if validateServiceProductSupplierName == false || addServiceProductSelectSupplierName == "" {
			messageHTML(w, validationMessageServiceProductSupplierName, "warning")
		} else if validateServiceProductContractLength == false {
			messageHTML(w, validationMessageContractLength, "warning")
		} else {

			dbDetail.table = "view___service_product"
			dbDetail.column = "service_product_name"
			dbDetail.columnWhere = "service_product_name"
			dbDetail.columnWhereValue = addServiceProductInputName

			checkServiceProductExist := selectWhere(dbDetail)

			if checkServiceProductExist == addServiceProductInputName {
				messageHTML(w, validationMessageServiceProductAlreadyExist, "warning")
			} else {

				dbDetail.connection.Query(`INSERT 
                                   INTO
                                 service_product (
                                   name,
                                   service_product_type,
                                   supplier_name,
                                   supplier_contract_length
                                 )
                                 VALUES(?, ?, ?, ?);`,
					addServiceProductInputName,
					addServiceProductSelectType,
					addServiceProductSelectSupplierName,
					nullSQL(addServiceProductSelectContractLength),
				)

				checkServiceProductCreated := selectWhere(dbDetail)

				if checkServiceProductCreated == addServiceProductInputName {
					messageHTML(w, validationMessageServiceProductCreated, "success")
				} else {
					messageHTML(w, validationMessageServiceProductNotCreated, "warning")
				}
			}
		}

		// Add new supplier code
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<form method=\"POST\" action=\"/service-product\">")
		fmt.Fprintf(w, "<table class=\"table-add\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>Add a New Supplier</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_supplier_input_name", "Supplier Name (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		confirmList := yesSlice()
		selectSingleHTML(w, "add_supplier_select_confirm", "yes to Confirm (Cannot Be Empty)", confirmList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><input type=\"submit\" value=\"Create Supplier\"></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</form>")

		addSupplierInputName := r.FormValue("add_supplier_input_name")
		addSupplierSelectConfirm := r.FormValue("add_supplier_select_confirm")

		// Validate the name
		validateSupplierName := validateInput(addSupplierInputName, "alphaNum")

		if addSupplierInputName == "" && addSupplierSelectConfirm == "" {
			// Do Nothing
		} else if validateSupplierName == false {
			messageHTML(w, validationMessageSupplierName, "warning")
		} else if validateSupplierName == false && addSupplierSelectConfirm == "yes" {
			messageHTML(w, validationMessageSupplierName, "warning")
		} else if validateSupplierName == true && addSupplierSelectConfirm != "yes" {
			messageHTML(w, validationMessageConfirmation, "warning")
		} else if addSupplierInputName == "⊛ YAP (Yet Another PBX) ⊛" {
			messageHTML(w, validationMessageSupplierYAP, "warning")
		} else if validateSupplierName == true && addSupplierSelectConfirm == "yes" {

			dbDetail.table = "supplier"
			dbDetail.column = "name"
			dbDetail.columnWhere = "name"
			dbDetail.columnWhereValue = addSupplierInputName

			checkSupplierExist := selectWhere(dbDetail)

			if checkSupplierExist == addSupplierInputName {
				messageHTML(w, validationMessageSupplierAlreadyExist, "warning")
			} else {

				dbDetail.connection.Query(`INSERT 
                                   INTO
                                 supplier (
                                   name
                                 )
                                 VALUES(?);`,
					addSupplierInputName,
				)

				checkSupplierCreated := selectWhere(dbDetail)

				if checkSupplierCreated == addSupplierInputName {
					messageHTML(w, validationMessageSupplierCreated, "success")
				} else {
					messageHTML(w, validationMessageSupplierNotCreated, "warning")
				}
			}
		}

		fmt.Fprintf(w, "<br>")
		// Add new sales tax rate code
		fmt.Fprintf(w, "<form method=\"POST\" action=\"/service-product\">")
		fmt.Fprintf(w, "<table class=\"table-add\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>Add a New Sales Tax Rate</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "add_sales_tax_rate_input_rate", "Sales Tax Rate (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		selectSingleHTML(w, "add_sales_tax_rate_select_confirm", "yes to Confirm (Cannot Be Empty)", confirmList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><input type=\"submit\" value=\"Create Sales Tax Rate\"></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</form>")

		addSalesTaxRateInputRate := r.FormValue("add_sales_tax_rate_input_rate")
		addSalesTaxRateSelectConfirm := r.FormValue("add_sales_tax_rate_select_confirm")

		// Validate the sales tax rate
		validateSalesTaxRate := validateInput(addSalesTaxRateInputRate, "tax")

		if addSalesTaxRateInputRate == "" && addSalesTaxRateSelectConfirm == "" {
			// Do Nothing
		} else if validateSalesTaxRate == false {
			messageHTML(w, validationMessageSalesTaxRate, "warning")
		} else if validateSalesTaxRate == false && addSalesTaxRateSelectConfirm == "yes" {
			messageHTML(w, validationMessageSalesTaxRate, "warning")
		} else if validateSalesTaxRate == true && addSalesTaxRateSelectConfirm != "yes" {
			messageHTML(w, validationMessageConfirmation, "warning")
		} else if validateSalesTaxRate == true && addSalesTaxRateSelectConfirm == "yes" {

			addSalesTaxRateInputRateFloat64 := stringToFloat64(addSalesTaxRateInputRate)
			addSalesTaxRateInputRateFloat64 = math.Round(addSalesTaxRateInputRateFloat64*100) / 100
			addSalesTaxRateInputRate = strconv.FormatFloat(addSalesTaxRateInputRateFloat64, 'f', 2, 64)

			dbDetail.table = "sales_tax_rate_lookup"
			dbDetail.column = "sales_tax_rate"
			dbDetail.columnWhere = "sales_tax_rate"
			dbDetail.columnWhereValue = addSalesTaxRateInputRate

			checkSalesTaxRateExist := selectWhere(dbDetail)

			if checkSalesTaxRateExist == addSalesTaxRateInputRate {
				messageHTML(w, validationMessageSalesTaxRateAlreadyExist, "warning")
			} else {

				dbDetail.connection.Query(`INSERT 
                                   INTO
                                 sales_tax_rate_lookup (
                                   sales_tax_rate
                                 )
                                 VALUES(?);`,
					addSalesTaxRateInputRate,
				)

				checkSalesTaxRateCreated := selectWhere(dbDetail)

				if checkSalesTaxRateCreated == addSalesTaxRateInputRate {
					messageHTML(w, validationMessageSalesTaxRateCreated, "success")
				} else {
					messageHTML(w, validationMessageSalesTaxRateNotCreated, "warning")
				}
			}
		}

	} else {
		panic("serviceProductAdd function shoud only be called with account type ID 100")
	}
}

// Service/Product edit function
func serviceProductEdit(w http.ResponseWriter, r *http.Request, dbDetail databaseFunctionParameter, genDetail generalFunctionParameter) {

	// Only account type ID 100 should be able to use this function
	if genDetail.userTypeID == "100" {

		// List of all column names to edit in the service_product table
		serviceProductColumnList := [][]string{
			{"name", "Service/Product Name"},
			{"service_product_type", "Service/Product Type"},
			{"supplier_name", "Supplier Name"},
			{"supplier_contract_length", "Supplier Contract Length"},
		}

		fmt.Fprintf(w, "<form method=\"POST\" action=\"/service-product\">")
		fmt.Fprintf(w, "<table class=\"table-service-product\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>Edit Service/Product Details</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td style=\"text-align: left;\">")
		fmt.Fprintf(w, "      <b><u>Acceptable Values for Columns</u></b><br><br>")
		fmt.Fprintf(w, "      <b>Service/Product Name:</b> text<br>")
		fmt.Fprintf(w, "      <b>Service/Product Type:</b> Services, Products<br>")
		fmt.Fprintf(w, "      <b>Supplier Name:</b> text<br>")
		dbDetail.table = "contract_length_lookup"
		dbDetail.column = "contract_length"
		contractLengthList := singleColumnSlice(dbDetail)
		fmt.Fprintf(w, "      <b>Supplier Contract Length: </b>EMPTY, ")
		fmt.Fprintf(w, strings.Join(contractLengthList, ", "))
		fmt.Fprintf(w, "      <br>")
		fmt.Fprintf(w, "    </td>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "edit_service_product_input_id", "Service/Product (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		selectDoubleHiddenHTML(w, "edit_service_product_select_column", "Column to Edit (Cannot Be Empty)", serviceProductColumnList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "edit_service_product_input_new_value", "New Value")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><input type=\"submit\" value=\"Update Service/Product\"></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</form>")

		editServiceProductInputID := r.FormValue("edit_service_product_input_id")
		editServiceProductSelectColumn := r.FormValue("edit_service_product_select_column")
		editServiceProductInputNewValue := r.FormValue("edit_service_product_input_new_value")

		// Validate service/product ID
		validateServiceProductID := validateInput(editServiceProductInputID, "number")

		if editServiceProductInputID == "" && editServiceProductSelectColumn == "" && editServiceProductInputNewValue == "" {
			// Do Nothing
		} else if validateServiceProductID == false || editServiceProductInputID == "" {
			messageHTML(w, validationMessageServiceProductID, "warning")
		} else if editServiceProductSelectColumn == "" {
			messageHTML(w, validationMessageServiceProductColumn, "warning")
		} else if editServiceProductInputNewValue == "⊛ YAP (Yet Another PBX) ⊛" || editServiceProductInputNewValue == "⊛ YAP PBX Setup ⊛" || editServiceProductInputNewValue == "⊛ YAP PBX Rental ⊛" || editServiceProductInputNewValue == "⊛ YAP PBX Cease ⊛" || editServiceProductInputNewValue == "⊛ YAP Extension Setup ⊛" || editServiceProductInputNewValue == "⊛ YAP Extension Rental ⊛" || editServiceProductInputNewValue == "⊛ YAP Extension Cease ⊛" {
			messageHTML(w, validationMessageServiceProductYAP, "warning")
		} else if editServiceProductSelectColumn == "name" || editServiceProductSelectColumn == "service_product_type" || editServiceProductSelectColumn == "supplier_name" {
			// Validate editServiceProductInputNewValue is a string and not empty
			validateNewValue := validateInput(editServiceProductInputNewValue, "alphaNum")
			if validateNewValue == true {
				dbDetail.connection.Query("UPDATE service_product SET "+editServiceProductSelectColumn+" = ? WHERE id = ?;", editServiceProductInputNewValue, editServiceProductInputID)
			} else {
				messageHTML(w, validationMessageGenericAlphaNumEmpty, "warning")
			}
		} else if editServiceProductSelectColumn == "supplier_contract_length" {
			// Validate editServiceProductInputNewValue is a string
			validateNewValue := validateInput(editServiceProductInputNewValue, "alphaNumEmpty")
			if validateNewValue == true {
				dbDetail.connection.Query("UPDATE service_product SET "+editServiceProductSelectColumn+" = ? WHERE id = ?;", editServiceProductInputNewValue, editServiceProductInputID)
			} else {
				messageHTML(w, validationMessageGenericAlphaNum, "warning")
			}
		} else {
			messageHTML(w, validationMessageServiceProductColumn, "warning")
		}

		// Supplier edit code

		// List of column names to edit in the supplier table
		supplierColumnList := [][]string{
			{"name", "Supplier Name"},
		}

		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<form method=\"POST\" action=\"/service-product\">")
		fmt.Fprintf(w, "<table class=\"table-service-product\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>Edit Supplier Details</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td style=\"text-align: left;\">")
		fmt.Fprintf(w, "      <b><u>Acceptable Values for Columns</u></b><br><br>")
		fmt.Fprintf(w, "      <b>Supplier Name:</b> text<br>")
		fmt.Fprintf(w, "    </td>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "edit_supplier_input_existing_value", "Supplier Name (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		selectDoubleHiddenHTML(w, "edit_supplier_select_column", "Column to Edit (Cannot Be Empty)", supplierColumnList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "edit_supplier_input_new_value", "New Value (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><input type=\"submit\" value=\"Update Supplier\"></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</form>")

		editSupplierInputExistingValue := r.FormValue("edit_supplier_input_existing_value")
		editSupplierSelectColumn := r.FormValue("edit_supplier_select_column")
		editSupplierInputNewValue := r.FormValue("edit_supplier_input_new_value")

		// Validate supplier existing value
		validateSupplierExistingValue := validateInput(editSupplierInputExistingValue, "alphaNum")

		if editSupplierInputExistingValue == "" && editSupplierSelectColumn == "" && editSupplierInputNewValue == "" {
			// Do Nothing
		} else if validateSupplierExistingValue == false || editSupplierInputExistingValue == "" {
			messageHTML(w, validationMessageSupplierExistingValue, "warning")
		} else if editSupplierSelectColumn == "" {
			messageHTML(w, validationMessageSupplierColumn, "warning")
		} else if editSupplierInputNewValue == "⊛ YAP (Yet Another PBX) ⊛" {
			messageHTML(w, validationMessageSupplierYAP, "warning")
		} else if editSupplierSelectColumn == "name" {
			// Validate editSupplierInputNewValue is a string and not empty
			validateNewValue := validateInput(editSupplierInputNewValue, "alphaNum")
			if validateNewValue == true {
				dbDetail.connection.Query("UPDATE supplier SET "+editSupplierSelectColumn+" = ? WHERE name = ?;", editSupplierInputNewValue, editSupplierInputExistingValue)
			} else {
				messageHTML(w, validationMessageGenericAlphaNum, "warning")
			}
		} else {
			messageHTML(w, validationMessageSupplierColumn, "warning")
		}

		// Sales Tax Rate name edit code

		// List of column names to edit in the sales_tax_rate_lookup table
		salesTaxRateColumnList := [][]string{
			{"sales_tax_rate", "Sales Tax Rate"},
		}

		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<form method=\"POST\" action=\"/service-product\">")
		fmt.Fprintf(w, "<table class=\"table-service-product\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>Edit Sales Tax Rate</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td style=\"text-align: left;\">")
		fmt.Fprintf(w, "      <b><u>Acceptable Values for Columns</u></b><br><br>")
		fmt.Fprintf(w, "      <b>Sales Tax Rate:</b> decimal number<br>")
		fmt.Fprintf(w, "    </td>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "edit_sales_tax_rate_input_existing_value", "Sales Tax Rate (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		selectDoubleHiddenHTML(w, "edit_sales_tax_rate_select_column", "Column to Edit (Cannot Be Empty)", salesTaxRateColumnList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "edit_sales_tax_rate_input_new_value", "New Value (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><input type=\"submit\" value=\"Update Sales Tax Rate\"></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</form>")

		editSalesTaxRateInputExistingValue := r.FormValue("edit_sales_tax_rate_input_existing_value")
		editSalesTaxRateSelectColumn := r.FormValue("edit_sales_tax_rate_select_column")
		editSalesTaxRateInputNewValue := r.FormValue("edit_sales_tax_rate_input_new_value")

		// Validate sales tax rate existing value
		validateSalesTaxRateExistingValue := validateInput(editSalesTaxRateInputExistingValue, "tax")

		if editSalesTaxRateInputExistingValue == "" && editSalesTaxRateSelectColumn == "" && editSalesTaxRateInputNewValue == "" {
			// Do Nothing
		} else if validateSalesTaxRateExistingValue == false {
			messageHTML(w, validationMessageSalesTaxRate, "warning")
		} else if editSalesTaxRateSelectColumn == "" {
			messageHTML(w, validationMessageSalesTaxRateColumn, "warning")
		} else if editSalesTaxRateSelectColumn == "sales_tax_rate" {
			// Validate editSalesTaxRateInputNewValue is a tax value and not empty
			validateNewValue := validateInput(editSalesTaxRateInputNewValue, "tax")
			editSalesTaxRateInputNewValueFloat64 := stringToFloat64(editSalesTaxRateInputNewValue)
			editSalesTaxRateInputExistingValueFloat64 := stringToFloat64(editSalesTaxRateInputExistingValue)
			if validateNewValue == true {
				dbDetail.connection.Query("UPDATE sales_tax_rate_lookup SET "+editSalesTaxRateSelectColumn+" = ? WHERE sales_tax_rate = ?;", math.Round(editSalesTaxRateInputNewValueFloat64*100)/100, math.Round(editSalesTaxRateInputExistingValueFloat64*100)/100)
			} else {
				messageHTML(w, validationMessageSalesTaxRateTax, "warning")
			}
		} else {
			messageHTML(w, validationMessageSalesTaxRateColumn, "warning")
		}
	} else {
		panic("serviceProductEdit function should only be called with account type ID 100")
	}
}

// Service/product delete function
func serviceProductDelete(w http.ResponseWriter, r *http.Request, dbDetail databaseFunctionParameter, genDetail generalFunctionParameter) {

	// Only account type ID 100 should be able to use this function
	if genDetail.userTypeID == "100" {

		// Delete a service/product
		fmt.Fprintf(w, "<form method=\"POST\" action=\"/service-product\">")
		fmt.Fprintf(w, "<table class=\"table-delete\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>Delete a Service/Product</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "delete_service_product_input_name", "Service/Product (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		confirmList := yesSlice()
		selectSingleHTML(w, "delete_service_product_select_confirm", "yes to Confirm (Cannot Be Empty)", confirmList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><input type=\"submit\" value=\"Delete Service/Product\"></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</form>")

		deleteServiceProductInputName := r.FormValue("delete_service_product_input_name")
		deleteServiceProductSelectConfirm := r.FormValue("delete_service_product_select_confirm")

		// Validate service/product input
		validateServiceProductName := validateInput(deleteServiceProductInputName, "alphaNum")

		if deleteServiceProductInputName == "" && deleteServiceProductSelectConfirm == "" {
			// Do Nothing
		} else if deleteServiceProductInputName == "" {
			messageHTML(w, "", "warning")
		} else if validateServiceProductName == false && deleteServiceProductSelectConfirm == "yes" {
			messageHTML(w, "", "warning")
		} else if validateServiceProductName == true && deleteServiceProductSelectConfirm != "yes" {
			messageHTML(w, "", "warning")
		} else if deleteServiceProductInputName == "⊛ YAP PBX Setup ⊛" || deleteServiceProductInputName == "⊛ YAP PBX Rental ⊛" || deleteServiceProductInputName == "⊛ YAP PBX Cease ⊛" || deleteServiceProductInputName == "⊛ YAP Extension Setup ⊛" || deleteServiceProductInputName == "⊛ YAP Extension Rental ⊛" || deleteServiceProductInputName == "⊛ YAP Extension Cease ⊛" {
			messageHTML(w, "", "warning")
		} else if validateServiceProductName == true && deleteServiceProductSelectConfirm == "yes" {

			dbDetail.table = "view___service_product"
			dbDetail.column = "service_product_name"
			dbDetail.columnWhere = "service_product_name"
			dbDetail.columnWhereValue = deleteServiceProductInputName

			checkServiceProductExist := selectWhere(dbDetail)

			if checkServiceProductExist == "" {
				messageHTML(w, validationMessageServiceProductDoesNotExist, "warning")
			} else {

				dbDetail.connection.Query(`DELETE FROM service_product WHERE name = ?;`, deleteServiceProductInputName)

				checkServiceProductDeleted := selectWhere(dbDetail)

				if checkServiceProductDeleted == "" {
					messageHTML(w, validationMessageServiceProductDeleted, "success")
				} else {
					messageHTML(w, validationMessageServiceProductNotDeleted, "warning")
				}
			}

		} else {
			messageHTML(w, validationMessageInvalid, "warning")
		}

		// Delete a supplier code
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<form method=\"POST\" action=\"/service-product\">")
		fmt.Fprintf(w, "<table class=\"table-delete\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>Delete a Supplier</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "delete_supplier_input_name", "Supplier (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		selectSingleHTML(w, "delete_supplier_select_confirm", "yes to Confirm (Cannot Be Empty)", confirmList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><input type=\"submit\" value=\"Delete Supplier\"></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</form>")

		deleteSupplierInputName := r.FormValue("delete_supplier_input_name")
		deleteSupplierSelectConfirm := r.FormValue("delete_supplier_select_confirm")

		// Validate supplier input
		validateSupplierName := validateInput(deleteSupplierInputName, "alphaNum")

		if deleteSupplierInputName == "" && deleteSupplierSelectConfirm == "" {
			// Do Nothing
		} else if deleteSupplierInputName == "" {
			messageHTML(w, validationMessageSupplierName, "warning")
		} else if validateSupplierName == false && deleteSupplierSelectConfirm == "yes" {
			messageHTML(w, validationMessageSupplierName, "warning")
		} else if validateSupplierName == true && deleteSupplierSelectConfirm != "yes" {
			messageHTML(w, validationMessageConfirmation, "warning")
		} else if deleteSupplierInputName == "⊛ YAP (Yet Another PBX) ⊛" {
			messageHTML(w, validationMessageSupplierYAP, "warning")
		} else if validateSupplierName == true && deleteSupplierSelectConfirm == "yes" {

			dbDetail.table = "supplier"
			dbDetail.column = "name"
			dbDetail.columnWhere = "name"
			dbDetail.columnWhereValue = deleteSupplierInputName

			checkSupplierExist := selectWhere(dbDetail)

			if checkSupplierExist == "" {
				messageHTML(w, validationMessageSupplierDoesNotExist, "warning")
			} else {

				dbDetail.connection.Query(`DELETE FROM supplier WHERE name = ?;`, deleteSupplierInputName)

				checkSupplierDeleted := selectWhere(dbDetail)

				if checkSupplierDeleted == "" {
					messageHTML(w, validationMessageSupplierDeleted, "success")
				} else {
					messageHTML(w, validationMessageSupplierNotDeleted, "warning")
				}
			}

		} else {
			messageHTML(w, validationMessageInvalid, "warning")
		}

		// Delete a sales tax rate code
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<form method=\"POST\" action=\"/service-product\">")
		fmt.Fprintf(w, "<table class=\"table-delete\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th class=\"table-title\";>Delete Sales Tax Rate</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>")
		inputHTML(w, "delete_sales_tax_rate_input_rate", "Sales Tax Rate (Cannot Be Empty)")
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "          <td>")
		selectSingleHTML(w, "delete_sales_tax_rate_select_confirm", "yes to Confirm (Cannot Be Empty)", confirmList)
		fmt.Fprintf(w, "          </td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><input type=\"submit\" value=\"Delete Sales Tax Rate\"></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</form>")

		deleteSalesTaxRateInputRate := r.FormValue("delete_sales_tax_rate_input_rate")
		deleteSalesTaxRateSelectConfirm := r.FormValue("delete_sales_tax_rate_select_confirm")

		// Validate supplier input
		validateSalesTaxRate := validateInput(deleteSalesTaxRateInputRate, "tax")

		if deleteSalesTaxRateInputRate == "" && deleteSalesTaxRateSelectConfirm == "" {
			// Do Nothing
		} else if deleteSalesTaxRateInputRate == "" {
			messageHTML(w, validationMessageSalesTaxRate, "warning")
		} else if validateSalesTaxRate == false && deleteSalesTaxRateSelectConfirm == "yes" {
			messageHTML(w, validationMessageSalesTaxRate, "warning")
		} else if validateSalesTaxRate == true && deleteSalesTaxRateSelectConfirm != "yes" {
			messageHTML(w, validationMessageConfirmation, "warning")
		} else if validateSalesTaxRate == true && deleteSalesTaxRateSelectConfirm == "yes" {

			dbDetail.table = "sales_tax_rate_lookup"
			dbDetail.column = "sales_tax_rate"
			dbDetail.columnWhere = "sales_tax_rate"
			dbDetail.columnWhereValue = deleteSalesTaxRateInputRate

			checkSalesTaxRateExist := selectWhere(dbDetail)

			if checkSalesTaxRateExist == "" {
				messageHTML(w, validationMessageSalesTaxRateDoesNotExist, "warning")
			} else {

				dbDetail.connection.Query(`DELETE FROM sales_tax_rate_lookup WHERE sales_tax_rate = ?;`, deleteSalesTaxRateInputRate)

				checkSalesTaxRateDeleted := selectWhere(dbDetail)

				if checkSalesTaxRateDeleted == "" {
					messageHTML(w, validationMessageSalesTaxRateDeleted, "success")
				} else {
					messageHTML(w, validationMessageSalesTaxRateNotDeleted, "warning")
				}
			}

		} else {
			messageHTML(w, validationMessageInvalid, "warning")
		}
	} else {
		panic("serviceProductDelete function should only be called with account type ID 100")
	}
}

//----------------------------------------------------------------------------------------------------

func main() {

	// Get the values from inside the YAP configuration file
	err := godotenv.Load("/etc/yap/yap.env")
	if err != nil {
		panic("Error loading yap.env file for database details")
	}

	// Get the database connection details
	dbUsername := os.Getenv("dbUsername")
	dbPassword := os.Getenv("dbPassword")
	dbName := os.Getenv("dbName")
	dbTransport := os.Getenv("dbTransport")
	dbAddress := os.Getenv("dbAddress")
	dbPort := os.Getenv("dbPort")
	dbTLS := os.Getenv("dbTLS")
	defaultExtLimit := os.Getenv("defaultExtLimit")
	currencySymbol := os.Getenv("currencySymbol")
	yapAdminUKVATRegistered := os.Getenv("yapAdminUKVATRegistered")
	extraButtonName := os.Getenv("extraButtonName")
	extraButtonURL := os.Getenv("extraButtonURL")
	/*
		accountingSoftwareURL :=
		accountingSoftwareClientID :=
		accountingSoftwareClientSecret :=
		accountingSoftwareRefreshToken :=
		accountingSoftwareCurrencyCode :=
	*/

	// Values allowed for dbTransport Variable
	var transportList = []string{"tcp", "udp"}
	validDbTransport := slices.Contains(transportList, dbTransport)

	validateDbAddress := validator.New()
	validateDbAddressErr := validateDbAddress.Var(dbAddress, "required,ip_addr")

	dbPortInt, err := strconv.Atoi(dbPort)
	if err != nil {
		panic("DATABASE PORT MUST BE A NUMBER IN /etc/yap/yap.env")
	}

	// Values allowed for dbTls Variable
	var dbTLSList = []string{"false", "true"}
	validDbTLS := slices.Contains(dbTLSList, dbTLS)

	// Values allowed for defaultExtLimit Variable
	var defaultExtLimitList = []string{"1", "2", "3", "4", "5", "10", "25", "50", "75", "100", "150", "200", "250", "500", "750", "1000", "1500", "2000", "2500", "5000"}
	validDefaultExtLimit := slices.Contains(defaultExtLimitList, defaultExtLimit)

	// Values allowed for currencySymbol Variable
	var currencySymbolList = []string{"", "£", "€", "$", "¥"}
	validCurrencySymbol := slices.Contains(currencySymbolList, currencySymbol)

	// Values allowed for ukVATRegistered Variable
	var yapAdminUKVATRegisteredList = []string{"no", "yes"}
	validYAPAdminUKVATRegistered := slices.Contains(yapAdminUKVATRegisteredList, yapAdminUKVATRegistered)

	validateExtraButtonURL := validator.New()
	validateExtraButtonURLErr := validateExtraButtonURL.Var(extraButtonURL, "required,http_url")

	// Catch if any errors were made in yap.env and feed back where to correct the error
	if dbUsername == "" {
		panic("DATABASE USERNAME CANNOT BE EMPTY IN /etc/yap/yap.env")
	} else if dbPassword == "" {
		panic("DATABASE PASSOWRD CANNOT BE EMPTY IN /etc/yap/yap.env")
	} else if dbName == "" {
		panic("DATABASE NAME CANNOT BE EMPTY IN /etc/yap/yap.env")
	} else if dbTransport == "" {
		panic("DATABASE TRANSPORT OPTION CANNOT BE EMPTY IN /etc/yap/yap.env")
	} else if validDbTransport == false {
		panic("DATABASE TRANSPORT OPTION MUST BE udp OR tcp IN /etc/yap/yap.env")
	} else if validateDbAddressErr != nil && dbAddress != "localhost" {
		panic("DATABASE ADDRESS MUST BE A VALID INTERENT PROTOCOL (IP) ADDRESS OR localhost IN /etc/yap/yap.env")
	} else if dbPortInt <= 0 || dbPortInt >= 65536 {
		panic("DATABASE PORT MUST BE IN THE NUMBER RANGE 1-65535 IN /etc/yap/yap.env")
	} else if dbTLS == "" {
		panic("DATABASE TLS OPTION CANNOT BE EMPTY IN /etc/yap/yap.env")
	} else if validDbTLS == false {
		panic("DATABASE TRANSPORT OPTION MUST BE false OR true IN /etc/yap/yap.env")
	} else if defaultExtLimit == "" {
		panic("DEFAULT SIP EXT OPTION CANNOT BE EMPTY IN /etc/yap/yap.env")
	} else if validDefaultExtLimit == false {
		panic(" DEFAULT SIP EXT OPTION MUST BE SET TO A VALID OPTION IN /etc/yap/yap.env\nVALID OPTIONS: 1, 2, 3, 4, 5, 10, 25, 50, 75, 100, 150, 200, 250, 500, 750, 1000, 1500, 2000, 2500, 5000")
	} else if validCurrencySymbol == false {
		panic(" CURRENCY SYMBOL OPTION MUST BE SET TO £, €, $, ¥ OR EMPTY IN /etc/yap/yap.env")
	} else if yapAdminUKVATRegistered == "" {
		panic("UK YAP ADMIN VAT REGISTERED OPTION CANNOT BE EMPTY IN /etc/yap/yap.env")
	} else if validYAPAdminUKVATRegistered == false {
		panic("UK YAP ADMIN VAT REGISTERED OPTION MUST BE no OR yes IN /etc/yap/yap.env")
	} else if extraButtonName != "" {
		if validateExtraButtonURLErr != nil {
			panic("THE EXTRA BUTTON URL VALUE MUST BE A VALID URL IN /etc/yap/yap.env")
		}
	}

	startHTML := csvcell.FileData(dirHTML, fileStartHTML)
	endHTML := csvcell.FileData(dirHTML, fileEndHTML)

	// Home Page
	go http.HandleFunc("/yap", func(w http.ResponseWriter, r *http.Request) {

		// Open database connection
		dbConnection, err := sql.Open("mysql", dbUsername+":"+dbPassword+"@"+dbTransport+"("+dbAddress+":"+dbPort+")/"+dbName+"?tls="+dbTLS)
		defer dbConnection.Close()

		// Error
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, startHTML)

		// Wallpaper
		wallpaper(w, "wallpaper-main-menu")

		// Code to call the emailHeaderHTTP function
		email := emailHeaderHTTP(r)

		var dbDetail databaseFunctionParameter
		dbDetail.connection = dbConnection
		dbDetail.database = dbName
		dbDetail.columnWhereValue = email

		userTypeID := userAccountData(dbDetail, "type_id")
		userCustomerID := userAccountData(dbDetail, "customer_id")
		userCustomerName := userAccountData(dbDetail, "customer_name")
		userPBXID := userAccountData(dbDetail, "pbx_id")
		userPBXName := userAccountData(dbDetail, "pbx_name")

		var genDetail generalFunctionParameter
		genDetail.userTypeID = userTypeID

		if userTypeID == "" {
			errorBox(w, "email_error", "header-main-menu", "")
		} else {
			if userTypeID == "100" {
				header(w, "YAP Admin Account<br>Main Menu", "", extraButtonName, extraButtonURL)
				mainMenuUserInformation(w, dbDetail, genDetail)
				fmt.Fprintf(w, "<br>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonOne := mainMenuParameter{writeHTTP: w, buttonName: "Customer & PBX<br>User Accounts<br>&#128100", hyperlink: "/user-account", headerCSS: "header-user-account", buttonCSS: "button-user-account"}
				mainMenuButton(mainMenuButtonOne)
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				mainMenuButtonTwo := mainMenuParameter{writeHTTP: w, buttonName: "All<br>Customers<br>&#128101", hyperlink: "/customer", headerCSS: "header-customer", buttonCSS: "button-customer"}
				mainMenuButton(mainMenuButtonTwo)
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				mainMenuButtonThree := mainMenuParameter{writeHTTP: w, buttonName: "All<br>PBXs<br>&#128222", hyperlink: "/pbx", headerCSS: "header-pbx", buttonCSS: "button-pbx"}
				mainMenuButton(mainMenuButtonThree)
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				mainMenuButtonFour := mainMenuParameter{writeHTTP: w, buttonName: "All PBX SIP<br>Extensions<br>&#128241", hyperlink: "/extension", headerCSS: "header-ext", buttonCSS: "button-ext"}
				mainMenuButton(mainMenuButtonFour)
				fmt.Fprintf(w, "</div>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonFive := mainMenuParameter{writeHTTP: w, buttonName: "All Customer<br>Invoicing<br>&#129534", hyperlink: "/invoice", headerCSS: "header-invoice", buttonCSS: "button-invoice"}
				mainMenuButton(mainMenuButtonFive)
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				mainMenuButtonSix := mainMenuParameter{writeHTTP: w, buttonName: "Accounting<br>Software<br>&#128183 &#128182 &#128181 &#128180", hyperlink: "/accounting-software", headerCSS: "header-accounting-software", buttonCSS: "button-accounting-software"}
				mainMenuButton(mainMenuButtonSix)
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				mainMenuButtonSeven := mainMenuParameter{writeHTTP: w, buttonName: "All Services &<br>Products<br>&#128230", hyperlink: "/service-product", headerCSS: "header-service-product", buttonCSS: "button-service-product"}
				mainMenuButton(mainMenuButtonSeven)
				fmt.Fprintf(w, "</div>")
				footer(w, "", "")
			} else if userTypeID == "200" {
				header(w, userCustomerName+"<br>[Customer ID: "+userCustomerID+"]<br>Main Menu", "", extraButtonName, extraButtonURL)
				mainMenuUserInformation(w, dbDetail, genDetail)
				fmt.Fprintf(w, "<br>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonOne := mainMenuParameter{writeHTTP: w, buttonName: "Customer & PBX<br>User Accounts<br>&#128100", hyperlink: "/user-account", headerCSS: "header-user-account", buttonCSS: "button-user-account"}
				mainMenuButton(mainMenuButtonOne)
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				mainMenuButtonTwo := mainMenuParameter{writeHTTP: w, buttonName: "Customer<br>Information<br>&#128101", hyperlink: "/customer", headerCSS: "header-customer", buttonCSS: "button-customer"}
				mainMenuButton(mainMenuButtonTwo)
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				mainMenuButtonThree := mainMenuParameter{writeHTTP: w, buttonName: "PBXs for the<br>Customer<br>&#128222", hyperlink: "/pbx", headerCSS: "header-pbx", buttonCSS: "button-pbx"}
				mainMenuButton(mainMenuButtonThree)
				fmt.Fprintf(w, "</div>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonFour := mainMenuParameter{writeHTTP: w, buttonName: "PBX<br>Extensions<br>&#128241", hyperlink: "/extension", headerCSS: "header-ext", buttonCSS: "button-ext"}
				mainMenuButton(mainMenuButtonFour)
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				mainMenuButtonFive := mainMenuParameter{writeHTTP: w, buttonName: "Customer<br>Invoice<br>&#129534", hyperlink: "/invoice", headerCSS: "header-invoice", buttonCSS: "button-invoice"}
				mainMenuButton(mainMenuButtonFive)
				fmt.Fprintf(w, "</div>")
				footer(w, "", "")
			} else if userTypeID == "201" {
				header(w, userCustomerName+"<br>[Customer ID: "+userCustomerID+"]<br>Main Menu", "", extraButtonName, extraButtonURL)
				mainMenuUserInformation(w, dbDetail, genDetail)
				fmt.Fprintf(w, "<br>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonOne := mainMenuParameter{writeHTTP: w, buttonName: "Own & PBX<br>User Accounts<br>&#128100", hyperlink: "/user-account", headerCSS: "header-user-account", buttonCSS: "button-user-account"}
				mainMenuButton(mainMenuButtonOne)
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				mainMenuButtonTwo := mainMenuParameter{writeHTTP: w, buttonName: "Customer<br>Information<br>&#128101", hyperlink: "/customer", headerCSS: "header-customer", buttonCSS: "button-customer"}
				mainMenuButton(mainMenuButtonTwo)
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				mainMenuButtonThree := mainMenuParameter{writeHTTP: w, buttonName: "PBXs for the<br>Customer<br>&#128222", hyperlink: "/pbx", headerCSS: "header-pbx", buttonCSS: "button-pbx"}
				mainMenuButton(mainMenuButtonThree)
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				mainMenuButtonFour := mainMenuParameter{writeHTTP: w, buttonName: "PBX<br>Extensions<br>&#128241", hyperlink: "/extension", headerCSS: "header-ext", buttonCSS: "button-ext"}
				mainMenuButton(mainMenuButtonFour)
				fmt.Fprintf(w, "</div>")
				footer(w, "", "")
			} else if userTypeID == "300" {
				header(w, userPBXName+"<br>[PBX ID: "+userPBXID+"]<br>Main Menu", "", extraButtonName, extraButtonURL)
				mainMenuUserInformation(w, dbDetail, genDetail)
				fmt.Fprintf(w, "<br>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonOne := mainMenuParameter{writeHTTP: w, buttonName: "Own & PBX<br>User Accounts<br>&#128100", hyperlink: "/user-account", headerCSS: "header-user-account", buttonCSS: "button-user-account"}
				mainMenuButton(mainMenuButtonOne)
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				mainMenuButtonTwo := mainMenuParameter{writeHTTP: w, buttonName: "PBX<br>Information<br>&#128222", hyperlink: "/pbx", headerCSS: "header-pbx", buttonCSS: "button-pbx"}
				mainMenuButton(mainMenuButtonTwo)
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				mainMenuButtonThree := mainMenuParameter{writeHTTP: w, buttonName: "PBX<br>Extensions<br>&#128241", hyperlink: "/extension", headerCSS: "header-ext", buttonCSS: "button-ext"}
				mainMenuButton(mainMenuButtonThree)
				fmt.Fprintf(w, "</div>")
				footer(w, "", "")
			} else if userTypeID == "301" || userTypeID == "302" {
				header(w, userPBXName+"<br>[PBX ID: "+userPBXID+"]<br>Main Menu", "", extraButtonName, extraButtonURL)
				mainMenuUserInformation(w, dbDetail, genDetail)
				fmt.Fprintf(w, "<br>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonOne := mainMenuParameter{writeHTTP: w, buttonName: "Own<br>User Account<br>&#128100", hyperlink: "/user-account", headerCSS: "header-user-account", buttonCSS: "button-user-account"}
				mainMenuButton(mainMenuButtonOne)
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				mainMenuButtonTwo := mainMenuParameter{writeHTTP: w, buttonName: "PBX<br>Information<br>&#128222", hyperlink: "/pbx", headerCSS: "header-pbx", buttonCSS: "button-pbx"}
				mainMenuButton(mainMenuButtonTwo)
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				mainMenuButtonThree := mainMenuParameter{writeHTTP: w, buttonName: "PBX<br>Extensions<br>&#128241", hyperlink: "/extension", headerCSS: "header-ext", buttonCSS: "button-ext"}
				mainMenuButton(mainMenuButtonThree)
				fmt.Fprintf(w, "</div>")
				footer(w, "", "")
			} else if userTypeID == "400" {
				header(w, userCustomerName+"<br>[Customer ID: "+userCustomerID+"]<br>Main Menu", "", extraButtonName, extraButtonURL)
				mainMenuUserInformation(w, dbDetail, genDetail)
				fmt.Fprintf(w, "<br>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonOne := mainMenuParameter{writeHTTP: w, buttonName: "Own<br>User Account<br>&#128100", hyperlink: "/user-account", headerCSS: "header-user-account", buttonCSS: "button-user-account"}
				mainMenuButton(mainMenuButtonOne)
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				mainMenuButtonTwo := mainMenuParameter{writeHTTP: w, buttonName: "Customer<br>Information<br>&#128101", hyperlink: "/customer", headerCSS: "header-customer", buttonCSS: "button-customer"}
				mainMenuButton(mainMenuButtonTwo)
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				fmt.Fprintf(w, "&nbsp")
				mainMenuButtonThree := mainMenuParameter{writeHTTP: w, buttonName: "Customer<br>Invoice<br>&#129534", hyperlink: "/invoice", headerCSS: "header-invoice", buttonCSS: "button-invoice"}
				mainMenuButton(mainMenuButtonThree)
				fmt.Fprintf(w, "</div>")
				footer(w, "", "")
			} else {
				errorBox(w, "account_type_error", "header-main-menu", "")
			}
		}
		fmt.Fprintf(w, endHTML)

	})

	// User Account Page
	http.HandleFunc("/user-account", func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
		}

		// Open database connection
		dbConnection, err := sql.Open("mysql", dbUsername+":"+dbPassword+"@"+dbTransport+"("+dbAddress+":"+dbPort+")/"+dbName+"?tls="+dbTLS)
		defer dbConnection.Close()

		// Error
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, startHTML)

		// Wallpaper
		wallpaper(w, "wallpaper-user-account")

		// Code to call the emailHeaderHTTP function
		email := emailHeaderHTTP(r)

		var dbDetail databaseFunctionParameter
		dbDetail.connection = dbConnection
		dbDetail.database = dbName
		dbDetail.columnWhereValue = email

		userID := userAccountData(dbDetail, "id")
		userTypeID := userAccountData(dbDetail, "type_id")
		userCustomerID := userAccountData(dbDetail, "customer_id")
		userCustomerName := userAccountData(dbDetail, "customer_name")
		userPBXID := userAccountData(dbDetail, "pbx_id")
		userPBXName := userAccountData(dbDetail, "pbx_name")

		var genDetail generalFunctionParameter
		genDetail.userID = userID
		genDetail.userTypeID = userTypeID

		if userTypeID == "" {
			errorBox(w, "email_error", "header-user-account", "button-user-account")
		} else {
			if userTypeID == "100" {
				header(w, "YAP Admin Account<br>All User Accounts on YAP", "header-user-account", extraButtonName, extraButtonURL)
				userAccountList(w, dbDetail, genDetail)
				fmt.Fprint(w, "<br>")
				userAccountAdd(w, r, dbDetail, genDetail)
				fmt.Fprint(w, "<br>")
				userAccountEdit(w, r, dbDetail, genDetail)
				fmt.Fprint(w, "<br>")
				userAccountDelete(w, r, dbDetail, genDetail)
				footer(w, "header-user-account", "button-user-account")
			} else if userTypeID == "200" {
				header(w, userCustomerName+"<br>[Customer ID: "+userCustomerID+"]<br>All User Accounts for the Customer", "header-user-account", extraButtonName, extraButtonURL)
				userAccountList(w, dbDetail, genDetail)
				footer(w, "header-user-account", "button-user-account")
			} else if userTypeID == "201" {
				header(w, userCustomerName+"<br>[Customer ID: "+userCustomerID+"]<br>All PBX User Accounts for the Customer", "header-user-account", extraButtonName, extraButtonURL)
				userAccountList(w, dbDetail, genDetail)
				footer(w, "header-user-account", "button-user-account")
			} else if userTypeID == "300" {
				header(w, userPBXName+"<br>[PBX ID: "+userPBXID+"]<br>All User Accounts Within the PBX", "header-user-account", extraButtonName, extraButtonURL)
				userAccountList(w, dbDetail, genDetail)
				footer(w, "header-user-account", "button-user-account")
			} else if userTypeID == "301" {
				header(w, userPBXName+"<br>[PBX ID: "+userPBXID+"]<br>Own User Account for PBX", "header-user-account", extraButtonName, extraButtonURL)
				userAccountList(w, dbDetail, genDetail)
				footer(w, "header-user-account", "button-user-account")
			} else if userTypeID == "302" {
				header(w, userPBXName+"<br>[PBX ID: "+userPBXID+"]<br>Own Read Only User Account for PBX", "header-user-account", extraButtonName, extraButtonURL)
				userAccountList(w, dbDetail, genDetail)
				footer(w, "header-user-account", "button-user-account")
			} else if userTypeID == "400" {
				header(w, userCustomerName+"<br>[Customer ID: "+userCustomerID+"]<br>Own Invoice Customer Account", "header-user-account", extraButtonName, extraButtonURL)
				userAccountList(w, dbDetail, genDetail)
				footer(w, "header-user-account", "button-user-account")
			} else {
				errorBox(w, "account_type_error", "header-user-account", "button-user-account")
			}
		}
		fmt.Fprintf(w, endHTML)
	})

	// Customer Page
	go http.HandleFunc("/customer", func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
		}

		// Open database connection
		dbConnection, err := sql.Open("mysql", dbUsername+":"+dbPassword+"@"+dbTransport+"("+dbAddress+":"+dbPort+")/"+dbName+"?tls="+dbTLS)
		defer dbConnection.Close()

		// Error
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, startHTML)

		// Wallpaper
		wallpaper(w, "wallpaper-customer")

		// Code to call the emailHeaderHTTP function
		email := emailHeaderHTTP(r)

		var dbDetail databaseFunctionParameter
		dbDetail.connection = dbConnection
		dbDetail.database = dbName
		dbDetail.columnWhereValue = email

		userTypeID := userAccountData(dbDetail, "type_id")
		userCustomerID := userAccountData(dbDetail, "customer_id")
		userCustomerName := userAccountData(dbDetail, "customer_name")

		var genDetail generalFunctionParameter
		genDetail.userTypeID = userTypeID
		genDetail.userCustomerID = userCustomerID
		genDetail.defaultExtLimit = defaultExtLimit
		genDetail.currencySymbol = currencySymbol

		if userTypeID == "" {
			errorBox(w, "email_error", "header-customer", "button-customer")
		} else {
			if userTypeID == "100" {
				header(w, "YAP Admin Account<br>All Customers on YAP", "header-customer", extraButtonName, extraButtonURL)
				customerList(w, dbDetail, genDetail)
				fmt.Fprint(w, "<br>")
				customerAdd(w, r, dbDetail, genDetail)
				fmt.Fprint(w, "<br>")
				customerEdit(w, r, dbDetail, genDetail)
				fmt.Fprint(w, "<br>")
				customerDelete(w, r, dbDetail, genDetail)
				footer(w, "header-customer", "button-customer")
			} else if userTypeID == "200" || userTypeID == "201" || userTypeID == "400" {
				header(w, userCustomerName+"<br>[Customer ID: "+userCustomerID+"]<br>Own Customer Information", "header-customer", extraButtonName, extraButtonURL)
				customerList(w, dbDetail, genDetail)
				footer(w, "header-customer", "button-customer")
			} else {
				errorBox(w, "account_type_error", "header-customer", "button-customer")
			}
		}
		fmt.Fprintf(w, endHTML)
	})

	// PBX Page
	go http.HandleFunc("/pbx", func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
		}

		// Open database connection
		dbConnection, err := sql.Open("mysql", dbUsername+":"+dbPassword+"@"+dbTransport+"("+dbAddress+":"+dbPort+")/"+dbName+"?tls="+dbTLS)
		defer dbConnection.Close()

		// Error
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, startHTML)

		// Wallpaper
		wallpaper(w, "wallpaper-pbx")

		// Code to call the emailHeaderHTTP function
		email := emailHeaderHTTP(r)

		var dbDetail databaseFunctionParameter
		dbDetail.connection = dbConnection
		dbDetail.database = dbName
		dbDetail.columnWhereValue = email

		userTypeID := userAccountData(dbDetail, "type_id")
		userCustomerID := userAccountData(dbDetail, "customer_id")
		userCustomerName := userAccountData(dbDetail, "customer_name")
		userPBXID := userAccountData(dbDetail, "pbx_id")
		userPBXName := userAccountData(dbDetail, "pbx_name")

		var genDetail generalFunctionParameter
		genDetail.userTypeID = userTypeID
		genDetail.userCustomerID = userCustomerID
		genDetail.userPBXID = userPBXID
		genDetail.defaultExtLimit = defaultExtLimit

		if userTypeID == "" {
			errorBox(w, "email_error", "header-pbx", "button-pbx")
		} else {
			if userTypeID == "100" {
				header(w, "YAP Admin Account<br>All BPXs on YAP", "header-pbx", extraButtonName, extraButtonURL)
				pbxList(w, dbDetail, genDetail)
				fmt.Fprintf(w, "<br>")
				pbxAdd(w, r, dbDetail, genDetail)
				fmt.Fprintf(w, "<br>")
				pbxEdit(w, r, dbDetail, genDetail)
				fmt.Fprintf(w, "<br>")
				pbxDelete(w, r, dbDetail, genDetail)
				footer(w, "header-pbx", "button-pbx")
			} else if userTypeID == "200" {
				header(w, userCustomerName+"<br>[Customer ID: "+userCustomerID+"]<br>All PBXs for the Customer", "header-pbx", extraButtonName, extraButtonURL)
				pbxList(w, dbDetail, genDetail)
				fmt.Fprintf(w, "<br>")
				pbxAdd(w, r, dbDetail, genDetail)
				fmt.Fprintf(w, "<br>")
				pbxEdit(w, r, dbDetail, genDetail)
				fmt.Fprintf(w, "<br>")
				pbxDelete(w, r, dbDetail, genDetail)
				footer(w, "header-pbx", "button-pbx")
			} else if userTypeID == "201" {
				header(w, userCustomerName+"<br>[Customer ID: "+userCustomerID+"]<br>All PBXs for the Customer", "header-pbx", extraButtonName, extraButtonURL)
				pbxList(w, dbDetail, genDetail)
				fmt.Fprintf(w, "<br>")
				pbxAdd(w, r, dbDetail, genDetail)
				fmt.Fprintf(w, "<br>")
				pbxEdit(w, r, dbDetail, genDetail)
				footer(w, "header-pbx", "button-pbx")
			} else if userTypeID == "300" || userTypeID == "301" || userTypeID == "302" {
				header(w, userPBXName+"<br>[PBX ID: "+userPBXID+"]<br>PBX Information", "header-pbx", extraButtonName, extraButtonURL)
				pbxList(w, dbDetail, genDetail)
				footer(w, "header-pbx", "button-pbx")
			} else {
				errorBox(w, "account_type_error", "header-pbx", "button-pbx")
			}
		}
		fmt.Fprintf(w, endHTML)
	})

	// SIP Extension Page
	go http.HandleFunc("/extension", func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
		}

		// Open database connection
		dbConnection, err := sql.Open("mysql", dbUsername+":"+dbPassword+"@"+dbTransport+"("+dbAddress+":"+dbPort+")/"+dbName+"?tls="+dbTLS)
		defer dbConnection.Close()

		// Error
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, startHTML)

		// Wallpaper
		wallpaper(w, "wallpaper-ext")

		// Code to call the emailHeaderHTTP function
		email := emailHeaderHTTP(r)

		var dbDetail databaseFunctionParameter
		dbDetail.connection = dbConnection
		dbDetail.database = dbName
		dbDetail.columnWhereValue = email

		userTypeID := userAccountData(dbDetail, "type_id")
		userCustomerID := userAccountData(dbDetail, "customer_id")
		userCustomerName := userAccountData(dbDetail, "customer_name")
		userPBXID := userAccountData(dbDetail, "pbx_id")
		userPBXName := userAccountData(dbDetail, "pbx_name")

		var genDetail generalFunctionParameter
		genDetail.userTypeID = userTypeID
		genDetail.userCustomerID = userCustomerID
		genDetail.userPBXID = userPBXID

		if userTypeID == "" {
			errorBox(w, "email_error", "header-ext", "button-ext")
		} else {
			if userTypeID == "100" {
				header(w, "YAP Admin Account<br>All Extensions on the Server", "header-ext", extraButtonName, extraButtonURL)
				extList(w, dbDetail, genDetail)
				fmt.Fprintf(w, "<br>")
				extAdd(w, r, dbDetail, genDetail)
				fmt.Fprintf(w, "<br>")
				extEdit(w, r, dbDetail, genDetail)
				fmt.Fprintf(w, "<br>")
				extDelete(w, r, dbDetail, genDetail)
				footer(w, "header-ext", "button-ext")
			} else if userTypeID == "200" || userTypeID == "201" {
				header(w, userCustomerName+"<br>[Customer ID: "+userCustomerID+"]<br>All Extensions for the Customer", "header-ext", extraButtonName, extraButtonURL)
				extList(w, dbDetail, genDetail)
				fmt.Fprintf(w, "<br>")
				extAdd(w, r, dbDetail, genDetail)
				fmt.Fprintf(w, "<br>")
				extEdit(w, r, dbDetail, genDetail)
				fmt.Fprintf(w, "<br>")
				extDelete(w, r, dbDetail, genDetail)
				footer(w, "header-ext", "button-ext")
			} else if userTypeID == "300" {
				header(w, userPBXName+"<br>[PBX ID: "+userPBXID+"]<br>All Extensions Within the PBX", "header-ext", extraButtonName, extraButtonURL)
				extList(w, dbDetail, genDetail)
				fmt.Fprintf(w, "<br>")
				extAdd(w, r, dbDetail, genDetail)
				fmt.Fprintf(w, "<br>")
				extEdit(w, r, dbDetail, genDetail)
				fmt.Fprintf(w, "<br>")
				extDelete(w, r, dbDetail, genDetail)
				footer(w, "header-ext", "button-ext")
			} else if userTypeID == "301" {
				header(w, userPBXName+"<br>[PBX ID: "+userPBXID+"]<br>All Extensions Within the PBX", "header-ext", extraButtonName, extraButtonURL)
				extList(w, dbDetail, genDetail)
				fmt.Fprintf(w, "<br>")
				extAdd(w, r, dbDetail, genDetail)
				fmt.Fprintf(w, "<br>")
				extEdit(w, r, dbDetail, genDetail)
				footer(w, "header-ext", "button-ext")
			} else if userTypeID == "302" {
				header(w, userPBXName+"<br>[PBX ID: "+userPBXID+"]<br>All Extensions Within the PBX", "header-ext", extraButtonName, extraButtonURL)
				extList(w, dbDetail, genDetail)
				footer(w, "header-ext", "button-ext")
			} else {
				errorBox(w, "account_type_error", "header-ext", "button-ext")
			}
		}
		fmt.Fprintf(w, endHTML)
	})

	// Invoicing Page
	go http.HandleFunc("/invoice", func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
		}

		// Open database connection
		dbConnection, err := sql.Open("mysql", dbUsername+":"+dbPassword+"@"+dbTransport+"("+dbAddress+":"+dbPort+")/"+dbName+"?tls="+dbTLS)
		defer dbConnection.Close()

		// Error
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, startHTML)

		// Wallpaper
		wallpaper(w, "wallpaper-invoice")

		// Code to call the emailHeaderHTTP function
		email := emailHeaderHTTP(r)

		var dbDetail databaseFunctionParameter
		dbDetail.connection = dbConnection
		dbDetail.database = dbName
		dbDetail.columnWhereValue = email

		userTypeID := userAccountData(dbDetail, "type_id")
		userCustomerID := userAccountData(dbDetail, "customer_id")
		userCustomerName := userAccountData(dbDetail, "customer_name")

		var genDetail generalFunctionParameter
		genDetail.userTypeID = userTypeID
		genDetail.userCustomerID = userCustomerID
		genDetail.currencySymbol = currencySymbol
		genDetail.yapAdminUKVATRegistered = yapAdminUKVATRegistered

		if userTypeID == "" {
			errorBox(w, "email_error", "header-invoice", "button-invoice")
		} else {
			if userTypeID == "100" {
				header(w, "YAP Admin Account<br>All Customer Invoices", "header-invoice", extraButtonName, extraButtonURL)
				invoiceList(w, dbDetail, genDetail)
				fmt.Fprintf(w, "<br>")
				invoiceAdd(w, r, dbDetail, genDetail)
				fmt.Fprintf(w, "<br>")
				invoiceDelete(w, r, dbDetail, genDetail)
				footer(w, "header-invoice", "button-invoice")
			} else if userTypeID == "200" || userTypeID == "400" {
				header(w, userCustomerName+"<br>[Customer ID: "+userCustomerID+"]<br>Customer Invoice", "header-invoice", extraButtonName, extraButtonURL)
				invoiceList(w, dbDetail, genDetail)
				footer(w, "header-invoice", "button-invoice")
			} else {
				errorBox(w, "account_type_error", "header-invoice", "button-invoice")
			}
		}
		fmt.Fprintf(w, endHTML)
	})

	go http.HandleFunc("/accounting-software", func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
		}

		// Open database connection
		dbConnection, err := sql.Open("mysql", dbUsername+":"+dbPassword+"@"+dbTransport+"("+dbAddress+":"+dbPort+")/"+dbName+"?tls="+dbTLS)
		defer dbConnection.Close()

		// Error
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, startHTML)

		// Wallpaper
		wallpaper(w, "wallpaper-accounting-software")

		// Code to call the emailHeaderHTTP function
		email := emailHeaderHTTP(r)

		var dbDetail databaseFunctionParameter
		dbDetail.connection = dbConnection
		dbDetail.database = dbName
		dbDetail.columnWhereValue = email

		userTypeID := userAccountData(dbDetail, "type_id")

		var genDetail generalFunctionParameter
		genDetail.userTypeID = userTypeID

		var accountingSoftware accountingSoftwareParameter
		accountingSoftware.url = ""
		accountingSoftware.clientID = ""
		accountingSoftware.clientSecret = ""
		accountingSoftware.refreshToken = ""
		accountingSoftware.currencyCode = ""

		if userTypeID == "" {
			errorBox(w, "email_error", "header-accounting-software", "button-accounting-software")
		} else {
			if userTypeID == "100" {
				if accountingSoftware.url == "" {
					errorBox(w, "url_error", "header-accounting-software", "button-accounting-software")
				} else if accountingSoftware.clientID == "" {
					errorBox(w, "client_id_error", "header-accounting-software", "button-accounting-software")
				} else if accountingSoftware.clientSecret == "" {
					errorBox(w, "client_secret_error", "header-accounting-software", "button-accounting-software")
				} else if accountingSoftware.refreshToken == "" {
					errorBox(w, "refresh_token_error", "header-accounting-software", "button-accounting-software")
				} else if accountingSoftware.currencyCode == "" {
					errorBox(w, "currency_code_error", "header-accounting-software", "button-accounting-software")
				} else {
					header(w, "YAP Admin Account<br>Send Invoices/Customer Details", "header-accounting-software", extraButtonName, extraButtonURL)
					sendCustomerInvoice(dbDetail, accountingSoftware)
					footer(w, "header-accounting-software", "button-accounting-software")
				}
			} else {
				errorBox(w, "account_type_error", "header-accounting-software", "button-accounting-software")
			}
		}
		fmt.Fprintf(w, endHTML)
	})

	// Service product Page
	http.HandleFunc("/service-product", func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
		}

		// Open database connection
		dbConnection, err := sql.Open("mysql", dbUsername+":"+dbPassword+"@"+dbTransport+"("+dbAddress+":"+dbPort+")/"+dbName+"?tls="+dbTLS)
		defer dbConnection.Close()

		// Error
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, startHTML)

		// Wallpaper
		wallpaper(w, "wallpaper-service-product")

		// Code to call the emailHeaderHTTP function
		email := emailHeaderHTTP(r)

		var dbDetail databaseFunctionParameter
		dbDetail.connection = dbConnection
		dbDetail.database = dbName
		dbDetail.columnWhereValue = email

		userTypeID := userAccountData(dbDetail, "type_id")
		userCustomerID := userAccountData(dbDetail, "customer_id")

		var genDetail generalFunctionParameter
		genDetail.userTypeID = userTypeID
		genDetail.userCustomerID = userCustomerID
		genDetail.defaultExtLimit = defaultExtLimit

		if userTypeID == "" {
			errorBox(w, "email_error", "header-service-product", "button-service-product")
		} else {
			if userTypeID == "100" {
				header(w, "YAP Admin Account<br>Service/Product Information", "header-service-product", extraButtonName, extraButtonURL)
				serviceProductList(w, dbDetail, genDetail)
				fmt.Fprintf(w, "<br>")
				serviceProductAdd(w, r, dbDetail, genDetail)
				fmt.Fprintf(w, "<br>")
				serviceProductEdit(w, r, dbDetail, genDetail)
				fmt.Fprintf(w, "<br>")
				serviceProductDelete(w, r, dbDetail, genDetail)
				footer(w, "header-service-product", "button-service-product")
			} else {
				errorBox(w, "account_type_error", "header-service-product", "button-service-product")
			}
		}
	})

	yapPort := os.Getenv("yapPort")
	yapPortInt, err := strconv.Atoi(yapPort)

	if err != nil {
		panic("YAP PORT MUST BE A NUMBER IN /etc/yap/yap.env")
	}

	yapAddress := os.Getenv("yapAddress")
	validateYapAddress := validator.New()
	validateYapAddressErr := validateYapAddress.Var(yapAddress, "required,ip_addr")

	if yapPortInt <= 1023 || yapPortInt >= 49152 {
		panic("YAP LISTENING PORT MUST BE IN THE NUMBER RANGE 1024-49151 IN /etc/yap/yap.env")
	} else if validateYapAddressErr != nil && yapAddress != "localhost" {
		panic("YAP ADDRESS MUST BE A VALID INTERENT PROTOCOL (IP) ADDRESS OR localhost IN /etc/yap/yap.env")
	} else if dbAddress == yapAddress && dbPort == yapPort {
		panic("YAP ADDRESS & PORT NUMBER CANNOT BE THE SAME AS DATABASE ADDRESS & PORT NUMBER IN /etc/yap/yap.env")
	} else {
		socket := yapAddress + ":" + yapPort
		fmt.Println("YAP is running on: " + socket)
		// Start YAP server on port specified above
		log.Fatal(http.ListenAndServe(socket, nil))
	}
}

// Contributor(s):
// Elliot Michael Keavney
