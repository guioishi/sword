CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(191) DEFAULT NULL,
  `password` longtext,
  `role` longtext,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_users_username` (`username`)
);

INSERT IGNORE INTO users (username, `password`, `role`) VALUES ('manager', '$2a$10$6CZ.RhfOxgzT1ogJM0IZ.u52jy8RHFypB08r9b1TQSl.PC0HZRr7a', 'manager');
INSERT IGNORE INTO users (username, `password`, `role`) VALUES ('technician', '$2a$10$osXJ238Kj4i4kuai.mMUheUhZVZG.7NkQIsVGlltENOvBRB1B7kdq', 'technician');
