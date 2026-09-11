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
	return p.db.Create(payment).Error
}

func (p *PaymentRepo) GetByID(id uint) (*domain.Payment, error) {
	var pay domain.Payment
	if err := p.db.Where("id = ?", id).First(&pay).Error; err != nil {
		return nil, err
	}
	return &pay, nil
}

func (p *PaymentRepo) GetByRideID(rideID uint) (*domain.Payment, error) {
	var pay domain.Payment
	if err := p.db.Where("ride_id = ?", rideID).First(&pay).Error; err != nil {
		return nil, err
	}
	return &pay, nil
}

func (p *PaymentRepo) Update(payment *domain.Payment) error {
	return p.db.Save(payment).Error
}

func (p *PaymentRepo) UpdateStatus(paymentID uint, status string) error {
	return p.db.Model(&domain.Payment{}).Where("id = ?", paymentID).Update("status", status).Error
}

func (p *PaymentRepo) Delete(id uint) error {
	return p.db.Where("id = ?", id).Delete(&domain.Payment{}).Error
}
