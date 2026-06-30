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
	"database/sql"
	"fmt"
	"github.com/ellwould/csvcell"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/sony/sonyflake"
	"log"
	"net/http"
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
	fmt.Fprintf(w, "    <a href=\"https://ell.today\" target=\"_blank\" class=\"button-general button-footer "+buttonCSS+"\">Other Software</a>")
	fmt.Fprintf(w, "  </h1>")
	fmt.Fprintf(w, "</div>")
}

// Function to format ISO DATE format
func formatDate(date string) string {
	const (
		iso    = "2006-01-02"
		layout = "2/January/2006"
	)

	parse, _ := time.Parse(iso, date)
	return string(parse.Format(layout))
}

// Function to format ISO DATETIME format
func formatDateTime(dateTime string) string {
	const (
		iso    = "2006-01-02 15:04:05"
		layout = "2/January/2006 15:04:05 PM"
	)

	parse, _ := time.Parse(iso, dateTime)
	return string(parse.Format(layout))
}

//----------------------------------------------------------------------------------------------------

// Pure embedded HTML Go functions

func inputHTML(w http.ResponseWriter, inputValue string, labelMessage string, inputType string) {
	fmt.Fprintf(w, "  <label for=\""+inputValue+"\"><b>Enter "+labelMessage+"</b>")
	fmt.Fprintf(w, "  </label><br>")
	fmt.Fprintf(w, "  <input type=\""+inputType+"\" id=\""+inputValue+"\" name=\""+inputValue+"\">")
	fmt.Fprintf(w, "<br>")
}

func selectSingleHTML(w http.ResponseWriter, selectValue string, labelMessage string, optionValue []string) {
	fmt.Fprintf(w, "  <label for=\""+selectValue+"\"><b>Select "+labelMessage+"</b>")
	fmt.Fprintf(w, "  </label><br>")
	fmt.Fprintf(w, "  <select id=\""+selectValue+"\" name=\""+selectValue+"\">")
	fmt.Fprintf(w, "<option value></option>")
	for _, value := range optionValue {
		fmt.Fprintf(w, "<option value=\""+string(value)+"\">"+string(value)+"</option>")
	}
	fmt.Fprintf(w, "  </select>")
}

func selectDoubleHTML(w http.ResponseWriter, selectValue string, labelMessage string, optionValue [][]string) {
	fmt.Fprintf(w, "  <label for=\""+selectValue+"\"><b>Select "+labelMessage+"</b>")
	fmt.Fprintf(w, "  </label><br>")
	fmt.Fprintf(w, "  <select id=\""+selectValue+"\" name=\""+selectValue+"\">")
	fmt.Fprintf(w, "<option value></option>")
	for _, value := range optionValue {
		fmt.Fprintf(w, "<option value=\""+string(value[0:][0])+"\">&nbsp "+labelMessage+" ID :"+string(value[0:][0])+" | "+labelMessage+" Name: "+string(value[1:][0])+"</option>")
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

func totalTableCount(w http.ResponseWriter, dbTotalTableCount databaseFunctionParameter) string {
	if dbTotalTableCount.countMinusOne == true {
		var countMinusOne string
		countMinusOneQuery, err := dbTotalTableCount.connection.Query(`SELECT
									       COUNT(*) -1
									       FROM
									         ` + dbTotalTableCount.database + `.` + dbTotalTableCount.table)
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
								         ` + dbTotalTableCount.database + `.` + dbTotalTableCount.table)
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

func totalTableCountWhere(w http.ResponseWriter, dbTotalTableCountWhere databaseFunctionParameter) string {
	if dbTotalTableCountWhere.countMinusOne == true {
		var countMinusOne string
		countMinusOneQuery, err := dbTotalTableCountWhere.connection.Query(`SELECT
                                                                    COUNT(*) -1
                                                                    FROM
                                                                      `+dbTotalTableCountWhere.database+`.`+dbTotalTableCountWhere.table+`
                                                                    WHERE
                                                                      `+dbTotalTableCountWhere.columnWhere+` =?`, dbTotalTableCountWhere.columnWhereValue)
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
								      `+dbTotalTableCountWhere.columnWhere+` =?`, dbTotalTableCountWhere.columnWhereValue)
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

func totalTableCountWhereAnd(w http.ResponseWriter, dbTotalTableCountWhereAnd databaseFunctionParameter) string {
	var count string
	countQuery, err := dbTotalTableCountWhereAnd.connection.Query(`SELECT
								       COUNT(*)
								       FROM
								         `+dbTotalTableCountWhereAnd.database+`.`+dbTotalTableCountWhereAnd.table+`
								       WHERE
								         `+dbTotalTableCountWhereAnd.columnWhere+` =?`+`
								       AND
								         `+dbTotalTableCountWhereAnd.columnWhereAnd+` =?`, dbTotalTableCountWhereAnd.columnWhereValue, dbTotalTableCountWhereAnd.columnWhereValueAnd)
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
		panic("The function userAccountData can only accept the following arguments: type_id, customer_id, customer_name or pbx_id")
	}
	dbSelectWhere.columnWhere = "user_account_email"
	dbSelectWhere.columnWhereValue = dbUserAccountData.columnWhereValue

	return selectWhere(dbSelectWhere)
}

// Function to retrive account type name  and ID or just the account type ID
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
                                                                      yap.user_account_type`)

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

// Function to retrive customer name and ID or just the customer ID
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
                                                               yap.view___customer_detail`)

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
                                                          yap.view___pbx_detail`)

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

//----------------------------------------------------------------------------------------------------

// Function to validate user input utlising the Go Validator package

func validateInput(value string, valueType string) (validation bool) {
	validateInput := validator.New()
	// Conditional statments are used for each type of value inputted from a user
	if valueType == "email" {
		validateInputErr := validateInput.Var(value, "email,required,min=6,max=200,excludes=0x2C")
		if validateInputErr != nil {
			validation = false
			return
		} else {
			validation = true
			return
		}

	} else if valueType == "alpha" {
		validateInputAlphaspaceErr := validateInput.Var(value, "alphaspace,min=1,max=30")
		validateInputSymbolErr := validateInput.Var(value, "excludes=`!\"£$%^&*()_=+{}[];:@'#~\\<>/?")
		if validateInputAlphaspaceErr != nil || validateInputSymbolErr != nil {
			validation = false
			return
		} else {
			validation = true
			return
		}

	} else if valueType == "alphanum" {
		validateInputAlphanumspaceErr := validateInput.Var(value, "alphanumspace,min=1,max=30")
		validateInputSymbolErr := validateInput.Var(value, "excludes=`!\"£$%^&*()_=+{}[];:@'#~\\<>/?")
		if validateInputAlphanumspaceErr != nil || validateInputSymbolErr != nil {
			validation = false
			return
		} else {
			validation = true
			return
		}

	} else {
		validation = false
		return
	}
}

//----------------------------------------------------------------------------------------------------

// Main menu page functions

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
	fmt.Fprintf(w, "    <th>Total SIP Extensions</th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "  <tr>")
	dbTotalTableCount.table = "view___customer_detail"
	dbTotalTableCount.countMinusOne = true
	fmt.Fprintf(w, "    <td>"+totalTableCount(w, dbTotalTableCount)+"</td>")
	dbTotalTableCount.table = "view___pbx_detail"
	dbTotalTableCount.countMinusOne = true
	fmt.Fprintf(w, "    <td>"+totalTableCount(w, dbTotalTableCount)+"</td>")
	dbTotalTableCount.table = "view___sip_extension_detail"
	dbTotalTableCount.countMinusOne = false
	fmt.Fprintf(w, "    <td>"+totalTableCount(w, dbTotalTableCount)+"</td>")
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
	fmt.Fprintf(w, "    <td>"+totalTableCountWhere(w, dbTotalTableCountWhere)+"</td>")
	dbTotalTableCountWhere.columnWhereValue = "200"
	fmt.Fprintf(w, "    <td>"+totalTableCountWhere(w, dbTotalTableCountWhere)+"</td>")
	dbTotalTableCountWhere.columnWhereValue = "201"
	fmt.Fprintf(w, "    <td>"+totalTableCountWhere(w, dbTotalTableCountWhere)+"</td>")
	dbTotalTableCountWhere.columnWhereValue = "300"
	fmt.Fprintf(w, "    <td>"+totalTableCountWhere(w, dbTotalTableCountWhere)+"</td>")
	dbTotalTableCountWhere.columnWhereValue = "301"
	fmt.Fprintf(w, "    <td>"+totalTableCountWhere(w, dbTotalTableCountWhere)+"</td>")
	dbTotalTableCountWhere.columnWhereValue = "302"
	fmt.Fprintf(w, "    <td>"+totalTableCountWhere(w, dbTotalTableCountWhere)+"</td>")
	dbTotalTableCountWhere.columnWhereValue = "400"
	fmt.Fprintf(w, "    <td>"+totalTableCountWhere(w, dbTotalTableCountWhere)+"</td>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "</table>")

}

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
		fmt.Fprintf(w, "    <td>"+totalTableCountWhere(w, dbTotalTableCountWhere)+"</td>")
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
					                pbx_site_contact_number,
					                pbx_invoice_address_line_1,
					                pbx_invoice_address_line_2,
					                pbx_invoice_city_town_village,
					                pbx_invoice_county_state_region,
					                pbx_invoice_postcode_zip_code,
					                pbx_invoice_country,
					                pbx_invoice_contact_email,
					                pbx_invoice_contact_number
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
			pbxName                     string
			pbxID                       string
			pbxSiteAddressLine1         string
			pbxSiteAddressLine2         string
			pbxSiteCityTownVillage      string
			pbxSiteCountyStateRegion    string
			pbxSitePostcodeZipCode      string
			pbxSiteCountry              string
			pbxSiteContactEmail         string
			pbxSiteContactNumber        string
			pbxInvoiceAddressLine1      string
			pbxInvoiceAddressLine2      string
			pbxInvoiceCityTownVillage   string
			pbxInvoiceCountyStateRegion string
			pbxInvoicePostcodeZipCode   string
			pbxInvoiceCountry           string
			pbxInvoiceContactEmail      string
			pbxInvoiceContactNumber     string
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
			&pbxInvoiceAddressLine1,
			&pbxInvoiceAddressLine2,
			&pbxInvoiceCityTownVillage,
			&pbxInvoiceCountyStateRegion,
			&pbxInvoicePostcodeZipCode,
			&pbxInvoiceCountry,
			&pbxInvoiceContactEmail,
			&pbxInvoiceContactNumber,
		)

		// Error
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, "<table id=\"table\" class=\"table-main-menu\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>PBX Name and ID</th>")
		fmt.Fprintf(w, "    <th>Total SIP Extensions in PBX</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td>PBX Name: "+pbxName+"<br><br>PBX ID: "+pbxID+"</td>")
		var dbTotalTableCountWhere databaseFunctionParameter
		dbTotalTableCountWhere.connection = dbPBXAccount.connection
		dbTotalTableCountWhere.database = dbPBXAccount.database
		dbTotalTableCountWhere.table = "view___sip_extension_detail"
		dbTotalTableCountWhere.columnWhere = "pbx_id"
		dbTotalTableCountWhere.columnWhereValue = pbxID
		fmt.Fprintf(w, "    <td>"+totalTableCountWhere(w, dbTotalTableCountWhere)+"</td>")
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
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<table id=\"table\" class=\"table-main-menu\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>PBX Invoice Address</th>")
		fmt.Fprintf(w, "    <th>PBX Invoice Email</th>")
		fmt.Fprintf(w, "    <th>PBX Invoice Phone Number</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td style=\"text-align: left;\">"+pbxInvoiceAddressLine1+"&nbsp<br>"+pbxInvoiceAddressLine2+"<br>"+pbxInvoiceCityTownVillage+"<br>"+pbxInvoiceCountyStateRegion+"<br><br>"+pbxInvoicePostcodeZipCode+"<br><br>"+pbxInvoiceCountry+"</td>")
		fmt.Fprintf(w, "    <td><a href=\"mailto:"+pbxInvoiceContactEmail+"\">"+pbxInvoiceContactEmail+"</a></td>")
		fmt.Fprintf(w, "    <td><a href=\"tel:"+pbxInvoiceContactNumber+"\">"+pbxInvoiceContactNumber+"</a></td>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")

	}

}

