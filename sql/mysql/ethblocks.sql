/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `block_uncles` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `block_number` bigint(20) unsigned DEFAULT NULL,
  `block_time` bigint(20) unsigned DEFAULT NULL,
  `parent_hash` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `uncle_hash` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `block_root` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `tx_hash` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `receipt_hash` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `mix_digest` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `block_nonce` bigint(20) unsigned DEFAULT NULL,
  `coinbase` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `gas_limit` bigint(20) unsigned DEFAULT NULL,
  `gas_used` bigint(20) unsigned DEFAULT NULL,
  `difficulty` bigint(20) unsigned DEFAULT NULL,
  `block_size` double DEFAULT NULL,
  `block_id` int(10) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_block_uncles_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `blocks` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `block_number` bigint(20) unsigned DEFAULT NULL,
  `block_time` bigint(20) unsigned DEFAULT NULL,
  `parent_hash` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `uncle_hash` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `block_root` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `tx_hash` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `receipt_hash` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `mix_digest` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `block_nonce` bigint(20) unsigned DEFAULT NULL,
  `coinbase` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `gas_limit` bigint(20) unsigned DEFAULT NULL,
  `gas_used` bigint(20) unsigned DEFAULT NULL,
  `difficulty` bigint(20) unsigned DEFAULT NULL,
  `block_size` double DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_blocks_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `transaction_log_topics` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `topic` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `block_id` int(10) unsigned DEFAULT NULL,
  `transaction_id` int(10) unsigned DEFAULT NULL,
  `transaction_receipt_id` int(10) unsigned DEFAULT NULL,
  `transaction_log_id` int(10) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_transaction_log_topics_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `transaction_logs` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `block_number` bigint(20) unsigned DEFAULT NULL,
  `block_hash` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `address` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `log_data` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `tx_hash` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `tx_index` int(10) unsigned DEFAULT NULL,
  `log_index` int(10) unsigned DEFAULT NULL,
  `removed` tinyint(1) DEFAULT NULL,
  `block_id` int(10) unsigned DEFAULT NULL,
  `transaction_id` int(10) unsigned DEFAULT NULL,
  `transaction_receipt_id` int(10) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_transaction_logs_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `transaction_receipts` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `block_number` bigint(20) unsigned DEFAULT NULL,
  `block_hash` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `tx_hash` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `tx_index` int(10) unsigned DEFAULT NULL,
  `tx_status` bigint(20) unsigned DEFAULT NULL,
  `cumulative_gas_used` bigint(20) unsigned DEFAULT NULL,
  `gas_used` bigint(20) unsigned DEFAULT NULL,
  `contract_address` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `post_state` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `block_id` int(10) unsigned DEFAULT NULL,
  `transaction_id` int(10) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_transaction_receipts_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `transactions` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `block_number` bigint(20) unsigned DEFAULT NULL,
  `block_hash` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `account_nonce` bigint(20) unsigned DEFAULT NULL,
  `price` bigint(20) unsigned DEFAULT NULL,
  `gas_limit` bigint(20) unsigned DEFAULT NULL,
  `tx_amount` bigint(20) unsigned DEFAULT NULL,
  `payload` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `tx_v` bigint(20) unsigned DEFAULT NULL,
  `tx_r` bigint(20) unsigned DEFAULT NULL,
  `tx_s` bigint(20) unsigned DEFAULT NULL,
  `block_id` int(10) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_transactions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
