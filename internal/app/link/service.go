package link

import (
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"github.com/teris-io/shortid"
	"github.com/wvoliveira/corgi/internal/pkg/common"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"github.com/wvoliveira/corgi/internal/pkg/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	Add(*gin.Context, model.Link) (model.Link, error)
	FindByID(*gin.Context, string, string) (model.Link, error)
	FindAll(*gin.Context, findAllRequest) (int64, int, []model.Link, error)
	Update(*gin.Context, model.Link) (model.Link, error)
	Delete(*gin.Context, string, string) (err error)

	NewHTTP(*gin.RouterGroup)
	HTTPAdd(*gin.Context)
	HTTPFindByID(*gin.Context)
	HTTPFindAll(*gin.Context)
	HTTPUpdate(*gin.Context)
	HTTPDelete(*gin.Context)
}

type service struct {
	db *sql.DB
}

// NewService creates a new authentication service.
func NewService(db *sql.DB) Service {
	return service{db}
}

// Add create a new shortener link.
func (s service) Add(c *gin.Context, link model.Link) (m model.Link, err error) {
	l := logger.Logger(c)

	if err = checkLink(link); err != nil {
		l.Error().Caller().Msg(err.Error())
		return
	}

	// If user is anonymous, create a random ID and blank another fields.
	if link.UserID == "anonymous" {
		sid, _ := shortid.New(1, shortid.DefaultABC, 2342)
		link.Keyword, _ = sid.Generate()
	}

	if link.UserID != "anonymous" {
		if link.Keyword == "" {
			sid, _ := shortid.New(1, shortid.DefaultABC, 2342)
			link.Keyword, _ = sid.Generate()
		}
	}

	query := "SELECT * FROM links WHERE domain = $1 AND keyword = $2 LIMIT 1"
	rows, err := s.db.QueryContext(c, query, link.Domain, link.Keyword)

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&m)

		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return
		}

		log.Info().Caller().Msg(fmt.Sprintf("link ID: %s", m.ID))
	}

	if m.ID != "" {
		l.Warn().Caller().Msg("domain with keyword already exists")
		return m, e.ErrAlreadyExists
	}

	query = `
		INSERT INTO links(id, created_at, domain, keyword, url, title, active, user_id) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err = s.db.ExecContext(
		c,
		query,
		ulid.Make(),
		time.Now(),
		link.Domain,
		link.Keyword,
		link.URL,
		link.Title,
		true,
		link.UserID,
	)

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	return
}

// FindByID get a shortener link from ID.
func (s service) FindByID(c *gin.Context, linkID, userID string) (link model.Link, err error) {
	log := logger.Logger(c)

	query := "SELECT * FROM links WHERE id = $1 AND user_id = $2 LIMIT 1"
	rows, err := s.db.QueryContext(c, query, linkID, userID)

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&link)

		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return
		}

		log.Info().Caller().Msg(fmt.Sprintf("link ID: %s", link.ID))
		return
	}

	log.Warn().Caller().Msg("link not found")
	return link, e.ErrLinkNotFound
}

// FindAll get a list of links from database.
func (s service) FindAll(c *gin.Context, r findAllRequest) (total int64, pages int, links []model.Link, err error) {
	log := logger.Logger(c)

	query := fmt.Sprintf("SELECT count(*) FROM links WHERE user_id = %s", r.UserID)
	// rows, err := s.db.QueryContext(c, query, r.UserID, r.Offset, r.Limit)

	if len(r.SearchText) >= 3 {
		query = query + fmt.Sprintf(" AND domain LIKE %%%[1]s%% OR keyword LIKE %%%[1]s%%", r.SearchText)
	}

	domain, keyword := common.SplitURL(r.ShortenedURL)
	if domain != "" && keyword != "" {
		query = query + fmt.Sprintf(" AND domain = %s AND keyword = %s", domain, keyword)
	}

	// TODO: add order by another field.
	query = query + fmt.Sprintf(" ORDER BY created_at LIMIT %d OFFSET %d", r.Limit, r.Offset)

	rows, err := s.db.QueryContext(c, query)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	defer rows.Close()
	var link model.Link

	if rows.Next() {
		err = rows.Scan(&link)

		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return
		}

		links = append(links, link)
		log.Info().Caller().Msg(fmt.Sprintf("link ID: %s", link.ID))
	}

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return total, pages, links, e.ErrInternalServerError
	}

	pages = int(math.Ceil(float64(total) / float64(r.Limit)))

	return
}

// Update change specific link by ID.
func (s service) Update(c *gin.Context, link model.Link) (m model.Link, err error) {
	log := logger.Logger(c)

	link.UpdatedAt = time.Now()

	err = s.db.Model(&model.Link{}).
		Where("id = ? AND user_id = ?", link.ID, link.UserID).
		Updates(&link).
		First(&m).Error

	if err == gorm.ErrRecordNotFound {
		log.Info().Caller().Msg(err.Error())
		return m, e.ErrLinkNotFound
	}

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return m, e.ErrInternalServerError
	}

	return
}

// Delete delete a link by ID.
func (s service) Delete(c *gin.Context, linkID, userID string) (err error) {
	log := logger.Logger(c)

	err = s.db.
		Model(&model.Link{}).
		Clauses(clause.Returning{}).
		Where("id = ? AND user_id = ?", linkID, userID).
		Delete(&model.Link{ID: linkID, UserID: userID}).Error

	if err == gorm.ErrRecordNotFound {
		log.Info().Caller().Msg(err.Error())
		return e.ErrLinkNotFound
	}

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return e.ErrInternalServerError
	}

	return
}
