package quote

import (
	"context"
	"encoding/json"
	"github.com/MeM0rd/q-quote-svc/pkg/client/postgres"
	"github.com/MeM0rd/q-quote-svc/pkg/logger"
	quotePbServer "github.com/MeM0rd/q-quote-svc/pkg/pb/quote"
)

type Server struct {
	quotePbServer.UnimplementedQuotePbServiceServer
	Logger *logger.Logger
}

func (s *Server) GetList(ctx context.Context, in *quotePbServer.GetListRequest) (*quotePbServer.GetListResponse, error) {
	var quotes []Quote

	q := `SELECT id, user_id, title, text FROM quotes`

	rows, err := postgres.DB.Query(q)
	s.Logger.Infof("error query postgres: %v", err)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var quote Quote

		err := rows.Scan(&quote.Id, &quote.UserId, &quote.Title, &quote.Text)
		s.Logger.Infof("error rows scan: %v", err)
		if err != nil {
			return nil, err
		}

		quotes = append(quotes, quote)
	}

	j, err := json.Marshal(quotes)
	s.Logger.Infof("error marhsal quotes: %v", err)
	if err != nil {
		return nil, err
	}

	return &quotePbServer.GetListResponse{
		Quotes: j,
	}, nil
}

func (s *Server) Create(ctx context.Context, in *quotePbServer.CreateRequest) (*quotePbServer.CreateResponse, error) {
	quote := Quote{
		UserId: int(in.GetUserId()),
		Title:  in.GetTitle(),
		Text:   in.GetText(),
	}

	q := `INSERT INTO quotes (user_id, title, text) VALUES ($1, $2, $3) RETURNING id, user_id, title, text`

	err := postgres.DB.QueryRow(q, quote.UserId, quote.Title, quote.Text).Scan(&quote.Id, &quote.UserId, &quote.Title, &quote.Text)
	if err != nil {
		s.Logger.Infof("error postgres query: %v", err)
		return nil, err
	}

	j, err := json.Marshal(quote)
	if err != nil {
		s.Logger.Infof("error marshal: %v", err)
		return nil, err
	}

	return &quotePbServer.CreateResponse{
		Status: true,
		Quote:  j,
		Err:    "",
	}, nil
}

func (s *Server) Delete(ctx context.Context, in *quotePbServer.DeleteRequest) (*quotePbServer.DeleteResponse, error) {
	const adminLevel = 2
	var userLevel, quoteOwner int
	userId := int(in.GetUserId())
	quoteId := int(in.GetQuoteId())

	q := `
	SELECT r.level  FROM roles r
	INNER JOIN users u ON r.id = u.role_id
	WHERE u.id = $1 
	`

	err := postgres.DB.QueryRow(q, userId).Scan(&userLevel)
	if err != nil {
		s.Logger.Infof("error checking user lvl in delete quote func: %v", err)
		return nil, err
	}

	if userLevel < adminLevel {
		q := `SELECT user_id FROM quotes WHERE  id = $1`
		s.Logger.Infof("%v,  %v", q, quoteId)
		err = postgres.DB.QueryRow(q, quoteId).Scan(&quoteOwner)
		if err != nil {
			s.Logger.Infof("error postgres query: %v", err)
			return nil, err
		}

		if quoteOwner != userId {
			s.Logger.Infof("Not enough rights")
			return &quotePbServer.DeleteResponse{
				Status: false,
				Err:    "Not enough rights",
			}, nil
		}
	}

	q = `DELETE FROM quotes WHERE id = $1`

	err = postgres.DB.QueryRow(q, quoteId).Err()
	if err != nil {
		s.Logger.Infof("error deleting quote: %v", err)
		return nil, err
	}

	return &quotePbServer.DeleteResponse{
		Status: true,
		Msg:    "Success",
		Err:    "",
	}, nil
}
