package information

import (
	"database/sql"

	"github.com/Megidy/e-commerce/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}
func (s *Store) CreateQuestion(question types.Question) error {
	_, err := s.db.Exec("insert into questions values(?,?,?,?)", question.Id, question.UserID, question.Title, question.Body)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) CreateRespond(respond types.Respond) error {
	_, err := s.db.Exec("insert into responds values(?,?,?)", respond.QuestionID, respond.RespondID, respond.Body)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) GetAllQuestions() ([]types.Question, error) {
	var questions []types.Question
	rows, err := s.db.Query("select * from questions")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var q types.Question
		err := rows.Scan(&q.Id, &q.UserID, &q.Title, &q.Body)
		if err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}
	return questions, err
}

func (s *Store) GetSingleQuestion(questionID string) (types.Question, error) {
	var q types.Question
	row, err := s.db.Query("select * from questions where id=?", questionID)
	if err != nil {
		return types.Question{}, err
	}
	for row.Next() {
		err = row.Scan(&q.Id, &q.UserID, &q.Title, &q.Body)
		if err != nil {
			return types.Question{}, err
		}
	}
	return q, nil
}
func (s *Store) GetAllResponds(questionID string) ([]types.Respond, error) {
	var responds []types.Respond
	rows, err := s.db.Query("select * from responds where question_id=?", questionID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var r types.Respond
		err = rows.Scan(&r.QuestionID, &r.RespondID, &r.Body)
		if err != nil {
			return nil, err
		}
		responds = append(responds, r)
	}
	return responds, err
}
func (s *Store) DeleteQuestion(questionID string) error {
	_, err := s.db.Exec("delete from questions where id=?", questionID)
	if err != nil {
		return err
	}
	return nil
}
