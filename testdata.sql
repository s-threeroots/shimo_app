insert into estimations
(
	created_at,
	updated_at,
	deleted_at,
	client_name,
	estimation_name,
	sub_total,
	tax,
	total
)
values (
	current_timestamp,
	current_timestamp,
	null,
	'client',
	'estimation',
	3000,
	300,
	3300
)

select * from estimations


insert into groups 
(
	created_at,
	updated_at,
	deleted_at,
	estimation_id,
	name
)
values (
	current_timestamp,
	current_timestamp,
	null,
	1,
	'group'
)

insert into items
(
	created_at,
	updated_at,
	deleted_at,
	group_id,
	name,
	amount,
	unit,
	unit_price,
	price
)
values (
	current_timestamp,
	current_timestamp,
	null,
	1,
	'name1',
	10,
	'unit1',
	200,
	2000
)

insert into items
(
	created_at,
	updated_at,
	deleted_at,
	group_id,
	name,
	amount,
	unit,
	unit_price,
	price
)
values (
	current_timestamp,
	current_timestamp,
	null,
	1,
	'name2',
	10,
	'unit2',
	100,
	1000
)
