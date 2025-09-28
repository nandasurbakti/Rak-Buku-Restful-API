-- +migrate Up
-- +migrate StatementBegin

create table books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    image_url VARCHAR(255),
    release_year INT NOT NULL,
    price INT NOT NULL,
    total_page INT NOT NULL CHECK (total_page > 0),
    thickness VARCHAR(10) NOT NULL CHECK (thickness IN ('tipis', 'tebal')),
    category_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255) NOT NULL,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_by VARCHAR(255) NOT NULL, 

    CONSTRAINT fk_books_category
    FOREIGN KEY (category_id)
    REFERENCES categories(id)
    ON DELETE SET NULL 
    ON UPDATE CASCADE
)

-- +migrate StatementEnd