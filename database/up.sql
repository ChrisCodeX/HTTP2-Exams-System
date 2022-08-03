DROP TABLE IF EXISTS students;

CREATE TABLE students (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    age INTEGER NOT NULL
);

DROP TABLE IF EXISTS exams;

CREATE TABLE exams (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

DROP TABLE IF EXISTS questions;

CREATE TABLE questions (
    id VARCHAR(32) PRIMARY KEY,
    exam_id VARCHAR(32) NOT NULL,
    question VARCHAR(255) NOT NULL,
    answer VARCHAR(255) NOT NULL,
    FOREIGN KEY (exam_id) REFERENCES exams(id)
);

DROP TABLE IF EXISTS enrollments;

CREATE TABLE enrollments (
    student_id VARCHAR(32) NOT NULL,
    exam_id VARCHAR(32) NOT NULL,
    FOREIGN KEY (student_id) REFERENCES students(id),
    FOREIGN KEY (exam_id) REFERENCES exams(id)
);