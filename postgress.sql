-- 1. COluna de assinatura 
ALTER TABLE certificates 
ADD COLUMN signature TEXT;

-- 2. Sem duplos hashs
ALTER TABLE certificates 
ADD CONSTRAINT unique_file_hash UNIQUE (file_hash);

-- 3. Cria um Índice para buscas por hash
CREATE INDEX idx_certificates_file_hash 
ON certificates (file_hash);

CREATE TABLE IF NOT EXISTS blocks (
    index INTEGER PRIMARY KEY,
    timestamp BIGINT NOT NULL,
    prev_hash TEXT NOT NULL,
    hash TEXT NOT NULL,
    nonce INTEGER NOT NULL,
    transactions JSONB NOT NULL -- Aqui guardamos a lista de certificados do bloco
);