-- migrate:up
INSERT INTO operation_types (description, entry_type) VALUES
('Normal Purchase','debit'),
('Purchase with installments','debit'),
('Withdrawal','debit'),
('Credit Voucher','credit');


-- migrate:down

DELETE FROM operation_types
WHERE description IN ('Normal Purchase', 'Purchase with installments', 'Withdrawal', 'Credit Voucher');

