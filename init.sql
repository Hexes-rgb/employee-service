-- init.sql
CREATE TABLE IF NOT EXISTS departments (
    id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    UNIQUE(company_id, name)
);

CREATE TABLE IF NOT EXISTS employees (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    company_id INTEGER NOT NULL,
    department_id INTEGER REFERENCES departments(id) ON DELETE SET NULL,
    passport_type VARCHAR(20),
    passport_number VARCHAR(50) UNIQUE
);