// Function for main menu page user information
func mainMenuUserInformation(w http.ResponseWriter, dbUserInformation databaseFunctionParameter, userTypeID string) {

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
		//Account detail tables
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

	if userTypeID == "100" {
		fmt.Fprintf(w, "<br>")
		mainMenuYapAccount(w, dbDetail)
	} else if userTypeID == "200" || userTypeID == "201" {
		fmt.Fprintf(w, "<br>")
		mainMenuCustomerAccount(w, dbDetail)
	} else if userTypeID == "300" || userTypeID == "301" || userTypeID == "302" {
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

// User account page functions

func userAccountList(w http.ResponseWriter, dbDetail databaseFunctionParameter, userTypeID string) {

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

		if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" || userTypeID == "300" {
			fmt.Fprintf(w, "<table id=\"table\" class=\"table-user-account\">")
			fmt.Fprintf(w, "  <tr>")
			if userTypeID == "100" {
				fmt.Fprintf(w, "    <th>Total YAP<br>Admin<br>Accounts<br>(Type ID: 100)</th>")
			}
			if userTypeID == "100" || userTypeID == "200" {
				fmt.Fprintf(w, "    <th>Total Customer<br>Admin<br>Accounts<br>(Type ID: 200)</th>")
				fmt.Fprintf(w, "    <th>Total Customer<br>Regular<br>Accounts<br>(Type ID: 201)</th>")
			}
			if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" || userTypeID == "300" {
				fmt.Fprintf(w, "    <th>Total PBX<br>Admin<br>Accounts<br>(Type ID: 300)</th>")
				fmt.Fprintf(w, "    <th>Total PBX<br>Regular<br>Accounts<br>(Type ID: 301)</th>")
				fmt.Fprintf(w, "    <th>Total PBX<br>Read Only<br>Accounts<br>(Type ID: 302)</th>")
			}
			if userTypeID == "100" || userTypeID == "200" {
				fmt.Fprintf(w, "    <th>Total Customer<br>Invoice<br>Accounts<br>(Type ID: 400)</th>")
			}
			fmt.Fprintf(w, "  </tr>")
			fmt.Fprintf(w, "  <tr>")
			if userTypeID == "100" {
				dbTotalTableCountWhere.columnWhereValue = "100"
				fmt.Fprintf(w, "    <td>"+totalTableCountWhere(w, dbTotalTableCountWhere)+"</td>")
				dbTotalTableCountWhere.columnWhereValue = "200"
				fmt.Fprintf(w, "    <td>"+totalTableCountWhere(w, dbTotalTableCountWhere)+"</td>")
				dbTotalTableCountWhere.columnWhereValue = "201"
				fmt.Fprintf(w, "    <td>"+totalTableCountWhere(w, dbTotalTableCountWhere)+"</td>")
				dbTotalTableCountWhere.columnWhereValue = "300"
				fmt.Fprintf(w, "    <td>"+totalTableCountWhere(w, dbTotalTableCountWhere)+"</td>")
				dbTotalTableCountWhere.columnWhereValue = "301"
				fmt.Fprintf(w, "    <td>"+totalTableCountWhere(w, dbTotalTableCountWhere)+"</td>")
				dbTotalTableCountWhere.columnWhereValue = "302"
				fmt.Fprintf(w, "    <td>"+totalTableCountWhere(w, dbTotalTableCountWhere)+"</td>")
				dbTotalTableCountWhere.columnWhereValue = "400"
				fmt.Fprintf(w, "    <td>"+totalTableCountWhere(w, dbTotalTableCountWhere)+"</td>")
			} else if userTypeID == "200" || userTypeID == "201" {
				dbTotalTableCountWhere.columnWhereAnd = "customer_id"
				dbTotalTableCountWhere.columnWhereValueAnd = customerID
				if userTypeID == "200" {
					dbTotalTableCountWhere.columnWhereValue = "200"
					fmt.Fprintf(w, "    <td>"+totalTableCountWhereAnd(w, dbTotalTableCountWhere)+"</td>")
					dbTotalTableCountWhere.columnWhereValue = "201"
					fmt.Fprintf(w, "    <td>"+totalTableCountWhereAnd(w, dbTotalTableCountWhere)+"</td>")
				}
				dbTotalTableCountWhere.columnWhereValue = "300"
				fmt.Fprintf(w, "    <td>"+totalTableCountWhereAnd(w, dbTotalTableCountWhere)+"</td>")
				dbTotalTableCountWhere.columnWhereValue = "301"
				fmt.Fprintf(w, "    <td>"+totalTableCountWhereAnd(w, dbTotalTableCountWhere)+"</td>")
				dbTotalTableCountWhere.columnWhereValue = "302"
				fmt.Fprintf(w, "    <td>"+totalTableCountWhereAnd(w, dbTotalTableCountWhere)+"</td>")
				if userTypeID == "200" {
					dbTotalTableCountWhere.columnWhereValue = "400"
					fmt.Fprintf(w, "    <td>"+totalTableCountWhereAnd(w, dbTotalTableCountWhere)+"</td>")
				}
			} else if userTypeID == "300" {
				dbTotalTableCountWhere.columnWhereAnd = "pbx_id"
				dbTotalTableCountWhere.columnWhereValueAnd = pbxID
				dbTotalTableCountWhere.columnWhereValue = "300"
				fmt.Fprintf(w, "    <td>"+totalTableCountWhereAnd(w, dbTotalTableCountWhere)+"</td>")
				dbTotalTableCountWhere.columnWhereValue = "301"
				fmt.Fprintf(w, "    <td>"+totalTableCountWhereAnd(w, dbTotalTableCountWhere)+"</td>")
				dbTotalTableCountWhere.columnWhereValue = "302"
				fmt.Fprintf(w, "    <td>"+totalTableCountWhereAnd(w, dbTotalTableCountWhere)+"</td>")
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
		if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" || userTypeID == "300" {
			fmt.Fprintf(w, "  <tr>")
			fmt.Fprintf(w, "    <th><button onclick=\"toggleOtherAccount() \"class=\"button-general button-user-account\">&nbsp Show/Hide Other Accounts &nbsp</button></th>")
			fmt.Fprintf(w, "  </tr>")
		}
		fmt.Fprintf(w, "</table>")
	}

	if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" || userTypeID == "300" {

		userCustomerID := userAccountData(dbDetail, "customer_id")
		userCustomerName := userAccountData(dbDetail, "customer_name")
		userPBXID := userAccountData(dbDetail, "pbx_id")
		userPBXName := userAccountData(dbDetail, "pbx_name")

		fmt.Fprintf(w, "<div id=\"other-account-div\" style=\"display:none\">")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<table id=\"table\" class=\"table-user-account\">")
		fmt.Fprintf(w, "  <tr>")
		if userTypeID == "100" {
			fmt.Fprintf(w, "    <th class=\"table-title\";>All User Account Details on the Server:</th>")
		} else if userTypeID == "200" {
			fmt.Fprintf(w, "    <th class=\"table-title\";>User Account Details for the Customer<br>"+userCustomerName+"<br>(Customer ID: "+userCustomerID+")</th>")
		} else if userTypeID == "201" {
			fmt.Fprintf(w, "    <th class=\"table-title\";>PBX User Account Details for the Customer<br>"+userCustomerName+"<br>(Customer ID: "+userCustomerID+")</th>")
		} else if userTypeID == "300" {
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
		if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
			inputTableHTMLArgument.inputID = "other-account-input-pbx-name"
			inputTableHTMLArgument.funcNameJS = "otherAccountSearchPBXName"
			inputTableHTMLArgument.placeholder = "PBX Name"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "other-account-input-pbx-id"
			inputTableHTMLArgument.funcNameJS = "otherAccountSearchPBXID"
			inputTableHTMLArgument.placeholder = "PBX ID"
			inputTableHTML(w, inputTableHTMLArgument)
		}
		if userTypeID == "100" {
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "other-account-input-customer-name"
			inputTableHTMLArgument.funcNameJS = "otherAccountSearchCustomerName"
			inputTableHTMLArgument.placeholder = "Customer Name"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			fmt.Fprintf(w, "    <br>")
			fmt.Fprintf(w, "    <br>")
			inputTableHTMLArgument.inputID = "other-account-input-customer-id"
			inputTableHTMLArgument.funcNameJS = "otherAccountSearchCustomerID"
			inputTableHTMLArgument.placeholder = "Customer ID"
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
		if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
			fmt.Fprintf(w, "          <th>PBX Name</th>")
			fmt.Fprintf(w, "          <th>PBX ID</th>")
		}
		if userTypeID == "100" {
			fmt.Fprintf(w, "          <th>Customer Name</th>")
			fmt.Fprintf(w, "          <th>Customer ID</th>")
		}

		fmt.Fprintf(w, "        </tr>")

		var whereClause string

		if userTypeID == "100" {
			whereClause = "WHERE customer_id != ? AND pbx_id != ?;"
			userCustomerID = "0"
			userPBXID = "0"
		} else if userTypeID == "200" || userTypeID == "201" {
			whereClause = "WHERE customer_id = ? AND pbx_id != ?;"
			userPBXID = "0"
		} else if userTypeID == "300" {
			whereClause = "WHERE customer_id = ? AND pbx_id = ?;"
		}

		otherUserAccountSQL, err := dbDetail.connection.Query(`SELECT
									 user_account_id,
						     			 user_account_first_name,
						     			 user_account_last_name,  
						     			 user_account_email,                                                   
						     			 user_account_type,  
						     			 user_account_type_id,
						     			 user_account_date_time_added, 
						     			 customer_id,
						     			 customer_name,
						     			 pbx_id,
						     			 pbx_name						     
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
				&customerID,
				&customerName,
				&pbxID,
				&pbxName,
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

			if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
				if pbxName != "system" {
					userAccountListTdHTML(w, userAccountID, userAccountTypeID, pbxName)
				} else {
					userAccountListTdHTML(w, userAccountID, userAccountTypeID, "-")
				}
				if pbxID != "1" {
					userAccountListTdHTML(w, userAccountID, userAccountTypeID, pbxID)
				} else {
					userAccountListTdHTML(w, userAccountID, userAccountTypeID, "-")
				}
			}
			if userTypeID == "100" {
				if customerName != "system" {
					userAccountListTdHTML(w, userAccountID, userAccountTypeID, customerName)
				} else {
					userAccountListTdHTML(w, userAccountID, userAccountTypeID, "-")
				}
				if customerID != "1" {
					userAccountListTdHTML(w, userAccountID, userAccountTypeID, customerID)
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
		if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
			// JS filter function for PBX name in the other account table
			filterTableJSArgument.funcNameJS = "otherAccountSearchPBXName"
			filterTableJSArgument.inputID = "other-account-input-pbx-name"
			filterTableJSArgument.columnNumber = 5
			filterTableJS(w, filterTableJSArgument)
			// JS filter function for PBX ID in the other account table
			filterTableJSArgument.funcNameJS = "otherAccountSearchPBXID"
			filterTableJSArgument.inputID = "other-account-input-pbx-id"
			filterTableJSArgument.columnNumber = 6
			filterTableJS(w, filterTableJSArgument)
		}
		if userTypeID == "100" {
			// JS filter function for the customer name in the other account table
			filterTableJSArgument.funcNameJS = "otherAccountSearchCustomerName"
			filterTableJSArgument.inputID = "other-account-input-customer-name"
			filterTableJSArgument.columnNumber = 7
			filterTableJS(w, filterTableJSArgument)
			// JS filter function for the customer ID in the other account table
			filterTableJSArgument.funcNameJS = "otherAccountSearchCustomerID"
			filterTableJSArgument.inputID = "other-account-input-customer-id"
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
}

// Function to add new user account
func userAccountAdd(w http.ResponseWriter, dbDetail databaseFunctionParameter, userID string, r *http.Request) {

	fmt.Fprintf(w, "<form method=\"POST\" action=\"/user-account\">")
	fmt.Fprintf(w, "<table class=\"table-user-account\">")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th class=\"table-title\";>Add New User Account</th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>")
	fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
	fmt.Fprintf(w, "        <tr>")
	fmt.Fprintf(w, "          <td>")
	inputHTML(w, "add_account_input_first_name", "First Name:", "text")
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "          <td>")
	inputHTML(w, "add_account_input_last_name", "Last Name:", "text")
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "          <td>")
	inputHTML(w, "add_account_input_email", "Email Address:", "text")
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
	validateFirstNameAlpha := validateInput(addAccountInputFirstName, "alpha")

	// Validate the last name string
	validateLastNameAlpha := validateInput(addAccountInputLastName, "alpha")

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

	if addAccountInputFirstName == "" {
	} else if validateFirstNameAlpha == false || validateLastNameAlpha == false {
		messageHTML(w, "Name length must be 1 to 30 characters and must only contain characters: a-z A-Z or -", "warning")
	} else if validateEmail == false {
		messageHTML(w, "A valid email address must be entered", "warning")
	} else if validateUserAccountTypeID == false {
		messageHTML(w, "User account type must be either 100, 200, 201, 300, 301, 302 or 400", "warning")
	} else if validatePBXID == false {
		messageHTML(w, "PBX ID invalid", "warning")
	} else if validateCustomerID == false {
		messageHTML(w, "Customer ID invalid", "warning")
	} else if addAccountSelectAccountType == "100" && userID != "1" {
		messageHTML(w, "Must be a YAP Admin (100) account with account ID 1 to create other YAP Admin (100) accounts", "warning")
	} else if addAccountSelectAccountType == "200" && addAccountSelectCustomerID == "" {
		messageHTML(w, "A customer ID must be selected when creating a 200 type account", "warning")
	} else if addAccountSelectAccountType == "201" && addAccountSelectCustomerID == "" {
		messageHTML(w, "A customer ID must be selected when creating a 201 type account", "warning")
	} else if addAccountSelectAccountType == "400" && addAccountSelectCustomerID == "" {
		messageHTML(w, "A customer ID must be selected when creating a 400 type account", "warning")
	} else if addAccountSelectAccountType == "300" && addAccountSelectPBXID == "" {
		messageHTML(w, "A PBX ID must be selected when creating a 300 type account", "warning")
	} else if addAccountSelectAccountType == "301" && addAccountSelectPBXID == "" {
		messageHTML(w, "A PBX ID must be selected when creating a 301 type account", "warning")
	} else if addAccountSelectAccountType == "302" && addAccountSelectPBXID == "" {
		messageHTML(w, "A PBX ID must be selected when creating a 302 type account", "warning")
	} else {
		if addAccountSelectAccountType == "100" {
			addAccountSelectCustomerID = "1"
			addAccountSelectPBXID = "1"
			messageHTML(w, "YAP Admin (100) Account: "+addAccountInputEmail+" Created", "success")
		} else if addAccountSelectAccountType == "200" || addAccountSelectAccountType == "201" || addAccountSelectAccountType == "400" {
			addAccountSelectPBXID = "1"
			if addAccountSelectAccountType == "200" {
				messageHTML(w, "Customer Admin (200) Account: "+addAccountInputEmail+" Created", "success")
			} else if addAccountSelectAccountType == "201" {
				messageHTML(w, "Customer Regular (201) Account: "+addAccountInputEmail+" Created", "success")
			} else if addAccountSelectAccountType == "400" {
				messageHTML(w, "Customer Invoice (400) Account: "+addAccountInputEmail+" Created", "success")
			}
		} else if addAccountSelectAccountType == "300" || addAccountSelectAccountType == "301" || addAccountSelectAccountType == "302" {
			var customerID string
			customerIDSQL, err := dbDetail.connection.Query(`SELECT
                                                                    customer_id
                                                                FROM
                                                                    yap.view___pbx_detail
                                                                WHERE pbx_id = ?`, addAccountSelectPBXID)

			// Error
			if err != nil {
				panic(err)
			}

			for customerIDSQL.Next() {

				err = customerIDSQL.Scan(
					&customerID,
				)

				// Error
				if err != nil {
					panic(err)
				}
				addAccountSelectCustomerID = customerID

			}

			messageHTML(w, "Account "+addAccountInputEmail+" Created", "success")
		}

		dbDetail.connection.Query(`INSERT 
        	                   INTO
	       		       user_account (
			           email,
			           first_name,
			           last_name,
			           user_account_type_id,
			           customer_id,
			           pbx_id,
			           account_active)
			       VALUES(?, ?, ?, ?, ?, ?, ?);`,
			addAccountInputEmail,
			addAccountInputFirstName,
			addAccountInputLastName,
			addAccountSelectAccountType,
			addAccountSelectCustomerID,
			addAccountSelectPBXID,
			"1")
	}
}

func userAccountDelete(w http.ResponseWriter, dbDetail databaseFunctionParameter, userID string, r *http.Request) {

	// Get account type ID and email from the database and append to slice
	var userAccountIDList []string
	var userAccountID string

	var userAccountEmailList []string
	var userAccountEmail string

	userAccountIDEmailSQL, err := dbDetail.connection.Query(`SELECT
	                                                             user_account_id,
	                                                             user_account_email
	                                                         FROM
	                                                             yap.view___account_detail`)

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

	fmt.Fprintf(w, "<form method=\"POST\" action=\"/user-account\">")
	fmt.Fprintf(w, "<table class=\"table-delete\">")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th class=\"table-title\";>Delete User Account</th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>")
	fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
	fmt.Fprintf(w, "        <tr>")
	fmt.Fprintf(w, "          <td>")
	inputHTML(w, "delete_account_input_account_id", "Account ID:", "text")
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "          <td>")
	inputHTML(w, "delete_account_input_account_email", "Account Email:", "text")
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "          <td>")
	userAccountTypeIDNameList, _ := userAccountTypeSlice(dbDetail)
	selectDoubleHTML(w, "delete_account_select_account_type", "Account Type", userAccountTypeIDNameList)
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

	deleteUserAccountInputAccountID := r.FormValue("delete_account_input_account_id")
	deleteUserAccountInputEmail := r.FormValue("delete_account_input_account_email")
	deleteUserAccountSelectAccountType := r.FormValue("delete_account_select_account_type")

	// Check user account ID is contained in the slice
	validateUserAccountID := slices.Contains(userAccountIDList, deleteUserAccountInputAccountID)

	// Check user email is contained in the slice
	validateUserAccountEmail := slices.Contains(userAccountEmailList, deleteUserAccountInputEmail)

	// Check user type ID is contained in the slice
	_, userAccountTypeIDList := userAccountTypeSlice(dbDetail)
	validateUserAccountTypeID := slices.Contains(userAccountTypeIDList, deleteUserAccountSelectAccountType)

	if deleteUserAccountInputAccountID == "" {
		// Do nothing
	} else if validateUserAccountID == false {
		messageHTML(w, "Account ID does not exist", "warning")
	} else if validateUserAccountEmail == false {
		messageHTML(w, "Account email address does not exist", "warning")
	} else if validateUserAccountTypeID == false {
		messageHTML(w, "User account type must be either 100, 200, 201, 300, 301, 302 or 400", "warning")
	} else if deleteUserAccountInputAccountID == "1" {
		messageHTML(w, "YAP Admin account with ID 1 cannot be deleted", "warning")
	} else if deleteUserAccountSelectAccountType == "100" && userID != "1" {
		messageHTML(w, "Must be a YAP Admin (100) account with account ID 1 to delete other YAP Admin (100) accounts", "warning")
	} else {
		dbDetail.connection.Query(`DELETE FROM user_account WHERE id = ? AND user_account_type_id = ?;`, deleteUserAccountInputAccountID, deleteUserAccountSelectAccountType)

		// Close connection
		defer dbDetail.connection.Close()

		var checkUserAccountDeleted string

		checkUserAccountDeletedSQL, err := dbDetail.connection.Query(`SELECT
									        user_account_id
									      FROM view___account_detail
									      WHERE user_account_id = ?;`,
			deleteUserAccountInputAccountID)

		// Error
		if err != nil {
			panic(err)
		}

		for checkUserAccountDeletedSQL.Next() {

			err = checkUserAccountDeletedSQL.Scan(
				&checkUserAccountDeleted,
			)

			// Error
			if err != nil {
				panic(err)
			}

		}

		// Close connection
		defer dbDetail.connection.Close()

		if checkUserAccountDeleted == "" {
			messageHTML(w, "Account "+deleteUserAccountInputEmail+" deleted", "success")
		} else {
			messageHTML(w, "Wrong account type ("+deleteUserAccountSelectAccountType+") selected for "+deleteUserAccountInputEmail, "warning")
		}
	}
}

//----------------------------------------------------------------------------------------------------

// Customer page functions

func customerList(w http.ResponseWriter, dbDetail databaseFunctionParameter, userTypeID string, userCustomerID string) {

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

	if userTypeID == "100" {
		fmt.Fprintf(w, "<table id=\"table\" class=\"table-customer\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table id=\"table\" class=\"table-customer\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <th>Total Customers On YAP</th>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		dbTableCountUserCustomer.countMinusOne = true
		fmt.Fprintf(w, "          <td>"+totalTableCount(w, dbTableCountUserCustomer)+"</td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><button onclick=\"toggleCustomer() \"class=\"button-general button-customer\">&nbsp Show/Hide Customers &nbsp</button></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
	}

	if userTypeID == "100" {
		fmt.Fprintf(w, "<div id=\"customer-div\" style=\"display:none\">")
		fmt.Fprintf(w, "<br>")
	} else {
		fmt.Fprintf(w, "<div id=\"customer-div\">")
	}
	fmt.Fprintf(w, "<table id=\"table\" class=\"table-customer\">")
	fmt.Fprintf(w, "  <tr>")
	if userTypeID == "100" {
		fmt.Fprintf(w, "    <th class=\"table-title\";>All Customer Contact Details on the Server:</th>")
	} else {
		fmt.Fprintf(w, "    <th class=\"table-title\";>Customer Contact Details</th>")
	}
	fmt.Fprintf(w, "  </tr>")
	if userTypeID == "100" {
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
	fmt.Fprintf(w, "          <th>Customer<br>ID</th>")
	fmt.Fprintf(w, "          <th>Customer<br>Name</th>")
	fmt.Fprintf(w, "          <th>Customer Site<br>Address</th>")
	fmt.Fprintf(w, "          <th>Customer Site<br>Email Address</th>")
	fmt.Fprintf(w, "          <th>Customer Site<br>Phone Number</th>")
	fmt.Fprintf(w, "          <th>Customer Invoice<br>Address</th>")
	fmt.Fprintf(w, "          <th>Customer Invoice<br>Email Address</th>")
	fmt.Fprintf(w, "          <th>Customer Invoice<br>Phone Number</th>")
	fmt.Fprintf(w, "        </tr>")

	var whereClause string

	if userTypeID == "100" {
		whereClause = "WHERE customer_id != ?;"
		userCustomerID = "1"
	} else if userTypeID == "200" || userTypeID == "201" {
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
						      `+whereClause, userCustomerID)

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
		fmt.Fprintf(w, "          <td>&nbsp<a href=\"mailto:"+customerInvoiceContactEmail+"\">"+customerInvoiceContactEmail+"</a></td>")
		fmt.Fprintf(w, "          <td><a href=\"tel:"+customerInvoiceContactNumber+"\">"+customerInvoiceContactNumber+"</a></td>")
		fmt.Fprintf(w, "        </tr>")
	}

	fmt.Fprintf(w, "      </table>")
	if userTypeID == "100" {
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

	// Customer resource table
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<table id=\"table\" class=\"table-customer\">")
	fmt.Fprintf(w, "  <tr>")
	if userTypeID == "100" {
		fmt.Fprintf(w, "    <th class=\"table-title\";>All PBX  Default Resource Limits on the Server for Each Customer:</th>")
	} else {
		fmt.Fprintf(w, "    <th class=\"table-title\";>PBX Default Resource Limits for Customer</th>")
	}
	fmt.Fprintf(w, "  </tr>")
	if userTypeID == "100" {
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "    <br>")
		var inputTableHTMLArgument jsFunctionParameter
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "customer-resource-input-customer-id"
		inputTableHTMLArgument.funcNameJS = "customerResourceSearchCustomerID"
		inputTableHTMLArgument.placeholder = "Customer ID"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "customer-resource-input-customer-name"
		inputTableHTMLArgument.funcNameJS = "customerResourceSearchCustomerName"
		inputTableHTMLArgument.placeholder = "Customer Name"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "customer-resource-input-date-time"
		inputTableHTMLArgument.funcNameJS = "customerResourceSearchDateTime"
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
	exportCSVButtonHTMLArgument.funcNameJS = "CustomerResource"
	exportCSVButtonHTMLArgument.buttonCSS = "button-customer"
	exportCSVButtonHTML(w, exportCSVButtonHTMLArgument)
	fmt.Fprintf(w, "    </th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>")
	fmt.Fprintf(w, "      <table id=\"customer-resource-table\" class=\"table-customer\">")
	fmt.Fprintf(w, "        <tr>")
	fmt.Fprintf(w, "          <th>Customer<br>ID</th>")
	fmt.Fprintf(w, "          <th>Customer<br>Name</th>")
	fmt.Fprintf(w, "          <th>Date & Time Customer<br>Added</th>")
	fmt.Fprintf(w, "          <th>UK Based</th>")
	fmt.Fprintf(w, "          <th>Consumer<br>Type</th>")
	fmt.Fprintf(w, "          <th>UK VAT Registered</th>")
	fmt.Fprintf(w, "          <th>UK VAT Number</th>")
	fmt.Fprintf(w, "          <th>Reselling<br>Minutes</th>")
	fmt.Fprintf(w, "          <th>PBX Limit</th>")
	fmt.Fprintf(w, "        </tr>")

	customerResourceSQL, err := dbDetail.connection.Query(`SELECT
							customer_id,
							customer_name,
							customer_date_time_added,
							customer_uk_based,
							customer_consumer_type,
							customer_uk_vat_registered,
							customer_uk_vat_number,
							customer_reselling_minutes,
							customer_pbx_limit
					              FROM
					  	        yap.view___customer_detail
						      `+whereClause, userCustomerID)

	// Error
	if err != nil {
		panic(err)

	}

	for customerResourceSQL.Next() {

		err = customerResourceSQL.Scan(
			&customerID,
			&customerName,
			&customerDateTimeAdded,
			&customerUKBased,
			&customerConsumerType,
			&customerUKVATRegistered,
			&customerUKVATNumber,
			&customerResellingMinutes,
			&customerPBXLimit,
		)

		// Error
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>"+customerID+"</td>")
		fmt.Fprintf(w, "          <td>"+customerName+"</td>")
		fmt.Fprintf(w, "          <td>"+formatDateTime(customerDateTimeAdded)+"</td>")
		fmt.Fprintf(w, "          <td>"+customerUKBased+"</td>")
		fmt.Fprintf(w, "          <td>"+customerConsumerType+"</td>")
		fmt.Fprintf(w, "          <td>"+customerUKVATRegistered+"</td>")
		fmt.Fprintf(w, "          <td>"+customerUKVATNumber+"</td>")
		fmt.Fprintf(w, "          <td>"+customerResellingMinutes+"</td>")
		fmt.Fprintf(w, "          <td>"+customerPBXLimit+"</td>")
		fmt.Fprintf(w, "        </tr>")
	}

	fmt.Fprintf(w, "      </table>")
	if userTypeID == "100" {
		var filterTableJSArgument jsFunctionParameter
		filterTableJSArgument.tableID = "customer-resource-table"
		// Call JS filter function for the customer name in the customer resource table
		filterTableJSArgument.funcNameJS = "customerResourceSearchCustomerID"
		filterTableJSArgument.inputID = "customer-resource-input-customer-id"
		filterTableJSArgument.columnNumber = 0
		filterTableJS(w, filterTableJSArgument)
		// Call JS filter function for the customer ID in the customer resource table
		filterTableJSArgument.funcNameJS = "customerResourceSearchCustomerName"
		filterTableJSArgument.inputID = "customer-resource-input-customer-name"
		filterTableJSArgument.columnNumber = 1
		filterTableJS(w, filterTableJSArgument)
		// Call JS filter function for date and time in the customer resource table
		filterTableJSArgument.funcNameJS = "customerResourceSearchDateTime"
		filterTableJSArgument.inputID = "customer-resource-input-date-time"
		filterTableJSArgument.columnNumber = 2
		filterTableJS(w, filterTableJSArgument)
	}
	exportCSVJSArgument.funcNameJS = "customerResource"
	exportCSVJSArgument.tableID = "customer-resource-table"
	exportCSVJSArgument.fileName = "YAP_customer_resource_details"
	exportCSVJSArgument.pathURL = "customer"
	exportCSVJS(w, exportCSVJSArgument)
	fmt.Fprintf(w, "    </th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "</table>")
	fmt.Fprintf(w, "</div>")
	if userTypeID == "100" {
		var toggleDivJSArgument jsFunctionParameter
		toggleDivJSArgument.funcNameJS = "toggleCustomer"
		toggleDivJSArgument.divID = "customer-div"
		toggleDivJS(w, toggleDivJSArgument)
	}

}

