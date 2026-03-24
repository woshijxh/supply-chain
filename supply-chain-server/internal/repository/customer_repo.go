package repository

import (
	"supply-chain-server/internal/model"

	"gorm.io/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) Create(customer *model.Customer) error {
	return r.db.Create(customer).Error
}

func (r *CustomerRepository) GetByID(id uint) (*model.Customer, error) {
	var customer model.Customer
	err := r.db.First(&customer, id).Error
	return &customer, err
}

func (r *CustomerRepository) Update(customer *model.Customer) error {
	return r.db.Save(customer).Error
}

func (r *CustomerRepository) Delete(id uint) error {
	return r.db.Delete(&model.Customer{}, id).Error
}

func (r *CustomerRepository) List(page, pageSize int, keyword string) ([]model.Customer, int64, error) {
	var customers []model.Customer
	var total int64

	query := r.db.Model(&model.Customer{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR contact LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&customers).Error
	return customers, total, err
}

func (r *CustomerRepository) GetActiveCustomers() ([]model.Customer, error) {
	var customers []model.Customer
	err := r.db.Where("status = ?", 1).Find(&customers).Error
	return customers, err
}

func (r *CustomerRepository) ExistsByName(name string) bool {
	var count int64
	r.db.Model(&model.Customer{}).Where("name = ?", name).Count(&count)
	return count > 0
}