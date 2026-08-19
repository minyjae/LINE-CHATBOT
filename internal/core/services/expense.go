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

func (s *expenseService) Create(userID uint, expense *entities.Expense) (*entities.Expense, error) {
	now := time.Now()
	expense.UserID = userID
	expense.Currency = defaultString(expense.Currency, "THB")
	expense.Category = defaultString(expense.Category, "uncategorized")
	if expense.SpentAt.IsZero() {
		expense.SpentAt = now
	}
	expense.CreatedAt = now
	expense.UpdatedAt = now
	return s.repo.Create(expense)
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

func (s *expenseService) Update(userID, id uint, expense *entities.Expense) (*entities.Expense, error) {
	current, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if current.UserID != userID {
		return nil, servicePort.ErrForbidden
	}

	expense.ID = id
	expense.UserID = userID
	expense.CreatedAt = current.CreatedAt
	expense.UpdatedAt = time.Now()
	if expense.Currency == "" {
		expense.Currency = current.Currency
	}
	if expense.Category == "" {
		expense.Category = current.Category
	}
	if expense.SpentAt.IsZero() {
		expense.SpentAt = current.SpentAt
	}
	return s.repo.Update(expense)
}

func (s *expenseService) Delete(userID, id uint) error {
	current, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if current.UserID != userID {
		return servicePort.ErrForbidden
	}
	return s.repo.Delete(id)
}
