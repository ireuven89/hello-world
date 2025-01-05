alter table `users`
add column `players_id`;
alter table `users` fk_item_to_user
foreign key (item_to_user)
references item_to_user(id);
