#!/bin/bash
host=localhost
name=hangchen
password=hangchen

mysql -h$host -u$name -p$password << EOF
use football;
INSERT INTO players(id, name, cellphone, status) VALUES($1,'$2',$3,$4) ON DUPLICATE KEY UPDATE name='$2', cellphone=$3, status=$4;
EOF


