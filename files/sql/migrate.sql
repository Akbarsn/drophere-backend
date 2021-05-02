DROP TABLE IF EXISTS links ;
DROP TABLE IF EXISTS user_storage_credentials ;
DROP TABLE IF EXISTS users ;

CREATE TABLE `users` ( 
`id` int unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
`email` varchar(255) NOT NULL UNIQUE,
`name` varchar(255) NOT NULL, 
`password` varchar(80) NULL,
`dropbox_token` varchar(255) DEFAULT NULL,
`drive_token` varchar(255) DEFAULT NULL,
`recover_password_token` varchar(255) NULL,
`recover_password_token_expiry` datetime NULL) ;


CREATE TABLE `links` (
  `id` int unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `user_id` int unsigned NOT NULL,
  `title` varchar(255) CHARACTER SET utf8mb4 NOT NULL,
  `password` varchar(255) NOT NULL,
  `slug` varchar(255) NOT NULL UNIQUE,
  `description` text CHARACTER SET utf8mb4 NOT NULL,
  `deadline` datetime NULL DEFAULT NULL,
  KEY `links_user_id_users_id_foreign` (`user_id`),
  CONSTRAINT `links_user_id_users_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ;

CREATE TABLE `user_storage_credentials` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(10) unsigned NOT NULL,
  `provider_id` int(10) unsigned NOT NULL,
  `provider_credential` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL DEFAULT '',
  `photo` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_provider_unique` (`user_id`, `provider_id`),
  KEY `usc_user_id_users_id_foreign` (`user_id`),
  CONSTRAINT `usc_user_id_users_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ;

ALTER TABLE `links`
ADD `user_storage_credential_id` int(10) unsigned NULL,
ADD KEY `links_usc_id_foreign` (`user_storage_credential_id`),
ADD CONSTRAINT `links_usc_id_foreign` FOREIGN KEY (`user_storage_credential_id`) REFERENCES `user_storage_credentials` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;



