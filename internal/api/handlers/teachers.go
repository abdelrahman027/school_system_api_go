package handlers

import (
	"encoding/json"
	"net/http"
	"schoolapi/internal/models"
	"schoolapi/internal/repository/sqlconnect"
	"strconv"
)

// Creating Teacher or teachers
func CreateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	// Store The New values in a Slice
	var newTeachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	//handling error
	if err != nil {
		// utils.CheckHttpError(err, w, "Error Decoding Data", http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//Connecting DB
	addedTeachers, err := sqlconnect.AddTeacherDB(w, newTeachers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//Creating the Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}
	//handling encode error
	err = json.NewEncoder(w).Encode(response)
	// utils.CheckHttpError(err, w, "Error decoding data", http.StatusInternalServerError)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

// geting Teachers  Many
func GetTeacherHandler(w http.ResponseWriter, r *http.Request) {
	//Connecting DB
	err, teachersList := sqlconnect.GetAllTeachersDB(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(teachersList),
		Data:   teachersList,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		// utils.CheckHttpError(err, w, "error Encoding Response", http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

}

// geting Teacher Single
func GetOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
	IDstr := r.PathValue("id")
	// if there is no id

	//with id provided
	numID, err := strconv.Atoi(IDstr)
	// utils.CheckHttpError(err, w, "Error Getting ID", http.StatusBadRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//Connecting DB
	teacher, err := sqlconnect.GetTeacherByID(w, numID)
	if err != nil {
		// log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	//set application header
	w.Header().Set("Content-Type", "application/json")
	//encoding json
	err = json.NewEncoder(w).Encode(teacher)
	//checking error
	if err != nil {
		// utils.CheckHttpError(err, w, "Error encoding response", http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
}

// Put Teacher Handler
func UpdateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	//geting id from Path
	IDstr := r.PathValue("id")
	id, err := strconv.Atoi(IDstr)
	if err != nil {
		// utils.CheckHttpError(err, w, "Invalid ID", http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var updatedTeacher models.Teacher
	err = json.NewDecoder(r.Body).Decode(&updatedTeacher)
	if err != nil {
		// utils.CheckHttpError(err, w, "Error Decoding From Json", http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusBadRequest)

	}

	err = sqlconnect.PutOneTeacherDB(w, id, updatedTeacher)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(updatedTeacher)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// Patch Teacher Handler
func PatchTeacherHandler(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]any
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		// utils.CheckHttpError(err, w, "Invalid Payload", http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = sqlconnect.PatchTeacherDB(w, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Patch One Teacher Handler
func PatchOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
	strID := r.PathValue("id")
	//Convert it to INT
	id, err := strconv.Atoi(strID)
	if err != nil {
		// utils.CheckHttpError(err, w, "Invalid Teacher ID", http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//Make map of updated FIelds and Decode it to json
	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		// utils.CheckHttpError(err, w, "Invalid Payload", http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
	//connect to DB
	existingTeacher, err := sqlconnect.PatchOneTeacherDB(w, id, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(existingTeacher)
	if err != nil {
		// utils.CheckHttpError(err, w, "Invalid Request", http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

// Delete Teacher Handler
func DeleteOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
	strID := r.PathValue("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		// utils.CheckHttpError(err, w, "Invalid Id", http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	err = sqlconnect.DeleteOneTeacherDB(w, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	// w.WriteHeader(http.StatusNoContent)

	w.Header().Set("Content-Type", "application/json")

	response := struct {
		Status string `json:"status"`
		ID     int    `json:"id"`
	}{
		Status: "Teacher Successfully Deleted",
		ID:     id,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

}

// Delete Teach Handler MANY
func DeleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
	var ids []int
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		// utils.CheckHttpError(err, w, "Invalid Payload", http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	deletedIds, err := sqlconnect.DeleteTeachersDB(w, ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status       string `json:"status"`
		DeletedCount int    `json:"number_deleted"`
		DeletedIds   []int  `json:"deleted_ids"`
	}{
		Status:       "Success",
		DeletedCount: len(deletedIds),
		DeletedIds:   deletedIds,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

}
