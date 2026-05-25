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
func errorBox(w http.ResponseWriter, errorType string, headerCSS string) {
	fmt.Fprintf(w, "<div class=\"div-error-box\">")
	fmt.Fprintf(w, "  <h1 class=\""+headerCSS+"\">")
	if errorType == "email_error" {
		fmt.Fprintf(w, "    User Account Not Found<br>")
		fmt.Fprintf(w, "    <a href=\"/oauth2/sign_out?rd=https://github.com/logout\" class=\"button-general button-header button-logout\">Logout</a>")
	} else if errorType == "account_type_error" {
		fmt.Fprintf(w, "    Account Type Forbidden<br>")
		fmt.Fprintf(w, "    <a href=\"/main-menu\" class=\"button-general button-header\">Main Menu</a>")
	} else {
		fmt.Fprintf(w, "    Unknown Error<br>")
		fmt.Fprintf(w, "    <a href=\"/main-menu\" class=\"button-general button-header\">Main Menu</a>")
	}
	fmt.Fprintf(w, "</h1>")
	fmt.Fprintf(w, "</div>")
}

// Function for the header
func header(w http.ResponseWriter, headerTitle string, headerCSS string) {
	fmt.Fprintf(w, "<div class=\"div-header\">")
	fmt.Fprintf(w, "  <h1 class=\""+headerCSS+"\">")
	fmt.Fprintf(w, "    ⊛ YAP [Yet Another PBX] ⊛")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "    <a href=\"/oauth2/sign_out?rd=https://github.com/logout\" class=\"button-general button-header button-logout\">Logout</a>")
	fmt.Fprintf(w, "    <a href=\"/\" class=\"button-general button-header button-home\">Home</a>")
	fmt.Fprintf(w, "    <a href=\"https://github.com/ellwould/yap/blob/main/LICENSE\" target=\"_blank\" class=\"button-general button-header button-license\">License</a>")
	fmt.Fprintf(w, "    <br>")
	fmt.Fprintf(w, "    "+headerTitle)
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

// Embedded JavaScript and associated HTML functions

// JavaScript toggle function
func toggleDivJS(w http.ResponseWriter, functionName string, divID string) {
	fmt.Fprintf(w, "<script>")
	fmt.Fprintf(w, "  function "+functionName+"() {")
	fmt.Fprintf(w, "  var x = document.getElementById(\""+divID+"\");")
	fmt.Fprintf(w, "  if (x.style.display === \"none\") {")
	fmt.Fprintf(w, "    x.style.display = \"table\";")
	fmt.Fprintf(w, "  } else {")
	fmt.Fprintf(w, "    x.style.display = \"none\";")
	fmt.Fprintf(w, "  }")
	fmt.Fprintf(w, "}")
	fmt.Fprintf(w, "</script>")
}

// Input for filtering a HTML table
func inputTableHTML(w http.ResponseWriter, functionName string, inputID string, placeholder string) {
	fmt.Fprintf(w, "<input type=\"text\" id=\""+inputID+"\" onkeyup=\""+functionName+"()\" placeholder=\"Filter Via "+placeholder+"...\" title=\""+placeholder+"\">")
}

