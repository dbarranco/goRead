CREATE TABLE `Usuario` (
	`idUsuario` int NOT NULL,
	`email` varchar(100) NOT NULL UNIQUE,
	`passwd` varchar(40) NOT NULL,
	PRIMARY KEY (`idUsuario`)
);

CREATE TABLE `Libros` (
	`idLibro` int NOT NULL,
	`Descripcion` varchar(140) NOT NULL,
	`Creador` varchar(100) NOT NULL,
	`Idioma` varchar(20) NOT NULL,
	`Ano` int NOT NULL,
	PRIMARY KEY (`idLibro`)
);

CREATE TABLE `userLibros` (
	`idUsuario` int NOT NULL,
	`idLibro` int NOT NULL
);

ALTER TABLE `userLibros` ADD CONSTRAINT `userLibros_fk0` FOREIGN KEY (`idUsuario`) REFERENCES `Usuario`(`idUsuario`);

ALTER TABLE `userLibros` ADD CONSTRAINT `userLibros_fk1` FOREIGN KEY (`idLibro`) REFERENCES `Libros`(`idLibro`);

