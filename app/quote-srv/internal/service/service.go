package service

import "wq-fotune-backend/app/quote-srv/internal/dao"

type QuoteService struct {
	dao *dao.Dao
}

func NewQuoteService() *QuoteService {
	handler := &QuoteService{dao: dao.New()}
	return handler
}