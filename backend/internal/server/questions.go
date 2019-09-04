package server

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"log"
)

// JSONData stores array of QuestionDTOs
type JSONData struct {
	Array []QuestionDTO `json:"questions"`
}

// QuestionDTO is used to store questions and answers from JSON
type QuestionDTO struct {
	Text    string   `db:"question" json:"text"`
	Answers []string `db:"answers" json:"answers"`
}

func (s *Server) questionsGETHandler(w http.ResponseWriter, r *http.Request) {
	data := JSONData{
		Array: []QuestionDTO{},
	}

	var err error

	data.Array, err = s.questionsService.GetAllQuestions()
	if err != nil {
		_ = NewAPIError(err.Error(), 500).Send(w)
		return
	}

	body, err := json.Marshal(data)
	if err != nil {
		_ = NewAPIError(err.Error(), http.StatusInternalServerError).Send(w)
		return
	}

	_, _ = w.Write(body)

}

func (s *Server) questionsPOSTHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		_ = NewAPIError("No body", http.StatusBadRequest).Send(w)
		return
	}

	defer func() { _ = r.Body.Close() }()

	data := JSONData{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		_ = NewAPIError(err.Error(), http.StatusInternalServerError).Send(w)
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		_ = NewAPIError(err.Error(), http.StatusBadRequest).Send(w)
		return
	}

	err = s.questionsService.DropTable()
	if err != nil {
		log.Println("Not dropping questions: ", err)
	}

	err = s.questionsService.CreateTable()
	if err != nil {
		_ = NewAPIError(err.Error(), http.StatusInternalServerError).Send(w)
		return
	}

	for iterator := range data.Array {
		err = s.questionsService.AddQuestion(data.Array[iterator].Text, data.Array[iterator].Answers)
		if err != nil {
			_ = NewAPIError(err.Error(), http.StatusInternalServerError).Send(w)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}