// Add customer function

func customerAdd(w http.ResponseWriter, dbDetail databaseFunctionParameter, r *http.Request) {

	fmt.Fprintf(w, "<form method=\"POST\" action=\"/customer\">")
	fmt.Fprintf(w, "<table class=\"table-customer\">")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th class=\"table-title\";>Add New Customer</th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>")
	fmt.Fprintf(w, "      <table style=\"border-style:hidden\">")
	fmt.Fprintf(w, "        <tr>")
	fmt.Fprintf(w, "          <td>")
	inputHTML(w, "add_customer_input_id", "Customer ID:", "text")
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "          <td>")
	inputHTML(w, "add_customer_input_name", "Customer Name:", "text")
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "          <td>")
	ukBasedList := []string{"yes", "no"}
	selectSingleHTML(w, "add_customer_select_uk_based", "Customer UK Based:", ukBasedList)
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "          <td>")
	resellingMinutesList := []string{"yes", "no"}
	selectSingleHTML(w, "add_customer_select_reselling_minutes", "Reselling Minutes:", resellingMinutesList)
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "        </tr>")
	fmt.Fprintf(w, "        <tr>")
	fmt.Fprintf(w, "          <td>")
	consumerTypeList := []string{"Residentail", "Sole Trader", "Partnership", "Limited Liability Partnership (LLP)", "Private Limited Company (LTD)", "Public Limited Company (PLC)", "Community Interest Company (CIC)"}
	selectSingleHTML(w, "add_customer_select_consumer_type", "Consumer Type:", consumerTypeList)
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "          <td>")
	ukVATRegisteredList := []string{"yes", "no"}
	selectSingleHTML(w, "add_customer_select_uk_vat_registered", "UK VAT Registered:", ukVATRegisteredList)
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "          <td>")
	inputHTML(w, "add_customer_input_uk_vat_number", "UK VAT Number:", "text")
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "          <td>")
	pbxLimitList := []string{"1", "2", "3", "4", "5", "10", "25", "50", "75", "100"}
	selectSingleHTML(w, "add_customer_select_pbx_limit", "PBX Limit:", pbxLimitList)
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "        </tr>")
	fmt.Fprintf(w, "      </tr>")
	fmt.Fprintf(w, "        <td style=\"border: none;\">")
	fmt.Fprintf(w, "          <br>")
	fmt.Fprintf(w, "          <br>")
	fmt.Fprintf(w, "        </td>")
	fmt.Fprintf(w, "      <tr>")
	fmt.Fprintf(w, "        <tr>")
	fmt.Fprintf(w, "          <td>")
	inputHTML(w, "add_customer_input_site_address_line_1", "Site Address Line One:", "text")
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "          <td>")
	inputHTML(w, "add_customer_input_site_address_line_2", "Site Address Line Two:", "text")
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "          <td>")
	inputHTML(w, "add_customer_input_site_city_town_village", "Site City/Town/Village:", "text")
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "          <td>")
	inputHTML(w, "add_customer_input_site_county_state_region", "Site County/State/Region:", "text")
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "        </tr>")
	fmt.Fprintf(w, "        <tr>")
	fmt.Fprintf(w, "          <td>")
	inputHTML(w, "add_customer_input_site_postcode_zip_code", "Site Postcode/Zip Code:", "text")
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "          <td>")
	inputHTML(w, "add_customer_input_site_country", "Site Country:", "text")
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "          <td>")
	inputHTML(w, "add_customer_input_site_contact_email", "Site Contact Email Address:", "text")
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "          <td>")
	inputHTML(w, "add_customer_input_site_contact_number", "Site Contact Phone Number:", "text")
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "        </tr>")
	fmt.Fprintf(w, "      </tr>")
	fmt.Fprintf(w, "        <td style=\"border: none;\">")
	fmt.Fprintf(w, "          <br>")
	fmt.Fprintf(w, "          <br>")
	fmt.Fprintf(w, "        </td>")
	fmt.Fprintf(w, "      <tr>")
	fmt.Fprintf(w, "        <tr>")
	fmt.Fprintf(w, "          <td>")
	inputHTML(w, "add_customer_input_invoice_address_line_1", "Invoice Address Line One:", "text")
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "          <td>")
	inputHTML(w, "add_customer_input_invoice_address_line_2", "Invoice Address Line Two:", "text")
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "          <td>")
	inputHTML(w, "add_customer_input_invoice_city_town_village", "Invoice City/Town/Village:", "text")
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "          <td>")
	inputHTML(w, "add_customer_input_invoice_county_state_region", "Invoice County/State/Region:", "text")
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "        </tr>")
	fmt.Fprintf(w, "        <tr>")
	fmt.Fprintf(w, "          <td>")
	inputHTML(w, "add_customer_input_invoice_postcode_zip_code", "Invoice Postcode/Zip Code:", "text")
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "          <td>")
	inputHTML(w, "add_customer_input_invoice_country", "Invoice Country:", "text")
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "          <td>")
	inputHTML(w, "add_customer_input_invoice_contact_email", "Invoice Contact Email Address:", "text")
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "          <td>")
	inputHTML(w, "add_customer_input_invoice_contact_number", "Invoice Contact Phone Number:", "text")
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
	/*
	           addCustomerInputID := r.FormValue("add_customer_input_id")
	           addCustomerInputName := r.FormValue("add_customer_input_name")
	           addCustomerSelectUKBased := r.FormValue("add_customer_select_uk_based")
	           addCustomerSelectResellingMinutes := r.FormValue("add_customer_select_reselling_minutes")
	           addCustomerSelectConsumerType := r.FormValue("add_customer_select_consumer_type")
	           addCustomerSelectUKVATRegistered := r.FormValue("add_customer_select_uk_vat_registered")
	           addCustomerInputUKVATNumber := r.FormValue("add_customer_input_uk_vat_number")
	           addCustomerSelectPBXLimit := r.FormValue("add_customer_select_pbx_limit")

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

	           // Validate the ID
	           validateIDAlphaNumeric := validateInput(addCustomerInputID, "alphanumeric")
	           // Validate the name
	           validateNameAlphaNumeric := validateInput(addCustomerInputName, "alphanumeric")
	           // Validate UK based
	   	validateUKBased := slices.Contains(ukBasedList, addCustomerSelectUKBased)
	   	// Validate reselling minutes
	   	validateResellingMinutes := slices.Contains(resellingMinutesList, addCustomerSelectResellingMinutes)
	   	// Validate consumer type
	   	validateConsumerType := slices.Contains(consumerTypeList, addCustomerSelectConsumerType)
	   	// Validate UK VAT registered status
	   	validateUKVATRegistered := slices.Contains(ukVATRegisteredList, addCustomerSelectUKVATRegistered)
	   	// Validate UK VAT number
	           validateUKVATNumber := validateInput(addCustomerInputUKVATNumber, "alphanumeric")
	           // Validate PBX limit
	           validatePBXLimit := slices.Contains(pbxLimitList, addCustomerSelectPBXLimit)

	           // Validate site address line one
	           validateSiteAddressLine1 := validateInput(addCustomerInputSiteAddressLine1, "alphanumeric")
	           // Validate site address line two
	           validateSiteAddressLine2 := validateInput(addCustomerInputSiteAddressLine2, "alphanumeric")
	           // Validate site city/town/village
	           validateSiteCityTownVillage := validateInput(addCustomerInputSiteCityTownVillage, "alphanumeric")
	           // Validate site county/state/region
	           validateSiteCountyStateRegion := validateInput(addCustomerInputSiteCountyStateRegion, "alphanumeric")
	           // Validate site postcode/zip code
	           validateSitePostcodeZipCode := validateInput(addCustomerInputSitePostcodeZipCode, "alphanumeric")
	           // Validate site country
	           validateSiteCountry := validateInput(addCustomerInputSiteCountry, "alphanumeric")
	           // Validate site contact emial
	           validateSiteContactEmail := validateInput(addCustomerInputSiteContactEmail, "alphanumeric")
	           // Validate Site contact phone number
	           validateSiteContactNumber := validateInput(addCustomerInputSiteContactNumber, "alphanumeric")

	   	// Validate invoice address line one
	           validateInvoiceAddressLine1 := validateInput(addCustomerInputInvoiceAddressLine1, "alphanumeric")
	           // Validate invoice address line two
	           validateInvoiceAddressLine2 := validateInput(addCustomerInputInvoiceAddressLine2, "alphanumeric")
	           // Validate invoice city/town/village
	           validateInvoiceCityTownVillage := validateInput(addCustomerInputInvoiceCityTownVillage, "alphanumeric")
	           // Validate invoice county/state/region
	           validateInvoiceCountyStateRegion := validateInput(addCustomerInputInvoiceCountyStateRegion, "alphanumeric")
	           // Validate invoice postcode/zip code
	           validateInvoicePostcodeZipCode := validateInput(addCustomerInputInvoicePostcodeZipCode, "alphanumeric")
	           // Validate invoice country
	           validateInvoiceCountry := validateInput(addCustomerInputInvoiceCountry, "alphanumeric")
	           // Validate invoice contact emial
	           validateInvoiceContactEmail := validateInput(addCustomerInputInvoiceContactEmail, "alphanumeric")
	           // Validate invoice contact phone number
	           validateInvoiceContactNumber := validateInput(addCustomerInputInvoiceContactNumber, "alphanumeric")

	           if addCustomerInputID == "" {
	           	// Do Nothing
	           } else if validateIDAlphaNumeric == false {
	   		messageHTML(w, "ID length must be 1 to 30 characters and must only contain characters: a-z A-Z or numbers", "warning")
	   	} else if validateName == false {
	                   messageHTML(w, "A valid email address must be entered", "warning")

	*/
}

