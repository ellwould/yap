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

//----------------------------------------------------------------------------------------------------

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

// Embedded JavaScript and associated HTML functions

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
	if data == "type_id" {
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
	dbTotalTableCount.table = "customer"
	dbTotalTableCount.countMinusOne = true
	fmt.Fprintf(w, "    <td>"+totalTableCount(w, dbTotalTableCount)+"</td>")
	dbTotalTableCount.table = "pbx"
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
					                     user_account_first_name,
					                     user_account_last_name,
					                     user_account_email,
					                     user_account_type,
					                     user_account_date_added,
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
			userAccountFirstName      string
			userAccountLastName       string
			userAccountEmail          string
			userAccountType           string
			userAccountDateAdded      string
			userAccountTypePermission string
		)

		err = result.Scan(
			&userAccountFirstName,
			&userAccountLastName,
			&userAccountEmail,
			&userAccountType,
			&userAccountDateAdded,
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
		fmt.Fprintf(w, "          <th>Name</th>")
		fmt.Fprintf(w, "          <th>Email</th>")
		fmt.Fprintf(w, "          <th>Account Type</th>")
		fmt.Fprintf(w, "          <th>Account Created</th>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>"+userAccountFirstName+"<br>"+userAccountLastName+"</td>")
		fmt.Fprintf(w, "          <td>"+userAccountEmail+"</td>")
		fmt.Fprintf(w, "          <td>"+userAccountType+"</td>")
		fmt.Fprintf(w, "          <td>"+userAccountDateAdded+"</td>")
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

// User account page functions

func userAccountList(w http.ResponseWriter, dbDetail databaseFunctionParameter, userTypeID string) {

	var (
		userAccountFirstName string
		userAccountLastName  string
		userAccountEmail     string
		userAccountType      string
		userAccountDateAdded string
		customerID           string
		customerName         string
		pbxID                string
		pbxName              string
	)

	ownUserAccountSQL, err := dbDetail.connection.Query(`SELECT
							       user_account_first_name,
							       user_account_last_name,
							       user_account_email,
							       user_account_type,
							       user_account_date_added,
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
			&userAccountFirstName,
			&userAccountLastName,
			&userAccountEmail,
			&userAccountType,
			&userAccountDateAdded,
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
		fmt.Fprintf(w, "          <th>Name</th>")
		fmt.Fprintf(w, "          <th>Email</th>")
		fmt.Fprintf(w, "          <th>Account Type</th>")
		fmt.Fprintf(w, "          <th>Account Created</th>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>"+userAccountFirstName+"<br>"+userAccountLastName+"</td>")
		fmt.Fprintf(w, "          <td>"+userAccountEmail+"</td>")
		fmt.Fprintf(w, "          <td>"+userAccountType+"</td>")
		fmt.Fprintf(w, "          <td>"+userAccountDateAdded+"</td>")
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
		inputTableHTMLArgument.inputID = "other-account-input-date"
		inputTableHTMLArgument.funcNameJS = "otherAccountSearchDate"
		inputTableHTMLArgument.placeholder = "Date Created"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
			fmt.Fprintf(w, "    <br>")
			fmt.Fprintf(w, "    <br>")
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "other-account-input-pbx-name"
			inputTableHTMLArgument.funcNameJS = "otherAccountSearchPBXName"
			inputTableHTMLArgument.placeholder = "PBX Name"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "other-account-input-pbx-id"
			inputTableHTMLArgument.funcNameJS = "otherAccountSearchPBXID"
			inputTableHTMLArgument.placeholder = "PBX ID"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		}
		if userTypeID == "100" {
			inputTableHTMLArgument.inputID = "other-account-input-customer-name"
			inputTableHTMLArgument.funcNameJS = "otherAccountSearchCustomerName"
			inputTableHTMLArgument.placeholder = "Customer Name"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "other-account-input-customer-id"
			inputTableHTMLArgument.funcNameJS = "otherAccountSearchCustomerID"
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
		exportCSVButtonHTMLArgument.funcNameJS = "OtherAccount"
		exportCSVButtonHTMLArgument.buttonCSS = "button-user-account"
		exportCSVButtonHTML(w, exportCSVButtonHTMLArgument)
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table id=\"other-account-table\" class=\"table-user-account\">")
		fmt.Fprintf(w, "        <tr>")
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
						     			 user_account_first_name,
						     			 user_account_last_name,  
						     			 user_account_email,                                                   
						     			 user_account_type,  
						     			 user_account_date_added, 
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
				&userAccountFirstName,
				&userAccountLastName,
				&userAccountEmail,
				&userAccountType,
				&userAccountDateAdded,
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
			fmt.Fprintf(w, "          <td>"+userAccountFirstName+" "+userAccountLastName+"</td>")
			fmt.Fprintf(w, "          <td>"+userAccountEmail+"</td>")
			fmt.Fprintf(w, "          <td>"+userAccountType+"</td>")
			fmt.Fprintf(w, "          <td>"+userAccountDateAdded+"</td>")
			if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
				if pbxName != "system" {
					fmt.Fprintf(w, "          <td>"+pbxName+"</td>")
				} else {
					fmt.Fprintf(w, "          <td>-</td>")
				}
				if pbxID != "1" {
					fmt.Fprintf(w, "          <td>"+pbxID+"</td>")
				} else {
					fmt.Fprintf(w, "          <td>-</td>")
				}
			}
			if userTypeID == "100" {
				if customerName != "system" {
					fmt.Fprintf(w, "          <td>"+customerName+"</td>")
				} else {
					fmt.Fprintf(w, "          <td>-</td>")
				}
				if customerID != "1" {
					fmt.Fprintf(w, "          <td>"+customerID+"</td>")
				} else {
					fmt.Fprintf(w, "          <td>-</td>")
				}
			}
			fmt.Fprintf(w, "        </tr>")
		}

		fmt.Fprintf(w, "      </table>")
		var filterTableJSArgument jsFunctionParameter
		filterTableJSArgument.tableID = "other-account-table"
		// JS filter function for name in the other account table
		filterTableJSArgument.funcNameJS = "otherAccountSearchName"
		filterTableJSArgument.inputID = "other-account-input-name"
		filterTableJSArgument.columnNumber = 0
		filterTableJS(w, filterTableJSArgument)
		// JS filter function for email in the other account table
		filterTableJSArgument.funcNameJS = "otherAccountSearchEmail"
		filterTableJSArgument.inputID = "other-account-input-email"
		filterTableJSArgument.columnNumber = 1
		filterTableJS(w, filterTableJSArgument)
		// JS filter function for type in the other account table
		filterTableJSArgument.funcNameJS = "otherAccountSearchType"
		filterTableJSArgument.inputID = "other-account-input-type"
		filterTableJSArgument.columnNumber = 2
		filterTableJS(w, filterTableJSArgument)
		// JS filter function for date in the other account table
		filterTableJSArgument.funcNameJS = "otherAccountSearchDate"
		filterTableJSArgument.inputID = "other-account-input-date"
		filterTableJSArgument.columnNumber = 3
		filterTableJS(w, filterTableJSArgument)
		if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
			// JS filter function for PBX name in the other account table
			filterTableJSArgument.funcNameJS = "otherAccountSearchPBXName"
			filterTableJSArgument.inputID = "other-account-input-pbx-name"
			filterTableJSArgument.columnNumber = 4
			filterTableJS(w, filterTableJSArgument)
			// JS filter function for PBX ID in the other account table
			filterTableJSArgument.funcNameJS = "otherAccountSearchPBXID"
			filterTableJSArgument.inputID = "other-account-input-pbx-id"
			filterTableJSArgument.columnNumber = 5
			filterTableJS(w, filterTableJSArgument)
		}
		if userTypeID == "100" {
			// JS filter function for the customer name in the other account table
			filterTableJSArgument.funcNameJS = "otherAccountSearchCustomerName"
			filterTableJSArgument.inputID = "other-account-input-customer-name"
			filterTableJSArgument.columnNumber = 6
			filterTableJS(w, filterTableJSArgument)
			// JS filter function for the customer ID in the other account table
			filterTableJSArgument.funcNameJS = "otherAccountSearchCustomerID"
			filterTableJSArgument.inputID = "other-account-input-customer-id"
			filterTableJSArgument.columnNumber = 7
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

// Function to add HTML input fields
func inputHTML(w http.ResponseWriter, inputValue string, labelMessage string, inputType string) {
	fmt.Fprintf(w, "  <label for=\""+inputValue+"\"><b>Enter "+labelMessage+"</b>")
	fmt.Fprintf(w, "  </label><br>")
	fmt.Fprintf(w, "  <input type=\""+inputType+"\" id=\""+inputValue+"\" name=\""+inputValue+"\">")
	fmt.Fprintf(w, "<br>")
}

func selectHTML(w http.ResponseWriter, selectValue string, labelMessage string, optionValue []string) {
	fmt.Fprintf(w, "  <label for=\""+selectValue+"\"><b>Select "+labelMessage+"</b>")
	fmt.Fprintf(w, "  </label><br>")
	fmt.Fprintf(w, "  <select id=\""+selectValue+"\" name=\""+selectValue+"\">")
	for i := 0; i < len(optionValue); i++ {
		fmt.Fprintf(w, "<option value="+optionValue[i]+">"+optionValue[i]+"</option>")
	}
	fmt.Fprintf(w, "  </select>")
}

func userAccountAdd(w http.ResponseWriter) {

	fmt.Fprintf(w, "<table id=\"table\" class=\"table-user-account\">")
	fmt.Fprintf(w, "  <form method=\"POST\" action=\"/user-account\">")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th class=\"table-title\";>Add New User:<br>(May need to update /etc/oauth2-proxy/email.txt)</th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>")
	fmt.Fprintf(w, "      <table id=\"table\" class=\"table-user-account\" style=\"border-style:hidden\">")
	fmt.Fprintf(w, "        <tr>")
	fmt.Fprintf(w, "          <td>")
	inputHTML(w, "firstName", "First Name:<br>(Cannot Be Blank)", "text")
	fmt.Fprintf(w, "	  </td>")
	fmt.Fprintf(w, "          <td>")
	inputHTML(w, "lastName", "Last Name:<br>(Cannot Be Blank)", "text")
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "          <td>")
	inputHTML(w, "email", "Email Address:<br>(Must Be a Valid Email Address)", "text")
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "        </tr>")
	fmt.Fprintf(w, "        <tr>")
	fmt.Fprintf(w, "          <td>")
	var accountTypeList = []string{"100 - YAP Admin", "200 - Group Admin", "201 - Group Regular", "300 - PBX Admin", "301 - PBX Regular", "302 - PBX Read Only"}
	selectHTML(w, "accountType", "Account Type:<br>(100, 200, 201, 300, 301, 302)", accountTypeList)
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "          <td>")
	inputHTML(w, "pbxID", "PBX ID:<br>(Blank if Account Type is 100, 200 or 201)", "text")
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "          <td>")
	inputHTML(w, "groupID", "Group ID:<br>(Blank if Account Type is 100)", "text")
	fmt.Fprintf(w, "          </td>")
	fmt.Fprintf(w, "        </tr>")
	fmt.Fprintf(w, "      </table>")
	fmt.Fprintf(w, "    </th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th><input type=\"submit\" value=\"Create User Account\"></th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "  </form>")
	fmt.Fprintf(w, "</table>")

}

func userAccountEdit() {

}

func userAccountDelete() {

}

//----------------------------------------------------------------------------------------------------

// Customer page functions

func customerList(w http.ResponseWriter, dbDetail databaseFunctionParameter, userTypeID string, userCustomerID string) {

	var (
		customerName                     string
		customerID                       string
		customerDateAdded                string
		customerActive                   string
		ukBased                          string
		consumerType                     string
		ukVATStatus                      string
		resellingMiniutes                string
		pbxLimit                         string
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
	dbTableCountUserCustomer.table = "customer"
	dbTableCountUserCustomer.columnWhere = "active"

	if userTypeID == "100" {
		fmt.Fprintf(w, "<table id=\"table\" class=\"table-customer\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table id=\"table\" class=\"table-customer\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <th>Total Customers On YAP</th>")
		fmt.Fprintf(w, "          <th>Total Active Customers</th>")
		fmt.Fprintf(w, "          <th>Total Inactive Customers</th>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		dbTableCountUserCustomer.countMinusOne = true
		fmt.Fprintf(w, "    <td>"+totalTableCount(w, dbTableCountUserCustomer)+"</td>")
		dbTableCountUserCustomer.columnWhere = "active"
		dbTableCountUserCustomer.countMinusOne = false
		dbTableCountUserCustomer.columnWhereValue = "1"
		fmt.Fprintf(w, "    <td>"+totalTableCountWhere(w, dbTableCountUserCustomer)+"</td>")
		dbTableCountUserCustomer.countMinusOne = true
		dbTableCountUserCustomer.columnWhereValue = "0"
		fmt.Fprintf(w, "    <td>"+totalTableCountWhere(w, dbTableCountUserCustomer)+"</td>")
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
		inputTableHTMLArgument.inputID = "customer-contact-input-customer-name"
		inputTableHTMLArgument.funcNameJS = "customerContactSearchCustomerName"
		inputTableHTMLArgument.placeholder = "Customer Name"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "customer-contact-input-customer-id"
		inputTableHTMLArgument.funcNameJS = "customerContactSearchCustomerID"
		inputTableHTMLArgument.placeholder = "Customer ID"
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
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    <br>")
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
	fmt.Fprintf(w, "          <th>Customer<br>Name</th>")
	fmt.Fprintf(w, "          <th>Customer<br>ID</th>")
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
		)

		// Error
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>"+customerName+"</td>")
		fmt.Fprintf(w, "          <td>"+customerID+"</td>")
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
		// JS filter function for the customer name in the customer contact table
		filterTableJSArgument.funcNameJS = "customerContactSearchCustomerName"
		filterTableJSArgument.inputID = "customer-contact-input-customer-name"
		filterTableJSArgument.columnNumber = 0
		filterTableJS(w, filterTableJSArgument)
		// JS filter function for the customer ID in the customer contact table
		filterTableJSArgument.funcNameJS = "customerContactSearchCustomerID"
		filterTableJSArgument.inputID = "customer-contact-input-customer-id"
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
		inputTableHTMLArgument.inputID = "customer-resource-input-customer-name"
		inputTableHTMLArgument.funcNameJS = "customerResourceSearchCustomerName"
		inputTableHTMLArgument.placeholder = "Customer Name"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "customer-resource-input-customer-id"
		inputTableHTMLArgument.funcNameJS = "customerResourceSearchCustomerID"
		inputTableHTMLArgument.placeholder = "Customer ID"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "customer-resource-input-date"
		inputTableHTMLArgument.funcNameJS = "customerResourceSearchDate"
		inputTableHTMLArgument.placeholder = "Date Created"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "customer-resource-input-active"
		inputTableHTMLArgument.funcNameJS = "customerResourceSearchActive"
		inputTableHTMLArgument.placeholder = "Customer Active Status"
		inputTableHTML(w, inputTableHTMLArgument)
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
	fmt.Fprintf(w, "          <th>Customer<br>Name</th>")
	fmt.Fprintf(w, "          <th>Customer<br>ID</th>")
	fmt.Fprintf(w, "          <th>Date Customer<br>Added</th>")
	fmt.Fprintf(w, "          <th>Customer<br>Active</th>")
	fmt.Fprintf(w, "          <th>UK Based</th>")
	fmt.Fprintf(w, "          <th>Consumer<br>Type</th>")
	fmt.Fprintf(w, "          <th>UK VAT Status</th>")
	fmt.Fprintf(w, "          <th>Reselling<br>Miniutes</th>")
	fmt.Fprintf(w, "          <th>PBX Limit</th>")
	fmt.Fprintf(w, "        </tr>")

	customerResourceSQL, err := dbDetail.connection.Query(`SELECT
							customer_name,
							customer_id,
							customer_date_added,
							customer_active,
							customer_uk_based,
							customer_consumer_type,
							customer_uk_vat_status,
							customer_reselling_miniutes,
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
			&customerName,
			&customerID,
			&customerDateAdded,
			&customerActive,
			&ukBased,
			&consumerType,
			&ukVATStatus,
			&resellingMiniutes,
			&pbxLimit,
		)

		// Error
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>"+customerName+"</td>")
		fmt.Fprintf(w, "          <td>"+customerID+"</td>")
		fmt.Fprintf(w, "          <td>"+customerDateAdded+"</td>")
		if customerActive == "1" {
			fmt.Fprintf(w, "          <td>Yes</td>")
		} else {
			fmt.Fprintf(w, "	  <td>No</td>")
		}
		fmt.Fprintf(w, "          <td>"+ukBased+"</td>")
		fmt.Fprintf(w, "          <td>"+consumerType+"</td>")
		fmt.Fprintf(w, "          <td>"+ukVATStatus+"</td>")
		fmt.Fprintf(w, "          <td>"+resellingMiniutes+"</td>")
		fmt.Fprintf(w, "          <td>"+pbxLimit+"</td>")
		fmt.Fprintf(w, "        </tr>")
	}

	fmt.Fprintf(w, "      </table>")
	if userTypeID == "100" {
		var filterTableJSArgument jsFunctionParameter
		filterTableJSArgument.tableID = "customer-resource-table"
		// Call JS filter function for the customer name in the customer resource table
		filterTableJSArgument.funcNameJS = "customerResourceSearchCustomerName"
		filterTableJSArgument.inputID = "customer-resource-input-customer-name"
		filterTableJSArgument.columnNumber = 0
		filterTableJS(w, filterTableJSArgument)
		// Call JS filter function for the customer ID in the customer resource table
		filterTableJSArgument.funcNameJS = "customerResourceSearchCustomerID"
		filterTableJSArgument.inputID = "customer-resource-input-customer-id"
		filterTableJSArgument.columnNumber = 1
		filterTableJS(w, filterTableJSArgument)
		// Call JS filter function for date in the customer resource table
		filterTableJSArgument.funcNameJS = "customerResourceSearchDate"
		filterTableJSArgument.inputID = "customer-resource-input-date"
		filterTableJSArgument.columnNumber = 2
		filterTableJS(w, filterTableJSArgument)
		// Call JS filter function for active status in the customer resource table
		filterTableJSArgument.funcNameJS = "customerResourceSearchActive"
		filterTableJSArgument.inputID = "customer-resource-input-active"
		filterTableJSArgument.columnNumber = 3
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

//----------------------------------------------------------------------------------------------------

// PBX page functions

func pbxList(w http.ResponseWriter, dbDetail databaseFunctionParameter, userTypeID string, userCustomerID string, userPBXID string) {

	var (
		pbxName                     string
		pbxID                       string
		pbxDateAdded                string
		pbxActive                   string
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
		customerName                string
		customerID                  string
	)

	var dbTableCountUserPBX databaseFunctionParameter
	dbTableCountUserPBX.connection = dbDetail.connection
	dbTableCountUserPBX.database = dbDetail.database
	dbTableCountUserPBX.table = "pbx"

	if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
		fmt.Fprintf(w, "<table id=\"table\" class=\"table-pbx\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table id=\"table\" class=\"table-pbx\">")
		fmt.Fprintf(w, "        <tr>")
		if userTypeID == "100" {
			fmt.Fprintf(w, "          <th>Total PBXs On YAP</th>")
		}
		fmt.Fprintf(w, "          <th>Total Active PBXs</th>")
		fmt.Fprintf(w, "          <th>Total Inactive PBXs</th>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		if userTypeID == "100" {
			dbTableCountUserPBX.countMinusOne = true
			fmt.Fprintf(w, "    <td>"+totalTableCount(w, dbTableCountUserPBX)+"</td>")
			dbTableCountUserPBX.columnWhere = "active"
			dbTableCountUserPBX.countMinusOne = false
			dbTableCountUserPBX.columnWhereValue = "1"
			fmt.Fprintf(w, "    <td>"+totalTableCountWhere(w, dbTableCountUserPBX)+"</td>")
			dbTableCountUserPBX.countMinusOne = true
			dbTableCountUserPBX.columnWhereValue = "0"
			fmt.Fprintf(w, "    <td>"+totalTableCountWhere(w, dbTableCountUserPBX)+"</td>")
		} else if userTypeID == "200" || userTypeID == "201" {
			dbTableCountUserPBX.columnWhere = "customer_id"
			dbTableCountUserPBX.columnWhereValue = userCustomerID
			dbTableCountUserPBX.columnWhereAnd = "active"
			dbTableCountUserPBX.columnWhereValueAnd = "1"
			fmt.Fprintf(w, "    <td>"+totalTableCountWhereAnd(w, dbTableCountUserPBX)+"</td>")
			dbTableCountUserPBX.columnWhereValueAnd = "0"
			fmt.Fprintf(w, "    <td>"+totalTableCountWhereAnd(w, dbTableCountUserPBX)+"</td>")
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
		inputTableHTMLArgument.inputID = "pbx-contact-input-pbx-name"
		inputTableHTMLArgument.funcNameJS = "pbxContactSearchPBXName"
		inputTableHTMLArgument.placeholder = "PBX Name"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "pbx-contact-input-pbx-id"
		inputTableHTMLArgument.funcNameJS = "pbxContactSearchPBXID"
		inputTableHTMLArgument.placeholder = "PBX ID"
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
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    <br>")
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
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    <br>")
		if userTypeID == "100" {
			inputTableHTMLArgument.inputID = "pbx-contact-input-customer-name"
			inputTableHTMLArgument.funcNameJS = "pbxContactSearchCustomerName"
			inputTableHTMLArgument.placeholder = "Customer Name"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "pbx-contact-input-customer-id"
			inputTableHTMLArgument.funcNameJS = "pbxContactSearchCustomerID"
			inputTableHTMLArgument.placeholder = "Customer ID"
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
	if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
		fmt.Fprintf(w, "          <th>PBX Name</th>")
		fmt.Fprintf(w, "          <th>PBX ID</th>")
	}
	fmt.Fprintf(w, "          <th>PBX Site<br> Address</th>")
	fmt.Fprintf(w, "          <th>PBX Site<br> Email Address</th>")
	fmt.Fprintf(w, "          <th>PBX Site<br> Phone Number</th>")
	fmt.Fprintf(w, "          <th>PBX Invoice<br> Address</th>")
	fmt.Fprintf(w, "          <th>PBX Invoice<br> Email Address</th>")
	fmt.Fprintf(w, "          <th>PBX Invoice<br> Phone Number</th>")
	if userTypeID == "100" {
		fmt.Fprintf(w, "          <th>Customer Name</th>")
		fmt.Fprintf(w, "          <th>Customer ID</th>")
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
					                pbx_invoice_contact_number,
					                customer_name,
					                customer_id
					              FROM
					  	        yap.view___pbx_detail
						      `+whereClause, userWhereID)

	// Error
	if err != nil {
		panic(err)

	}

	for pbxContactSQL.Next() {

		err = pbxContactSQL.Scan(
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
			&customerName,
			&customerID,
		)

		// Error
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "        <tr>")
		if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
			fmt.Fprintf(w, "          <td>"+pbxName+"</td>")
			fmt.Fprintf(w, "          <td>"+pbxID+"</td>")
		}
		fmt.Fprintf(w, "          <td style=\"text-align: left;\">"+pbxSiteAddressLine1+"&nbsp<br>"+pbxSiteAddressLine2+"&nbsp<br>"+pbxSiteCityTownVillage+"&nbsp<br>"+pbxSiteCountyStateRegion+"&nbsp<br><br>"+pbxSitePostcodeZipCode+"&nbsp<br><br>"+pbxSiteCountry+"&nbsp</td>")
		fmt.Fprintf(w, "          <td><a href=\"mailto:"+pbxSiteContactEmail+"\">"+pbxSiteContactEmail+"</a></td>")
		fmt.Fprintf(w, "          <td><a href=\"tel:"+pbxSiteContactNumber+"\">"+pbxSiteContactNumber+"</a></td>")
		fmt.Fprintf(w, "          <td style=\"text-align: left;\">"+pbxInvoiceAddressLine1+"&nbsp<br>"+pbxInvoiceAddressLine2+"&nbsp<br>"+pbxInvoiceCityTownVillage+"&nbsp<br>"+pbxInvoiceCountyStateRegion+"&nbsp<br><br>"+pbxInvoicePostcodeZipCode+"&nbsp<br><br>"+pbxInvoiceCountry+"&nbsp</td>")
		fmt.Fprintf(w, "          <td>&nbsp<a href=\"mailto:"+pbxInvoiceContactEmail+"\">"+pbxInvoiceContactEmail+"</a></td>")
		fmt.Fprintf(w, "          <td><a href=\"tel:"+pbxInvoiceContactNumber+"\">"+pbxInvoiceContactNumber+"</a></td>")
		if userTypeID == "100" {
			fmt.Fprintf(w, "          <td>"+customerName+"</td>")
			fmt.Fprintf(w, "          <td>"+customerID+"</td>")
		}
		fmt.Fprintf(w, "        </tr>")
	}

	fmt.Fprintf(w, "      </table>")
	if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
		var filterTableJSArgument jsFunctionParameter
		filterTableJSArgument.tableID = "pbx-contact-table"
		// Call JS filter function for PBX name in the PBX contact table
		filterTableJSArgument.funcNameJS = "pbxContactSearchPBXName"
		filterTableJSArgument.inputID = "pbx-contact-input-pbx-name"
		filterTableJSArgument.columnNumber = 0
		filterTableJS(w, filterTableJSArgument)
		// Call JS filter function for PBX ID in the PBX contact table
		filterTableJSArgument.funcNameJS = "pbxContactSearchPBXID"
		filterTableJSArgument.inputID = "pbx-contact-input-pbx-id"
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
			// Call JS filter function for the customer name in the PBX contact table
			filterTableJSArgument.funcNameJS = "pbxContactSearchCustomerName"
			filterTableJSArgument.inputID = "pbx-contact-input-customer-name"
			filterTableJSArgument.columnNumber = 8
			filterTableJS(w, filterTableJSArgument)
			// Call JS filter function for the customer ID in the PBX contact table
			filterTableJSArgument.funcNameJS = "pbxContactSearchCustomerID"
			filterTableJSArgument.inputID = "pbx-contact-input-customer-id"
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
		inputTableHTMLArgument.inputID = "pbx-resource-input-pbx-name"
		inputTableHTMLArgument.funcNameJS = "pbxResourceSearchPBXName"
		inputTableHTMLArgument.placeholder = "PBX Name"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "pbx-resource-input-pbx-id"
		inputTableHTMLArgument.funcNameJS = "pbxResourceSearchPBXID"
		inputTableHTMLArgument.placeholder = "PBX ID"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "pbx-resource-input-date"
		inputTableHTMLArgument.funcNameJS = "pbxResourceSearchDate"
		inputTableHTMLArgument.placeholder = "Date Created"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "pbx-resource-input-active"
		inputTableHTMLArgument.funcNameJS = "pbxResourceSearchActive"
		inputTableHTMLArgument.placeholder = "PBX Active Status"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    <br>")
		if userTypeID == "100" {
			inputTableHTMLArgument.inputID = "pbx-resource-input-customer-name"
			inputTableHTMLArgument.funcNameJS = "pbxResourceSearchCustomerName"
			inputTableHTMLArgument.placeholder = "Customer Name"
			inputTableHTML(w, inputTableHTMLArgument)
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTMLArgument.inputID = "pbx-resource-input-customer-id"
			inputTableHTMLArgument.funcNameJS = "pbxResourceSearchCustomerID"
			inputTableHTMLArgument.placeholder = "Customer ID"
			inputTableHTML(w, inputTableHTMLArgument)
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
		fmt.Fprintf(w, "          <th>PBX Name</th>")
		fmt.Fprintf(w, "          <th>PBX ID</th>")
	}
	fmt.Fprintf(w, "          <th>PBX Date</th>")
	fmt.Fprintf(w, "          <th>PBX Active<br>Status</th>")
	fmt.Fprintf(w, "          <th>SIP Extension<br>Limit for PBX</th>")
	if userTypeID == "100" {
		fmt.Fprintf(w, "          <th>Customer Name</th>")
		fmt.Fprintf(w, "          <th>Customer ID</th>")
	}
	fmt.Fprintf(w, "        </tr>")

	pbxResourceSQL, err := dbDetail.connection.Query(`SELECT
							pbx_name,
							pbx_id,
							pbx_date_added,
							pbx_active,
							pbx_sip_extension_limit,
							customer_name,
							customer_id
					              FROM
					  	        yap.view___pbx_detail
						      `+whereClause, userWhereID)

	// Error
	if err != nil {
		panic(err)

	}

	for pbxResourceSQL.Next() {

		err = pbxResourceSQL.Scan(
			&pbxName,
			&pbxID,
			&pbxDateAdded,
			&pbxActive,
			&pbxSIPExtensionLimit,
			&customerName,
			&customerID,
		)

		// Error
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "        <tr>")
		if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
			fmt.Fprintf(w, "          <td>"+pbxName+"</td>")
			fmt.Fprintf(w, "          <td>"+pbxID+"</td>")
		}
		fmt.Fprintf(w, "          <td>"+pbxDateAdded+"</td>")
		if pbxActive == "1" {
			fmt.Fprintf(w, "          <td>YES</td>")
		} else {
			fmt.Fprintf(w, "          <td>NO</td>")
		}
		fmt.Fprintf(w, "          <td>"+pbxSIPExtensionLimit+"</td>")
		if userTypeID == "100" {
			fmt.Fprintf(w, "          <td>"+customerName+"</td>")
			fmt.Fprintf(w, "          <td>"+customerID+"</td>")
		}
		fmt.Fprintf(w, "        </tr>")
	}

	fmt.Fprintf(w, "      </table>")
	if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
		var filterTableJSArgument jsFunctionParameter
		filterTableJSArgument.tableID = "pbx-resource-table"
		// Call JS filter function for PBX name in the PBX resource table
		filterTableJSArgument.funcNameJS = "pbxResourceSearchPBXName"
		filterTableJSArgument.inputID = "pbx-resource-input-pbx-name"
		filterTableJSArgument.columnNumber = 0
		filterTableJS(w, filterTableJSArgument)
		// Call JS filter function for PBX ID in the PBX resource table
		filterTableJSArgument.funcNameJS = "pbxResourceSearchPBXID"
		filterTableJSArgument.inputID = "pbx-resource-input-pbx-id"
		filterTableJSArgument.columnNumber = 1
		filterTableJS(w, filterTableJSArgument)
		// Call JS filter function for date in the PBX resource table
		filterTableJSArgument.funcNameJS = "pbxResourceSearchDate"
		filterTableJSArgument.inputID = "pbx-resource-input-date"
		filterTableJSArgument.columnNumber = 2
		filterTableJS(w, filterTableJSArgument)
		// Call JS filter function for active status in the PBX resource table
		filterTableJSArgument.funcNameJS = "pbxResourceSearchActive"
		filterTableJSArgument.inputID = "pbx-resource-input-active"
		filterTableJSArgument.columnNumber = 3
		filterTableJS(w, filterTableJSArgument)
		if userTypeID == "100" {
			// Call JS filter function for the customer name in the PBX resource table
			filterTableJSArgument.funcNameJS = "pbxResourceSearchCustomerName"
			filterTableJSArgument.inputID = "pbx-resource-input-customer-name"
			filterTableJSArgument.columnNumber = 10
			filterTableJS(w, filterTableJSArgument)
			// Call JS filter function for the customer ID in the PBX resource table
			filterTableJSArgument.funcNameJS = "pbxResourceSearchCustomerID"
			filterTableJSArgument.inputID = "pbx-resource-input-customer-id"
			filterTableJSArgument.columnNumber = 11
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
		pbxName                string
		pbxID                  string
		customerName           string
		customerID             string
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
		inputTableHTMLArgument.inputID = "sip-extension-detail-input-customer-name"
		inputTableHTMLArgument.funcNameJS = "sipExtensionDetailSearchCustomerName"
		inputTableHTMLArgument.placeholder = "Customer Name"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "sip-extension-detail-input-customer-id"
		inputTableHTMLArgument.funcNameJS = "sipExtensionDetailSearchCustomerID"
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
		fmt.Fprintf(w, "          <th>Customer Name</th>")
		fmt.Fprintf(w, "          <th>Customer ID</th>")
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
					                pbx_name,
					                pbx_id,
					                customer_name,
					                customer_id
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
			&pbxName,
			&pbxID,
			&customerName,
			&customerID,
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
			fmt.Fprintf(w, "          <td>"+customerName+"</td>")
			fmt.Fprintf(w, "          <td>"+customerID+"</td>")
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
		// Call JS filter function for the customer name in the SIP extension detail table
		filterTableJSArgument.funcNameJS = "sipExtensionDetailSearchCustomerName"
		filterTableJSArgument.inputID = "sip-extension-detail-input-customer-name"
		filterTableJSArgument.columnNumber = 5
		filterTableJS(w, filterTableJSArgument)
		// Call JS filter function for the customer ID in the SIP extension detail table
		filterTableJSArgument.funcNameJS = "sipExtensionDetailSearchCustomerID"
		filterTableJSArgument.inputID = "sip-extension-detail-input-customer-id"
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
		inputTableHTMLArgument.inputID = "sip-extension-reg-input-customer-name"
		inputTableHTMLArgument.funcNameJS = "sipExtensionRegSearchCustomerName"
		inputTableHTMLArgument.placeholder = "Customer Name"
		inputTableHTML(w, inputTableHTMLArgument)
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTMLArgument.inputID = "sip-extension-reg-input-customer-id"
		inputTableHTMLArgument.funcNameJS = "sipExtensionRegSearchCustomerID"
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
	fmt.Fprintf(w, "      <table id=\"sip-extension-reg-table\" class=\"table-sip-extension\">")
	fmt.Fprintf(w, "        <tr>")
	fmt.Fprintf(w, "          <th>SIP Username</th>")
	fmt.Fprintf(w, "          <th>URI</th>")
	fmt.Fprintf(w, "          <th>User Agent<br>SIP Client</th>")
	if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
		fmt.Fprintf(w, "          <th>PBX Name</th>")
	}
	if userTypeID == "100" {
		fmt.Fprintf(w, "          <th>Customer Name</th>")
		fmt.Fprintf(w, "          <th>Customer ID</th>")
	}
	fmt.Fprintf(w, "        </tr>")

	sipExtensionRegSQL, err := dbDetail.connection.Query(`SELECT
							sip_username,
							uri,
							user_agent,
					                pbx_name,
					                customer_name,
					                customer_id
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
			&customerName,
			&customerID,
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
		fmt.Fprintf(w, "          <td>"+customerName+"</td>")
		fmt.Fprintf(w, "          <td>"+customerID+"</td>")
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
		// Call JS filter function for the customer name in the SIP extension registration (reg) table
		filterTableJSArgument.funcNameJS = "sipExtensionRegSearchCustomerName"
		filterTableJSArgument.inputID = "sip-extension-reg-input-customer-name"
		filterTableJSArgument.columnNumber = 4
		filterTableJS(w, filterTableJSArgument)
		// Call JS filter function for the customer ID in the SIP extension registration (reg) table
		filterTableJSArgument.funcNameJS = "sipExtensionRegSearchCustomerID"
		filterTableJSArgument.inputID = "sip-extension-reg-input-customer-id"
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
				mainMenuButtonFive := mainMenuParameter{writeHTTP: w, buttonName: "All Customer<br>Invoicing<br>&#129534", hyperlink: "/invoicing", headerCSS: "header-invoicing", buttonCSS: "button-invoicing"}
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
				mainMenuButtonFive := mainMenuParameter{writeHTTP: w, buttonName: "Customer<br>Invoice<br>&#129534", hyperlink: "/invoicing", headerCSS: "header-invoicing", buttonCSS: "button-invoicing"}
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
				mainMenuButtonFive := mainMenuParameter{writeHTTP: w, buttonName: "Customer<br>Invoice<br>&#129534", hyperlink: "/invoicing", headerCSS: "header-invoicing", buttonCSS: "button-invoicing"}
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
				mainMenuButtonOne := mainMenuParameter{writeHTTP: w, buttonName: "Customer<br>Invoice<br>&#129534", hyperlink: "/invoicing", headerCSS: "header-invoicing", buttonCSS: "button-invoicing"}
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
	go http.HandleFunc("/user-account", func(w http.ResponseWriter, r *http.Request) {

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
				userAccountAdd(w)
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

	go http.HandleFunc("/invoicing", func(w http.ResponseWriter, r *http.Request) {

		// Open database connection
		dbConnection, err := sql.Open("mysql", dbUsername+":"+dbPassword+"@"+dbTransport+"("+dbAddress+":"+dbPort+")/"+dbName+"?tls="+dbTls)
		defer dbConnection.Close()

		// Error
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, startHTML)

		// Wallpaper
		wallpaper(w, "wallpaper-invoicing")

		// Code to call the emailHeaderHTTP function
		email := emailHeaderHTTP(r)

		var dbDetail databaseFunctionParameter
		dbDetail.connection = dbConnection
		dbDetail.database = dbName
		dbDetail.columnWhereValue = email

		userTypeID := userAccountData(dbDetail, "type_id")

		if userTypeID == "" {
			errorBox(w, "email_error", "header-sip-trunk", "button-sip-trunk")
		} else {
			if userTypeID == "100" {
				header(w, "YAP Admin Account<br>Customer Invoicing", "header-invoicing", extraButtonName, extraButtonURL)

				footer(w, "header-invoicing", "button-invoicing")
			} else {
				errorBox(w, "account_type_error", "header-invoicing", "button-invoicing")
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
