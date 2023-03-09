CREATE TABLE books (
  id SERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  main_character TEXT NOT NULL,
  prologue TEXT NOT NULL
);
/*
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
*/
INSERT INTO books (title, main_character, prolugue)
VALUES ('Boghach Khan Son of Dirse Khan','Boghach Khan','One day Bayindir Khan son of Kam Ghan rose up from his
place. He had his striped parasol set up on the earth’s face, his
many-coloured pavilion reared up to the face of the sky. In a
thousand places silken rugs were spread. Once a year the Khan
of Khans, Bayindir Khan, used to make a feast and entertain
the Oghuz nobles. This year again he made a feast and had his
men slaughter of horses the stallions, of camels the males, of
sheep the rams');
