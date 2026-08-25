package services

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
)

// expenseService จัดการ business logic ของรายจ่ายก่อนส่งต่อไป repository
// input: สร้างจาก NewExpenseServiceImpl พร้อม ExpenseRepository
// output: service ที่ทำ create/list/summary/update/delete รายจ่ายโดยผูกข้อมูลกับ userID
type expenseService struct {
	repo repoPort.ExpenseRepository
}

var _ servicePort.ExpenseService = (*expenseService)(nil)

// NewExpenseServiceImpl สร้าง expense service implementation
// input: repo repository สำหรับอ่าน/เขียนรายจ่าย
// output: *expenseService ที่พร้อมถูกใช้ผ่าน ExpenseService interface
func NewExpenseServiceImpl(repo repoPort.ExpenseRepository) *expenseService {
	return &expenseService{repo: repo}
}

// Create สร้างรายจ่ายใหม่ให้ user และเติมค่า default ที่จำเป็น
// input: userID เจ้าของรายจ่าย, expense ข้อมูลรายจ่ายจาก controller/assistant
// output: *Expense ที่บันทึกแล้ว หรือ error ถ้า repository บันทึกไม่สำเร็จ
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

// List ดึงรายจ่ายของ user แบบแบ่งหน้า
// input: userID เจ้าของข้อมูล, limit จำนวนที่ต้องการ, offset ตำแหน่งเริ่มต้น
// output: []*Expense รายการรายจ่าย หรือ error จาก repository
func (s *expenseService) List(userID uint, limit, offset int) ([]*entities.Expense, error) {
	return s.repo.ListByUserID(userID, limit, offset)
}

// SummaryByPeriod สรุปยอดรายจ่ายตามช่วงเวลา
// input: userID เจ้าของข้อมูล, start เวลาเริ่มช่วง, end เวลาสิ้นสุดแบบ exclusive
// output: *ExpenseSummary ที่มี total/currency/start/end หรือ error จาก repository
func (s *expenseService) SummaryByPeriod(userID uint, start, end time.Time) (*servicePort.ExpenseSummary, error) {
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

// SummaryByMonth สรุปยอดรายจ่ายรายเดือน
// input: userID เจ้าของข้อมูล, year ปี, month เดือน, loc timezone สำหรับคำนวณขอบเขตเดือน
// output: *ExpenseSummary ของเดือนนั้น หรือ error จาก repository
func (s *expenseService) SummaryByMonth(userID uint, year int, month time.Month, loc *time.Location) (*servicePort.ExpenseSummary, error) {
	if loc == nil {
		loc = time.Local
	}
	start := time.Date(year, month, 1, 0, 0, 0, 0, loc)
	end := start.AddDate(0, 1, 0)
	return s.SummaryByPeriod(userID, start, end)
}

// Update แก้ไขรายจ่ายโดยตรวจ ownership ก่อนบันทึก
// input: userID เจ้าของข้อมูล, id รายจ่ายที่ต้องการแก้, expense ค่าใหม่
// output: *Expense ที่แก้แล้ว หรือ ErrForbidden ถ้ารายการไม่ใช่ของ user นี้
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

// Delete ลบรายจ่ายโดยตรวจ ownership ก่อนลบ
// input: userID เจ้าของข้อมูล, id รายจ่ายที่ต้องการลบ
// output: error nil เมื่อลบสำเร็จ หรือ ErrForbidden ถ้ารายการไม่ใช่ของ user นี้
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
