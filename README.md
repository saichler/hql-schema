# hql-schema
Habitat Schema structure for Structs/Objects for HQL Queries

The schema is a tree structure describing the relations between structs via a key path to structs attributes inside structs. for example: If i define a struct named "Employee" with an inner struct named "Address" with a field name of "Addr", hence the schema key will be "Employee.Addr" and the type will be "Address".
