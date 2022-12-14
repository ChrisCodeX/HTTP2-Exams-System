package database

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/ChrisCodeX/gRPC/models"
	_ "github.com/lib/pq"
)

// PostgresRepository
type PostgresRepository struct {
	db *sql.DB
}

// Constructor of Postgres Repository
func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{db}, nil
}

// Methods

/*Student Service*/
// Get student from database by id
func (repo *PostgresRepository) GetStudent(ctx context.Context, id string) (*models.Student, error) {
	// Query
	rows, err := repo.db.QueryContext(ctx, "SELECT id, name, age FROM students WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	// Stop reading rows
	defer CloseReadingRows(rows)

	// Map rows values of the query into the students struct
	var student = models.Student{}

	for rows.Next() {
		if err = rows.Scan(&student.Id, &student.Name, &student.Age); err != nil {
			return &student, nil
		}
	}
	if err = rows.Err(); err != nil {
		return &student, err
	}

	return &student, nil
}

// Insert a student into the database
func (repo *PostgresRepository) SetStudent(ctx context.Context, student *models.Student) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO students (id, name, age) VALUES ($1, $2, $3)",
		student.Id, student.Name, student.Age)
	return err
}

/*Exams Service*/
// Get exam by id
func (repo *PostgresRepository) GetExam(ctx context.Context, id string) (*models.Exam, error) {
	// Query
	rows, err := repo.db.QueryContext(ctx, "SELECT id, name FROM exams WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	// Stop reading rows
	defer CloseReadingRows(rows)

	// Map rows values of the query into the students struct
	var exam = models.Exam{}

	for rows.Next() {
		if err = rows.Scan(&exam.Id, &exam.Name); err != nil {
			return &exam, nil
		}
	}
	if err = rows.Err(); err != nil {
		return &exam, err
	}

	return &exam, nil
}

// Insert a exam into the database
func (repo *PostgresRepository) SetExam(ctx context.Context, exam *models.Exam) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO exams (id, name) VALUES ($1, $2)",
		exam.Id, exam.Name)
	return err
}

// Enroll a student
func (repo *PostgresRepository) SetEnrollment(ctx context.Context, enrollment *models.Enrollment) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO enrollments (id, fk_exam_id, fk_student_id) VALUES ($1, $2, $3)", enrollment.Id, enrollment.ExamId, enrollment.StudentId)
	return err
}

