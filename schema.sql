-- DROP TABLE IF EXISTS COMPTES;
-- DROP TABLE IF EXISTS AGENDA;
-- DROP TABLE IF EXISTS MESSAGES;

CREATE TABLE COMPTES (
  id TEXT NOT NULL,
  mdp  TEXT NOT NULL,
  description  TEXT NOT NULL,
  PRIMARY KEY(id));
CREATE TABLE AGENDA (
  numero INTEGER PRIMARY KEY AUTOINCREMENT,
  id  TEXT NOT NULL,
  date  TEXT NOT NULL,
  contenu TEXT NOT NULL,
  FOREIGN KEY(id) REFERENCES comptes(id));
CREATE TABLE MESSAGES (
  numero INTEGER PRIMARY KEY AUTOINCREMENT,
  id TEXT NOT NULL,
  dest TEXT NOT NULL,
  date TEXT NOT NULL,  
  contenu TEXT NOT NULL,
  FOREIGN KEY(id) REFERENCES comptes(id),
  FOREIGN KEY(dest) REFERENCES comptes(id));
CREATE TABLE ACTIVITES (
  nom TEXT NOT NULL,
  physique INTEGER NOT NULL,
  intellectuel INTEGER NOT NULL,
  social INTEGER NOT NULL,
  description TEXT NOT NULL,
  PRIMARY KEY (nom));

INSERT INTO COMPTES VALUES ('bot', '123', 'le compte du chatbot');
INSERT INTO ACTIVITES VALUES ('Tai Chi', 1, 0, 0, "Le Tai Chi est un art martial chinois traditionnel pratiqué pour ses bienfaits sur la santé physique. Il consiste en une série de mouvements lents et fluides associée à un exercice de respiration. Pratiquer cette activité permettra d'améliorer votre flexibilité et votre équilibre tout en réduisant votre stress.");
INSERT INTO ACTIVITES VALUES ('Marche nordique', 1, 0, 1, "C'est une marche sportive dans laquelle vous disposez de bâtons spécifiques semblables à ceux du ski de fond qui vous permettent de vous projetez plus rapidement vers l'avant, sollicitant ainsi l'ensemble de vos muscles. Cela vous permettra ainsi de travailler votre système cardiovasculaire mais également de vous aérer l'esprit.");
INSERT INTO ACTIVITES VALUES ('Jardinage', 1, 0, 1, "Cette activité consiste en la culture des végétaux (plantes, fleurs, fruits, etc...) dans un espace vert. Elle vous procurera naturellement une sensation de bien-être et vous permettra de vous connecter avec la nature et d'améliorer l'environnement.");
INSERT INTO ACTIVITES VALUES ('Théâtre', 1, 1, 1, "Le théâtre est une pratique artistique dans laquelle vous et plusieurs autres personnes interprétez des rôles différents pour raconter une histoire devant un public. Cette activité vous permettra de maintenir vos liens sociaux et/ou de faire de nouvelles rencontres en plus d'être un excellent exercice de mémoire. ");
INSERT INTO ACTIVITES VALUES ('Chant', 0, 0, 1, "C'est une autre pratique artistique dans laquelle vous chantez et produisez des mélodies. Vous pouvez la pratiquer en solo ou en groupe (une chorale ou autre) et avec ou sans accompagnement musical. Pratiquer cette activité renforcera votre respiration et votre posture en plus de travailler votre mémoire. Cela favorise de plus la socialisation lorsque pratiquée en groupe.");
