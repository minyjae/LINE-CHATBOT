package services

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
)

type expenseService struct {
	repo repoPort.ExpenseRepository
}

var _ servicePort.ExpenseService = (*expenseService)(nil)

func NewExpenseServiceImpl(repo repoPort.ExpenseRepository) *expenseService {
	return &expenseService{repo: repo}
}

func (s *expenseService) List(userID uint, limit, offset int) ([]*entities.Expense, error) {
	return s.repo.ListByUserID(userID, limit, offset)
}

func (s *expenseService) SummaryByMonth(userID uint, year int, month time.Month, loc *time.Location) (*servicePort.ExpenseSummary, error) {
	if loc == nil {
		loc = time.Local
	}
	start := time.Date(year, month, 1, 0, 0, 0, 0, loc)
	end := start.AddDate(0, 1, 0)
	total, err := s.repo.SumBySpentAtBetween(userID, start, end)
	if err != nil {
		return nil, err
	}

	return &servicePort.ExpenseSummary{
		Total:    total,
		Currency: "THB",
		Start:    start,
		End:      end,
	}, nil
}
