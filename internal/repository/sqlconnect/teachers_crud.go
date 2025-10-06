package sqlconnect

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"schoolapi/internal/models"
	"schoolapi/pkg/utils"
	"strconv"
	"strings"
)

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
func AddTeacherDB(w http.ResponseWriter, newTeachers []models.Teacher) ([]models.Teacher, error) {
	db, err := ConnectDB("school")
	//check connecntion Error
	if err != nil {
		utils.CheckHttpError(err, w, "Error Connecting to the Database", http.StatusInternalServerError)
		return nil, err
	}
	defer db.Close()

	//Preparing Statement
	stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES(?,?,?,?,?)")
	//handling err
	utils.CheckHttpError(err, w, "Error in preparing sql query ", http.StatusInternalServerError)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	//making a new array of added teachers
	addedTeachers := make([]models.Teacher, len(newTeachers))
	//looping and adding teacher in a new teachers
	for i, teacher := range newTeachers {
		res, err := stmt.Exec(teacher.FirstName, teacher.LastName, teacher.Email, teacher.Class, teacher.Subject)
		utils.CheckHttpError(err, w, "Error adding teacher to db", http.StatusInternalServerError)
		if err != nil {
			return nil, err
		}

		lastID, err := res.LastInsertId()
		// utils.CheckHttpError(err, w, "errro geting last ID", http.StatusInternalServerError)

		utils.CheckHttpError(err, w, "Error Getting Last ID", http.StatusInternalServerError)
		if err != nil {
			return nil, err
		}
		// pushing the enered values to added teacher slice
		teacher.ID = int(lastID)
		addedTeachers[i] = teacher
	}
	return addedTeachers, err
}

func GetAllTeachersDB(w http.ResponseWriter, r *http.Request) (error, []models.Teacher) {
	db, err := ConnectDB("school")
	//check connecntion Error
	// utils.CheckHttpError(err, w, "Error Connecting to the Database", http.StatusInternalServerError)
	if err != nil {
		return utils.ErrorHandler(err, "Error Connecting to DB"), nil
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
		// utils.CheckHttpError(err, w, "Error Query Teachers", http.StatusInternalServerError)
		return utils.ErrorHandler(err, "Error Finding Table to DB"), nil
	}
	defer rows.Close()
	teachersList := make([]models.Teacher, 0)

	for rows.Next() {
		var teacher models.Teacher
		err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {
			// utils.CheckHttpError(err, w, "Error Scaning Row", http.StatusInternalServerError)
			return utils.ErrorHandler(err, "Error scanning to DB"), nil
		}
		teachersList = append(teachersList, teacher)
	}
	return utils.ErrorHandler(err, "Somthing Went Wrong"), teachersList
}
func GetTeacherByID(w http.ResponseWriter, numID int) (models.Teacher, error) {
	db, err := ConnectDB("school")
	//check connecntion Error
	if err != nil {
		// utils.CheckHttpError(err, w, "Error Connecting to the Database", http.StatusInternalServerError)
		return models.Teacher{}, utils.ErrorHandler(err, "Error Connecting to the Database")
	}
	defer db.Close()
	//geting id from Path

	//making new var to hold the data
	var teacher models.Teacher
	//query through the row
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", numID).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	//handling errors
	if err == sql.ErrNoRows {
		// utils.CheckHttpError(err, w, "Teacher not Found", http.StatusNotFound)
		return models.Teacher{}, utils.ErrorHandler(err, "Teacher not Found")
	} else if err != nil {
		// utils.CheckHttpError(err, w, "Database Query Error", http.StatusInternalServerError)
		return models.Teacher{}, utils.ErrorHandler(err, "Database Query Error")
	}
	return teacher, err
}

