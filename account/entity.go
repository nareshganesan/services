package account

import (
	"encoding/json"
	// "fmt"
	g "github.com/nareshganesan/services/globals"
	"github.com/nareshganesan/services/shared"
	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
)

// Entity for account index
type Entity struct {
	ID                string
	Email             string `form:"email" json:"email"`
	Password          string `form:"password" json:"password"`
	Name              string `form:"name" json:"name,omitempty"`
	Title             string `form:"title" json:"title,omitempty"`
	Roles             string `form:"roles" json:"roles,omitempty"`
	VerificationToken string `form:"verification_token" json:"verification_token,omitempty"`
	IsArchived        bool   `form:"is_archived" json:"is_archived" binding:"omitempty"`
	IsVerified        bool   `form:"is_verified" json:"is_verified" binding:"omitempty"`
	IsLocked          bool   `json:"is_locked" binding:"omitempty"`
	FailedAttempts    int    `json:"failed_attempts" binding:"omitempty"`
}

// Create for account entity
func (a *Entity) Create() (string, bool) {
	es := g.GetGlobals()
	a.IsVerified = true
	a.IsArchived = false
	a.IsLocked = false
	pHash, err := shared.GetHash(a.Password)
	if err != nil {
		return "", false
	}
	a.Password = pHash
	doc, err := es.Index(
		g.Config.ES.Index.Accounts.Name,
		g.Config.ES.Index.Accounts.DocType, a)
	if err != nil {
		return "", false
	}
	return doc.Id, true
}

// GetByID for account entity.
func (a *Entity) GetByID() bool {
	es := g.GetGlobals()
	_, err := es.Get(
		g.Config.ES.Index.Accounts.Name,
		g.Config.ES.Index.Accounts.DocType, a.ID)
	if err != nil {
		return false
	}
	return true
}

// Update for account entity
func (a *Entity) Update() (string, bool) {
	es := g.GetGlobals()
	// a.IsVerified = true
	data := EntityToMap(a)
	doc, err := es.Update(
		g.Config.ES.Index.Accounts.Name,
		g.Config.ES.Index.Accounts.DocType,
		a.ID, data)
	if err != nil {
		return "", false
	}
	return doc.Id, true
}

// Delete for account entity
func (a *Entity) Delete() (string, bool) {
	es := g.GetGlobals()
	var delAcc Entity
	delAcc.ID = a.ID
	delAcc.IsArchived = true
	data := EntityToMap(&delAcc)
	doc, err := es.Delete(
		g.Config.ES.Index.Accounts.Name,
		g.Config.ES.Index.Accounts.DocType, delAcc.ID, data)
	if err != nil {
		return "", false
	}
	return doc.Id, true
}

// IsExistingUser method, identifies if accountcredentials given in the query
// and returns true if exists else false
func (a *Entity) IsExistingUser() bool {
	query := GetSearchByFieldQuery("email", a.Email)
	isExisting := false
	existingAccount := a.FetchOne(query)
	if existingAccount == nil {
		isExisting = false
	} else {
		isExisting = true
	}
	return isExisting
}

// List method, lists all the accounts given page number and size per page
// Uses match_all query
// page number starts from 0
func (a *Entity) List(page, size int, query elastic.Query) []Entity {
	es := g.GetGlobals()
	l := es.Log
	docs := es.Search(
		g.Config.ES.Index.Accounts.Name,
		g.Config.ES.Index.Accounts.DocType, page, size, query)
	if docs == nil {
		return nil
	}
	var count = docs.Hits.TotalHits
	var accounts = make([]Entity, count)

	for idx, hit := range docs.Hits.Hits {
		var account Entity
		err := json.Unmarshal(*hit.Source, &account)
		if err != nil {
			l.WithFields(logrus.Fields{
				"error": err,
			}).Error("Failed to unmarshal account result")
		}
		account.ID = hit.Id
		accounts[idx] = account
	}
	return accounts

}

