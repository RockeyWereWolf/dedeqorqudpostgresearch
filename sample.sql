CREATE TABLE books (
  id SERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  main_character TEXT NOT NULL,
  content TEXT NOT NULL
);
INSERT INTO books (title, main_character, content)
VALUES ('Boghach Khan Son of Dirse Khan','Boghach Khan','One day Bayindir Khan son of Kam Ghan rose up from his
place. He had his striped parasol set up on the earth’s face, his
many-coloured pavilion reared up to the face of the sky. In a
thousand places silken rugs were spread. Once a year the Khan
of Khans, Bayindir Khan, used to make a feast and entertain
the Oghuz nobles. This year again he made a feast and had his
men slaughter of horses the stallions, of camels the males, of
sheep the rams'); 
