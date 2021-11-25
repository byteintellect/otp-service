SET
SQL_MODE= "NO_AUTO_VALUE_ON_ZERO";

CREATE TABLE otps
(
    id          BIGINT      NOT NULL AUTO_INCREMENT,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP,
    deleted_at  TIMESTAMP,
    phone       VARCHAR(45) NOT NULL,
    token       VARCHAR(45) NOT NULL,
    external_id VARCHAR(45) NOT NULL,
    status      TINYINT   DEFAULT 0,
    PRIMARY KEY (id),
    UNIQUE (external_id)
);