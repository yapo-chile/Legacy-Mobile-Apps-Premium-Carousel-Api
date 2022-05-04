package repository

import (
	"fmt"

	"gitlab.com/yapo_team/legacy/mobile-apps/premium-carousel-api/pkg/domain"
	"gitlab.com/yapo_team/legacy/mobile-apps/premium-carousel-api/pkg/usecases"
)

// productRepo holds connections to get user products
type purchaseRepo struct {
	handler DbHandler
}

// MakePurchaseRepository creates a new instance of PurchaseRepository
func MakePurchaseRepository(handler DbHandler) usecases.PurchaseRepository {
	return &purchaseRepo{
		handler: handler,
	}
}

// CreatePurchase creates a new purchase
func (repo *purchaseRepo) CreatePurchase(purchaseNumber, price int,
	purchaseType domain.PurchaseType) (purchase domain.Purchase, err error) {
	result, err := repo.handler.Query(
		`INSERT INTO purchase(purchase_number, price, purchase_type)
			VALUES (
				$1, $2, $3
			) RETURNING id, created_at, purchase_status`, purchaseNumber, price, purchaseType)
	if err != nil {
		return domain.Purchase{}, err
	}
	defer result.Close()
	if result.Next() {
		result.Scan(&purchase.ID, &purchase.CreatedAt, &purchase.Status)
	} else {
		return domain.Purchase{},
			fmt.Errorf("next error: getting purchaseID from database")
	}
	return domain.Purchase{
		ID:        purchase.ID,
		Number:    purchaseNumber,
		Price:     price,
		Type:      purchaseType,
		Status:    purchase.Status,
		CreatedAt: purchase.CreatedAt,
	}, nil
}

// AcceptePurchase changes the purchase status to Accepted
func (repo *purchaseRepo) AcceptPurchase(purchase domain.Purchase) (domain.Purchase, error) {
	if err := repo.setStatus(purchase.ID, domain.AcceptedPurchase); err != nil {
		return domain.Purchase{}, err
	}
	purchase.Status = domain.AcceptedPurchase
	return purchase, nil
}

// setStatus sets the purchase status
func (repo *purchaseRepo) setStatus(purchaseID int, status domain.PurchaseStatus) error {
	result, err := repo.handler.
		Query(
			`UPDATE purchase SET purchase_status=$1 WHERE id=$2`,
			status,
			purchaseID,
		)
	if err != nil {
		return err
	}
	return result.Close()
}
