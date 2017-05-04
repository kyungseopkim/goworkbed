convert2https

this is a tool to convert urls started with "http://" in account setting and published urls.

==CAUTION:==
convert2https and database.yml files should be in the same folder.
convert2https is looking database connection info at database.yml 

Syntax:
$ ./convert2https  <account> <environment> [confirm]

example:

$ ./convert2https captora.com production 
it will print SQL statements without modification. if manual modification is needed, 
redirect statements into file and review it.

$ ./convert2https captora.com staging  > https.sql 
$ vi https.sql  and modify it
$ mysql -h ... -u ... -p .. < https.sql

if you are sure that everything is correct.

$ ./convert2https captora.com production y 
it will modify all statements directly into Database.




