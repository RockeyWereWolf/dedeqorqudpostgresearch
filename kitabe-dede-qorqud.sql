CREATE TABLE books (
  id SERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  main_character TEXT NOT NULL,
  prologue TEXT NOT NULL
);

CREATE TABLE episodes (
  id SERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  book_id INTEGER NOT NULL,
  episode_text TEXT NOT NULL,
  FOREIGN KEY (book_id) REFERENCES books(id)
);

CREATE TABLE characters (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  book_id INTEGER NOT NULL,
  episode_id INTEGER NOT NULL,
  FOREIGN KEY (book_id) REFERENCES books(id),
  FOREIGN KEY (episode_id) REFERENCES episodes(id)
);
