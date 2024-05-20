package service

import (
	"alc/model/store"
	"context"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type Comment struct {
	db *pgxpool.Pool
}

func NewCommentService(db *pgxpool.Pool) Comment {
	return Comment{
		db: db,
	}
}

func (cs Comment) GetComments(item store.Item, limit int, page int) ([]store.ItemComment, error) {
	sql := `SELECT ic.id, ic.title, ic.message, ic.rating, ic.up_votes, ic.down_votes,
	ic.is_edited, ic.created_at, ic.edited_at,
	us.user_id, us.name, us.email, us.role
	FROM item_comments AS ic
	INNER JOIN users AS us ON ic.commented_by = us.user_id
	WHERE ic.item_id = $1
	ORDER BY ic.created_at DESC, ic.title
	LIMIT ($2 + 1) OFFSET ($3 - 1) * $2`
	rows, err := cs.db.Query(context.Background(), sql, item.Id, limit, page)
	if err != nil {
		return []store.ItemComment{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer rows.Close()

	comments, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (store.ItemComment, error) {
		var comment store.ItemComment
		if err := row.Scan(&comment.Id, &comment.Title, &comment.Message, &comment.Rating, &comment.UpVotes, &comment.DownVotes,
			&comment.IsEdited, &comment.CreatedAt, &comment.EditedAt,
			&comment.CommentedBy.Id, &comment.CommentedBy.Name, &comment.CommentedBy.Email, &comment.CommentedBy.Role); err != nil {
			return store.ItemComment{}, err
		}
		comment.Item = item
		return comment, nil
	})
	if err != nil {
		return []store.ItemComment{}, echo.NewHTTPError(http.StatusInternalServerError)
	}

	return comments, nil
}

func (cs Comment) GetCommentsByRating(item store.Item, limit int, page int, rate int) ([]store.ItemComment, error) {
	sql := `SELECT ic.id, ic.title, ic.message, ic.rating, ic.up_votes, ic.down_votes,
	ic.is_edited, ic.created_at, ic.edited_at,
	us.user_id, us.name, us.email, us.role
	FROM item_comments AS ic
	INNER JOIN users AS us ON ic.commented_by = us.user_id
	WHERE ic.item_id = $1 AND ic.rating = $4
	ORDER BY ic.created_at DESC, ic.title
	LIMIT ($2 + 1) OFFSET ($3 - 1) * $2`
	rows, err := cs.db.Query(context.Background(), sql, item.Id, limit, page, rate)
	if err != nil {
		return []store.ItemComment{}, echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer rows.Close()

	comments, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (store.ItemComment, error) {
		var comment store.ItemComment
		if err := row.Scan(&comment.Id, &comment.Title, &comment.Message, &comment.Rating, &comment.UpVotes, &comment.DownVotes,
			&comment.IsEdited, &comment.CreatedAt, &comment.EditedAt,
			&comment.CommentedBy.Id, &comment.CommentedBy.Name, &comment.CommentedBy.Email, &comment.CommentedBy.Role); err != nil {
			return store.ItemComment{}, err
		}
		comment.Item = item
		return comment, nil
	})
	if err != nil {
		return []store.ItemComment{}, echo.NewHTTPError(http.StatusInternalServerError)
	}

	return comments, nil
}

func (cs Comment) InsertComment(comment store.ItemComment) (int, error) {
	var id int
	sql := `INSERT INTO item_comments (item_id, commented_by,
	title, message, rating)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id`
	if err := cs.db.QueryRow(context.Background(), sql, comment.Item.Id, comment.CommentedBy.Id,
		comment.Title, comment.Message, comment.Rating).Scan(&id); err != nil {
		return 0, echo.NewHTTPError(http.StatusInternalServerError)
	}

	return id, nil
}

func (cs Comment) EditComment(id int, comment store.ItemComment) error {
	sql := `UPDATE item_comments
	SET title = $1, message = $2, rating = $3, is_edited = TRUE, edited_at = NOW()
	WHERE id = $4`
	if _, err := cs.db.Exec(context.Background(), sql, comment.Title, comment.Message, comment.Rating, id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

func (cs Comment) DeleteComment(id int) error {
	sql := `DELETE FROM item_comments WHERE id = $1`
	if _, err := cs.db.Exec(context.Background(), sql, id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

func (cs Comment) UpvoteComment(id int) error {
	sql := `UPDATE item_comments SET up_votes = up_votes + 1 WHERE id = $1`
	if _, err := cs.db.Exec(context.Background(), sql, id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

func (cs Comment) DownvoteComment(id int) error {
	sql := `UPDATE item_comments SET down_votes = down_votes + 1 WHERE id = $1`
	if _, err := cs.db.Exec(context.Background(), sql, id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}
