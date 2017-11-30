# Testing

## Test case

|No|function|pattern|test case|
|---|:--|:--|:--|
|a-1|Register|success|Test_Register3|
|a-2||number of parameters<expect|Test_Register1|
|a-3||number of parameters>expect|Test_Register1|
|a-4||duplicate key|Test_Register2|
|b-1|Get|success|Test_Get1|
|b-2||not found||
|c-1|Update|success||
|c-2||number of parameters<expect||
|c-3||number of parameters>expect||
|c-4||not found||
|c-5||not own||
|d-1|Deregister|success|Test_Deregister4|
|d-2||number of parameters<expect|Test_Deregister1|
|d-3||number of parameters>expect|Test_Deregister2|
|d-4||not found||
|d-5||not own|Test_Deregister3|
