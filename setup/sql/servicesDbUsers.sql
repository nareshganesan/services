-- MySQL Server Version 5.7.15

-- Set Password validation policy: Medium or STRONG
-- Medium: Length >= 8, numeric, mixed case, and special characters
-- Strong: Length >= 8, numeric, mixed case, special characters and dictionary
-- Default: Medium

-- -----------------------------------------------------
-- Create services admin user
-- -----------------------------------------------------

CREATE USER IF NOT EXISTS 'admin'@'localhost' IDENTIFIED BY 'Admin#2017';

GRANT ALL ON services.* TO 'admin'@'localhost' ;

GRANT ALL ON TABLE services.* TO 'admin'@'localhost' ;

FLUSH PRIVILEGES;

-- -----------------------------------------------------
-- Create services app user
-- -----------------------------------------------------

CREATE USER IF NOT EXISTS 'appUser'@'localhost' IDENTIFIED BY 'Services#2018';

GRANT CREATE ON services.* TO 'appUser'@'localhost' ;

GRANT INSERT ON services.* TO 'appUser'@'localhost' ;

GRANT SELECT ON services.* TO 'appUser'@'localhost' ;

GRANT UPDATE ON services.* TO 'appUser'@'localhost' ;

GRANT DELETE ON services.* TO 'appUser'@'localhost' ;

FLUSH PRIVILEGES;
