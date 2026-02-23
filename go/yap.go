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

// Function to retrive HTTP email header
func emailHeaderHTTP(r *http.Request) (email string) {
	email = r.Header.Get("X-Email")
	return email
}

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
func errorBox(w http.ResponseWriter, errorType string) {
	fmt.Fprintf(w, "<div class=\"div-error-box\">")
	fmt.Fprintf(w, "  <h1>")
	if errorType == "email_error" {
		fmt.Fprintf(w, "    User Account Not Found")
	} else if errorType == "account_type_error" {
		fmt.Fprintf(w, "    Account Type Error")
	} else {
		fmt.Fprintf(w, "    Unknown Error")
	}
	fmt.Fprintf(w, "    <br>")
	fmt.Fprintf(w, "    <a href=\"/oauth2/sign_out?rd=https://github.com/logout\" class=\"button-general button-header button-logout\">Logout</a>")
	fmt.Fprintf(w, "  </h1>")
	fmt.Fprintf(w, "</div>")
}

// Function for the header
func header(w http.ResponseWriter, headerTitle string, headerCSS string) {
	fmt.Fprintf(w, "<div class=\"div-header\">")
	fmt.Fprintf(w, "  <h1 class=\""+headerCSS+"\">")
	fmt.Fprintf(w, "    <a href=\"/oauth2/sign_out?rd=https://github.com/logout\" class=\"button-general button-header button-logout\">Logout</a>")
	fmt.Fprintf(w, "    <a href=\"/\" class=\"button-general button-header button-home\">Home</a>")
	fmt.Fprintf(w, "    <a href=\"https://github.com/ellwould/yap/blob/main/LICENSE\" target=\"_blank\" class=\"button-general button-header button-license\">License</a>")
	fmt.Fprintf(w, "    <br>")
	fmt.Fprintf(w, "    ⊛ YAP [Yet Another PBX] ⊛")
	fmt.Fprintf(w, "    <br>")
	fmt.Fprintf(w, "    "+headerTitle+"")
	fmt.Fprintf(w, "  </h1>")
	fmt.Fprintf(w, "</div>")
}

// Function for the footer
func footer(w http.ResponseWriter, headerCSS string, buttonCSS string) {
	fmt.Fprintf(w, "<div class=\"div-footer\">")
	fmt.Fprintf(w, "  <h2 class=\""+headerCSS+"\">")
	fmt.Fprintf(w, "    <a href=\"https://github.com/ellwould/yap\" target=\"_blank\" class=\"button-general button-footer "+buttonCSS+"\">YAP Source Code</a>")
	fmt.Fprintf(w, "    <a href=\"https://ell.today\" target=\"_blank\" class=\"button-general button-footer "+buttonCSS+"\">Other Software</a>")
	fmt.Fprintf(w, "  </h2>")
	fmt.Fprintf(w, "</div>")
}

//----------------------------------------------------------------------------------------------------

type databaseFunctionParameter struct {
	connection       *sql.DB
	database         string
	table            string
	column           string
	columnWhere      string
	columnWhereValue string
	countMinusOne    bool
}

