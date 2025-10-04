package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"schoolapi/internal/models"
	"schoolapi/internal/repository/sqlconnect"
	"schoolapi/pkg/utils"
	"strconv"
	"strings"
)

// Creating Teacher or teachers
func CreateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	//Connecting DB
	db, err := sqlconnect.ConnectDB("school")
	//check connecntion Error
	utils.CheckHttpError(err, w, "Error Connecting to the Database", http.StatusInternalServerError)
	if err != nil {
		return
	}
	defer db.Close()
	// Store The New values in a Slice
	var newTeachers []models.Teacher
	err = json.NewDecoder(r.Body).Decode(&newTeachers)
	//handling error
	utils.CheckHttpError(err, w, "Error Decoding Data", http.StatusBadRequest)
	if err != nil {
		return
	}

	//Preparing Statement
	stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES(?,?,?,?,?)")
	//handling err
	utils.CheckHttpError(err, w, "Error in preparing sql query ", http.StatusInternalServerError)
	if err != nil {
		return
	}
	defer stmt.Close()
	//making a new array of added teachers
	addedTeachers := make([]models.Teacher, len(newTeachers))
	//looping and adding teacher in a new teachers
	for i, teacher := range newTeachers {
		res, err := stmt.Exec(teacher.FirstName, teacher.LastName, teacher.Email, teacher.Class, teacher.Subject)
		utils.CheckHttpError(err, w, "Error adding teacher to db", http.StatusInternalServerError)
		if err != nil {
			return
		}

		lastID, err := res.LastInsertId()
		// utils.CheckHttpError(err, w, "errro geting last ID", http.StatusInternalServerError)

		utils.CheckHttpError(err, w, "Error Getting Last ID", http.StatusInternalServerError)
		if err != nil {
			return
		}
		// pushing the enered values to added teacher slice
		teacher.ID = int(lastID)
		addedTeachers[i] = teacher
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
	utils.CheckHttpError(err, w, "Error decoding data", http.StatusInternalServerError)
	if err != nil {
		return
	}

}

