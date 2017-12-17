create schema versus;

CREATE TABLE `versus`.`db_operations` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `stringColumn` VARCHAR(500) NOT NULL,
  `intColumn` INT NOT NULL,
  PRIMARY KEY (`id`));

INSERT INTO `versus`.`db_operations` (`stringColumn`, `id`, `intColumn`) VALUES ('some string', '1', '1');
