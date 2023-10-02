CREATE TABLE IF NOT EXISTS "operation_types" (
    "operation_type_id" INT PRIMARY KEY,
    "description" TEXT NOT NULL
);

INSERT INTO operation_types (operation_type_id, description) VALUES (1, 'Normal Purchase') ON CONFLICT (operation_type_id) DO NOTHING;
INSERT INTO operation_types (operation_type_id, description) VALUES (2, 'Purchase with installments') ON CONFLICT (operation_type_id) DO NOTHING;
INSERT INTO operation_types (operation_type_id, description) VALUES (3, 'Withdrawal') ON CONFLICT (operation_type_id) DO NOTHING;
INSERT INTO operation_types (operation_type_id, description) VALUES (4, 'Credit Voucher') ON CONFLICT (operation_type_id) DO NOTHING;
