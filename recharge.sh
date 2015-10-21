#!/bin/bash
host=nporsche.com
name=ty
password=ty789

mysql -h$host -u$name -p$password << EOF
use football;
set @player_id=(select id from players where name='$1');
INSERT INTO revenue_log(player_id, amount, reason) VALUES(@player_id,$2,'$3');
EOF


