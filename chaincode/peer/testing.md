# Testing

## Test case

|No|function|pattern|test case|
|---|:--|:--|:--|
|a-1|Register|success|Test_Get1|
|a-2||number of parameters<expect|Test_Register1|
|a-3||number of parameters>expect|Test_Register1|
|a-4||duplicate key|Test_Register2|
|b-1|Get|success|Test_Get|
|b-2||not found|Test_Get|
|c-1|Update|success|Test_Update|
|c-2||number of parameters<expect|Test_Update|
|c-3||number of parameters>expect|Test_Update|
|c-4||not found|Test_Update|
|c-5||not own|Test_Update|
|c-6||same data|Test_Update|
|d-1|Deregister|success|Test_Deregister4|
|d-2||number of parameters>expect|Test_Deregister1|
|d-3||not found|Test_Deregister2|