// JavaScript filter HTML table function
func filterTableJS(w http.ResponseWriter, functionName string, inputID string, tableID string, columnNumber int) {
	if columnNumber > 11 {
		panic("Table column number cannot exceed 11")
	} else if columnNumber < 0 {
		panic("Table column number cannot be a negative number")
	} else {
		fmt.Fprintf(w, "<script>")
		fmt.Fprintf(w, "function "+functionName+"() {")
		fmt.Fprintf(w, "  var input, filter, table, tr, td, i, txtValue;")
		fmt.Fprintf(w, "  input = document.getElementById(\""+inputID+"\");")
		fmt.Fprintf(w, "  filter = input.value.toUpperCase();")
		fmt.Fprintf(w, "  table = document.getElementById(\""+tableID+"\");")
		fmt.Fprintf(w, "  tr = table.getElementsByTagName(\"tr\");")
		fmt.Fprintf(w, "  for (i = 0; i < tr.length; i++) {")
		fmt.Fprintf(w, "    td = tr[i].getElementsByTagName(\"td\")["+strconv.Itoa(columnNumber)+"];")
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
func copyButtonJS(w http.ResponseWriter, data string) {
	jsFuncName := strings.Replace(data, "-", "", -1)
	fmt.Fprintf(w, "<div class=\"button-data-space\"></div><button class=\"button-data\" onclick=cp"+jsFuncName+"()>Copy &#10697</button><br>")
	fmt.Fprintf(w, "<script>")
	fmt.Fprintf(w, "  function cp"+jsFuncName+"() {")
	fmt.Fprintf(w, "    navigator.clipboard.writeText('"+data+"');")
	fmt.Fprintf(w, "  }")
	fmt.Fprintf(w, "</script>")
}

// HTML button to call the JavaScript exportCSVJS function
func exportCSVButtonHTML(w http.ResponseWriter, jsFuncName string, buttonCSS string) {
	fmt.Fprintf(w, "<button class=\"button-general "+buttonCSS+"\" onclick=\"exportTable"+jsFuncName+"ToCSV()\">Export to CSV</button><br>")
}

// JavaScript to download or view a HTML table as a CSV file
func exportCSVJS(w http.ResponseWriter, jsFuncName string, tableID string, fileName string, path string) {
	fmt.Fprintf(w, "<script>")
	fmt.Fprintf(w, "  function exportTable"+jsFuncName+"ToCSV() {")
	fmt.Fprintf(w, "    const table = document.getElementById('"+tableID+"');")
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
	fmt.Fprintf(w, "    csv.download = '"+fileName+".csv';")
	fmt.Fprintf(w, "    document.body.append(csv);")
	fmt.Fprintf(w, "    csv.click();")
	fmt.Fprintf(w, "    document.body.remove(csv);")
	fmt.Fprintf(w, "    window.open('/"+path+"', '_blank');")
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
	} else if data == "group_id" {
		dbSelectWhere.column = "group_id"
	} else if data == "group_name" {
		dbSelectWhere.column = "group_name"
	} else if data == "pbx_id" {
		dbSelectWhere.column = "pbx_id"
	} else if data == "pbx_name" {
		dbSelectWhere.column = "pbx_name"
	} else {
		panic("The function userAccountData can only accept the following arguments: type_id, group_id, group_name or pbx_id")
	}
	dbSelectWhere.columnWhere = "user_account_email"
	dbSelectWhere.columnWhereValue = dbUserAccountData.columnWhereValue

	return selectWhere(dbSelectWhere)
}

// Function to create a HTTP server that serves files with path named download
func download(id string, file string) {
	http.HandleFunc("/download/"+id+"/"+file, func(w http.ResponseWriter, r *http.Request) {
		dirPathB := "/var/lib/yap/call-recording/" + file + ""
		http.ServeFile(w, r, dirPathB)
	})
}

// Function to create a HTTP server with path named play-audio
func playAudio(id string, file string) {
	http.HandleFunc("/play-audio", func(w http.ResponseWriter, r *http.Request) {
	})
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
	fmt.Fprintf(w, "    <th>Total Groups</th>")
	fmt.Fprintf(w, "    <th>Total PBXs</th>")
	fmt.Fprintf(w, "    <th>Total SIP Endpoints</th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "  <tr>")
	dbTotalTableCount.table = "user_group"
	dbTotalTableCount.countMinusOne = true
	fmt.Fprintf(w, "    <td>"+totalTableCount(w, dbTotalTableCount)+"</td>")
	dbTotalTableCount.table = "pbx"
	dbTotalTableCount.countMinusOne = true
	fmt.Fprintf(w, "    <td>"+totalTableCount(w, dbTotalTableCount)+"</td>")
	dbTotalTableCount.table = "ps_endpoints"
	dbTotalTableCount.countMinusOne = false
	fmt.Fprintf(w, "    <td>"+totalTableCount(w, dbTotalTableCount)+"</td>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "</table>")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<table id=\"table\" class=\"table-main-menu\">")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>Total YAP<br>Admin<br>Accounts<br>(Type ID: 100)</th>")
	fmt.Fprintf(w, "    <th>Total Group<br>Admin<br>Accounts<br>(Type ID: 200)</th>")
	fmt.Fprintf(w, "    <th>Total Group<br>Regular<br>Accounts<br>(Type ID: 201)</th>")
	fmt.Fprintf(w, "    <th>Total PBX<br>Admin<br>Accounts<br>(Type ID: 300)</th>")
	fmt.Fprintf(w, "    <th>Total PBX<br>Regular<br>Accounts<br>(Type ID: 301)</th>")
	fmt.Fprintf(w, "    <th>Total PBX<br>Read Only<br>Accounts<br>(Type ID: 302)</th>")
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
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "</table>")

}

func mainMenuGroupAccount(w http.ResponseWriter, dbGroupAccount databaseFunctionParameter) {

	result, err := dbGroupAccount.connection.Query(`SELECT
					                  group_name,
					                  group_id,
					                  group_site_address_line_1,
					                  group_site_address_line_2,
					                  group_site_city_town_village,
					                  group_site_county_state_region,
					                  group_site_postcode_zip_code,
					                  group_site_country,
					                  group_site_contact_email,
					                  group_site_contact_number,
					                  group_invoice_address_line_1,
					                  group_invoice_address_line_2,
					                  group_invoice_city_town_village,
					                  group_invoice_county_state_region,
					                  group_invoice_postcode_zip_code,
					                  group_invoice_country,
					                  group_invoice_contact_email,
					                  group_invoice_contact_number,
					                  pbx_id
					                FROM
					                  yap.view___account_detail
					                WHERE
					                  user_account_email = ?;`, dbGroupAccount.columnWhereValue)

	// Error
	if err != nil {
		panic(err)
	}

	for result.Next() {
		var (
			groupName                     string
			groupID                       string
			groupSiteAddressLine1         string
			groupSiteAddressLine2         string
			groupSiteCityTownVillage      string
			groupSiteCountyStateRegion    string
			groupSitePostcodeZipCode      string
			groupSiteCountry              string
			groupSiteContactEmail         string
			groupSiteContactNumber        string
			groupInvoiceAddressLine1      string
			groupInvoiceAddressLine2      string
			groupInvoiceCityTownVillage   string
			groupInvoiceCountyStateRegion string
			groupInvoicePostcodeZipCode   string
			groupInvoiceCountry           string
			groupInvoiceContactEmail      string
			groupInvoiceContactNumber     string
			pbxID                         string
		)

		err = result.Scan(
			&groupName,
			&groupID,
			&groupSiteAddressLine1,
			&groupSiteAddressLine2,
			&groupSiteCityTownVillage,
			&groupSiteCountyStateRegion,
			&groupSitePostcodeZipCode,
			&groupSiteCountry,
			&groupSiteContactEmail,
			&groupSiteContactNumber,
			&groupInvoiceAddressLine1,
			&groupInvoiceAddressLine2,
			&groupInvoiceCityTownVillage,
			&groupInvoiceCountyStateRegion,
			&groupInvoicePostcodeZipCode,
			&groupInvoiceCountry,
			&groupInvoiceContactEmail,
			&groupInvoiceContactNumber,
			&pbxID,
		)

		// Error
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, "<table id=\"table\" class=\"table-main-menu\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>Group Name and ID</th>")
		fmt.Fprintf(w, "    <th>Total PBX(s) in Group</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td>Group Name: "+groupName+"<br><br>Group ID: "+groupID+"</td>")
		var dbTotalTableCountWhere databaseFunctionParameter
		dbTotalTableCountWhere.connection = dbGroupAccount.connection
		dbTotalTableCountWhere.database = dbGroupAccount.database
		dbTotalTableCountWhere.table = "pbx"
		dbTotalTableCountWhere.columnWhere = "id"
		dbTotalTableCountWhere.columnWhereValue = pbxID
		fmt.Fprintf(w, "    <td>"+totalTableCountWhere(w, dbTotalTableCountWhere)+"</td>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<table id=\"table\" class=\"table-main-menu\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>Group Site Address</th>")
		fmt.Fprintf(w, "    <th>Group Site Email</th>")
		fmt.Fprintf(w, "    <th>Group Site Phone Number</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td style=\"text-align: left;\">"+groupSiteAddressLine1+"<br>"+groupSiteAddressLine2+"<br>"+groupSiteCityTownVillage+"<br>"+groupSiteCountyStateRegion+"<br><br>"+groupSitePostcodeZipCode+"<br><br>"+groupSiteCountry+"</td>")
		fmt.Fprintf(w, "    <td><a href=\"mailto:"+groupSiteContactEmail+"\">"+groupSiteContactEmail+"</a></td>")
		fmt.Fprintf(w, "    <td><a href=\"tel:"+groupSiteContactNumber+"\">"+groupSiteContactNumber+"</a></td>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<table id=\"table\" class=\"table-main-menu\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>Group Invoice Address</th>")
		fmt.Fprintf(w, "    <th>Group Invoice Email</th>")
		fmt.Fprintf(w, "    <th>Group Invoice Phone Number</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td style=\"text-align: left;\">"+groupInvoiceAddressLine1+"<br>"+groupInvoiceAddressLine2+"<br>"+groupInvoiceCityTownVillage+"<br>"+groupInvoiceCountyStateRegion+"<br><br>"+groupInvoicePostcodeZipCode+"<br><br>"+groupInvoiceCountry+"</td>")
		fmt.Fprintf(w, "    <td><a href=\"mailto:"+groupInvoiceContactEmail+"\">"+groupInvoiceContactEmail+"</a></td>")
		fmt.Fprintf(w, "    <td><a href=\"tel:"+groupInvoiceContactNumber+"\">"+groupInvoiceContactNumber+"</a></td>")
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
		dbTotalTableCountWhere.table = "ps_endpoints"
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
	fmt.Fprintf(w, "<br>")

	var dbDetail databaseFunctionParameter
	dbDetail.connection = dbUserInformation.connection
	dbDetail.database = dbUserInformation.database
	dbDetail.columnWhereValue = dbUserInformation.columnWhereValue

	if userTypeID == "100" {
		mainMenuYapAccount(w, dbDetail)
	} else if userTypeID == "200" || userTypeID == "201" {
		mainMenuGroupAccount(w, dbDetail)
	} else if userTypeID == "300" || userTypeID == "301" || userTypeID == "302" {
		mainMenuPBXAccount(w, dbDetail)
	} else {
	}
	fmt.Fprintf(w, "</div>")
	toggleDivJS(w, "toggleAccountDetail", "account-detail-div")
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
		userAccountTypeID    string
		userAccountFirstName string
		userAccountLastName  string
		userAccountEmail     string
		userAccountType      string
		userAccountDateAdded string
		groupID              string
		groupName            string
		pbxID                string
		pbxName              string
	)

	ownUserAccountSQL, err := dbDetail.connection.Query(`SELECT
							       user_account_first_name,
							       user_account_last_name,
							       user_account_email,
							       user_account_type,
							       user_account_date_added,
							       group_id,
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
			&groupID,
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
				fmt.Fprintf(w, "    <th>Total Group<br>Admin<br>Accounts<br>(Type ID: 200)</th>")
				fmt.Fprintf(w, "    <th>Total Group<br>Regular<br>Accounts<br>(Type ID: 201)</th>")
			}
			if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" || userTypeID == "300" {
				fmt.Fprintf(w, "    <th>Total PBX<br>Admin<br>Accounts<br>(Type ID: 300)</th>")
				fmt.Fprintf(w, "    <th>Total PBX<br>Regular<br>Accounts<br>(Type ID: 301)</th>")
				fmt.Fprintf(w, "    <th>Total PBX<br>Read Only<br>Accounts<br>(Type ID: 302)</th>")
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
			} else if userTypeID == "200" || userTypeID == "201" {
				dbTotalTableCountWhere.columnWhereAnd = "group_id"
				dbTotalTableCountWhere.columnWhereValueAnd = groupID
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
			fmt.Fprintf(w, "    <th><button onclick=\"toggleOtherAccount() \"class=\"button-general button-user-account\">&nbsp Show/Hide Other Account(s) &nbsp</button></th>")
			fmt.Fprintf(w, "  </tr>")
		}
		fmt.Fprintf(w, "</table>")
	}

	if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" || userTypeID == "300" {

		userGroupID := userAccountData(dbDetail, "group_id")
		userGroupName := userAccountData(dbDetail, "group_name")
		userPBXID := userAccountData(dbDetail, "pbx_id")
		userPBXName := userAccountData(dbDetail, "pbx_name")

		fmt.Fprintf(w, "<div id=\"other-account-div\" style=\"display:none\">")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<table id=\"table\" class=\"table-user-account\">")
		fmt.Fprintf(w, "  <tr>")
		if userTypeID == "100" {
			fmt.Fprintf(w, "    <th class=\"table-title\";>All User Account Details on the Server:</th>")
		} else if userTypeID == "200" {
			fmt.Fprintf(w, "    <th class=\"table-title\";>User Account Details Within the Group<br>"+userGroupName+"<br>(Group ID: "+userGroupID+")</th>")
		} else if userTypeID == "201" {
			fmt.Fprintf(w, "    <th class=\"table-title\";>PBX User Account Details Within the Group<br>"+userGroupName+"<br>(Group ID: "+userGroupID+")</th>")
		} else if userTypeID == "300" {
			fmt.Fprintf(w, "    <th class=\"table-title\";>PBX User Account Details Within the PBX<br>"+userPBXName+"<br>(PBX ID: "+userPBXID+")</th>")
		}
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTML(w, "otherAccountSearchName", "other-account-input-name", "Name")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTML(w, "otherAccountSearchEmail", "other-account-input-email", "Email")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTML(w, "otherAccountSearchDate", "other-account-input-date", "Date Created")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTML(w, "otherAccountSearchType", "other-account-input-type", "Account Type")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
			fmt.Fprintf(w, "    <br><br>")
			inputTableHTML(w, "otherAccountSearchPBXName", "other-account-input-pbx-name", "PBX Name")
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTML(w, "otherAccountSearchPBXID", "other-account-input-pbx-id", "PBX ID")
		}
		if userTypeID == "100" {
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTML(w, "otherAccountSearchGroupName", "other-account-input-group-name", "Group Name")
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTML(w, "otherAccountSearchGroupID", "other-account-input-group-id", "Group ID")
		}
		fmt.Fprintf(w, "    <br><br>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		exportCSVButtonHTML(w, "OtherAccount", "button-user-account")
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
			fmt.Fprintf(w, "          <th>Group Name</th>")
			fmt.Fprintf(w, "          <th>Group ID</th>")
		}

		fmt.Fprintf(w, "        </tr>")

		if userTypeID == "100" {

			otherUserAccountSQL, err := dbDetail.connection.Query(`SELECT
						     			 user_account_first_name,
						     			 user_account_last_name,  
						     			 user_account_email,                                                   
						     			 user_account_type,  
						     			 user_account_date_added, 
						     			 group_id,
						     			 group_name,
						     			 pbx_id,
						     			 pbx_name						     
								       FROM
								         yap.view___account_detail;`)

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
					&groupID,
					&groupName,
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
				if groupName != "system" {
					fmt.Fprintf(w, "          <td>"+groupName+"</td>")
				} else {
					fmt.Fprintf(w, "          <td>-</td>")
				}
				if groupID != "1" {
					fmt.Fprintf(w, "          <td>"+groupID+"</td>")
				} else {
					fmt.Fprintf(w, "          <td>-</td>")
				}
				fmt.Fprintf(w, "        </tr>")
			}

		} else if userTypeID == "200" || userTypeID == "201" {

			otherUserAccountSQL, err := dbDetail.connection.Query(`SELECT
                                                                    	 user_account_type_id,
                                                                         user_account_first_name,
                                          				 user_account_last_name,
                                          				 user_account_email,
                                          				 user_account_type,
                                          				 user_account_date_added,
                                          				 pbx_id,
                                          				 pbx_name
								       FROM
								         yap.view___account_detail
								       WHERE
								         group_id =?;`, userGroupID)

			// Error
			if err != nil {
				panic(err)
			}

			for otherUserAccountSQL.Next() {

				err = otherUserAccountSQL.Scan(
					&userAccountTypeID,
					&userAccountFirstName,
					&userAccountLastName,
					&userAccountEmail,
					&userAccountType,
					&userAccountDateAdded,
					&pbxID,
					&pbxName,
				)

				// Error
				if err != nil {
					panic(err)
				}

				if userTypeID == "200" {
					fmt.Fprintf(w, "        <tr>")
					fmt.Fprintf(w, "          <td>"+userAccountFirstName+" "+userAccountLastName+"</td>")
					fmt.Fprintf(w, "          <td>"+userAccountEmail+"</td>")
					fmt.Fprintf(w, "          <td>"+userAccountType+"</td>")
					fmt.Fprintf(w, "          <td>"+userAccountDateAdded+"</td>")
					if pbxID != "1" {
						fmt.Fprintf(w, "          <td>"+pbxName+"</td>")
						fmt.Fprintf(w, "          <td>"+pbxID+"</td>")
					} else {
						fmt.Fprintf(w, "          <td>-</td>")
						fmt.Fprintf(w, "          <td>-</td>")
					}
					fmt.Fprintf(w, "        </tr>")
				} else if userTypeID == "201" {
					if userAccountTypeID == "300" || userAccountTypeID == "301" || userAccountTypeID == "302" {
						fmt.Fprintf(w, "        <tr>")
						fmt.Fprintf(w, "          <td>"+userAccountFirstName+" "+userAccountLastName+"</td>")
						fmt.Fprintf(w, "          <td>"+userAccountEmail+"</td>")
						fmt.Fprintf(w, "          <td>"+userAccountType+"</td>")
						fmt.Fprintf(w, "          <td>"+userAccountDateAdded+"</td>")
						fmt.Fprintf(w, "          <td>"+pbxName+"</td>")
						fmt.Fprintf(w, "          <td>"+pbxID+"</td>")
						fmt.Fprintf(w, "        </tr>")
					}
				}
			}

		} else if userTypeID == "300" {

			otherUserAccountSQL, err := dbDetail.connection.Query(`SELECT
                                                                         user_account_first_name,
                                                                         user_account_last_name,
                                                                         user_account_email,
                                                                         user_account_type,
                                                                         user_account_date_added
                                                                       FROM
                                                                         yap.view___account_detail
                                                                       WHERE
                                                                         group_id =? AND pbx_id =?;`, userGroupID, userPBXID)

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
				fmt.Fprintf(w, "        </tr>")
			}
		}
		fmt.Fprintf(w, "      </table>")
		filterTableJS(w, "otherAccountSearchName", "other-account-input-name", "other-account-table", 0)
		filterTableJS(w, "otherAccountSearchEmail", "other-account-input-email", "other-account-table", 1)
		filterTableJS(w, "otherAccountSearchType", "other-account-input-type", "other-account-table", 2)
		filterTableJS(w, "otherAccountSearchDate", "other-account-input-date", "other-account-table", 3)
		if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
			filterTableJS(w, "otherAccountSearchPBXName", "other-account-input-pbx-name", "other-account-table", 4)
			filterTableJS(w, "otherAccountSearchPBXID", "other-account-input-pbx-id", "other-account-table", 5)
		}
		if userTypeID == "100" {
			filterTableJS(w, "otherAccountSearchGroupName", "other-account-input-group-name", "other-account-table", 6)
			filterTableJS(w, "otherAccountSearchGroupID", "other-account-input-group-id", "other-account-table", 7)
		}
		exportCSVJS(w, "OtherAccount", "other-account-table", "YAP_user_account_details", "user-account")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "</div>")
		toggleDivJS(w, "toggleOtherAccount", "other-account-div")
	}
}

func userAccountAdd() {

}

func userAccountEdit() {

}

func userAccountDelete() {

}

//----------------------------------------------------------------------------------------------------

// Group page functions

func groupList(w http.ResponseWriter, dbDetail databaseFunctionParameter, userTypeID string, userGroupID string) {

	var (
		groupName                               string
		groupID                                 string
		groupDateAdded                          string
		groupActive                             string
		pbxLimit                                string
		newPBXSIPEndpointDefaultLimit           string
		newPBXSIPTrunkDefaultLimit              string
		newPBXPhoneNumberDefaultLimit           string
		newPBXCDRDefaultLimit                   string
		newPBXVoicemailDefaultMegabyteLimit     string
		newPBXCallRecordingDefaultMegabyteLimit string
		groupSiteAddressLine1                   string
		groupSiteAddressLine2                   string
		groupSiteCityTownVillage                string
		groupSiteCountyStateRegion              string
		groupSitePostcodeZipCode                string
		groupSiteCountry                        string
		groupSiteContactEmail                   string
		groupSiteContactNumber                  string
		groupInvoiceAddressLine1                string
		groupInvoiceAddressLine2                string
		groupInvoiceCityTownVillage             string
		groupInvoiceCountyStateRegion           string
		groupInvoicePostcodeZipCode             string
		groupInvoiceCountry                     string
		groupInvoiceContactEmail                string
		groupInvoiceContactNumber               string
	)

	var dbTableCountUserGroup databaseFunctionParameter
	dbTableCountUserGroup.connection = dbDetail.connection
	dbTableCountUserGroup.database = dbDetail.database
	dbTableCountUserGroup.table = "user_group"
	dbTableCountUserGroup.columnWhere = "group_active"

	if userTypeID == "100" {
		fmt.Fprintf(w, "<table id=\"table\" class=\"table-group\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table id=\"table\" class=\"table-group\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <th>Total Groups On YAP</th>")
		fmt.Fprintf(w, "          <th>Total Active Groups</th>")
		fmt.Fprintf(w, "          <th>Total Inactive Groups</th>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		dbTableCountUserGroup.countMinusOne = true
		fmt.Fprintf(w, "          <td>"+totalTableCount(w, dbTableCountUserGroup)+"</td>")
		dbTableCountUserGroup.columnWhereValue = "1"
		fmt.Fprintf(w, "          <td>"+totalTableCountWhere(w, dbTableCountUserGroup)+"</td>")
		dbTableCountUserGroup.columnWhereValue = "0"
		dbTableCountUserGroup.countMinusOne = false
		fmt.Fprintf(w, "          <td>"+totalTableCountWhere(w, dbTableCountUserGroup)+"</td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><button onclick=\"toggleGroup() \"class=\"button-general button-group\">&nbsp Show/Hide Group(s) &nbsp</button></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
	}

	if userTypeID == "100" {
		fmt.Fprintf(w, "<div id=\"group-div\" style=\"display:none\">")
		fmt.Fprintf(w, "<br>")
	} else {
		fmt.Fprintf(w, "<div id=\"group-div\">")
	}
	fmt.Fprintf(w, "<table id=\"table\" class=\"table-group\">")
	fmt.Fprintf(w, "  <tr>")
	if userTypeID == "100" {
		fmt.Fprintf(w, "    <th class=\"table-title\";>All Group Contact Details on the Server:</th>")
	} else {
		fmt.Fprintf(w, "    <th class=\"table-title\";>Group Contact Details</th>")
	}
	fmt.Fprintf(w, "  </tr>")
	if userTypeID == "100" {
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "    <br>")
		inputTableHTML(w, "groupContactSearchGroupName", "group-contact-input-group-name", "Group Name")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTML(w, "groupContactSearchGroupID", "group-contact-input-group-id", "Group ID")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTML(w, "groupContactSearchSiteAddress", "group-contact-input-site-address", "Group Site Address")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTML(w, "groupContactSearchSiteEmail", "group-contact-input-site-email", "Group Site Email Address")
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    <br>")
		inputTableHTML(w, "groupContactSearchSitePhone", "group-contact-input-site-phone", "Group Site Phone Number")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTML(w, "groupContactSearchInvoiceAddress", "group-contact-input-invoice-address", "Group Invoice Address")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTML(w, "groupContactSearchInvoiceEmail", "group-contact-input-invoice-email", "Group Invoice Email Address")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTML(w, "groupContactSearchInvoicePhone", "group-contact-input-invoice-phone", "Group Invoice Phone Number")
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
	}
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>")
	exportCSVButtonHTML(w, "GroupContact", "button-group")
	fmt.Fprintf(w, "    </th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>")
	fmt.Fprintf(w, "      <table id=\"group-contact-table\" class=\"table-group\">")
	fmt.Fprintf(w, "        <tr>")
	fmt.Fprintf(w, "          <th>Group Name</th>")
	fmt.Fprintf(w, "          <th>Group ID</th>")
	fmt.Fprintf(w, "          <th>Group Site<br> Address</th>")
	fmt.Fprintf(w, "          <th>Group Site<br> Email Address</th>")
	fmt.Fprintf(w, "          <th>Group Site<br> Phone Number</th>")
	fmt.Fprintf(w, "          <th>Group Invoice<br> Address</th>")
	fmt.Fprintf(w, "          <th>Group Invoice<br> Email Address</th>")
	fmt.Fprintf(w, "          <th>Group Invoice<br> Phone Number</th>")
	fmt.Fprintf(w, "        </tr>")
	if userTypeID == "100" {
		groupAllSQL, err := dbDetail.connection.Query(`SELECT
							group_name,
							group_id,
							group_site_address_line_1,
					                group_site_address_line_2,
					                group_site_city_town_village,
					                group_site_county_state_region,
					                group_site_postcode_zip_code,
					                group_site_country,
					                group_site_contact_email,
					                group_site_contact_number,
					                group_invoice_address_line_1,
					                group_invoice_address_line_2,
					                group_invoice_city_town_village,
					                group_invoice_county_state_region,
					                group_invoice_postcode_zip_code,
					                group_invoice_country,
					                group_invoice_contact_email,
					                group_invoice_contact_number
					              FROM
					  	        yap.view___group_detail
						      WHERE
						        group_id != 1;`)

		// Error
		if err != nil {
			panic(err)

		}

		for groupAllSQL.Next() {

			err = groupAllSQL.Scan(
				&groupName,
				&groupID,
				&groupSiteAddressLine1,
				&groupSiteAddressLine2,
				&groupSiteCityTownVillage,
				&groupSiteCountyStateRegion,
				&groupSitePostcodeZipCode,
				&groupSiteCountry,
				&groupSiteContactEmail,
				&groupSiteContactNumber,
				&groupInvoiceAddressLine1,
				&groupInvoiceAddressLine2,
				&groupInvoiceCityTownVillage,
				&groupInvoiceCountyStateRegion,
				&groupInvoicePostcodeZipCode,
				&groupInvoiceCountry,
				&groupInvoiceContactEmail,
				&groupInvoiceContactNumber,
			)

			// Error
			if err != nil {
				panic(err)
			}
			fmt.Fprintf(w, "        <tr>")
			fmt.Fprintf(w, "          <td>"+groupName+"</td>")
			fmt.Fprintf(w, "          <td>"+groupID+"</td>")
			fmt.Fprintf(w, "          <td style=\"text-align: left;\">"+groupSiteAddressLine1+"&nbsp<br>"+groupSiteAddressLine2+"&nbsp<br>"+groupSiteCityTownVillage+"&nbsp<br>"+groupSiteCountyStateRegion+"&nbsp<br><br>"+groupSitePostcodeZipCode+"&nbsp<br><br>"+groupSiteCountry+"&nbsp</td>")
			fmt.Fprintf(w, "          <td><a href=\"mailto:"+groupSiteContactEmail+"\">"+groupSiteContactEmail+"</a></td>")
			fmt.Fprintf(w, "          <td><a href=\"tel:"+groupSiteContactNumber+"\">"+groupSiteContactNumber+"</a></td>")
			fmt.Fprintf(w, "          <td style=\"text-align: left;\">"+groupInvoiceAddressLine1+"&nbsp<br>"+groupInvoiceAddressLine2+"&nbsp<br>"+groupInvoiceCityTownVillage+"&nbsp<br>"+groupInvoiceCountyStateRegion+"&nbsp<br><br>"+groupInvoicePostcodeZipCode+"&nbsp<br><br>"+groupInvoiceCountry+"&nbsp</td>")
			fmt.Fprintf(w, "          <td>&nbsp<a href=\"mailto:"+groupInvoiceContactEmail+"\">"+groupInvoiceContactEmail+"</a></td>")
			fmt.Fprintf(w, "          <td><a href=\"tel:"+groupInvoiceContactNumber+"\">"+groupInvoiceContactNumber+"</a></td>")
			fmt.Fprintf(w, "        </tr>")
		}
	} else {
		groupWhereSQL, err := dbDetail.connection.Query(`SELECT
							group_name,
							group_id,
							group_site_address_line_1,
					                group_site_address_line_2,
					                group_site_city_town_village,
					                group_site_county_state_region,
					                group_site_postcode_zip_code,
					                group_site_country,
					                group_site_contact_email,
					                group_site_contact_number,
					                group_invoice_address_line_1,
					                group_invoice_address_line_2,
					                group_invoice_city_town_village,
					                group_invoice_county_state_region,
					                group_invoice_postcode_zip_code,
					                group_invoice_country,
					                group_invoice_contact_email,
					                group_invoice_contact_number
					            FROM
					  	        yap.view___group_detail
						    WHERE
					  	        group_id = ?;`, userGroupID)

		// Error
		if err != nil {
			panic(err)

		}

		for groupWhereSQL.Next() {

			err = groupWhereSQL.Scan(
				&groupName,
				&groupID,
				&groupSiteAddressLine1,
				&groupSiteAddressLine2,
				&groupSiteCityTownVillage,
				&groupSiteCountyStateRegion,
				&groupSitePostcodeZipCode,
				&groupSiteCountry,
				&groupSiteContactEmail,
				&groupSiteContactNumber,
				&groupInvoiceAddressLine1,
				&groupInvoiceAddressLine2,
				&groupInvoiceCityTownVillage,
				&groupInvoiceCountyStateRegion,
				&groupInvoicePostcodeZipCode,
				&groupInvoiceCountry,
				&groupInvoiceContactEmail,
				&groupInvoiceContactNumber,
			)

			// Error
			if err != nil {
				panic(err)
			}
			fmt.Fprintf(w, "        <tr>")
			fmt.Fprintf(w, "          <td>"+groupName+"</td>")
			fmt.Fprintf(w, "          <td>"+groupID+"</td>")
			fmt.Fprintf(w, "          <td style=\"text-align: left;\">"+groupSiteAddressLine1+"&nbsp<br>"+groupSiteAddressLine2+"<br>"+groupSiteCityTownVillage+"<br>"+groupSiteCountyStateRegion+"<br><br>"+groupSitePostcodeZipCode+"<br><br>"+groupSiteCountry+"</td>")
			fmt.Fprintf(w, "          <td>&nbsp<a href=\"mailto:"+groupSiteContactEmail+"\">"+groupSiteContactEmail+"</a></td>")
			fmt.Fprintf(w, "          <td>&nbsp<a href=\"tel:"+groupSiteContactNumber+"\">"+groupSiteContactNumber+"</a></td>")
			fmt.Fprintf(w, "          <td style=\"text-align: left;\">"+groupInvoiceAddressLine1+"&nbsp<br>"+groupInvoiceAddressLine2+"<br>"+groupInvoiceCityTownVillage+"<br>"+groupInvoiceCountyStateRegion+"<br><br>"+groupInvoicePostcodeZipCode+"<br><br>"+groupInvoiceCountry+"</td>")
			fmt.Fprintf(w, "          <td>&nbsp<a href=\"mailto:"+groupInvoiceContactEmail+"\">"+groupInvoiceContactEmail+"</a></td>")
			fmt.Fprintf(w, "          <td>&nbsp<a href=\"tel:"+groupInvoiceContactNumber+"\">"+groupInvoiceContactNumber+"</a></td>")
			fmt.Fprintf(w, "        </tr>")
		}
	}
	fmt.Fprintf(w, "      </table>")
	if userTypeID == "100" {
		filterTableJS(w, "groupContactSearchGroupName", "group-contact-input-group-name", "group-contact-table", 0)
		filterTableJS(w, "groupContactSearchGroupID", "group-contact-input-group-id", "group-contact-table", 1)
		filterTableJS(w, "groupContactSearchSiteAddress", "group-contact-input-site-address", "group-contact-table", 2)
		filterTableJS(w, "groupContactSearchSiteEmail", "group-contact-input-site-email", "group-contact-table", 3)
		filterTableJS(w, "groupContactSearchSitePhone", "group-contact-input-site-phone", "group-contact-table", 4)
		filterTableJS(w, "groupContactSearchInvoiceAdddress", "group-contact-input-Invoice-address", "group-contact-table", 5)
		filterTableJS(w, "groupContactSearchInvoiceEmail", "group-contact-input-invoice-email", "group-contact-table", 6)
		filterTableJS(w, "groupContactSearchInvoicePhone", "group-contact-input-invoice-phone", "group-contact-table", 7)
	}
	exportCSVJS(w, "GroupContact", "group-contact-table", "YAP_group_contact_details", "group")
	fmt.Fprintf(w, "    </th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "</table>")

	// Group resource table
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<table id=\"table\" class=\"table-group\">")
	fmt.Fprintf(w, "  <tr>")
	if userTypeID == "100" {
		fmt.Fprintf(w, "    <th class=\"table-title\";>All PBX  Default Resource Limits on the Server for Each Group:</th>")
	} else {
		fmt.Fprintf(w, "    <th class=\"table-title\";>PBX Default Resource Limits for Group</th>")
	}
	fmt.Fprintf(w, "  </tr>")
	if userTypeID == "100" {
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "    <br>")
		inputTableHTML(w, "groupResourceSearchGroupName", "group-resource-input-group-name", "Group Name")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTML(w, "groupResourceSearchGroupID", "group-resource-input-group-id", "Group ID")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTML(w, "groupResourceSearchDate", "group-resource-input-date", "Date Created")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTML(w, "groupResourceSearchActive", "group-resource-input-active", "Group Active Status")
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
	}
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>")
	exportCSVButtonHTML(w, "GroupResource", "button-group")
	fmt.Fprintf(w, "    </th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>")
	fmt.Fprintf(w, "      <table id=\"group-resource-table\" class=\"table-group\">")
	fmt.Fprintf(w, "        <tr>")
	fmt.Fprintf(w, "          <th>Group Name</th>")
	fmt.Fprintf(w, "          <th>Group ID</th>")
	fmt.Fprintf(w, "          <th>Group Date</th>")
	fmt.Fprintf(w, "          <th>Group Active <br>Status</th>")
	fmt.Fprintf(w, "          <th>PBX Limit</th>")
	fmt.Fprintf(w, "          <th>SIP Endpoint <br>Limit for <br>a New PBX</th>")
	fmt.Fprintf(w, "          <th>SIP Trunk <br>Limit for <br>a New PBX</th>")
	fmt.Fprintf(w, "          <th>Phone Number <br>Limit for <br>a new PBX</th>")
	fmt.Fprintf(w, "          <th>CDR Limit <br>for a New PBX</th>")
	fmt.Fprintf(w, "          <th>Voicemail <br>Limit for a New PBX <br>(Megabytes)</th>")
	fmt.Fprintf(w, "          <th>Call Recording <br>Limit for a New PBX <br>(Megabytes)</th>")
	fmt.Fprintf(w, "        </tr>")
	if userTypeID == "100" {
		groupAllSQL, err := dbDetail.connection.Query(`SELECT
							group_name,
							group_id,
							group_date_added,
							group_active,
							pbx_limit,
							new_pbx_sip_endpoint_default_limit,
							new_pbx_sip_trunk_default_limit,
							new_pbx_phone_number_default_limit,
							new_pbx_cdr_default_limit,
							new_pbx_voicemail_default_megabyte_limit,
							new_pbx_call_recording_default_megabyte_limit
					              FROM
					  	        yap.view___group_detail
						      WHERE
						        group_id != 1;`)

		// Error
		if err != nil {
			panic(err)

		}

		for groupAllSQL.Next() {

			err = groupAllSQL.Scan(
				&groupName,
				&groupID,
				&groupDateAdded,
				&groupActive,
				&pbxLimit,
				&newPBXSIPEndpointDefaultLimit,
				&newPBXSIPTrunkDefaultLimit,
				&newPBXPhoneNumberDefaultLimit,
				&newPBXCDRDefaultLimit,
				&newPBXVoicemailDefaultMegabyteLimit,
				&newPBXCallRecordingDefaultMegabyteLimit,
			)

			// Error
			if err != nil {
				panic(err)
			}
			fmt.Fprintf(w, "        <tr>")
			fmt.Fprintf(w, "          <td>"+groupName+"</td>")
			fmt.Fprintf(w, "          <td>"+groupID+"</td>")
			fmt.Fprintf(w, "          <td>"+groupDateAdded+"</td>")
			if groupActive == "1" {
				fmt.Fprintf(w, "          <td>YES</td>")
			} else {
				fmt.Fprintf(w, "          <td>NO</td>")
			}
			fmt.Fprintf(w, "          <td>"+pbxLimit+"</td>")
			fmt.Fprintf(w, "          <td>"+newPBXSIPEndpointDefaultLimit+"</td>")
			fmt.Fprintf(w, "          <td>"+newPBXSIPTrunkDefaultLimit+"</td>")
			fmt.Fprintf(w, "          <td>"+newPBXPhoneNumberDefaultLimit+"</td>")
			fmt.Fprintf(w, "          <td>"+newPBXCDRDefaultLimit+"</td>")
			fmt.Fprintf(w, "          <td>"+newPBXVoicemailDefaultMegabyteLimit+"</td>")
			fmt.Fprintf(w, "          <td>"+newPBXCallRecordingDefaultMegabyteLimit+"</td>")
			fmt.Fprintf(w, "        </tr>")
		}
	} else {
		groupWhereSQL, err := dbDetail.connection.Query(`SELECT
							group_name,
                                                        group_id,
                                                        group_date_added,
                                                        group_active,
                                                        pbx_limit,
                                                        new_pbx_sip_endpoint_default_limit,
                                                        new_pbx_sip_trunk_default_limit,
                                                        new_pbx_phone_number_default_limit,
                                                        new_pbx_cdr_default_limit,
                                                        new_pbx_voicemail_default_megabyte_limit,
                                                        new_pbx_call_recording_default_megabyte_limit
					            FROM
					  	        yap.view___group_detail
						    WHERE
					  	        group_id = ?;`, userGroupID)

		// Error
		if err != nil {
			panic(err)

		}

		for groupWhereSQL.Next() {

			err = groupWhereSQL.Scan(
				&groupName,
				&groupID,
				&groupDateAdded,
				&groupActive,
				&pbxLimit,
				&newPBXSIPEndpointDefaultLimit,
				&newPBXSIPTrunkDefaultLimit,
				&newPBXPhoneNumberDefaultLimit,
				&newPBXCDRDefaultLimit,
				&newPBXVoicemailDefaultMegabyteLimit,
				&newPBXCallRecordingDefaultMegabyteLimit,
			)

			// Error
			if err != nil {
				panic(err)
			}
			fmt.Fprintf(w, "        <tr>")
			fmt.Fprintf(w, "          <td>"+groupName+"</td>")
			fmt.Fprintf(w, "          <td>"+groupID+"</td>")
			fmt.Fprintf(w, "          <td>"+groupDateAdded+"</td>")
			fmt.Fprintf(w, "          <td>"+groupActive+"</td>")
			fmt.Fprintf(w, "          <td>"+pbxLimit+"</td>")
			fmt.Fprintf(w, "          <td>"+newPBXSIPEndpointDefaultLimit+"</td>")
			fmt.Fprintf(w, "          <td>"+newPBXSIPTrunkDefaultLimit+"</td>")
			fmt.Fprintf(w, "          <td>"+newPBXPhoneNumberDefaultLimit+"</td>")
			fmt.Fprintf(w, "          <td>"+newPBXCDRDefaultLimit+"</td>")
			fmt.Fprintf(w, "          <td>"+newPBXVoicemailDefaultMegabyteLimit+"</td>")
			fmt.Fprintf(w, "          <td>"+newPBXCallRecordingDefaultMegabyteLimit+"</td>")
			fmt.Fprintf(w, "        </tr>")
		}
	}
	fmt.Fprintf(w, "      </table>")
	if userTypeID == "100" {
		filterTableJS(w, "groupResourceSearchGroupName", "group-resource-input-group-name", "group-resource-table", 0)
		filterTableJS(w, "groupResourceSearchGroupID", "group-resource-input-group-id", "group-resource-table", 1)
		filterTableJS(w, "groupResourceSearchDate", "group-resource-input-date", "group-resource-table", 2)
		filterTableJS(w, "groupResourceSearchActive", "group-resource-input-active", "group-resource-table", 3)
	}
	exportCSVJS(w, "GroupResource", "group-resource-table", "YAP_group_resource_details", "group")
	fmt.Fprintf(w, "    </th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "</table>")
	fmt.Fprintf(w, "</div>")
	if userTypeID == "100" {
		toggleDivJS(w, "toggleGroup", "group-div")
	}

}

//----------------------------------------------------------------------------------------------------

// PBX page functions

func pbxList(w http.ResponseWriter, dbDetail databaseFunctionParameter, userTypeID string, userGroupID string, userPBXID string) {

	var (
		pbxName                       string
		pbxID                         string
		pbxDateAdded                  string
		pbxActive                     string
		pbxSIPEndpointLimit           string
		pbxSIPTrunkLimit              string
		pbxPhoneNumberLimit           string
		pbxCDRLimit                   string
		pbxVoicemailMegabyteLimit     string
		pbxCallRecordingMegabyteLimit string
		pbxSiteAddressLine1           string
		pbxSiteAddressLine2           string
		pbxSiteCityTownVillage        string
		pbxSiteCountyStateRegion      string
		pbxSitePostcodeZipCode        string
		pbxSiteCountry                string
		pbxSiteContactEmail           string
		pbxSiteContactNumber          string
		pbxInvoiceAddressLine1        string
		pbxInvoiceAddressLine2        string
		pbxInvoiceCityTownVillage     string
		pbxInvoiceCountyStateRegion   string
		pbxInvoicePostcodeZipCode     string
		pbxInvoiceCountry             string
		pbxInvoiceContactEmail        string
		pbxInvoiceContactNumber       string
		groupName                     string
		groupID                       string
	)

	var dbTableCountUserPBX databaseFunctionParameter
	dbTableCountUserPBX.connection = dbDetail.connection
	dbTableCountUserPBX.database = dbDetail.database
	dbTableCountUserPBX.table = "pbx"
	dbTableCountUserPBX.columnWhere = "id"

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
			fmt.Fprintf(w, "          <td>"+totalTableCount(w, dbTableCountUserPBX)+"</td>")
			dbTableCountUserPBX.columnWhereValue = "1"
			fmt.Fprintf(w, "          <td>"+totalTableCountWhere(w, dbTableCountUserPBX)+"</td>")
			dbTableCountUserPBX.columnWhereValue = "0"
			dbTableCountUserPBX.countMinusOne = false
			fmt.Fprintf(w, "          <td>"+totalTableCountWhere(w, dbTableCountUserPBX)+"</td>")
		} else if userTypeID == "200" || userTypeID == "201" {
			var dbTableCountUserPBXWhere databaseFunctionParameter
			dbTableCountUserPBXWhere.connection = dbDetail.connection
			dbTableCountUserPBXWhere.database = dbDetail.database
			dbTableCountUserPBXWhere.table = "pbx"
			dbTableCountUserPBXWhere.columnWhere = "group_id"
			dbTableCountUserPBXWhere.columnWhereValue = userGroupID
			dbTableCountUserPBXWhere.columnWhereAnd = "pbx_active"
			dbTableCountUserPBXWhere.columnWhereValueAnd = "1"
			fmt.Fprintf(w, "    <td>"+totalTableCountWhereAnd(w, dbTableCountUserPBXWhere)+"</td>")
			dbTableCountUserPBXWhere.columnWhereValueAnd = "0"
			fmt.Fprintf(w, "    <td>"+totalTableCountWhereAnd(w, dbTableCountUserPBXWhere)+"</td>")
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
		inputTableHTML(w, "pbxContactSearchPBXName", "pbx-contact-input-pbx-name", "PBX Name")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTML(w, "pbxContactSearchPBXID", "pbx-contact-input-pbx-id", "PBX ID")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTML(w, "pbxContactSearchSiteAddress", "pbx-contact-input-site-address", "PBX Site Address")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTML(w, "pbxContactSearchSiteEmail", "pbx-contact-input-site-email", "PBX Site Email Address")
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    <br>")
		inputTableHTML(w, "pbxContactSearchSitePhone", "pbx-contact-input-site-phone", "PBX Site Phone Number")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTML(w, "pbxContactSearchInvoiceAddress", "pbx-contact-input-invoice-address", "PBX Invoice Address")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTML(w, "pbxContactSearchInvoiceEmail", "pbx-contact-input-invoice-email", "PBX Invoice Email Address")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTML(w, "pbxContactSearchInvoicePhone", "pbx-contact-input-invoice-phone", "PBX Invoice Phone Number")
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    <br>")
		if userTypeID == "100" {
			inputTableHTML(w, "pbxContactSearchGroupName", "pbx-contact-input-group-name", "Group Name")
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTML(w, "pbxContactSearchGroupID", "pbx-contact-input-group-id", "Group ID")
			fmt.Fprintf(w, "    <br>")
			fmt.Fprintf(w, "    <br>")
		}
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
	}
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>")
	exportCSVButtonHTML(w, "PBXContact", "button-pbx")
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
		fmt.Fprintf(w, "          <th>Group Name</th>")
		fmt.Fprintf(w, "          <th>Group ID</th>")
	}
	fmt.Fprintf(w, "        </tr>")
	if userTypeID == "100" {
		pbxAllSQL, err := dbDetail.connection.Query(`SELECT
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
					                group_name,
					                group_id
					              FROM
					  	        yap.view___pbx_detail
						      WHERE
						        pbx_id != 1;`)

		// Error
		if err != nil {
			panic(err)

		}

		for pbxAllSQL.Next() {

			err = pbxAllSQL.Scan(
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
				&groupName,
				&groupID,
			)

			// Error
			if err != nil {
				panic(err)
			}
			fmt.Fprintf(w, "        <tr>")
			fmt.Fprintf(w, "          <td>"+pbxName+"</td>")
			fmt.Fprintf(w, "          <td>"+pbxID+"</td>")
			fmt.Fprintf(w, "          <td style=\"text-align: left;\">"+pbxSiteAddressLine1+"&nbsp<br>"+pbxSiteAddressLine2+"&nbsp<br>"+pbxSiteCityTownVillage+"&nbsp<br>"+pbxSiteCountyStateRegion+"&nbsp<br><br>"+pbxSitePostcodeZipCode+"&nbsp<br><br>"+pbxSiteCountry+"&nbsp</td>")
			fmt.Fprintf(w, "          <td><a href=\"mailto:"+pbxSiteContactEmail+"\">"+pbxSiteContactEmail+"</a></td>")
			fmt.Fprintf(w, "          <td><a href=\"tel:"+pbxSiteContactNumber+"\">"+pbxSiteContactNumber+"</a></td>")
			fmt.Fprintf(w, "          <td style=\"text-align: left;\">"+pbxInvoiceAddressLine1+"&nbsp<br>"+pbxInvoiceAddressLine2+"&nbsp<br>"+pbxInvoiceCityTownVillage+"&nbsp<br>"+pbxInvoiceCountyStateRegion+"&nbsp<br><br>"+pbxInvoicePostcodeZipCode+"&nbsp<br><br>"+pbxInvoiceCountry+"&nbsp</td>")
			fmt.Fprintf(w, "          <td>&nbsp<a href=\"mailto:"+pbxInvoiceContactEmail+"\">"+pbxInvoiceContactEmail+"</a></td>")
			fmt.Fprintf(w, "          <td><a href=\"tel:"+pbxInvoiceContactNumber+"\">"+pbxInvoiceContactNumber+"</a></td>")
			fmt.Fprintf(w, "          <td>"+groupName+"</td>")
			fmt.Fprintf(w, "          <td>"+groupID+"</td>")
			fmt.Fprintf(w, "        </tr>")
		}
	} else if userTypeID == "200" || userTypeID == "201" {
		pbxWhereUserGroupSQL, err := dbDetail.connection.Query(`SELECT
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
					  	        yap.view___pbx_detail
						      WHERE
						        group_id = ?;`, userGroupID)

		// Error
		if err != nil {
			panic(err)

		}

		for pbxWhereUserGroupSQL.Next() {

			err = pbxWhereUserGroupSQL.Scan(
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
			fmt.Fprintf(w, "        <tr>")
			fmt.Fprintf(w, "          <td>"+pbxName+"</td>")
			fmt.Fprintf(w, "          <td>"+pbxID+"</td>")
			fmt.Fprintf(w, "          <td style=\"text-align: left;\">"+pbxSiteAddressLine1+"&nbsp<br>"+pbxSiteAddressLine2+"&nbsp<br>"+pbxSiteCityTownVillage+"&nbsp<br>"+pbxSiteCountyStateRegion+"&nbsp<br><br>"+pbxSitePostcodeZipCode+"&nbsp<br><br>"+pbxSiteCountry+"&nbsp</td>")
			fmt.Fprintf(w, "          <td><a href=\"mailto:"+pbxSiteContactEmail+"\">"+pbxSiteContactEmail+"</a></td>")
			fmt.Fprintf(w, "          <td><a href=\"tel:"+pbxSiteContactNumber+"\">"+pbxSiteContactNumber+"</a></td>")
			fmt.Fprintf(w, "          <td style=\"text-align: left;\">"+pbxInvoiceAddressLine1+"&nbsp<br>"+pbxInvoiceAddressLine2+"&nbsp<br>"+pbxInvoiceCityTownVillage+"&nbsp<br>"+pbxInvoiceCountyStateRegion+"&nbsp<br><br>"+pbxInvoicePostcodeZipCode+"&nbsp<br><br>"+pbxInvoiceCountry+"&nbsp</td>")
			fmt.Fprintf(w, "          <td>&nbsp<a href=\"mailto:"+pbxInvoiceContactEmail+"\">"+pbxInvoiceContactEmail+"</a></td>")
			fmt.Fprintf(w, "          <td><a href=\"tel:"+pbxInvoiceContactNumber+"\">"+pbxInvoiceContactNumber+"</a></td>")
			fmt.Fprintf(w, "        </tr>")
		}
	} else if userTypeID == "300" || userTypeID == "301" || userTypeID == "302" {
		pbxWhereUserPBXSQL, err := dbDetail.connection.Query(`SELECT
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
					  	        yap.view___pbx_detail
						    WHERE
					  	        pbx_id = ?;`, userPBXID)

		// Error
		if err != nil {
			panic(err)

		}

		for pbxWhereUserPBXSQL.Next() {

			err = pbxWhereUserPBXSQL.Scan(
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
			fmt.Fprintf(w, "        <tr>")
			fmt.Fprintf(w, "          <td style=\"text-align: left;\">"+pbxSiteAddressLine1+"&nbsp<br>"+pbxSiteAddressLine2+"<br>"+pbxSiteCityTownVillage+"<br>"+pbxSiteCountyStateRegion+"<br><br>"+pbxSitePostcodeZipCode+"<br><br>"+pbxSiteCountry+"</td>")
			fmt.Fprintf(w, "          <td>&nbsp<a href=\"mailto:"+pbxSiteContactEmail+"\">"+pbxSiteContactEmail+"</a></td>")
			fmt.Fprintf(w, "          <td>&nbsp<a href=\"tel:"+pbxSiteContactNumber+"\">"+pbxSiteContactNumber+"</a></td>")
			fmt.Fprintf(w, "          <td style=\"text-align: left;\">"+pbxInvoiceAddressLine1+"&nbsp<br>"+pbxInvoiceAddressLine2+"<br>"+pbxInvoiceCityTownVillage+"<br>"+pbxInvoiceCountyStateRegion+"<br><br>"+pbxInvoicePostcodeZipCode+"<br><br>"+pbxInvoiceCountry+"</td>")
			fmt.Fprintf(w, "          <td>&nbsp<a href=\"mailto:"+pbxInvoiceContactEmail+"\">"+pbxInvoiceContactEmail+"</a></td>")
			fmt.Fprintf(w, "          <td>&nbsp<a href=\"tel:"+pbxInvoiceContactNumber+"\">"+pbxInvoiceContactNumber+"</a></td>")
			fmt.Fprintf(w, "        </tr>")
		}
	}
	fmt.Fprintf(w, "      </table>")
	if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
		filterTableJS(w, "pbxContactSearchPBXName", "pbx-contact-input-pbx-name", "pbx-contact-table", 0)
		filterTableJS(w, "pbxContactSearchPBXID", "pbx-contact-input-pbx-id", "pbx-contact-table", 1)
		filterTableJS(w, "pbxContactSearchSiteAddress", "pbx-contact-input-site-address", "pbx-contact-table", 2)
		filterTableJS(w, "pbxContactSearchSiteEmail", "pbx-contact-input-site-email", "pbx-contact-table", 3)
		filterTableJS(w, "pbxContactSearchSitePhone", "pbx-contact-input-site-phone", "pbx-contact-table", 4)
		filterTableJS(w, "pbxContactSearchInvoiceAdddress", "pbx-contact-input-Invoice-address", "pbx-contact-table", 5)
		filterTableJS(w, "pbxContactSearchInvoiceEmail", "pbx-contact-input-invoice-email", "pbx-contact-table", 6)
		filterTableJS(w, "pbxContactSearchInvoicePhone", "pbx-contact-input-invoice-phone", "pbx-contact-table", 7)
		if userTypeID == "100" {
			filterTableJS(w, "pbxContactSearchGroupName", "pbx-contact-input-group-name", "pbx-contact-table", 8)
			filterTableJS(w, "pbxContactSearchGroupID", "pbx-contact-input-group-id", "pbx-contact-table", 9)
		}
	}
	exportCSVJS(w, "PBXContact", "pbx-contact-table", "YAP_pbx_contact_details", "pbx")
	fmt.Fprintf(w, "    </th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "</table>")

	// Group resource table
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
		inputTableHTML(w, "pbxResourceSearchPBXName", "pbx-resource-input-pbx-name", "PBX Name")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTML(w, "pbxResourceSearchPBXID", "pbx-resource-input-pbx-id", "PBX ID")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTML(w, "pbxResourceSearchDate", "pbx-resource-input-date", "Date Created")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTML(w, "pbxResourceSearchActive", "pbx-resource-input-active", "PBX Active Status")
		fmt.Fprintf(w, "    <br>")
		fmt.Fprintf(w, "    <br>")
		if userTypeID == "100" {
			inputTableHTML(w, "pbxResourceSearchGroupName", "pbx-resource-input-group-name", "Group Name")
			fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
			inputTableHTML(w, "pbxResourceSearchGroupID", "pbx-resource-input-group-id", "Group ID")
			fmt.Fprintf(w, "    <br>")
			fmt.Fprintf(w, "    <br>")
		}
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
	}
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>")
	exportCSVButtonHTML(w, "PBXResource", "button-pbx")
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
	fmt.Fprintf(w, "          <th>PBX Active <br>Status</th>")
	fmt.Fprintf(w, "          <th>SIP Endpoint <br>Limit for <br>PBX</th>")
	fmt.Fprintf(w, "          <th>SIP Trunk <br>Limit for <br>PBX</th>")
	fmt.Fprintf(w, "          <th>Phone Number <br>Limit for <br>PBX</th>")
	fmt.Fprintf(w, "          <th>CDR Limit <br>for PBX</th>")
	fmt.Fprintf(w, "          <th>Voicemail <br>Limit for PBX <br>(Megabytes)</th>")
	fmt.Fprintf(w, "          <th>Call Recording <br>Limit for PBX <br>(Megabytes)</th>")
	if userTypeID == "100" {
		fmt.Fprintf(w, "          <th>Group Name</th>")
		fmt.Fprintf(w, "          <th>Group ID</th>")
	}
	fmt.Fprintf(w, "        </tr>")
	if userTypeID == "100" {
		pbxAllSQL, err := dbDetail.connection.Query(`SELECT
							pbx_name,
							pbx_id,
							pbx_date_added,
							pbx_active,
							pbx_sip_endpoint_limit,
							pbx_sip_trunk_limit,
							pbx_phone_number_limit,
							pbx_cdr_limit,
							pbx_voicemail_megabyte_limit,
							pbx_call_recording_megabyte_limit,
							group_name,
							group_id
					              FROM
					  	        yap.view___pbx_detail
						      WHERE
						        pbx_id != 1;`)

		// Error
		if err != nil {
			panic(err)

		}

		for pbxAllSQL.Next() {

			err = pbxAllSQL.Scan(
				&pbxName,
				&pbxID,
				&pbxDateAdded,
				&pbxActive,
				&pbxSIPEndpointLimit,
				&pbxSIPTrunkLimit,
				&pbxPhoneNumberLimit,
				&pbxCDRLimit,
				&pbxVoicemailMegabyteLimit,
				&pbxCallRecordingMegabyteLimit,
				&groupName,
				&groupID,
			)

			// Error
			if err != nil {
				panic(err)
			}
			fmt.Fprintf(w, "        <tr>")
			fmt.Fprintf(w, "          <td>"+pbxName+"</td>")
			fmt.Fprintf(w, "          <td>"+pbxID+"</td>")
			fmt.Fprintf(w, "          <td>"+pbxDateAdded+"</td>")
			if pbxActive == "1" {
				fmt.Fprintf(w, "          <td>YES</td>")
			} else {
				fmt.Fprintf(w, "          <td>NO</td>")
			}
			fmt.Fprintf(w, "          <td>"+pbxSIPEndpointLimit+"</td>")
			fmt.Fprintf(w, "          <td>"+pbxSIPTrunkLimit+"</td>")
			fmt.Fprintf(w, "          <td>"+pbxPhoneNumberLimit+"</td>")
			fmt.Fprintf(w, "          <td>"+pbxCDRLimit+"</td>")
			fmt.Fprintf(w, "          <td>"+pbxVoicemailMegabyteLimit+"</td>")
			fmt.Fprintf(w, "          <td>"+pbxCallRecordingMegabyteLimit+"</td>")
			fmt.Fprintf(w, "          <td>"+groupName+"</td>")
			fmt.Fprintf(w, "          <td>"+groupID+"</td>")
			fmt.Fprintf(w, "        </tr>")
		}
	} else if userTypeID == "200" || userTypeID == "201" {
		pbxWhereUserGroupSQL, err := dbDetail.connection.Query(`SELECT
                                                        pbx_name,
                                                        pbx_id,
                                                        pbx_date_added,
                                                        pbx_active,
                                                        pbx_sip_endpoint_limit,
                                                        pbx_sip_trunk_limit,
                                                        pbx_phone_number_limit,
                                                        pbx_cdr_limit,
                                                        pbx_voicemail_megabyte_limit,
                                                        pbx_call_recording_megabyte_limit
                                                      FROM
                                                        yap.view___pbx_detail
                                                      WHERE
                                                        group_id = ?;`, userGroupID)

		// Error
		if err != nil {
			panic(err)

		}

		for pbxWhereUserGroupSQL.Next() {

			err = pbxWhereUserGroupSQL.Scan(
				&pbxName,
				&pbxID,
				&pbxDateAdded,
				&pbxActive,
				&pbxSIPEndpointLimit,
				&pbxSIPTrunkLimit,
				&pbxPhoneNumberLimit,
				&pbxCDRLimit,
				&pbxVoicemailMegabyteLimit,
				&pbxCallRecordingMegabyteLimit,
			)

			// Error
			if err != nil {
				panic(err)
			}
			fmt.Fprintf(w, "        <tr>")
			fmt.Fprintf(w, "          <td>"+pbxName+"</td>")
			fmt.Fprintf(w, "          <td>"+pbxID+"</td>")
			fmt.Fprintf(w, "          <td>"+pbxDateAdded+"</td>")
			if pbxActive == "1" {
				fmt.Fprintf(w, "          <td>YES</td>")
			} else {
				fmt.Fprintf(w, "          <td>NO</td>")
			}
			fmt.Fprintf(w, "          <td>"+pbxSIPEndpointLimit+"</td>")
			fmt.Fprintf(w, "          <td>"+pbxSIPTrunkLimit+"</td>")
			fmt.Fprintf(w, "          <td>"+pbxPhoneNumberLimit+"</td>")
			fmt.Fprintf(w, "          <td>"+pbxCDRLimit+"</td>")
			fmt.Fprintf(w, "          <td>"+pbxVoicemailMegabyteLimit+"</td>")
			fmt.Fprintf(w, "          <td>"+pbxCallRecordingMegabyteLimit+"</td>")
			fmt.Fprintf(w, "        </tr>")
		}
	} else if userTypeID == "300" || userTypeID == "301" || userTypeID == "302" {
		pbxWhereUserPBXSQL, err := dbDetail.connection.Query(`SELECT
                                                        pbx_date_added,
                                                        pbx_active,
                                                        pbx_sip_endpoint_limit,
                                                        pbx_sip_trunk_limit,
                                                        pbx_phone_number_limit,
                                                        pbx_cdr_limit,
                                                        pbx_voicemail_megabyte_limit,
                                                        pbx_call_recording_megabyte_limit
					            FROM
					  	        yap.view___pbx_detail
						    WHERE
					  	        pbx_id = ?;`, userPBXID)

		// Error
		if err != nil {
			panic(err)

		}

		for pbxWhereUserPBXSQL.Next() {

			err = pbxWhereUserPBXSQL.Scan(
				&pbxDateAdded,
				&pbxActive,
				&pbxSIPEndpointLimit,
				&pbxSIPTrunkLimit,
				&pbxPhoneNumberLimit,
				&pbxCDRLimit,
				&pbxVoicemailMegabyteLimit,
				&pbxCallRecordingMegabyteLimit,
			)

			// Error
			if err != nil {
				panic(err)
			}
			fmt.Fprintf(w, "        <tr>")
			fmt.Fprintf(w, "          <td>"+pbxDateAdded+"</td>")
			fmt.Fprintf(w, "          <td>"+pbxActive+"</td>")
			fmt.Fprintf(w, "          <td>"+pbxSIPEndpointLimit+"</td>")
			fmt.Fprintf(w, "          <td>"+pbxSIPTrunkLimit+"</td>")
			fmt.Fprintf(w, "          <td>"+pbxPhoneNumberLimit+"</td>")
			fmt.Fprintf(w, "          <td>"+pbxCDRLimit+"</td>")
			fmt.Fprintf(w, "          <td>"+pbxVoicemailMegabyteLimit+"</td>")
			fmt.Fprintf(w, "          <td>"+pbxCallRecordingMegabyteLimit+"</td>")
			fmt.Fprintf(w, "        </tr>")
		}
	}
	fmt.Fprintf(w, "      </table>")
	if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
		filterTableJS(w, "pbxResourceSearchPBXName", "pbx-resource-input-pbx-name", "pbx-resource-table", 0)
		filterTableJS(w, "pbxResourceSearchPBXID", "pbx-resource-input-pbx-id", "pbx-resource-table", 1)
		filterTableJS(w, "pbxResourceSearchDate", "pbx-resource-input-date", "pbx-resource-table", 2)
		filterTableJS(w, "pbxResourceSearchActive", "pbx-resource-input-active", "pbx-resource-table", 3)
		if userTypeID == "100" {
			filterTableJS(w, "pbxResourceSearchGroupName", "pbx-resource-input-group-name", "pbx-resource-table", 10)
			filterTableJS(w, "pbxResourceSearchGroupID", "pbx-resource-input-group-id", "pbx-resource-table", 11)
		}
	}
	exportCSVJS(w, "PBXResource", "pbx-resource-table", "YAP_pbx_resource_details", "pbx")
	fmt.Fprintf(w, "    </th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "</table>")
	fmt.Fprintf(w, "</div>")
	if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
		toggleDivJS(w, "togglePBX", "pbx-div")
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

// SIP endpoint page Functions

func sipEndpointList(w http.ResponseWriter, dbDetail databaseFunctionParameter, userTypeID string, userGroupID string, userPBXID string) {

	var (
		sipUsername  string
		sipPassword  string
		codecAllowed string
		dtmfUsed     string
		callGroup    string
		pickupGroup  string
		registered   string
		pbxName      string
		groupName    string
		groupID      string
	)

	var dbTableCountUserSIPEndpoint databaseFunctionParameter
	dbTableCountUserSIPEndpoint.connection = dbDetail.connection
	dbTableCountUserSIPEndpoint.database = dbDetail.database
	dbTableCountUserSIPEndpoint.table = "view___sip_endpoint_detail"
	dbTableCountUserSIPEndpoint.columnWhere = "sip_username"

	fmt.Fprintf(w, "<table id=\"table\" class=\"table-sip-endpoint\">")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>")
	fmt.Fprintf(w, "      <table id=\"table\" class=\"table-sip-endpoint\">")
	fmt.Fprintf(w, "        <tr>")
	if userTypeID == "100" {
		fmt.Fprintf(w, "          <th>Total SIP Endpoints On YAP</th>")
	} else if userTypeID == "200" || userTypeID == "201" {
		fmt.Fprintf(w, "          <th>Total SIP Endpoints Within the Group</th>")
	} else if userTypeID == "300" || userTypeID == "301" || userTypeID == "302" {
		fmt.Fprintf(w, "          <th>Total SIP Endpoints Within the PBX</th>")
	}
	fmt.Fprintf(w, "        </tr>")
	fmt.Fprintf(w, "        <tr>")
	if userTypeID == "100" {
		dbTableCountUserSIPEndpoint.countMinusOne = false
		fmt.Fprintf(w, "          <td>"+totalTableCount(w, dbTableCountUserSIPEndpoint)+"</td>")
	} else if userTypeID == "200" || userTypeID == "201" {
		var dbTableCountUserSIPEndpointWhere databaseFunctionParameter
		dbTableCountUserSIPEndpointWhere.connection = dbDetail.connection
		dbTableCountUserSIPEndpointWhere.database = dbDetail.database
		dbTableCountUserSIPEndpointWhere.table = "view___sip_endpoint_detail"
		dbTableCountUserSIPEndpointWhere.columnWhere = "group_id"
		dbTableCountUserSIPEndpointWhere.columnWhereValue = userGroupID
		fmt.Fprintf(w, "    <td>"+totalTableCountWhere(w, dbTableCountUserSIPEndpointWhere)+"</td>")
	} else if userTypeID == "300" || userTypeID == "301" || userTypeID == "302" {
		var dbTableCountUserSIPEndpointWhere databaseFunctionParameter
		dbTableCountUserSIPEndpointWhere.connection = dbDetail.connection
		dbTableCountUserSIPEndpointWhere.database = dbDetail.database
		dbTableCountUserSIPEndpointWhere.table = "view___sip_endpoint_detail"
		dbTableCountUserSIPEndpointWhere.columnWhere = "pbx_id"
		dbTableCountUserSIPEndpointWhere.columnWhereValue = userPBXID
		fmt.Fprintf(w, "    <td>"+totalTableCountWhere(w, dbTableCountUserSIPEndpointWhere)+"</td>")
	}
	fmt.Fprintf(w, "        </tr>")
	fmt.Fprintf(w, "      </table>")
	fmt.Fprintf(w, "    </th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th><button onclick=\"toggleSIPEndpoint() \"class=\"button-general button-sip-endpoint\">&nbsp Show/Hide SIP Endpoint(s) &nbsp</button></th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "</table>")

	fmt.Fprintf(w, "<div id=\"sip-endpoint-div\" style=\"display:none\">")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<table id=\"table\" class=\"table-sip-endpoint\">")
	fmt.Fprintf(w, "  <tr>")
	if userTypeID == "100" {
		fmt.Fprintf(w, "    <th class=\"table-title\";>All SIP Endpoint Details on the Server:</th>")
	} else if userTypeID == "200" || userTypeID == "201" {
		fmt.Fprintf(w, "    <th class=\"table-title\";>All SIP Endpoint Details Within the Group:</th>")
	} else if userTypeID == "300" || userTypeID == "301" || userTypeID == "302" {
		fmt.Fprintf(w, "    <th class=\"table-title\";>All SIP Endpoint Details Within the PBX:</th>")
	}
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>")
	fmt.Fprintf(w, "    <br>")
	fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
	inputTableHTML(w, "sipEndpointSearchSIPUsername", "sip-endpoint-input-sip-username", "SIP Username/PBX ID")
	fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
	inputTableHTML(w, "sipEndpointSearchCodec", "sip-endpoint-input-codec", "Codec(s) Allowed")
	fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
	inputTableHTML(w, "sipEndpointSearchDTMF", "sip-endpoint-input-dtmf", "DTMF Method Used")
	fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
	inputTableHTML(w, "sipEndpointSearchCallGroup", "sip-endpoint-input-call-group", "Call Group")
	fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
	fmt.Fprintf(w, "    <br>")
	fmt.Fprintf(w, "    <br>")
	fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
	inputTableHTML(w, "sipEndpointSearchPickupGroup", "sip-endpoint-input-pickup-group", "Pickup Group")
	fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
	if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
		inputTableHTML(w, "sipEndpointSearchPBXName", "sip-endpoint-input-pbx-name", "PBX Name")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
	}
	if userTypeID == "100" {
		inputTableHTML(w, "sipEndpointSearchGroupName", "sip-endpoint-input-group-name", "Group Name")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
		inputTableHTML(w, "sipEndpointSearchGroupID", "sip-endpoint-input-group-id", "Group ID")
		fmt.Fprintf(w, "    &nbsp &nbsp &nbsp")
	}
	fmt.Fprintf(w, "    <br>")
	fmt.Fprintf(w, "    <br>")
	fmt.Fprintf(w, "    </th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>")
	fmt.Fprintf(w, "      <table id=\"sip-endpoint-table\" class=\"table-sip-endpoint\">")
	fmt.Fprintf(w, "        <tr>")
	fmt.Fprintf(w, "          <th>SIP Username</th>")
	fmt.Fprintf(w, "          <th>SIP Password</th>")
	fmt.Fprintf(w, "          <th>Codec(s)<br>Allowed</th>")
	fmt.Fprintf(w, "          <th>DTMF<br>Method<br>Used</th>")
	fmt.Fprintf(w, "          <th>Call<br>Group</th>")
	fmt.Fprintf(w, "          <th>Pickup<br>Group</th>")
	fmt.Fprintf(w, "          <th>Registered</th>")
	if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
		fmt.Fprintf(w, "          <th>PBX Name</th>")
	}
	if userTypeID == "100" {
		fmt.Fprintf(w, "          <th>Group Name</th>")
		fmt.Fprintf(w, "          <th>Group ID</th>")
	}
	fmt.Fprintf(w, "        </tr>")
	if userTypeID == "100" {
		sipEndpointAllSQL, err := dbDetail.connection.Query(`SELECT
							sip_username,
							sip_password,
							codec_allowed,
					                dtmf_mode,
					                named_call_group,
					                named_pickup_group,
					                registered,
					                pbx_name,
					                group_name,
					                group_id
					              FROM
					  	        yap.view___sip_endpoint_detail`)

		// Error
		if err != nil {
			panic(err)

		}

		for sipEndpointAllSQL.Next() {

			err = sipEndpointAllSQL.Scan(
				&sipUsername,
				&sipPassword,
				&codecAllowed,
				&dtmfUsed,
				&callGroup,
				&pickupGroup,
				&registered,
				&pbxName,
				&groupName,
				&groupID,
			)

			// Error
			if err != nil {
				panic(err)
			}
			fmt.Fprintf(w, "        <tr>")
			fmt.Fprintf(w, "          <td>"+sipUsername)
			copyButtonJS(w, sipUsername)
			fmt.Fprintf(w, "	  </td>")
			fmt.Fprintf(w, "          <td>"+sipPassword)
			copyButtonJS(w, sipPassword)
			fmt.Fprintf(w, "          </td>")
			fmt.Fprintf(w, "          <td>"+codecAllowed+"</td>")
			fmt.Fprintf(w, "          <td>"+dtmfUsed+"</td>")
			fmt.Fprintf(w, "          <td>"+callGroup+"</td>")
			fmt.Fprintf(w, "          <td>"+pickupGroup+"</td>")
			if registered == "1" {
				fmt.Fprintf(w, "          <td>&#128994</td>")
			} else {
				fmt.Fprintf(w, "          <td>&#128308</td>")
			}
			fmt.Fprintf(w, "          <td>"+pbxName+"</td>")
			fmt.Fprintf(w, "          <td>"+groupName+"</td>")
			fmt.Fprintf(w, "          <td>"+groupID+"</td>")
			fmt.Fprintf(w, "        </tr>")
		}
	} else if userTypeID == "200" || userTypeID == "201" {
		sipEndpointWhereGroupSQL, err := dbDetail.connection.Query(`SELECT
                                                        sip_username,
                                                        sip_password,
                                                        codec_allowed,
                                                        dtmf_mode,
                                                        named_call_group,
                                                        named_pickup_group,
                                                        registered,
                                                        pbx_name
                                                      FROM
                                                        yap.view___sip_endpoint_detail
						      WHERE
                                                        group_id = ?;`, userGroupID)

		// Error
		if err != nil {
			panic(err)

		}

		for sipEndpointWhereGroupSQL.Next() {

			err = sipEndpointWhereGroupSQL.Scan(
				&sipUsername,
				&sipPassword,
				&codecAllowed,
				&dtmfUsed,
				&callGroup,
				&pickupGroup,
				&registered,
				&pbxName,
			)

			// Error
			if err != nil {
				panic(err)
			}
			fmt.Fprintf(w, "        <tr>")
			fmt.Fprintf(w, "          <td>"+sipUsername)
			copyButtonJS(w, sipUsername)
			fmt.Fprintf(w, "          </td>")
			fmt.Fprintf(w, "          <td>"+sipPassword)
			copyButtonJS(w, sipPassword)
			fmt.Fprintf(w, "          </td>")
			fmt.Fprintf(w, "          <td>"+codecAllowed+"</td>")
			fmt.Fprintf(w, "          <td>"+dtmfUsed+"</td>")
			fmt.Fprintf(w, "          <td>"+callGroup+"</td>")
			fmt.Fprintf(w, "          <td>"+pickupGroup+"</td>")
			if registered == "1" {
				fmt.Fprintf(w, "          <td>&#128994</td>")
			} else {
				fmt.Fprintf(w, "          <td>&#128308</td>")
			}
			fmt.Fprintf(w, "          <td>"+pbxName+"</td>")
			fmt.Fprintf(w, "        </tr>")
		}

	} else if userTypeID == "300" || userTypeID == "301" || userTypeID == "302" {
		sipEndpointWherePBXSQL, err := dbDetail.connection.Query(`SELECT
                                                        sip_username,
                                                        sip_password,
                                                        codec_allowed,
                                                        dtmf_mode,
                                                        named_call_group,
                                                        named_pickup_group,
                                                        registered
                                                      FROM
                                                        yap.view___sip_endpoint_detail
                                                      WHERE
                                                        pbx_id = ?;`, userPBXID)

		// Error
		if err != nil {
			panic(err)

		}

		for sipEndpointWherePBXSQL.Next() {

			err = sipEndpointWherePBXSQL.Scan(
				&sipUsername,
				&sipPassword,
				&codecAllowed,
				&dtmfUsed,
				&callGroup,
				&pickupGroup,
				&registered,
			)

			// Error
			if err != nil {
				panic(err)
			}
			fmt.Fprintf(w, "        <tr>")
			fmt.Fprintf(w, "          <td>"+sipUsername)
			copyButtonJS(w, sipUsername)
			fmt.Fprintf(w, "          </td>")
			fmt.Fprintf(w, "          <td>"+sipPassword)
			copyButtonJS(w, sipPassword)
			fmt.Fprintf(w, "          </td>")
			fmt.Fprintf(w, "          <td>"+codecAllowed+"</td>")
			fmt.Fprintf(w, "          <td>"+dtmfUsed+"</td>")
			fmt.Fprintf(w, "          <td>"+callGroup+"</td>")
			fmt.Fprintf(w, "          <td>"+pickupGroup+"</td>")
			if registered == "1" {
				fmt.Fprintf(w, "          <td>&#128994</td>")
			} else {
				fmt.Fprintf(w, "          <td>&#128308</td>")
			}
			fmt.Fprintf(w, "        </tr>")
		}

	}

	fmt.Fprintf(w, "      </table>")
	filterTableJS(w, "sipEndpointSearchSIPUsername", "sip-endpoint-input-sip-username", "sip-endpoint-table", 0)
	filterTableJS(w, "sipEndpointSearchCodec", "sip-endpoint-input-codec", "sip-endpoint-table", 2)
	filterTableJS(w, "sipEndpointSearchDTMF", "sip-endpoint-input-dtmf", "sip-endpoint-table", 3)
	filterTableJS(w, "sipEndpointSearchCallGroup", "sip-endpoint-input-call-group", "sip-endpoint-table", 4)
	filterTableJS(w, "sipEndpointSearchPickupGroup", "sip-endpoint-input-pickup-group", "sip-endpoint-table", 5)
	if userTypeID == "100" || userTypeID == "200" || userTypeID == "201" {
		filterTableJS(w, "sipEndpointSearchPBXName", "sip-endpoint-input-pbx-name", "sip-endpoint-table", 7)
	}
	if userTypeID == "100" {
		filterTableJS(w, "sipEndpointSearchGroupName", "sip-endpoint-input-group-name", "sip-endpoint-table", 8)
		filterTableJS(w, "sipEndpointSearchGroupID", "sip-endpoint-input-group-id", "sip-endpoint-table", 9)
	}
	fmt.Fprintf(w, "    </th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "</table>")
	fmt.Fprintf(w, "</div>")
	toggleDivJS(w, "toggleSIPEndpoint", "sip-endpoint-div")
}

//----------------------------------------------------------------------------------------------------

// SIP trunk page functions

//----------------------------------------------------------------------------------------------------

// Phone number page functions

//----------------------------------------------------------------------------------------------------

// CDR page functions

//----------------------------------------------------------------------------------------------------

// Voicemail page functions

//----------------------------------------------------------------------------------------------------

// Call recording page functions

//----------------------------------------------------------------------------------------------------

// MoH / AA music page functions

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
	}

	startHTML := csvcell.FileData(dirHTML, fileStartHTML)
	endHTML := csvcell.FileData(dirHTML, fileEndHTML)

	// Home Page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

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
		userGroupID := userAccountData(dbDetail, "group_id")
		userGroupName := userAccountData(dbDetail, "group_name")
		userPBXID := userAccountData(dbDetail, "pbx_id")
		userPBXName := userAccountData(dbDetail, "pbx_name")

		if userTypeID == "" {
			errorBox(w, "email_error", "header-main-menu")
		} else {
			if userTypeID == "100" {
				header(w, "Main Menu<br>YAP Admin Account", "")
				mainMenuUserInformation(w, dbDetail, userTypeID)
				fmt.Fprintf(w, "<br>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonOne := mainMenuParameter{writeHTTP: w, buttonName: "All Group & PBX<br>User Account(s)<br>&#128100", hyperlink: "/user-account", headerCSS: "header-user-account", buttonCSS: "button-user-account"}
				mainMenuButton(mainMenuButtonOne)
				mainMenuButtonTwo := mainMenuParameter{writeHTTP: w, buttonName: "All<br>Groups<br>&#128101", hyperlink: "/group", headerCSS: "header-group", buttonCSS: "button-group"}
				mainMenuButton(mainMenuButtonTwo)
				mainMenuButtonThree := mainMenuParameter{writeHTTP: w, buttonName: "All<br>PBXs<br>&#128222", hyperlink: "/pbx", headerCSS: "header-pbx", buttonCSS: "button-pbx"}
				mainMenuButton(mainMenuButtonThree)
				mainMenuButtonFour := mainMenuParameter{writeHTTP: w, buttonName: "All PBX SIP<br>Endpoints<br>&#128241", hyperlink: "/sip-endpoint", headerCSS: "header-sip-endpoint", buttonCSS: "button-sip-endpoint"}
				mainMenuButton(mainMenuButtonFour)
				fmt.Fprintf(w, "</div>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonFive := mainMenuParameter{writeHTTP: w, buttonName: "All PBX SIP<br>Trunks<br>&#8596", hyperlink: "/sip-trunk", headerCSS: "header-sip-trunk", buttonCSS: "button-sip-trunk"}
				mainMenuButton(mainMenuButtonFive)
				mainMenuButtonSix := mainMenuParameter{writeHTTP: w, buttonName: "All PBX Phone<br>Numbers<br>&#128290", hyperlink: "/phone-number", headerCSS: "header-phone-number", buttonCSS: "button-phone-number"}
				mainMenuButton(mainMenuButtonSix)
				mainMenuButtonSeven := mainMenuParameter{writeHTTP: w, buttonName: "All PBX<br>CDRs<br>&#128202", hyperlink: "/cdr", headerCSS: "header-cdr", buttonCSS: "button-cdr"}
				mainMenuButton(mainMenuButtonSeven)
				mainMenuButtonEight := mainMenuParameter{writeHTTP: w, buttonName: "All PBX<br>Voicemails<br>&#127897", hyperlink: "/voicemail", headerCSS: "header-voicemail", buttonCSS: "button-voicemail"}
				mainMenuButton(mainMenuButtonEight)
				fmt.Fprintf(w, "</div>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonNine := mainMenuParameter{writeHTTP: w, buttonName: "All PBX Call<br>Recordings<br>&#128252", hyperlink: "/call-recording", headerCSS: "header-call-recording", buttonCSS: "button-call-recording"}
				mainMenuButton(mainMenuButtonNine)
				mainMenuButtonTen := mainMenuParameter{writeHTTP: w, buttonName: "All PBX MoH &<br>AA Music<br>&#127925", hyperlink: "/music", headerCSS: "header-music", buttonCSS: "button-music"}
				mainMenuButton(mainMenuButtonTen)
				mainMenuButtonEleven := mainMenuParameter{writeHTTP: w, buttonName: "All Server<br>Logs<br>&#128195", hyperlink: "/server-log", headerCSS: "header-server-log", buttonCSS: "button-server-log"}
				mainMenuButton(mainMenuButtonEleven)
				mainMenuButtonTweleve := mainMenuParameter{writeHTTP: w, buttonName: "YAP Server<br>Information<br>&#128421", hyperlink: "/server-information", headerCSS: "header-server-information", buttonCSS: "button-server-information"}
				mainMenuButton(mainMenuButtonTweleve)
				fmt.Fprintf(w, "</div>")
				footer(w, "", "")
			} else if userTypeID == "200" {
				header(w, "Main Menu<br>"+userGroupName+"<br>[Group ID: "+userGroupID+"]", "")
				mainMenuUserInformation(w, dbDetail, userTypeID)
				fmt.Fprintf(w, "<br>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonOne := mainMenuParameter{writeHTTP: w, buttonName: "Group & PBX<br>User Account(s)<br>&#128100", hyperlink: "/user-account", headerCSS: "header-user-account", buttonCSS: "button-user-account"}
				mainMenuButton(mainMenuButtonOne)
				mainMenuButtonTwo := mainMenuParameter{writeHTTP: w, buttonName: "Group<br>Information<br>&#128101", hyperlink: "/group", headerCSS: "header-group", buttonCSS: "button-group"}
				mainMenuButton(mainMenuButtonTwo)
				mainMenuButtonThree := mainMenuParameter{writeHTTP: w, buttonName: "PBX(s) Within<br>The Group<br>&#128222", hyperlink: "/pbx", headerCSS: "header-pbx", buttonCSS: "button-pbx"}
				mainMenuButton(mainMenuButtonThree)
				mainMenuButtonFour := mainMenuParameter{writeHTTP: w, buttonName: "PBX SIP<br>Endpoint(s)<br>&#128241", hyperlink: "/sip-endpoint", headerCSS: "header-sip-endpoint", buttonCSS: "button-sip-endpoint"}
				mainMenuButton(mainMenuButtonFour)
				fmt.Fprintf(w, "</div>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonFive := mainMenuParameter{writeHTTP: w, buttonName: "PBX SIP<br>Trunk(s)<br>&#8596", hyperlink: "/sip-trunk", headerCSS: "header-sip-trunk", buttonCSS: "button-sip-trunk"}
				mainMenuButton(mainMenuButtonFive)
				mainMenuButtonSix := mainMenuParameter{writeHTTP: w, buttonName: "PBX Phone<br>Number(s)<br>&#128290", hyperlink: "/phone-number", headerCSS: "header-phone-number", buttonCSS: "button-phone-number"}
				mainMenuButton(mainMenuButtonSix)
				mainMenuButtonSeven := mainMenuParameter{writeHTTP: w, buttonName: "PBX<br>CDR(s)<br>&#128202", hyperlink: "/cdr", headerCSS: "header-cdr", buttonCSS: "button-cdr"}
				mainMenuButton(mainMenuButtonSeven)
				mainMenuButtonEight := mainMenuParameter{writeHTTP: w, buttonName: "PBX<br>Voicemail(s)<br>&#127897", hyperlink: "/voicemail", headerCSS: "header-voicemail", buttonCSS: "button-voicemail"}
				mainMenuButton(mainMenuButtonEight)
				fmt.Fprintf(w, "</div>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonNine := mainMenuParameter{writeHTTP: w, buttonName: "PBX Call<br>Recording(s)<br>&#128252", hyperlink: "/call-recording", headerCSS: "header-call-recording", buttonCSS: "button-call-recording"}
				mainMenuButton(mainMenuButtonNine)
				mainMenuButtonTen := mainMenuParameter{writeHTTP: w, buttonName: "PBX MoH & AA<br>Music<br>&#127925", hyperlink: "/music", headerCSS: "header-music", buttonCSS: "button-music"}
				mainMenuButton(mainMenuButtonTen)
				mainMenuButtonEleven := mainMenuParameter{writeHTTP: w, buttonName: "Group & PBX<br>Server Log<br>&#128195", hyperlink: "/server-log", headerCSS: "header-server-log", buttonCSS: "button-server-log"}
				mainMenuButton(mainMenuButtonEleven)
				fmt.Fprintf(w, "</div>")
				footer(w, "", "")
			} else if userTypeID == "201" {
				header(w, "Main Menu<br>"+userGroupName+"<br>[Group ID: "+userGroupID+"]", "")
				mainMenuUserInformation(w, dbDetail, userTypeID)
				fmt.Fprintf(w, "<br>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonOne := mainMenuParameter{writeHTTP: w, buttonName: "Own & PBX<br>User Account(s)<br>&#128100", hyperlink: "/user-account", headerCSS: "header-user-account", buttonCSS: "button-user-account"}
				mainMenuButton(mainMenuButtonOne)
				mainMenuButtonTwo := mainMenuParameter{writeHTTP: w, buttonName: "Group<br>Information<br>&#128101", hyperlink: "/group", headerCSS: "header-group", buttonCSS: "button-group"}
				mainMenuButton(mainMenuButtonTwo)
				mainMenuButtonThree := mainMenuParameter{writeHTTP: w, buttonName: "PBX(s) Within<br>The Group<br>&#128222", hyperlink: "/pbx", headerCSS: "header-pbx", buttonCSS: "button-pbx"}
				mainMenuButton(mainMenuButtonThree)
				mainMenuButtonFour := mainMenuParameter{writeHTTP: w, buttonName: "PBX SIP<br>Endpoint(s)<br>&#128241", hyperlink: "/sip-endpoint", headerCSS: "header-sip-endpoint", buttonCSS: "button-sip-endpoint"}
				mainMenuButton(mainMenuButtonFour)
				fmt.Fprintf(w, "</div>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonFive := mainMenuParameter{writeHTTP: w, buttonName: "PBX SIP<br>Trunk(s)<br>&#8596", hyperlink: "/sip-trunk", headerCSS: "header-sip-trunk", buttonCSS: "button-sip-trunk"}
				mainMenuButton(mainMenuButtonFive)
				mainMenuButtonSix := mainMenuParameter{writeHTTP: w, buttonName: "PBX Phone<br>Number(s)<br>&#128290", hyperlink: "/phone-number", headerCSS: "header-phone-number", buttonCSS: "button-phone-number"}
				mainMenuButton(mainMenuButtonSix)
				mainMenuButtonSeven := mainMenuParameter{writeHTTP: w, buttonName: "PBX<br>CDR(s)<br>&#128202", hyperlink: "/cdr", headerCSS: "header-cdr", buttonCSS: "button-cdr"}
				mainMenuButton(mainMenuButtonSeven)
				mainMenuButtonEight := mainMenuParameter{writeHTTP: w, buttonName: "PBX<br>Voicemail(s)<br>&#127897", hyperlink: "/voicemail", headerCSS: "header-voicemail", buttonCSS: "button-voicemail"}
				mainMenuButton(mainMenuButtonEight)
				fmt.Fprintf(w, "</div>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonNine := mainMenuParameter{writeHTTP: w, buttonName: "PBX Call<br>Recording(s)<br>&#128252", hyperlink: "/call-recording", headerCSS: "header-call-recording", buttonCSS: "button-call-recording"}
				mainMenuButton(mainMenuButtonNine)
				mainMenuButtonTen := mainMenuParameter{writeHTTP: w, buttonName: "PBX MoH & AA<br>Music<br>&#127925", hyperlink: "/music", headerCSS: "header-music", buttonCSS: "button-music"}
				mainMenuButton(mainMenuButtonTen)
				mainMenuButtonEleven := mainMenuParameter{writeHTTP: w, buttonName: "PBX<br>Server Log<br>&#128195", hyperlink: "/server-log", headerCSS: "header-server-log", buttonCSS: "button-server-log"}
				mainMenuButton(mainMenuButtonEleven)
				fmt.Fprintf(w, "</div>")
				footer(w, "", "")
			} else if userTypeID == "300" {
				header(w, "Main Menu<br>"+userPBXName+"<br>[PBX ID: "+userPBXID+"]", "")
				mainMenuUserInformation(w, dbDetail, userTypeID)
				fmt.Fprintf(w, "<br>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonOne := mainMenuParameter{writeHTTP: w, buttonName: "Own & PBX<br>User Account(s)<br>&#128100", hyperlink: "/user-account", headerCSS: "header-user-account", buttonCSS: "button-user-account"}
				mainMenuButton(mainMenuButtonOne)
				mainMenuButtonTwo := mainMenuParameter{writeHTTP: w, buttonName: "PBX<br>Information<br>&#128222", hyperlink: "/pbx", headerCSS: "header-pbx", buttonCSS: "button-pbx"}
				mainMenuButton(mainMenuButtonTwo)
				mainMenuButtonThree := mainMenuParameter{writeHTTP: w, buttonName: "PBX SIP<br>Endpoint(s)<br>&#128241", hyperlink: "/sip-endpoint", headerCSS: "header-sip-endpoint", buttonCSS: "button-sip-endpoint"}
				mainMenuButton(mainMenuButtonThree)
				mainMenuButtonFour := mainMenuParameter{writeHTTP: w, buttonName: "PBX SIP<br>Trunk(s)<br>&#8596", hyperlink: "/sip-trunk", headerCSS: "header-sip-trunk", buttonCSS: "button-sip-trunk"}
				mainMenuButton(mainMenuButtonFour)
				fmt.Fprintf(w, "</div>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonFive := mainMenuParameter{writeHTTP: w, buttonName: "PBX Phone<br>Number(s)<br>&#128290", hyperlink: "/phone-number", headerCSS: "header-phone-number", buttonCSS: "button-phone-number"}
				mainMenuButton(mainMenuButtonFive)
				mainMenuButtonSix := mainMenuParameter{writeHTTP: w, buttonName: "PBX<br>CDR(s)<br>&#128202", hyperlink: "/cdr", headerCSS: "header-cdr", buttonCSS: "button-cdr"}
				mainMenuButton(mainMenuButtonSix)
				mainMenuButtonSeven := mainMenuParameter{writeHTTP: w, buttonName: "PBX<br>Voicemail(s)<br>&#127897", hyperlink: "/voicemail", headerCSS: "header-voicemail", buttonCSS: "button-voicemail"}
				mainMenuButton(mainMenuButtonSeven)
				mainMenuButtonEight := mainMenuParameter{writeHTTP: w, buttonName: "PBX Call<br>Recording(s)<br>&#128252", hyperlink: "/call-recording", headerCSS: "header-call-recording", buttonCSS: "button-call-recording"}
				mainMenuButton(mainMenuButtonEight)
				fmt.Fprintf(w, "</div>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonNine := mainMenuParameter{writeHTTP: w, buttonName: "PBX MoH & AA<br>Music<br>&#127925", hyperlink: "/music", headerCSS: "header-music", buttonCSS: "button-music"}
				mainMenuButton(mainMenuButtonNine)
				mainMenuButtonTen := mainMenuParameter{writeHTTP: w, buttonName: "PBX<br>Server Log<br>&#128195", hyperlink: "/server-log", headerCSS: "header-server-log", buttonCSS: "button-server-log"}
				mainMenuButton(mainMenuButtonTen)
				fmt.Fprintf(w, "</div>")
				footer(w, "", "")
			} else if userTypeID == "301" || userTypeID == "302" {
				header(w, "Main Menu<br>"+userPBXName+"<br>[PBX ID: "+userPBXID+"]", "")
				mainMenuUserInformation(w, dbDetail, userTypeID)
				fmt.Fprintf(w, "<br>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonOne := mainMenuParameter{writeHTTP: w, buttonName: "Own<br>User Account<br>&#128100", hyperlink: "/user-account", headerCSS: "header-user-account", buttonCSS: "button-user-account"}
				mainMenuButton(mainMenuButtonOne)
				mainMenuButtonTwo := mainMenuParameter{writeHTTP: w, buttonName: "PBX<br>Information<br>&#128222", hyperlink: "/pbx", headerCSS: "header-pbx", buttonCSS: "button-pbx"}
				mainMenuButton(mainMenuButtonTwo)
				mainMenuButtonThree := mainMenuParameter{writeHTTP: w, buttonName: "PBX SIP<br>Endpoint(s)<br>&#128241", hyperlink: "/sip-endpoint", headerCSS: "header-sip-endpoint", buttonCSS: "button-sip-endpoint"}
				mainMenuButton(mainMenuButtonThree)
				mainMenuButtonFour := mainMenuParameter{writeHTTP: w, buttonName: "PBX SIP<br>Trunk(s)<br>&#8596", hyperlink: "/sip-trunk", headerCSS: "header-sip-trunk", buttonCSS: "button-sip-trunk"}
				mainMenuButton(mainMenuButtonFour)
				fmt.Fprintf(w, "</div>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonFive := mainMenuParameter{writeHTTP: w, buttonName: "PBX Phone<br>Number(s)<br>&#128290", hyperlink: "/phone-number", headerCSS: "header-phone-number", buttonCSS: "button-phone-number"}
				mainMenuButton(mainMenuButtonFive)
				mainMenuButtonSix := mainMenuParameter{writeHTTP: w, buttonName: "PBX<br>CDR(s)<br>&#128202", hyperlink: "/cdr", headerCSS: "header-cdr", buttonCSS: "button-cdr"}
				mainMenuButton(mainMenuButtonSix)
				mainMenuButtonSeven := mainMenuParameter{writeHTTP: w, buttonName: "PBX<br>Voicemail(s)<br>&#127897", hyperlink: "/voicemail", headerCSS: "header-voicemail", buttonCSS: "button-voicemail"}
				mainMenuButton(mainMenuButtonSeven)
				mainMenuButtonEight := mainMenuParameter{writeHTTP: w, buttonName: "PBX Call<br>Recording(s)<br>&#128252", hyperlink: "/call-recording", headerCSS: "header-call-recording", buttonCSS: "button-call-recording"}
				mainMenuButton(mainMenuButtonEight)
				fmt.Fprintf(w, "</div>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonNine := mainMenuParameter{writeHTTP: w, buttonName: "PBX MoH & AA<br>Music<br>&#127925", hyperlink: "/music", headerCSS: "header-music", buttonCSS: "button-music"}
				mainMenuButton(mainMenuButtonNine)
				fmt.Fprintf(w, "</div>")
				footer(w, "", "")
			} else {
				errorBox(w, "account_type_error", "header-main-menu")
			}
		}
		fmt.Fprintf(w, endHTML)

	})

	// User Account Page
	http.HandleFunc("/user-account", func(w http.ResponseWriter, r *http.Request) {

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
		userGroupID := userAccountData(dbDetail, "group_id")
		userGroupName := userAccountData(dbDetail, "group_name")
		userPBXID := userAccountData(dbDetail, "pbx_id")
		userPBXName := userAccountData(dbDetail, "pbx_name")

		if userTypeID == "" {
			errorBox(w, "email_error", "header-user-account")
		} else {
			if userTypeID == "100" {
				header(w, "All User Accounts on the Server<br>YAP Admin Account", "header-user-account")
				userAccountList(w, dbDetail, userTypeID)
				footer(w, "header-user-account", "button-user-account")
			} else if userTypeID == "200" {
				header(w, "All User Accounts Within the Group<br>"+userGroupName+"<br>[Group ID: "+userGroupID+"]", "header-user-account")
				userAccountList(w, dbDetail, userTypeID)
				footer(w, "header-user-account", "button-user-account")
			} else if userTypeID == "201" {
				header(w, "All PBX User Accounts Within the Group<br>"+userGroupName+"<br>[Group ID: "+userGroupID+"]", "header-user-account")
				userAccountList(w, dbDetail, userTypeID)
				footer(w, "header-user-account", "button-user-account")
			} else if userTypeID == "300" {
				header(w, "All User Accounts Within the PBX<br>"+userPBXName+"<br>[PBX ID: "+userPBXID+"]", "header-user-account")
				userAccountList(w, dbDetail, userTypeID)
				footer(w, "header-user-account", "button-user-account")
			} else if userTypeID == "301" {
				header(w, "Own User Account for PBX<br>"+userPBXName+"<br>[PBX ID: "+userPBXID+"]", "header-user-account")
				userAccountList(w, dbDetail, userTypeID)
				footer(w, "header-user-account", "button-user-account")
			} else if userTypeID == "302" {
				header(w, "Own Read Only User Account for PBX<br>"+userPBXName+"<br>[PBX ID: "+userPBXID+"]", "header-user-account")
				userAccountList(w, dbDetail, userTypeID)
				footer(w, "header-user-account", "button-user-account")
			} else {
				errorBox(w, "account_type_error", "header-user-account")
			}
		}
		fmt.Fprintf(w, endHTML)
	})

	// Group Page
	http.HandleFunc("/group", func(w http.ResponseWriter, r *http.Request) {

		// Open database connection
		dbConnection, err := sql.Open("mysql", dbUsername+":"+dbPassword+"@"+dbTransport+"("+dbAddress+":"+dbPort+")/"+dbName+"?tls="+dbTls)
		defer dbConnection.Close()

		// Error
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, startHTML)

		// Wallpaper
		wallpaper(w, "wallpaper-group")

		// Code to call the emailHeaderHTTP function
		email := emailHeaderHTTP(r)

		var dbDetail databaseFunctionParameter
		dbDetail.connection = dbConnection
		dbDetail.database = dbName
		dbDetail.columnWhereValue = email

		userTypeID := userAccountData(dbDetail, "type_id")
		userGroupID := userAccountData(dbDetail, "group_id")
		userGroupName := userAccountData(dbDetail, "group_name")

		if userTypeID == "" {
			errorBox(w, "email_error", "header-group")
		} else {
			if userTypeID == "100" {
				header(w, "All Groups on the Server<br>YAP Admin Account", "header-group")
				groupList(w, dbDetail, userTypeID, userGroupID)
				footer(w, "header-group", "button-group")
			} else if userTypeID == "200" || userTypeID == "201" {
				header(w, "Own Group Information<br>"+userGroupName+"<br>[Group ID: "+userGroupID+"]", "header-group")
				groupList(w, dbDetail, userTypeID, userGroupID)
				footer(w, "header-group", "button-group")
			} else {
				errorBox(w, "account_type_error", "header-group")
			}
		}
		fmt.Fprintf(w, endHTML)
	})

	// PBX Page

	http.HandleFunc("/pbx", func(w http.ResponseWriter, r *http.Request) {

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
		userGroupID := userAccountData(dbDetail, "group_id")
		userGroupName := userAccountData(dbDetail, "group_name")
		userPBXID := userAccountData(dbDetail, "pbx_id")
		userPBXName := userAccountData(dbDetail, "pbx_name")

		if userTypeID == "" {
			errorBox(w, "email_error", "header-pbx")
		} else {
			if userTypeID == "100" {
				header(w, "All PBXs on the Server<br>YAP Admin Account", "header-pbx")
				pbxList(w, dbDetail, userTypeID, userGroupID, userPBXID)
				footer(w, "header-pbx", "button-pbx")
			} else if userTypeID == "200" || userTypeID == "201" {
				header(w, "All PBXs Within the Group<br>"+userGroupName+"<br>[Group ID: "+userGroupID+"]", "header-pbx")
				pbxList(w, dbDetail, userTypeID, userGroupID, userPBXID)
				footer(w, "header-pbx", "button-pbx")
			} else if userTypeID == "300" || userTypeID == "301" || userTypeID == "302" {
				header(w, "PBX Information<br>"+userPBXName+"<br>[PBX ID: "+userPBXID+"]", "header-pbx")
				pbxList(w, dbDetail, userTypeID, userGroupID, userPBXID)
				footer(w, "header-pbx", "button-pbx")
			} else {
				errorBox(w, "account_type_error", "header-pbx")
			}
		}
		fmt.Fprintf(w, endHTML)
	})

	// SIP Endpoint Page
	http.HandleFunc("/sip-endpoint", func(w http.ResponseWriter, r *http.Request) {

		// Open database connection
		dbConnection, err := sql.Open("mysql", dbUsername+":"+dbPassword+"@"+dbTransport+"("+dbAddress+":"+dbPort+")/"+dbName+"?tls="+dbTls)
		defer dbConnection.Close()

		// Error
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, startHTML)

		// Wallpaper
		wallpaper(w, "wallpaper-sip-endpoint")

		// Code to call the emailHeaderHTTP function
		email := emailHeaderHTTP(r)

		var dbDetail databaseFunctionParameter
		dbDetail.connection = dbConnection
		dbDetail.database = dbName
		dbDetail.columnWhereValue = email

		userTypeID := userAccountData(dbDetail, "type_id")
		userGroupID := userAccountData(dbDetail, "group_id")
		userGroupName := userAccountData(dbDetail, "group_name")
		userPBXID := userAccountData(dbDetail, "pbx_id")
		userPBXName := userAccountData(dbDetail, "pbx_name")

		if userTypeID == "" {
			errorBox(w, "email_error", "header-sip-endpoint")
		} else {
			if userTypeID == "100" {
				header(w, "All SIP Endpoints on the Server<br>YAP Admin Account", "header-sip-endpoint")
				sipEndpointList(w, dbDetail, userTypeID, userGroupID, userPBXID)
				footer(w, "header-sip-endpoint", "button-sip-endpoint")
			} else if userTypeID == "200" || userTypeID == "201" {
				header(w, "All SIP Endpoints Within the Group<br>"+userGroupName+"<br>[Group ID: "+userGroupID+"]", "header-sip-endpoint")
				sipEndpointList(w, dbDetail, userTypeID, userGroupID, userPBXID)
				footer(w, "header-sip-endpoint", "button-sip-endpoint")
			} else if userTypeID == "300" || userTypeID == "301" || userTypeID == "302" {
				header(w, "All SIP Endpoints Within the PBX<br>"+userPBXName+"<br>[PBX ID: "+userPBXID+"]", "header-sip-endpoint")
				sipEndpointList(w, dbDetail, userTypeID, userGroupID, userPBXID)
				footer(w, "header-sip-endpoint", "button-sip-endpoint")
			} else {
				errorBox(w, "account_type_error", "header-sip-endpoint")
			}
		}
	})

	// SIP Trunk Page
	http.HandleFunc("/sip-trunk", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, startHTML)
		header(w, "SIP Trunks", "header-sip-trunk")
		// Wallpaper
		wallpaper(w, "wallpaper-sip-trunk")

		footer(w, "header-sip-trunk", "button-sip-trunk")
		fmt.Fprintf(w, endHTML)
	})

	// Phone Number Page
	http.HandleFunc("/phone-number", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, startHTML)
		header(w, "Phone Numbers", "header-phone-number")
		// Wallpaper
		wallpaper(w, "wallpaper-phone-number")

		footer(w, "header-phone-number", "button-phone-number")
		fmt.Fprintf(w, endHTML)
	})

	// CDR Page
	http.HandleFunc("/cdr", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, startHTML)
		header(w, "CDRs", "header-cdr")
		// Wallpaper
		wallpaper(w, "wallpaper-cdr")

		footer(w, "header-cdr", "button-cdr")
		fmt.Fprintf(w, endHTML)
	})

	// Voicemail Page
	http.HandleFunc("/voicemail", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, startHTML)
		header(w, "Voicemails", "header-voicemail")
		// Wallpaper
		wallpaper(w, "wallpaper-voicemail")

		footer(w, "header-voicemail", "button-voicemail")
		fmt.Fprintf(w, endHTML)
	})

	// Call Recording Page
	http.HandleFunc("/call-recording", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, startHTML)
		// Call Recording Wallpaper
		wallpaper(w, "wallpaper-call-recording")
		header(w, "Call Recordings", "header-call-recording")
		fmt.Fprintf(w, "<table id=\"table\" class=\"table-call-recording\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>From (Caller)</th>")
		fmt.Fprintf(w, "    <th>To (Callie)</th>")
		fmt.Fprintf(w, "    <th>Date (YYYY-MM-DD)<br>Time (HH:MM:SS)</th>")
		fmt.Fprintf(w, "    <th>Listen</th>")
		fmt.Fprintf(w, "    <th>Download</th>")
		fmt.Fprintf(w, "  </tr>")
		randomString := genID()
		audioList := csvcell.DirList("/var/lib/yap/call-recording")
		for _, audioListLoop := range audioList {
			audioListSplit := strings.Split(audioListLoop, "_")
			/*
				check := "vm"
				str := mohListLoop
				if strings.Contains(str, check) {
			*/
			fmt.Fprintf(w, "  <tr>")
			fmt.Fprintf(w, "    <td>"+audioListSplit[0]+"</td>")
			fmt.Fprintf(w, "    <td>"+audioListSplit[1]+"</td>")
			fmt.Fprintf(w, "    <td>"+audioListSplit[2]+"<br>"+strings.Replace(audioListSplit[3], "-", ":", -1)+"</td>")
			fmt.Fprintf(w, "    <td><a href=\"/download/"+randomString+"/"+audioListLoop+"\" style=\"text-decoration: none;\" target=\"_blank\">&#9654</a></td>")
			fmt.Fprintf(w, "    <td><a href=\"/download/"+randomString+"/"+audioListLoop+"\" style=\"text-decoration: none;\" download=\""+audioListLoop+"\">&#11015</a></td>")
			fmt.Fprintf(w, "  </tr>")
			download(randomString, audioListLoop)
		}
		fmt.Fprintf(w, "</table>")
		footer(w, "header-call-recording", "button-call-recording")
		fmt.Fprintf(w, endHTML)
	})

	// MoH/AA (Music) Page
	http.HandleFunc("/music", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, startHTML)
		header(w, "MoH and AA", "header-music")
		// Wallpaper
		wallpaper(w, "wallpaper-music")

		footer(w, "header-music", "button-music")
		fmt.Fprintf(w, endHTML)
	})

	// Server Log Page
	http.HandleFunc("/server-log", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, startHTML)
		header(w, "Server Logs", "header-server-log")
		// Wallpaper
		wallpaper(w, "wallpaper-server-log")

		footer(w, "header-server-log", "button-server-log")
		fmt.Fprintf(w, endHTML)
	})

	// Server Information Page
	http.HandleFunc("/server-information", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, startHTML)
		header(w, "Server Information", "header-server-information")
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
		fmt.Println("YAP is running on port " + socket)
		//Start server on port specified above
		log.Fatal(http.ListenAndServe(socket, nil))
	}
}

// Contributor(s):
// Elliot Michael Keavney
