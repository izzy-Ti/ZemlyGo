package postgres

import (
	"github.com/izzy-Ti/ZemlyGo/internals/domain"
	"github.com/izzy-Ti/ZemlyGo/internals/repository/interfaces"
	"gorm.io/gorm"
)

type PaymentRepo struct {
	db *gorm.DB
}

func NewPaymentRepo(db *gorm.DB) interfaces.PaymentRepository {
	return &PaymentRepo{db: db}
}

func (p *PaymentRepo) Create(payment *domain.Payment) error {
	if err := p.db.Create(payment).Error; err != nil {
		return err
	}
	return nil
}
func (p *PaymentRepo) GetByID(id uint) (*domain.Payment, error) {
	var pay *domain.Payment
	if err := p.db.Where("id = ?", id).First(&pay).Error; err != nil {
		return nil, err
	}
	return pay, nil
}
func (p *PaymentRepo) GetByRideID(rideID uint) (*domain.Payment, error) {
	var pay *domain.Payment
	if err := p.db.Where("ride_id = ?", rideID).First(&pay).Error; err != nil {
		return nil, err
	}
	return pay, nil
}
func (p *PaymentRepo) Update(payment *domain.Payment) error {
	if err := p.db.Save(payment).Error; err != nil {
		return err
	}
	return nil
}
func (p *PaymentRepo) UpdateStatus(paymentID uint, status string) error {
	var pay *domain.Payment
	if err := p.db.Where("id = ?", paymentID).First(&pay).Error; err != nil {
		return err
	}
	pay.Status = status
	if err := p.db.Save(pay).Error; err != nil {
		return err
	}
	return nil
}
func (p *PaymentRepo) Delete(id uint) error {
	if err := p.db.Where("id = ?", id).Delete(&domain.Payment{}).Error; err != nil {
		return err
	}
	return nil
}
