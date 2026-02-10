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
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
)

// Constant for directory path that contains the files yap-start.html and yap-end.html
const dirHTML string = "/etc/yap/html-css"

// Constant for fileStartHTML file
const fileStartHTML string = "yap-start.html"

// Constant for fileEndHTML file
const fileEndHTML string = "yap-end.html"

// Function to retrive HTTP email header
func emailHeaderHTTP(r *http.Request) (email string) {
	email = r.Header.Get("X-Email")
	return email
}

// Function for error message
func errorBox(w http.ResponseWriter) {
	fmt.Fprintf(w, "<div class=\"error-box\">")
	fmt.Fprintf(w, "  <h1>")
	fmt.Fprintf(w, "    User Account Not Found")
	fmt.Fprintf(w, "    <br>")
	fmt.Fprintf(w, "    <a href=\"/oauth2/sign_out?rd=https://github.com/logout\" class=\"general-button header-button logout-button\">Logout</a>")
	fmt.Fprintf(w, "  </h1>")
	fmt.Fprintf(w, "</div>")
}

// Function for the header
func header(w http.ResponseWriter) {
	fmt.Fprintf(w, "<div class=\"header\">")
	fmt.Fprintf(w, "  <h1>")
	fmt.Fprintf(w, "    <a href=\"/oauth2/sign_out?rd=https://github.com/logout\" class=\"general-button header-button logout-button\">Logout</a>")
	fmt.Fprintf(w, "    <a href=\"/\" class=\"general-button header-button home-button\">Home</a>")
	fmt.Fprintf(w, "    <a href=\"https://github.com/ellwould/yap/blob/main/LICENSE\" target=\"_blank\" class=\"general-button header-button license-button\">License</a>")
	fmt.Fprintf(w, "    <br>")
	fmt.Fprintf(w, "    ✱ YAP (Yet Another PBX) ✱")
	fmt.Fprintf(w, "  </h1>")
	fmt.Fprintf(w, "</div>")
}

// Function for the footer
func footer(w http.ResponseWriter) {
	fmt.Fprintf(w, "<div class=\"footer\">")
	fmt.Fprintf(w, "  <h2>")
	fmt.Fprintf(w, "    <a href=\"https://github.com/ellwould/yap\" target=\"_blank\" class=\"general-button footer-button\">YAP Source Code</a>")
	fmt.Fprintf(w, "    <a href=\"https://ell.today\" target=\"_blank\" class=\"general-button footer-button\">Other Software</a>")
	fmt.Fprintf(w, "  </h2>")
	fmt.Fprintf(w, "</div>")
}