func PutOneTeacherDB(w http.ResponseWriter, id int, updatedTeacher models.Teacher) error {
	db, err := ConnectDB("school")
	if err != nil {
		// utils.CheckHttpError(err, w, "Data base Not Connecting", http.StatusInternalServerError)

		return utils.ErrorHandler(err, "Error ConnectingS to DB")
	}
	defer db.Close()
	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email , class, subject FROM teachers WHERE id = ?", id).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)

	if err != nil {
		if err == sql.ErrNoRows {
			// utils.CheckHttpError(err, w, "Teahcer Not Found", http.StatusNotFound)
			return utils.ErrorHandler(err, "Error Connecting to DB")
		}
		// utils.CheckHttpError(err, w, "Error Just Happend", http.StatusInternalServerError)
		return utils.ErrorHandler(err, "Internal Server Error")
	}
	updatedTeacher.ID = existingTeacher.ID
	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ? , email = ? , class = ? , subject= ? WHERE id = ?", updatedTeacher.FirstName, updatedTeacher.LastName, updatedTeacher.Email, updatedTeacher.Class, updatedTeacher.Subject, updatedTeacher.ID)
	if err != nil {
		// utils.CheckHttpError(err, w, "Error Updating Teacher", http.StatusInternalServerError)
		return utils.ErrorHandler(err, "Internal Server Error")
	}
	return err
}

func PatchTeacherDB(w http.ResponseWriter, updates []map[string]any) error {
	db, err := ConnectDB("school")
	if err != nil {
		// utils.CheckHttpError(err, w, "Error COnnecting to DB", http.StatusInternalServerError)
		return utils.ErrorHandler(err, "Internal Server Error")
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		// utils.CheckHttpError(err, w, "Error Starting Transaction", http.StatusInternalServerError)
		return utils.ErrorHandler(err, "Internal Server Error")
	}
	for _, update := range updates {
		idStr, ok := update["id"].(string)
		log.Println(update["id"])
		if !ok {
			tx.Rollback()
			// http.Error(w, "Error Parsing ID", http.StatusBadRequest)
			return utils.ErrorHandler(err, "Internal Server Error")
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			tx.Rollback()
			// utils.CheckHttpError(err, w, "Invalid ID ", http.StatusBadRequest)
			return utils.ErrorHandler(err, "Internal Server Error")
		}
		var teacherFromDB models.Teacher
		err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(&teacherFromDB.ID, &teacherFromDB.FirstName, &teacherFromDB.LastName, &teacherFromDB.Email, &teacherFromDB.Class, &teacherFromDB.Subject)
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				// utils.CheckHttpError(err, w, "Teacher not Found", http.StatusNotFound)
				return utils.ErrorHandler(err, "Internal Server Error")
			}
			// utils.CheckHttpError(err, w, "Error Fetching DB", http.StatusInternalServerError)
			return utils.ErrorHandler(err, "Internal Server Error")
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
							// utils.CheckHttpError(err, w, "Error Converting Value of Field", http.StatusBadRequest)
							return utils.ErrorHandler(err, "Internal Server Error")
						}
						break
					}
				}

			}
			_, err := tx.Exec("UPDATE teachers SET first_name=? , last_name=?, email=? ,subject=?, class=? WHERE id = ?", teacherFromDB.FirstName, teacherFromDB.LastName, teacherFromDB.Email, teacherFromDB.Class, teacherFromDB.Subject, teacherFromDB.ID)
			if err != nil {
				tx.Rollback()
				// utils.CheckHttpError(err, w, "Error Updating Row", http.StatusInternalServerError)
				return utils.ErrorHandler(err, "Internal Server Error")
			}

		}

	}
	err = tx.Commit()
	if err != nil {
		// utils.CheckHttpError(err, w, "Error Commiting Updates", http.StatusInternalServerError)
		return utils.ErrorHandler(err, "Internal Server Error")
	}
	return nil
}
func PatchOneTeacherDB(w http.ResponseWriter, id int, updates map[string]interface{}) (models.Teacher, error) {
	db, err := ConnectDB("school")
	if err != nil {
		// utils.CheckHttpError(err, w, "Error COnnecting databade", http.StatusInternalServerError)
		return models.Teacher{}, utils.ErrorHandler(err, "Internal Server Error")
	}
	defer db.Close()
	//Getting Existing Teacher To a new Var
	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT first_name , last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(&existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)
	if err != nil {
		if err == sql.ErrNoRows {
			// utils.CheckHttpError(err, w, "Teacher not Found", http.StatusNotFound)
			return models.Teacher{}, utils.ErrorHandler(err, "Internal Server Error")
		}
		// utils.CheckHttpError(err, w, "Cannot Retrieve Data", http.StatusInternalServerError)
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
		return models.Teacher{}, utils.ErrorHandler(err, "Internal Server Error")
	}

	//Advanced solutions (MORE DYNAMIC)

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
		// utils.CheckHttpError(err, w, "ERROR Updating Teahcer", http.StatusInternalServerError)
		return models.Teacher{}, utils.ErrorHandler(err, "Internal Server Error")
	}
	return existingTeacher, nil
}