//----------------------------------------------------------------------------------------------------

// PBX page functions

func pbxList(w http.ResponseWriter, dbDetail databaseFunctionParameter, userTypeID string, userCustomerID string, userPBXID string) {

	var (
		pbxID                       string
		pbxName                     string
		pbxDateTimeAdded            string
		pbxSIPExtensionLimit        string
		pbxSiteAddressLine1         string
		pbxSiteAddressLine2         string
		pbxSiteCityTownVillage      string
		pbxSiteCountyStateRegion    string
		pbxSitePostcodeZipCode      string
		pbxSiteCountry              string
		pbxSiteContactEmail         string
		pbxSiteContactNumber        string
		pbxInvoiceAddressLine1      string
		pbxInvoiceAddressLine2      string
		pbxInvoiceCityTownVillage   string
		pbxInvoiceCountyStateRegion string
		pbxInvoicePostcodeZipCode   string
		pbxInvoiceCountry           string
		pbxInvoiceContactEmail      string
		pbxInvoiceContactNumber     string
		customerID                  string
		customerName                string
	)

	var dbTableCountUserPBX databaseFunctionParameter
	dbTableCountUserPBX.connection = dbDetail.connection
	dbTableCountUserPBX.database = dbDetail.database
	dbTableCountUserPBX.table = "view___pbx_detail"

	if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
		fmt.Fprintf(w, "<table id=\"table\" class=\"table-pbx\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table id=\"table\" class=\"table-pbx\">")
		fmt.Fprintf(w, "        <tr>")
		if userTypeID == "100" {
			fmt.Fprintf(w, "          <th>Total PBXs On YAP</th>")
		}
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		if userTypeID == "100" {
			dbTableCountUserPBX.countMinusOne = true
			fmt.Fprintf(w, "    <td>"+totalTableCount(w, dbTableCountUserPBX)+"</td>")
		}
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><button onclick=\"togglePBX() \"class=\"button-general button-pbx\">&nbsp Show/Hide PBX(s) &nbsp</button></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
	}

	if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
		fmt.Fprintf(w, "<div id=\"pbx-div\" style=\"display:none\">")
		fmt.Fprintf(w, "<br>")
	} else {
		fmt.Fprintf(w, "<div id=\"pbx-div\">")
	}
	fmt.Fprintf(w, "<table id=\"table\" class=\"table-pbx\">")
	fmt.Fprintf(w, "  <tr>")
	if userTypeID == "100" {
		fmt.Fprintf(w, "    <th class=\"table-title\";>All PBX Contact Details on the Server:</th>")
	} else {
		fmt.Fprintf(w, "    <th class=\"table-title\";>PBX Contact Details</th>")
	}
	fmt.Fprintf(w, "  </tr>")
	if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
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
		inputTableHTMLArgument.inputID = "pbx-contact-input-invoice-address"
		inputTableHTMLArgument.funcNameJS = "pbxContactSearchInvoiceAddress"
		inputTableHTMLArgument.placeholder = "PBX Invoice Address"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "pbx-contact-input-invoice-email"
		inputTableHTMLArgument.funcNameJS = "pbxContactSearchInvoiceEmail"
		inputTableHTMLArgument.placeholder = "PBX Invoice Email Address"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "pbx-contact-input-invoice-phone"
		inputTableHTMLArgument.funcNameJS = "pbxContactSearchInvoicePhone"
		inputTableHTMLArgument.placeholder = "PBX Invoice Phone Number"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    <br>")
		if userTypeID == "100" {
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "pbx-contact-input-customer-id"
			inputTableHTMLArgument.funcNameJS = "pbxContactSearchCustomerID"
			inputTableHTMLArgument.placeholder = "Customer ID"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "pbx-contact-input-customer-name"
			inputTableHTMLArgument.funcNameJS = "pbxContactSearchCustomerName"
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
	if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
		fmt.Fprintf(w, "          <th>PBX ID</th>")
		fmt.Fprintf(w, "          <th>PBX Name</th>")
	}
	fmt.Fprintf(w, "          <th>PBX Site<br> Address</th>")
	fmt.Fprintf(w, "          <th>PBX Site<br> Email Address</th>")
	fmt.Fprintf(w, "          <th>PBX Site<br> Phone Number</th>")
	fmt.Fprintf(w, "          <th>PBX Invoice<br> Address</th>")
	fmt.Fprintf(w, "          <th>PBX Invoice<br> Email Address</th>")
	fmt.Fprintf(w, "          <th>PBX Invoice<br> Phone Number</th>")
	if userTypeID == "100" {
		fmt.Fprintf(w, "          <th>Customer ID</th>")
		fmt.Fprintf(w, "          <th>Customer Name</th>")
	}
	fmt.Fprintf(w, "        </tr>")

	var whereClause string
	var userWhereID string

	if userTypeID == "100" {
		whereClause = "WHERE pbx_id != ?;"
		userWhereID = "1"
	} else if userTypeID == "200" || userTypeID == "201" {
		whereClause = "WHERE customer_id = ?;"
		userWhereID = userCustomerID
	} else if userTypeID == "300" || userTypeID == "301" || userTypeID == "302" {
		whereClause = "WHERE pbx_id = ?;"
		userWhereID = userPBXID
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
					                pbx_invoice_address_line_1,
					                pbx_invoice_address_line_2,
					                pbx_invoice_city_town_village,
					                pbx_invoice_county_state_region,
					                pbx_invoice_postcode_zip_code,
					                pbx_invoice_country,
					                pbx_invoice_contact_email,
					                pbx_invoice_contact_number,
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
			&pbxInvoiceAddressLine1,
			&pbxInvoiceAddressLine2,
			&pbxInvoiceCityTownVillage,
			&pbxInvoiceCountyStateRegion,
			&pbxInvoicePostcodeZipCode,
			&pbxInvoiceCountry,
			&pbxInvoiceContactEmail,
			&pbxInvoiceContactNumber,
			&customerID,
			&customerName,
		)

		// Error
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "        <tr>")
		if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
			fmt.Fprintf(w, "          <td>"+pbxID+"</td>")
			fmt.Fprintf(w, "          <td>"+pbxName+"</td>")
		}
		fmt.Fprintf(w, "          <td style=\"text-align: left;\">"+pbxSiteAddressLine1+"&nbsp<br>"+pbxSiteAddressLine2+"&nbsp<br>"+pbxSiteCityTownVillage+"&nbsp<br>"+pbxSiteCountyStateRegion+"&nbsp<br><br>"+pbxSitePostcodeZipCode+"&nbsp<br><br>"+pbxSiteCountry+"&nbsp</td>")
		fmt.Fprintf(w, "          <td><a href=\"mailto:"+pbxSiteContactEmail+"\">"+pbxSiteContactEmail+"</a></td>")
		fmt.Fprintf(w, "          <td><a href=\"tel:"+pbxSiteContactNumber+"\">"+pbxSiteContactNumber+"</a></td>")
		fmt.Fprintf(w, "          <td style=\"text-align: left;\">"+pbxInvoiceAddressLine1+"&nbsp<br>"+pbxInvoiceAddressLine2+"&nbsp<br>"+pbxInvoiceCityTownVillage+"&nbsp<br>"+pbxInvoiceCountyStateRegion+"&nbsp<br><br>"+pbxInvoicePostcodeZipCode+"&nbsp<br><br>"+pbxInvoiceCountry+"&nbsp</td>")
		fmt.Fprintf(w, "          <td>&nbsp<a href=\"mailto:"+pbxInvoiceContactEmail+"\">"+pbxInvoiceContactEmail+"</a></td>")
		fmt.Fprintf(w, "          <td><a href=\"tel:"+pbxInvoiceContactNumber+"\">"+pbxInvoiceContactNumber+"</a></td>")
		if userTypeID == "100" {
			fmt.Fprintf(w, "          <td>"+customerID+"</td>")
			fmt.Fprintf(w, "          <td>"+customerName+"</td>")
		}
		fmt.Fprintf(w, "        </tr>")
	}

	fmt.Fprintf(w, "      </table>")
	if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
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
		// Call JS filter function for invoice address in the PBX contact table
		filterTableJSArgument.funcNameJS = "pbxContactSearchInvoiceAddress"
		filterTableJSArgument.inputID = "pbx-contact-input-invoice-address"
		filterTableJSArgument.columnNumber = 5
		filterTableJS(w, filterTableJSArgument)
		// Call JS filter function for invoice email in the PBX contact table
		filterTableJSArgument.funcNameJS = "pbxContactSearchInvoiceEmail"
		filterTableJSArgument.inputID = "pbx-contact-input-invoice-email"
		filterTableJSArgument.columnNumber = 6
		filterTableJS(w, filterTableJSArgument)
		// Call JS filter function for invoice phone in the PBX contact table
		filterTableJSArgument.funcNameJS = "pbxContactSearchInvoicePhone"
		filterTableJSArgument.inputID = "pbx-contact-input-invoice-phone"
		filterTableJSArgument.columnNumber = 7
		filterTableJS(w, filterTableJSArgument)
		if userTypeID == "100" {
			// Call JS filter function for the customer ID in the PBX contact table
			filterTableJSArgument.funcNameJS = "pbxContactSearchCustomerID"
			filterTableJSArgument.inputID = "pbx-contact-input-customer-id"
			filterTableJSArgument.columnNumber = 8
			filterTableJS(w, filterTableJSArgument)
			// Call JS filter function for the customer name in the PBX contact table
			filterTableJSArgument.funcNameJS = "pbxContactSearchCustomerName"
			filterTableJSArgument.inputID = "pbx-contact-input-customer-name"
			filterTableJSArgument.columnNumber = 9
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
	if userTypeID == "100" {
		fmt.Fprintf(w, "    <th class=\"table-title\";>All PBX Resource Limits on the Server:</th>")
	} else {
		fmt.Fprintf(w, "    <th class=\"table-title\";>PBX Resource Limits</th>")
	}
	fmt.Fprintf(w, "  </tr>")
	if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
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
		if userTypeID == "100" {
			inputTableHTMLArgument.inputID = "pbx-resource-input-customer-id"
			inputTableHTMLArgument.funcNameJS = "pbxResourceSearchCustomerID"
			inputTableHTMLArgument.placeholder = "Customer ID"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			fmt.Fprintf(w, "    <br>")
			fmt.Fprintf(w, "    <br>")
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
	if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
		fmt.Fprintf(w, "          <th>PBX ID</th>")
		fmt.Fprintf(w, "          <th>PBX Name</th>")
	}
	fmt.Fprintf(w, "          <th>PBX Date & Time</th>")
	fmt.Fprintf(w, "          <th>SIP Extension<br>Limit for PBX</th>")
	if userTypeID == "100" {
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
		if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
			fmt.Fprintf(w, "          <td>"+pbxID+"</td>")
			fmt.Fprintf(w, "          <td>"+pbxName+"</td>")
		}
		fmt.Fprintf(w, "          <td>"+formatDateTime(pbxDateTimeAdded)+"</td>")
		fmt.Fprintf(w, "          <td>"+pbxSIPExtensionLimit+"</td>")
		if userTypeID == "100" {
			fmt.Fprintf(w, "          <td>"+customerID+"</td>")
			fmt.Fprintf(w, "          <td>"+customerName+"</td>")
		}
		fmt.Fprintf(w, "        </tr>")
	}

	fmt.Fprintf(w, "      </table>")
	if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
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
		if userTypeID == "100" {
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
	if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
		var toggleDivJSArgument jsFunctionParameter
		toggleDivJSArgument.funcNameJS = "togglePBX"
		toggleDivJSArgument.divID = "pbx-div"
		toggleDivJS(w, toggleDivJSArgument)
	}

}

