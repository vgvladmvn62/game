package questions

import (
	"github.com/lib/pq"
	"github.wdf.sap.corp/Magikarpie/bullseye/internal/db"
	"github.wdf.sap.corp/Magikarpie/bullseye/internal/server"
)

// Repository holds a client of database.
type Repository struct {
	database *db.Database
}

// NewRepository creates new Repository.
func NewRepository(db *db.Database) *Repository {
	return &Repository{db}
}

// CreateTable creates table in database.
func (repository *Repository) CreateTable() error {
	_, err := repository.database.Exec("CREATE TABLE questions (question TEXT, answers TEXT[])")

	return err
}

// DropTable drops table.
func (repository *Repository) DropTable() error {
	_, err := repository.database.Exec("DROP TABLE questions")

	return err
}

// AddQuestion adds question with answers to database.
func (repository *Repository) AddQuestion(text string, answers []string) error {
	stmt, err := repository.database.Prepare("INSERT INTO questions(question, answers) VALUES ($1, $2);")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(text, pq.Array(answers))
	if err != nil {
		return err
	}

	return nil
}

// GetAllQuestions returns all questions from database.
func (repository *Repository) GetAllQuestions() ([]server.QuestionDTO, error) {
	var questions []server.QuestionDTO
	rows, err := repository.database.Query("SELECT question, answers FROM questions")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var question string
	var answers pq.StringArray
	for rows.Next() {
		err = rows.Scan(&question, &answers)

		if err != nil {
			return nil, err
		}

		questions = append(questions, server.QuestionDTO{Text: question, Answers: answers})

	}
	return questions, nil
}
