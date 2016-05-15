CREATE TABLE `Usuario` (
	`idUsuario` int NOT NULL,
	`email` varchar NOT NULL UNIQUE,
	`passwd` varchar(40) NOT NULL,
	PRIMARY KEY (`idUsuario`)
);

CREATE TABLE `Libros` (
	`idLibro` int NOT NULL,
	`DC.Description` char NOT NULL,
	`DC.Creator` char NOT NULL,
	`DC.Language` char NOT NULL,
	PRIMARY KEY (`idLibro`)
);

CREATE TABLE `userLibros` (
	`idUsuario` int NOT NULL,
	`idLibro` int NOT NULL
);

ALTER TABLE `userLibros` ADD CONSTRAINT `userLibros_fk0` FOREIGN KEY (`idUsuario`) REFERENCES `Usuario`(`idUsuario`);

ALTER TABLE `userLibros` ADD CONSTRAINT `userLibros_fk1` FOREIGN KEY (`idLibro`) REFERENCES `Libros`(`idLibro`);

