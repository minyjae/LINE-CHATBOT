package services

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
)

// incomeService จัดการ business logic ของรายรับก่อนส่งต่อไป repository
// input: สร้างจาก NewIncomeServiceImpl พร้อม IncomeRepository
// output: service ที่ทำ create/list/summary/update/delete รายรับโดยผูกข้อมูลกับ userID
type incomeService struct {
	repo repoPort.IncomeRepository
}

var _ servicePort.IncomeService = (*incomeService)(nil)

// NewIncomeServiceImpl สร้าง income service implementation
// input: repo repository สำหรับอ่าน/เขียนรายรับ
// output: *incomeService ที่พร้อมถูกใช้ผ่าน IncomeService interface
func NewIncomeServiceImpl(repo repoPort.IncomeRepository) *incomeService {
	return &incomeService{repo: repo}
}

// Create สร้างรายรับใหม่ให้ user และเติมค่า default ที่จำเป็น
// input: userID เจ้าของรายรับ, income ข้อมูลรายรับจาก controller/assistant
// output: *Income ที่บันทึกแล้ว หรือ error ถ้า repository บันทึกไม่สำเร็จ
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

// List ดึงรายรับของ user แบบแบ่งหน้า
// input: userID เจ้าของข้อมูล, limit จำนวนที่ต้องการ, offset ตำแหน่งเริ่มต้น
// output: []*Income รายการรายรับ หรือ error จาก repository
func (s *incomeService) List(userID uint, limit, offset int) ([]*entities.Income, error) {
	return s.repo.ListByUserID(userID, limit, offset)
}

// SummaryByPeriod สรุปยอดรายรับตามช่วงเวลา
// input: userID เจ้าของข้อมูล, start เวลาเริ่มช่วง, end เวลาสิ้นสุดแบบ exclusive
// output: *IncomeSummary ที่มี total/currency/start/end หรือ error จาก repository
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

// SummaryByMonth สรุปยอดรายรับรายเดือน
// input: userID เจ้าของข้อมูล, year ปี, month เดือน, loc timezone สำหรับคำนวณขอบเขตเดือน
// output: *IncomeSummary ของเดือนนั้น หรือ error จาก repository
func (s *incomeService) SummaryByMonth(userID uint, year int, month time.Month, loc *time.Location) (*servicePort.IncomeSummary, error) {
	if loc == nil {
		loc = time.Local
	}
	start := time.Date(year, month, 1, 0, 0, 0, 0, loc)
	end := start.AddDate(0, 1, 0)
	return s.SummaryByPeriod(userID, start, end)
}

// Update แก้ไขรายรับโดยตรวจ ownership ก่อนบันทึก
// input: userID เจ้าของข้อมูล, id รายรับที่ต้องการแก้, income ค่าใหม่
// output: *Income ที่แก้แล้ว หรือ ErrForbidden ถ้ารายการไม่ใช่ของ user นี้
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

// Delete ลบรายรับโดยตรวจ ownership ก่อนลบ
// input: userID เจ้าของข้อมูล, id รายรับที่ต้องการลบ
// output: error nil เมื่อลบสำเร็จ หรือ ErrForbidden ถ้ารายการไม่ใช่ของ user นี้
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