func selectWhere(dbSelectWhere databaseFunctionParameter) string {
	var selectWhere string
	selectWhereQuery, err := dbSelectWhere.connection.Query("SELECT "+dbSelectWhere.column+" FROM "+dbSelectWhere.database+"."+dbSelectWhere.table+" WHERE "+dbSelectWhere.columnWhere+" = ?;", dbSelectWhere.columnWhereValue)
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
		countMinusOneQuery, err := dbTotalTableCount.connection.Query("SELECT COUNT(*) -1 FROM " + dbTotalTableCount.database + "." + dbTotalTableCount.table)
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
		countQuery, err := dbTotalTableCount.connection.Query("SELECT COUNT(*) FROM " + dbTotalTableCount.database + "." + dbTotalTableCount.table)
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
	var count string
	countQuery, err := dbTotalTableCountWhere.connection.Query("SELECT COUNT(*) FROM "+dbTotalTableCountWhere.database+"."+dbTotalTableCountWhere.table+" WHERE "+dbTotalTableCountWhere.columnWhere+" =?", dbTotalTableCountWhere.columnWhereValue)
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

func userAccountTypeID(dbUserAccountTypeID databaseFunctionParameter) string {

	var dbSelectWhere databaseFunctionParameter
	dbSelectWhere.connection = dbUserAccountTypeID.connection
	dbSelectWhere.database = dbUserAccountTypeID.database
	dbSelectWhere.table = "view___account_detail"
	dbSelectWhere.column = "user_account_type_id"
	dbSelectWhere.columnWhere = "user_account_email"
	dbSelectWhere.columnWhereValue = dbUserAccountTypeID.columnWhereValue

	return selectWhere(dbSelectWhere)
}

//----------------------------------------------------------------------------------------------------

func yapAccount(w http.ResponseWriter, dbYapAccount databaseFunctionParameter) {

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
	fmt.Fprintf(w, "    <th>&nbsp Total Groups &nbsp</th>")
	fmt.Fprintf(w, "    <th>&nbsp Total PBXs &nbsp</th>")
	fmt.Fprintf(w, "    <th>&nbsp Total SIP Endpoints &nbsp</th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "  <tr>")
	dbTotalTableCount.table = "group"
	dbTotalTableCount.countMinusOne = true
	fmt.Fprintf(w, "    <td>&nbsp"+totalTableCount(w, dbTotalTableCount)+"&nbsp</td>")
	dbTotalTableCount.table = "pbx"
	dbTotalTableCount.countMinusOne = true
	fmt.Fprintf(w, "    <td>&nbsp"+totalTableCount(w, dbTotalTableCount)+"&nbsp</td>")
	dbTotalTableCount.table = "ps_endpoints"
	dbTotalTableCount.countMinusOne = false
	fmt.Fprintf(w, "    <td>&nbsp"+totalTableCount(w, dbTotalTableCount)+"&nbsp</td>")
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
	fmt.Fprintf(w, "    <td>&nbsp"+totalTableCountWhere(w, dbTotalTableCountWhere)+"&nbsp</td>")
	dbTotalTableCountWhere.columnWhereValue = "200"
	fmt.Fprintf(w, "    <td>&nbsp"+totalTableCountWhere(w, dbTotalTableCountWhere)+"&nbsp</td>")
	dbTotalTableCountWhere.columnWhereValue = "201"
	fmt.Fprintf(w, "    <td>&nbsp"+totalTableCountWhere(w, dbTotalTableCountWhere)+"&nbsp</td>")
	dbTotalTableCountWhere.columnWhereValue = "300"
	fmt.Fprintf(w, "    <td>&nbsp"+totalTableCountWhere(w, dbTotalTableCountWhere)+"&nbsp</td>")
	dbTotalTableCountWhere.columnWhereValue = "301"
	fmt.Fprintf(w, "    <td>&nbsp"+totalTableCountWhere(w, dbTotalTableCountWhere)+"&nbsp</td>")
	dbTotalTableCountWhere.columnWhereValue = "302"
	fmt.Fprintf(w, "    <td>&nbsp"+totalTableCountWhere(w, dbTotalTableCountWhere)+"&nbsp</td>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "</table>")

}

func groupAccount(w http.ResponseWriter, dbGroupAccount databaseFunctionParameter) {

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
					   FROM yap.view___account_detail
					   WHERE user_account_email = ?;`, dbGroupAccount.columnWhereValue)

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
			pbxID                         string
		)

		err = result.Scan(
			&groupName,
			&groupID,
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
			&pbxID,
		)

		// Error
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, "<table id=\"table\" class\"table-main-menu\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>&nbsp Group Name and ID &nbsp</th>")
		fmt.Fprintf(w, "    <th>&nbsp Total PBX(s) in Group &nbsp</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td>&nbsp Group Name: "+groupName+"&nbsp<br><br>Group ID: "+groupID+"</td>")
		var dbTotalTableCountWhere databaseFunctionParameter
		dbTotalTableCountWhere.connection = dbGroupAccount.connection
		dbTotalTableCountWhere.database = dbGroupAccount.database
		dbTotalTableCountWhere.table = "pbx"
		dbTotalTableCountWhere.columnWhere = "id"
		dbTotalTableCountWhere.columnWhereValue = pbxID
		fmt.Fprintf(w, "    <td>&nbsp"+totalTableCountWhere(w, dbTotalTableCountWhere)+"&nbsp</td>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<table id=\"table\" class\"table-main-menu\">")
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
		fmt.Fprintf(w, "<table id=\"table\" class\"table-main-menu\">")
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

func pbxAccount(w http.ResponseWriter, dbPBXAccount databaseFunctionParameter) {

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
					   FROM yap.view___account_detail
					   WHERE user_account_email = ?;`, dbPBXAccount.columnWhereValue)

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

		fmt.Fprintf(w, "<table id=\"table\" class\"table-main-menu\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>&nbsp PBX Name and ID &nbsp</th>")
		fmt.Fprintf(w, "    <th>&nbsp Total SIP Extensions in PBX &nbsp</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td>&nbsp PBX Name: "+pbxName+"&nbsp<br><br>PBX ID: "+pbxID+"</td>")
		var dbTotalTableCountWhere databaseFunctionParameter
		dbTotalTableCountWhere.connection = dbPBXAccount.connection
		dbTotalTableCountWhere.database = dbPBXAccount.database
		dbTotalTableCountWhere.table = "ps_endpoints"
		dbTotalTableCountWhere.columnWhere = "pbx_id"
		dbTotalTableCountWhere.columnWhereValue = pbxID
		fmt.Fprintf(w, "    <td>&nbsp"+totalTableCountWhere(w, dbTotalTableCountWhere)+"&nbsp</td>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<table id=\"table\" class\"table-main-menu\">")
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
		fmt.Fprintf(w, "<table id=\"table\" class\"table-main-menu\">")
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

func userInformation(w http.ResponseWriter, dbUserInformation databaseFunctionParameter, userTypeID string) {

	result, err := dbUserInformation.connection.Query(`SELECT
					     user_account_type_id,
					     user_account_first_name,
					     user_account_last_name,
					     user_account_email,
					     user_account_type,
					     user_account_date_added,
					     user_account_type_permission
					   FROM yap.view___account_detail
					   WHERE user_account_email = ?;`, dbUserInformation.columnWhereValue)

	// Error
	if err != nil {
		panic(err)
	}

	for result.Next() {
		var (
			userAccountTypeId         string
			userAccountFirstName      string
			userAccountLastName       string
			userAccountEmail          string
			userAccountType           string
			userAccountDateAdded      string
			userAccountTypePermission string
		)

		err = result.Scan(
			&userAccountTypeId,
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
		fmt.Fprintf(w, "    <th><button onclick=\"toggleAccountDetail() \"class=\"button-general\">&nbsp Show / Hide More Account Details &nbsp</button></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		//Account detail tables
		fmt.Fprintf(w, "</div>")
		fmt.Fprintf(w, "<div id=\"accountDetailDiv\" style=\"display:none\">")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<table id=\"table\" class=\"table-main-menu\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>&nbsp Account Type Permissions - "+userAccountType+"&nbsp</th>")
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

	if userTypeID == "100" || userTypeID == "101" {
		yapAccount(w, dbDetail)
	} else if userTypeID == "200" || userTypeID == "201" {
		groupAccount(w, dbDetail)
	} else if userTypeID == "300" || userTypeID == "301" {
		pbxAccount(w, dbDetail)
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

type mainMenuParameter struct {
	writeHTTP  http.ResponseWriter
	buttonName string
	hyperlink  string
	headerCSS  string
	buttonCSS  string
}

func mainMenuButton(mainMenu mainMenuParameter) {
	fmt.Fprintf(mainMenu.writeHTTP, "&nbsp")
	fmt.Fprintf(mainMenu.writeHTTP, "<h2 class=\""+mainMenu.headerCSS+"\">")
	fmt.Fprintf(mainMenu.writeHTTP, "  <a href=\""+mainMenu.hyperlink+"\" class=\"button-general button-main-menu "+mainMenu.buttonCSS+"\"><p>"+mainMenu.buttonName+"</p></a>")
	fmt.Fprintf(mainMenu.writeHTTP, "</h2>")
	fmt.Fprintf(mainMenu.writeHTTP, "&nbsp")
}

//----------------------------------------------------------------------------------------------------

func download(id string, file string) {
	http.HandleFunc("/download/"+id+"/"+file, func(w http.ResponseWriter, r *http.Request) {
		dirPathB := "/var/lib/yap/call-recording/" + file + ""
		http.ServeFile(w, r, dirPathB)
	})
}

func playAudio(id string, file string) {
	http.HandleFunc("/play-audio", func(w http.ResponseWriter, r *http.Request) {
	})
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

		// Main Menu Wallpaper
		fmt.Fprintf(w, "<div class=\"wallpaper\"></div>")
		fmt.Fprintf(w, "<div class=\"wallpaper wallpaper2\"></div>")
		fmt.Fprintf(w, "<div class=\"wallpaper wallpaper3\"></div>")

		// Code to call the emailHeaderHTTP function
		email := emailHeaderHTTP(r)

		var dbDetail databaseFunctionParameter
		dbDetail.connection = dbConnection
		dbDetail.database = dbName
		dbDetail.columnWhereValue = email

		userTypeID := userAccountTypeID(dbDetail)

		if userTypeID == "" {
			errorBox(w, "email_error")
		} else {
			if userTypeID == "100" {
				header(w, "Main Menu", "")
				userInformation(w, dbDetail, userTypeID)
				fmt.Fprintf(w, "<br>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonOne := mainMenuParameter{writeHTTP: w, buttonName: "All User<br>Accounts<br>&#128100", hyperlink: "/", headerCSS: "header-user-account", buttonCSS: "button-user-account"}
				mainMenuButton(mainMenuButtonOne)
				mainMenuButtonTwo := mainMenuParameter{writeHTTP: w, buttonName: "All<br>Groups<br>&#128101", hyperlink: "/", headerCSS: "header-group", buttonCSS: "button-group"}
				mainMenuButton(mainMenuButtonTwo)
				mainMenuButtonThree := mainMenuParameter{writeHTTP: w, buttonName: "All<br>PBXs<br>&#128222", hyperlink: "/", headerCSS: "header-pbx", buttonCSS: "button-pbx"}
				mainMenuButton(mainMenuButtonThree)
				fmt.Fprintf(w, "</div>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonFour := mainMenuParameter{writeHTTP: w, buttonName: "All SIP<br>Endpoints<br>&#128241", hyperlink: "/", headerCSS: "header-sip-endpoint", buttonCSS: "button-sip-endpoint"}
				mainMenuButton(mainMenuButtonFour)
				mainMenuButtonFive := mainMenuParameter{writeHTTP: w, buttonName: "All SIP<br>Trunks<br>&#8596", hyperlink: "/", headerCSS: "header-sip-trunk", buttonCSS: "button-sip-trunk"}
				mainMenuButton(mainMenuButtonFive)
				mainMenuButtonSix := mainMenuParameter{writeHTTP: w, buttonName: "All Phone<br>Numbers<br>&#128290", hyperlink: "/", headerCSS: "header-phone-number", buttonCSS: "button-phone-number"}
				mainMenuButton(mainMenuButtonSix)
				fmt.Fprintf(w, "</div>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonSeven := mainMenuParameter{writeHTTP: w, buttonName: "All<br>CDRs<br>&#128202", hyperlink: "/", headerCSS: "header-cdr", buttonCSS: "button-cdr"}
				mainMenuButton(mainMenuButtonSeven)
				mainMenuButtonEight := mainMenuParameter{writeHTTP: w, buttonName: "All<br>Voicemails<br>&#127897", hyperlink: "/", headerCSS: "header-voicemail", buttonCSS: "button-voicemail"}
				mainMenuButton(mainMenuButtonEight)
				mainMenuButtonNine := mainMenuParameter{writeHTTP: w, buttonName: "All Call<br>Recordings<br>&#128252", hyperlink: "/call-recording", headerCSS: "header-call-recording", buttonCSS: "button-call-recording"}
				mainMenuButton(mainMenuButtonNine)
				fmt.Fprintf(w, "</div>")
				fmt.Fprintf(w, "<div class=\"div-main-menu\">")
				mainMenuButtonTen := mainMenuParameter{writeHTTP: w, buttonName: "All MoH / AA<br>Music<br>&#127925", hyperlink: "/", headerCSS: "header-music", buttonCSS: "button-music"}
				mainMenuButton(mainMenuButtonTen)
				mainMenuButtonEleven := mainMenuParameter{writeHTTP: w, buttonName: "All Server<br>Logs<br>&#128195", hyperlink: "/", headerCSS: "header-server-log", buttonCSS: "button-server-log"}
				mainMenuButton(mainMenuButtonEleven)
				mainMenuButtonTweleve := mainMenuParameter{writeHTTP: w, buttonName: "Server<br>Information<br>&#128421", hyperlink: "/", headerCSS: "header-server-information", buttonCSS: "button-server-information"}
				mainMenuButton(mainMenuButtonTweleve)
				fmt.Fprintf(w, "</div>")
				footer(w, "", "")
			} else if userTypeID == "200" || userTypeID == "201" {
				header(w, "", "")
				userInformation(w, dbDetail, userTypeID)
				fmt.Fprintf(w, "<br>")
				fmt.Fprintf(w, "<div class=\"main-menu\">")
				mainMenuButtonOne := mainMenuParameter{writeHTTP: w, buttonName: "All User<br>Accounts<br>&#128100", hyperlink: "/", headerCSS: "user-main-menu-header", buttonCSS: "user-main-menu-button"}
				mainMenuButton(mainMenuButtonOne)
				mainMenuButtonTwo := mainMenuParameter{writeHTTP: w, buttonName: "All<br>Groups<br>&#128101", hyperlink: "/", headerCSS: "group-main-menu-header", buttonCSS: "group-main-menu-button"}
				mainMenuButton(mainMenuButtonTwo)
				mainMenuButtonThree := mainMenuParameter{writeHTTP: w, buttonName: "All<br>PBXs<br>&#128222", hyperlink: "/", headerCSS: "pbx-main-menu-header", buttonCSS: "pbx-main-menu-button"}
				mainMenuButton(mainMenuButtonThree)
				fmt.Fprintf(w, "</div>")
				fmt.Fprintf(w, "<div class=\"main-menu\">")
				mainMenuButtonFour := mainMenuParameter{writeHTTP: w, buttonName: "All SIP<br>Endpoints<br>&#128241", hyperlink: "/", headerCSS: "endpoint-main-menu-header", buttonCSS: "endpoint-main-menu-button"}
				mainMenuButton(mainMenuButtonFour)
				mainMenuButtonFive := mainMenuParameter{writeHTTP: w, buttonName: "All SIP<br>Trunks<br>&#8596", hyperlink: "/", headerCSS: "trunk-main-menu-header", buttonCSS: "trunk-main-menu-button"}
				mainMenuButton(mainMenuButtonFive)
				mainMenuButtonSix := mainMenuParameter{writeHTTP: w, buttonName: "All Phone<br>Numbers<br>&#128290", hyperlink: "/", headerCSS: "number-main-menu-header", buttonCSS: "number-main-menu-button"}
				mainMenuButton(mainMenuButtonSix)
				fmt.Fprintf(w, "</div>")
				fmt.Fprintf(w, "<div class=\"main-menu\">")
				mainMenuButtonSeven := mainMenuParameter{writeHTTP: w, buttonName: "All <br>CDRs<br>&#128202", hyperlink: "/", headerCSS: "cdr-main-menu-header", buttonCSS: "cdr-main-menu-button"}
				mainMenuButton(mainMenuButtonSeven)
				mainMenuButtonEight := mainMenuParameter{writeHTTP: w, buttonName: "All Server<br>Logs<br>&#128210", hyperlink: "/", headerCSS: "log-main-menu-header", buttonCSS: "log-main-menu-button"}
				mainMenuButton(mainMenuButtonEight)
				mainMenuButtonNine := mainMenuParameter{writeHTTP: w, buttonName: "Server<br>Information<br>&#128421", hyperlink: "/", headerCSS: "information-main-menu-header", buttonCSS: "information-main-menu-button"}
				mainMenuButton(mainMenuButtonNine)
				fmt.Fprintf(w, "</div>")
				footer(w, "", "")
			} else if userTypeID == "300" || userTypeID == "301" || userTypeID == "302" {
				header(w, "", "")
				userInformation(w, dbDetail, userTypeID)
				fmt.Fprintf(w, "<br>")
				fmt.Fprintf(w, "<div class=\"main-menu\">")
				mainMenuButtonOne := mainMenuParameter{writeHTTP: w, buttonName: "All User<br>Accounts<br>&#128100", hyperlink: "/", headerCSS: "user-main-menu-header", buttonCSS: "user-main-menu-button"}
				mainMenuButton(mainMenuButtonOne)
				mainMenuButtonTwo := mainMenuParameter{writeHTTP: w, buttonName: "All<br>Groups<br>&#128101", hyperlink: "/", headerCSS: "group-main-menu-header", buttonCSS: "group-main-menu-button"}
				mainMenuButton(mainMenuButtonTwo)
				mainMenuButtonThree := mainMenuParameter{writeHTTP: w, buttonName: "All<br>PBXs<br>&#128222", hyperlink: "/", headerCSS: "pbx-main-menu-header", buttonCSS: "pbx-main-menu-button"}
				mainMenuButton(mainMenuButtonThree)
				fmt.Fprintf(w, "</div>")
				fmt.Fprintf(w, "<div class=\"main-menu\">")
				mainMenuButtonFour := mainMenuParameter{writeHTTP: w, buttonName: "All SIP<br>Endpoints<br>&#128241", hyperlink: "/", headerCSS: "endpoint-main-menu-header", buttonCSS: "endpoint-main-menu-button"}
				mainMenuButton(mainMenuButtonFour)
				mainMenuButtonFive := mainMenuParameter{writeHTTP: w, buttonName: "All SIP<br>Trunks<br>&#8596", hyperlink: "/", headerCSS: "trunk-main-menu-header", buttonCSS: "trunk-main-menu-button"}
				mainMenuButton(mainMenuButtonFive)
				mainMenuButtonSix := mainMenuParameter{writeHTTP: w, buttonName: "All Phone<br>Numbers<br>&#128290", hyperlink: "/", headerCSS: "number-main-menu-header", buttonCSS: "number-main-menu-button"}
				mainMenuButton(mainMenuButtonSix)
				fmt.Fprintf(w, "</div>")
				fmt.Fprintf(w, "<div class=\"main-menu\">")
				mainMenuButtonSeven := mainMenuParameter{writeHTTP: w, buttonName: "All <br>CDRs<br>&#128202", hyperlink: "/", headerCSS: "cdr-main-menu-header", buttonCSS: "cdr-main-menu-button"}
				mainMenuButton(mainMenuButtonSeven)
				mainMenuButtonEight := mainMenuParameter{writeHTTP: w, buttonName: "All Server<br>Logs<br>&#128210", hyperlink: "/", headerCSS: "log-main-menu-header", buttonCSS: "log-main-menu-button"}
				mainMenuButton(mainMenuButtonEight)
				mainMenuButtonNine := mainMenuParameter{writeHTTP: w, buttonName: "Server<br>Information<br>&#128421", hyperlink: "/", headerCSS: "information-main-menu-header", buttonCSS: "information-main-menu-button"}
				mainMenuButton(mainMenuButtonNine)
				fmt.Fprintf(w, "</div>")
				footer(w, "", "")
			} else {
				errorBox(w, "account_type_error")
			}
		}
		fmt.Fprintf(w, endHTML)

	})

	// Call Recording Page

	http.HandleFunc("/call-recording", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, startHTML)
		// Call Recording Wallpaper
		fmt.Fprintf(w, "<div class=\"wallpaper wallpaper-call-recording\"></div>")
		fmt.Fprintf(w, "<div class=\"wallpaper wallpaper-call-recording wallpaper2\"></div>")
		fmt.Fprintf(w, "<div class=\"wallpaper wallpaper-call-recording wallpaper3\"></div>")
		header(w, "Call Recordings", "header-call-recording")
		fmt.Fprintf(w, "<table id=\"table\" class=\"table-call-recording\">")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>&nbsp From (Caller) &nbsp</th>")
		fmt.Fprintf(w, "    <th>&nbsp To (Callie) &nbsp</th>")
		fmt.Fprintf(w, "    <th>&nbsp Date (YYYY-MM-DD)&nbsp<br>Time (HH:MM:SS)</th>")
		fmt.Fprintf(w, "    <th>&nbsp Listen &nbsp</th>")
		fmt.Fprintf(w, "    <th>&nbsp Download &nbsp</th>")
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
