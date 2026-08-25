package services

import (
	"strings"
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
)

// noteService จัดการ business logic ของ note ก่อนส่งต่อไป repository
// input: สร้างจาก NewNoteServiceImpl พร้อม NoteRepository
// output: service ที่ทำ create/list/search/update/delete note โดยผูกข้อมูลกับ userID
type noteService struct {
	repo repoPort.NoteRepository
}

var _ servicePort.NoteService = (*noteService)(nil)

// NewNoteServiceImpl สร้าง note service implementation
// input: repo repository สำหรับอ่าน/เขียน note
// output: *noteService ที่พร้อมถูกใช้ผ่าน NoteService interface
func NewNoteServiceImpl(repo repoPort.NoteRepository) *noteService {
	return &noteService{repo: repo}
}

// Create สร้าง note ใหม่ให้ user และเติม timestamp
// input: userID เจ้าของ note, note ข้อมูล note ที่รับจาก controller/assistant
// output: *Note ที่บันทึกแล้ว หรือ error ถ้า repository บันทึกไม่สำเร็จ
func (s *noteService) Create(userID uint, note *entities.Note) (*entities.Note, error) {
	now := time.Now()
	note.UserID = userID
	note.CreatedAt = now
	note.UpdatedAt = now
	return s.repo.Create(note)
}

// List ดึง note ของ user แบบแบ่งหน้า
// input: userID เจ้าของข้อมูล, limit จำนวนที่ต้องการ, offset ตำแหน่งเริ่มต้น
// output: []*Note รายการ note หรือ error จาก repository
func (s *noteService) List(userID uint, limit, offset int) ([]*entities.Note, error) {
	return s.repo.ListByUserID(userID, limit, offset)
}

// Search ค้น note จาก content หรือ fallback เป็น list เมื่อ query ว่าง
// input: userID เจ้าของข้อมูล, query คำค้นหา, limit จำนวนสูงสุด
// output: []*Note ที่ตรงกับคำค้น หรือรายการล่าสุดเมื่อ query ว่าง
func (s *noteService) Search(userID uint, query string, limit int) ([]*entities.Note, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return s.repo.ListByUserID(userID, limit, 0)
	}
	return s.repo.SearchByContent(userID, query, limit)
}

// Update แก้ไข note โดยตรวจ ownership ก่อนบันทึก
// input: userID เจ้าของข้อมูล, id note ที่ต้องการแก้, note ค่าใหม่
// output: *Note ที่แก้แล้ว หรือ ErrForbidden ถ้า note ไม่ใช่ของ user นี้
func (s *noteService) Update(userID, id uint, note *entities.Note) (*entities.Note, error) {
	current, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if current.UserID != userID {
		return nil, servicePort.ErrForbidden
	}

	note.ID = id
	note.UserID = userID
	note.CreatedAt = current.CreatedAt
	note.UpdatedAt = time.Now()
	return s.repo.Update(note)
}

// Delete ลบ note โดยตรวจ ownership ก่อนลบ
// input: userID เจ้าของข้อมูล, id note ที่ต้องการลบ
// output: error nil เมื่อลบสำเร็จ หรือ ErrForbidden ถ้า note ไม่ใช่ของ user นี้
func (s *noteService) Delete(userID, id uint) error {
	current, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if current.UserID != userID {
		return servicePort.ErrForbidden
	}
	return s.repo.Delete(id)
}
