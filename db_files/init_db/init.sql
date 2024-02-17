CREATE TABLE `url_maps` (
  `id` int NOT NULL AUTO_INCREMENT,
  `url` varchar(36) NOT NULL,
  `long_url` varchar(500) NOT NULL,
  `created_at` datetime default current_timestamp,
  `visit_count` int DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `url` (`url`)
);
CREATE TABLE `url_meta` (
  `url_id` int NOT NULL,
  `ip` varchar(40) NOT NULL,
  `location` varchar(60) NOT NULL,
  `device_type` varchar(50) NOT NULL,
  FOREIGN KEY (`url_id`) REFERENCES `url_maps` (`id`)
);

CREATE UNIQUE INDEX idx_urls ON url_maps (url);