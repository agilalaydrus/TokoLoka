package service

import "main.go/entity"

// ConvertToTransactionResponse - Mengubah Transaction menjadi TransactionResponse
func ConvertToTransactionResponse(transaction *entity.Transaction) entity.TransactionResponse {
	transactionResponse := entity.TransactionResponse{
		ID:                transaction.ID,
		UserID:            transaction.UserID,
		DestinationNumber: transaction.DestinationNumber,
		TotalPrice:        transaction.TotalPrice,
		Status:            transaction.Status,
		CreatedAt:         transaction.CreatedAt,
		UpdatedAt:         transaction.UpdatedAt,
		User: entity.UserSafeResponse{
			ID:       transaction.User.ID,
			Username: transaction.User.Username,
			Email:    transaction.User.Email,
			Role:     transaction.User.Role,
		},
	}

	for _, item := range transaction.Items {
		transactionResponse.Items = append(transactionResponse.Items, ConvertToTransactionItemResponse(item))
	}

	return transactionResponse
}

// ConvertToTransactionItemResponse - Mengubah TransactionItem menjadi TransactionItemResponse
func ConvertToTransactionItemResponse(item entity.TransactionItem) entity.TransactionItemResponse {
	return entity.TransactionItemResponse{
		ID:        item.ID,
		ProductID: item.ProductID,
		Quantity:  item.Quantity,
		Price:     item.Price,
		Product: entity.ProductResponse{
			ID:          item.Product.ID,
			Name:        item.Product.Name,
			Description: item.Product.Description,
			Price:       item.Product.Price,
			Stock:       item.Product.Stock,
			Category: entity.CategoryResponse{
				ID:          item.Product.Category.ID,
				Name:        item.Product.Category.Name,
				Description: item.Product.Category.Description,
			},
		},
	}
}
