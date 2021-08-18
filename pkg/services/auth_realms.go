// Copyright Red Hat

package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/identitatem/idp-configs-api/pkg/common"
	"github.com/identitatem/idp-configs-api/pkg/db"
	"github.com/identitatem/idp-configs-api/pkg/errors"
	"github.com/identitatem/idp-configs-api/pkg/models"
)

func GetAuthRealmsForAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var authRealms []models.AuthRealm

	// Get account from request header
	account, err := common.GetAccount(r)

	if (err != nil) {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return		
	}
	
	// Fetch Auth Realms for specific account from the DB
	result := db.DB.Where("Account = ?", account).Find(&authRealms)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusBadRequest)
		return
	}

	// TODO: support filtering and searching by name (query param)

	// Respond with auth realms for the account
	json.NewEncoder(w).Encode(&authRealms)	
}

func CreateAuthRealmForAccount(w http.ResponseWriter, r *http.Request) {
    var authRealm models.AuthRealm

	fmt.Println(r.Context())

	w.Header().Set("Content-Type", "application/json")

    err := json.NewDecoder(r.Body).Decode(&authRealm)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	// Request body must contain the auth-realm Name and Custom Resource
	if (authRealm.Name == "" || authRealm.CustomResource == nil) {
		// Bad request
		err := errors.NewBadRequest("The request body must contain 'Name' and 'CustomResource'")
		w.WriteHeader(err.Status)
		json.NewEncoder(w).Encode(&err)
		return	
	}

	// TODO: Additional validations on Custom Resource/ evaluating checksum

	// Validate account from request
	err = checkAccount(&authRealm, r)

	if (err != nil) {
		// Bad request
		err := errors.NewBadRequest(err.Error())
		w.WriteHeader(err.Status)
		json.NewEncoder(w).Encode(&err)
		return			
	}

	// Temporarily responding with the auth realm object that will be submitted to the DB
	authRealmJSON, _ := json.Marshal(authRealm)

	// Create record for auth_realm in the DB		
	tx := db.DB.Create(&authRealm)
	if tx.Error != nil {
		// Error updating the DB
		err := errors.NewInternalServerError("Error creating record in the DB: " + tx.Error.Error())
		w.WriteHeader(err.Status)
		json.NewEncoder(w).Encode(&err)
		return			
	}

	fmt.Fprint(w, string(authRealmJSON))
}

func GetAuthRealmByID(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "GetAuthRealmByID")
}

func UpdateAuthRealmByID(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "UpdateAuthRealmByID")
}

func DeleteAuthRealmByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DeleteAuthRealmByID")
}

func checkAccount (authRealm *models.AuthRealm, r *http.Request) (error) {	
	// Get account from request header
	account, err := common.GetAccount(r)

	if err != nil {
		return err
	}

	// Check if request body contains an Account (optional field)
	if (authRealm.Account != "" && authRealm.Account != account) {
		// Account in the request body must match the account of the authenticated user
		return fmt.Errorf("account in the request body does not match account for the authenticated user")
	}

	// Set the account from the request header as the account for the auth realm
	authRealm.Account = account

	return nil
}