func selectWhere(dbConnection *sql.DB, dbName string, tableName string, column string, columnWhere string, value string) string {
	var selectWhere string
	selectWhereQuery, err := dbConnection.Query("SELECT "+column+" FROM "+dbName+"."+tableName+" WHERE "+columnWhere+" = ?;", value)
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

func totalTableCount(dbConnection *sql.DB, w http.ResponseWriter, dbName string, tableName string, minus1 bool) string {
	if minus1 == true {
		var countMinus1 string
		countMinus1Query, err := dbConnection.Query("SELECT COUNT(*) -1 FROM " + dbName + "." + tableName)
		if err != nil {
			panic(err)
		}
		for countMinus1Query.Next() {
			err = countMinus1Query.Scan(&countMinus1)
			// Error
			if err != nil {
				panic(err)
			}
		}
		return countMinus1
	} else {
		var count string
		countQuery, err := dbConnection.Query("SELECT COUNT(*) FROM " + dbName + "." + tableName)
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

func totalTableCountWhere(dbConnection *sql.DB, w http.ResponseWriter, dbName string, tableName string, whereColumn string, value string) string {
	var count string
	countQuery, err := dbConnection.Query("SELECT COUNT(*) FROM "+dbName+"."+tableName+" WHERE "+whereColumn+" =?", value)
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

func userAccountTypeId(dbConnection *sql.DB, dbName string, email string) string {

	return selectWhere(dbConnection, dbName, "view___account_detail", "user_account_type_id", "user_account_email", email)
}

func yapAccount(dbConnection *sql.DB, w http.ResponseWriter, dbName string, email string) {

	fmt.Fprintf(w, "</table>")
	fmt.Fprintf(w, "<table id=\"table\">")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>&nbsp Total Groups &nbsp</th>")
	fmt.Fprintf(w, "    <th>&nbsp Total PBXs &nbsp</th>")
	fmt.Fprintf(w, "    <th>&nbsp Total SIP Extensions &nbsp</th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <td>&nbsp"+totalTableCount(dbConnection, w, dbName, "group", true)+"&nbsp</td>")
	fmt.Fprintf(w, "    <td>&nbsp"+totalTableCount(dbConnection, w, dbName, "pbx", true)+"&nbsp</td>")
	fmt.Fprintf(w, "    <td>&nbsp"+totalTableCount(dbConnection, w, dbName, "ps_endpoints", false)+"&nbsp</td>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "</table>")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<table id=\"table\">")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>Total YAP<br>Admin<br>Accounts<br>(Type ID: 100)</th>")
	fmt.Fprintf(w, "    <th>Total YAP<br>Read Only<br>Accounts<br>(Type ID: 101)</th>")
	fmt.Fprintf(w, "    <th>Total Group<br>Admin<br>Accounts<br>(Type ID: 200)</th>")
	fmt.Fprintf(w, "    <th>Total Group<br>Read Only<br>Accounts<br>(Type ID: 201)</th>")
	fmt.Fprintf(w, "    <th>Total PBX<br>Admin<br>Accounts<br>(Type ID: 300)</th>")
	fmt.Fprintf(w, "    <th>Total PBX<br>Read Only<br>Accounts<br>(Type ID: 301)</th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <td>&nbsp"+totalTableCountWhere(dbConnection, w, dbName, "user_account", "user_account_type_id", "100")+"&nbsp</td>")
	fmt.Fprintf(w, "    <td>&nbsp"+totalTableCountWhere(dbConnection, w, dbName, "user_account", "user_account_type_id", "101")+"&nbsp</td>")
	fmt.Fprintf(w, "    <td>&nbsp"+totalTableCountWhere(dbConnection, w, dbName, "user_account", "user_account_type_id", "200")+"&nbsp</td>")
	fmt.Fprintf(w, "    <td>&nbsp"+totalTableCountWhere(dbConnection, w, dbName, "user_account", "user_account_type_id", "201")+"&nbsp</td>")
	fmt.Fprintf(w, "    <td>&nbsp"+totalTableCountWhere(dbConnection, w, dbName, "user_account", "user_account_type_id", "300")+"&nbsp</td>")
	fmt.Fprintf(w, "    <td>&nbsp"+totalTableCountWhere(dbConnection, w, dbName, "user_account", "user_account_type_id", "301")+"&nbsp</td>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "</table>")

}

func groupAccount(dbConnection *sql.DB, w http.ResponseWriter, dbName string, email string) {

	result, err := dbConnection.Query(`SELECT
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
					   FROM yap.view___account_detail
					   WHERE user_account_email = ?;`, email)

	// Error
	if err != nil {
		panic(err)
	}

	for result.Next() {
		var (
			groupName                     string
			groupId                       string
			groupSiteAddressLine1         string
			groupSiteAddressLine2         string
			groupSiteCityTownCillage      string
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
			pbxId                         string
		)

		err = result.Scan(
			&groupName,
			&groupId,
			&groupSiteAddressLine1,
			&groupSiteAddressLine2,
			&groupSiteCityTownCillage,
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
			&pbxId,
		)

		// Error
		if err != nil {
			panic("Error in userInfomation function")
		}

		fmt.Fprintf(w, "<table id=\"table\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>&nbsp Group Name and ID &nbsp</th>")
		fmt.Fprintf(w, "    <th>&nbsp Total PBX(s) in Group &nbsp</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td>&nbsp Group Name: "+groupName+"&nbsp<br><br>Group ID: "+groupId+"</td>")
		fmt.Fprintf(w, "    <td>&nbsp"+totalTableCountWhere(dbConnection, w, dbName, "pbx", "id", pbxId)+"&nbsp</td>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<table id=\"table\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>&nbsp Group Site Address &nbsp</th>")
		fmt.Fprintf(w, "    <th>&nbsp Group Site Email &nbsp</th>")
		fmt.Fprintf(w, "    <th>&nbsp Group Site Phone Number &nbsp</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td style=\"text-align: left;\">"+groupSiteAddressLine1+"&nbsp<br>"+groupSiteAddressLine2+"<br>"+groupSiteCityTownCillage+"<br>"+groupSiteCountyStateRegion+"<br><br>"+groupSitePostcodeZipCode+"<br><br>"+groupSiteCountry+"</td>")
		fmt.Fprintf(w, "    <td>&nbsp<a href=\"mailto:"+groupSiteContactEmail+"\">"+groupSiteContactEmail+"</a>&nbsp<br><br>(Click to email)</td>")
		fmt.Fprintf(w, "    <td>&nbsp<a href=\"tel:"+groupSiteContactNumber+"\">"+groupSiteContactNumber+"</a>&nbsp<br><br>(Click to call)</td>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<table id=\"table\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>&nbsp Group Invoice Address &nbsp</th>")
		fmt.Fprintf(w, "    <th>&nbsp Group Invoice Email &nbsp</th>")
		fmt.Fprintf(w, "    <th>&nbsp Group Invoice Phone Number &nbsp</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td style=\"text-align: left;\">"+groupInvoiceAddressLine1+"&nbsp<br>"+groupInvoiceAddressLine2+"<br>"+groupInvoiceCityTownVillage+"<br>"+groupInvoiceCountyStateRegion+"<br><br>"+groupInvoicePostcodeZipCode+"<br><br>"+groupInvoiceCountry+"</td>")
		fmt.Fprintf(w, "    <td>&nbsp<a href=\"mailto:"+groupInvoiceContactEmail+"\">"+groupInvoiceContactEmail+"</a>&nbsp<br><br>(Click to email)</td>")
		fmt.Fprintf(w, "    <td>&nbsp<a href=\"tel:"+groupInvoiceContactNumber+"\">"+groupInvoiceContactNumber+"</a>&nbsp<br><br>(Click to call)</td>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
	}
}

func pbxAccount(dbConnection *sql.DB, w http.ResponseWriter, dbName string, email string) {

	result, err := dbConnection.Query(`SELECT
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
					   FROM yap.view___account_detail
					   WHERE user_account_email = ?;`, email)

	// Error
	if err != nil {
		panic(err)
	}

	for result.Next() {

		var (
			pbxName                     string
			pbxId                       string
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
			&pbxId,
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

		fmt.Fprintf(w, "<table id=\"table\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>&nbsp PBX Name and ID &nbsp</th>")
		fmt.Fprintf(w, "    <th>&nbsp Total SIP Extensions in PBX &nbsp</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td>&nbsp PBX Name: "+pbxName+"&nbsp<br><br>PBX ID: "+pbxId+"</td>")
		fmt.Fprintf(w, "    <td>&nbsp"+totalTableCountWhere(dbConnection, w, dbName, "ps_endpoints", "pbx_id", pbxId)+"&nbsp</td>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<table id=\"table\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>&nbsp PBX Site Address &nbsp</th>")
		fmt.Fprintf(w, "    <th>&nbsp PBX Site Email &nbsp</th>")
		fmt.Fprintf(w, "    <th>&nbsp PBX Site Phone Number &nbsp</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td style=\"text-align: left;\">"+pbxSiteAddressLine1+"&nbsp<br>"+pbxSiteAddressLine2+"<br>"+pbxSiteCityTownVillage+"<br>"+pbxSiteCountyStateRegion+"<br><br>"+pbxSitePostcodeZipCode+"<br><br>"+pbxSiteCountry+"</td>")
		fmt.Fprintf(w, "    <td>&nbsp<a href=\"mailto:"+pbxSiteContactEmail+"\">"+pbxSiteContactEmail+"</a>&nbsp<br><br>(Click to email)</td>")
		fmt.Fprintf(w, "    <td>&nbsp<a href=\"tel:"+pbxSiteContactNumber+"\">"+pbxSiteContactNumber+"</a>&nbsp<br><br>(Click to call)</td>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<table id=\"table\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>&nbsp PBX Invoice Address &nbsp</th>")
		fmt.Fprintf(w, "    <th>&nbsp PBX Invoice Email &nbsp</th>")
		fmt.Fprintf(w, "    <th>&nbsp PBX Invoice Phone Number &nbsp</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td style=\"text-align: left;\">"+pbxInvoiceAddressLine1+"&nbsp<br>"+pbxInvoiceAddressLine2+"<br>"+pbxInvoiceCityTownVillage+"<br>"+pbxInvoiceCountyStateRegion+"<br><br>"+pbxInvoicePostcodeZipCode+"<br><br>"+pbxInvoiceCountry+"</td>")
		fmt.Fprintf(w, "    <td>&nbsp<a href=\"mailto:"+pbxInvoiceContactEmail+"\">"+pbxInvoiceContactEmail+"</a>&nbsp<br><br>(Click to email)</td>")
		fmt.Fprintf(w, "    <td>&nbsp<a href=\"tel:"+pbxInvoiceContactNumber+"\">"+pbxInvoiceContactNumber+"</a>&nbsp<br><br>(Click to call)</td>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")

	}

}

func userInfomation(dbConnection *sql.DB, w http.ResponseWriter, dbName string, email string) {

	userTypeId := userAccountTypeId(dbConnection, dbName, email)
	result, err := dbConnection.Query(`SELECT
					     user_account_type_id,
					     user_account_first_name,
					     user_account_last_name,
					     user_account_email,
					     user_account_type,
					     user_account_date_added,
					     user_account_type_description
					   FROM yap.view___account_detail
					   WHERE user_account_email = ?;`, email)

	// Error
	if err != nil {
		panic(err)
	}

	for result.Next() {
		var (
			userAccountTypeId          string
			userAccountFirstName       string
			userAccountLastName        string
			userAccountEmail           string
			userAccountType            string
			userAccountDateAdded       string
			userAccountTypeDescription string
		)

		err = result.Scan(
			&userAccountTypeId,
			&userAccountFirstName,
			&userAccountLastName,
			&userAccountEmail,
			&userAccountType,
			&userAccountDateAdded,
			&userAccountTypeDescription,
		)

		// Error
		if err != nil {
			panic("Error in userInfomation function")
		}

		fmt.Fprintf(w, "<div>")
		fmt.Fprintf(w, "<table id=\"table\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>")
		fmt.Fprintf(w, "      <table id=\"table\">")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <th>&nbsp Name &nbsp</th>")
		fmt.Fprintf(w, "          <th>&nbsp Email &nbsp</th>")
		fmt.Fprintf(w, "          <th>&nbsp Account Type &nbsp</th>")
		fmt.Fprintf(w, "          <th>&nbsp Account Created &nbsp</th>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "        <tr>")
		fmt.Fprintf(w, "          <td>&nbsp"+userAccountFirstName+"&nbsp<br>"+userAccountLastName+"</td>")
		fmt.Fprintf(w, "          <td>&nbsp"+userAccountEmail+"&nbsp</td>")
		fmt.Fprintf(w, "          <td>&nbsp"+userAccountType+"&nbsp</td>")
		fmt.Fprintf(w, "          <td>&nbsp"+userAccountDateAdded+"&nbsp</td>")
		fmt.Fprintf(w, "        </tr>")
		fmt.Fprintf(w, "      </table>")
		fmt.Fprintf(w, "    </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><button onclick=\"toggleAccountDetail() \"class=\"general-button\">&nbsp Show / Hide More Account Details &nbsp</button></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		//Account detail tables
		fmt.Fprintf(w, "</div>")
		fmt.Fprintf(w, "<div id=\"accountDetailDiv\" style=\"display:none\">")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<table id=\"table\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>&nbsp Account Type Description &nbsp</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td>&nbsp"+userAccountTypeDescription+"&nbsp</td>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
	}
	fmt.Fprintf(w, "<br>")
	if userTypeId == "100" || userTypeId == "101" {
		yapAccount(dbConnection, w, dbName, email)
	} else if userTypeId == "200" || userTypeId == "201" {
		groupAccount(dbConnection, w, dbName, email)
	} else if userTypeId == "300" || userTypeId == "301" {
		pbxAccount(dbConnection, w, dbName, email)
	} else {
	}
	fmt.Fprintf(w, "</div>")
	//Java Script embedded to toggle account detail div
	fmt.Fprintf(w, "<script>")
	fmt.Fprintf(w, "  function toggleAccountDetail() {")
	fmt.Fprintf(w, "  var x = document.getElementById(\"accountDetailDiv\");")
	fmt.Fprintf(w, "  if (x.style.display === \"none\") {")
	fmt.Fprintf(w, "    x.style.display = \"table\";")
	fmt.Fprintf(w, "  } else {")
	fmt.Fprintf(w, "    x.style.display = \"none\";")
	fmt.Fprintf(w, "  }")
	fmt.Fprintf(w, "}")
	fmt.Fprintf(w, "</script>")
}

func mainMenuButton(w http.ResponseWriter, buttonName string, hyperlink string) {
	fmt.Fprintf(w, "&nbsp")
	fmt.Fprintf(w, "<h2>")
	fmt.Fprintf(w, "  <a href=\""+hyperlink+"\" class=\"general-button main-menu-button\"><p>"+buttonName+"</p></a>")
	fmt.Fprintf(w, "</h2>")
	fmt.Fprintf(w, "&nbsp")
}

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
		// Code to call the emailHeaderHTTP function
		email := emailHeaderHTTP(r)
		userTypeId := userAccountTypeId(dbConnection, dbName, email)
		if userTypeId == "100" || userTypeId == "101" {
			header(w)
			userInfomation(dbConnection, w, dbName, email)
			fmt.Fprintf(w, "<br>")
			fmt.Fprintf(w, "<div class=\"main-menu\">")
			mainMenuButton(w, "All User<br>Accounts<br>&#128100", "/all-user-accounts")
			mainMenuButton(w, "All<br>Groups<br>&#128101", "/all-groups")
			mainMenuButton(w, "All Phone<br>PBXs<br>&#128222", "/all-pbxs")
			fmt.Fprintf(w, "</div>")
			footer(w)
		} else if userTypeId == "200" || userTypeId == "201" {
			header(w)
			userInfomation(dbConnection, w, dbName, email)
			fmt.Fprintf(w, "<div class=\"main-menu\">")
			mainMenuButton(w, "Edit<br>Group", "/group")
			mainMenuButton(w, "All<br>PBXs", "/all-pbxs")
			mainMenuButton(w, "", "")
			fmt.Fprintf(w, "</div>")
			footer(w)
		} else if userTypeId == "300" || userTypeId == "301" {
			header(w)
			userInfomation(dbConnection, w, dbName, email)
			fmt.Fprintf(w, "<br>")
			fmt.Fprintf(w, "<div class=\"main-menu\">")
			mainMenuButton(w, "", "")
			mainMenuButton(w, "", "")
			mainMenuButton(w, "", "")
			fmt.Fprintf(w, "</div>")
			footer(w)
		} else {
			errorBox(w)
		}
		fmt.Fprintf(w, endHTML)

	})

	// User Accounts Page
	http.HandleFunc("/user-accounts", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, startHTML)
		header(w)
		footer(w)
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
