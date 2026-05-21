CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE document_requests (
                                   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                   student_id UUID NOT NULL,
                                   type VARCHAR(100) NOT NULL,
                                   status VARCHAR(50) DEFAULT 'pending'
);

CREATE TABLE finance_accounts (
                                  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                  user_id UUID NOT NULL UNIQUE,
                                  balance FLOAT DEFAULT 0.0
);

CREATE TABLE transactions (
                              id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                              account_id UUID REFERENCES finance_accounts(id) ON DELETE CASCADE,
                              amount FLOAT NOT NULL,
                              type VARCHAR(50) NOT NULL,
                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE dormitories (
                             id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                             room_number VARCHAR(20) NOT NULL,
                             occupant_id UUID UNIQUE
);