// Search returns list of account entity.
// Given any query object, page number and size (number of accounts per page)
// page number starts from 0
func (a *Entity) Search(page, size int, query elastic.Query) []Entity {
	es := g.GetGlobals()
	l := es.Log
	docs := es.Search(
		g.Config.ES.Index.Accounts.Name,
		g.Config.ES.Index.Accounts.DocType, page, size, query)
	if docs == nil {
		return nil
	}
	var count = docs.Hits.TotalHits
	var accounts = make([]Entity, count)

	for idx, hit := range docs.Hits.Hits {
		var account Entity
		err := json.Unmarshal(*hit.Source, &account)
		if err != nil {
			l.WithFields(logrus.Fields{
				"error": err,
			}).Error("Failed to unmarshal account result")
		}
		account.ID = hit.Id
		accounts[idx] = account
	}
	return accounts
}

// FetchOne returns an account, only if one account is matched for the query
// else always returns nil
func (a *Entity) FetchOne(query elastic.Query) *Entity {
	es := g.GetGlobals()
	l := es.Log
	doc := es.FetchOne(
		g.Config.ES.Index.Accounts.Name,
		g.Config.ES.Index.Accounts.DocType, query)
	if doc == nil {
		return nil
	}
	var account Entity
	hit := doc.Hits.Hits[0]
	err := json.Unmarshal(*hit.Source, &account)
	if err != nil {
		l.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to unmarshal account result")
	}
	account.ID = hit.Id
	return &account
}

// Authenticate authenticates account entity
// using auth query
func (a *Entity) Authenticate() bool {
	es := g.GetGlobals()
	l := es.Log
	query := GetSearchByFieldQuery("email", a.Email)
	isAuthenticated := false
	existingAccount := a.FetchOne(query)
	if existingAccount == nil {
		isAuthenticated = false
		// either there are multiple accounts or no accounts
		// fetchone returns nil if results > 1 or results == 0
		l.Error("Account does not exist")
	} else {
		if existingAccount.IsArchived {
			l.Error("Account is archived")
		} else {
			l.Info("Account is not archived")
			if !existingAccount.IsLocked {
				l.Info("Account is not locked")
				if shared.VerifyHash(a.Password, existingAccount.Password) {
					if existingAccount.IsVerified {
						l.Info("Account is verified")
						isAuthenticated = true
						existingAccount.FailedAttempts = 0
						existingAccount.Update()
					} else {
						l.Error("Account is not verified")
					}
				} else {
					l.Info("Account is locked")
					existingAccount.loginAttempts()
					existingAccount.Update()
				}
			} else {
				l.Info("Account is locked")
			}
		}
		if existingAccount.ID != "" {
			a.ID = existingAccount.ID
			a.IsArchived = existingAccount.IsArchived
			a.IsLocked = existingAccount.IsLocked
			a.FailedAttempts = existingAccount.FailedAttempts
		}
	}
	return isAuthenticated
}

// loginAttempts calculates the no of failed login attempts
func (a *Entity) loginAttempts() {
	es := g.GetGlobals()
	l := es.Log
	attempts := a.FailedAttempts
	if attempts > 0 {
		a.FailedAttempts++
	} else {
		a.FailedAttempts = 1
	}
	if a.FailedAttempts > g.Config.Tokens.Auth.Attempts {
		l.Error("Account is locked")
	} else if a.FailedAttempts == g.Config.Tokens.Auth.Attempts {
		a.IsLocked = true
		a.IsVerified = false
		l.WithFields(logrus.Fields{
			"failedattempts": a.FailedAttempts,
			"maxattempts":    g.Config.Tokens.Auth.Attempts,
		}).Error("No of failed attempts exceeded maximum allowed attempts!")
	}
	l.WithFields(logrus.Fields{
		"failedattempts":     a.FailedAttempts,
		"maxattempts":        g.Config.Tokens.Auth.Attempts,
		"remaining attempts": (g.Config.Tokens.Auth.Attempts - a.FailedAttempts),
	}).Error("No of failed attempts exceeded maximum allowed attempts!")
}