// geting Teachers  Many
func GetTeacherHandler(w http.ResponseWriter, r *http.Request) {
	//Connecting DB
	db, err := sqlconnect.ConnectDB("school")
	//check connecntion Error
	utils.CheckHttpError(err, w, "Error Connecting to the Database", http.StatusInternalServerError)
	if err != nil {
		return
	}
	defer db.Close()

	//DB Query
	query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1"
	//Slice of args to add in query
	var args []interface{}
	//
	query, args = FilterParams(r, query, args)
	query = addSorting(r, query)
	rows, err := db.Query(query, args...)
	if err != nil {
		utils.CheckHttpError(err, w, "Error Query Teachers", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	teachersList := make([]models.Teacher, 0)

	for rows.Next() {
		var teacher models.Teacher
		err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {
			utils.CheckHttpError(err, w, "Error Scaning Row", http.StatusInternalServerError)
			return
		}
		teachersList = append(teachersList, teacher)
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
		utils.CheckHttpError(err, w, "error Encoding Response", http.StatusInternalServerError)
		return
	}

}

// geting Teacher Single
func GetOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
	//Connecting DB
	db, err := sqlconnect.ConnectDB("school")
	//check connecntion Error
	utils.CheckHttpError(err, w, "Error Connecting to the Database", http.StatusInternalServerError)
	if err != nil {
		return
	}
	defer db.Close()
	//geting id from Path

	IDstr := r.PathValue("id")
	// if there is no id

	//with id provided
	numID, err := strconv.Atoi(IDstr)
	utils.CheckHttpError(err, w, "Error Getting ID", http.StatusBadRequest)
	if err != nil {
		return
	}
	//making new var to hold the data
	var teacher models.Teacher
	//query through the row
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", numID).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	//handling errors
	if err == sql.ErrNoRows {
		utils.CheckHttpError(err, w, "Teacher not Found", http.StatusNotFound)
		return
	} else if err != nil {
		utils.CheckHttpError(err, w, "Database Query Error", http.StatusInternalServerError)
		return
	}
	//set application header
	w.Header().Set("Content-Type", "application/json")
	//encoding json
	err = json.NewEncoder(w).Encode(teacher)
	//checking error
	if err != nil {
		utils.CheckHttpError(err, w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// Put Teacher Handler
func UpdateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	//geting id from Path
	IDstr := r.PathValue("id")
	id, err := strconv.Atoi(IDstr)
	if err != nil {
		utils.CheckHttpError(err, w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var updatedTeacher models.Teacher
	err = json.NewDecoder(r.Body).Decode(&updatedTeacher)
	if err != nil {
		utils.CheckHttpError(err, w, "Error Decoding From Json", http.StatusInternalServerError)
	}

	db, err := sqlconnect.ConnectDB("school")
	if err != nil {
		utils.CheckHttpError(err, w, "Data base Not Connecting", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email , class, subject FROM teachers WHERE id = ?", id).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.CheckHttpError(err, w, "Teahcer Not Found", http.StatusNotFound)
			return
		}
		utils.CheckHttpError(err, w, "Error Just Happend", http.StatusInternalServerError)
		return
	}
	updatedTeacher.ID = existingTeacher.ID
	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ? , email = ? , class = ? , subject= ? WHERE id = ?", updatedTeacher.FirstName, updatedTeacher.LastName, updatedTeacher.Email, updatedTeacher.Class, updatedTeacher.Subject, updatedTeacher.ID)
	if err != nil {
		utils.CheckHttpError(err, w, "Error Updating Teacher", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeacher)
}

// Patch Teacher Handler
func PatchTeacherHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sqlconnect.ConnectDB("school")
	if err != nil {
		utils.CheckHttpError(err, w, "Error COnnecting to DB", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	var updates []map[string]any
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		utils.CheckHttpError(err, w, "Invalid Payload", http.StatusBadRequest)
		return
	}
	tx, err := db.Begin()
	if err != nil {
		utils.CheckHttpError(err, w, "Error Starting Transaction", http.StatusInternalServerError)
		return
	}
	for _, update := range updates {
		idStr, ok := update["id"].(string)
		log.Println(update["id"])
		if !ok {
			tx.Rollback()
			http.Error(w, "Error Parsing ID", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			tx.Rollback()
			utils.CheckHttpError(err, w, "Invalid ID ", http.StatusBadRequest)
			return
		}
		var teacherFromDB models.Teacher
		err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(&teacherFromDB.ID, &teacherFromDB.FirstName, &teacherFromDB.LastName, &teacherFromDB.Email, &teacherFromDB.Class, &teacherFromDB.Subject)
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				utils.CheckHttpError(err, w, "Teacher not Found", http.StatusNotFound)
				return
			}
			utils.CheckHttpError(err, w, "Error Fetching DB", http.StatusInternalServerError)
			return
		}
		teacherVal := reflect.ValueOf(&teacherFromDB).Elem()
		teacherType := teacherVal.Type()
		teacherFieldCount := teacherVal.NumField()
		for k, v := range update {
			if k == "id" {
				continue
			}
			for i := range teacherFieldCount {
				field := teacherType.Field(i)
				if field.Tag.Get("json") == k+",omitempty" {
					fieldVal := teacherVal.Field(i)
					if fieldVal.CanSet() {
						val := reflect.ValueOf(v)
						if val.Type().ConvertibleTo(fieldVal.Type()) {

							fieldVal.Set(val.Convert(fieldVal.Type()))
						} else {
							tx.Rollback()
							utils.CheckHttpError(err, w, "Error Converting Value of Field", http.StatusBadRequest)
							return
						}
						break
					}
				}

			}
			_, err := tx.Exec("UPDATE teachers SET first_name=? , last_name=?, email=? ,subject=?, class=? WHERE id = ?", teacherFromDB.FirstName, teacherFromDB.LastName, teacherFromDB.Email, teacherFromDB.Class, teacherFromDB.Subject, teacherFromDB.ID)
			if err != nil {
				tx.Rollback()
				utils.CheckHttpError(err, w, "Error Updating Row", http.StatusInternalServerError)
				return
			}

		}

	}
	err = tx.Commit()
	if err != nil {
		utils.CheckHttpError(err, w, "Error Commiting Updates", http.StatusInternalServerError)
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
		utils.CheckHttpError(err, w, "Invalid Teacher ID", http.StatusBadRequest)
		return
	}
	//Make map of updated FIelds and Decode it to json
	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		utils.CheckHttpError(err, w, "Invalid Payload", http.StatusInternalServerError)
		return
	}
	//connect to DB
	db, err := sqlconnect.ConnectDB("school")
	if err != nil {
		utils.CheckHttpError(err, w, "Error COnnecting databade", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	//Getting Existing Teacher To a new Var
	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT first_name , last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(&existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.CheckHttpError(err, w, "Teacher not Found", http.StatusNotFound)
			return
		}
		utils.CheckHttpError(err, w, "Cannot Retrieve Data", http.StatusInternalServerError)
		return
	}
	//Option Long one (BASIC)
	// for key, value := range updates {
	// 	switch key {
	// 	case "first_name":
	// 		existingTeacher.FirstName = value.(string)
	// 	case "last_name":
	// 		existingTeacher.LastName = value.(string)
	// 	case "email":
	// 		existingTeacher.Email = value.(string)
	// 	case "class":
	// 		existingTeacher.Class = value.(string)
	// 	case "subject":
	// 		existingTeacher.Subject = value.(string)
	// 	}
	// }

	//Advanced solutions (MORE DYNAMIC)
	fmt.Println("Old", existingTeacher)
	teacherVal := reflect.ValueOf(&existingTeacher).Elem()
	teacherType := teacherVal.Type()
	teacherTypeCount := teacherVal.NumField()

	for k, v := range updates {
		for i := range teacherTypeCount {
			field := teacherType.Field(i)
			if field.Tag.Get("json") == k+",omitempty" {
				if teacherVal.Field(i).CanSet() {
					teacherVal.Field(i).Set(reflect.ValueOf(v).Convert(teacherVal.Field(i).Type()))
				}

			}
		}
	}
	fmt.Println("Updated", existingTeacher)
	_, err = db.Exec("UPDATE teachers SET first_name= ?,last_name= ?,email= ?,class= ?,subject= ? WHERE id = ?", &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject, &existingTeacher.ID)
	if err != nil {
		utils.CheckHttpError(err, w, "ERROR Updating Teahcer", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(existingTeacher)
	if err != nil {
		utils.CheckHttpError(err, w, "Invalid Request", http.StatusInternalServerError)
		return
	}

}

// Delete Teacher Handler
func DeleteOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
	strID := r.PathValue("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		utils.CheckHttpError(err, w, "Invalid Id", http.StatusBadRequest)
		return
	}

	db, err := sqlconnect.ConnectDB("school")
	if err != nil {
		utils.CheckHttpError(err, w, "Cant Connect DB", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)
	if err != nil {
		utils.CheckHttpError(err, w, "Cant Connect DB", http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.CheckHttpError(err, w, "Cant Connect DB", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Teacher Not Found", http.StatusNotFound)
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
		utils.CheckHttpError(err, w, "ErrorReturning Value", http.StatusBadRequest)
		return
	}

}

// Delete Teach Handler MANY
func DeleteTeacherHandler(w http.ResponseWriter, r *http.Request) {

	db, err := sqlconnect.ConnectDB("school")
	if err != nil {
		utils.CheckHttpError(err, w, "Cant Connect DB", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var ids []int
	err = json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		utils.CheckHttpError(err, w, "Invalid Payload", http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		utils.CheckHttpError(err, w, "Error Starting Transaction", http.StatusInternalServerError)
		return
	}

	stmt, err := tx.Prepare("DELETE FROM teachers WHERE id = ?")
	if err != nil {
		utils.CheckHttpError(err, w, "Error Intializing stmt", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	var deletedIds []int
	for _, id := range ids {

		res, err := stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			utils.CheckHttpError(err, w, "Row Not Found", http.StatusInternalServerError)
			return
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			tx.Rollback()
			utils.CheckHttpError(err, w, "Error Getting Result", http.StatusInternalServerError)
			return
		}
		if rowsAffected > 0 {
			deletedIds = append(deletedIds, id)
		}
		if rowsAffected < 1 {
			tx.Rollback()
			http.Error(w, "No Teachers to deletet", http.StatusNotFound)
			return
		}

	}
	err = tx.Commit()
	if err != nil {

		utils.CheckHttpError(err, w, "Error Getting Result", http.StatusInternalServerError)
		return
	}
	if len(deletedIds) < 1 {
		http.Error(w, "Id does not exist", http.StatusNotFound)
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
		utils.CheckHttpError(err, w, "Error Getting Result", http.StatusInternalServerError)
		return
	}

}

// validation order in sorting (GET METHOD)
func isValidSortOrder(order string) bool {
	return order == "asc" || order == "desc"
}

// validation Field in sorting (GET METHOD)
func isValidSortField(field string) bool {
	validFields := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
		"subject":    true,
	}
	return validFields[field]
}

func addSorting(r *http.Request, query string) string {
	sortParams := r.URL.Query()["sortBy"]
	if len(sortParams) > 0 {
		query += " ORDER BY"
	}
	for i, param := range sortParams {
		parts := strings.Split(param, ":")
		if len(parts) != 2 {
			continue
		}
		sortField, sortOrder := parts[0], parts[1]
		if !isValidSortField(sortField) && !isValidSortOrder(sortOrder) {
			continue
		}
		if i > 0 {
			query += " , "
		}
		query += " " + sortField + " " + sortOrder
	}
	return query
}

func FilterParams(r *http.Request, query string, args []interface{}) (string, []interface{}) {
	params := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
		"subject":    "subject",
	}

	for param, db_field := range params {
		value := r.URL.Query().Get(param)
		if value != "" {
			query += " AND " + db_field + " =?"
			args = append(args, value)
		}
	}
	return query, args
}
