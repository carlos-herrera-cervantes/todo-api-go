package models

import (
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type QueryParameters struct {
	Page     int    `schema:"page"`
	PageSize int    `schema:"pageSize"`
	Sort     string `schema:"sort"`
	Filter   string `schema:"filter"`
	Relation string `schema:"relation"`
}

// Takes the string with relationships and transforms it to Bson document
func parseRelationships(model string, relation string, pipeline mongo.Pipeline) mongo.Pipeline {
	if len(relation) == 0 {
		return pipeline
	}

	relationships := GetRelationModel(model)
	lookupStage := relationships[relation]["lookupStage"]
	unwindStage := relationships[relation]["unwindStage"]

	if lookupStage != nil {
		pipeline = append(pipeline, lookupStage)
	}

	if unwindStage != nil {
		pipeline = append(pipeline, unwindStage)
	}

	return pipeline
}

// Takes the string with filter and transform it to bson document
func parseFilter(filter string) bson.D {
	if len(filter) == 0 {
		return bson.D{}
	}

	splitedByComma := strings.Split(filter, ",")
	var matchStage bson.D

	for i := range splitedByComma {
		splitedByEqual := strings.Split(splitedByComma[i], "=")
		value, error := primitive.ObjectIDFromHex(splitedByEqual[1])

		if error != nil {
			matchStage = append(matchStage, bson.E{Key: splitedByEqual[0], Value: splitedByEqual[1]})
		} else {
			matchStage = append(matchStage, bson.E{Key: splitedByEqual[0], Value: value})
		}
	}

	return matchStage
}

// Takes the page and page size and trasform they to bson document
func parsePagination(page int, pageSize int, pipeline mongo.Pipeline) mongo.Pipeline {
	var offset int
	var limit int

	if page <= 1 {
		offset = 0
	} else {
		offset = page - 1
	}

	if pageSize > 100 {
		limit = 100
	} else if pageSize < 1 {
		limit = 10
	} else {
		limit = pageSize
	}

	pipeline = append(pipeline, bson.D{{Key: "$skip", Value: offset * limit}})
	pipeline = append(pipeline, bson.D{{Key: "$limit", Value: limit}})

	return pipeline
}

// Takes the sorting criteria and transform it to bson document
func parseSort(sort string, pipeline mongo.Pipeline) mongo.Pipeline {
	var sortStage bson.D

	if len(sort) == 0 {
		sortStage = bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}}
	} else {
		splitedByComma := strings.Split(sort, ",")

		if strings.Contains(splitedByComma[0], "-") {
			splitedByHyphen := strings.Split(splitedByComma[0], "-")
			sortStage = bson.D{{Key: "$sort", Value: bson.D{{Key: splitedByHyphen[1], Value: -1}}}}
		} else {
			sortStage = bson.D{{Key: "$sort", Value: bson.D{{Key: splitedByComma[0], Value: 1}}}}
		}
	}

	pipeline = append(pipeline, sortStage)

	return pipeline
}

// Sets the values for query parameters
func (query *QueryParameters) SetValues(model string) mongo.Pipeline {
	pipeline := make([]bson.D, 0)
	pipeline = parseRelationships(model, query.Relation, pipeline)
	matchStage := parseFilter(query.Filter)

	pipeline = append(pipeline, bson.D{{Key: "$match", Value: matchStage}})
	pipeline = parseSort(query.Sort, pipeline)
	pipeline = parsePagination(query.Page, query.PageSize, pipeline)

	return pipeline
}

// Converts the string filter to bson D
func (query *QueryParameters) SetFilter() bson.D {
	return parseFilter(query.Filter)
}

// Returns the page and page size
func (query *QueryParameters) SetPagination() (int, int) {
	var offset int
	var limit int

	if query.Page <= 1 {
		offset = 1
	} else {
		offset = query.Page
	}

	if query.PageSize > 100 {
		limit = 100
	} else if query.PageSize < 1 {
		limit = 10
	} else {
		limit = query.PageSize
	}

	return offset, limit
}