func DeleteOneTeacherDB(w http.ResponseWriter, id int) error {
	db, err := ConnectDB("school")
	if err != nil {
		// utils.CheckHttpError(err, w, "Cant Connect DB", http.StatusInternalServerError)
		return utils.ErrorHandler(err, "Internal Server Error")
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)
	if err != nil {
		// utils.CheckHttpError(err, w, "Cant Connect DB", http.StatusInternalServerError)
		return utils.ErrorHandler(err, "Internal Server Error")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// utils.CheckHttpError(err, w, "Cant Connect DB", http.StatusInternalServerError)
		return utils.ErrorHandler(err, "Internal Server Error")
	}
	if rowsAffected == 0 {
		// http.Error(w, "Teacher Not Found", http.StatusNotFound)
		return utils.ErrorHandler(err, "Internal Server Error")
	}
	return utils.ErrorHandler(err, "Internal Server Error")
}
func DeleteTeachersDB(w http.ResponseWriter, ids []int) ([]int, error) {
	db, err := ConnectDB("school")
	if err != nil {
		// utils.CheckHttpError(err, w, "Cant Connect DB", http.StatusInternalServerError)
		return nil, utils.ErrorHandler(err, "Internal Server Error")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		// utils.CheckHttpError(err, w, "Error Starting Transaction", http.StatusInternalServerError)
		return nil, utils.ErrorHandler(err, "Internal Server Error")
	}

	stmt, err := tx.Prepare("DELETE FROM teachers WHERE id = ?")
	if err != nil {
		// utils.CheckHttpError(err, w, "Error Intializing stmt", http.StatusInternalServerError)
		return nil, utils.ErrorHandler(err, "Internal Server Error")
	}
	defer stmt.Close()
	var deletedIds []int
	for _, id := range ids {

		res, err := stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			// utils.CheckHttpError(err, w, "Row Not Found", http.StatusInternalServerError)
			return nil, utils.ErrorHandler(err, "Internal Server Error")
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			tx.Rollback()
			// utils.CheckHttpError(err, w, "Error Getting Result", http.StatusInternalServerError)
			return nil, utils.ErrorHandler(err, "Internal Server Error")
		}
		if rowsAffected > 0 {
			deletedIds = append(deletedIds, id)
		}
		if rowsAffected < 1 {
			tx.Rollback()
			// http.Error(w, "No Teachers to deletet", http.StatusNotFound)
			return nil, utils.ErrorHandler(err, "Internal Server Error")
		}

	}
	err = tx.Commit()
	if err != nil {

		// utils.CheckHttpError(err, w, "Error Getting Result", http.StatusInternalServerError)
		return nil, utils.ErrorHandler(err, "Internal Server Error")
	}
	if len(deletedIds) < 1 {
		// http.Error(w, "Id does not exist", http.StatusNotFound)
		return deletedIds, utils.ErrorHandler(err, "Internal Server Error")
	}
	return deletedIds, nil
}
