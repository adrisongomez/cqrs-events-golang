# CQRS

*CQRS* means Command Query Responibility Segregation

RS -> Every service should have a unique Responibility.

"There are only way to access or write that kind of data and it's delegate to the X service"

C -> Write on DB or Datasource
Q -> Read from DB or Datasource

This separation of concern allow us to grow our services depending of the current need from write/read of a service 
without affecting the whole structure of the system.



## CQRS on Event driven Architecture...


        [Commands]
            A
            |
|||||||| BUFFER |||||||||||
            |
            V
        [Queries]
