CREATE TABLE users (
          id        SERIAL PRIMARY KEY,
          username   VARCHAR(255) NOT NULL UNIQUE,
          password   VARCHAR(255) NOT NULL,
          created_at TIMESTAMP    NOT NULL DEFAULT now()
);

CREATE TABLE notes (
          id SERIAL PRIMARY KEY,
          user_id INT NOT NULL,
          title VARCHAR(255) NOT NULL,
          content TEXT NOT NULL,
          created_at TIMESTAMP DEFAULT NOW(),
          FOREIGN KEY (user_id) REFERENCES users(id)
);


