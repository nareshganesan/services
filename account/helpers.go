package account

import (
	"fmt"
	"github.com/olivere/elastic"
)

// GetListAccountQuery returns match_all es query for account index
func GetListAccountQuery() *elastic.MatchAllQuery {
	query := elastic.NewMatchAllQuery()
	query = query.QueryName("ListAccountsQuery")
	return query
}

// GetSearchByFieldQuery returns term es query on given field for account index
func GetSearchByFieldQuery(field, value string) *elastic.TermQuery {
	query := elastic.NewTermQuery(field, value)
	qName := fmt.Sprintf("%sTermQuery", field)
	query = query.QueryName(qName)
	return query
}