// Get students by exam id
func (repo *PostgresRepository) GetStudentsPerExam(ctx context.Context, examId string) ([]*models.Student, error) {
	// Query
	rows, err := repo.db.QueryContext(ctx, "SELECT id, name, age FROM students 	WHERE id IN (SELECT fk_student_id FROM enrollments WHERE fk_exam_id = $1)", examId)
	if err != nil {
		return nil, err
	}

	// Stop reading rows
	defer CloseReadingRows(rows)

	// Map query results into student struct and add into students struct
	var students []*models.Student

	for rows.Next() {
		var student = models.Student{}
		if err = rows.Scan(&student.Id, &student.Name, &student.Age); err == nil {
			students = append(students, &student)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

// Insert a question into the database
func (repo *PostgresRepository) SetQuestion(ctx context.Context, question *models.Question) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO questions (id, question, answer, fk_exam_id) VALUES ($1, $2, $3, $4)",
		question.Id, question.Question, question.Answer, question.ExamId)
	return err
}

// Get Question Per Exam
func (repo *PostgresRepository) GetQuestionPerExam(ctx context.Context, examId string) ([]*models.Question, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, question, answer FROM questions WHERE fk_exam_id = $1", examId)
	if err != nil {
		return nil, err
	}

	// Stop reading rows
	defer CloseReadingRows(rows)

	// Map rows of query into question struct
	var questions []*models.Question
	for rows.Next() {
		var question = models.Question{}
		if err = rows.Scan(&question.Id, &question.Question, &question.Answer); err == nil {
			questions = append(questions, &question)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return questions, nil
}

// Insert Answer into database
func (repo *PostgresRepository) GetCountQuestionsByExamId(ctx context.Context, examId string) (*uint16, error) {
	// Query
	row, err := repo.db.QueryContext(ctx, "SELECT COUNT(*) FROM questions")
	if err != nil {
		return nil, err
	}
	defer CloseReadingRows(row)
	var countString string
	for row.Next() {
		if err = row.Scan(&countString); err != nil {
			return nil, err
		}
	}
	if err = row.Err(); err != nil {
		return nil, err
	}
	countInt, err := strconv.Atoi(countString)
	if err != nil {
		return nil, err
	}
	count := uint16(countInt)
	return &count, nil
}

// Get enrollment by id
func (repo *PostgresRepository) GetEnrollmentById(ctx context.Context, id string) (*models.Enrollment, error) {
	// Query
	rows, err := repo.db.QueryContext(ctx, "SELECT id, fk_exam_id, fk_student_id FROM enrollments WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	// Stop reading rows
	defer CloseReadingRows(rows)

	// Map rows values of the query into the students struct
	var enrollment = models.Enrollment{}

	for rows.Next() {
		if err = rows.Scan(&enrollment.Id, &enrollment.ExamId, &enrollment.StudentId); err != nil {
			return &enrollment, nil
		}
	}
	if err = rows.Err(); err != nil {
		return &enrollment, err
	}

	return &enrollment, nil
}

// Set Qualifications
func (repo *PostgresRepository) SetQualifications(ctx context.Context, qualification *models.Qualification) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO qualifications (id, score, fk_enrollment_id) VALUES ($1, $2, $3)",
		qualification.Id, qualification.Score, qualification.EnrollmentId)
	return err
}

// Get Qualifications
func (repo *PostgresRepository) GetQualificationsByEnrollmentId(ctx context.Context, enrollmentId string) (*models.Qualification, error) {
	// Query
	rows, err := repo.db.QueryContext(ctx, "SELECT id, fk_enrollment_id, score FROM qualifications WHERE fk_enrollment_id = $1", enrollmentId)
	if err != nil {
		return nil, err
	}

	// Stop reading rows
	defer CloseReadingRows(rows)

	// Map rows values of the query into the students struct
	var qualification = models.Qualification{}

	for rows.Next() {
		if err = rows.Scan(&qualification.Id, &qualification.EnrollmentId, &qualification.Score); err != nil {
			return &qualification, nil
		}
	}
	if err = rows.Err(); err != nil {
		return &qualification, err
	}

	return &qualification, nil
}

// Set Student Answers
func (repo *PostgresRepository) SetStudentAnswers(ctx context.Context, studentAnswers *models.StudentAnswers) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO student_answers (id, fk_enrollment_id, question_id, student_answer, correct) VALUES ($1, $2, $3, $4, $5)",
		studentAnswers.Id, studentAnswers.EnrollmentId, studentAnswers.QuestionId, studentAnswers.StudentAnswer, studentAnswers.Correct)
	return err
}

// Get answers of students per enrollments
func (repo *PostgresRepository) GetAnswersPerEnrollment(ctx context.Context, enrollmentId string) ([]*models.StudentAnswers, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, fk_enrollment_id, question_id, student_answer, correct FROM student_answers WHERE fk_enrollment_id = $1", enrollmentId)
	if err != nil {
		return nil, err
	}

	// Stop reading rows
	defer CloseReadingRows(rows)

	// Map rows of query into question struct
	var studentAnswers []*models.StudentAnswers
	for rows.Next() {
		var studentAnswer = models.StudentAnswers{}
		if err = rows.Scan(&studentAnswer.Id, &studentAnswer.EnrollmentId, &studentAnswer.QuestionId, &studentAnswer.StudentAnswer, &studentAnswer.Correct); err == nil {
			studentAnswers = append(studentAnswers, &studentAnswer)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return studentAnswers, nil
}
