package fmp

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// StatementType identifies an FMP financial statement endpoint.
type StatementType string

const (
	StatementIncome       StatementType = "income"
	StatementBalanceSheet StatementType = "balance"
	StatementCashFlow     StatementType = "cash-flow"
)

// FinancialStatement is one raw FMP financial statement row.
type FinancialStatement map[string]any

// StockStatementsRequest describes a financial statement request.
type StockStatementsRequest struct {
	Symbol    string
	Statement StatementType
	Period    string
	Limit     int
}

// StockStatements returns financial statements for one stock symbol.
func (c *Client) StockStatements(ctx context.Context, request StockStatementsRequest) ([]FinancialStatement, error) {
	symbol := strings.ToUpper(strings.TrimSpace(request.Symbol))
	if symbol == "" {
		return nil, fmt.Errorf("symbol is required")
	}

	path, err := statementPath(request.Statement)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("symbol", symbol)
	if request.Period != "" {
		query.Set("period", request.Period)
	}
	if request.Limit > 0 {
		query.Set("limit", strconv.Itoa(request.Limit))
	}

	var statements []FinancialStatement
	if err := c.get(ctx, path, query, &statements); err != nil {
		return nil, err
	}
	return statements, nil
}

func statementPath(statement StatementType) (string, error) {
	switch statement {
	case StatementIncome:
		return "/income-statement", nil
	case StatementBalanceSheet:
		return "/balance-sheet-statement", nil
	case StatementCashFlow:
		return "/cash-flow-statement", nil
	default:
		return "", fmt.Errorf("unsupported statement type %q", statement)
	}
}
