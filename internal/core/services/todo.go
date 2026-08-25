package services

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
)

// todoService จัดการ business logic ของ todo ก่อนส่งต่อไป repository
// input: สร้างจาก NewTodoServiceImpl พร้อม TodoRepository
// output: service ที่ทำ create/list/update/delete todo โดยผูกข้อมูลกับ userID
type todoService struct {
	repo repoPort.TodoRepository
}

var _ servicePort.TodoService = (*todoService)(nil)

// NewTodoServiceImpl สร้าง todo service implementation
// input: repo repository สำหรับอ่าน/เขียน todo
// output: *todoService ที่พร้อมถูกใช้ผ่าน TodoService interface
func NewTodoServiceImpl(repo repoPort.TodoRepository) *todoService {
	return &todoService{repo: repo}
}

// Create สร้าง todo ใหม่ให้ user และเติมค่า default ที่จำเป็น
// input: userID เจ้าของ todo, todo ข้อมูล todo ที่รับจาก controller/assistant
// output: *Todo ที่บันทึกแล้ว หรือ error ถ้า repository บันทึกไม่สำเร็จ
func (s *todoService) Create(userID uint, todo *entities.Todo) (*entities.Todo, error) {
	now := time.Now()
	todo.UserID = userID
	todo.Status = defaultString(todo.Status, "pending")
	todo.Priority = defaultString(todo.Priority, "normal")
	todo.CreatedAt = now
	todo.UpdatedAt = now
	return s.repo.Create(todo)
}

// List ดึง todo ของ user แบบแบ่งหน้า
// input: userID เจ้าของข้อมูล, limit จำนวนที่ต้องการ, offset ตำแหน่งเริ่มต้น
// output: []*Todo รายการ todo หรือ error จาก repository
func (s *todoService) List(userID uint, limit, offset int) ([]*entities.Todo, error) {
	return s.repo.ListByUserID(userID, limit, offset)
}

// ListPending ดึง todo ที่ยัง pending ของ user
// input: userID เจ้าของข้อมูล
// output: []*Todo เฉพาะรายการที่ยังไม่เสร็จ หรือ error จาก repository
func (s *todoService) ListPending(userID uint) ([]*entities.Todo, error) {
	return s.repo.ListPendingByUserID(userID)
}

// Update แก้ไข todo โดยตรวจ ownership ก่อนบันทึก
// input: userID เจ้าของข้อมูล, id todo ที่ต้องการแก้, todo ค่าใหม่
// output: *Todo ที่แก้แล้ว หรือ ErrForbidden ถ้า todo ไม่ใช่ของ user นี้
func (s *todoService) Update(userID, id uint, todo *entities.Todo) (*entities.Todo, error) {
	current, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if current.UserID != userID {
		return nil, servicePort.ErrForbidden
	}

	todo.ID = id
	todo.UserID = userID
	todo.CreatedAt = current.CreatedAt
	todo.UpdatedAt = time.Now()
	if todo.Status == "" {
		todo.Status = current.Status
	}
	if todo.Priority == "" {
		todo.Priority = current.Priority
	}
	return s.repo.Update(todo)
}

// Delete ลบ todo โดยตรวจ ownership ก่อนลบ
// input: userID เจ้าของข้อมูล, id todo ที่ต้องการลบ
// output: error nil เมื่อลบสำเร็จ หรือ ErrForbidden ถ้า todo ไม่ใช่ของ user นี้
func (s *todoService) Delete(userID, id uint) error {
	current, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if current.UserID != userID {
		return servicePort.ErrForbidden
	}
	return s.repo.Delete(id)
}