//----------------------------------------------------------------------------------------------------

// Function to create a PBX dialplan table in MariaDB
func createDailplanTable(dbDetail databaseFunctionParameter) {
	_, err := dbDetail.connection.Exec("CREATE TABLE " + dbDetail.table + " (id BIGINT(20) NOT NULL, context VARCHAR(40) NOT NULL, exten VARCHAR(40) NOT NULL, priority INT(11) NOT NULL, app VARCHAR(40) NOT NULL, appdata VARCHAR(256) NOT NULL)")
	if err != nil {
		panic(err)
	}
}

// Function to drop a PBX dialplan table in MariaDB
func dropDailplanTable(dbDetail databaseFunctionParameter) {
	_, err := dbDetail.connection.Exec("DROP TABLE " + dbDetail.table + "")
	if err != nil {
		panic(err)
	}
}

//----------------------------------------------------------------------------------------------------

// SIP extension page functions

func sipExtensionList(w http.ResponseWriter, dbDetail databaseFunctionParameter, userTypeID string, userCustomerID string, userPBXID string) {

	var (
		sipUsername            string
		sipPassword            string
		callerID               string
		callerIDPrivacy        string
		callGroup              string
		codecAllowed           string
		directMedia            string
		directMediaMethod      string
		dtmfMode               string
		forceRPort             string
		fromSIPHeaderUser      string
		fromSIPHeaderDomain    string
		ipAddressAllowed       string
		pickupGroup            string
		mediaEncryptionEnabled string
		stirShakenEnabled      string
		stirShakenProfile      string
		registered             string
		pbxID                  string
		pbxName                string
		customerID             string
		customerName           string
	)

	// Registered table
	var (
		uri       string
		userAgent string
	)

	var dbTableCountUserSIPExtension databaseFunctionParameter
	dbTableCountUserSIPExtension.connection = dbDetail.connection
	dbTableCountUserSIPExtension.database = dbDetail.database
	dbTableCountUserSIPExtension.table = "view___sip_extension_detail"
	dbTableCountUserSIPExtension.columnWhere = "sip_username"

	fmt.Fprintf(w, "<table id=\"table\" class=\"table-sip-extension\">")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>")
	fmt.Fprintf(w, "      <table id=\"table\" class=\"table-sip-extension\">")
	fmt.Fprintf(w, "        <tr>")
	if userTypeID == "100" {
		fmt.Fprintf(w, "          <th>Total SIP Extensions On YAP</th>")
	} else if userTypeID == "200" || userTypeID == "201" {
		fmt.Fprintf(w, "          <th>Total SIP Extensions for the Customer</th>")
	} else if userTypeID == "300" || userTypeID == "301" || userTypeID == "302" {
		fmt.Fprintf(w, "          <th>Total SIP Extensions Within the PBX</th>")
	}
	fmt.Fprintf(w, "        </tr>")
	fmt.Fprintf(w, "        <tr>")
	if userTypeID == "100" {
		dbTableCountUserSIPExtension.countMinusOne = false
		fmt.Fprintf(w, "          <td>"+totalTableCount(w, dbTableCountUserSIPExtension)+"</td>")
	} else if userTypeID == "200" || userTypeID == "201" {
		var dbTableCountUserSIPExtensionWhere databaseFunctionParameter
		dbTableCountUserSIPExtensionWhere.connection = dbDetail.connection
		dbTableCountUserSIPExtensionWhere.database = dbDetail.database
		dbTableCountUserSIPExtensionWhere.table = "view___sip_extension_detail"
		dbTableCountUserSIPExtensionWhere.columnWhere = "customer_id"
		dbTableCountUserSIPExtensionWhere.columnWhereValue = userCustomerID
		fmt.Fprintf(w, "    <td>"+totalTableCountWhere(w, dbTableCountUserSIPExtensionWhere)+"</td>")
	} else if userTypeID == "300" || userTypeID == "301" || userTypeID == "302" {
		var dbTableCountUserSIPExtensionWhere databaseFunctionParameter
		dbTableCountUserSIPExtensionWhere.connection = dbDetail.connection
		dbTableCountUserSIPExtensionWhere.database = dbDetail.database
		dbTableCountUserSIPExtensionWhere.table = "view___sip_extension_detail"
		dbTableCountUserSIPExtensionWhere.columnWhere = "pbx_id"
		dbTableCountUserSIPExtensionWhere.columnWhereValue = userPBXID
		fmt.Fprintf(w, "    <td>"+totalTableCountWhere(w, dbTableCountUserSIPExtensionWhere)+"</td>")
	}
	fmt.Fprintf(w, "        </tr>")
	fmt.Fprintf(w, "      </table>")
	fmt.Fprintf(w, "    </th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th><button onclick=\"toggleSIPExtension() \"class=\"button-general button-sip-extension\">&nbsp Show/Hide SIP Extension &nbsp</button></th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "</table>")

	fmt.Fprintf(w, "<div id=\"sip-extension-div\" style=\"display:none\">")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<table id=\"table\" class=\"table-sip-extension\">")
	fmt.Fprintf(w, "  <tr>")
	if userTypeID == "100" {
		fmt.Fprintf(w, "    <th class=\"table-title\";>All SIP Extension Details on the Server:</th>")
	} else if userTypeID == "200" || userTypeID == "201" {
		fmt.Fprintf(w, "    <th class=\"table-title\";>All SIP Extension Details for the Customer:</th>")
	} else if userTypeID == "300" || userTypeID == "301" || userTypeID == "302" {
		fmt.Fprintf(w, "    <th class=\"table-title\";>All SIP Extension Details Within the PBX:</th>")
	}
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>")
	fmt.Fprintf(w, "    <br>")
	fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
	var inputTableHTMLArgument jsFunctionParameter
	inputTableHTMLArgument.inputID = "sip-extension-detail-input-sip-username"
	inputTableHTMLArgument.funcNameJS = "sipExtensionDetailSearchSIPUsername"
	inputTableHTMLArgument.placeholder = "SIP Username/PBX ID"
	inputTableHTML(w, inputTableHTMLArgument)
	fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
	inputTableHTMLArgument.inputID = "sip-extension-detail-input-option"
	inputTableHTMLArgument.funcNameJS = "sipExtensionDetailSearchOption"
	inputTableHTMLArgument.placeholder = "Options"
	inputTableHTML(w, inputTableHTMLArgument)
	fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
	if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
		inputTableHTMLArgument.inputID = "sip-extension-detail-input-pbx-name"
		inputTableHTMLArgument.funcNameJS = "sipExtensionDetailSearchPBXName"
		inputTableHTMLArgument.placeholder = "PBX Name"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
	}
	if userTypeID == "100" {
		inputTableHTMLArgument.inputID = "sip-extension-detail-input-customer-id"
		inputTableHTMLArgument.funcNameJS = "sipExtensionDetailSearchCustomerID"
		inputTableHTMLArgument.placeholder = "Customer ID"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "sip-extension-detail-input-customer-name"
		inputTableHTMLArgument.funcNameJS = "sipExtensionDetailSearchCustomerName"
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
	fmt.Fprintf(w, "      <table id=\"sip-extension-detail-table\" class=\"table-sip-extension\">")
	fmt.Fprintf(w, "        <tr>")
	fmt.Fprintf(w, "          <th>SIP Username</th>")
	fmt.Fprintf(w, "          <th>SIP Password</th>")
	fmt.Fprintf(w, "          <th>Registered</th>")
	fmt.Fprintf(w, "          <th>Options</th>")
	if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
		fmt.Fprintf(w, "          <th>PBX Name</th>")
	}
	if userTypeID == "100" {
		fmt.Fprintf(w, "          <th>Customer ID</th>")
		fmt.Fprintf(w, "          <th>Customer Name</th>")
	}
	fmt.Fprintf(w, "        </tr>")

	var whereClause string
	var userWhereID string

	if userTypeID == "100" {
		whereClause = "WHERE customer_id != ?;"
		userWhereID = "1"
	} else if userTypeID == "200" || userTypeID == "201" {
		whereClause = "WHERE customer_id = ?;"
		userWhereID = userCustomerID
	} else if userTypeID == "300" || userTypeID == "301" || userTypeID == "302" {
		whereClause = "WHERE pbx_id = ?;"
		userWhereID = userPBXID
	}

	sipExtensionDetailSQL, err := dbDetail.connection.Query(`SELECT
							sip_username,
							sip_password,
							caller_id,
							caller_id_privacy,
							call_group,
							codec_allowed,
							direct_media,
							direct_media_method,
					                dtmf_mode,
					                force_rport,
					                from_sip_header_user,
					                from_sip_Header_domain,
					                ip_address_allowed,
					                pickup_group,
					                media_encryption_enabled,
					                stir_shaken_enabled,
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

	for sipExtensionDetailSQL.Next() {

		err = sipExtensionDetailSQL.Scan(
			&sipUsername,
			&sipPassword,
			&callerID,
			&callerIDPrivacy,
			&callGroup,
			&codecAllowed,
			&directMedia,
			&directMediaMethod,
			&dtmfMode,
			&forceRPort,
			&fromSIPHeaderUser,
			&fromSIPHeaderDomain,
			&ipAddressAllowed,
			&pickupGroup,
			&mediaEncryptionEnabled,
			&stirShakenEnabled,
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
		copyButtonJSArgument.buttonCSS = "button-sip-extension"
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
		fmt.Fprintf(w, "          <b>Caller ID:</b> "+callerID+"<br>")
		fmt.Fprintf(w, "          <b>Caller ID Privacy:</b> "+callerIDPrivacy+"<br>")
		fmt.Fprintf(w, "          <b>Call Group:</b> "+callGroup+"<br>")
		fmt.Fprintf(w, "          <b>Codecs Allowed:</b> "+codecAllowed+"<br>")
		fmt.Fprintf(w, "          <b>Direct Media:</b> "+directMedia+"<br>")
		fmt.Fprintf(w, "          <b>Direct Media Method:</b> "+directMediaMethod+"<br>")
		fmt.Fprintf(w, "          <b>DTMF Mode:</b> "+dtmfMode+"<br>")
		fmt.Fprintf(w, "          <b>Force Rport:</b> "+forceRPort+"<br>")
		fmt.Fprintf(w, "          <b>From SIP Header User:</b> "+fromSIPHeaderUser+"<br>")
		fmt.Fprintf(w, "          <b>From SIP Header Domain:</b> "+fromSIPHeaderDomain+"<br>")
		fmt.Fprintf(w, "          <b>IP Address Allowed:</b> "+ipAddressAllowed+"<br>")
		fmt.Fprintf(w, "          <b>Pickup Group:</b> "+pickupGroup+"<br>")
		fmt.Fprintf(w, "          <b>Media Encryption Enabled:</b> "+mediaEncryptionEnabled+"<br>")
		fmt.Fprintf(w, "          <b>STIR/SHAKEN Enabled:</b> "+stirShakenEnabled+"<br>")
		fmt.Fprintf(w, "          <b>STIR/SHAKEN Profile:</b> "+stirShakenProfile+"<br>")
		fmt.Fprintf(w, "          </td>")
		if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
			fmt.Fprintf(w, "          <td>"+pbxName+"</td>")
		}
		if userTypeID == "100" {
			fmt.Fprintf(w, "          <td>"+customerID+"</td>")
			fmt.Fprintf(w, "          <td>"+customerName+"</td>")
		}
		fmt.Fprintf(w, "        </tr>")
	}

	fmt.Fprintf(w, "      </table>")
	var filterTableJSArgument jsFunctionParameter
	filterTableJSArgument.tableID = "sip-extension-detail-table"
	// Call JS filter function for SIP username in the SIP extension detail table
	filterTableJSArgument.funcNameJS = "sipExtensionDetailSearchSIPUsername"
	filterTableJSArgument.inputID = "sip-extension-detail-input-sip-username"
	filterTableJSArgument.columnNumber = 0
	filterTableJS(w, filterTableJSArgument)
	// Call JS filter function for options in the SIP extension detail table
	filterTableJSArgument.funcNameJS = "sipExtensionDetailSearchOption"
	filterTableJSArgument.inputID = "sip-extension-detail-input-option"
	filterTableJSArgument.columnNumber = 3
	filterTableJS(w, filterTableJSArgument)
	if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
		// Call JS filter function for PBX name in the SIP extension detail table
		filterTableJSArgument.funcNameJS = "sipExtensionDetailSearchPBXName"
		filterTableJSArgument.inputID = "sip-extension-detail-input-pbx-name"
		filterTableJSArgument.columnNumber = 4
		filterTableJS(w, filterTableJSArgument)
	}
	if userTypeID == "100" {
		// Call JS filter function for the customer ID in the SIP extension detail table
		filterTableJSArgument.funcNameJS = "sipExtensionDetailSearchCustomerID"
		filterTableJSArgument.inputID = "sip-extension-detail-input-customer-id"
		filterTableJSArgument.columnNumber = 5
		filterTableJS(w, filterTableJSArgument)
		// Call JS filter function for the customer name in the SIP extension detail table
		filterTableJSArgument.funcNameJS = "sipExtensionDetailSearchCustomerName"
		filterTableJSArgument.inputID = "sip-extension-detail-input-customer-name"
		filterTableJSArgument.columnNumber = 6
		filterTableJS(w, filterTableJSArgument)
	}
	fmt.Fprintf(w, "    </th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "</table>")
	fmt.Fprintf(w, "<br>")

	// Registered SIP extension Table
	fmt.Fprintf(w, "<table id=\"table\" class=\"table-sip-extension\">")
	fmt.Fprintf(w, "  <tr>")
	if userTypeID == "100" {
		fmt.Fprintf(w, "    <th class=\"table-title\";>All SIP Extensions Registered on the Server:</th>")
	} else if userTypeID == "200" || userTypeID == "201" {
		fmt.Fprintf(w, "    <th class=\"table-title\";>All SIP Extensions Registered for the Customer:</th>")
	} else if userTypeID == "300" || userTypeID == "301" || userTypeID == "302" {
		fmt.Fprintf(w, "    <th class=\"table-title\";>All SIP Extensions Registered Within the PBX:</th>")
	}
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>")
	fmt.Fprintf(w, "    <br>")
	fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
	inputTableHTMLArgument.inputID = "sip-extension-reg-input-sip-username"
	inputTableHTMLArgument.funcNameJS = "sipExtensionRegSearchSIPUsername"
	inputTableHTMLArgument.placeholder = "SIP Username/PBX ID"
	inputTableHTML(w, inputTableHTMLArgument)
	fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
	inputTableHTMLArgument.inputID = "sip-extension-reg-input-uri"
	inputTableHTMLArgument.funcNameJS = "sipExtensionRegSearchURI"
	inputTableHTMLArgument.placeholder = "URI"
	inputTableHTML(w, inputTableHTMLArgument)
	fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
	inputTableHTMLArgument.inputID = "sip-extension-reg-input-user-agent"
	inputTableHTMLArgument.funcNameJS = "sipExtensionRegSearchUserAgent"
	inputTableHTMLArgument.placeholder = "User Agent"
	inputTableHTML(w, inputTableHTMLArgument)
	fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
	if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
		inputTableHTMLArgument.inputID = "sip-extension-reg-input-pbx-name"
		inputTableHTMLArgument.funcNameJS = "sipExtensionRegSearchPBXName"
		inputTableHTMLArgument.placeholder = "PBX Name"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
	}
	if userTypeID == "100" {
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "sip-extension-reg-input-customer-id"
		inputTableHTMLArgument.funcNameJS = "sipExtensionRegSearchCustomerID"
		inputTableHTMLArgument.placeholder = "Customer ID"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "sip-extension-reg-input-customer-name"
		inputTableHTMLArgument.funcNameJS = "sipExtensionRegSearchCustomerName"
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
	fmt.Fprintf(w, "      <table id=\"sip-extension-reg-table\" class=\"table-sip-extension\">")
	fmt.Fprintf(w, "        <tr>")
	fmt.Fprintf(w, "          <th>SIP Username</th>")
	fmt.Fprintf(w, "          <th>URI</th>")
	fmt.Fprintf(w, "          <th>User Agent<br>SIP Client</th>")
	if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
		fmt.Fprintf(w, "          <th>PBX Name</th>")
	}
	if userTypeID == "100" {
		fmt.Fprintf(w, "          <th>Customer ID</th>")
		fmt.Fprintf(w, "          <th>Customer Name</th>")
	}
	fmt.Fprintf(w, "        </tr>")

	sipExtensionRegSQL, err := dbDetail.connection.Query(`SELECT
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

	for sipExtensionRegSQL.Next() {

		err = sipExtensionRegSQL.Scan(
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
	filterTableJSArgument.tableID = "sip-extension-reg-table"
	// Call JS filter function for SIP username in the SIP extension registration (reg) table
	filterTableJSArgument.funcNameJS = "sipExtensionRegSearchSIPUsername"
	filterTableJSArgument.inputID = "sip-extension-reg-input-sip-username"
	filterTableJSArgument.columnNumber = 0
	filterTableJS(w, filterTableJSArgument)
	// Call JS filter function for URI in the SIP extension registration (reg) table
	filterTableJSArgument.funcNameJS = "sipExtensionRegSearchURI"
	filterTableJSArgument.inputID = "sip-extension-reg-input-uri"
	filterTableJSArgument.columnNumber = 1
	filterTableJS(w, filterTableJSArgument)
	// Call JS filter function for user agent in the SIP extension registration (reg) table
	filterTableJSArgument.funcNameJS = "sipExtensionRegSearchUserAgent"
	filterTableJSArgument.inputID = "sip-extension-reg-input-user-agent"
	filterTableJSArgument.columnNumber = 2
	filterTableJS(w, filterTableJSArgument)
	if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
		// Call JS filter function for PBX name in the SIP extension registration (reg) table
		filterTableJSArgument.funcNameJS = "sipExtensionRegSearchPBXName"
		filterTableJSArgument.inputID = "sip-extension-reg-input-pbx-name"
		filterTableJSArgument.columnNumber = 3
		filterTableJS(w, filterTableJSArgument)
	}
	if userTypeID == "100" {
		// Call JS filter function for the customer ID in the SIP extension registration (reg) table
		filterTableJSArgument.funcNameJS = "sipExtensionRegSearchCustomerID"
		filterTableJSArgument.inputID = "sip-extension-reg-input-customer-id"
		filterTableJSArgument.columnNumber = 4
		filterTableJS(w, filterTableJSArgument)
		// Call JS filter function for the customer name in the SIP extension registration (reg) table
		filterTableJSArgument.funcNameJS = "sipExtensionRegSearchCustomerName"
		filterTableJSArgument.inputID = "sip-extension-reg-input-customer-name"
		filterTableJSArgument.columnNumber = 5
		filterTableJS(w, filterTableJSArgument)
	}
	fmt.Fprintf(w, "    </th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "</table>")
	fmt.Fprintf(w, "</div>")
	var toggleDivJSArgument jsFunctionParameter
	toggleDivJSArgument.funcNameJS = "toggleSIPExtension"
	toggleDivJSArgument.divID = "sip-extension-div"
	toggleDivJS(w, toggleDivJSArgument)
}

//----------------------------------------------------------------------------------------------------

// Invoice page functions

func invoiceList(w http.ResponseWriter, dbDetail databaseFunctionParameter, userTypeID string, userCustomerID string, yapAdminUKVATRegistered string) {

	var (
		customerName                 string
		customerID                   string
		customerUKBased              string
		customerResellingMinutes     string
		customerUKVATRegistered      string
		customerUKVATNumber          string
		invoiceItemTag               string
		invoiceItemSellPrice         string
		invoiceItemDateTimeAdded     string
		invoiceItemSalesTaxRate      string
		invoiceItemSalesTaxStatus    string
		invoiceBillItemOnce          string
		invoiceItemOnHold            string
		invoiceItemContractLength    string
		invoiceItemContractStartDate string
		goodServiceName              string
		goodServiceType              string
		goodServiceSupplierName      string
		goodServiceContractLength    string
	)

	var dbTableCountInvoice databaseFunctionParameter
	dbTableCountInvoice.connection = dbDetail.connection
	dbTableCountInvoice.database = dbDetail.database
	dbTableCountInvoice.table = "view___invoice_item"

	if userTypeID == "100" {
		fmt.Fprintf(w, "<table id=\"table\" class=\"table-invoice\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table id=\"table\" class=\"table-invoice\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <th>Total Invoice Services/Products</th>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		dbTableCountInvoice.countMinusOne = false
		fmt.Fprintf(w, "          <td>"+totalTableCount(w, dbTableCountInvoice)+"</td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><button onclick=\"toggleInvoice() \"class=\"button-general button-invoice\">&nbsp Show/Hide Invoice &nbsp</button></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
	}

	if userTypeID == "100" {
		fmt.Fprintf(w, "<div id=\"invoice-div\" style=\"display:none\">")
		fmt.Fprintf(w, "<br>")
	} else {
		fmt.Fprintf(w, "<div id=\"invoice-div\">")
	}
	fmt.Fprintf(w, "<table id=\"table\" class=\"table-invoice\">")
	fmt.Fprintf(w, "  <tr>")
	if userTypeID == "100" {
		fmt.Fprintf(w, "    <th class=\"table-title\";>All Customer Service/Product Invoice Items on YAP:</th>")
	} else {
		fmt.Fprintf(w, "    <th class=\"table-title\";>Customer Invoice Services/Products</th>")
	}
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>")
	fmt.Fprintf(w, "    <br>")

	var inputTableHTMLArgument jsFunctionParameter
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
	if userTypeID == "100" {
		inputTableHTMLArgument.inputID = "invoice-input-detail"
		inputTableHTMLArgument.funcNameJS = "invoiceSearchDetail"
		inputTableHTMLArgument.placeholder = "Invoice Item Details"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "invoice-input-customer-name"
		inputTableHTMLArgument.funcNameJS = "invoiceSearchCustomerName"
		inputTableHTMLArgument.placeholder = "Customer Name"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		fmt.Fprintf(w, "    <br><br>")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "invoice-input-customer-id"
		inputTableHTMLArgument.funcNameJS = "invoiceSearchCustomerID"
		inputTableHTMLArgument.placeholder = "Customer ID"
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
	fmt.Fprintf(w, "          <th>Service/Product<br>Name & Information</th>")
	fmt.Fprintf(w, "          <th>Service/Product<br>Sale Price</th>")
	if userTypeID == "100" {
		fmt.Fprintf(w, "          <th>Service/Product<br>Details</th>")
		fmt.Fprintf(w, "          <th>Customer<br>Name</th>")
		fmt.Fprintf(w, "          <th>Customer<br>ID</th>")
	}
	fmt.Fprintf(w, "        </tr>")

	var whereClause string

	if userTypeID == "100" {
		whereClause = "WHERE customer_id != ?;"
		userCustomerID = "1"
	} else if userTypeID == "200" || userTypeID == "400" {
		whereClause = "WHERE customer_id = ?;"
	}

	invoiceSQL, err := dbDetail.connection.Query(`SELECT
			                                customer_name,
                                                        customer_id,
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
							good_service_name,
							good_service_type,
							good_service_supplier_name,
							good_service_contract_length
					              FROM
					  	        yap.view___invoice_item
						      `+whereClause, userCustomerID)

	// Error
	if err != nil {
		panic(err)

	}

	for invoiceSQL.Next() {

		err = invoiceSQL.Scan(
			&customerName,
			&customerID,
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
			&goodServiceName,
			&goodServiceType,
			&goodServiceSupplierName,
			&goodServiceContractLength,
		)

		// Error
		if err != nil {
			panic(err)
		}
		if userTypeID != "100" && invoiceItemOnHold == "yes" {
		} else {
			fmt.Fprintf(w, "        <tr>")
			fmt.Fprintf(w, "          <td style=\"text-align: left;\">")
			fmt.Fprintf(w, "            "+goodServiceName+"<br><br>")
			fmt.Fprintf(w, "            <b>Service/Product Tag:</b> "+invoiceItemTag+"<br>")
			if invoiceItemContractLength == "n/a" {
				fmt.Fprintf(w, "          <b>Contract Start Date:</b> n/a<br><b>Contract Length:</b> n/a")
			} else {
				fmt.Fprintf(w, "          <b>Contract Start Date:</b> "+formatDate(invoiceItemContractStartDate)+"<br><b>Contract Length:</b> "+invoiceItemContractLength)
			}
			fmt.Fprintf(w, "          </td>")
			fmt.Fprintf(w, "          <td style=\"text-align: left;\">")
			//If the YAP Admin is UK VAT registered and the service/product is taxable
			if yapAdminUKVATRegistered == "yes" && invoiceItemSalesTaxStatus == "TAXABLE" {
				fmt.Fprintf(w, "            <b>Price (exVAT):</b> £"+invoiceItemSellPrice+"<br>")
				// Convert UK sales VAT rate to float64
				invoiceItemSalesTaxRateFloat64, _ := strconv.ParseFloat(invoiceItemSalesTaxRate, 64)
				fmt.Fprintf(w, "            <b>VAT Rate:</b> "+strconv.FormatFloat(invoiceItemSalesTaxRateFloat64, 'f', -1, 64)+"&#37;<br>")
				// Convert item sell price to float64
				invoiceItemSellPriceExVATFloat64, _ := strconv.ParseFloat(invoiceItemSellPrice, 64)
				var invoiceItemSellPriceIncVATFloat64 float64 = invoiceItemSellPriceExVATFloat64 * (invoiceItemSalesTaxRateFloat64/100 + 1)
				var invoiceItemSellVATFloat64 float64 = invoiceItemSellPriceIncVATFloat64 - invoiceItemSellPriceExVATFloat64
				fmt.Fprintf(w, "            <b>VAT:</b> £"+strconv.FormatFloat(invoiceItemSellVATFloat64, 'f', 2, 64)+"<br>")
				fmt.Fprintf(w, "            <b>Total Price (incVAT):</b> £"+strconv.FormatFloat(invoiceItemSellPriceIncVATFloat64, 'f', 2, 64)+"<br>")
				//If the YAP Admin is UK VAT registered and the service/product is exempt
			} else if yapAdminUKVATRegistered == "yes" && invoiceItemSalesTaxStatus == "EXEMPT" {
				fmt.Fprintf(w, "            <b>Price (exVAT):</b> £"+invoiceItemSellPrice+"<br>")
				fmt.Fprintf(w, "            <b>VAT Rate:</b> Exempt<br>")
				fmt.Fprintf(w, "            <b>VAT:</b> £0.00</b><br>")
				fmt.Fprintf(w, "            <b>Total Price (incVAT):</b> £"+invoiceItemSellPrice)
				//If the YAP Admin is not UK VAT registered
			} else if yapAdminUKVATRegistered == "no" {
				fmt.Fprintf(w, "            £"+invoiceItemSellPrice+"<br>")
			}
			fmt.Fprintf(w, "          </td>")
			if userTypeID == "100" {
				fmt.Fprintf(w, "          <td style=\"text-align: left; vertical-align: top;\">")
				fmt.Fprintf(w, "            <b><u>Item Details</u></b><br><br>")
				fmt.Fprintf(w, " 	    <b>Item Added Date & Time: </b>"+formatDateTime(invoiceItemDateTimeAdded)+"<br>")
				fmt.Fprintf(w, "            <b>Sale VAT Status: </b>"+invoiceItemSalesTaxStatus+"<br>")
				fmt.Fprintf(w, "            <b>Bill Item Once: </b>"+invoiceBillItemOnce+"<br>")
				fmt.Fprintf(w, "            <b>Item on Hold: </b>"+invoiceItemOnHold+"<br>")
				fmt.Fprintf(w, "            <b>Item Type: </b>"+goodServiceType+"<br>")
				fmt.Fprintf(w, "            <hr class=\"line-table\"></h>")
				fmt.Fprintf(w, "            <b><u>Customer Details</u></b><br><br>")
				fmt.Fprintf(w, "            <b>Reselling Minutes: </b>"+customerResellingMinutes+"<br>")
				fmt.Fprintf(w, "            <b>UK Based: </b>"+customerUKBased+"<br>")
				fmt.Fprintf(w, "            <b>UK VAT Registered: </b>"+customerUKVATRegistered+"<br>")
				fmt.Fprintf(w, "            <b>UK VAT Number: </b>"+customerUKVATNumber+"<br>")
				fmt.Fprintf(w, "            <hr class=\"line-table\"></h>")
				fmt.Fprintf(w, "            <b><u>Supplier Details</u></b><br><br>")
				fmt.Fprintf(w, "            <b>Name: </b>"+goodServiceSupplierName+"<br>")
				fmt.Fprintf(w, "            <b>Supplier Contract Length: </b>"+goodServiceContractLength)
				fmt.Fprintf(w, "          </td>")
				fmt.Fprintf(w, "          <td>"+customerName+"</td>")
				fmt.Fprintf(w, "          <td>"+customerID+"</td>")
			}
			fmt.Fprintf(w, "        </tr>")
		}
	}

	fmt.Fprintf(w, "      </table>")
	var filterTableJSArgument jsFunctionParameter
	filterTableJSArgument.tableID = "invoice-table"

	filterTableJSArgument.funcNameJS = "invoiceSearchNameInformation"
	filterTableJSArgument.inputID = "invoice-input-name-information"
	filterTableJSArgument.columnNumber = 0
	filterTableJS(w, filterTableJSArgument)

	filterTableJSArgument.funcNameJS = "invoiceSearchSalePrice"
	filterTableJSArgument.inputID = "invoice-input-sale-price"
	filterTableJSArgument.columnNumber = 1
	filterTableJS(w, filterTableJSArgument)

	if userTypeID == "100" {
		filterTableJSArgument.funcNameJS = "invoiceSearchDetail"
		filterTableJSArgument.inputID = "invoice-input-detail"
		filterTableJSArgument.columnNumber = 2
		filterTableJS(w, filterTableJSArgument)

		filterTableJSArgument.funcNameJS = "invoiceSearchCustomerName"
		filterTableJSArgument.inputID = "invoice-input-customer-name"
		filterTableJSArgument.columnNumber = 3
		filterTableJS(w, filterTableJSArgument)

		filterTableJSArgument.funcNameJS = "invoiceSearchCustomerID"
		filterTableJSArgument.inputID = "invoice-input-customer-id"
		filterTableJSArgument.columnNumber = 4
		filterTableJS(w, filterTableJSArgument)
	}
	var exportCSVJSArgument jsFunctionParameter
	exportCSVJSArgument.funcNameJS = "Invoice"
	exportCSVJSArgument.tableID = "invoice-table"
	exportCSVJSArgument.fileName = "YAP_customer_contact_details"
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
}

//----------------------------------------------------------------------------------------------------

// Server log page functions

//----------------------------------------------------------------------------------------------------

// Server information functions

//----------------------------------------------------------------------------------------------------

func main() {

	//Get the values from inside the YAP configuration file
	err := godotenv.Load("/etc/yap/yap.env")
	if err != nil {
		panic("Error loading yap.env file for database details")
	}

	//Get the database connection details
	dbUsername := os.Getenv("dbUsername")
	dbPassword := os.Getenv("dbPassword")
	dbName := os.Getenv("dbName")
	dbTransport := os.Getenv("dbTransport")
	dbAddress := os.Getenv("dbAddress")
	dbPort := os.Getenv("dbPort")
	dbTls := os.Getenv("dbTls")
	extraButtonName := os.Getenv("extraButtonName")
	extraButtonURL := os.Getenv("extraButtonURL")
	yapAdminUKVATRegistered := os.Getenv("yapAdminUKVATRegistered")

	//Values allowed for dbTransport Variable
	var allowedTransportValue = []string{"tcp", "udp"}
	validDbTransport := slices.Contains(allowedTransportValue, dbTransport)

	validateDbAddress := validator.New()
	validateDbAddressErr := validateDbAddress.Var(dbAddress, "required,ip_addr")

	dbPortInt, err := strconv.Atoi(dbPort)
	if err != nil {
		panic("DATABASE PORT MUST BE A NUMBER IN /etc/yap/yap.env")
	}

	//Values allowed for dbTls Variable
	var allowedTlsValue = []string{"false", "true"}
	validDbTls := slices.Contains(allowedTlsValue, dbTls)

	validateExtraButtonURL := validator.New()
	validateExtraButtonURLErr := validateExtraButtonURL.Var(extraButtonURL, "required,http_url")

	//Values allowed for ukVATRegistered Variable
	var allowedYAPAdminUKVATRegisteredValue = []string{"no", "yes"}
	validYAPAdminUKVATRegistered := slices.Contains(allowedYAPAdminUKVATRegisteredValue, yapAdminUKVATRegistered)

	//Catch if any errors were made in yap.env and feed back where to correct error
	if dbUsername == "" {
		panic("DATABASE USERNAME CANNOT BE BLANK IN /etc/yap/yap.env")
	} else if dbPassword == "" {
		panic("DATABASE PASSOWRD CANNOT BE BLANK IN /etc/yap/yap.env")
	} else if dbName == "" {
		panic("DATABASE NAME CANNOT BE BLANK IN /etc/yap/yap.env")
	} else if dbTransport == "" {
		panic("DATABASE TRANSPORT OPTION CANNOT BE BLANK IN /etc/yap/yap.env")
	} else if validDbTransport == false {
		panic("DATABASE TRANSPORT OPTION MUST BE udp OR tcp IN /etc/yap/yap.env")
	} else if validateDbAddressErr != nil && dbAddress != "localhost" {
		panic("DATABASE ADDRESS MUST BE A VALID INTERENT PROTOCOL (IP) ADDRESS OR localhost IN /etc/yap/yap.env")
	} else if dbPortInt <= 0 || dbPortInt >= 65536 {
		panic("DATABASE PORT MUST BE IN THE NUMBER RANGE 1-65535 IN /etc/yap/yap.env")
	} else if dbTls == "" {
		panic("DATABASE TLS OPTION CANNOT BE BLANK IN /etc/yap/yap.env")
	} else if validDbTls == false {
		panic("DATABASE TRANSPORT OPTION MUST BE false OR true IN /etc/yap/yap.env")
	} else if extraButtonName != "" {
		if validateExtraButtonURLErr != nil {
			panic("THE EXTRA BUTTON URL VALUE MUST BE A VALID URL IN /etc/yap/yap.env")
		}
	} else if yapAdminUKVATRegistered == "" {
		panic("UK YAP ADMIN VAT REGISTERED OPTION CANNOT BE BLANK IN /etc/yap/yap.env")
	} else if validYAPAdminUKVATRegistered == false {
		panic("UK YAP ADMIN VAT REGISTERED OPTION MUST BE no OR yes IN /etc/yap/yap.env")
	}

	startHTML := csvcell.FileData(dirHTML, fileStartHTML)
	endHTML := csvcell.FileData(dirHTML, fileEndHTML)

	// Home Page
	go http.HandleFunc("/yap", func(w http.ResponseWriter, r *http.Request) {

		// Open database connection
		dbConnection, err := sql.Open("mysql", dbUsername+":"+dbPassword+"@"+dbTransport+"("+dbAddress+":"+dbPort+")/"+dbName+"?tls="+dbTls)
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

		if userTypeID == "" {
			errorBox(w, "email_error", "header-main-menu", "")
		} else {
			if userTypeID == "100" {
				header(w, "YAP Admin Account<br>Main Menu", "", extraButtonName, extraButtonURL)
				mainMenuUserInformation(w, dbDetail, userTypeID)
				fmt.Fprintf(w, "<br>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonOne := mainMenuParameter{writeHTTP: w, buttonName: "Customer & PBX<br>User Accounts<br>&#128100", hyperlink: "/user-account", headerCSS: "header-user-account", buttonCSS: "button-user-account"}
				mainMenuButton(mainMenuButtonOne)
				mainMenuButtonTwo := mainMenuParameter{writeHTTP: w, buttonName: "All<br>Customers<br>&#128101", hyperlink: "/customer", headerCSS: "header-customer", buttonCSS: "button-customer"}
				mainMenuButton(mainMenuButtonTwo)
				mainMenuButtonThree := mainMenuParameter{writeHTTP: w, buttonName: "All<br>PBXs<br>&#128222", hyperlink: "/pbx", headerCSS: "header-pbx", buttonCSS: "button-pbx"}
				mainMenuButton(mainMenuButtonThree)
				mainMenuButtonFour := mainMenuParameter{writeHTTP: w, buttonName: "All PBX SIP<br>Extensions<br>&#128241", hyperlink: "/sip-extension", headerCSS: "header-sip-extension", buttonCSS: "button-sip-extension"}
				mainMenuButton(mainMenuButtonFour)
				fmt.Fprintf(w, "</div>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonFive := mainMenuParameter{writeHTTP: w, buttonName: "All Customer<br>Invoicing<br>&#129534", hyperlink: "/invoice", headerCSS: "header-invoice", buttonCSS: "button-invoice"}
				mainMenuButton(mainMenuButtonFive)
				mainMenuButtonSix := mainMenuParameter{writeHTTP: w, buttonName: "All Server<br>Logs<br>&#128221", hyperlink: "/server-log", headerCSS: "header-server-log", buttonCSS: "button-server-log"}
				mainMenuButton(mainMenuButtonSix)
				mainMenuButtonSeven := mainMenuParameter{writeHTTP: w, buttonName: "YAP Server<br>Information<br>&#128421", hyperlink: "/server-information", headerCSS: "header-server-information", buttonCSS: "button-server-information"}
				mainMenuButton(mainMenuButtonSeven)
				fmt.Fprintf(w, "</div>")
				footer(w, "", "")
			} else if userTypeID == "200" {
				header(w, userCustomerName+"<br>[Customer ID: "+userCustomerID+"]<br>Main Menu", "", extraButtonName, extraButtonURL)
				mainMenuUserInformation(w, dbDetail, userTypeID)
				fmt.Fprintf(w, "<br>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonOne := mainMenuParameter{writeHTTP: w, buttonName: "Customer & PBX<br>User Accounts<br>&#128100", hyperlink: "/user-account", headerCSS: "header-user-account", buttonCSS: "button-user-account"}
				mainMenuButton(mainMenuButtonOne)
				mainMenuButtonTwo := mainMenuParameter{writeHTTP: w, buttonName: "Customer<br>Information<br>&#128101", hyperlink: "/customer", headerCSS: "header-customer", buttonCSS: "button-customer"}
				mainMenuButton(mainMenuButtonTwo)
				mainMenuButtonThree := mainMenuParameter{writeHTTP: w, buttonName: "PBXs for the<br>Customer<br>&#128222", hyperlink: "/pbx", headerCSS: "header-pbx", buttonCSS: "button-pbx"}
				mainMenuButton(mainMenuButtonThree)
				fmt.Fprintf(w, "</div>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonFour := mainMenuParameter{writeHTTP: w, buttonName: "PBX SIP<br>Extensions<br>&#128241", hyperlink: "/sip-extension", headerCSS: "header-sip-extension", buttonCSS: "button-sip-extension"}
				mainMenuButton(mainMenuButtonFour)
				mainMenuButtonFive := mainMenuParameter{writeHTTP: w, buttonName: "Customer<br>Invoice<br>&#129534", hyperlink: "/invoice", headerCSS: "header-invoice", buttonCSS: "button-invoice"}
				mainMenuButton(mainMenuButtonFive)
				fmt.Fprintf(w, "</div>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonSix := mainMenuParameter{writeHTTP: w, buttonName: "Customer & PBX<br>Server Logs<br>&#128221", hyperlink: "/server-log", headerCSS: "header-server-log", buttonCSS: "button-server-log"}
				mainMenuButton(mainMenuButtonSix)
				fmt.Fprintf(w, "</div>")
				footer(w, "", "")
			} else if userTypeID == "201" {
				header(w, userCustomerName+"<br>[Customer ID: "+userCustomerID+"]<br>Main Menu", "", extraButtonName, extraButtonURL)
				mainMenuUserInformation(w, dbDetail, userTypeID)
				fmt.Fprintf(w, "<br>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonOne := mainMenuParameter{writeHTTP: w, buttonName: "Own & PBX<br>User Accounts<br>&#128100", hyperlink: "/user-account", headerCSS: "header-user-account", buttonCSS: "button-user-account"}
				mainMenuButton(mainMenuButtonOne)
				mainMenuButtonTwo := mainMenuParameter{writeHTTP: w, buttonName: "Customer<br>Information<br>&#128101", hyperlink: "/customer", headerCSS: "header-customer", buttonCSS: "button-customer"}
				mainMenuButton(mainMenuButtonTwo)
				mainMenuButtonThree := mainMenuParameter{writeHTTP: w, buttonName: "PBXs for the<br>Customer<br>&#128222", hyperlink: "/pbx", headerCSS: "header-pbx", buttonCSS: "button-pbx"}
				mainMenuButton(mainMenuButtonThree)
				fmt.Fprintf(w, "</div>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonFour := mainMenuParameter{writeHTTP: w, buttonName: "PBX SIP<br>Extensions<br>&#128241", hyperlink: "/sip-extension", headerCSS: "header-sip-extension", buttonCSS: "button-sip-extension"}
				mainMenuButton(mainMenuButtonFour)
				mainMenuButtonFive := mainMenuParameter{writeHTTP: w, buttonName: "Customer<br>Invoice<br>&#129534", hyperlink: "/invoice", headerCSS: "header-invoice", buttonCSS: "button-invoice"}
				mainMenuButton(mainMenuButtonFive)
				fmt.Fprintf(w, "</div>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonSix := mainMenuParameter{writeHTTP: w, buttonName: "PBX<br>Server Log<br>&#128221", hyperlink: "/server-log", headerCSS: "header-server-log", buttonCSS: "button-server-log"}
				mainMenuButton(mainMenuButtonSix)
				fmt.Fprintf(w, "</div>")
				footer(w, "", "")
			} else if userTypeID == "300" {
				header(w, userPBXName+"<br>[PBX ID: "+userPBXID+"]<br>Main Menu", "", extraButtonName, extraButtonURL)
				mainMenuUserInformation(w, dbDetail, userTypeID)
				fmt.Fprintf(w, "<br>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonOne := mainMenuParameter{writeHTTP: w, buttonName: "Own & PBX<br>User Accounts<br>&#128100", hyperlink: "/user-account", headerCSS: "header-user-account", buttonCSS: "button-user-account"}
				mainMenuButton(mainMenuButtonOne)
				mainMenuButtonTwo := mainMenuParameter{writeHTTP: w, buttonName: "PBX<br>Information<br>&#128222", hyperlink: "/pbx", headerCSS: "header-pbx", buttonCSS: "button-pbx"}
				mainMenuButton(mainMenuButtonTwo)
				mainMenuButtonThree := mainMenuParameter{writeHTTP: w, buttonName: "PBX SIP<br>Extensions<br>&#128241", hyperlink: "/sip-extension", headerCSS: "header-sip-extension", buttonCSS: "button-sip-extension"}
				mainMenuButton(mainMenuButtonThree)
				mainMenuButtonFour := mainMenuParameter{writeHTTP: w, buttonName: "PBX<br>Server Log<br>&#128221", hyperlink: "/server-log", headerCSS: "header-server-log", buttonCSS: "button-server-log"}
				mainMenuButton(mainMenuButtonFour)
				fmt.Fprintf(w, "</div>")
				footer(w, "", "")
			} else if userTypeID == "301" || userTypeID == "302" {
				header(w, userPBXName+"<br>[PBX ID: "+userPBXID+"]<br>Main Menu", "", extraButtonName, extraButtonURL)
				mainMenuUserInformation(w, dbDetail, userTypeID)
				fmt.Fprintf(w, "<br>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonOne := mainMenuParameter{writeHTTP: w, buttonName: "Own<br>User Account<br>&#128100", hyperlink: "/user-account", headerCSS: "header-user-account", buttonCSS: "button-user-account"}
				mainMenuButton(mainMenuButtonOne)
				mainMenuButtonTwo := mainMenuParameter{writeHTTP: w, buttonName: "PBX<br>Information<br>&#128222", hyperlink: "/pbx", headerCSS: "header-pbx", buttonCSS: "button-pbx"}
				mainMenuButton(mainMenuButtonTwo)
				mainMenuButtonThree := mainMenuParameter{writeHTTP: w, buttonName: "PBX SIP<br>Extensions<br>&#128241", hyperlink: "/sip-extension", headerCSS: "header-sip-extension", buttonCSS: "button-sip-extension"}
				mainMenuButton(mainMenuButtonThree)
				fmt.Fprintf(w, "</div>")
				footer(w, "", "")
			} else if userTypeID == "400" {
				header(w, userCustomerName+"<br>[Customer ID: "+userCustomerID+"]<br>Main Menu", "", extraButtonName, extraButtonURL)
				mainMenuUserInformation(w, dbDetail, userTypeID)
				fmt.Fprintf(w, "<br>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonTwo := mainMenuParameter{writeHTTP: w, buttonName: "Own<br>User Account<br>&#128100", hyperlink: "/user-account", headerCSS: "header-user-account", buttonCSS: "button-user-account"}
				mainMenuButton(mainMenuButtonTwo)
				mainMenuButtonOne := mainMenuParameter{writeHTTP: w, buttonName: "Customer<br>Invoice<br>&#129534", hyperlink: "/invoice", headerCSS: "header-invoice", buttonCSS: "button-invoice"}
				mainMenuButton(mainMenuButtonOne)
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
		dbConnection, err := sql.Open("mysql", dbUsername+":"+dbPassword+"@"+dbTransport+"("+dbAddress+":"+dbPort+")/"+dbName+"?tls="+dbTls)
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

		if userTypeID == "" {
			errorBox(w, "email_error", "header-user-account", "button-user-account")
		} else {
			if userTypeID == "100" {
				header(w, "YAP Admin Account<br>All User Accounts on YAP", "header-user-account", extraButtonName, extraButtonURL)
				userAccountList(w, dbDetail, userTypeID)
				fmt.Fprint(w, "<br>")
				userAccountAdd(w, dbDetail, userID, r)
				fmt.Fprint(w, "<br>")
				userAccountDelete(w, dbDetail, userID, r)
				footer(w, "header-user-account", "button-user-account")
			} else if userTypeID == "200" {
				header(w, userCustomerName+"<br>[Customer ID: "+userCustomerID+"]<br>All User Accounts for the Customer", "header-user-account", extraButtonName, extraButtonURL)
				userAccountList(w, dbDetail, userTypeID)
				footer(w, "header-user-account", "button-user-account")
			} else if userTypeID == "201" {
				header(w, userCustomerName+"<br>[Customer ID: "+userCustomerID+"]<br>All PBX User Accounts for the Customer", "header-user-account", extraButtonName, extraButtonURL)
				userAccountList(w, dbDetail, userTypeID)
				footer(w, "header-user-account", "button-user-account")
			} else if userTypeID == "300" {
				header(w, userPBXName+"<br>[PBX ID: "+userPBXID+"]<br>All User Accounts Within the PBX", "header-user-account", extraButtonName, extraButtonURL)
				userAccountList(w, dbDetail, userTypeID)
				footer(w, "header-user-account", "button-user-account")
			} else if userTypeID == "301" {
				header(w, userPBXName+"<br>[PBX ID: "+userPBXID+"]<br>Own User Account for PBX", "header-user-account", extraButtonName, extraButtonURL)
				userAccountList(w, dbDetail, userTypeID)
				footer(w, "header-user-account", "button-user-account")
			} else if userTypeID == "302" {
				header(w, userPBXName+"<br>[PBX ID: "+userPBXID+"]<br>Own Read Only User Account for PBX", "header-user-account", extraButtonName, extraButtonURL)
				userAccountList(w, dbDetail, userTypeID)
				footer(w, "header-user-account", "button-user-account")
			} else if userTypeID == "400" {
				header(w, userCustomerName+"<br>[Customer ID: "+userCustomerID+"]<br>Own Invoice Customer Account", "header-user-account", extraButtonName, extraButtonURL)
				userAccountList(w, dbDetail, userTypeID)
				footer(w, "header-user-account", "button-user-account")
			} else {
				errorBox(w, "account_type_error", "header-user-account", "button-user-account")
			}
		}
		fmt.Fprintf(w, endHTML)
	})

	// Customer Page
	go http.HandleFunc("/customer", func(w http.ResponseWriter, r *http.Request) {

		// Open database connection
		dbConnection, err := sql.Open("mysql", dbUsername+":"+dbPassword+"@"+dbTransport+"("+dbAddress+":"+dbPort+")/"+dbName+"?tls="+dbTls)
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

		if userTypeID == "" {
			errorBox(w, "email_error", "header-customer", "button-customer")
		} else {
			if userTypeID == "100" {
				header(w, "YAP Admin Account<br>All Customers on YAP", "header-customer", extraButtonName, extraButtonURL)
				customerList(w, dbDetail, userTypeID, userCustomerID)
				fmt.Fprint(w, "<br>")
				customerAdd(w, dbDetail, r)
				footer(w, "header-customer", "button-customer")
			} else if userTypeID == "200" || userTypeID == "201" {
				header(w, userCustomerName+"<br>[Customer ID: "+userCustomerID+"]<br>Own Customer Information", "header-customer", extraButtonName, extraButtonURL)
				customerList(w, dbDetail, userTypeID, userCustomerID)
				footer(w, "header-customer", "button-customer")
			} else {
				errorBox(w, "account_type_error", "header-customer", "button-customer")
			}
		}
		fmt.Fprintf(w, endHTML)
	})

	// PBX Page

	go http.HandleFunc("/pbx", func(w http.ResponseWriter, r *http.Request) {

		// Open database connection
		dbConnection, err := sql.Open("mysql", dbUsername+":"+dbPassword+"@"+dbTransport+"("+dbAddress+":"+dbPort+")/"+dbName+"?tls="+dbTls)
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

		if userTypeID == "" {
			errorBox(w, "email_error", "header-pbx", "button-pbx")
		} else {
			if userTypeID == "100" {
				header(w, "YAP Admin Account<br>All BPXs on YAP", "header-pbx", extraButtonName, extraButtonURL)
				pbxList(w, dbDetail, userTypeID, userCustomerID, userPBXID)
				footer(w, "header-pbx", "button-pbx")
			} else if userTypeID == "200" || userTypeID == "201" {
				header(w, userCustomerName+"<br>[Customer ID: "+userCustomerID+"]<br>All PBXs for the Customer", "header-pbx", extraButtonName, extraButtonURL)
				pbxList(w, dbDetail, userTypeID, userCustomerID, userPBXID)
				footer(w, "header-pbx", "button-pbx")
			} else if userTypeID == "300" || userTypeID == "301" || userTypeID == "302" {
				header(w, userPBXName+"<br>[PBX ID: "+userPBXID+"]<br>PBX Information", "header-pbx", extraButtonName, extraButtonURL)
				pbxList(w, dbDetail, userTypeID, userCustomerID, userPBXID)
				footer(w, "header-pbx", "button-pbx")
			} else {
				errorBox(w, "account_type_error", "header-pbx", "button-pbx")
			}
		}
		fmt.Fprintf(w, endHTML)
	})

	// SIP Extension Page
	go http.HandleFunc("/sip-extension", func(w http.ResponseWriter, r *http.Request) {

		// Open database connection
		dbConnection, err := sql.Open("mysql", dbUsername+":"+dbPassword+"@"+dbTransport+"("+dbAddress+":"+dbPort+")/"+dbName+"?tls="+dbTls)
		defer dbConnection.Close()

		// Error
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, startHTML)

		// Wallpaper
		wallpaper(w, "wallpaper-sip-extension")

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

		if userTypeID == "" {
			errorBox(w, "email_error", "header-sip-extension", "button-sip-extension")
		} else {
			if userTypeID == "100" {
				header(w, "YAP Admin Account<br>All SIP Extensions on the Server", "header-sip-extension", extraButtonName, extraButtonURL)
				sipExtensionList(w, dbDetail, userTypeID, userCustomerID, userPBXID)
				footer(w, "header-sip-extension", "button-sip-extension")
			} else if userTypeID == "200" || userTypeID == "201" {
				header(w, userCustomerName+"<br>[Customer ID: "+userCustomerID+"]<br>All SIP Extensions for the Customer", "header-sip-extension", extraButtonName, extraButtonURL)
				sipExtensionList(w, dbDetail, userTypeID, userCustomerID, userPBXID)
				footer(w, "header-sip-extension", "button-sip-extension")
			} else if userTypeID == "300" || userTypeID == "301" || userTypeID == "302" {
				header(w, userPBXName+"<br>[PBX ID: "+userPBXID+"]<br>All SIP Extensions Within the PBX", "header-sip-extension", extraButtonName, extraButtonURL)
				sipExtensionList(w, dbDetail, userTypeID, userCustomerID, userPBXID)
				footer(w, "header-sip-extension", "button-sip-extension")
			} else {
				errorBox(w, "account_type_error", "header-sip-extension", "button-sip-extension")
			}
		}
		fmt.Fprintf(w, endHTML)
	})

	// Invoicing Page

	go http.HandleFunc("/invoice", func(w http.ResponseWriter, r *http.Request) {

		// Open database connection
		dbConnection, err := sql.Open("mysql", dbUsername+":"+dbPassword+"@"+dbTransport+"("+dbAddress+":"+dbPort+")/"+dbName+"?tls="+dbTls)
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

		if userTypeID == "" {
			errorBox(w, "email_error", "header-invoice", "button-invoice")
		} else {
			if userTypeID == "100" {
				header(w, "YAP Admin Account<br>All Customer Invoices", "header-invoice", extraButtonName, extraButtonURL)
				invoiceList(w, dbDetail, userTypeID, userCustomerID, yapAdminUKVATRegistered)
				footer(w, "header-invoice", "button-invoice")
			} else if userTypeID == "200" {
				header(w, userCustomerName+"<br>[Customer ID: "+userCustomerID+"]<br>Customer Invoice", "header-invoice", extraButtonName, extraButtonURL)
				invoiceList(w, dbDetail, userTypeID, userCustomerID, yapAdminUKVATRegistered)
				footer(w, "header-invoice", "button-invoice")
			} else if userTypeID == "400" {
				header(w, userCustomerName+"<br>[Customer ID: "+userCustomerID+"]<br>Customer Invoice", "header-invoice", extraButtonName, extraButtonURL)
				invoiceList(w, dbDetail, userTypeID, userCustomerID, yapAdminUKVATRegistered)
				footer(w, "header-invoice", "button-invoice")

			} else {
				errorBox(w, "account_type_error", "header-invoice", "button-invoice")
			}
		}

		//footer(w, "header-sip-trunk", "button-sip-trunk")
	})

	// Server Log Page
	http.HandleFunc("/server-log", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, startHTML)
		header(w, "Server Logs", "header-server-log", extraButtonName, extraButtonURL)

		// Wallpaper
		wallpaper(w, "wallpaper-server-log")

		footer(w, "header-server-log", "button-server-log")
		fmt.Fprintf(w, endHTML)
	})

	// Server Information Page
	http.HandleFunc("/server-information", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, startHTML)
		header(w, "Server Information", "header-server-information", extraButtonName, extraButtonURL)

		// Wallpaper
		wallpaper(w, "wallpaper-server-information")

		footer(w, "header-server-information", "button-server-information")
		fmt.Fprintf(w, endHTML)
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
		//Start server on port specified above
		log.Fatal(http.ListenAndServe(socket, nil))
	}
}

// Contributor(s):
// Elliot Michael Keavney
