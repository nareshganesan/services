package globals

import (
	"context"
	"fmt"
	// "github.com/nareshganesan/services/shared"
	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func newESDB() (*elastic.Client, error) {
	l := Gbl.Log
	opts := []elastic.ClientOptionFunc{
		elastic.SetURL(Config.ES.Urls...),
		elastic.SetBasicAuth(Config.ES.Username, Config.ES.Password),
		elastic.SetSniff(Config.ES.Sniff),
		elastic.SetHealthcheckInterval(
			time.Duration(Config.ES.Health) * time.Second),
		elastic.SetGzip(Config.ES.Gzip),
		elastic.SetMaxRetries(Config.ES.Retries),
		// logrus logging client will take care of info, error and trace logs internally
		elastic.SetErrorLog(GetGlobals().ESLog),
		elastic.SetInfoLog(GetGlobals().ESLog),
	}
	// Create an Elasticsearch client
	client, err := elastic.NewClient(opts...)
	if err != nil {
		l.WithFields(logrus.Fields{
			"error": err,
		}).Error("Could not create elastic client.")
		return nil, err
	}
	return client, nil
}

// ConfigureElasticDB configures globals application with ElasticSearch client
func ConfigureElasticDB() {
	l := Gbl.Log
	es, err := newESDB()
	if err != nil {
		l.WithFields(logrus.Fields{
			"error": err,
		}).Error("Could not create elastic connection.")
		os.Exit(1)
	}
	Gbl.ES = es
	return
}

func (g *Globals) indexExists(index string) bool {
	l := Gbl.Log
	esCtx := context.Background()
	exists, err := g.ES.IndexExists(index).Do(esCtx)
	if err != nil {
		// Handle error
		l.WithFields(logrus.Fields{
			"error": err,
			"index": index,
		}).Error("Error checking index")
		return exists
	}
	if exists {
		l.WithFields(logrus.Fields{
			"index": index,
		}).Info("Index exists")
		fmt.Println(exists)
		return exists
	}
	return exists

}

// Following methods are ElasticSearch client API in global application object

// Index method for ES given index and doctype
func (g *Globals) Index(index, docType string, obj interface{}) (*elastic.IndexResponse, error) {
	l := Gbl.Log
	esCtx := context.Background()
	doc, err := g.ES.Index().
		Index(index).
		Type(docType).
		BodyJson(obj).
		Refresh("true").
		Do(esCtx)
	if err != nil {
		l.WithFields(logrus.Fields{
			"error":   err,
			"index":   index,
			"doctype": docType,
		}).Error("Error creating document")
		return nil, err
	}
	l.WithFields(logrus.Fields{
		"index":   index,
		"doctype": docType,
	}).Info("created document")
	return doc, nil
}

// Get method for ES given ID, index and doctype
func (g *Globals) Get(index, docType, docID string) (*elastic.GetResult, error) {
	l := Gbl.Log
	esCtx := context.Background()
	doc, err := g.ES.Get().
		Index(index).
		Type(docType).
		Id(docID).
		Do(esCtx)
	if err != nil {
		l.WithFields(logrus.Fields{
			"error":   err,
			"index":   index,
			"doctype": docType,
			"id":      docID,
		}).Error("Error getting document")
		return nil, err
	}
	l.WithFields(logrus.Fields{
		"index":   index,
		"doctype": docType,
		"id":      docID,
	}).Info("found document")
	return doc, nil
}

// Update method for ES given ID, index and doctype
func (g *Globals) Update(index, docType, docID string, obj interface{}) (*elastic.UpdateResponse, error) {
	l := Gbl.Log
	esCtx := context.Background()
	doc, err := g.ES.Update().Index(index).Type(docType).Id(docID).
		DocAsUpsert(true).
		Doc(obj).
		Refresh("true").
		Do(esCtx)
	if err != nil {
		l.WithFields(logrus.Fields{
			"error":   err,
			"index":   index,
			"doctype": docType,
		}).Error("Error updating document")
		return nil, err
	}
	l.WithFields(logrus.Fields{
		"index":   index,
		"doctype": docType,
		"id":      docID,
	}).Info("updated document")
	return doc, nil
}

// Delete method is application level soft delete for ES given ID, index and doctype
func (g *Globals) Delete(index, docType, docID string, obj interface{}) (*elastic.UpdateResponse, error) {
	return g.Update(index, docType, docID, obj)
}

// HardDelete method deletes the document from ES given ID, index and doctype
func (g *Globals) HardDelete(index, docType, docID string) (string, error) {
	l := Gbl.Log
	esCtx := context.Background()
	res, err := g.ES.Delete().
		Index(index).
		Type(docType).
		Id(docID).
		Refresh("true").
		Do(esCtx)
	if err != nil {
		l.WithFields(logrus.Fields{
			"error":   err,
			"index":   index,
			"doctype": docType,
		}).Error("Error deleting document")
		return "", err
	}
	l.WithFields(logrus.Fields{
		"index":   index,
		"doctype": docType,
		"id":      docID,
		"res":     res,
	}).Info("document deleted")
	return docID, nil
}

// Search method for ES given index, doctype, page number, size and query
func (g *Globals) Search(index, docType string, page, size int, query elastic.Query) *elastic.SearchResult {
	l := Gbl.Log
	esCtx := context.Background()
	results, err := g.ES.Search().
		Index(index).
		Query(query).
		From(page).
		Size(size).
		Pretty(true).
		Do(esCtx)
	if err != nil {
		l.WithFields(logrus.Fields{
			"error":   err,
			"index":   index,
			"doctype": docType,
			"query":   query,
		}).Error("Error searching document")
	}
	l.WithFields(logrus.Fields{
		"index":      index,
		"doctype":    docType,
		"query":      query,
		"time":       results.TookInMillis,
		"total docs": results.Hits.TotalHits,
	}).Info("Query time in milliseconds.")
	if results.Hits.TotalHits > 0 {
		l.Info("Multiple documents found")
		return results
	}
	l.Info("No documents found")
	return nil
}

// FetchOne is wrapper method of ES search, except it returns only one document
func (g *Globals) FetchOne(index, docType string, query elastic.Query) *elastic.SearchResult {
	l := Gbl.Log
	esCtx := context.Background()
	page := 0
	size := 5
	results, err := g.ES.Search().
		Index(index).
		Query(query).
		From(page).
		Size(size).
		Pretty(true).
		Do(esCtx)
	if err != nil {
		l.WithFields(logrus.Fields{
			"error":   err,
			"index":   index,
			"doctype": docType,
			"query":   query,
		}).Error("Error fetching document")
	}
	fmt.Printf("Query took %d milliseconds\n", results.TookInMillis)
	l.WithFields(logrus.Fields{
		"index":      index,
		"doctype":    docType,
		"query":      query,
		"time":       results.TookInMillis,
		"total docs": results.Hits.TotalHits,
	}).Info("Query time in milliseconds.")
	if results.Hits.TotalHits == 1 {
		l.Info("Exactly one document found")
		return results
	} else if results.Hits.TotalHits > 1 {
		// Multiple hits found
		l.Info("More than one document found")
		return nil
	}
	// No hits
	l.Info("No documents found")
	return nil
}

// CreateIndex method of ES, creates a fresh index given index name and mapping
func (g *Globals) CreateIndex(index string, mapping interface{}, forceCreate bool) bool {
	l := Gbl.Log
	created := false
	esCtx := context.Background()
	exists, err := g.ES.IndexExists(index).Do(esCtx)
	if err != nil {
		l.WithFields(logrus.Fields{
			"error":   err,
			"index":   index,
			"mapping": mapping,
		}).Error("Error while checking existing index")
	}
	if exists && forceCreate {
		status := g.DeleteIndex(index)
		if !status {
			l.WithFields(logrus.Fields{
				"index":   index,
				"mapping": mapping,
			}).Error("Error deleting existing index")
		} else {
			l.Info("Force creating new index. (Deleting existing index with same name)")
		}
	}
	_, err = g.ES.CreateIndex(index).BodyJson(mapping).Do(esCtx)
	if err != nil {
		l.WithFields(logrus.Fields{
			"error":   err,
			"index":   index,
			"mapping": mapping,
		}).Error("Error creating index")
	} else {
		l.WithFields(logrus.Fields{
			"index":   index,
			"mapping": mapping,
		}).Info("index created")
		created = true
	}
	return created
}

// DeleteIndex deletes the index
func (g *Globals) DeleteIndex(index string) bool {
	l := Gbl.Log
	deleted := false
	esCtx := context.Background()
	deleteIndex, err := g.ES.DeleteIndex(index).Do(esCtx)
	if err != nil {
		// Handle error
		l.WithFields(logrus.Fields{
			"error": err,
			"index": index,
		}).Error("Error deleting index")
	}
	if !deleteIndex.Acknowledged {
		// Not acknowledged
		l.WithFields(logrus.Fields{
			"error": err,
			"index": index,
		}).Error("index delete failed!")
	} else {
		l.WithFields(logrus.Fields{
			"index": index,
		}).Info("index deleted!")
		deleted = true
	}
	return deleted
}

// CreateAlias method of ES, creates an alias for index
func (g *Globals) CreateAlias(index, alias string, forceCreate bool) bool {
	l := Gbl.Log
	created := false
	esCtx := context.Background()
	indexExists, err := g.ES.IndexExists(index).Do(esCtx)
	if err != nil {
		l.WithFields(logrus.Fields{
			"error": err,
			"index": index,
		}).Error("Error while checking existing index")
	}
	if indexExists {
		res, err := g.GetAlias(index)
		if err != nil {
			l.WithFields(logrus.Fields{
				"error": err,
				"index": index,
			}).Error("Error getting alias for index")
		} else {
			if (len(res.IndicesByAlias(alias)) > 0) && !forceCreate {
				l.Info("Alias exists for index, do not force create")
				created = true
				return created
			} else {
				if (len(res.IndicesByAlias(alias)) > 0) && forceCreate {
					l.WithFields(logrus.Fields{
						"index": index,
						"alias": alias,
					}).Info("Force deleting existing alias for index")
					status, err := g.ES.Alias().
						Remove(index, alias).
						// Pretty(true).
						Do(esCtx)
					if err != nil {
						l.WithFields(logrus.Fields{
							"error": err,
							"index": index,
							"alias": alias,
						}).Error("Error deleting alias for index")
					}
					if !status.Acknowledged {
						l.WithFields(logrus.Fields{
							"index":        index,
							"alias":        alias,
							"status":       status,
							"acknowledged": status.Acknowledged,
						}).Error("Alias could not be removed!")
					} else {
						l.WithFields(logrus.Fields{
							"index":        index,
							"alias":        alias,
							"status":       status,
							"acknowledged": status.Acknowledged,
						}).Info("Existing alias removed!")
					}
				}
				_, err = g.ES.Alias().Add(index, alias).Do(esCtx)
				if err != nil {
					l.WithFields(logrus.Fields{
						"error": err,
						"index": index,
						"alias": alias,
					}).Error("Error creating alias for index")
				} else {
					l.WithFields(logrus.Fields{
						"index": index,
						"alias": alias,
					}).Info("alias for index created")
					created = true
				}
			}
		}

	}
	return created
}

// GetAlias retrieves aliases for the index
func (g *Globals) GetAlias(index string) (*elastic.AliasesResult, error) {
	esCtx := context.Background()
	res, err := g.ES.Aliases().Index(index).Do(esCtx)
	if err != nil {
		return nil, err
	}
	return res, nil
}
