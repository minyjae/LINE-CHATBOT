package services

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
)

type incomeService struct {
	repo repoPort.IncomeRepository
}

var _ servicePort.IncomeService = (*incomeService)(nil)

func NewIncomeServiceImpl(repo repoPort.IncomeRepository) *incomeService {
	return &incomeService{repo: repo}
}

func (s *incomeService) Create(userID uint, income *entities.Income) (*entities.Income, error) {
	now := time.Now()
	income.UserID = userID
	income.Currency = defaultString(income.Currency, "THB")
	income.Category = defaultString(income.Category, "uncategorized")
	if income.ReceivedAt.IsZero() {
		income.ReceivedAt = now
	}
	income.CreatedAt = now
	income.UpdatedAt = now
	return s.repo.Create(income)
}

func (s *incomeService) List(userID uint, limit, offset int) ([]*entities.Income, error) {
	return s.repo.ListByUserID(userID, limit, offset)
}

func (s *incomeService) SummaryByPeriod(userID uint, start, end time.Time) (*servicePort.IncomeSummary, error) {
	total, err := s.repo.SumByReceivedAtBetween(userID, start, end)
	if err != nil {
		return nil, err
	}
	return &servicePort.IncomeSummary{
		Total:    total,
		Currency: "THB",
		Start:    start,
		End:      end,
	}, nil
}

func (s *incomeService) SummaryByMonth(userID uint, year int, month time.Month, loc *time.Location) (*servicePort.IncomeSummary, error) {
	if loc == nil {
		loc = time.Local
	}
	start := time.Date(year, month, 1, 0, 0, 0, 0, loc)
	end := start.AddDate(0, 1, 0)
	return s.SummaryByPeriod(userID, start, end)
}

func (s *incomeService) Update(userID, id uint, income *entities.Income) (*entities.Income, error) {
	current, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if current.UserID != userID {
		return nil, servicePort.ErrForbidden
	}

	income.ID = id
	income.UserID = userID
	income.CreatedAt = current.CreatedAt
	income.UpdatedAt = time.Now()
	if income.Currency == "" {
		income.Currency = current.Currency
	}
	if income.Category == "" {
		income.Category = current.Category
	}
	if income.ReceivedAt.IsZero() {
		income.ReceivedAt = current.ReceivedAt
	}
	return s.repo.Update(income)
}

func (s *incomeService) Delete(userID, id uint) error {
	current, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if current.UserID != userID {
		return servicePort.ErrForbidden
	}
	return s.repo.Delete(id)
}
