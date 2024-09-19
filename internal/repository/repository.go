package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/captainlettuce/sales-stats/pkg/types"
	"time"
)

// PostgresDB is a subset of *slqx.DB
type PostgresDB interface {
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type PostgresRepository struct {
	db PostgresDB
}

var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func New(db PostgresDB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) AggregateSalesPrice(ctx context.Context, c types.AggregationRequest) ([]types.AggregationResult, error) {
	if c.Filter.OrderDate == nil {
		return nil, errors.Join(types.ErrInvalidArgument, errors.New("OrderDate is required"))
	}

	if len(c.AggregationMethods) == 0 {
		return nil, errors.Join(types.ErrInvalidArgument, errors.New("at least one AggregationMethodType required"))
	}

	q := sq.
		Select("order_date::DATE as od").
		From("orders").
		GroupBy("od").
		OrderBy("od ASC")

	var err error
	q, err = applyAggregations(q, "price", c.AggregationMethods)
	if err != nil {
		return nil, fmt.Errorf("%w applying aggregations %w", types.ErrInvalidArgument, err)
	}

	q = applyDateRange(q, "order_date", c.Filter.OrderDate)

	q = applyFilter(q, c.Filter)

	q = applyGrouping(q, c.GroupBy)

	q.OrderBy("order_date ASC")
	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w creating sql query and args %w", types.ErrUnexpectedError, err)
	}

	var res []types.AggregationResult
	err = r.db.SelectContext(ctx, &res, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w querying db %w", types.ErrUnexpectedError, err)
	}

	return res, nil
}

func applyFilter(q squirrel.SelectBuilder, filter types.FilterCriteria) squirrel.SelectBuilder {
	if len(filter.OrderStatus) > 0 {
		q = whereOrWhereIn(q, "order_status", filter.OrderStatus)
	}

	if len(filter.Make) > 0 {
		q = whereOrWhereIn(q, "make", filter.Make)
	}

	if len(filter.Model) > 0 {
		q = whereOrWhereIn(q, "model", filter.Model)
	}

	if len(filter.Version) > 0 {
		q = whereOrWhereIn(q, "version", filter.Version)
	}

	if len(filter.Color) > 0 {
		q = whereOrWhereIn(q, "color", filter.Color)
	}

	if filter.ModelYear != nil {
		q = applyIntRange(q, "model_year", *filter.ModelYear)
	}

	if filter.MileageKilometers != nil {
		q = applyIntRange(q, "mileage_kilometers", *filter.MileageKilometers)
	}

	return q
}

func whereOrWhereIn[T any](q squirrel.SelectBuilder, col string, filter []T) squirrel.SelectBuilder {
	switch len(filter) {
	case 0:
		break
	case 1:
		q = q.Where(squirrel.Eq{col: filter[0]})
	default:
		q = q.Where(squirrel.Eq{col: filter})
	}
	return q
}

func applyDateRange(q squirrel.SelectBuilder, col string, r *types.DateRange) squirrel.SelectBuilder {
	switch {
	case r.Before != nil && r.After != nil:
		return q.Where(col+"::DATE BETWEEN ?::DATE AND ?::DATE", r.After.Format(time.DateOnly), r.Before.Format(time.DateOnly))
	case r.After != nil:
		return q.Where(col+"::DATE > ?::DATE", r.After.Format(time.DateOnly))
	case r.Before != nil:
		return q.Where(col+"::DATE < ?::DATE", r.Before.Format(time.DateOnly))
	default:
		return q
	}
}

func applyIntRange(q squirrel.SelectBuilder, col string, r types.IntegerRange) squirrel.SelectBuilder {
	if r.Over != nil {
		q = q.Where(col+" > ?", r.Over)
	}
	if r.Under != nil {
		q = q.Where(col+" < ?", r.Under)
	}

	return q
}

func applyGrouping(q squirrel.SelectBuilder, g types.GroupingCriteria) squirrel.SelectBuilder {
	selectAndGroup := func(q squirrel.SelectBuilder, col string) squirrel.SelectBuilder {
		return q.Columns(col).GroupBy(col)
	}

	if g.Make {
		q = selectAndGroup(q, "make")
	}

	if g.Model {
		q = selectAndGroup(q, "model")
	}

	if g.Version {
		q = selectAndGroup(q, "version")
	}

	if g.Color {
		q = selectAndGroup(q, "color")
	}

	if g.ModelYear {
		q = selectAndGroup(q, "model_year")
	}

	return q
}

func applyAggregations(q squirrel.SelectBuilder, col string, methods []types.AggregationMethodType) (squirrel.SelectBuilder, error) {
	for _, m := range methods {
		var queryString string
		switch m {
		case types.AggregationMethodAvg:
			// It's assumed that we don't care about decimals here,
			// so we let postgres do the rounding and avoid floats in go
			queryString = "ROUND(AVG(" + col + "))"
		case types.AggregationMethodMax:
			queryString = "MAX(" + col + ")"
		case types.AggregationMethodMin:
			queryString = "MIN(" + col + ")"
		case types.AggregationMethodSum:
			queryString = "SUM(" + col + ")"
		default:
			return q, fmt.Errorf("unknown aggregation method '%s'", m)
		}
		q = q.Columns(queryString + " as " + string(m))
	}

	return q, nil
}
